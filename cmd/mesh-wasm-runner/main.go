package main

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/plaid/go-envvar/envvar"

	log "github.com/sirupsen/logrus"
)

type WasmRunnerConfig struct {
	// Path to the directory containing the payload that the file server should serve from.
	WasmPayloadPath      string `envvar:"WASM_PAYLOAD_PATH" default:"./dist"`
	ChromeDevProtocolUrl string `envvar:"CHROME_DEV_PROTOCOL_URL" default:"http://localhost:9222"`
	FileServerPort       string `envvar:"FILE_SERVER_PORT" default:"8888"`
	// NOTE: This should be set to http://host.docker.internal:8888 if developing on MacOS
	NavigateToUrl    string `envvar:"NAVIGATE_TO_URL" default:"http://localhost:8888"`
	UseExecAllocator bool   `envvar:"USE_EXEC_ALLOCATOR" default:"true"`
}

// NOTE(oskar) - we copy the config from core.Confg and supplement it with JSON
// struct tags
type Config struct {
	// Verbosity is the logging verbosity: 0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug 6=trace
	Verbosity                        int           `envvar:"VERBOSITY" default:"2" json:"verbosity"`
	DataDir                          string        `envvar:"DATA_DIR" default:"0x_mesh"`
	P2PTCPPort                       int           `envvar:"P2P_TCP_PORT" default:"60558"`
	P2PWebSocketsPort                int           `envvar:"P2P_WEBSOCKETS_PORT" default:"60559"`
	EthereumRPCURL                   string        `envvar:"ETHEREUM_RPC_URL" json:"ethereumRPCURL"`
	EthereumChainID                  int           `envvar:"ETHEREUM_CHAIN_ID" json:"ethereumChainID"`
	UseBootstrapList                 bool          `envvar:"USE_BOOTSTRAP_LIST" default:"true" json:"useBootstrapList"`
	BootstrapList                    string        `envvar:"BOOTSTRAP_LIST" default:"" json:"bootstrapList"`
	BlockPollingInterval             time.Duration `envvar:"BLOCK_POLLING_INTERVAL" default:"5s" json:"blockPollingIntervalSeconds"`
	EthereumRPCMaxContentLength      int           `envvar:"ETHEREUM_RPC_MAX_CONTENT_LENGTH" default:"524288" json:"ethereumRPCMaxContentLength"`
	EnableEthereumRPCRateLimiting    bool          `envvar:"ENABLE_ETHEREUM_RPC_RATE_LIMITING" default:"true" json:"enableEthereumRPCRateLimiting"`
	EthereumRPCMaxRequestsPer24HrUTC int           `envvar:"ETHEREUM_RPC_MAX_REQUESTS_PER_24_HR_UTC" default:"200000" json:"ethereumRPCMaxRequestsPer24HrUTC"`
	EthereumRPCMaxRequestsPerSecond  float64       `envvar:"ETHEREUM_RPC_MAX_REQUESTS_PER_SECOND" default:"30" json:"ethereumRPCMaxRequestsPerSecond"`
	CustomContractAddresses          string        `envvar:"CUSTOM_CONTRACT_ADDRESSES" default:"" json:"customContractAddresses"`
	MaxOrdersInStorage               int           `envvar:"MAX_ORDERS_IN_STORAGE" default:"100000" json:"maxOrdersInStorage"`
	CustomOrderFilter                string        `envvar:"CUSTOM_ORDER_FILTER" default:"{}" json:"customOrderFilter"`
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	var wasmRunnerConfig WasmRunnerConfig
	if err := envvar.Parse(&wasmRunnerConfig); err != nil {
		log.Error(err)
	}

	var coreConfig Config
	if err := envvar.Parse(&coreConfig); err != nil {
		log.WithField("error", err.Error()).Fatal("could not parse environment variables")
	}

	var ctx context.Context
	if wasmRunnerConfig.UseExecAllocator {
		allocatorContext, cancel := chromedp.NewExecAllocator(context.Background(), chromedp.NoSandbox, chromedp.Headless)
		defer cancel()
		ctx, _ = chromedp.NewContext(allocatorContext)
	} else {
		allocatedUrl := getAllocatedBrowserURL(&wasmRunnerConfig)
		allocatorContext, cancel := chromedp.NewRemoteAllocator(context.Background(), allocatedUrl)
		defer cancel()
		ctx, _ = chromedp.NewContext(allocatorContext)
	}

	log.Info("starting file server")
	go startFileServer(ctx, &wasmRunnerConfig)

	browserLogMessages := make(chan string, 1024)
	log.Info("starting wasm mesh node")
	go startNode(ctx, wasmRunnerConfig.NavigateToUrl, &coreConfig, browserLogMessages)
	for {
		select {
		case <-ctx.Done():
			return
		case m := <-browserLogMessages:
			fmt.Println(m)

		}
	}
}
