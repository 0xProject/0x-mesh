// +build !js

package main

import (
	"context"
	"flag"
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

var browserConversionTestsEnabled bool

// The test `TestBrowserConversions` has a non-standard timeout, so it needs to be
// run seperately from other go tests.
func init() {
	flag.BoolVar(&browserConversionTestsEnabled, "enable-browser-conversion-tests", false, "enable browser conversion tests")
	testing.Init()
	flag.Parse()
}

// This test is the entry-point to the Browser Conversion Tests. The other relevant
// files are "../../conversion-test/conversion_test.ts" and "./main.go."
//
// This entry-point builds the repository so that the current codebase is tested
// rather than an old version. After building, this test will register any test
// cases that are expected to be executed. Finally, a file server and a headless
// browser are initialized. The file server serves the "../../dist" directory,
// which will contain a webpage that contains a test script that will instantiate
// the Wasm buffer and execute the tests. As the tests are executed, test results
// are logged to the browser console, which this test will be able to access. As
// these logs are received, they are verified against registered test cases. Test
// failures, unexpected tests, and missing tests are all failure conditions for this
// test.
func TestBrowserConversions(t *testing.T) {
	if !browserConversionTestsEnabled {
		t.Skip("Browser conversion tests are disabled. You can enable them with the --enable-browser-conversion-tests flag")
	}

	// Declare a context that will be used for all child processes, servers, and
	// other goroutines.
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	ctx, _ = chromedp.NewContext(ctx, chromedp.WithErrorf(t.Errorf))
	defer cancel()

	buildForTests(t, ctx)

	// Register the Go --> Typescript test cases that should be logged.
	registerContractEventTest()
	registerGetOrdersResponseTest("EmptyOrderInfo", 0)
	registerGetOrdersResponseTest("OneOrderInfo", 1)
	registerGetOrdersResponseTest("TwoOrderInfos", 2)
	registerOrderEventTest("EmptyContractEvents", 0)
	registerOrderEventTest("ExchangeFillContractEvent", 1)
	registerSignedOrderTest("NullAssetData")
	registerSignedOrderTest("NonNullAssetData")
	registerStatsTest("RealisticStats")
	registerValidationResultsTest("EmptyValidationResults", 0, 0)
	registerValidationResultsTest("OneAcceptedResult", 1, 0)
	registerValidationResultsTest("OneRejectedResult", 0, 1)
	registerValidationResultsTest("RealisticValidationResults", 2, 1)

	// Register the Typescript --> Go test cases that should be logged.
	registerConvertConfigTest("NullConfig")
	registerConvertConfigTest("UndefinedConfig")
	registerConvertConfigTest("EmptyConfig")
	registerConvertConfigTest("MinimalConfig")
	registerConvertConfigTest("FullConfig")

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
	registerContractEventField(description, fmt.Sprintf("parameters.%s", param))
}

func registerContractEventField(description string, field string) {
	registerTest(fmt.Sprintf("(contractEvent | %s | %s)", description, field))
}

func registerConvertConfigTest(description string) {
	registerConvertConfigField(description, "config")
	registerConvertConfigField(description, "web3Provider")
	registerConvertConfigField(description, "err")
}

func registerConvertConfigField(description string, field string) {
	registerTest(fmt.Sprintf("(convertConfig | %s | %s)", description, field))
}

func registerGetOrdersResponseTest(description string, orderInfoLength int) {
	registerGetOrdersResponseField(description, "snapshotID")
	registerGetOrdersResponseField(description, "snapshotTimestamp")
	registerGetOrdersResponseField(description, "orderInfo.length")
	for i := 0; i < orderInfoLength; i++ {
		registerGetOrdersResponseField(description, "orderInfo.orderHash")
		registerGetOrdersResponseField(description, "orderInfo.signedOrder.chainId")
		registerGetOrdersResponseField(description, "orderInfo.signedOrder.makerAddress")
		registerGetOrdersResponseField(description, "orderInfo.signedOrder.takerAddress")
		registerGetOrdersResponseField(description, "orderInfo.signedOrder.senderAddress")
		registerGetOrdersResponseField(description, "orderInfo.signedOrder.feeRecipientAddress")
		registerGetOrdersResponseField(description, "orderInfo.signedOrder.exchangeAddress")
		registerGetOrdersResponseField(description, "orderInfo.signedOrder.makerAssetData")
		registerGetOrdersResponseField(description, "orderInfo.signedOrder.makerAssetAmount")
		registerGetOrdersResponseField(description, "orderInfo.signedOrder.makerFeeAssetData")
		registerGetOrdersResponseField(description, "orderInfo.signedOrder.makerFee")
		registerGetOrdersResponseField(description, "orderInfo.signedOrder.takerAssetData")
		registerGetOrdersResponseField(description, "orderInfo.signedOrder.takerAssetAmount")
		registerGetOrdersResponseField(description, "orderInfo.signedOrder.takerFeeAssetData")
		registerGetOrdersResponseField(description, "orderInfo.signedOrder.takerFee")
		registerGetOrdersResponseField(description, "orderInfo.signedOrder.expirationTimeSeconds")
		registerGetOrdersResponseField(description, "orderInfo.signedOrder.salt")
		registerGetOrdersResponseField(description, "orderInfo.fillableTakerAssetAmount")
	}
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
	boilerplate := "contractEvents."
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
	boilerplate := "signedOrder."
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

func registerStatsTest(description string) {
	registerStatsField(description, "version")
	registerStatsField(description, "pubSubTopic")
	registerStatsField(description, "rendezvous")
	registerStatsField(description, "secondaryRendezvous")
	registerStatsField(description, "peerID")
	registerStatsField(description, "ethereumChainID")
	registerStatsField(description, "latestBlock | hash")
	registerStatsField(description, "latestBlock | number")
	registerStatsField(description, "numOrders")
	registerStatsField(description, "numPeers")
	registerStatsField(description, "numOrdersIncludingRemoved")
	registerStatsField(description, "numPinnedOrders")
	registerStatsField(description, "maxExpirationTime")
	registerStatsField(description, "startOfCurrentUTCDay")
	registerStatsField(description, "ethRPCRequestsSentInCurrentUTCDay")
	registerStatsField(description, "ethRPCRateLimitExpiredRequests")
}

func registerValidationResultsTest(description string, acceptedLength int, rejectedLength int) {
	registerValidationResultsField(description, "accepted.length")
	for i := 0; i < acceptedLength; i++ {
		registerValidationResultsField(description, "accepted.orderHash")
		registerValidationResultsField(description, "accepted.signedOrder.chainId")
		registerValidationResultsField(description, "accepted.signedOrder.makerAddress")
		registerValidationResultsField(description, "accepted.signedOrder.takerAddress")
		registerValidationResultsField(description, "accepted.signedOrder.senderAddress")
		registerValidationResultsField(description, "accepted.signedOrder.feeRecipientAddress")
		registerValidationResultsField(description, "accepted.signedOrder.exchangeAddress")
		registerValidationResultsField(description, "accepted.signedOrder.makerAssetData")
		registerValidationResultsField(description, "accepted.signedOrder.makerAssetAmount")
		registerValidationResultsField(description, "accepted.signedOrder.makerFeeAssetData")
		registerValidationResultsField(description, "accepted.signedOrder.makerFee")
		registerValidationResultsField(description, "accepted.signedOrder.takerAssetData")
		registerValidationResultsField(description, "accepted.signedOrder.takerAssetAmount")
		registerValidationResultsField(description, "accepted.signedOrder.takerFeeAssetData")
		registerValidationResultsField(description, "accepted.signedOrder.takerFee")
		registerValidationResultsField(description, "accepted.signedOrder.expirationTimeSeconds")
		registerValidationResultsField(description, "accepted.signedOrder.salt")
		registerValidationResultsField(description, "accepted.signedOrder.signature")
		registerValidationResultsField(description, "accepted.fillableTakerAssetAmount")
		registerValidationResultsField(description, "accepted.isNew")
	}

	registerValidationResultsField(description, "rejected.length")
	for i := 0; i < rejectedLength; i++ {
		registerValidationResultsField(description, "rejected.orderHash")
		registerValidationResultsField(description, "rejected.signedOrder.chainId")
		registerValidationResultsField(description, "rejected.signedOrder.makerAddress")
		registerValidationResultsField(description, "rejected.signedOrder.takerAddress")
		registerValidationResultsField(description, "rejected.signedOrder.senderAddress")
		registerValidationResultsField(description, "rejected.signedOrder.feeRecipientAddress")
		registerValidationResultsField(description, "rejected.signedOrder.exchangeAddress")
		registerValidationResultsField(description, "rejected.signedOrder.makerAssetData")
		registerValidationResultsField(description, "rejected.signedOrder.makerAssetAmount")
		registerValidationResultsField(description, "rejected.signedOrder.makerFeeAssetData")
		registerValidationResultsField(description, "rejected.signedOrder.makerFee")
		registerValidationResultsField(description, "rejected.signedOrder.takerAssetData")
		registerValidationResultsField(description, "rejected.signedOrder.takerAssetAmount")
		registerValidationResultsField(description, "rejected.signedOrder.takerFeeAssetData")
		registerValidationResultsField(description, "rejected.signedOrder.takerFee")
		registerValidationResultsField(description, "rejected.signedOrder.expirationTimeSeconds")
		registerValidationResultsField(description, "rejected.signedOrder.salt")
		registerValidationResultsField(description, "rejected.signedOrder.signature")
		registerValidationResultsField(description, "rejected.kind")
		registerValidationResultsField(description, "rejected.status.code")
		registerValidationResultsField(description, "rejected.status.message")
	}
}

func registerGetOrdersResponseField(description string, field string) {
	registerTest(fmt.Sprintf("(getOrdersResponse | %s | %s)", description, field))
}

func registerOrderEventField(description string, field string) {
	registerTest(fmt.Sprintf("(orderEvent | %s | %s)", description, field))
}

func registerSignedOrderField(description string, field string) {
	registerTest(fmt.Sprintf("(signedOrder | %s | %s)", description, field))
}

func registerStatsField(description string, field string) {
	registerTest(fmt.Sprintf("(stats | %s | %s)", description, field))
}

func registerValidationResultsField(description string, field string) {
	registerTest(fmt.Sprintf("(validationResults | %s | %s)", description, field))
}

func registerTest(test string) {
	testCases = append(testCases, fmt.Sprintf(`"%s: true"`, test))
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
					} else if arg.Type == runtime.TypeString && len(string(arg.Value)) > 1 {
						t.Errorf("Unexpected extra test results: %s\n", string(arg.Value))
					} else {
						t.Errorf("Unexpected non-string output: %s %s\n", arg.Type, arg.Value)
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
	if count < len(testCases) {
		for i := count; i < len(testCases); i++ {
			t.Errorf("expected: %s actual: no response", testCases[i])
		}
	}
}

func buildForTests(t *testing.T, ctx context.Context) {
	fmt.Println("Clear yarn cache...")
	cmd := exec.CommandContext(ctx, "yarn", "cache", "clean")
	cmd.Dir = "../../../../"
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "could not clean yarn cache: %s", string(output))

	fmt.Println("Installing dependencies for Wasm binary and Typescript bindings...")
	cmd = exec.CommandContext(ctx, "yarn", "install")
	cmd.Dir = "../../../../"
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "could not install depedencies for TypeScript bindings: %s", string(output))

	fmt.Println("Building Wasm binary and Typescript bindings...")
	cmd = exec.CommandContext(ctx, "yarn", "build")
	cmd.Dir = "../../../../"
	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "could not build Wasm binary and Typescript bindings: %s", string(output))
	fmt.Println("Finished building for tests")
}
