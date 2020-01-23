package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

func TestBrowserConversions(t *testing.T) {
	// Declare a context that will be used for all child processes, servers, and
	// other goroutines.
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	ctx, _ = chromedp.NewContext(ctx, chromedp.WithErrorf(t.Errorf))
	defer cancel()

	// Start a simple HTTP server to serve the web page for the browser node.
	ts := httptest.NewServer(http.FileServer(http.Dir("../../dist")))
	defer ts.Close()

	browserLogs := make(chan string, 1024)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		startBrowserInstance(t, ctx, ts.URL, browserLogs)
	}()

	messageWg := &sync.WaitGroup{}
	messageWg.Add(1)
	go func() {
		defer messageWg.Done()
		waitForLogSubstring(ctx, browserLogs, "test started")
		waitForLogSubstring(ctx, browserLogs, "orderEvents")
		waitForLogSubstring(ctx, browserLogs, "test ended")
	}()

	messageWg.Wait()
	cancel()
	wg.Wait()
}

// FIXME(jalextowle): This is a direct copy from integration-tests. I should find a way to avoid duplication.
func startBrowserInstance(t *testing.T, ctx context.Context, url string, browserLogMessages chan<- string) {
	// Use chromedp to visit the web page for the browser node.
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			switch ev.Type {
			case runtime.APITypeLog:
				// Send console.log events through the channel.
				for _, arg := range ev.Args {
					if arg.Type == runtime.TypeString {
						fmt.Println("[browser]: " + string(arg.Value))
						browserLogMessages <- string(arg.Value)
					}
				}
			case runtime.APITypeError:
				// Report any console.error events as test failures.
				for _, arg := range ev.Args {
					t.Errorf("JavaScript console error: (%s) %s", arg.Type, arg.Value)
				}
			}
		}
	})

	if err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		// The #jsFinished element is one specifically created to indicate that the
		// JavaScript code is done running.
		chromedp.WaitVisible("#jsFinished", chromedp.ByID),
	); err != nil && err != context.Canceled {
		t.Error(err)
	}
}

// waitForLogMessage blocks until a message is logged that psses the given
// filter or the context is done. If the message is logged before the context is
// done, it will return the entire message. Otherwise it returns an error.
func waitForLogMessage(ctx context.Context, logMessages <-chan string, filter func(string) bool) (string, error) {
	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("(%s) timed out waiting for message", ctx.Err())
		case msg := <-logMessages:
			if filter(msg) {
				return msg, nil
			}
		}
	}
}

func waitForLogSubstring(ctx context.Context, logMessages <-chan string, substr string) (string, error) {
	return waitForLogMessage(ctx, logMessages, func(msg string) bool {
		return strings.Contains(msg, substr)
	})
}
