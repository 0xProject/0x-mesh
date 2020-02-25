// Package core contains everything needed to configure and run a 0x Mesh node.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/core/ordersync"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/encoding"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/blockwatch"
	"github.com/0xProject/0x-mesh/ethereum/ethrpcclient"
	"github.com/0xProject/0x-mesh/ethereum/ratelimit"
	"github.com/0xProject/0x-mesh/ethereum/simplestack"
	"github.com/0xProject/0x-mesh/expirationwatch"
	"github.com/0xProject/0x-mesh/keys"
	"github.com/0xProject/0x-mesh/loghooks"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/orderfilter"
	"github.com/0xProject/0x-mesh/p2p"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch"
	"github.com/albrow/stringset"
	"github.com/benbjohnson/clock"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

const (
	ethereumRPCRequestTimeout     = 30 * time.Second
	peerConnectTimeout            = 60 * time.Second
	checkNewAddrInterval          = 20 * time.Second
	expirationPollingInterval     = 50 * time.Millisecond
	rateLimiterCheckpointInterval = 1 * time.Minute
	// estimatedNonPollingEthereumRPCRequestsPer24Hrs is an estimate of the
	// minimum number of RPC requests Mesh needs to send (not including block
	// polling). It's based on real-world data from a mainnet Mesh node. This
	// estimate won't necessarily hold true as network activity grows over time or
	// for different Ethereum networks, but it should be good enough.
	estimatedNonPollingEthereumRPCRequestsPer24Hrs = 50000
	// logStatsInterval is how often to log stats for this node.
	logStatsInterval = 5 * time.Minute
	version          = "development"
	// ordersyncMinPeers is the minimum amount of peers to receive orders from
	// before considering the ordersync process finished.
	ordersyncMinPeers            = 5
	paginationSubprotocolPerPage = 500
)

// Note(albrow): The Config type is currently copied to browser/ts/index.ts. We
// need to keep both definitions in sync, so if you change one you must also
// change the other.

