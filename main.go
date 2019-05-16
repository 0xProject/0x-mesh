// +build !js

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/0xProject/0x-mesh/blockwatch"
	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/ws"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/plaid/go-envvar/envvar"
	log "github.com/sirupsen/logrus"
)

const (
	pubsubTopic                 = "0x-orders:v0"
	blockWatcherPollingInterval = 5 * time.Second
	blockWatcherRetentionLimit  = 20
	ethereumRPCRequestTimeout   = 30 * time.Second
	ethWatcherPollingInterval   = 5 * time.Second
)

var (
	errInternal = errors.New("internal error")
)

type meshEnvVars struct {
	// Verbosity is the logging verbosity: 0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug 6=trace
	Verbosity int `envvar:"VERBOSITY" default:"2"`
	// DatabaseDir is the directory to use for persisting the database.
	DatabaseDir string `envvar:"DATABASE_DIR" default:"./0x_mesh_db"`
	// RPCPort is the port to use for the JSON RPC API over WebSockets. By
	// default, 0x Mesh will let the OS select a randomly available port.
	RPCPort int `envvar:"RPC_PORT" default:"0"`
	// P2PListenPort is the port on which to listen for new peer connections. By
	// default, 0x Mesh will let the OS select a randomly available port.
	P2PListenPort int `envvar:"P2P_LISTEN_PORT" default:"0"`
	// EthereumRPCURL is the URL of an Etheruem node which supports the JSON RPC
	// API.
	EthereumRPCURL string `envvar:"ETHEREUM_RPC_URL"`
	// EthereumNetworkID is the network ID to use when communicating with
	// Ethereum.
	EthereumNetworkID int `envvar:"ETHEREUM_NETWORK_ID"`
}

type application struct {
	env            meshEnvVars
	db             *meshdb.MeshDB
	node           *core.Node
	blockWatcher   *blockwatch.Watcher
	orderWatcher   *orderwatch.Watcher
	ethWathcher    *ethereum.ETHWatcher
	orderValidator *zeroex.OrderValidator
	wsServer       *ws.Server
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	// Main entry point for the 0x Mesh node
	app, err := newApp()
	if err != nil {
		log.WithField("err", err.Error()).Fatal("could not initialize app")
	}
	if err := app.start(); err != nil {
		log.WithField("err", err.Error()).Fatal("fatal error while starting app")
	}
	defer app.close()

	// Block forever
	select {}
}

func newApp() (*application, error) {
	// Parse environment variables.
	env := meshEnvVars{}
	if err := envvar.Parse(&env); err != nil {
		return nil, err
	}

	// Configure logger
	// TODO(albrow): Don't use global variables for these settings.
	log.SetLevel(log.Level(env.Verbosity))
	log.WithFields(map[string]interface{}{
		"VERBOSITY":           env.Verbosity,
		"DATABASE_DIR":        env.DatabaseDir,
		"RPC_PORT":            env.RPCPort,
		"P2P_LISTEN_PORT":     env.P2PListenPort,
		"ETHEREUM_RPC_URL":    env.EthereumRPCURL,
		"ETHEREUM_NETWORK_ID": env.EthereumNetworkID,
	}).Info("parsed environment variables")

	// Initialize db
	db, err := meshdb.NewMeshDB(env.DatabaseDir)
	if err != nil {
		return nil, err
	}

	// Initialize the ETH client, which will be used by various watchers.
	ethClient, err := ethclient.Dial(env.EthereumRPCURL)
	if err != nil {
		return nil, err
	}

	// Initialize block watcher (but don't start it yet).
	blockWatcherClient, err := blockwatch.NewRpcClient(ethClient, ethereumRPCRequestTimeout)
	if err != nil {
		return nil, err
	}
	topics := orderwatch.GetRelevantTopics()
	blockWatcherConfig := blockwatch.Config{
		MeshDB:              db,
		PollingInterval:     blockWatcherPollingInterval,
		StartBlockDepth:     rpc.LatestBlockNumber,
		BlockRetentionLimit: blockWatcherRetentionLimit,
		WithLogs:            true,
		Topics:              topics,
		Client:              blockWatcherClient,
	}
	blockWatcher := blockwatch.New(blockWatcherConfig)

	// Initialize order watcher (but don't start it yet).
	orderWatcher, err := orderwatch.New(db, blockWatcher, ethClient, env.EthereumNetworkID)
	if err != nil {
		return nil, err
	}

	// Initialize the ETH balance watcher (but don't start it yet).
	ethWatcher, err := ethereum.NewETHWatcher(ethWatcherPollingInterval, ethClient, env.EthereumNetworkID)
	if err != nil {
		return nil, err
	}
	// TODO(albrow): Call Add for all existing makers/signers in the database.

	// Initialize the order validator
	orderValidator, err := zeroex.NewOrderValidator(ethClient, env.EthereumNetworkID)
	if err != nil {
		return nil, err
	}

	app := &application{
		env:            env,
		db:             db,
		blockWatcher:   blockWatcher,
		orderWatcher:   orderWatcher,
		ethWathcher:    ethWatcher,
		orderValidator: orderValidator,
	}

	// Initialize the core node.
	nodeConfig := core.Config{
		Topic:          pubsubTopic,
		ListenPort:     env.P2PListenPort,
		Insecure:       false,
		RandSeed:       0,
		MessageHandler: app,
	}
	node, err := core.New(nodeConfig)
	if err != nil {
		return nil, err
	}
	app.node = node

	// Initialize the JSON RPC WebSocket server (but don't start it yet).
	rpcAddr := fmt.Sprintf(":%d", env.RPCPort)
	wsServer, err := ws.NewServer(rpcAddr, app)
	if err != nil {
		return nil, err
	}
	app.wsServer = wsServer

	return app, nil
}

