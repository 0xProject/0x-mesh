package main

import (
	"context"
	"fmt"

	"github.com/0xProject/0x-mesh/core"
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

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	var wasmRunnerConfig WasmRunnerConfig
	if err := envvar.Parse(&wasmRunnerConfig); err != nil {
		log.Error(err)
	}

	var coreConfig core.Config
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
