package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/0xProject/0x-mesh/core"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/plaid/go-envvar/envvar"

	log "github.com/sirupsen/logrus"
)

type WasmRunnerConfig struct {
	// WasmPayloadPath
	WasmPayloadPath string `envvar:"WASM_PAYLOAD_PATH" default:"./packages/runner-wasm/dist"`

	ChromeDevProtocolUrl string `envvar:"CHROME_DEV_PROTOCOL_URL" default:"http://localhost:9222"`
}

func wrapHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithField("Path", r.URL.Path).Info("serving request")
		h.ServeHTTP(w, r)
	}
}

func startFileServer(ctx context.Context, config *WasmRunnerConfig) {
	log.Error(http.ListenAndServe(":8888", wrapHandler(http.FileServer(http.Dir(config.WasmPayloadPath)))))
}

func getJSONConfig(nodeConfig *core.Config) (string, error) {
	jsonBytes, err := json.Marshal(nodeConfig)
	return string(jsonBytes), err
}

func startNode(ctx context.Context, url string, nodeConfig *core.Config, browserLogMessages chan string) {
	// Use chromedp to visit the web page for the browser node.
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			switch ev.Type {
			case runtime.APITypeLog:
				// Send console.log events through the channel.
				for _, arg := range ev.Args {
					if arg.Type == runtime.TypeString {
						s, err := strconv.Unquote(string(arg.Value))
						if err != nil {
							log.Error(err)
							continue
						}
						browserLogMessages <- s
					}
				}
			case runtime.APITypeError:
				// Report any console.error events as test failures.
				for _, arg := range ev.Args {
					log.Errorf("JavaScript console error: (%s) %s %s", arg.Type, arg.Value, arg.Description)
				}
			}
		}
	})

	config, err := getJSONConfig(nodeConfig)
	if err != nil {
		log.Error(err)
	}
	log.Info("running chromedp commands")
	var res []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("#jsFinishedLoading", chromedp.ByID),
		chromedp.Evaluate(fmt.Sprintf("window.startMesh(%s);", config), &res),
		chromedp.WaitVisible("#jsFinished", chromedp.ByID),
	); err != nil && err != context.Canceled {
		log.Error(err)
	}
	log.Trace("result:", res)
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

	// coreConfig.EthereumRPCURL = "http://localhost:8545"
	// coreConfig.EthereumChainID = 1337

	allocatedUrl := getAllocatedBrowserURL(wasmRunnerConfig.ChromeDevProtocolUrl)
	allocatorContext, cancel := chromedp.NewRemoteAllocator(context.Background(), allocatedUrl)
	defer cancel()
	ctx, _ := chromedp.NewContext(allocatorContext)

	log.Info("starting file server")
	go startFileServer(ctx, &wasmRunnerConfig)

	browserLogMessages := make(chan string, 1024)
	log.Info("starting wasm mesh node")
	go startNode(ctx, "http://host.docker.internal:8888", &coreConfig, browserLogMessages)
	for {
		select {
		case <-ctx.Done():
			return
		case m := <-browserLogMessages:
			fmt.Println(m)

		}
	}
}