func (app *application) GetMessagesToShare(max int) ([][]byte, error) {
	// For now, we just select a random set of orders from those we have stored.
	// TODO(albrow): This could be made more efficient if the db package supported
	// a `Count` method for counting the number of models in a collection or
	// counting the number of models that satisfy some query.
	// TODO(albrow): Add an index for IsDeleted and don't return messages that
	// have already been deleted.
	// TODO: This will need to change when we add support for WeijieSub.
	var allOrders []*meshdb.Order
	if err := app.db.Orders.FindAll(&allOrders); err != nil {
		return nil, err
	}
	if len(allOrders) == 0 {
		return nil, nil
	}
	start := rand.Intn(len(allOrders))
	end := start + max
	if end > len(allOrders) {
		end = len(allOrders)
	}
	selectedOrders := allOrders[start:end]

	log.WithFields(map[string]interface{}{
		"maxNumberToShare":    max,
		"actualNumberToShare": len(selectedOrders),
	}).Debug("preparing to share orders with peers")

	// After we have selected all the orders to share, we need to encode them to
	// the message data format.
	messageData := make([][]byte, len(selectedOrders))
	for i, order := range selectedOrders {
		log.WithFields(map[string]interface{}{
			"order": order,
		}).Debug("selected order to share")
		encoded, err := encodeOrder(order.SignedOrder)
		if err != nil {
			return nil, err
		}
		messageData[i] = encoded
	}
	return messageData, nil
}

func (app *application) ValidateAndStore(messages []*core.Message) ([]*core.Message, error) {
	orders := []*zeroex.SignedOrder{}
	orderHashToMessage := map[common.Hash]*core.Message{}
	for _, msg := range messages {
		order, err := decodeOrder(msg.Data)
		if err != nil {
			return nil, err
		}
		orderHash, err := order.ComputeOrderHash()
		if err != nil {
			return nil, err
		}
		// Validate doesn't guarantee there are no duplicates so we keep track of
		// which orders we've already seen.
		if _, alreadySeen := orderHashToMessage[orderHash]; alreadySeen {
			continue
		}
		log.WithFields(map[string]interface{}{
			"order":     order,
			"orderHash": orderHash,
			"from":      msg.From.String(),
		}).Debug("received order from peer")
		orders = append(orders, order)
		orderHashToMessage[orderHash] = msg
	}

	// Validate the orders in a single batch.
	orderHashToOrderInfo := app.orderValidator.BatchValidate(orders)

	// Filter out the invalid messages (i.e. messages which correspond to invalid
	// orders).
	validMessages := []*core.Message{}
	for orderHash, msg := range orderHashToMessage {
		orderInfo, found := orderHashToOrderInfo[orderHash]
		if !found {
			continue
		}
		if zeroex.IsOrderValid(orderInfo) {
			validMessages = append(validMessages, msg)
			alreadyStored, err := app.orderAlreadyStored(orderInfo.OrderHash)
			if err != nil {
				return nil, err
			}
			if alreadyStored {
				log.WithFields(map[string]interface{}{
					"orderInfo": orderInfo,
					"from":      msg.From.String(),
				}).Debug("order received from peer is valid but already stored")
			} else {
				log.WithFields(map[string]interface{}{
					"orderInfo": orderInfo,
					"from":      msg.From.String(),
				}).Debug("storing valid order received from peer")
				// Watch stores the message in the database.
				if err := app.orderWatcher.Watch(orderInfo); err != nil {
					return nil, err
				}
			}
		} else {
			log.WithFields(map[string]interface{}{
				"orderInfo": orderInfo,
				"from":      msg.From.String(),
			}).Debug("not storing invalid order received from peer")
		}
	}
	return validMessages, nil
}

