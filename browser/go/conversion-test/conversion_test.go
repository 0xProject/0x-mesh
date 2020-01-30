package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/require"
)

var testCases []string

func TestBrowserConversions(t *testing.T) {
	// Declare a context that will be used for all child processes, servers, and
	// other goroutines.
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	ctx, _ = chromedp.NewContext(ctx, chromedp.WithErrorf(t.Errorf))
	defer cancel()

	buildForTests(t, ctx)
	registerOrderEventTest("EmptyContractEvents", 0)
	registerOrderEventTest("ExchangeFillContractEvent", 1)
	registerSignedOrderTest("NullAssetData")
	registerSignedOrderTest("NonNullAssetData")
	registerContractEventTest()

	// Start a simple HTTP server to serve the web page for the browser node.
	ts := httptest.NewServer(http.FileServer(http.Dir("../../dist")))
	defer ts.Close()

	done := make(chan interface{})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		startBrowserInstance(t, ctx, ts.URL, done)
	}()

	go func() {
		select {
		case <-done:
			// NOTE(jalextowle): It is somewhat useful to know whether or not
			// there are test results that were logged in the typescript but were
			// not registered in this test file. For these purposes, we wait for
			// last logs to appear before closing the test. Logs that are logged
			// after the sleeping period will still be ignored.
			time.Sleep(2 * time.Second)
			cancel()
		}
	}()

	wg.Wait()
}