// Config is a set of configuration options for 0x Mesh.
type Config struct {
	// Verbosity is the logging verbosity: 0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug 6=trace
	Verbosity int `envvar:"VERBOSITY" default:"2"`
	// DataDir is the directory to use for persisting all data, including the
	// database and private key files.
	DataDir string `envvar:"DATA_DIR" default:"0x_mesh"`
	// P2PTCPPort is the port on which to listen for new TCP connections from
	// peers in the network. Set to 60558 by default.
	P2PTCPPort int `envvar:"P2P_TCP_PORT" default:"60558"`
	// P2PWebSocketsPort is the port on which to listen for new WebSockets
	// connections from peers in the network. Set to 60559 by default.
	P2PWebSocketsPort int `envvar:"P2P_WEBSOCKETS_PORT" default:"60559"`
	// EthereumRPCURL is the URL of an Etheruem node which supports the JSON RPC
	// API.
	EthereumRPCURL string `envvar:"ETHEREUM_RPC_URL" json:"-"`
	// EthereumChainID is the chain ID specifying which Ethereum chain you wish to
	// run your Mesh node for
	EthereumChainID int `envvar:"ETHEREUM_CHAIN_ID"`
	// UseBootstrapList is whether to bootstrap the DHT by connecting to a
	// specific set of peers.
	UseBootstrapList bool `envvar:"USE_BOOTSTRAP_LIST" default:"true"`
	// BootstrapList is a comma-separated list of multiaddresses to use for
	// bootstrapping the DHT (e.g.,
	// "/ip4/3.214.190.67/tcp/60558/ipfs/16Uiu2HAmGx8Z6gdq5T5AQE54GMtqDhDFhizywTy1o28NJbAMMumF").
	// If empty, the default bootstrap list will be used.
	BootstrapList string `envvar:"BOOTSTRAP_LIST" default:""`
	// BlockPollingInterval is the polling interval to wait before checking for a new Ethereum block
	// that might contain transactions that impact the fillability of orders stored by Mesh. Different
	// chains have different block producing intervals: POW chains are typically slower (e.g., Mainnet)
	// and POA chains faster (e.g., Kovan) so one should adjust the polling interval accordingly.
	BlockPollingInterval time.Duration `envvar:"BLOCK_POLLING_INTERVAL" default:"5s"`
	// EthereumRPCMaxContentLength is the maximum request Content-Length accepted by the backing Ethereum RPC
	// endpoint used by Mesh. Geth & Infura both limit a request's content length to 1024 * 512 Bytes. Parity
	// and Alchemy have much higher limits. When batch validating 0x orders, we will fit as many orders into a
	// request without crossing the max content length. The default value is appropriate for operators using Geth
	// or Infura. If using Alchemy or Parity, feel free to double the default max in order to reduce the
	// number of RPC calls made by Mesh.
	EthereumRPCMaxContentLength int `envvar:"ETHEREUM_RPC_MAX_CONTENT_LENGTH" default:"524288"`
	// EnableEthereumRPCRateLimiting determines whether or not Mesh should limit
	// the number of Ethereum RPC requests it sends. It defaults to true.
	// Disabling Ethereum RPC rate limiting can reduce latency for receiving order
	// events in some network conditions, but can also potentially lead to higher
	// costs or other rate limiting issues outside of Mesh, depending on your
	// Ethereum RPC provider. If set to false, ethereumRPCMaxRequestsPer24HrUTC
	// and ethereumRPCMaxRequestsPerSecond will have no effect.
	EnableEthereumRPCRateLimiting bool `envvar:"ENABLE_ETHEREUM_RPC_RATE_LIMITING" default:"true"`
	// EthereumRPCMaxRequestsPer24HrUTC caps the number of Ethereum JSON-RPC requests a Mesh node will make
	// per 24hr UTC time window (time window starts and ends at midnight UTC). It defaults to 200k but
	// can be increased well beyond this limit depending on your infrastructure or Ethereum RPC provider.
	EthereumRPCMaxRequestsPer24HrUTC int `envvar:"ETHEREUM_RPC_MAX_REQUESTS_PER_24_HR_UTC" default:"200000"`
	// EthereumRPCMaxRequestsPerSecond caps the number of Ethereum JSON-RPC requests a Mesh node will make per
	// second. This limits the concurrency of these requests and prevents the Mesh node from getting rate-limited.
	// It defaults to the recommended 30 rps for Infura's free tier, and can be increased to 100 rpc for pro users,
	// and potentially higher on alternative infrastructure.
	EthereumRPCMaxRequestsPerSecond float64 `envvar:"ETHEREUM_RPC_MAX_REQUESTS_PER_SECOND" default:"30"`
	// CustomContractAddresses is a JSON-encoded string representing a set of
	// custom addresses to use for the configured chain ID. The contract
	// addresses for most common chains/networks are already included by default, so this
	// is typically only needed for testing on custom chains/networks. The given
	// addresses are added to the default list of addresses for known chains/networks and
	// overriding any contract addresses for known chains/networks is not allowed. The
	// addresses for exchange, devUtils, erc20Proxy, erc721Proxy and erc1155Proxy are required
	// for each chain/network. For example:
	//
	//    {
	//        "exchange":"0x48bacb9266a570d521063ef5dd96e61686dbe788",
	//        "devUtils": "0x38ef19fdf8e8415f18c307ed71967e19aac28ba1",
	//        "erc20Proxy": "0x1dc4c1cefef38a777b15aa20260a54e584b16c48",
	//        "erc721Proxy": "0x1d7022f5b17d2f8b695918fb48fa1089c9f85401",
	//        "erc1155Proxy": "0x64517fa2b480ba3678a2a3c0cf08ef7fd4fad36f"
	//    }
	//
	CustomContractAddresses string `envvar:"CUSTOM_CONTRACT_ADDRESSES" default:""`
	// MaxOrdersInStorage is the maximum number of orders that Mesh will keep in
	// storage. As the number of orders in storage grows, Mesh will begin
	// enforcing a limit on maximum expiration time for incoming orders and remove
	// any orders with an expiration time too far in the future.
	MaxOrdersInStorage int `envvar:"MAX_ORDERS_IN_STORAGE" default:"100000"`
	// CustomOrderFilter is a stringified JSON Schema which will be used for
	// validating incoming orders. If provided, Mesh will only receive orders from
	// other peers in the network with the same filter.
	//
	// Here is an example filter which will only allow orders with a specific
	// makerAssetData:
	//
	//    {
	//        "properties": {
	//            "makerAssetData": {
	//                "const": "0xf47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"
	//            }
	//        }
	//    }
	//
	// Note that you only need to include the requirements for your specific
	// application in the filter. The default requirements for a valid order (e.g.
	// all the required fields) are automatically included. For more information
	// on JSON Schemas, see https://json-schema.org/
	CustomOrderFilter string `envvar:"CUSTOM_ORDER_FILTER" default:"{}"`
	// EthereumRPCClient is the client to use for all Ethereum RPC reuqests. It is only
	// settable in browsers and cannot be set via environment variable. If
	// provided, EthereumRPCURL will be ignored.
	EthereumRPCClient ethclient.RPCClient `envvar:"-"`
}

type snapshotInfo struct {
	Snapshot            *db.Snapshot
	CreatedAt           time.Time
	ExpirationTimestamp time.Time
}

type App struct {
	config                    Config
	peerID                    peer.ID
	privKey                   p2pcrypto.PrivKey
	node                      *p2p.Node
	chainID                   int
	blockWatcher              *blockwatch.Watcher
	orderWatcher              *orderwatch.Watcher
	orderValidator            *ordervalidator.OrderValidator
	orderFilter               *orderfilter.Filter
	snapshotExpirationWatcher *expirationwatch.Watcher
	muIdToSnapshotInfo        sync.Mutex
	idToSnapshotInfo          map[string]snapshotInfo
	ethRPCRateLimiter         ratelimit.RateLimiter
	ethRPCClient              ethrpcclient.Client
	db                        *meshdb.MeshDB
	ordersyncService          *ordersync.Service
	// FIXME(jalextowle): Maybe this should just be contractAddresses?
	chainIDToContractAddresses map[int]ethereum.ContractAddresses

	// started is closed to signal that the App has been started. Some methods
	// will block until after the App is started.
	started chan struct{}
}

var setupLoggerOnce = &sync.Once{}

