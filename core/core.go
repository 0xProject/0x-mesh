// +build !js

// package core contains everything needed to configure and run a 0x Mesh node.
package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/blockwatch"
	"github.com/0xProject/0x-mesh/expirationwatch"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/rpc"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch"
	"github.com/albrow/stringset"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

const (
	blockWatcherRetentionLimit = 20
	ethereumRPCRequestTimeout  = 30 * time.Second
	ethWatcherPollingInterval  = 1 * time.Minute
	peerConnectTimeout         = 60 * time.Second
	checkNewAddrInterval       = 20 * time.Second
	expirationPollingInterval  = 50 * time.Millisecond
)

// Config is a set of configuration options for 0x Mesh.
type Config struct {
	// Verbosity is the logging verbosity: 0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug 6=trace
	Verbosity int `envvar:"VERBOSITY" default:"2"`
	// DatabaseDir is the directory to use for persisting the database.
	DatabaseDir string `envvar:"DATABASE_DIR" default:"./0x_mesh/db"`
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
	// PrivateKeyPath is the path to a Secp256k1 private key which will be
	// used for signing messages and generating a peer ID. If empty, a randomly
	// generated key will be used.
	PrivateKeyPath string `envvar:"PRIVATE_KEY_PATH" default:"./0x_mesh/key/privkey"`
	// OrderExpirationBuffer is the amount of time before the order's stipulated expiration time
	// that you'd want it pruned from the Mesh node.
	OrderExpirationBuffer time.Duration `envvar:"ORDER_EXPIRATION_BUFFER" default:"10s"`
	// BlockPollingInterval is the polling interval to wait before checking for a new Ethereum block
	// that might contain transactions that impact the fillability of orders stored by Mesh. Different
	// networks have different block producing intervals: POW networks are typically slower (e.g., Mainnet)
	// and POA networks faster (e.g., Kovan) so one should adjust the polling interval accordingly.
	BlockPollingInterval time.Duration `envvar:"BLOCK_POLLING_INTERVAL" default:"5s"`
}

type snapshotInfo struct {
	Snapshot            *db.Snapshot
	ExpirationTimestamp time.Time
}

type App struct {
	config                    Config
	db                        *meshdb.MeshDB
	node                      *p2p.Node
	networkID                 int
	blockWatcher              *blockwatch.Watcher
	orderWatcher              *orderwatch.Watcher
	ethWatcher                *ethereum.ETHWatcher
	orderValidator            *zeroex.OrderValidator
	snapshotExpirationWatcher *expirationwatch.Watcher
	muIdToSnapshotInfo        sync.Mutex
	idToSnapshotInfo          map[string]snapshotInfo
}

func New(config Config) (*App, error) {
	// Configure logger
	// TODO(albrow): Don't use global variables for log settings.
	log.SetLevel(log.Level(config.Verbosity))
	log.WithField("config", config).Info("creating new App with config")

	// Initialize db
	meshDB, err := meshdb.NewMeshDB(config.DatabaseDir)
	if err != nil {
		return nil, err
	}

	// Initialize the ETH client, which will be used by various watchers.
	ethClient, err := ethclient.Dial(config.EthereumRPCURL)
	if err != nil {
		return nil, err
	}

	// Initialize block watcher (but don't start it yet).
	blockWatcherClient, err := blockwatch.NewRpcClient(config.EthereumRPCURL, ethereumRPCRequestTimeout)
	if err != nil {
		return nil, err
	}
	topics := orderwatch.GetRelevantTopics()
	blockWatcherConfig := blockwatch.Config{
		MeshDB:              meshDB,
		PollingInterval:     config.BlockPollingInterval,
		StartBlockDepth:     ethrpc.LatestBlockNumber,
		BlockRetentionLimit: blockWatcherRetentionLimit,
		WithLogs:            true,
		Topics:              topics,
		Client:              blockWatcherClient,
	}
	blockWatcher := blockwatch.New(blockWatcherConfig)
	go func() {
		for {
			err, isOpen := <-blockWatcher.Errors
			if isOpen {
				log.WithField("error", err).Error("BlockWatcher error encountered")
			} else {
				return // Exit when the error channel is closed
			}
		}
	}()

	// Initialize order watcher (but don't start it yet).
	orderWatcher, err := orderwatch.New(meshDB, blockWatcher, ethClient, config.EthereumNetworkID, config.OrderExpirationBuffer)
	if err != nil {
		return nil, err
	}

	// Initialize the ETH balance watcher (but don't start it yet).
	ethWatcher, err := ethereum.NewETHWatcher(ethWatcherPollingInterval, ethClient, config.EthereumNetworkID)
	if err != nil {
		return nil, err
	}
	// TODO(albrow): Call Add for all existing makers/signers in the database.

	snapshotExpirationWatcher := expirationwatch.New(0 * time.Second)

	// Initialize the order validator
	orderValidator, err := zeroex.NewOrderValidator(ethClient, config.EthereumNetworkID)
	if err != nil {
		return nil, err
	}

	app := &App{
		config:                    config,
		db:                        meshDB,
		networkID:                 config.EthereumNetworkID,
		blockWatcher:              blockWatcher,
		orderWatcher:              orderWatcher,
		ethWatcher:                ethWatcher,
		orderValidator:            orderValidator,
		snapshotExpirationWatcher: snapshotExpirationWatcher,
		idToSnapshotInfo:          map[string]snapshotInfo{},
	}

	// Initialize the p2p node.
	nodeConfig := p2p.Config{
		Topic:            getPubSubTopic(config.EthereumNetworkID),
		ListenPort:       config.P2PListenPort,
		Insecure:         false,
		PrivateKeyPath:   config.PrivateKeyPath,
		MessageHandler:   app,
		RendezvousString: getRendezvous(config.EthereumNetworkID),
		UseBootstrapList: config.UseBootstrapList,
	}
	node, err := p2p.New(nodeConfig)
	if err != nil {
		return nil, err
	}
	app.node = node

	return app, nil
}