func (app *application) start() error {
	go func() {
		err := app.node.Start()
		if err != nil {
			log.WithField("error", err.Error()).Error("core node returned error")
			app.close()
		}
	}()
	log.WithFields(map[string]interface{}{
		"multiaddress": app.node.Multiaddrs(),
		"peerID":       app.node.ID().String(),
	}).Info("started core node")

	// TODO(albrow) we might want to match the synchronous API of core.Node which
	// returns any fatal errors. As it currently stands, if one of these watchers
	// experiences a fatal error or crashes, it is difficult for us to tear down
	// correctly.
	if err := app.blockWatcher.StartPolling(); err != nil {
		return err
	}
	log.Info("started block watcher")
	if err := app.orderWatcher.Start(); err != nil {
		return err
	}
	log.Info("started order watcher")
	// TODO(fabio): Subscribe to the ETH balance updates and update them in the DB
	// for future use by the order storing algorithm.
	if err := app.ethWathcher.Start(); err != nil {
		return err
	}
	log.Info("started ETH balance watcher")

	go func() {
		err := app.wsServer.Listen()
		if err != nil {
			log.WithField("error", err.Error()).Error("RPC server returned error")
			app.close()
		}
	}()
	// Wait for the server to start listening and select an address.
	for app.wsServer.Addr() == nil {
		time.Sleep(10 * time.Millisecond)
	}
	log.WithField("address", app.wsServer.Addr().String()).Info("started RPC server")

	return nil
}

// AddOrder is called when an RPC client sends an AddOrder request.
func (app *application) AddOrder(order *zeroex.SignedOrder) error {
	log.Info("received order via RPC")
	orderHash, err := order.ComputeOrderHash()
	if err != nil {
		log.WithField("order", order).Error("received order via RPC but could not compute order hash")
		return errInternal
	}
	orderHashToOrderInfo := app.orderValidator.BatchValidate([]*zeroex.SignedOrder{order})
	orderInfo, found := orderHashToOrderInfo[orderHash]
	if !found {
		log.WithField("order", order).Error("received order via RPC but could not validate it")
		return errInternal
	}
	if !zeroex.IsOrderValid(orderInfo) {
		log.WithField("orderInfo", orderInfo).Error("received invalid order via RPC")
		return errors.New("invalid order")
	}

	log.WithField("orderInfo", orderInfo).Debug("order received via RPC is valid")

	alreadyStored, err := app.orderAlreadyStored(orderInfo.OrderHash)
	if err != nil {
		log.WithFields(map[string]interface{}{
			"orderInfo": orderInfo,
			"error":     err.Error(),
		}).Error("received valid order via RPC but could not determine if it is already stored")
		return errInternal
	}
	if alreadyStored {
		log.WithField("orderInfo", orderInfo).Debug("received valid order via RPC but it is already stored")
		return nil
	}

	if err := app.orderWatcher.Watch(orderInfo); err != nil {
		log.WithFields(map[string]interface{}{
			"orderInfo": orderInfo,
			"error":     err.Error(),
		}).Error("received valid order via RPC but could not watch it")
		return errInternal
	}

	return nil
}

func (app *application) close() {
	if err := app.wsServer.Close(); err != nil {
		log.WithField("error", err.Error()).Error("error while closing RPC server")
	}
	if err := app.node.Close(); err != nil {
		log.WithField("error", err.Error()).Error("error while closing node")
	}
	app.ethWathcher.Stop()
	if err := app.orderWatcher.Stop(); err != nil {
		log.WithField("error", err.Error()).Error("error while closing orderWatcher")
	}
	app.blockWatcher.StopPolling()
	app.db.Close()
}

type orderMessage struct {
	// TODO(albrow): Add additional metadata for the order? Signer?
	MessageType string
	Order       *zeroex.SignedOrder
}

func encodeOrder(order *zeroex.SignedOrder) ([]byte, error) {
	return json.Marshal(orderMessage{
		MessageType: "order",
		Order:       order,
	})
}

func decodeOrder(data []byte) (*zeroex.SignedOrder, error) {
	var orderMessage orderMessage
	if err := json.Unmarshal(data, &orderMessage); err != nil {
		return nil, err
	}
	if orderMessage.MessageType != "order" {
		return nil, fmt.Errorf("unexpected message type: %q", orderMessage.MessageType)
	}
	return orderMessage.Order, nil
}

// TODO(albrow): Either use the Exists method or check for a typed error after
// updating the db package.
func (app *application) orderAlreadyStored(orderHash common.Hash) (bool, error) {
	var order meshdb.Order
	err := app.db.Orders.FindByID(orderHash.Bytes(), &order)
	if err == nil {
		return true, nil
	} else if strings.Contains(err.Error(), "not found") {
		return false, nil
	}
	return false, err
}