func New(config Config) (*App, error) {
	// Configure logger
	// TODO(albrow): Don't use global variables for log settings.
	setupLoggerOnce.Do(func() {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.Level(config.Verbosity))
		log.AddHook(loghooks.NewKeySuffixHook())
	})

	// Add custom contract addresses if needed.
	var chainIDToContractAddresses map[int]ethereum.ContractAddresses
	if config.CustomContractAddresses != "" {
		if err := parseAndAddCustomContractAddresses(config.EthereumChainID, chainIDToContractAddresses, config.CustomContractAddresses); err != nil {
			return nil, err
		}
	}

	// Load private key and add peer ID hook.
	privKeyPath := filepath.Join(config.DataDir, "keys", "privkey")
	privKey, err := initPrivateKey(privKeyPath)
	if err != nil {
		return nil, err
	}
	peerID, err := peer.IDFromPrivateKey(privKey)
	if err != nil {
		return nil, err
	}
	log.AddHook(loghooks.NewPeerIDHook(peerID))

	if config.EthereumRPCMaxContentLength < constants.MaxOrderSizeInBytes {
		return nil, fmt.Errorf("Cannot set `EthereumRPCMaxContentLength` to be less then MaxOrderSizeInBytes: %d", constants.MaxOrderSizeInBytes)
	}
	config = unquoteConfig(config)

	if config.EnableEthereumRPCRateLimiting {
		// Ensure ETHEREUM_RPC_MAX_REQUESTS_PER_24_HR_UTC is reasonably set given BLOCK_POLLING_INTERVAL
		per24HrPollingRequests := int((24 * time.Hour) / config.BlockPollingInterval)
		minNumOfEthRPCRequestsIn24HrPeriod := per24HrPollingRequests + estimatedNonPollingEthereumRPCRequestsPer24Hrs
		if minNumOfEthRPCRequestsIn24HrPeriod > config.EthereumRPCMaxRequestsPer24HrUTC {
			return nil, fmt.Errorf(
				"Given BLOCK_POLLING_INTERVAL (%s), there are insufficient remaining ETH RPC requests in a 24hr period for Mesh to function properly. Increase ETHEREUM_RPC_MAX_REQUESTS_PER_24_HR_UTC to at least %d (currently configured to: %d)",
				config.BlockPollingInterval,
				minNumOfEthRPCRequestsIn24HrPeriod,
				config.EthereumRPCMaxRequestsPer24HrUTC,
			)
		}
	}

	// Initialize db
	databasePath := filepath.Join(config.DataDir, "db")
	meshDB, err := meshdb.New(databasePath)
	if err != nil {
		return nil, err
	}

	// Initialize metadata and check stored chain id (if any).
	metadata, err := initMetadata(config.EthereumChainID, meshDB)
	if err != nil {
		return nil, err
	}

	// Initialize ETH JSON-RPC RateLimiter
	var ethRPCRateLimiter ratelimit.RateLimiter
	if config.EnableEthereumRPCRateLimiting == false {
		ethRPCRateLimiter = ratelimit.NewUnlimited()
	} else {
		clock := clock.New()
		var err error
		ethRPCRateLimiter, err = ratelimit.New(config.EthereumRPCMaxRequestsPer24HrUTC, config.EthereumRPCMaxRequestsPerSecond, meshDB, clock)
		if err != nil {
			return nil, err
		}
	}

	// Initialize the ETH client, which will be used by various watchers.
	var ethRPCClient ethclient.RPCClient
	if config.EthereumRPCClient != nil {
		if config.EthereumRPCURL != "" {
			log.Warn("Ignoring EthereumRPCURL and using the provided EthereumRPCClient")
		}
		ethRPCClient = config.EthereumRPCClient
	} else if config.EthereumRPCURL != "" {
		ethRPCClient, err = rpc.Dial(config.EthereumRPCURL)
		if err != nil {
			log.WithError(err).Error("Could not dial EthereumRPCURL")
			return nil, err
		}
	} else {
		return nil, errors.New("cannot initialize core.App: neither EthereumRPCURL or EthereumRPCClient were provided")
	}
	ethClient, err := ethrpcclient.New(ethRPCClient, ethereumRPCRequestTimeout, ethRPCRateLimiter)
	if err != nil {
		return nil, err
	}

	// Initialize block watcher (but don't start it yet).
	blockWatcherClient, err := blockwatch.NewRpcClient(ethClient)
	if err != nil {
		return nil, err
	}

	// Remove any old mini headers that might be lingering in the database.
	// See https://github.com/0xProject/0x-mesh/issues/667 and https://github.com/0xProject/0x-mesh/pull/716
	// We need to leave this in place becuase:
	//
	// 1. It is still necessary for anyone upgrading from older versions to >= 9.0.1 in the future.
	// 2. There's still a chance there are old MiniHeaders in the database (e.g. due to a sudden
	//    unexpected shut down).
	//
	err = meshDB.PruneMiniHeadersAboveRetentionLimit()
	if err != nil {
		return nil, err
	}

	topics := orderwatch.GetRelevantTopics()
	miniHeaders, err := meshDB.FindAllMiniHeadersSortedByNumber()
	if err != nil {
		return nil, err
	}
	stack := simplestack.New(meshDB.MiniHeaderRetentionLimit, miniHeaders)
	blockWatcherConfig := blockwatch.Config{
		Stack:           stack,
		PollingInterval: config.BlockPollingInterval,
		WithLogs:        true,
		Topics:          topics,
		Client:          blockWatcherClient,
	}
	blockWatcher := blockwatch.New(blockWatcherConfig)

	// Initialize the order validator
	orderValidator, err := ordervalidator.New(
		ethClient,
		config.EthereumChainID,
		config.EthereumRPCMaxContentLength,
	)
	if err != nil {
		return nil, err
	}

	// Initialize order watcher (but don't start it yet).
	orderWatcher, err := orderwatch.New(orderwatch.Config{
		MeshDB:            meshDB,
		BlockWatcher:      blockWatcher,
		OrderValidator:    orderValidator,
		ChainID:           config.EthereumChainID,
		MaxOrders:         config.MaxOrdersInStorage,
		MaxExpirationTime: metadata.MaxExpirationTime,
	})
	if err != nil {
		return nil, err
	}

	// Initialize the order filter
	orderFilter, err := orderfilter.New(config.EthereumChainID, config.CustomOrderFilter)
	if err != nil {
		return nil, fmt.Errorf("invalid custom order filter: %s", err.Error())
	}

	// Initialize remaining fields.
	snapshotExpirationWatcher := expirationwatch.New()

	app := &App{
		started:                    make(chan struct{}),
		config:                     config,
		privKey:                    privKey,
		peerID:                     peerID,
		chainID:                    config.EthereumChainID,
		blockWatcher:               blockWatcher,
		orderWatcher:               orderWatcher,
		orderValidator:             orderValidator,
		orderFilter:                orderFilter,
		snapshotExpirationWatcher:  snapshotExpirationWatcher,
		idToSnapshotInfo:           map[string]snapshotInfo{},
		ethRPCRateLimiter:          ethRPCRateLimiter,
		ethRPCClient:               ethClient,
		db:                         meshDB,
		chainIDToContractAddresses: chainIDToContractAddresses,
	}

	log.WithFields(map[string]interface{}{
		"config":  config,
		"version": version,
	}).Info("finished initializing core.App")

	return app, nil
}