// NOTE(jalextowle): The ContractEvent tests and the event decoder tests are combined
// because it is important that `ContractEvents` are converted correctly, which directly
// tests event decoding.
func registerContractEventTest() {
	// ERC20ApprovalEvent
	registerContractEventPrelude("ERC20ApprovalEvent")
	registerContractEventParams("ERC20ApprovalEvent", "owner")
	registerContractEventParams("ERC20ApprovalEvent", "spender")
	registerContractEventParams("ERC20ApprovalEvent", "value")

	// ERC20TransferEvent
	registerContractEventPrelude("ERC20TransferEvent")
	registerContractEventParams("ERC20TransferEvent", "from")
	registerContractEventParams("ERC20TransferEvent", "to")
	registerContractEventParams("ERC20TransferEvent", "value")

	// ERC721ApprovalEvent
	registerContractEventPrelude("ERC721ApprovalEvent")
	registerContractEventParams("ERC721ApprovalEvent", "owner")
	registerContractEventParams("ERC721ApprovalEvent", "approved")
	registerContractEventParams("ERC721ApprovalEvent", "tokenId")

	// ERC721ApprovalForAllEvent
	registerContractEventPrelude("ERC721ApprovalForAllEvent")
	registerContractEventParams("ERC721ApprovalForAllEvent", "owner")
	registerContractEventParams("ERC721ApprovalForAllEvent", "operator")
	registerContractEventParams("ERC721ApprovalForAllEvent", "approved")

	// ERC721TransferEvent
	registerContractEventPrelude("ERC721TransferEvent")
	registerContractEventParams("ERC721TransferEvent", "from")
	registerContractEventParams("ERC721TransferEvent", "to")
	registerContractEventParams("ERC721TransferEvent", "tokenId")

	// ERC1155ApprovalForAllEvent
	registerContractEventPrelude("ERC1155ApprovalForAllEvent")
	registerContractEventParams("ERC1155ApprovalForAllEvent", "owner")
	registerContractEventParams("ERC1155ApprovalForAllEvent", "operator")
	registerContractEventParams("ERC1155ApprovalForAllEvent", "approved")

	// ERC1155TransferSingleEvent
	registerContractEventPrelude("ERC1155TransferSingleEvent")
	registerContractEventParams("ERC1155TransferSingleEvent", "operator")
	registerContractEventParams("ERC1155TransferSingleEvent", "from")
	registerContractEventParams("ERC1155TransferSingleEvent", "to")
	registerContractEventParams("ERC1155TransferSingleEvent", "id")
	registerContractEventParams("ERC1155TransferSingleEvent", "value")

	// ERC1155TransferBatchEvent
	registerContractEventPrelude("ERC1155TransferBatchEvent")
	registerContractEventParams("ERC1155TransferBatchEvent", "operator")
	registerContractEventParams("ERC1155TransferBatchEvent", "from")
	registerContractEventParams("ERC1155TransferBatchEvent", "to")
	registerContractEventParams("ERC1155TransferBatchEvent", "ids")
	registerContractEventParams("ERC1155TransferBatchEvent", "values")

	// ExchangeFillEvent
	registerContractEventPrelude("ExchangeFillEvent")
	registerContractEventParams("ExchangeFillEvent", "makerAddress")
	registerContractEventParams("ExchangeFillEvent", "takerAddress")
	registerContractEventParams("ExchangeFillEvent", "senderAddress")
	registerContractEventParams("ExchangeFillEvent", "feeRecipientAddress")
	registerContractEventParams("ExchangeFillEvent", "makerAssetFilledAmount")
	registerContractEventParams("ExchangeFillEvent", "takerAssetFilledAmount")
	registerContractEventParams("ExchangeFillEvent", "makerFeePaid")
	registerContractEventParams("ExchangeFillEvent", "takerFeePaid")
	registerContractEventParams("ExchangeFillEvent", "protocolFeePaid")
	registerContractEventParams("ExchangeFillEvent", "orderHash")
	registerContractEventParams("ExchangeFillEvent", "makerAssetData")
	registerContractEventParams("ExchangeFillEvent", "takerAssetData")
	registerContractEventParams("ExchangeFillEvent", "makerFeeAssetData")
	registerContractEventParams("ExchangeFillEvent", "takerFeeAssetData")

	// ExchangeCancelEvent
	registerContractEventPrelude("ExchangeCancelEvent")
	registerContractEventParams("ExchangeCancelEvent", "makerAddress")
	registerContractEventParams("ExchangeCancelEvent", "senderAddress")
	registerContractEventParams("ExchangeCancelEvent", "feeRecipientAddress")
	registerContractEventParams("ExchangeCancelEvent", "orderHash")
	registerContractEventParams("ExchangeCancelEvent", "makerAssetData")
	registerContractEventParams("ExchangeCancelEvent", "takerAssetData")

	// ExchangeCancelUpToEvent
	registerContractEventPrelude("ExchangeCancelUpToEvent")
	registerContractEventParams("ExchangeCancelUpToEvent", "makerAddress")
	registerContractEventParams("ExchangeCancelUpToEvent", "orderSenderAddress")
	registerContractEventParams("ExchangeCancelUpToEvent", "orderEpoch")

	// WethDepositEvent
	registerContractEventPrelude("WethDepositEvent")
	registerContractEventParams("WethDepositEvent", "owner")
	registerContractEventParams("WethDepositEvent", "value")

	// WethWithdrawalEvent
	registerContractEventPrelude("WethWithdrawalEvent")
	registerContractEventParams("WethWithdrawalEvent", "owner")
	registerContractEventParams("WethWithdrawalEvent", "value")

	// FooBarBazEvent
	registerContractEventPrelude("FooBarBazEvent")
	registerContractEventParams("FooBarBazEvent", "owner")
	registerContractEventParams("FooBarBazEvent", "spender")
	registerContractEventParams("FooBarBazEvent", "value")

	fmt.Println("Done registering ContractEvent test")
}

func registerContractEventPrelude(description string) {
	registerContractEventField(description, "blockHash")
	registerContractEventField(description, "txHash")
	registerContractEventField(description, "txIndex")
	registerContractEventField(description, "logIndex")
	registerContractEventField(description, "isRemoved")
	registerContractEventField(description, "address")
	registerContractEventField(description, "kind")
}

func registerContractEventParams(description string, param string) {
	registerContractEventField(description, fmt.Sprintf("parameter | %s", param))
}

func registerContractEventField(description string, field string) {
	registerTest(fmt.Sprintf("(contractEventTest | %s | %s)", description, field))
}

