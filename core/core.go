// +build !js

// package core contains everything needed to configure and run a 0x Mesh node.
package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/blockwatch"
	"github.com/0xProject/0x-mesh/keys"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch"
	"github.com/albrow/stringset"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	p2pcrypto "github.com/libp2p/go-libp2p-crypto"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
)

const (
	blockWatcherRetentionLimit = 20
	ethereumRPCRequestTimeout  = 30 * time.Second
	ethWatcherPollingInterval  = 1 * time.Minute
	peerConnectTimeout         = 60 * time.Second
	checkNewAddrInterval       = 20 * time.Second
)

// Config is a set of configuration options for 0x Mesh.
type Config struct {
	// Verbosity is the logging verbosity: 0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug 6=trace
	Verbosity int `envvar:"VERBOSITY" default:"2"`
	// DataDir is the directory to use for persisting all data, including the
	// database and private key files.
	DataDir string `envvar:"DATA_DIR" default:"0x_mesh"`
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
	// OrderExpirationBuffer is the amount of time before the order's stipulated expiration time
	// that you'd want it pruned from the Mesh node.
	OrderExpirationBuffer time.Duration `envvar:"ORDER_EXPIRATION_BUFFER" default:"10s"`
	// BlockPollingInterval is the polling interval to wait before checking for a new Ethereum block
	// that might contain transactions that impact the fillability of orders stored by Mesh. Different
	// networks have different block producing intervals: POW networks are typically slower (e.g., Mainnet)
	// and POA networks faster (e.g., Kovan) so one should adjust the polling interval accordingly.
	BlockPollingInterval time.Duration `envvar:"BLOCK_POLLING_INTERVAL" default:"5s"`
	// EthereumRPCMaxContentLength is the maximum request Content-Length accepted by the backing Ethereum RPC
	// endpoint used by Mesh. Geth & Infura both limit a request's content length to 1024 * 512 Bytes. Parity
	// and Alchemy have much higher limits. When batch validating 0x orders, we will fit as many orders into a
	// request without crossing the max content length. The default value is appropriate for operators using Geth
	// or Infura. If using Alchemy or Parity, feel free to double the default max in order to reduce the
	// number of RPC calls made by Mesh.
	EthereumRPCMaxContentLength int `envvar:"ETHEREUM_RPC_MAX_CONTENT_LENGTH" default:"524288"`
}

type App struct {
	config         Config
	db             *meshdb.MeshDB
	node           *p2p.Node
	networkID      int
	blockWatcher   *blockwatch.Watcher
	orderWatcher   *orderwatch.Watcher
	ethWatcher     *ethereum.ETHWatcher
	orderValidator *zeroex.OrderValidator
	orderJSONSchema     *gojsonschema.Schema
}