func getPubSubTopic(networkID int) string {
	return fmt.Sprintf("/0x-orders/network/%d/version/0.0.1", networkID)
}

func getRendezvous(networkID int) string {
	return fmt.Sprintf("/0x-mesh/network/%d/version/0.0.1", networkID)
}

func (app *App) Start() error {
	go func() {
		err := app.node.Start()
		if err != nil {
			log.WithField("error", err.Error()).Error("p2p node returned error")
			app.Close()
		}
	}()
	addrs := app.node.Multiaddrs()
	go app.periodicallyCheckForNewAddrs(addrs)
	log.WithFields(map[string]interface{}{
		"addresses": addrs,
		"peerID":    app.node.ID().String(),
	}).Info("started p2p node")

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
	if err := app.ethWatcher.Start(); err != nil {
		return err
	}
	log.Info("started ETH balance watcher")
	go func() {
		expiredSnapshotsChan := app.snapshotExpirationWatcher.Receive()
		for expiredSnapshots := range expiredSnapshotsChan {
			for _, expiredSnapshot := range expiredSnapshots {
				app.muIdToSnapshotInfo.Lock()
				delete(app.idToSnapshotInfo, expiredSnapshot.ID)
				app.muIdToSnapshotInfo.Unlock()
			}
		}
	}()
	if err := app.snapshotExpirationWatcher.Start(expirationPollingInterval); err != nil {
		return err
	}
	log.Info("started snapshot expiration watcher")

	return nil
}

func (app *App) periodicallyCheckForNewAddrs(startingAddrs []ma.Multiaddr) {
	seenAddrs := stringset.New()
	for _, addr := range startingAddrs {
		seenAddrs.Add(addr.String())
	}
	// TODO: There might be a more efficient way to do this if we have access to
	// an event bus. See: https://github.com/libp2p/go-libp2p/issues/467
	for {
		time.Sleep(checkNewAddrInterval)
		newAddrs := app.node.Multiaddrs()
		for _, addr := range newAddrs {
			if !seenAddrs.Contains(addr.String()) {
				log.WithFields(map[string]interface{}{
					"address": addr,
				}).Info("found new listen address")
				seenAddrs.Add(addr.String())
			}
		}
	}
}

// ErrSnapshotNotFound is the error returned when a snapshot not found with a particular id
type ErrSnapshotNotFound struct {
	id string
}

func (e ErrSnapshotNotFound) Error() string {
	return fmt.Sprintf("No snapshot found with id: %s. To create a new snapshot, send a request with an empty snapshotID", e.id)
}

