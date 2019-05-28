// +build !js

// package core contains everything needed to configure and run a 0x Mesh node.
package core

import (
	"context"
	"time"

	"github.com/0xProject/0x-mesh/db"

	"github.com/0xProject/0x-mesh/blockwatch"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	log "github.com/sirupsen/logrus"
)

const (
	pubsubTopic                 = "/0x-orders/0.0.1"
	blockWatcherPollingInterval = 5 * time.Second
	blockWatcherRetentionLimit  = 20
	ethereumRPCRequestTimeout   = 30 * time.Second
	ethWatcherPollingInterval   = 5 * time.Second
	peerConnectTimeout          = 60 * time.Second
)

// Config is a set of configuration options for 0x Mesh.
type Config struct {
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
	// UseBootstrapList is whether to use the predetermined list of peers to
	// bootstrap the DHT and peer discovery.
	UseBootstrapList bool `envvar:"USE_BOOTSTRAP_LIST" default:"false"`
}

type App struct {
	config         Config
	db             *meshdb.MeshDB
	node           *p2p.Node
	blockWatcher   *blockwatch.Watcher
	orderWatcher   *orderwatch.Watcher
	ethWathcher    *ethereum.ETHWatcher
	orderValidator *zeroex.OrderValidator
}

func New(config Config) (*App, error) {
	// Configure logger
	// TODO(albrow): Don't use global variables for log settings.
	log.SetLevel(log.Level(config.Verbosity))
	log.WithField("config", config).Info("creating new App with config")

	// Initialize db
	db, err := meshdb.NewMeshDB(config.DatabaseDir)
	if err != nil {
		return nil, err
	}

	// Initialize the ETH client, which will be used by various watchers.
	ethClient, err := ethclient.Dial(config.EthereumRPCURL)
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
		StartBlockDepth:     ethrpc.LatestBlockNumber,
		BlockRetentionLimit: blockWatcherRetentionLimit,
		WithLogs:            true,
		Topics:              topics,
		Client:              blockWatcherClient,
	}
	blockWatcher := blockwatch.New(blockWatcherConfig)

	// Initialize order watcher (but don't start it yet).
	orderWatcher, err := orderwatch.New(db, blockWatcher, ethClient, config.EthereumNetworkID)
	if err != nil {
		return nil, err
	}

	// Initialize the ETH balance watcher (but don't start it yet).
	ethWatcher, err := ethereum.NewETHWatcher(ethWatcherPollingInterval, ethClient, config.EthereumNetworkID)
	if err != nil {
		return nil, err
	}
	// TODO(albrow): Call Add for all existing makers/signers in the database.

	// Initialize the order validator
	orderValidator, err := zeroex.NewOrderValidator(ethClient, config.EthereumNetworkID)
	if err != nil {
		return nil, err
	}

	app := &App{
		config:         config,
		db:             db,
		blockWatcher:   blockWatcher,
		orderWatcher:   orderWatcher,
		ethWathcher:    ethWatcher,
		orderValidator: orderValidator,
	}

	// Initialize the p2p node.
	nodeConfig := p2p.Config{
		Topic:            pubsubTopic,
		ListenPort:       config.P2PListenPort,
		Insecure:         false,
		RandSeed:         0,
		MessageHandler:   app,
		RendezvousString: "/0x-mesh/0.0.1",
		UseBootstrapList: config.UseBootstrapList,
	}
	node, err := p2p.New(nodeConfig)
	if err != nil {
		return nil, err
	}
	app.node = node

	return app, nil
}

func (app *App) Start() error {
	go func() {
		err := app.node.Start()
		if err != nil {
			log.WithField("error", err.Error()).Error("core node returned error")
			app.Close()
		}
	}()
	log.WithFields(map[string]interface{}{
		"multiaddress": app.node.Multiaddrs(),
		"peerID":       app.node.ID().String(),
	}).Info("started core node")

	// TODO(albrow) we might want to match the synchronous API of p2p.Node which
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

	return nil
}

// TODO(albrow): Uset the more efficient Exists method instead of FindByID.
func (app *App) orderAlreadyStored(orderHash common.Hash) (bool, error) {
	var order meshdb.Order
	err := app.db.Orders.FindByID(orderHash.Bytes(), &order)
	if err == nil {
		return true, nil
	}
	if _, ok := err.(db.NotFoundError); ok {
		return false, nil
	}
	return false, err
}

// AddOrders can be used to add orders to Mesh. It validates the given orders
// and if they are valid, will store and eventually broadcast the orders to peers.
func (app *App) AddOrders(orders []*zeroex.SignedOrder) (*rpc.AddOrdersResponse, error) {
	orderHashToOrderInfo := app.orderValidator.BatchValidate(orders)
	addOrderResponse := &rpc.AddOrdersResponse{}
	for _, order := range orders {
		orderHash, err := order.ComputeOrderHash()
		if err != nil {
			return nil, err
		}
		orderInfo, ok := orderHashToOrderInfo[orderHash]
		if !ok {
			// Validation network request for this order failed
			addOrderResponse.FailedToAdd = append(addOrderResponse.FailedToAdd, orderHash)
			continue
		}
		succinctOrderInfo := &zeroex.SuccinctOrderInfo{
			OrderHash:                orderInfo.OrderHash,
			FillableTakerAssetAmount: orderInfo.FillableTakerAssetAmount,
			OrderStatus:              orderInfo.OrderStatus,
		}
		if !zeroex.IsOrderValid(orderInfo) {
			addOrderResponse.Invalid = append(addOrderResponse.Invalid, succinctOrderInfo)
			continue
		}
		addOrderResponse.Added = append(addOrderResponse.Added, succinctOrderInfo)

		alreadyStored, err := app.orderAlreadyStored(orderInfo.OrderHash)
		if err != nil {
			return nil, err
		}
		if alreadyStored {
			continue
		}

		err = app.orderWatcher.Watch(orderInfo)
		if err != nil {
			return nil, err
		}
	}
	return addOrderResponse, nil
}

// AddPeer can be used to manually connect to a new peer.
func (app *App) AddPeer(peerInfo peerstore.PeerInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), peerConnectTimeout)
	defer cancel()
	return app.node.Connect(ctx, peerInfo)
}

// SubscribeToOrderEvents let's one subscribe to order events emitted by the OrderWatcher
func (app *App) SubscribeToOrderEvents(sink chan<- []*zeroex.OrderInfo) event.Subscription {
	subscription := app.orderWatcher.Subscribe(sink)
	return subscription
}

// Close closes the app
func (app *App) Close() {
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