// unquoteConfig removes quotes (if needed) from each string field in config.
func unquoteConfig(config Config) Config {
	if unquotedEthereumRPCURL, err := strconv.Unquote(config.EthereumRPCURL); err == nil {
		config.EthereumRPCURL = unquotedEthereumRPCURL
	}
	if unquotedDataDir, err := strconv.Unquote(config.DataDir); err == nil {
		config.DataDir = unquotedDataDir
	}
	return config
}

func getPublishTopics(chainID int, customFilter *orderfilter.Filter) ([]string, error) {
	defaultTopic, err := orderfilter.GetDefaultTopic(chainID)
	if err != nil {
		return nil, err
	}
	customTopic := customFilter.Topic()
	if defaultTopic == customTopic {
		// If we're just using the default order filter, we don't need to publish to
		// multiple topics.
		return []string{defaultTopic}, nil
	} else {
		// If we are using a custom order filter, publish to *both* the default
		// topic and the custom topic. All orders that match the custom order filter
		// must necessarily match the default filter. This also allows us to
		// implement cross-topic forwarding in the future.
		// See https://github.com/0xProject/0x-mesh/pull/563
		return []string{defaultTopic, customTopic}, nil
	}
}

func getRendezvous(chainID int) string {
	return fmt.Sprintf("/0x-mesh/network/%d/version/2", chainID)
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

func initMetadata(chainID int, meshDB *meshdb.MeshDB) (*meshdb.Metadata, error) {
	metadata, err := meshDB.GetMetadata()
	if err != nil {
		if _, ok := err.(db.NotFoundError); ok {
			// No stored metadata found (first startup)
			metadata = &meshdb.Metadata{
				EthereumChainID:   chainID,
				MaxExpirationTime: constants.UnlimitedExpirationTime,
			}
			if err := meshDB.SaveMetadata(metadata); err != nil {
				return nil, err
			}
			return metadata, nil
		}
		return nil, err
	}

	// on subsequent startups, verify we are on the same chain
	if metadata.EthereumChainID != chainID {
		err := fmt.Errorf("expected chainID to be %d but got %d", metadata.EthereumChainID, chainID)
		log.WithError(err).Error("Mesh previously started on different Ethereum chain; switch chainId or remove DB")
		return nil, err
	}
	return metadata, nil
}

func (app *App) Start(ctx context.Context) error {
	// Get the publish topics depending on our custom order filter.
	publishTopics, err := getPublishTopics(app.config.EthereumChainID, app.orderFilter)
	if err != nil {
		return err
	}

	// Create a child context so that we can preemptively cancel if there is an
	// error.
	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Below, we will start several independent goroutines. We use separate
	// channels to communicate errors and a waitgroup to wait for all goroutines
	// to exit.
	wg := &sync.WaitGroup{}

	// Close the database when the context is canceled.
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			log.Debug("closing app.db")
		}()
		<-innerCtx.Done()
		app.db.Close()
	}()

	// Start rateLimiter
	ethRPCRateLimiterErrChan := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			log.Debug("closing eth RPC rate limiter")
		}()
		ethRPCRateLimiterErrChan <- app.ethRPCRateLimiter.Start(innerCtx, rateLimiterCheckpointInterval)
	}()

	// Set up the snapshot expiration watcher pruning logic
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			log.Debug("closing snapshot expiration watcher")
		}()
		ticker := time.NewTicker(expirationPollingInterval)
		for {
			select {
			case <-innerCtx.Done():
				return
			case now := <-ticker.C:
				expiredSnapshots := app.snapshotExpirationWatcher.Prune(now)
				for _, expiredSnapshot := range expiredSnapshots {
					app.muIdToSnapshotInfo.Lock()
					delete(app.idToSnapshotInfo, expiredSnapshot.ID)
					app.muIdToSnapshotInfo.Unlock()
				}
			}
		}
	}()

	// Start the order watcher.
	orderWatcherErrChan := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			log.Debug("closing order watcher")
		}()
		log.Info("starting order watcher")
		orderWatcherErrChan <- app.orderWatcher.Watch(innerCtx)
	}()

	// Note: this is a blocking call so we won't continue set up until its finished.
	blocksElapsed, err := app.blockWatcher.FastSyncToLatestBlock(innerCtx)
	if err != nil {
		return err
	}

	// Start the block watcher.
	blockWatcherErrChan := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			log.Debug("closing block watcher")
		}()
		log.Info("starting block watcher")
		blockWatcherErrChan <- app.blockWatcher.Watch(innerCtx)
	}()

	// If Mesh is not caught up with the latest block found via Ethereum RPC, ensure orderWatcher
	// has processed at least one recent block before starting the P2P node and completing app start,
	// so that Mesh does not validate any orders at outdated block heights
	isCaughtUp := app.IsCaughtUpToLatestBlock(innerCtx)
	if !isCaughtUp {
		if err := app.orderWatcher.WaitForAtLeastOneBlockToBeProcessed(ctx); err != nil {
			return err
		}
	}

	if blocksElapsed >= constants.MaxBlocksStoredInNonArchiveNode {
		log.WithField("blocksElapsed", blocksElapsed).Info("More than 128 blocks have elapsed since last boot. Re-validating all orders stored (this can take a while)...")
		// Re-validate all orders since too many blocks have elapsed to fast-sync events
		if err := app.orderWatcher.Cleanup(innerCtx, 0*time.Minute); err != nil {
			return err
		}
	}

	// Initialize the p2p node.
	// Note(albrow): The main reason that we need to use a `started` channel in
	// some methods is that we cannot call p2p.New without passing in a context
	// (due to how libp2p works). This means that before app.Start is called,
	// app.node will be nil and attempting to call any methods on app.node will
	// panic with a nil pointer exception. All the other fields of core.App that
	// we need to use will have already been initialized and are ready to use.
	bootstrapList := p2p.DefaultBootstrapList
	if app.config.BootstrapList != "" {
		bootstrapList = strings.Split(app.config.BootstrapList, ",")
	}
	nodeConfig := p2p.Config{
		SubscribeTopic:         app.orderFilter.Topic(),
		PublishTopics:          publishTopics,
		TCPPort:                app.config.P2PTCPPort,
		WebSocketsPort:         app.config.P2PWebSocketsPort,
		Insecure:               false,
		PrivateKey:             app.privKey,
		MessageHandler:         app,
		RendezvousString:       getRendezvous(app.config.EthereumChainID),
		UseBootstrapList:       app.config.UseBootstrapList,
		BootstrapList:          bootstrapList,
		DataDir:                filepath.Join(app.config.DataDir, "p2p"),
		CustomMessageValidator: app.orderFilter.ValidatePubSubMessage,
	}
	app.node, err = p2p.New(innerCtx, nodeConfig)
	if err != nil {
		return err
	}

	// Register and start ordersync service.
	ordersyncSubprotocols := []ordersync.Subprotocol{
		NewFilteredPaginationSubprotocol(app, paginationSubprotocolPerPage),
	}
	app.ordersyncService = ordersync.New(innerCtx, app.node, ordersyncSubprotocols)
	orderSyncErrChan := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			log.Debug("closing ordersync service")
		}()
		if err := app.ordersyncService.GetOrders(innerCtx, ordersyncMinPeers); err != nil {
			orderSyncErrChan <- err
		}
	}()

	// Start the p2p node.
	p2pErrChan := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			log.Debug("closing p2p node")
		}()
		addrs := app.node.Multiaddrs()
		log.WithFields(map[string]interface{}{
			"addresses": addrs,
			"topic":     app.orderFilter.Topic(),
		}).Info("starting p2p node")

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				log.Debug("closing new addrs checker")
			}()
			app.periodicallyCheckForNewAddrs(innerCtx, addrs)
		}()

		p2pErrChan <- app.node.Start()
	}()

	// Start loop for periodically logging stats.
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			log.Debug("closing periodic stats logger")
		}()
		app.periodicallyLogStats(innerCtx)
	}()

	// Signal that the app has been started.
	log.Info("core.App was started")
	close(app.started)

	// If any error channel returns a non-nil error, we cancel the inner context
	// and return the error. Note that this means we only return the first error
	// that occurs.
	select {
	case err := <-p2pErrChan:
		if err != nil {
			log.WithError(err).Error("p2p node exited with error")
			cancel()
			return err
		}
	case err := <-orderWatcherErrChan:
		if err != nil {
			log.WithError(err).Error("order watcher exited with error")
			cancel()
			return err
		}
	case err := <-blockWatcherErrChan:
		if err != nil {
			log.WithError(err).Error("block watcher exited with error")
			cancel()
			return err
		}
	case err := <-ethRPCRateLimiterErrChan:
		if err != nil {
			log.WithError(err).Error("ETH JSON-RPC ratelimiter exited with error")
			cancel()
			return err
		}
	case err := <-orderSyncErrChan:
		if err != nil {
			log.WithError(err).Error("ordersync service exited with error")
			cancel()
			return err
		}
	}

	// Wait for all goroutines to exit. If we reached here it means we are done
	// and there are no errors.
	wg.Wait()
	log.Debug("app successfully closed")
	return nil
}

