package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"

	log "github.com/sirupsen/logrus"
)

func wrapHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.WithField("path", r.URL.Path).Info("serving file server request")
		h.ServeHTTP(w, r)
	}
}

func startFileServer(ctx context.Context, config *WasmRunnerConfig) {
	log.Error(http.ListenAndServe(fmt.Sprintf(":%s", config.FileServerPort), wrapHandler(http.FileServer(http.Dir(config.WasmPayloadPath)))))
}

func getJSONConfig(nodeConfig *Config) (string, error) {
	jsonBytes, err := json.Marshal(nodeConfig)
	return string(jsonBytes), err
}

func getAllocatedBrowserURL(config *WasmRunnerConfig) string {
	resp, err := http.Get(fmt.Sprintf("%s/json/version", config.ChromeDevProtocolUrl))
	if err != nil {
		log.Fatal(err)
	}

	var result map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal(err)
	}
	return result["webSocketDebuggerUrl"].(string)
}

func startNode(ctx context.Context, url string, nodeConfig *Config, browserLogMessages chan string) {
	// Use chromedp to visit the web page for the browser node.
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			switch ev.Type {
			case runtime.APITypeLog:
				for _, arg := range ev.Args {
					if arg.Type == runtime.TypeString {
						// The logs arrive as quoted JSON strings
						s, err := strconv.Unquote(string(arg.Value))
						if err != nil {
							log.WithError(err).Error("failed to unquote log value")
							continue
						}
						browserLogMessages <- s
					}
				}
			case runtime.APITypeError:
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
	log.Info("chromedp commands finished")
	log.Trace("result:", res)
}
