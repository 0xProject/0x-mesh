package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/require"
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

	buildForTests(t, ctx)

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
	testContractEventPrelude(ctx, "ERC20ApprovalEvent", browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC20ApprovalEvent | parameter | owner): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC20ApprovalEvent | parameter | spender): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC20ApprovalEvent | parameter | value): true")

	// ERC20TransferEvent
	testContractEventPrelude(ctx, "ERC20TransferEvent", browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC20TransferEvent | parameter | from): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC20TransferEvent | parameter | to): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC20TransferEvent | parameter | value): true")

	// ERC721ApprovalEvent
	testContractEventPrelude(ctx, "ERC721ApprovalEvent", browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC721ApprovalEvent | parameter | owner): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC721ApprovalEvent | parameter | approved): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC721ApprovalEvent | parameter | tokenId): true")

	// ERC721ApprovalForAllEvent
	testContractEventPrelude(ctx, "ERC721ApprovalForAllEvent", browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC721ApprovalForAllEvent | parameter | owner): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC721ApprovalForAllEvent | parameter | operator): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC721ApprovalForAllEvent | parameter | approved): true")

	// ERC721TransferEvent
	testContractEventPrelude(ctx, "ERC721TransferEvent", browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC721TransferEvent | parameter | from): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC721TransferEvent | parameter | to): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC721TransferEvent | parameter | tokenId): true")

	// ERC1155ApprovalForAllEvent
	testContractEventPrelude(ctx, "ERC1155ApprovalForAllEvent", browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC1155ApprovalForAllEvent | parameter | owner): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC1155ApprovalForAllEvent | parameter | operator): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC1155ApprovalForAllEvent | parameter | approved): true")

	// ERC1155TransferSingleEvent
	testContractEventPrelude(ctx, "ERC1155TransferSingleEvent", browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC1155TransferSingleEvent | parameter | operator): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC1155TransferSingleEvent | parameter | from): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC1155TransferSingleEvent | parameter | to): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC1155TransferSingleEvent | parameter | id): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC1155TransferSingleEvent | parameter | value): true")

	// ERC1155TransferBatchEvent
	testContractEventPrelude(ctx, "ERC1155TransferBatchEvent", browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC1155TransferBatchEvent | parameter | operator): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC1155TransferBatchEvent | parameter | from): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC1155TransferBatchEvent | parameter | to): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC1155TransferBatchEvent | parameter | ids): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ERC1155TransferBatchEvent | parameter | values): true")

	// ExchangeFillEvent
	testContractEventPrelude(ctx, "ExchangeFillEvent", browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeFillEvent | parameter | makerAddress): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeFillEvent | parameter | takerAddress): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeFillEvent | parameter | senderAddress): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeFillEvent | parameter | feeRecipientAddress): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeFillEvent | parameter | makerAssetFilledAmount): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeFillEvent | parameter | takerAssetFilledAmount): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeFillEvent | parameter | makerFeePaid): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeFillEvent | parameter | takerFeePaid): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeFillEvent | parameter | protocolFeePaid): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeFillEvent | parameter | orderHash): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeFillEvent | parameter | makerAssetData): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeFillEvent | parameter | takerAssetData): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeFillEvent | parameter | makerFeeAssetData): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeFillEvent | parameter | takerFeeAssetData): true")

	// ExchangeCancelEvent
	testContractEventPrelude(ctx, "ExchangeCancelEvent", browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeCancelEvent | parameter | makerAddress): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeCancelEvent | parameter | senderAddress): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeCancelEvent | parameter | feeRecipientAddress): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeCancelEvent | parameter | orderHash): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeCancelEvent | parameter | makerAssetData): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeCancelEvent | parameter | takerAssetData): true")

	// ExchangeCancelUpToEvent
	testContractEventPrelude(ctx, "ExchangeCancelUpToEvent", browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeCancelUpToEvent | parameter | makerAddress): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeCancelUpToEvent | parameter | orderSenderAddress): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | ExchangeCancelUpToEvent | parameter | orderEpoch): true")

	// WethDepositEvent
	testContractEventPrelude(ctx, "WethDepositEvent", browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | WethDepositEvent | parameter | owner): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | WethDepositEvent | parameter | value): true")

	// WethWithdrawalEvent
	testContractEventPrelude(ctx, "WethWithdrawalEvent", browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | WethWithdrawalEvent | parameter | owner): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | WethWithdrawalEvent | parameter | value): true")

	// FooBarBazEvent
	testContractEventPrelude(ctx, "FooBarBazEvent", browserLogs)
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | FooBarBazEvent | parameter | owner): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | FooBarBazEvent | parameter | spender): true")
	waitForLogSubstring(ctx, browserLogs, "(contractEventTest | FooBarBazEvent | parameter | value): true")

	// NOTE(jalextowle): This logic ensures that tests that have been created in the
	// typescript file "conversion_test.ts" will fail without a corresponding log section
	// in this file.
	stop.Lock()
	stop.value = true
	stop.Unlock()
}

func testContractEventPrelude(ctx context.Context, description string, browserLogs chan string) {
	waitForLogSubstring(ctx, browserLogs, fmt.Sprintf("(contractEventTest | %s | blockHash): true", description))
	waitForLogSubstring(ctx, browserLogs, fmt.Sprintf("(contractEventTest | %s | txHash): true", description))
	waitForLogSubstring(ctx, browserLogs, fmt.Sprintf("(contractEventTest | %s | txIndex): true", description))
	waitForLogSubstring(ctx, browserLogs, fmt.Sprintf("(contractEventTest | %s | logIndex): true", description))
	waitForLogSubstring(ctx, browserLogs, fmt.Sprintf("(contractEventTest | %s | isRemoved): true", description))
	waitForLogSubstring(ctx, browserLogs, fmt.Sprintf("(contractEventTest | %s | address): true", description))
	waitForLogSubstring(ctx, browserLogs, fmt.Sprintf("(contractEventTest | %s | kind): true", description))
}

// FIXME(jalextowle): This is a modified copy from integration-tests. I should find a way to avoid duplication.
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

func buildForTests(t *testing.T, ctx context.Context) {
	fmt.Println("Clear yarn cache...")
	cmd := exec.CommandContext(ctx, "yarn", "cache", "clean")
	cmd.Dir = "../../"
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "could not clean yarn cache: %s", string(output))

	fmt.Println("Installing dependencies for Wasm binary and Typescript bindings...")
	cmd = exec.CommandContext(ctx, "yarn", "install")
	cmd.Dir = "../../"
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "could not install depedencies for TypeScript bindings: %s", string(output))

	fmt.Println("Building Wasm binary and Typescript bindings...")
	cmd = exec.CommandContext(ctx, "yarn", "build")
	cmd.Dir = "../../"
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "could not build Wasm binary and Typescript bindings: %s", string(output))
	fmt.Println("Done building everything")
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
