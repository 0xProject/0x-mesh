package main

import (
	"time"

	"github.com/0xProject/0x-mesh/blockwatch"
	"github.com/0xProject/0x-mesh/core"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/zeroex/orderwatch"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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
	// ETHBalanceCheckerAddress is the hex-encoded address of the ETH balance
	// checker smart contract.
	ETHBalanceCheckerAddress string `envvar:"ETH_BALANCE_CHECKER_ADDRESS"`
}

type application struct {
	env          meshEnvVars
	db           *meshdb.MeshDB
	node         *core.Node
	blockWatcher *blockwatch.Watcher
	orderWatcher *orderwatch.Watcher
	ethWathcher  *ethereum.ETHWatcher
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

	// Initialize block watcher (but don't start it yet).
	rpcClient, err := blockwatch.NewRpcClient(env.EthereumRPCURL, ethereumRPCRequestTimeout)
	if err != nil {
		return nil, err
	}
	blockWatcherConfig := blockwatch.Config{
		MeshDB:          db,
		PollingInterval: blockWatcherPollingInterval,
		// TODO(albrow): Start at current block or by checking the database for the
		// most recent block?
		StartBlockDepth:     0,
		BlockRetentionLimit: blockWatcherRetentionLimit,
		WithLogs:            true,
		// TODO(albrow): What should Topics be?
		Topics: nil,
		Client: rpcClient,
	}
	blockWatcher := blockwatch.New(blockWatcherConfig)

	// Initialize order watcher (but don't start it yet).
	orderWatcher, err := orderwatch.New(blockWatcher, rpcClient)
	if err != nil {
		return nil, err
	}

	// Initialize the ETH balance watcher (but don't start it yet).
	ethClient, err := ethclient.Dial(env.EthereumRPCURL)
	if err != nil {
		return nil, err
	}
	ethBalanceCheckerAddress := common.HexToAddress(env.ETHBalanceCheckerAddress)
	ethWatcher, err := ethereum.NewETHWatcher(ethWatcherPollingInterval, ethClient, ethBalanceCheckerAddress)
	if err != nil {
		return nil, err
	}

	app := &application{
		env:          env,
		db:           db,
		blockWatcher: blockWatcher,
		orderWatcher: orderWatcher,
		ethWathcher:  ethWatcher,
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

func (app *application) GetMessagesToShare(max int) ([][]byte, error) {
	// TODO(albrow): Implement this.
	return nil, nil
}

func (app *application) Validate(messages []*core.Message) ([]*core.Message, error) {
	// TODO(albrow): Implement this.
	return messages, nil
}

func (app *application) Store([]*core.Message) error {
	// TODO(albrow): Implement this.
	return nil
}

func (app *application) start() error {
	// TODO(albrow): Implement this.
	return nil
}