// GetOrders retrieves paginated orders from the Mesh DB at a specific snapshot in time. Passing an empty
// string as `snapshotID` creates a new snapshot and returns the first set of results. To fetch all orders,
// continue to make requests supplying the `snapshotID` returned from the first request. After 1 minute of not
// receiving further requests referencing a specific snapshot, the snapshot expires and can no longer be used.
func (app *App) GetOrders(page, perPage int, snapshotID string) (*rpc.GetOrdersResponse, error) {
	ordersInfos := []*zeroex.AcceptedOrderInfo{}
	if perPage <= 0 {
		return &rpc.GetOrdersResponse{
			OrdersInfos: ordersInfos,
			SnapshotID:  snapshotID,
		}, nil
	}

	var snapshot *db.Snapshot
	if snapshotID == "" {
		// Create a new snapshot
		snapshotID = uuid.New().String()
		var err error
		snapshot, err = app.db.Orders.GetSnapshot()
		if err != nil {
			return nil, err
		}
		expirationTimestamp := time.Now().Add(1 * time.Minute)
		app.snapshotExpirationWatcher.Add(expirationTimestamp, snapshotID)
		app.muIdToSnapshotInfo.Lock()
		app.idToSnapshotInfo[snapshotID] = snapshotInfo{
			Snapshot:            snapshot,
			ExpirationTimestamp: expirationTimestamp,
		}
		app.muIdToSnapshotInfo.Unlock()
	} else {
		// Try and find an existing snapshot
		app.muIdToSnapshotInfo.Lock()
		info, ok := app.idToSnapshotInfo[snapshotID]
		if !ok {
			app.muIdToSnapshotInfo.Unlock()
			return nil, ErrSnapshotNotFound{id: snapshotID}
		}
		snapshot = info.Snapshot
		// Reset the snapshot's expiry
		app.snapshotExpirationWatcher.Remove(info.ExpirationTimestamp, snapshotID)
		expirationTimestamp := time.Now().Add(1 * time.Minute)
		app.snapshotExpirationWatcher.Add(expirationTimestamp, snapshotID)
		app.idToSnapshotInfo[snapshotID] = snapshotInfo{
			Snapshot:            snapshot,
			ExpirationTimestamp: expirationTimestamp,
		}
		app.muIdToSnapshotInfo.Unlock()
	}

	notRemovedFilter := app.db.Orders.IsRemovedIndex.ValueFilter([]byte{0})
	var selectedOrders []*meshdb.Order
	err := snapshot.NewQuery(notRemovedFilter).Offset(page * perPage).Max(perPage).Run(&selectedOrders)
	if err != nil {
		return nil, err
	}
	for _, order := range selectedOrders {
		ordersInfos = append(ordersInfos, &zeroex.AcceptedOrderInfo{
			OrderHash:                order.Hash,
			SignedOrder:              order.SignedOrder,
			FillableTakerAssetAmount: order.FillableTakerAssetAmount,
		})
	}

	getOrdersResponse := &rpc.GetOrdersResponse{
		SnapshotID:  snapshotID,
		OrdersInfos: ordersInfos,
	}

	return getOrdersResponse, nil
}

// AddOrders can be used to add orders to Mesh. It validates the given orders
// and if they are valid, will store and eventually broadcast the orders to peers.
func (app *App) AddOrders(orders []*zeroex.SignedOrder) (*zeroex.ValidationResults, error) {
	validationResults, err := app.validateOrders(orders)
	if err != nil {
		return nil, err
	}
	for _, acceptedOrderInfo := range validationResults.Accepted {
		err = app.orderWatcher.Watch(acceptedOrderInfo)
		if err != nil {
			return nil, err
		}
	}
	return validationResults, nil
}

// AddPeer can be used to manually connect to a new peer.
func (app *App) AddPeer(peerInfo peerstore.PeerInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), peerConnectTimeout)
	defer cancel()
	return app.node.Connect(ctx, peerInfo)
}

// SubscribeToOrderEvents let's one subscribe to order events emitted by the OrderWatcher
func (app *App) SubscribeToOrderEvents(sink chan<- []*zeroex.OrderEvent) event.Subscription {
	subscription := app.orderWatcher.Subscribe(sink)
	return subscription
}

// Close closes the app
func (app *App) Close() {
	if err := app.node.Close(); err != nil {
		log.WithField("error", err.Error()).Error("error while closing node")
	}
	app.ethWatcher.Stop()
	if err := app.orderWatcher.Stop(); err != nil {
		log.WithField("error", err.Error()).Error("error while closing orderWatcher")
	}
	app.blockWatcher.Stop()
	app.db.Close()
}
