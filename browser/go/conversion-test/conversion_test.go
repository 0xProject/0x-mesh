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

var stop struct {
	value bool
	sync.Mutex
}

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
		testContractEvents(ctx, browserLogs)

		// NOTE(jalextowle): Sleep to wait for any late logs. This isn't a perfect
		// solution, but it should improve the DevEx of working with these tests since
		// this will allow the tests to fail for extra logs that aren't too late.
		time.Sleep(time.Second)
	}()

	messageWg.Wait()
	cancel()
	wg.Wait()
}

func testContractEvents(ctx context.Context, browserLogs chan string) {
	// ERC20ApprovalEvent
	testContractEventPrelude(ctx, 0, browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | 0 | parameter | owner): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | 0 | parameter | spender): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | 0 | parameter | value): true")

	// ERC20TransferEvent
	testContractEventPrelude(ctx, 1, browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | 1 | parameter | from): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | 1 | parameter | to): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | 1 | parameter | value): true")

	// ERC721ApprovalEvent
	testContractEventPrelude(ctx, 2, browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | 2 | parameter | owner): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | 2 | parameter | approved): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | 2 | parameter | tokenId): true")

	// ERC721ApprovalForAllEvent
	testContractEventPrelude(ctx, 3, browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | 3 | parameter | owner): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | 3 | parameter | operator): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | 3 | parameter | approved): true")

	// ERC721TransferEvent
	testContractEventPrelude(ctx, 4, browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | 4 | parameter | from): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | 4 | parameter | to): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | 4 | parameter | tokenId): true")

	// NOTE(jalextowle): This logic ensures that tests that have been created in the
	// typescript file "conversion_test.ts" will fail without a corresponding log section
	// in this file.
	stop.Lock()
	stop.value = true
	stop.Unlock()
}

func testContractEventPrelude(ctx context.Context, idx int, browserLogs chan string) {
	waitForLogSubstring(ctx, browserLogs, fmt.Sprintf("(contractEventTest | %d | blockHash): true", idx))
	waitForLogSubstring(ctx, browserLogs, fmt.Sprintf("(contractEventTest | %d | txHash): true", idx))
	waitForLogSubstring(ctx, browserLogs, fmt.Sprintf("(contractEventTest | %d | txIndex): true", idx))
	waitForLogSubstring(ctx, browserLogs, fmt.Sprintf("(contractEventTest | %d | logIndex): true", idx))
	waitForLogSubstring(ctx, browserLogs, fmt.Sprintf("(contractEventTest | %d | isRemoved): true", idx))
	waitForLogSubstring(ctx, browserLogs, fmt.Sprintf("(contractEventTest | %d | address): true", idx))
	waitForLogSubstring(ctx, browserLogs, fmt.Sprintf("(contractEventTest | %d | kind): true", idx))
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
					stop.Lock()
					shouldStop := stop.value
					stop.Unlock()
					if !shouldStop {
						if arg.Type == runtime.TypeString {
							fmt.Println("[browser]: " + string(arg.Value))
							browserLogMessages <- string(arg.Value)
						}
					} else {
						t.Errorf("Browser log after test: (%s)", arg.Value)
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

// TODO(jalextowle): This should be inlined in the tests so that we make sure that all of the logs contain our tests (and no extra).
// This will allow the "late log" logic to be removed.
//
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