func New(config Config) (*App, error) {
	// Configure logger
	// TODO(albrow): Don't use global variables for log settings.
	log.SetLevel(log.Level(config.Verbosity))
	log.WithField("config", config).Info("creating new App with config")

	if config.EthereumRPCMaxContentLength < maxOrderSizeInBytes {
		return nil, fmt.Errorf("Cannot set `EthereumRPCMaxContentLength` to be less then maxOrderSizeInBytes: %d", maxOrderSizeInBytes)
	}

	// Initialize db
	databasePath := filepath.Join(config.DataDir, "db")
	db, err := meshdb.NewMeshDB(databasePath)
	if err != nil {
		return nil, err
	}

	// Initialize the ETH client, which will be used by various watchers.
	ethClient, err := ethclient.Dial(config.EthereumRPCURL)
	if err != nil {
		return nil, err
	}

	// Initialize the order validator
	orderValidator, err := zeroex.NewOrderValidator(ethClient, config.EthereumNetworkID, config.EthereumRPCMaxContentLength)
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
		MeshDB:              db,
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
	orderWatcher, err := orderwatch.New(db, blockWatcher, orderValidator, config.EthereumNetworkID, config.OrderExpirationBuffer)
	if err != nil {
		return nil, err
	}

	// Initialize the ETH balance watcher (but don't start it yet).
	ethWatcher, err := ethereum.NewETHWatcher(ethWatcherPollingInterval, ethClient, config.EthereumNetworkID)
	if err != nil {
		return nil, err
	}
	// TODO(albrow): Call Add for all existing makers/signers in the database.

	orderJSONSchema, err := setupOrderSchemaValidator()
	if err != nil {
		return nil, err
	}
	app := &App{
		config:         config,
		db:             db,
		networkID:      config.EthereumNetworkID,
		blockWatcher:   blockWatcher,
		orderWatcher:   orderWatcher,
		ethWatcher:     ethWatcher,
		orderValidator: orderValidator,
		orderJSONSchema:     orderJSONSchema,
	}

	// Initialize the p2p node.
	privateKeyPath := filepath.Join(config.DataDir, "keys", "privkey")
	privKey, err := initPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}
	nodeConfig := p2p.Config{
		Topic:            getPubSubTopic(config.EthereumNetworkID),
		ListenPort:       config.P2PListenPort,
		Insecure:         false,
		PrivateKey:       privKey,
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

func initPrivateKey(path string) (p2pcrypto.PrivKey, error) {
	privKey, err := keys.GetPrivateKeyFromPath(path)
	if err == nil {
		return privKey, nil
	} else if os.IsNotExist(err) {
		// If the private key doesn't exist, generate one.
		log.Info("No private key found. Generating a new one.")
		return keys.GenerateAndSavePrivateKey(path)
	}

	// For any other type of error, return it.
	return nil, err
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

// AddOrders can be used to add orders to Mesh. It validates the given orders
// and if they are valid, will store and eventually broadcast the orders to peers.
func (app *App) AddOrders(signedOrdersRaw []*json.RawMessage) (*zeroex.ValidationResults, error) {
	allValidationResults := &zeroex.ValidationResults{
		Accepted: []*zeroex.AcceptedOrderInfo{},
		Rejected: []*zeroex.RejectedOrderInfo{},
	}
	schemaValidOrders := []*zeroex.SignedOrder{}
	for _, signedOrderRaw := range signedOrdersRaw {
		signedOrderBytes := []byte(*signedOrderRaw)
		result, err := app.schemaValidateOrder(signedOrderBytes)
		if err != nil {
			signedOrder := &zeroex.SignedOrder{}
			if err := signedOrder.UnmarshalJSON(signedOrderBytes); err != nil {
				signedOrder = nil
			}
			log.WithField("signedOrderRaw", string(signedOrderBytes)).Info("Unexpected error while attempting to validate signedOrderJSON against schema")
			allValidationResults.Rejected = append(allValidationResults.Rejected, &zeroex.RejectedOrderInfo{
				SignedOrder: signedOrder,
				Kind:        MeshValidation,
				Status: zeroex.RejectedOrderStatus{
					Code:    ROInvalidSchemaCode,
					Message: "order did not pass JSON-schema validation: Malformed JSON or empty payload",
				},
			})
			continue
		}
		if !result.Valid() {
			log.WithField("signedOrderRaw", string(signedOrderBytes)).Info("Order failed schema validation")
			status := zeroex.RejectedOrderStatus{
				Code:    ROInvalidSchemaCode,
				Message: fmt.Sprintf("order did not pass JSON-schema validation: %s", result.Errors()),
			}
			signedOrder := &zeroex.SignedOrder{}
			if err := signedOrder.UnmarshalJSON(signedOrderBytes); err != nil {
				signedOrder = nil
			}
			allValidationResults.Rejected = append(allValidationResults.Rejected, &zeroex.RejectedOrderInfo{
				SignedOrder: signedOrder,
				Kind:        MeshValidation,
				Status:      status,
			})
			continue
		}

		signedOrder := &zeroex.SignedOrder{}
		if err := signedOrder.UnmarshalJSON(signedOrderBytes); err != nil {
			// This error should never happen since the signedOrder already passed the JSON schema validation above
			log.WithField("signedOrderRaw", string(signedOrderBytes)).Panic("Failed to unmarshal SignedOrder")
		}
		schemaValidOrders = append(schemaValidOrders, signedOrder)
	}

	validationResults, err := app.validateOrders(schemaValidOrders)
	if err != nil {
		return nil, err
	}
	for _, orderInfo := range validationResults.Accepted {
		allValidationResults.Accepted = append(allValidationResults.Accepted, orderInfo)
	}
	for _, orderInfo := range validationResults.Rejected {
		allValidationResults.Rejected = append(allValidationResults.Rejected, orderInfo)
	}

	for _, acceptedOrderInfo := range allValidationResults.Accepted {
		err = app.orderWatcher.Watch(acceptedOrderInfo)
		if err != nil {
			return nil, err
		}
	}
	return allValidationResults, nil
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
