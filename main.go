// +build !js

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/0xProject/0x-mesh/blockwatch"
	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/meshdb"
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
	blockWatcherRetentionLimit  = 40
	ethereumRPCRequestTimeout   = 30 * time.Second
	ethWatcherPollingInterval   = 5 * time.Second
)

type meshEnvVars struct {
	// Verbosity is the logging verbosity: 0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug 6=trace
	Verbosity int `envvar:"VERBOSITY" default:"2"`
	// DatabaseDir is the directory to use for persisting the database.
	DatabaseDir string `envvar:"DATABASE_DIR" default:"./0x_mesh_db"`
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
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	// Main entry point for the 0x Mesh node
	app, err := newApp()
	if err != nil {
		log.WithField("err", err.Error()).Fatal("could not initialize app")
	}
	if err := app.start(); err != nil {
		log.WithField("err", err.Error()).Fatal("fatal error")
	}
}

func newApp() (*application, error) {
	// Parse environment variables.
	env := meshEnvVars{}
	if err := envvar.Parse(&env); err != nil {
		return nil, err
	}

	// Configure logger
	// TOOD(albrow): Don't use global veriables for these settings.
	log.SetLevel(log.Level(env.Verbosity))

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
	blockWatcherConfig := blockwatch.Config{
		MeshDB:              db,
		PollingInterval:     blockWatcherPollingInterval,
		StartBlockDepth:     rpc.LatestBlockNumber,
		BlockRetentionLimit: blockWatcherRetentionLimit,
		WithLogs:            true,
		// TODO(albrow): What should Topics be?
		Topics: nil,
		Client: blockWatcherClient,
	}
	blockWatcher := blockwatch.New(blockWatcherConfig)

	// Initialize order watcher (but don't start it yet).
	orderValidatorAddress, err := getOrderValidatorAddressForNetwork(env.EthereumNetworkID)
	if err != nil {
		return nil, err
	}
	orderWatcher, err := orderwatch.New(blockWatcher, ethClient, orderValidatorAddress)
	if err != nil {
		return nil, err
	}

	// Initialize the ETH balance watcher (but don't start it yet).
	ethBalanceCheckerAddress, err := getETHBalanceCheckerAddressForNetwork(env.EthereumNetworkID)
	if err != nil {
		return nil, err
	}
	ethWatcher, err := ethereum.NewETHWatcher(ethWatcherPollingInterval, ethClient, ethBalanceCheckerAddress)
	if err != nil {
		return nil, err
	}

	// Initialize the order validator
	orderValidator, err := zeroex.NewOrderValidator(orderValidatorAddress, ethClient)
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

	return app, nil
}

func getETHBalanceCheckerAddressForNetwork(networkID int) (common.Address, error) {
	switch networkID {
	case 1:
		return ethereum.MainnetEthBalanceCheckerAddress, nil
	case 50:
		return ethereum.GanacheEthBalanceCheckerAddress, nil
	default:
		return [common.AddressLength]byte{}, fmt.Errorf("unknown or unsupported network id: %d", networkID)
	}
}

func getOrderValidatorAddressForNetwork(networkID int) (common.Address, error) {
	switch networkID {
	case 1:
		return zeroex.MainnetOrderValidatorAddress, nil
	case 50:
		return zeroex.GanacheOrderValidatorAddress, nil
	default:
		return [common.AddressLength]byte{}, fmt.Errorf("unknown or unsupported network id: %d", networkID)
	}
}

func (app *application) GetMessagesToShare(max int) ([][]byte, error) {
	// For now, we just select a random set of orders from those we have stored.
	// TODO(albrow): This could be made more efficient if the db package supported
	// a `Count` method for counting the number of models in a collection or
	// counting the number of models that satisfy some query.
	// TODO: This will need to change when we add support for WeijieSub.
	var allOrders []*meshdb.Order
	if err := app.db.Orders.FindAll(&allOrders); err != nil {
		return nil, err
	}
	start := rand.Intn(len(allOrders))
	end := start + max
	if end > len(allOrders) {
		end = len(allOrders)
	}
	selectedOrders := allOrders[start:end]

	// After we have selected all the orders to share, we need to encode them to
	// the message data format.
	messageData := make([][]byte, len(selectedOrders))
	for i, order := range selectedOrders {
		encoded, err := encodeOrder(order.SignedOrder)
		if err != nil {
			return nil, err
		}
		messageData[i] = encoded
	}
	return messageData, nil
}

func (app *application) Validate(messages []*core.Message) ([]*core.Message, error) {
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
		if orderInfo.OrderStatus == zeroex.Fillable {
			validMessages = append(validMessages, msg)
		}
	}
	return validMessages, nil
}

func (app *application) Store(messages []*core.Message) error {
	for _, msg := range messages {
		order, err := decodeOrder(msg.Data)
		if err != nil {
			return err
		}
		orderHash, err := order.ComputeOrderHash()
		if err != nil {
			return err
		}
		// Calling Watch implicitly stores the order.
		if err := app.orderWatcher.Watch(order, orderHash); err != nil {
			return err
		}
	}
	return nil
}

func (app *application) start() error {
	// TODO(albrow): Implement this.
	return nil
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