func registerSignedOrderTest(description string) {
	registerSignedOrderField(description, "chainId")
	registerSignedOrderField(description, "makerAddress")
	registerSignedOrderField(description, "takerAddress")
	registerSignedOrderField(description, "senderAddress")
	registerSignedOrderField(description, "feeRecipientAddress")
	registerSignedOrderField(description, "exchangeAddress")
	registerSignedOrderField(description, "makerAssetData")
	registerSignedOrderField(description, "makerAssetAmount")
	registerSignedOrderField(description, "makerFeeAssetData")
	registerSignedOrderField(description, "makerFee")
	registerSignedOrderField(description, "takerAssetData")
	registerSignedOrderField(description, "takerAssetAmount")
	registerSignedOrderField(description, "takerFeeAssetData")
	registerSignedOrderField(description, "takerFee")
	registerSignedOrderField(description, "expirationTimeSeconds")
	registerSignedOrderField(description, "salt")
	registerSignedOrderField(description, "signature")
}

func registerOrderEventTest(description string, length int) {
	registerOrderEventField(description, "timestamp")
	registerOrderEventField(description, "orderHash")
	registerOrderEventField(description, "endState")
	registerOrderEventField(description, "fillableTakerAssetAmount")
	registerOrderEventSignedOrder(description)
	registerOrderEventContractEventsPrelude(description, length)
}

func registerOrderEventContractEventsPrelude(description string, length int) {
	boilerplate := "contractEvents | "
	registerOrderEventField(description, boilerplate+"length")
	if length == 0 {
		return
	}
	registerOrderEventField(description, boilerplate+"blockHash")
	registerOrderEventField(description, boilerplate+"txHash")
	registerOrderEventField(description, boilerplate+"txIndex")
	registerOrderEventField(description, boilerplate+"logIndex")
	registerOrderEventField(description, boilerplate+"isRemoved")
	registerOrderEventField(description, boilerplate+"address")
	registerOrderEventField(description, boilerplate+"kind")
}

func registerOrderEventSignedOrder(description string) {
	boilerplate := "signedOrder | "
	registerOrderEventField(description, boilerplate+"chainId")
	registerOrderEventField(description, boilerplate+"makerAddress")
	registerOrderEventField(description, boilerplate+"takerAddress")
	registerOrderEventField(description, boilerplate+"senderAddress")
	registerOrderEventField(description, boilerplate+"feeRecipientAddress")
	registerOrderEventField(description, boilerplate+"exchangeAddress")
	registerOrderEventField(description, boilerplate+"makerAssetData")
	registerOrderEventField(description, boilerplate+"makerAssetAmount")
	registerOrderEventField(description, boilerplate+"makerFeeAssetData")
	registerOrderEventField(description, boilerplate+"makerFee")
	registerOrderEventField(description, boilerplate+"takerAssetData")
	registerOrderEventField(description, boilerplate+"takerAssetAmount")
	registerOrderEventField(description, boilerplate+"takerFeeAssetData")
	registerOrderEventField(description, boilerplate+"takerFee")
	registerOrderEventField(description, boilerplate+"expirationTimeSeconds")
	registerOrderEventField(description, boilerplate+"salt")
}

func registerOrderEventField(description string, field string) {
	registerTest(fmt.Sprintf("(orderEventTest | %s | %s)", description, field))
}

func registerSignedOrderField(description string, field string) {
	registerTest(fmt.Sprintf("(signedOrderTest | %s | %s)", description, field))
}

func registerTest(test string) {
	testCases = append(testCases, fmt.Sprintf("\"%s: true\"", test))
}

func startBrowserInstance(t *testing.T, ctx context.Context, url string, done chan interface{}) {
	testLength := len(testCases)
	count := 0

	// Use chromedp to visit the web page for the browser node.
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			switch ev.Type {
			case runtime.APITypeLog:
				// Send console.log events through the channel.
				for _, arg := range ev.Args {
					if arg.Type == runtime.TypeString && count < testLength {
						if testCases[count] != string(arg.Value) {
							t.Errorf("expected: %s | actual: %s", testCases[count], string(arg.Value))
						}
						count++
						if count == testLength {
							done <- struct{}{}
						}
					} else {
						t.Errorf("Unexpected test results: %s", arg.Value)
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