func (app *App) periodicallyCheckForNewAddrs(ctx context.Context, startingAddrs []ma.Multiaddr) {
	<-app.started

	// TODO(albrow): There might be a more efficient way to do this if we have access to
	// an event bus. See: https://github.com/libp2p/go-libp2p/issues/467
	seenAddrs := stringset.New()
	for _, addr := range startingAddrs {
		seenAddrs.Add(addr.String())
	}
	ticker := time.NewTicker(checkNewAddrInterval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
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
}

// ErrSnapshotNotFound is the error returned when a snapshot not found with a particular id
type ErrSnapshotNotFound struct {
	id string
}

func (e ErrSnapshotNotFound) Error() string {
	return fmt.Sprintf("No snapshot found with id: %s. To create a new snapshot, send a request with an empty snapshotID", e.id)
}

// ErrPerPageZero is the error returned when a GetOrders request specifies perPage to 0
type ErrPerPageZero struct{}

func (e ErrPerPageZero) Error() string {
	return "perPage cannot be zero"
}

// GetOrders retrieves paginated orders from the Mesh DB at a specific snapshot in time. Passing an empty
// string as `snapshotID` creates a new snapshot and returns the first set of results. To fetch all orders,
// continue to make requests supplying the `snapshotID` returned from the first request. After 1 minute of not
// received further requests referencing a specific snapshot, the snapshot expires and can no longer be used.
func (app *App) GetOrders(page, perPage int, snapshotID string) (*types.GetOrdersResponse, error) {
	<-app.started

	if perPage <= 0 {
		return nil, ErrPerPageZero{}
	}

	ordersInfos := []*types.OrderInfo{}
	var snapshot *db.Snapshot
	var createdAt time.Time
	if snapshotID == "" {
		// Create a new snapshot
		snapshotID = uuid.New().String()
		var err error
		snapshot, err = app.db.Orders.GetSnapshot()
		if err != nil {
			return nil, err
		}
		createdAt = time.Now().UTC()
		expirationTimestamp := time.Now().Add(1 * time.Minute)
		app.snapshotExpirationWatcher.Add(expirationTimestamp, snapshotID)
		app.muIdToSnapshotInfo.Lock()
		app.idToSnapshotInfo[snapshotID] = snapshotInfo{
			Snapshot:            snapshot,
			CreatedAt:           createdAt,
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
		createdAt = info.CreatedAt
		// Reset the snapshot's expiry
		app.snapshotExpirationWatcher.Remove(info.ExpirationTimestamp, snapshotID)
		expirationTimestamp := time.Now().Add(1 * time.Minute)
		app.snapshotExpirationWatcher.Add(expirationTimestamp, snapshotID)
		app.idToSnapshotInfo[snapshotID] = snapshotInfo{
			Snapshot:            snapshot,
			CreatedAt:           createdAt,
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
		ordersInfos = append(ordersInfos, &types.OrderInfo{
			OrderHash:                order.Hash,
			SignedOrder:              order.SignedOrder,
			FillableTakerAssetAmount: order.FillableTakerAssetAmount,
		})
	}

	getOrdersResponse := &types.GetOrdersResponse{
		SnapshotID:        snapshotID,
		SnapshotTimestamp: createdAt,
		OrdersInfos:       ordersInfos,
	}

	return getOrdersResponse, nil
}

// AddOrders can be used to add orders to Mesh. It validates the given orders
// and if they are valid, will store and eventually broadcast the orders to
// peers. If pinned is true, the orders will be marked as pinned, which means
// they will only be removed if they become unfillable and will not be removed
// due to having a high expiration time or any incentive mechanisms.
func (app *App) AddOrders(ctx context.Context, signedOrdersRaw []*json.RawMessage, pinned bool) (*ordervalidator.ValidationResults, error) {
	<-app.started

	allValidationResults := &ordervalidator.ValidationResults{
		Accepted: []*ordervalidator.AcceptedOrderInfo{},
		Rejected: []*ordervalidator.RejectedOrderInfo{},
	}
	orderHashesSeen := map[common.Hash]struct{}{}
	schemaValidOrders := []*zeroex.SignedOrder{}
	for _, signedOrderRaw := range signedOrdersRaw {
		signedOrderBytes := []byte(*signedOrderRaw)
		result, err := app.orderFilter.ValidateOrderJSON(signedOrderBytes)
		if err != nil {
			signedOrder := &zeroex.SignedOrder{}
			if err := signedOrder.UnmarshalJSON(signedOrderBytes); err != nil {
				signedOrder = nil
			}
			log.WithField("signedOrderRaw", string(signedOrderBytes)).Info("Unexpected error while attempting to validate signedOrderJSON against schema")
			allValidationResults.Rejected = append(allValidationResults.Rejected, &ordervalidator.RejectedOrderInfo{
				SignedOrder: signedOrder,
				Kind:        ordervalidator.MeshValidation,
				Status: ordervalidator.RejectedOrderStatus{
					Code:    ordervalidator.ROInvalidSchemaCode,
					Message: "order did not pass JSON-schema validation: Malformed JSON or empty payload",
				},
			})
			continue
		}
		if !result.Valid() {
			log.WithField("signedOrderRaw", string(signedOrderBytes)).Info("Order failed schema validation")
			status := ordervalidator.RejectedOrderStatus{
				Code:    ordervalidator.ROInvalidSchemaCode,
				Message: fmt.Sprintf("order did not pass JSON-schema validation: %s", result.Errors()),
			}
			signedOrder := &zeroex.SignedOrder{}
			if err := signedOrder.UnmarshalJSON(signedOrderBytes); err != nil {
				signedOrder = nil
			}
			allValidationResults.Rejected = append(allValidationResults.Rejected, &ordervalidator.RejectedOrderInfo{
				SignedOrder: signedOrder,
				Kind:        ordervalidator.MeshValidation,
				Status:      status,
			})
			continue
		}

		signedOrder := &zeroex.SignedOrder{}
		if err := signedOrder.UnmarshalJSON(signedOrderBytes); err != nil {
			// This error should never happen since the signedOrder already passed the JSON schema validation above
			log.WithField("signedOrderRaw", string(signedOrderBytes)).Error("Failed to unmarshal SignedOrder")
			return nil, err
		}

		orderHash, err := signedOrder.ComputeOrderHash()
		if err != nil {
			return nil, err
		}
		if _, alreadySeen := orderHashesSeen[orderHash]; alreadySeen {
			continue
		}

		schemaValidOrders = append(schemaValidOrders, signedOrder)
		orderHashesSeen[orderHash] = struct{}{}
	}

	validationResults, err := app.orderWatcher.ValidateAndStoreValidOrders(ctx, schemaValidOrders, pinned, app.chainID)
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
		// If the order isn't new, we don't add to OrderWatcher, log it's receipt
		// or share the order with peers.
		if !acceptedOrderInfo.IsNew {
			continue
		}

		log.WithFields(log.Fields{
			"orderHash": acceptedOrderInfo.OrderHash.String(),
		}).Debug("added new valid order via RPC or browser callback")

		// Share the order with our peers.
		if err := app.shareOrder(acceptedOrderInfo.SignedOrder); err != nil {
			return nil, err
		}
	}

	return allValidationResults, nil
}

// shareOrder immediately shares the given order on the GossipSub network.
func (app *App) shareOrder(order *zeroex.SignedOrder) error {
	<-app.started

	encoded, err := encoding.OrderToRawMessage(app.orderFilter.Topic(), order)
	if err != nil {
		return err
	}
	return app.node.Send(encoded)
}

// AddPeer can be used to manually connect to a new peer.
func (app *App) AddPeer(peerInfo peerstore.PeerInfo) error {
	<-app.started

	return app.node.Connect(peerInfo, peerConnectTimeout)
}

// GetStats retrieves stats about the Mesh node
func (app *App) GetStats() (*types.Stats, error) {
	<-app.started

	latestBlockHeader, err := app.db.FindLatestMiniHeader()
	if err != nil {
		return nil, err
	}
	latestBlock := types.LatestBlock{
		Number: int(latestBlockHeader.Number.Int64()),
		Hash:   latestBlockHeader.Hash,
	}
	notRemovedFilter := app.db.Orders.IsRemovedIndex.ValueFilter([]byte{0})
	numOrders, err := app.db.Orders.NewQuery(notRemovedFilter).Count()
	if err != nil {
		return nil, err
	}
	numOrdersIncludingRemoved, err := app.db.Orders.Count()
	if err != nil {
		return nil, err
	}
	numPinnedOrders, err := app.db.CountPinnedOrders()
	if err != nil {
		return nil, err
	}
	metadata, err := app.db.GetMetadata()
	if err != nil {
		return nil, err
	}

	response := &types.Stats{
		Version:                           version,
		PubSubTopic:                       app.orderFilter.Topic(),
		Rendezvous:                        getRendezvous(app.config.EthereumChainID),
		PeerID:                            app.peerID.String(),
		EthereumChainID:                   app.config.EthereumChainID,
		LatestBlock:                       latestBlock,
		NumOrders:                         numOrders,
		NumPeers:                          app.node.GetNumPeers(),
		NumOrdersIncludingRemoved:         numOrdersIncludingRemoved,
		NumPinnedOrders:                   numPinnedOrders,
		MaxExpirationTime:                 app.orderWatcher.MaxExpirationTime().String(),
		StartOfCurrentUTCDay:              metadata.StartOfCurrentUTCDay,
		EthRPCRequestsSentInCurrentUTCDay: metadata.EthRPCRequestsSentInCurrentUTCDay,
		EthRPCRateLimitExpiredRequests:    app.ethRPCClient.GetRateLimitDroppedRequests(),
	}
	return response, nil
}

func (app *App) periodicallyLogStats(ctx context.Context) {
	<-app.started

	ticker := time.NewTicker(logStatsInterval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
		}

		stats, err := app.GetStats()
		if err != nil {
			log.WithError(err).Error("could not get stats")
			continue
		}
		log.WithFields(log.Fields{
			"version":                           stats.Version,
			"pubSubTopic":                       stats.PubSubTopic,
			"rendezvous":                        stats.Rendezvous,
			"ethereumChainID":                   stats.EthereumChainID,
			"latestBlock":                       stats.LatestBlock,
			"numOrders":                         stats.NumOrders,
			"numOrdersIncludingRemoved":         stats.NumOrdersIncludingRemoved,
			"numPinnedOrders":                   stats.NumPinnedOrders,
			"numPeers":                          stats.NumPeers,
			"maxExpirationTime":                 stats.MaxExpirationTime,
			"startOfCurrentUTCDay":              stats.StartOfCurrentUTCDay,
			"ethRPCRequestsSentInCurrentUTCDay": stats.EthRPCRequestsSentInCurrentUTCDay,
			"ethRPCRateLimitExpiredRequests":    stats.EthRPCRateLimitExpiredRequests,
		}).Info("current stats")
	}
}

// SubscribeToOrderEvents let's one subscribe to order events emitted by the OrderWatcher
func (app *App) SubscribeToOrderEvents(sink chan<- []*zeroex.OrderEvent) event.Subscription {
	// app.orderWatcher is guaranteed to be initialized. No need to wait.
	subscription := app.orderWatcher.Subscribe(sink)
	return subscription
}

// IsCaughtUpToLatestBlock returns whether or not the latest block stored by Mesh corresponds
// to the latest block retrieved from it's Ethereum RPC endpoint
func (app *App) IsCaughtUpToLatestBlock(ctx context.Context) bool {
	latestBlockStored, err := app.db.FindLatestMiniHeader()
	if err != nil {
		if _, ok := err.(meshdb.MiniHeaderCollectionEmptyError); ok {
			return false
		}
		log.WithFields(map[string]interface{}{
			"err": err.Error(),
		}).Warn("failed to fetch the latest miniHeader from DB")
		return false
	}
	latestBlock, err := app.ethRPCClient.HeaderByNumber(ctx, nil)
	if err != nil {
		log.WithFields(map[string]interface{}{
			"err": err.Error(),
		}).Warn("failed to fetch the latest block header via Ethereum RPC")
		return false
	}
	return latestBlock.Number.Cmp(latestBlockStored.Number) == 0
}

func parseAndAddCustomContractAddresses(chainID int, chainIDToContractAddresses map[int]ethereum.ContractAddresses, encodedContractAddresses string) error {
	customAddresses := ethereum.ContractAddresses{}
	if err := json.Unmarshal([]byte(encodedContractAddresses), &customAddresses); err != nil {
		return fmt.Errorf("config.CustomContractAddresses is invalid: %s", err.Error())
	}
	if err := addContractAddressesForChainID(chainID, chainIDToContractAddresses, customAddresses); err != nil {
		return fmt.Errorf("config.CustomContractAddresses is invalid: %s", err.Error())
	}
	return nil
}

func addContractAddressesForChainID(chainID int, chainIDToContractAddresses map[int]ethereum.ContractAddresses, addresses ethereum.ContractAddresses) error {
	if chainID == 1 {
		return fmt.Errorf("cannot add contract addresses for chainID 1: addresses for mainnet are hard-coded and cannot be changed")
	}
	if addresses.Exchange == constants.NullAddress {
		return fmt.Errorf("cannot add contract addresses for chain ID %d: Exchange address is required", chainID)
	}
	if addresses.DevUtils == constants.NullAddress {
		return fmt.Errorf("cannot add contract addresses for chain ID %d: DevUtils address is required", chainID)
	}
	if addresses.ERC20Proxy == constants.NullAddress {
		return fmt.Errorf("cannot add contract addresses for chain ID %d: ERC20Proxy address is required", chainID)
	}
	if addresses.ERC721Proxy == constants.NullAddress {
		return fmt.Errorf("cannot add contract addresses for chain ID %d: ERC721Proxy address is required", chainID)
	}
	if addresses.ERC1155Proxy == constants.NullAddress {
		return fmt.Errorf("cannot add contract addresses for chain ID %d: ERC1155Proxy address is required", chainID)
	}
	// TODO(albrow): Uncomment this if we re-add coordinator support.
	// if addresses.CoordinatorRegistry == constants.NullAddress {
	// 	return fmt.Errorf("cannot add contract addresses for chain ID %d: CoordinatorRegistry address is required", chainID)
	// }
	chainIDToContractAddresses[chainID] = addresses
	return nil
}
