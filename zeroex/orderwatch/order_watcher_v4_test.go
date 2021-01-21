// +build !js

package orderwatch

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/ethereum/blockwatch"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/scenario/orderopts"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestOrderWatcherV4TakerWhitelist(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)
	require.NoError(t, err)

	testCases := []*struct {
		order                     *zeroex.SignedOrderV4
		isTakerAddressWhitelisted bool
	}{
		{
			scenario.NewSignedTestOrderV4(t, orderopts.SetupMakerState(true)),
			true,
		},
		{
			scenario.NewSignedTestOrderV4(
				t,
				orderopts.SetupMakerState(true),
				orderopts.TakerAddress(ganacheAddresses.ExchangeProxyFlashWallet),
			),
			true,
		},
		{
			scenario.NewSignedTestOrderV4(
				t,
				orderopts.SetupMakerState(true),
				orderopts.TakerAddress(common.HexToAddress("0x1")),
			),
			false,
		},
	}

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	for _, testCase := range testCases {
		results, err := orderWatcher.ValidateAndStoreValidOrdersV4(ctx, []*zeroex.SignedOrderV4{testCase.order}, constants.TestChainID, false, &types.AddOrdersOpts{})
		require.NoError(t, err)
		if testCase.isTakerAddressWhitelisted {
			orderHash, err := testCase.order.ComputeOrderHash()
			require.NoError(t, err)
			assert.Len(t, results.Rejected, 0)
			require.Len(t, results.Accepted, 1)
			assert.Equal(t, results.Accepted[0].OrderHash, orderHash)
		} else {
			orderHash, err := testCase.order.ComputeOrderHash()
			require.NoError(t, err)
			assert.Len(t, results.Accepted, 0)
			require.Len(t, results.Rejected, 1)
			assert.Equal(t, results.Rejected[0].OrderHash, orderHash)
			assert.Equal(t, results.Rejected[0].Kind, ordervalidator.MeshValidation)
			assert.Equal(t, results.Rejected[0].Status, ordervalidator.ROTakerAddressNotAllowed)
		}
	}
}
func TestOrderWatcherV4DoesntStoreInvalidOrdersWithConfigurations(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	for _, testCase := range []*struct {
		description          string
		signedOrderGenerator func() *zeroex.SignedOrderV4
		addOrdersOpts        *types.AddOrdersOpts
	}{
		{
			description: "doesn't store cancelled orders when KeepCancelled is disabled",
			signedOrderGenerator: func() *zeroex.SignedOrderV4 {
				signedOrder := scenario.NewSignedTestOrderV4(t,
					orderopts.SetupMakerState(true),
					orderopts.MakerAssetData(scenario.ZRXAssetData),
				)
				// Cancel order
				opts := &bind.TransactOpts{
					From:   signedOrder.Maker,
					Signer: scenario.GetTestSignerFn(signedOrder.Maker),
				}
				trimmedOrder := signedOrder.EthereumAbiLimitOrder()
				txn, err := exchangeV4.CancelLimitOrder(opts, trimmedOrder)
				require.NoError(t, err)
				waitTxnSuccessfullyMined(t, ethClient, txn)
				return signedOrder
			},
			addOrdersOpts: &types.AddOrdersOpts{
				KeepCancelled:   false,
				KeepExpired:     true,
				KeepFullyFilled: true,
				KeepUnfunded:    true,
			},
		},
		{
			description: "doesn't store expired orders when KeepExpired is disabled",
			signedOrderGenerator: func() *zeroex.SignedOrderV4 {
				return scenario.NewSignedTestOrderV4(t,
					orderopts.SetupMakerState(true),
					orderopts.MakerAssetData(scenario.ZRXAssetData),
					orderopts.ExpirationTimeSeconds(big.NewInt(0)),
				)
			},
			addOrdersOpts: &types.AddOrdersOpts{
				KeepCancelled:   true,
				KeepExpired:     false,
				KeepFullyFilled: true,
				KeepUnfunded:    true,
			},
		},
		{
			description: "doesn't store fully filled orders when KeepFullyFilled is disabled",
			signedOrderGenerator: func() *zeroex.SignedOrderV4 {
				takerAddress := constants.GanacheAccount3
				signedOrder := scenario.NewSignedTestOrderV4(t,
					orderopts.SetupMakerState(true),
					orderopts.SetupTakerAddress(takerAddress),
					orderopts.MakerAssetData(scenario.ZRXAssetData),
				)
				// Fill order
				opts := &bind.TransactOpts{
					From:   takerAddress,
					Signer: scenario.GetTestSignerFn(takerAddress),
					Value:  big.NewInt(100000000000000000),
				}
				trimmedOrder := signedOrder.EthereumAbiLimitOrder()
				signature := signedOrder.EthereumAbiSignature()
				txn, err := exchangeV4.FillLimitOrder(opts, trimmedOrder, signature, signedOrder.TakerAmount)
				require.NoError(t, err)
				waitTxnSuccessfullyMined(t, ethClient, txn)
				return signedOrder
			},
			addOrdersOpts: &types.AddOrdersOpts{
				KeepCancelled:   true,
				KeepExpired:     true,
				KeepFullyFilled: false,
				KeepUnfunded:    true,
			},
		},
		{
			description: "doesn't store unfunded orders when KeepUnfunded is disabled",
			signedOrderGenerator: func() *zeroex.SignedOrderV4 {
				return scenario.NewSignedTestOrderV4(t,
					orderopts.MakerAssetData(scenario.ZRXAssetData),
					orderopts.MakerAssetAmount(big.NewInt(1)),
					orderopts.MakerFeeAssetData(scenario.WETHAssetData),
				)
			},
			addOrdersOpts: &types.AddOrdersOpts{
				KeepCancelled:   true,
				KeepExpired:     true,
				KeepFullyFilled: true,
				KeepUnfunded:    false,
			},
		},
	} {
		teardownSubTest := setupSubTest(t)

		ctx, cancel := context.WithCancel(context.Background())
		database, err := db.New(ctx, db.TestOptions())
		require.NoError(t, err, testCase.description)

		blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)

		signedOrder := testCase.signedOrderGenerator()

		err = blockWatcher.SyncToLatestBlock()
		require.NoError(t, err)

		validationResults, err := orderWatcher.ValidateAndStoreValidOrdersV4(ctx, []*zeroex.SignedOrderV4{signedOrder}, constants.TestChainID, false, testCase.addOrdersOpts)
		require.NoError(t, err)

		assert.Len(t, validationResults.Accepted, 0, testCase.description)
		assert.Len(t, validationResults.Rejected, 1, testCase.description)

		orders, err := database.FindOrdersV4(nil)
		require.NoError(t, err, testCase.description)
		require.Len(t, orders, 0)

		teardownSubTest(t)
		cancel()
	}
}

func TestOrderWatcherV4StoresValidOrdersWithConfigurations(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	for _, testCase := range []*struct {
		description            string
		expectedFillableAmount *big.Int
		signedOrderGenerator   func() *zeroex.SignedOrderV4
		addOrdersOpts          *types.AddOrdersOpts
		isExpired              bool
	}{
		{
			description:            "stores valid orders",
			expectedFillableAmount: big.NewInt(42),
			signedOrderGenerator: func() *zeroex.SignedOrderV4 {
				return scenario.NewSignedTestOrderV4(t,
					orderopts.SetupMakerState(true),
					orderopts.MakerAssetData(scenario.ZRXAssetData),
				)
			},
			addOrdersOpts: &types.AddOrdersOpts{},
		},
		{
			description:            "stores cancelled orders when KeepCancelled is enabled",
			expectedFillableAmount: big.NewInt(0),
			signedOrderGenerator: func() *zeroex.SignedOrderV4 {
				signedOrder := scenario.NewSignedTestOrderV4(t,
					orderopts.SetupMakerState(true),
					orderopts.MakerAssetData(scenario.ZRXAssetData),
				)
				// Cancel order
				opts := &bind.TransactOpts{
					From:   signedOrder.Maker,
					Signer: scenario.GetTestSignerFn(signedOrder.Maker),
				}
				trimmedOrder := signedOrder.EthereumAbiLimitOrder()
				txn, err := exchangeV4.CancelLimitOrder(opts, trimmedOrder)
				require.NoError(t, err)
				waitTxnSuccessfullyMined(t, ethClient, txn)
				return signedOrder
			},
			addOrdersOpts: &types.AddOrdersOpts{KeepCancelled: true},
		},
		{
			description:            "stores expired orders when KeepExpired is enabled",
			expectedFillableAmount: big.NewInt(0),
			signedOrderGenerator: func() *zeroex.SignedOrderV4 {
				return scenario.NewSignedTestOrderV4(t,
					orderopts.SetupMakerState(true),
					orderopts.MakerAssetData(scenario.ZRXAssetData),
					orderopts.ExpirationTimeSeconds(big.NewInt(0)),
				)
			},
			addOrdersOpts: &types.AddOrdersOpts{KeepExpired: true},
			isExpired:     true,
		},
		{
			description:            "stores fully filled orders when KeepFullyFilled is enabled",
			expectedFillableAmount: big.NewInt(0),
			signedOrderGenerator: func() *zeroex.SignedOrderV4 {
				takerAddress := constants.GanacheAccount3
				signedOrder := scenario.NewSignedTestOrderV4(t,
					orderopts.SetupMakerState(true),
					orderopts.SetupTakerAddress(takerAddress),
					orderopts.MakerAssetData(scenario.ZRXAssetData),
				)
				// Fill order
				opts := &bind.TransactOpts{
					From:   takerAddress,
					Signer: scenario.GetTestSignerFn(takerAddress),
					Value:  big.NewInt(100000000000000000),
				}
				trimmedOrder := signedOrder.EthereumAbiLimitOrder()
				signature := signedOrder.EthereumAbiSignature()
				txn, err := exchangeV4.FillLimitOrder(opts, trimmedOrder, signature, signedOrder.TakerAmount)
				require.NoError(t, err)
				waitTxnSuccessfullyMined(t, ethClient, txn)
				return signedOrder
			},
			addOrdersOpts: &types.AddOrdersOpts{KeepFullyFilled: true},
		},
		{
			description:            "stores unfunded orders when KeepUnfunded is enabled",
			expectedFillableAmount: big.NewInt(0),
			signedOrderGenerator: func() *zeroex.SignedOrderV4 {
				return scenario.NewSignedTestOrderV4(t,
					orderopts.MakerAssetData(scenario.ZRXAssetData),
					orderopts.MakerFee(big.NewInt(1)),
					orderopts.MakerFeeAssetData(scenario.WETHAssetData),
				)
			},
			addOrdersOpts: &types.AddOrdersOpts{KeepUnfunded: true},
		},
	} {
		teardownSubTest := setupSubTest(t)

		ctx, cancel := context.WithCancel(context.Background())
		database, err := db.New(ctx, db.TestOptions())
		require.NoError(t, err, testCase.description)

		blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)

		signedOrder := testCase.signedOrderGenerator()

		err = blockWatcher.SyncToLatestBlock()
		require.NoError(t, err)

		validationResults, err := orderWatcher.ValidateAndStoreValidOrdersV4(ctx, []*zeroex.SignedOrderV4{signedOrder}, constants.TestChainID, false, testCase.addOrdersOpts)
		require.NoError(t, err)

		isUnfillable := testCase.expectedFillableAmount.Cmp(big.NewInt(0)) == 0

		if isUnfillable {
			assert.Len(t, validationResults.Accepted, 0, testCase.description)
			assert.Len(t, validationResults.Rejected, 1, testCase.description)
		} else {
			assert.Len(t, validationResults.Accepted, 1, testCase.description)
			assert.Len(t, validationResults.Rejected, 0, testCase.description)
		}

		expectedOrderHash, err := signedOrder.ComputeOrderHash()
		require.NoError(t, err, testCase.description)

		latestStoredBlock, err := database.GetLatestMiniHeader()
		require.NoError(t, err, testCase.description)

		orders, err := database.FindOrdersV4(nil)
		require.NoError(t, err, testCase.description)
		require.Len(t, orders, 1)
		expectedOrderState := orderState{
			hash:               expectedOrderHash,
			isRemoved:          false,
			isUnfillable:       isUnfillable,
			isExpired:          testCase.isExpired,
			fillableAmount:     testCase.expectedFillableAmount,
			lastUpdated:        time.Now(),
			lastValidatedBlock: latestStoredBlock,
		}
		checkOrderState(t, expectedOrderState, orders[0])

		teardownSubTest(t)
		cancel()
	}
}

func TestOrderWatcherV4UnfundedInsufficientERC20Balance(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	for _, testCase := range []*struct {
		description     string
		addOrdersOpts   *types.AddOrdersOpts
		shouldBeRemoved bool
	}{
		{
			description:     "should be removed with no configurations",
			addOrdersOpts:   &types.AddOrdersOpts{},
			shouldBeRemoved: true,
		},
		{
			description:     "should be kept with KeepUnfunded",
			addOrdersOpts:   &types.AddOrdersOpts{KeepUnfunded: true},
			shouldBeRemoved: false,
		},
	} {
		teardownSubTest := setupSubTest(t)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		database, err := db.New(ctx, db.TestOptions())
		require.NoError(t, err, testCase.description)

		signedOrder := scenario.NewSignedTestOrderV4(t,
			orderopts.SetupMakerState(true),
			orderopts.MakerAssetData(scenario.ZRXAssetData),
		)
		expectedOrderHash, err := signedOrder.ComputeOrderHash()
		require.NoError(t, err, testCase.description)
		blockWatcher, orderEventsChan := setupOrderWatcherScenarioV4(ctx, t, database, signedOrder, testCase.addOrdersOpts)

		// Transfer makerAsset out of maker address
		opts := &bind.TransactOpts{
			From:   signedOrder.Maker,
			Signer: scenario.GetTestSignerFn(signedOrder.Maker),
		}
		txn, err := zrx.Transfer(opts, constants.GanacheAccount4, signedOrder.MakerAmount)
		require.NoError(t, err, testCase.description)
		waitTxnSuccessfullyMined(t, ethClient, txn)

		err = blockWatcher.SyncToLatestBlock()
		require.NoError(t, err, testCase.description)

		orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
		require.Len(t, orderEvents, 1, testCase.description)
		orderEvent := orderEvents[0]
		assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState, testCase.description)

		latestStoredBlock, err := database.GetLatestMiniHeader()
		require.NoError(t, err, testCase.description)
		orders, err := database.FindOrdersV4(nil)
		require.NoError(t, err, testCase.description)
		require.Len(t, orders, 1, testCase.description)
		expectedOrderState := orderState{
			hash:               expectedOrderHash,
			isRemoved:          testCase.shouldBeRemoved,
			isUnfillable:       true,
			fillableAmount:     big.NewInt(0),
			lastUpdated:        time.Now(),
			lastValidatedBlock: latestStoredBlock,
		}
		checkOrderState(t, expectedOrderState, orders[0])

		cancel()
		teardownSubTest(t)
	}
}

func TestOrderWatcherV4UnfundedInsufficientERC20Allowance(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	for _, testCase := range []*struct {
		description     string
		addOrdersOpts   *types.AddOrdersOpts
		shouldBeRemoved bool
	}{
		{
			description:     "should be removed with no configurations",
			addOrdersOpts:   &types.AddOrdersOpts{},
			shouldBeRemoved: true,
		},
		{
			description:     "should be kept with KeepUnfunded",
			addOrdersOpts:   &types.AddOrdersOpts{KeepUnfunded: true},
			shouldBeRemoved: false,
		},
	} {
		teardownSubTest := setupSubTest(t)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		database, err := db.New(ctx, db.TestOptions())
		require.NoError(t, err, testCase.description)

		signedOrder := scenario.NewSignedTestOrderV4(t,
			orderopts.SetupMakerState(true),
			orderopts.MakerAssetData(scenario.ZRXAssetData),
		)
		expectedOrderHash, err := signedOrder.ComputeOrderHash()
		require.NoError(t, err, testCase.description)
		blockWatcher, orderEventsChan := setupOrderWatcherScenarioV4(ctx, t, database, signedOrder, testCase.addOrdersOpts)

		// Remove Maker's ZRX approval to ExchangeProxy
		opts := &bind.TransactOpts{
			From:   signedOrder.Maker,
			Signer: scenario.GetTestSignerFn(signedOrder.Maker),
		}
		txn, err := zrx.Approve(opts, ganacheAddresses.ExchangeProxy, big.NewInt(0))
		require.NoError(t, err, testCase.description)
		waitTxnSuccessfullyMined(t, ethClient, txn)

		err = blockWatcher.SyncToLatestBlock()
		require.NoError(t, err, testCase.description)

		orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
		require.Len(t, orderEvents, 1, testCase.description)
		orderEvent := orderEvents[0]
		assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState, testCase.description)

		latestStoredBlock, err := database.GetLatestMiniHeader()
		require.NoError(t, err, testCase.description)
		orders, err := database.FindOrdersV4(nil)
		require.NoError(t, err, testCase.description)
		require.Len(t, orders, 1, testCase.description)
		expectedOrderState := orderState{
			hash:               expectedOrderHash,
			isRemoved:          testCase.shouldBeRemoved,
			isUnfillable:       true,
			fillableAmount:     big.NewInt(0),
			lastUpdated:        time.Now(),
			lastValidatedBlock: latestStoredBlock,
		}
		checkOrderState(t, expectedOrderState, orders[0])

		cancel()
		teardownSubTest(t)
	}
}

func TestOrderWatcherV4UnfundedThenFundedAgain(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	for _, testCase := range []*struct {
		description     string
		addOrdersOpts   *types.AddOrdersOpts
		shouldBeRemoved bool
	}{
		{
			description:     "should be removed with no configurations",
			addOrdersOpts:   &types.AddOrdersOpts{},
			shouldBeRemoved: true,
		},
		{
			description:     "should be kept with KeepUnfunded",
			addOrdersOpts:   &types.AddOrdersOpts{KeepUnfunded: true},
			shouldBeRemoved: false,
		},
	} {
		teardownSubTest := setupSubTest(t)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		database, err := db.New(ctx, db.TestOptions())
		require.NoError(t, err, testCase.description)

		signedOrder := scenario.NewSignedTestOrderV4(t,
			orderopts.SetupMakerState(true),
			orderopts.MakerAssetData(scenario.ZRXAssetData),
			orderopts.TakerAssetData(scenario.WETHAssetData),
		)
		expectedOrderHash, err := signedOrder.ComputeOrderHash()
		require.NoError(t, err, testCase.description)
		blockWatcher, orderEventsChan := setupOrderWatcherScenarioV4(ctx, t, database, signedOrder, testCase.addOrdersOpts)

		// Transfer makerAsset out of maker address
		opts := &bind.TransactOpts{
			From:   signedOrder.Maker,
			Signer: scenario.GetTestSignerFn(signedOrder.Maker),
		}
		txn, err := zrx.Transfer(opts, constants.GanacheAccount4, signedOrder.MakerAmount)
		require.NoError(t, err, testCase.description)
		waitTxnSuccessfullyMined(t, ethClient, txn)

		err = blockWatcher.SyncToLatestBlock()
		require.NoError(t, err, testCase.description)

		orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
		require.Len(t, orderEvents, 1, testCase.description)
		orderEvent := orderEvents[0]
		assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState, testCase.description)

		latestStoredBlock, err := database.GetLatestMiniHeader()
		require.NoError(t, err, testCase.description)
		orders, err := database.FindOrdersV4(nil)
		require.NoError(t, err, testCase.description)
		require.Len(t, orders, 1, testCase.description)
		expectedOrderState := orderState{
			hash:               expectedOrderHash,
			isRemoved:          testCase.shouldBeRemoved,
			isUnfillable:       true,
			fillableAmount:     big.NewInt(0),
			lastUpdated:        time.Now(),
			lastValidatedBlock: latestStoredBlock,
		}
		checkOrderState(t, expectedOrderState, orders[0])

		// Transfer makerAsset back to maker address
		zrxCoinbase := constants.GanacheAccount0
		opts = &bind.TransactOpts{
			From:   zrxCoinbase,
			Signer: scenario.GetTestSignerFn(zrxCoinbase),
		}
		txn, err = zrx.Transfer(opts, signedOrder.Maker, signedOrder.MakerAmount)
		require.NoError(t, err, testCase.description)
		waitTxnSuccessfullyMined(t, ethClient, txn)

		err = blockWatcher.SyncToLatestBlock()
		require.NoError(t, err, testCase.description)

		orderEvents = <-orderEventsChan
		require.Len(t, orderEvents, 1, testCase.description)
		orderEvent = orderEvents[0]
		assert.Equal(t, zeroex.ESOrderAdded, orderEvent.EndState, testCase.description)

		latestStoredBlock, err = database.GetLatestMiniHeader()
		require.NoError(t, err, testCase.description)
		newOrders, err := database.FindOrdersV4(nil)
		require.NoError(t, err, testCase.description)
		require.Len(t, newOrders, 1, testCase.description)
		expectedOrderState = orderState{
			hash:               expectedOrderHash,
			isRemoved:          false,
			isUnfillable:       false,
			fillableAmount:     signedOrder.TakerAmount,
			lastUpdated:        time.Now(),
			lastValidatedBlock: latestStoredBlock,
		}
		checkOrderState(t, expectedOrderState, newOrders[0])

		cancel()
		teardownSubTest(t)
	}
}

func TestOrderWatcherV4NoChange(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	for _, testCase := range []*struct {
		description     string
		addOrdersOpts   *types.AddOrdersOpts
		shouldBeRemoved bool
	}{
		{
			description:     "no change with empty configurations",
			addOrdersOpts:   &types.AddOrdersOpts{},
			shouldBeRemoved: true,
		},
		// NOTE(jalextowle): We use all of the configurations here since this test
		// doesn't rely on a particular configuration being set. This tests that
		// these configurations do not change behvior when there are no changes
		// to the orders validity.
		{
			description: "no change with all configurations",
			addOrdersOpts: &types.AddOrdersOpts{
				KeepCancelled:   true,
				KeepExpired:     true,
				KeepFullyFilled: true,
				KeepUnfunded:    true,
			},
			shouldBeRemoved: false,
		},
	} {
		teardownSubTest := setupSubTest(t)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		database, err := db.New(ctx, db.TestOptions())
		require.NoError(t, err)

		signedOrder := scenario.NewSignedTestOrderV4(t,
			orderopts.SetupMakerState(true),
			orderopts.MakerAssetData(scenario.ZRXAssetData),
			orderopts.TakerAssetData(scenario.WETHAssetData),
		)
		expectedOrderHash, err := signedOrder.ComputeOrderHash()
		require.NoError(t, err)
		blockWatcher, _ := setupOrderWatcherScenarioV4(ctx, t, database, signedOrder, testCase.addOrdersOpts)

		latestStoredBlock, err := database.GetLatestMiniHeader()
		require.NoError(t, err)
		orders, err := database.FindOrdersV4(nil)
		require.NoError(t, err)
		require.Len(t, orders, 1)
		expectedOrderState := orderState{
			hash:               expectedOrderHash,
			isRemoved:          false,
			isUnfillable:       false,
			fillableAmount:     signedOrder.TakerAmount,
			lastUpdated:        time.Now(),
			lastValidatedBlock: latestStoredBlock,
		}
		checkOrderState(t, expectedOrderState, orders[0])

		// Transfer more ZRX to makerAddress (doesn't impact the order)
		zrxCoinbase := constants.GanacheAccount0
		opts := &bind.TransactOpts{
			From:   zrxCoinbase,
			Signer: scenario.GetTestSignerFn(zrxCoinbase),
		}
		txn, err := zrx.Transfer(opts, signedOrder.Maker, signedOrder.MakerAmount)
		require.NoError(t, err)
		waitTxnSuccessfullyMined(t, ethClient, txn)

		err = blockWatcher.SyncToLatestBlock()
		require.NoError(t, err)

		// HACK(albrow): Normally we would wait for order events instead of sleeping here,
		// but in this case we don't *expect* any order events. Sleeping is a workaround.
		// We could potentially solve this by adding internal events inside of order watcher
		// that are only used for testing, but that would also incur some overhead.
		time.Sleep(processBlockSleepTime)

		latestStoredBlock, err = database.GetLatestMiniHeader()
		require.NoError(t, err)
		newOrders, err := database.FindOrdersV4(nil)
		require.NoError(t, err)
		require.Len(t, newOrders, 1)
		expectedOrderState = orderState{
			hash:               expectedOrderHash,
			isRemoved:          false,
			isUnfillable:       false,
			fillableAmount:     signedOrder.TakerAmount,
			lastUpdated:        time.Now(),
			lastValidatedBlock: latestStoredBlock,
		}
		checkOrderState(t, expectedOrderState, newOrders[0])

		cancel()
		teardownSubTest(t)
	}
}

func TestOrderWatcherV4WETHWithdrawAndDeposit(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}
	t.Skip("This test fails with a timeout. I'm not sure why, but we'll disable the test for now. It is likely an error in how test scenarios are run.")

	for _, testCase := range []*struct {
		description     string
		addOrdersOpts   *types.AddOrdersOpts
		shouldBeRemoved bool
	}{
		{
			description:     "should be removed with no configurations",
			addOrdersOpts:   &types.AddOrdersOpts{},
			shouldBeRemoved: true,
		},
		{
			description:     "should be kept with KeepUnfunded",
			addOrdersOpts:   &types.AddOrdersOpts{KeepUnfunded: true},
			shouldBeRemoved: false,
		},
	} {
		teardownSubTest := setupSubTest(t)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		database, err := db.New(ctx, db.TestOptions())
		require.NoError(t, err, testCase.description)

		signedOrder := scenario.NewSignedTestOrderV4(t,
			orderopts.SetupMakerState(true),
			orderopts.MakerAssetData(scenario.WETHAssetData),
			orderopts.TakerAssetData(scenario.ZRXAssetData),
		)
		expectedOrderHash, err := signedOrder.ComputeOrderHash()
		require.NoError(t, err, testCase.description)
		blockWatcher, orderEventsChan := setupOrderWatcherScenarioV4(ctx, t, database, signedOrder, testCase.addOrdersOpts)

		// Withdraw maker's WETH (i.e. decrease WETH balance)
		// HACK(fabio): For some reason the txn fails with "out of gas" error with the
		// estimated gas amount
		gasLimit := uint64(50000)
		opts := &bind.TransactOpts{
			From:     signedOrder.Maker,
			Signer:   scenario.GetTestSignerFn(signedOrder.Maker),
			GasLimit: gasLimit,
		}
		txn, err := weth.Withdraw(opts, signedOrder.MakerAmount)
		require.NoError(t, err, testCase.description)
		waitTxnSuccessfullyMined(t, ethClient, txn)

		err = blockWatcher.SyncToLatestBlock()
		require.NoError(t, err, testCase.description)
		// Timeout occurs in ths call
		orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
		require.Len(t, orderEvents, 1)
		orderEvent := orderEvents[0]
		assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState, testCase.description)

		latestStoredBlock, err := database.GetLatestMiniHeader()
		require.NoError(t, err, testCase.description)
		orders, err := database.FindOrdersV4(nil)
		require.NoError(t, err, testCase.description)
		require.Len(t, orders, 1, testCase.description)
		assert.Equal(t, orderEvent.OrderHash, orders[0].Hash, testCase.description)
		expectedOrderState := orderState{
			hash:               expectedOrderHash,
			isRemoved:          testCase.shouldBeRemoved,
			isUnfillable:       true,
			fillableAmount:     big.NewInt(0),
			lastUpdated:        time.Now(),
			lastValidatedBlock: latestStoredBlock,
		}
		checkOrderState(t, expectedOrderState, orders[0])

		// Deposit maker's ETH (i.e. increase WETH balance)
		opts = &bind.TransactOpts{
			From:   signedOrder.Maker,
			Signer: scenario.GetTestSignerFn(signedOrder.Maker),
			Value:  signedOrder.MakerAmount,
		}
		txn, err = weth.Deposit(opts)
		require.NoError(t, err, testCase.description)
		waitTxnSuccessfullyMined(t, ethClient, txn)

		err = blockWatcher.SyncToLatestBlock()
		require.NoError(t, err, testCase.description)

		orderEvents = <-orderEventsChan
		require.Len(t, orderEvents, 1)
		orderEvent = orderEvents[0]
		assert.Equal(t, zeroex.ESOrderAdded, orderEvent.EndState, testCase.description)

		latestStoredBlock, err = database.GetLatestMiniHeader()
		require.NoError(t, err, testCase.description)
		newOrders, err := database.FindOrdersV4(nil)
		require.NoError(t, err, testCase.description)
		require.Len(t, newOrders, 1, testCase.description)
		expectedOrderState = orderState{
			hash:               expectedOrderHash,
			isRemoved:          false,
			isUnfillable:       false,
			fillableAmount:     signedOrder.TakerAmount,
			lastUpdated:        time.Now(),
			lastValidatedBlock: latestStoredBlock,
		}
		checkOrderState(t, expectedOrderState, newOrders[0])

		cancel()
		teardownSubTest(t)
	}
}

func TestOrderWatcherV4Canceled(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	for _, testCase := range []*struct {
		description     string
		addOrdersOpts   *types.AddOrdersOpts
		shouldBeRemoved bool
	}{
		{
			description:     "should be removed with no configurations",
			addOrdersOpts:   &types.AddOrdersOpts{},
			shouldBeRemoved: true,
		},
		{
			description:     "should be kept with KeepCancelled",
			addOrdersOpts:   &types.AddOrdersOpts{KeepCancelled: true},
			shouldBeRemoved: false,
		},
	} {
		teardownSubTest := setupSubTest(t)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		database, err := db.New(ctx, db.TestOptions())
		require.NoError(t, err, testCase.description)

		signedOrder := scenario.NewSignedTestOrderV4(t, orderopts.SetupMakerState(true))
		expectedOrderHash, err := signedOrder.ComputeOrderHash()
		require.NoError(t, err, testCase.description)
		blockWatcher, orderEventsChan := setupOrderWatcherScenarioV4(ctx, t, database, signedOrder, testCase.addOrdersOpts)

		// Cancel order
		opts := &bind.TransactOpts{
			From:   signedOrder.Maker,
			Signer: scenario.GetTestSignerFn(signedOrder.Maker),
		}
		trimmedOrder := signedOrder.EthereumAbiLimitOrder()
		txn, err := exchangeV4.CancelLimitOrder(opts, trimmedOrder)
		require.NoError(t, err, testCase.description)
		waitTxnSuccessfullyMined(t, ethClient, txn)

		err = blockWatcher.SyncToLatestBlock()
		require.NoError(t, err, testCase.description)

		orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
		require.Len(t, orderEvents, 1)
		orderEvent := orderEvents[0]
		assert.Equal(t, zeroex.ESOrderCancelled, orderEvent.EndState, testCase.description)

		latestStoredBlock, err := database.GetLatestMiniHeader()
		require.NoError(t, err, testCase.description)
		orders, err := database.FindOrdersV4(nil)
		require.NoError(t, err, testCase.description)
		require.Len(t, orders, 1, testCase.description)
		expectedOrderState := orderState{
			hash:               expectedOrderHash,
			isRemoved:          testCase.shouldBeRemoved,
			isUnfillable:       true,
			fillableAmount:     big.NewInt(0),
			lastUpdated:        time.Now(),
			lastValidatedBlock: latestStoredBlock,
		}
		checkOrderState(t, expectedOrderState, orders[0])

		cancel()
		teardownSubTest(t)
	}
}

func TestOrderWatcherV4CancelUpTo(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	for _, testCase := range []*struct {
		description     string
		addOrdersOpts   *types.AddOrdersOpts
		shouldBeRemoved bool
	}{
		{
			description:     "should be removed with no configurations",
			addOrdersOpts:   &types.AddOrdersOpts{},
			shouldBeRemoved: true,
		},
		{
			description:     "should be kept with KeepCancelled",
			addOrdersOpts:   &types.AddOrdersOpts{KeepCancelled: true},
			shouldBeRemoved: false,
		},
	} {
		teardownSubTest := setupSubTest(t)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		database, err := db.New(ctx, db.TestOptions())
		require.NoError(t, err, testCase.description)

		signedOrder := scenario.NewSignedTestOrderV4(t, orderopts.SetupMakerState(true))
		expectedOrderHash, err := signedOrder.ComputeOrderHash()
		require.NoError(t, err, testCase.description)
		blockWatcher, orderEventsChan := setupOrderWatcherScenarioV4(ctx, t, database, signedOrder, testCase.addOrdersOpts)

		// Cancel order with epoch
		opts := &bind.TransactOpts{
			From:   signedOrder.Maker,
			Signer: scenario.GetTestSignerFn(signedOrder.Maker),
		}
		targetOrderEpoch := big.NewInt(0).Add(signedOrder.Salt, big.NewInt(1))
		txn, err := exchangeV4.CancelPairLimitOrders(opts,
			signedOrder.MakerToken,
			signedOrder.TakerToken,
			targetOrderEpoch)
		require.NoError(t, err, testCase.description)
		waitTxnSuccessfullyMined(t, ethClient, txn)

		err = blockWatcher.SyncToLatestBlock()
		require.NoError(t, err, testCase.description)

		orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
		require.Len(t, orderEvents, 1, testCase.description)
		orderEvent := orderEvents[0]
		assert.Equal(t, zeroex.ESOrderCancelled, orderEvent.EndState, testCase.description)

		latestStoredBlock, err := database.GetLatestMiniHeader()
		require.NoError(t, err, testCase.description)
		orders, err := database.FindOrdersV4(nil)
		require.NoError(t, err, testCase.description)
		require.Len(t, orders, 1, testCase.description)
		expectedOrderState := orderState{
			hash:               expectedOrderHash,
			isRemoved:          testCase.shouldBeRemoved,
			isUnfillable:       true,
			fillableAmount:     big.NewInt(0),
			lastUpdated:        time.Now(),
			lastValidatedBlock: latestStoredBlock,
		}
		checkOrderState(t, expectedOrderState, orders[0])

		cancel()
		teardownSubTest(t)
	}
}

func TestOrderWatcherV4ERC20Filled(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	for _, testCase := range []*struct {
		description     string
		addOrdersOpts   *types.AddOrdersOpts
		shouldBeRemoved bool
	}{
		{
			description:     "should be removed with no configurations",
			addOrdersOpts:   &types.AddOrdersOpts{},
			shouldBeRemoved: true,
		},
		{
			description:     "should be kept with KeepFullyFilled",
			addOrdersOpts:   &types.AddOrdersOpts{KeepFullyFilled: true},
			shouldBeRemoved: false,
		},
	} {
		teardownSubTest := setupSubTest(t)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		database, err := db.New(ctx, db.TestOptions())
		require.NoError(t, err, testCase.description)

		takerAddress := constants.GanacheAccount3
		signedOrder := scenario.NewSignedTestOrderV4(t,
			orderopts.SetupMakerState(true),
			orderopts.SetupTakerAddress(takerAddress),
		)
		expectedOrderHash, err := signedOrder.ComputeOrderHash()
		require.NoError(t, err, testCase.description)
		blockWatcher, orderEventsChan := setupOrderWatcherScenarioV4(ctx, t, database, signedOrder, testCase.addOrdersOpts)

		// Fill order
		opts := &bind.TransactOpts{
			From:   takerAddress,
			Signer: scenario.GetTestSignerFn(takerAddress),
			Value:  big.NewInt(100000000000000000),
		}
		trimmedOrder := signedOrder.EthereumAbiLimitOrder()
		signature := signedOrder.EthereumAbiSignature()
		txn, err := exchangeV4.FillLimitOrder(opts, trimmedOrder, signature, signedOrder.TakerAmount)
		require.NoError(t, err, testCase.description)
		waitTxnSuccessfullyMined(t, ethClient, txn)

		err = blockWatcher.SyncToLatestBlock()
		require.NoError(t, err, testCase.description)

		orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
		require.Len(t, orderEvents, 1, testCase.description)
		orderEvent := orderEvents[0]
		assert.Equal(t, zeroex.ESOrderFullyFilled, orderEvent.EndState, testCase.description)

		latestStoredBlock, err := database.GetLatestMiniHeader()
		require.NoError(t, err, testCase.description)
		orders, err := database.FindOrdersV4(nil)
		require.NoError(t, err, testCase.description)
		require.Len(t, orders, 1, testCase.description)
		expectedOrderState := orderState{
			hash:               expectedOrderHash,
			isRemoved:          testCase.shouldBeRemoved,
			isUnfillable:       true,
			fillableAmount:     big.NewInt(0),
			lastUpdated:        time.Now(),
			lastValidatedBlock: latestStoredBlock,
		}
		checkOrderState(t, expectedOrderState, orders[0])

		cancel()
		teardownSubTest(t)
	}
}

func TestOrderWatcherV4ERC20PartiallyFilled(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	for _, testCase := range []*struct {
		description   string
		addOrdersOpts *types.AddOrdersOpts
	}{
		{
			description:   "should be kept with no configurations",
			addOrdersOpts: &types.AddOrdersOpts{},
		},
		{
			description:   "should be kept with KeepFullyFilled",
			addOrdersOpts: &types.AddOrdersOpts{KeepFullyFilled: true},
		},
	} {
		teardownSubTest := setupSubTest(t)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		database, err := db.New(ctx, db.TestOptions())
		require.NoError(t, err, testCase.description)

		takerAddress := constants.GanacheAccount3
		signedOrder := scenario.NewSignedTestOrderV4(t,
			orderopts.SetupMakerState(true),
			orderopts.SetupTakerAddress(takerAddress),
		)
		expectedOrderHash, err := signedOrder.ComputeOrderHash()
		require.NoError(t, err, testCase.description)
		blockWatcher, orderEventsChan := setupOrderWatcherScenarioV4(ctx, t, database, signedOrder, testCase.addOrdersOpts)

		// Partially fill order
		opts := &bind.TransactOpts{
			From:   takerAddress,
			Signer: scenario.GetTestSignerFn(takerAddress),
			Value:  big.NewInt(100000000000000000),
		}
		trimmedOrder := signedOrder.EthereumAbiLimitOrder()
		signature := signedOrder.EthereumAbiSignature()
		halfAmount := new(big.Int).Div(signedOrder.TakerAmount, big.NewInt(2))
		txn, err := exchangeV4.FillLimitOrder(opts, trimmedOrder, signature, halfAmount)
		require.NoError(t, err, testCase.description)
		waitTxnSuccessfullyMined(t, ethClient, txn)

		err = blockWatcher.SyncToLatestBlock()
		require.NoError(t, err, testCase.description)

		orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
		require.Len(t, orderEvents, 1, testCase.description)
		orderEvent := orderEvents[0]
		assert.Equal(t, zeroex.ESOrderFilled, orderEvent.EndState, testCase.description)

		latestStoredBlock, err := database.GetLatestMiniHeader()
		require.NoError(t, err, testCase.description)
		orders, err := database.FindOrdersV4(nil)
		require.NoError(t, err, testCase.description)
		require.Len(t, orders, 1, testCase.description)
		expectedOrderState := orderState{
			hash:               expectedOrderHash,
			isRemoved:          false,
			isUnfillable:       false,
			isExpired:          false,
			fillableAmount:     halfAmount,
			lastUpdated:        time.Now(),
			lastValidatedBlock: latestStoredBlock,
		}
		checkOrderState(t, expectedOrderState, orders[0])

		cancel()
		teardownSubTest(t)
	}
}

func TestOrderWatcherOrderV4ExpiredThenUnexpired(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	for _, testCase := range []*struct {
		description     string
		addOrdersOpts   *types.AddOrdersOpts
		shouldBeRemoved bool
	}{
		{
			description:     "should be removed with no configurations",
			addOrdersOpts:   &types.AddOrdersOpts{},
			shouldBeRemoved: true,
		},
		{
			description:     "should be kept with KeepExpired",
			addOrdersOpts:   &types.AddOrdersOpts{KeepExpired: true},
			shouldBeRemoved: false,
		},
	} {
		// Set up test and orderWatcher
		teardownSubTest := setupSubTest(t)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		dbOptions := db.TestOptions()
		database, err := db.New(ctx, dbOptions)
		require.NoError(t, err, testCase.description)

		// Create and add an order (which will later become expired) to OrderWatcher
		expirationTime := time.Now().Add(24 * time.Hour)
		expirationTimeSeconds := big.NewInt(expirationTime.Unix())
		signedOrder := scenario.NewSignedTestOrderV4(t,
			orderopts.SetupMakerState(true),
			orderopts.ExpirationTimeSeconds(expirationTimeSeconds),
		)
		expectedOrderHash, err := signedOrder.ComputeOrderHash()
		require.NoError(t, err, testCase.description)
		blockwatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)
		watchOrderV4(ctx, t, orderWatcher, blockwatcher, signedOrder, false, testCase.addOrdersOpts)

		orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
		orderWatcher.Subscribe(orderEventsChan)

		// Simulate a block found with a timestamp past expirationTime
		latestBlock, err := database.GetLatestMiniHeader()
		require.NoError(t, err, testCase.description)
		nextBlock := &types.MiniHeader{
			Parent:    latestBlock.Hash,
			Hash:      common.HexToHash("0x1"),
			Number:    big.NewInt(0).Add(latestBlock.Number, big.NewInt(1)),
			Timestamp: expirationTime.Add(1 * time.Minute),
		}
		expiringBlockEvents := []*blockwatch.Event{
			{
				Type:        blockwatch.Added,
				BlockHeader: nextBlock,
			},
		}
		orderWatcher.blockEventsChan <- expiringBlockEvents

		// Await expired event
		orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
		require.Len(t, orderEvents, 1, testCase.description)
		orderEvent := orderEvents[0]
		assert.Equal(t, zeroex.ESOrderExpired, orderEvent.EndState, testCase.description)

		orders, err := database.FindOrdersV4(nil)
		require.NoError(t, err, testCase.description)
		require.Len(t, orders, 1, testCase.description)
		expectedOrderState := orderState{
			hash:               expectedOrderHash,
			isRemoved:          testCase.shouldBeRemoved,
			isUnfillable:       true,
			isExpired:          true,
			fillableAmount:     signedOrder.TakerAmount,
			lastUpdated:        time.Now(),
			lastValidatedBlock: nextBlock,
		}
		checkOrderState(t, expectedOrderState, orders[0])

		// Simulate a block re-org
		replacementBlockHash := common.HexToHash("0x2")
		reorgBlockEvents := []*blockwatch.Event{
			{
				Type:        blockwatch.Removed,
				BlockHeader: nextBlock,
			},
			{
				Type: blockwatch.Added,
				BlockHeader: &types.MiniHeader{
					Parent:    nextBlock.Parent,
					Hash:      replacementBlockHash,
					Number:    nextBlock.Number,
					Logs:      []ethtypes.Log{},
					Timestamp: expirationTime.Add(-2 * time.Hour),
				},
			},
			{
				Type: blockwatch.Added,
				BlockHeader: &types.MiniHeader{
					Parent:    replacementBlockHash,
					Hash:      common.HexToHash("0x3"),
					Number:    big.NewInt(0).Add(nextBlock.Number, big.NewInt(1)),
					Logs:      []ethtypes.Log{},
					Timestamp: expirationTime.Add(-1 * time.Hour),
				},
			},
		}
		orderWatcher.blockEventsChan <- reorgBlockEvents

		// Await unexpired event
		orderEvents = waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
		require.Len(t, orderEvents, 1, testCase.description)
		orderEvent = orderEvents[0]
		assert.Equal(t, zeroex.ESOrderUnexpired, orderEvent.EndState, testCase.description)

		newOrders, err := database.FindOrdersV4(nil)
		require.NoError(t, err, testCase.description)
		require.Len(t, newOrders, 1, testCase.description)
		expectedOrderState = orderState{
			hash:               expectedOrderHash,
			isRemoved:          false,
			isUnfillable:       false,
			isExpired:          false,
			fillableAmount:     signedOrder.TakerAmount,
			lastUpdated:        time.Now(),
			lastValidatedBlock: reorgBlockEvents[len(reorgBlockEvents)-1].BlockHeader,
		}
		checkOrderState(t, expectedOrderState, newOrders[0])

		cancel()
		teardownSubTest(t)
	}
}

func TestOrderWatcherOrderV4ExpiredWhenAddedThenUnexpired(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	// Set up test and orderWatcher
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	dbOptions := db.TestOptions()
	database, err := db.New(ctx, dbOptions)
	require.NoError(t, err)
	blockwatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)

	// Create an order which will be expired when added to the OrderWatcher
	expirationTime := time.Now().Add(-24 * time.Hour)
	expirationTimeSeconds := big.NewInt(expirationTime.Unix())
	signedOrder := scenario.NewSignedTestOrderV4(t,
		orderopts.SetupMakerState(true),
		orderopts.ExpirationTimeSeconds(expirationTimeSeconds),
	)
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)

	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// Add the order to Mesh
	err = blockwatcher.SyncToLatestBlock()
	require.NoError(t, err)
	validationResults, err := orderWatcher.ValidateAndStoreValidOrdersV4(ctx, []*zeroex.SignedOrderV4{signedOrder}, constants.TestChainID, false, &types.AddOrdersOpts{KeepExpired: true})
	require.NoError(t, err)

	assert.Len(t, validationResults.Accepted, 0)
	assert.Len(t, validationResults.Rejected, 1)
	assert.Equal(t, ordervalidator.ROExpired, validationResults.Rejected[0].Status)

	orders, err := database.FindOrdersV4(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	expectedValidationBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          false,
		isUnfillable:       true,
		isExpired:          true,
		fillableAmount:     big.NewInt(0),
		lastUpdated:        time.Now(),
		lastValidatedBlock: expectedValidationBlock,
	}
	checkOrderState(t, expectedOrderState, orders[0])

	// TODO(jalextowle): This code isn't needed with the current hacky test.
	// This could be improved by stubbing out validation.
	// Grep for SlowContractCaller in this file to see how the interface
	// works. The idea would be to create a contract caller that ensures that
	// DevUtils is called when the OrderWatcher is deciding whether or not to
	// unexpire the order.
	//
	// The reason why this can't be tested in a more straightforward
	// way is that it's non-trivial to implement a re-org at the level of ganache.
	// Using ganache to simulate the whole workflow would be the optimal testing
	// solution.
	//
	// Simulate a block re-org
	// replacementBlockHash := common.HexToHash("0x2")
	// reorgBlockEvents := []*blockwatch.Event{
	// 	{
	// 		Type:        blockwatch.Removed,
	// 		BlockHeader: expectedValidationBlock,
	// 	},
	// 	{
	// 		Type: blockwatch.Added,
	// 		BlockHeader: &types.MiniHeader{
	// 			Parent:    expectedValidationBlock.Parent,
	// 			Hash:      replacementBlockHash,
	// 			Number:    expectedValidationBlock.Number,
	// 			Logs:      []ethtypes.Log{},
	// 			Timestamp: expirationTime.Add(-2 * time.Hour),
	// 		},
	// 	},
	// 	{
	// 		Type: blockwatch.Added,
	// 		BlockHeader: &types.MiniHeader{
	// 			Parent:    replacementBlockHash,
	// 			Hash:      common.HexToHash("0x3"),
	// 			Number:    big.NewInt(0).Add(expectedValidationBlock.Number, big.NewInt(1)),
	// 			Logs:      []ethtypes.Log{},
	// 			Timestamp: expirationTime.Add(-1 * time.Hour),
	// 		},
	// 	},
	// }
	// orderWatcher.blockEventsChan <- reorgBlockEvents

	// HACK(jalextowle): The block events above don't actually exist in the Ganache
	// environment that Mesh uses to validate orders. This means that Mesh won't be
	// able to actually validate the order.
	possiblyUnexpiredOrders, err := orderWatcher.findOrdersToPossiblyUnexpire(expirationTime.Add(-2 * time.Hour))
	require.NoError(t, err)
	assert.Len(t, possiblyUnexpiredOrders, 1)
}

// NOTE(jalextowle): We don't need to implement a test for this with configurations
// as the configurations do not interact with the pinning system.
func TestOrderWatcherV4DecreaseExpirationTime(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	// Set up test and orderWatcher. Manually change maxOrders.
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	maxOrders := 10
	dbOpts := db.TestOptions()
	dbOpts.MaxOrders = maxOrders
	database, err := db.New(ctx, dbOpts)
	require.NoError(t, err)

	blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)
	orderWatcher.maxOrders = maxOrders

	// Create and watch maxOrders orders. Each order has a different expiration time.
	optionsForIndex := func(index int) []orderopts.Option {
		expirationTime := time.Now().Add(10*time.Minute + time.Duration(index)*time.Minute)
		expirationTimeSeconds := big.NewInt(expirationTime.Unix())
		return []orderopts.Option{
			orderopts.SetupMakerState(true),
			orderopts.ExpirationTimeSeconds(expirationTimeSeconds),
		}
	}
	signedOrders := scenario.NewSignedTestOrdersBatchV4(t, maxOrders, optionsForIndex)
	for _, signedOrder := range signedOrders {
		watchOrderV4(ctx, t, orderWatcher, blockWatcher, signedOrder, false, &types.AddOrdersOpts{})
	}

	// We don't care about the order events above for the purposes of this test,
	// so we only subscribe now.
	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// The next order should cause some orders to be removed and the appropriate
	// events to fire.
	expirationTime := time.Now().Add(10*time.Minute + 1*time.Second)
	expirationTimeSeconds := big.NewInt(expirationTime.Unix())
	signedOrder := scenario.NewSignedTestOrderV4(t,
		orderopts.SetupMakerState(true),
		orderopts.ExpirationTimeSeconds(expirationTimeSeconds),
	)
	watchOrderV4(ctx, t, orderWatcher, blockWatcher, signedOrder, false, &types.AddOrdersOpts{})
	expectedOrderEvents := 2
	orderEvents := waitForOrderEvents(t, orderEventsChan, expectedOrderEvents, 4*time.Second)
	require.Len(t, orderEvents, expectedOrderEvents, "wrong number of order events were fired")

	storedMaxExpirationTime, err := database.GetCurrentMaxExpirationTime()
	require.NoError(t, err)

	// One event should be STOPPED_WATCHING. The other event should be ADDED.
	// The order in which the events are emitted is not guaranteed.
	numAdded := 0
	numStoppedWatching := 0
	for _, orderEvent := range orderEvents {
		switch orderEvent.EndState {
		case zeroex.ESOrderAdded:
			numAdded += 1
			var orderExpirationTime *big.Int
			if orderEvent.SignedOrder != nil {
				orderExpirationTime = orderEvent.SignedOrder.ExpirationTimeSeconds
			} else {
				orderExpirationTime = orderEvent.SignedOrderV4.Expiry
			}
			assert.True(t, orderExpirationTime.Cmp(storedMaxExpirationTime) == -1, "ADDED order has an expiration time of %s which is *greater than* the maximum of %s", orderExpirationTime, storedMaxExpirationTime)
		case zeroex.ESStoppedWatching:
			numStoppedWatching += 1
			var orderExpirationTime *big.Int
			if orderEvent.SignedOrder != nil {
				orderExpirationTime = orderEvent.SignedOrder.ExpirationTimeSeconds
			} else {
				orderExpirationTime = orderEvent.SignedOrderV4.Expiry
			}
			assert.True(t, orderExpirationTime.Cmp(storedMaxExpirationTime) != -1, "STOPPED_WATCHING order has an expiration time of %s which is *less than* the maximum of %s", orderExpirationTime, storedMaxExpirationTime)
		default:
			t.Errorf("unexpected order event type: %s", orderEvent.EndState)
		}
	}
	assert.Equal(t, 1, numAdded, "wrong number of ADDED events")
	assert.Equal(t, 1, numStoppedWatching, "wrong number of STOPPED_WATCHING events")

	// Now we check that the correct number of orders remain and that all
	// remaining orders have an expiration time less than the current max.
	expectedRemainingOrders := orderWatcher.maxOrders
	remainingOrders, err := database.FindOrdersV4(nil)
	require.NoError(t, err)
	require.Len(t, remainingOrders, expectedRemainingOrders)
	for _, order := range remainingOrders {
		var expiry *big.Int
		if order.OrderV3 != nil {
			expiry = order.OrderV3.ExpirationTimeSeconds
		}
		if order.OrderV4 != nil {
			expiry = order.OrderV4.Expiry
		}
		assert.True(t, expiry.Cmp(storedMaxExpirationTime) != 1, "remaining order has an expiration time of %s which is *greater than* the maximum of %s", expiry, storedMaxExpirationTime)
	}

	// Confirm that a pinned order will be accepted even if its expiration
	// is greater than the current max.
	pinnedOrder := scenario.NewSignedTestOrderV4(t,
		orderopts.SetupMakerState(true),
		orderopts.ExpirationTimeSeconds(big.NewInt(0).Add(storedMaxExpirationTime, big.NewInt(10))),
	)
	pinnedOrderHash, err := pinnedOrder.ComputeOrderHash()
	require.NoError(t, err)
	watchOrderV4(ctx, t, orderWatcher, blockWatcher, pinnedOrder, true, &types.AddOrdersOpts{})

	expectedOrderEvents = 2
	orderEvents = waitForOrderEvents(t, orderEventsChan, expectedOrderEvents, 4*time.Second)
	require.Len(t, orderEvents, expectedOrderEvents, "wrong number of order events were fired")

	// One event should be STOPPED_WATCHING. The other event should be ADDED.
	// The order in which the events are emitted is not guaranteed.
	numAdded = 0
	numStoppedWatching = 0
	for _, orderEvent := range orderEvents {
		switch orderEvent.EndState {
		case zeroex.ESOrderAdded:
			numAdded += 1
			assert.Equal(t, pinnedOrderHash.Hex(), orderEvent.OrderHash.Hex(), "ADDED event had wrong order hash")
		case zeroex.ESStoppedWatching:
			numStoppedWatching += 1
		default:
			t.Errorf("unexpected order event type: %s", orderEvent.EndState)
		}
	}
	assert.Equal(t, 1, numAdded, "wrong number of ADDED events")
	assert.Equal(t, 1, numStoppedWatching, "wrong number of STOPPED_WATCHING events")
}

func TestOrderWatcherV4BatchEmitsAddedEvents(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)

	// Subscribe to OrderWatcher
	orderEventsChan := make(chan []*zeroex.OrderEvent, 10)
	orderWatcher.Subscribe(orderEventsChan)

	// Create numOrders test orders in a batch.
	numOrders := 2
	orderOptions := scenario.OptionsForAll(orderopts.SetupMakerState(true))
	signedOrders := scenario.NewSignedTestOrdersBatchV4(t, numOrders, orderOptions)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	validationResults, err := orderWatcher.ValidateAndStoreValidOrdersV4(ctx, signedOrders, constants.TestChainID, false, &types.AddOrdersOpts{})
	require.Len(t, validationResults.Rejected, 0)
	require.NoError(t, err)

	orderEvents := <-orderEventsChan
	require.Len(t, orderEvents, numOrders)
	for _, orderEvent := range orderEvents {
		assert.Equal(t, zeroex.ESOrderAdded, orderEvent.EndState)
	}

	orders, err := database.FindOrdersV4(nil)
	require.NoError(t, err)
	require.Len(t, orders, numOrders)
}

func TestOrderWatcherV4Cleanup(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)
	blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)

	// Create and add two orders to OrderWatcher
	orderOptions := scenario.OptionsForAll(orderopts.SetupMakerState(true))
	signedOrders := scenario.NewSignedTestOrdersBatchV4(t, 2, orderOptions)
	signedOrderOne := signedOrders[0]
	watchOrderV4(ctx, t, orderWatcher, blockWatcher, signedOrderOne, false, &types.AddOrdersOpts{})
	signedOrderTwo := signedOrders[1]
	watchOrderV4(ctx, t, orderWatcher, blockWatcher, signedOrderTwo, false, &types.AddOrdersOpts{})
	signedOrderOneHash, err := signedOrderTwo.ComputeOrderHash()
	require.NoError(t, err)

	// Set lastUpdate for signedOrderOne to more than defaultLastUpdatedBuffer so that signedOrderOne
	// does not get re-validated by the cleanup job
	err = database.UpdateOrder(signedOrderOneHash, func(orderToUpdate *types.OrderWithMetadata) (*types.OrderWithMetadata, error) {
		orderToUpdate.LastUpdated = time.Now().Add(-defaultLastUpdatedBuffer - 1*time.Minute)
		return orderToUpdate, nil
	})
	require.NoError(t, err)

	// Subscribe to OrderWatcher
	orderEventsChan := make(chan []*zeroex.OrderEvent, 10)
	orderWatcher.Subscribe(orderEventsChan)

	// Since no state changes occurred without corresponding events being emitted, we expect
	// cleanup not to result in any new events
	err = orderWatcher.Cleanup(ctx, defaultLastUpdatedBuffer)
	require.NoError(t, err)

	select {
	case <-orderEventsChan:
		t.Error("Expected no orderEvents to fire after calling Cleanup()")
	case <-time.After(100 * time.Millisecond):
		// Noop
	}
}

func TestOrderWatcherV4HandleOrderExpirationsExpired(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	for _, testCase := range []*struct {
		description     string
		addOrdersOpts   *types.AddOrdersOpts
		shouldBeRemoved bool
	}{
		{
			description:     "should be removed with no configurations",
			addOrdersOpts:   &types.AddOrdersOpts{},
			shouldBeRemoved: true,
		},
		{
			description:     "should be kept with KeepExpired",
			addOrdersOpts:   &types.AddOrdersOpts{KeepExpired: true},
			shouldBeRemoved: false,
		},
	} {
		// Set up test and orderWatcher
		teardownSubTest := setupSubTest(t)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		database, err := db.New(ctx, db.TestOptions())
		require.NoError(t, err, testCase.description)

		// Create and add an order (which will later become expired) to OrderWatcher
		expirationTime := time.Now().Add(24 * time.Hour)
		expirationTimeSeconds := big.NewInt(expirationTime.Unix())
		orderOptions := scenario.OptionsForAll(
			orderopts.SetupMakerState(true),
			orderopts.ExpirationTimeSeconds(expirationTimeSeconds),
		)
		signedOrders := scenario.NewSignedTestOrdersBatchV4(t, 2, orderOptions)
		signedOrderOne := signedOrders[0]
		signedOrderTwo := signedOrders[1]
		blockwatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)
		watchOrderV4(ctx, t, orderWatcher, blockwatcher, signedOrderOne, false, testCase.addOrdersOpts)
		watchOrderV4(ctx, t, orderWatcher, blockwatcher, signedOrderTwo, false, testCase.addOrdersOpts)

		signedOrderOneHash, err := signedOrderOne.ComputeOrderHash()
		require.NoError(t, err, testCase.description)
		orderOne, err := database.GetOrder(signedOrderOneHash)
		require.NoError(t, err, testCase.description)
		// Since we flag SignedOrderOne for revalidation, we expect `handleOrderExpirations` not to return an
		// expiry event for it.
		ordersToRevalidate := map[common.Hash]*types.OrderWithMetadata{
			signedOrderOneHash: orderOne,
		}

		// Make a "fake" block with a timestamp 1 second after expirationTime.
		latestBlock, err := database.GetLatestMiniHeader()
		require.NoError(t, err, testCase.description)
		latestBlock.Timestamp = expirationTime.Add(1 * time.Second)
		orderEvents, _, err := orderWatcher.handleOrderExpirations(latestBlock, ordersToRevalidate)
		require.NoError(t, err, testCase.description)

		require.Len(t, orderEvents, 1)
		orderEvent := orderEvents[0]
		signedOrderTwoHash, err := signedOrderTwo.ComputeOrderHash()
		require.NoError(t, err, testCase.description)
		assert.Equal(t, signedOrderTwoHash, orderEvent.OrderHash, testCase.description)
		assert.Equal(t, zeroex.ESOrderExpired, orderEvent.EndState, testCase.description)
		assert.Equal(t, big.NewInt(0), orderEvent.FillableTakerAssetAmount, testCase.description)
		assert.Len(t, orderEvent.ContractEvents, 0, testCase.description)

		orderTwo, err := database.GetOrder(signedOrderTwoHash)
		require.NoError(t, err, testCase.description)
		assert.Equal(t, testCase.shouldBeRemoved, orderTwo.IsRemoved, testCase.description)
		assert.Equal(t, true, orderTwo.IsUnfillable, testCase.description)
		assert.Equal(t, true, orderTwo.IsExpired, testCase.description)

		cancel()
		teardownSubTest(t)
	}
}

// NOTE(jalextowle): We don't need a corresponding test with configurations
// because this test does not test for removal.
func TestOrderWatcherV4HandleOrderExpirationsUnexpired(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	// Set up test and orderWatcher
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	// Create and add an order (which will later become expired) to OrderWatcher
	expirationTime := time.Now().Add(24 * time.Hour)
	expirationTimeSeconds := big.NewInt(expirationTime.Unix())
	orderOptions := scenario.OptionsForAll(
		orderopts.SetupMakerState(true),
		orderopts.ExpirationTimeSeconds(expirationTimeSeconds),
	)
	signedOrders := scenario.NewSignedTestOrdersBatchV4(t, 2, orderOptions)
	signedOrderOne := signedOrders[0]
	signedOrderTwo := signedOrders[1]
	blockwatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)
	watchOrderV4(ctx, t, orderWatcher, blockwatcher, signedOrderOne, false, &types.AddOrdersOpts{})
	watchOrderV4(ctx, t, orderWatcher, blockwatcher, signedOrderTwo, false, &types.AddOrdersOpts{})

	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// Simulate a block found with a timestamp past expirationTime
	latestBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	blockTimestamp := expirationTime.Add(1 * time.Minute)
	nextBlock := &types.MiniHeader{
		Parent:    latestBlock.Hash,
		Hash:      common.HexToHash("0x1"),
		Number:    big.NewInt(0).Add(latestBlock.Number, big.NewInt(1)),
		Timestamp: blockTimestamp,
	}
	expiringBlockEvents := []*blockwatch.Event{
		{
			Type:        blockwatch.Added,
			BlockHeader: nextBlock,
		},
	}
	orderWatcher.blockEventsChan <- expiringBlockEvents

	// Await expired event
	orderEvents := waitForOrderEvents(t, orderEventsChan, 2, 4*time.Second)
	require.Len(t, orderEvents, 2)
	for _, orderEvent := range orderEvents {
		assert.Equal(t, zeroex.ESOrderExpired, orderEvent.EndState)
	}

	signedOrderOneHash, err := signedOrderOne.ComputeOrderHash()
	require.NoError(t, err)
	orderOne, err := database.GetOrder(signedOrderOneHash)
	require.NoError(t, err)
	// Since we flag SignedOrderOne for revalidation, we expect `handleOrderExpirations` not to return an
	// unexpiry event for it.
	ordersToRevalidate := map[common.Hash]*types.OrderWithMetadata{
		signedOrderOneHash: orderOne,
	}

	// Make a "fake" block with a timestamp 1 minute before expirationTime. This simulates
	// block-reorg where new latest block has an earlier timestamp than the last
	latestBlock, err = database.GetLatestMiniHeader()
	require.NoError(t, err)
	latestBlock.Timestamp = expirationTime.Add(-1 * time.Minute)
	orderEvents, _, err = orderWatcher.handleOrderExpirations(latestBlock, ordersToRevalidate)
	require.NoError(t, err)

	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	signedOrderTwoHash, err := signedOrderTwo.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, signedOrderTwoHash, orderEvent.OrderHash)
	assert.Equal(t, zeroex.ESOrderUnexpired, orderEvent.EndState)
	assert.Equal(t, signedOrderTwo.TakerAmount, orderEvent.FillableTakerAssetAmount)
	assert.Len(t, orderEvent.ContractEvents, 0)

	orderTwo, err := database.GetOrder(signedOrderTwoHash)
	require.NoError(t, err)
	assert.Equal(t, false, orderTwo.IsRemoved)
	assert.Equal(t, false, orderTwo.IsUnfillable)
	assert.Equal(t, false, orderTwo.IsExpired)
}

// Scenario: Order has become unexpired and filled in the same block events processed. We test this case using
// `convertValidationResultsIntoOrderEvents` since we cannot properly time-travel using Ganache.
// Source: https://github.com/trufflesuite/ganache-cli/issues/708
func TestConvertValidationResultsIntoOrderV4EventsUnexpired(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	// Set up test and orderWatcher
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	// Create and add an order (which will later become expired) to OrderWatcher
	expirationTime := time.Now().Add(24 * time.Hour)
	expirationTimeSeconds := big.NewInt(expirationTime.Unix())
	signedOrder := scenario.NewSignedTestOrderV4(t,
		orderopts.SetupMakerState(true),
		orderopts.ExpirationTimeSeconds(expirationTimeSeconds),
	)
	blockwatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)
	watchOrderV4(ctx, t, orderWatcher, blockwatcher, signedOrder, false, &types.AddOrdersOpts{})

	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// Simulate a block found with a timestamp past expirationTime. This will mark the order as removed
	// and will remove it from the expiration watcher.
	latestBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	blockTimestamp := expirationTime.Add(1 * time.Minute)
	nextBlock := &types.MiniHeader{
		Parent:    latestBlock.Hash,
		Hash:      common.HexToHash("0x1"),
		Number:    big.NewInt(0).Add(latestBlock.Number, big.NewInt(1)),
		Timestamp: blockTimestamp,
	}
	expiringBlockEvents := []*blockwatch.Event{
		{
			Type:        blockwatch.Added,
			BlockHeader: nextBlock,
		},
	}
	orderWatcher.blockEventsChan <- expiringBlockEvents

	// Await expired event
	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	assert.Equal(t, zeroex.ESOrderExpired, orderEvents[0].EndState)

	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	orderOne, err := database.GetOrder(orderHash)
	require.NoError(t, err)

	validationResults := ordervalidator.ValidationResults{
		Accepted: []*ordervalidator.AcceptedOrderInfo{
			{
				OrderHash:                orderHash,
				SignedOrderV4:            signedOrder,
				FillableTakerAssetAmount: big.NewInt(1).Div(signedOrder.TakerAmount, big.NewInt(2)),
				IsNew:                    false,
			},
		},
		Rejected: []*ordervalidator.RejectedOrderInfo{},
	}
	orderHashToDBOrder := map[common.Hash]*types.OrderWithMetadata{
		orderHash: orderOne,
	}
	exchangeFillEvent := "ExchangeFillEvent"
	orderHashToEvents := map[common.Hash][]*zeroex.ContractEvent{
		orderHash: {
			&zeroex.ContractEvent{
				Kind: exchangeFillEvent,
			},
		},
	}
	// Make a "fake" block with a timestamp 1 minute before expirationTime. This simulates
	// block-reorg where new latest block has an earlier timestamp than the last
	validationBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	validationBlock.Timestamp = expirationTime.Add(-1 * time.Minute)
	orderEvents, err = orderWatcher.convertValidationResultsIntoOrderEvents(&validationResults, orderHashToDBOrder, orderHashToEvents, map[common.Hash]struct{}{}, validationBlock)
	require.NoError(t, err)

	require.Len(t, orderEvents, 2)
	orderEventTwo := orderEvents[0]
	assert.Equal(t, orderHash, orderEventTwo.OrderHash)
	assert.Equal(t, zeroex.ESOrderUnexpired, orderEventTwo.EndState)
	assert.Len(t, orderEventTwo.ContractEvents, 0)
	orderEventOne := orderEvents[1]
	assert.Equal(t, orderHash, orderEventOne.OrderHash)
	assert.Equal(t, zeroex.ESOrderFilled, orderEventOne.EndState)
	assert.Len(t, orderEventOne.ContractEvents, 1)
	assert.Equal(t, orderEventOne.ContractEvents[0].Kind, exchangeFillEvent)

	existingOrder, err := database.GetOrder(orderHash)
	require.NoError(t, err)
	assert.Equal(t, false, existingOrder.IsRemoved)
	assert.Equal(t, false, existingOrder.IsUnfillable)
	assert.Equal(t, false, existingOrder.IsExpired)
}

func TestRevalidateOrdersV4ForMissingEvents(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	// Set up test and orderWatcher
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)
	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// Create a new order
	signedOrder := scenario.NewSignedTestOrderV4(t, orderopts.SetupMakerState(true))
	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)
	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)

	// Cancel the order
	opts := &bind.TransactOpts{
		From:   signedOrder.Maker,
		Signer: scenario.GetTestSignerFn(signedOrder.Maker),
	}
	trimmedOrder := signedOrder.EthereumAbiLimitOrder()
	txn, err := exchangeV4.CancelLimitOrder(opts, trimmedOrder)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	validationResultsChan := make(chan *ordervalidator.ValidationResults, 1)
	g, innerCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		// NOTE(jalextowle): Sleep to allow the call to ValidateAndStoreValidOrders
		// to begin before syncing to latest block.
		time.Sleep(time.Second)
		err := blockWatcher.SyncToLatestBlock()
		return err
	})
	g.Go(func() error {
		validationResults, err := orderWatcher.ValidateAndStoreValidOrdersV4(innerCtx, []*zeroex.SignedOrderV4{signedOrder}, constants.TestChainID, false, &types.AddOrdersOpts{})
		if err != nil {
			return err
		}
		validationResultsChan <- validationResults
		return nil
	})
	err = g.Wait()
	require.NoError(t, err)

	select {
	case validationResults := <-validationResultsChan:
		require.Equal(t, len(validationResults.Accepted), 1)
		assert.Equal(t, len(validationResults.Rejected), 0)
		orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
		assert.Equal(t, zeroex.ESOrderAdded, orderEvents[0].EndState)
		assert.Equal(t, orderHash, orderEvents[0].OrderHash)
	default:
		t.Fatal("No validation results received")
	}

	err = orderWatcher.RevalidateOrdersForMissingEvents(ctx)
	require.NoError(t, err)
	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	assert.Equal(t, zeroex.ESOrderCancelled, orderEvents[0].EndState)
	assert.Equal(t, orderHash, orderEvents[0].OrderHash)
}

// TestMissingOrderEventsWithMissingBlocks tests that the orderwatcher will not
// miss block events for orders that were originally validated in a block that
// currently exists in the database.
func TestMissingOrderV4Events(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	// Set up test and orderWatcher
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	// TODO(jalextowle): This test will fail with "context canceled" if a context
	// with a timeout is used here.
	ctx := context.Background()
	dbOpts := db.TestOptions()
	database, err := db.New(ctx, dbOpts)
	require.NoError(t, err)

	validator, err := ordervalidator.New(
		&SlowContractCaller{
			caller:            ethRPCClient,
			contractCallDelay: time.Second,
		},
		constants.TestChainID,
		ethereumRPCMaxContentLength,
		ganacheAddresses,
	)
	require.NoError(t, err)

	blockWatcher, orderWatcher := setupOrderWatcherWithValidator(ctx, t, ethRPCClient, database, dbOpts.MaxMiniHeaders, validator)
	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// Create a new order
	signedOrder := scenario.NewSignedTestOrderV4(t, orderopts.SetupMakerState(true))
	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)
	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)

	// Cancel the order
	opts := &bind.TransactOpts{
		From:   signedOrder.Maker,
		Signer: scenario.GetTestSignerFn(signedOrder.Maker),
	}
	trimmedOrder := signedOrder.EthereumAbiLimitOrder()
	txn, err := exchangeV4.CancelLimitOrder(opts, trimmedOrder)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	validationResultsChan := make(chan *ordervalidator.ValidationResults, 1)
	g, innerCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		// NOTE(jalextowle): Sleep to allow the call to ValidateAndStoreValidOrders
		// to begin before syncing to latest block.
		time.Sleep(time.Second)
		err := blockWatcher.SyncToLatestBlock()
		return err
	})
	g.Go(func() error {
		validationResults, err := orderWatcher.ValidateAndStoreValidOrdersV4(innerCtx, []*zeroex.SignedOrderV4{signedOrder}, constants.TestChainID, false, &types.AddOrdersOpts{})
		if err != nil {
			return err
		}
		validationResultsChan <- validationResults
		return nil
	})
	err = g.Wait()
	require.NoError(t, err)

	select {
	case validationResults := <-validationResultsChan:
		require.Equal(t, len(validationResults.Accepted), 1)
		assert.Equal(t, len(validationResults.Rejected), 0)
		orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
		assert.Equal(t, zeroex.ESOrderAdded, orderEvents[0].EndState)
		assert.Equal(t, orderHash, orderEvents[0].OrderHash)
	default:
		t.Fatal("No validation results received")
	}

	// Add new block events and then check to see if the order has been removed from the blockwatcher
	latestBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	nextBlock := &types.MiniHeader{
		Parent:    latestBlock.Hash,
		Hash:      common.HexToHash("0x1"),
		Number:    big.NewInt(0).Add(latestBlock.Number, big.NewInt(1)),
		Timestamp: latestBlock.Timestamp.Add(15 * time.Second),
	}
	newBlockEvents := []*blockwatch.Event{
		{
			Type:        blockwatch.Added,
			BlockHeader: nextBlock,
		},
	}
	orderWatcher.blockEventsChan <- newBlockEvents

	// Await canceled event
	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 10*time.Second)
	assert.Equal(t, zeroex.ESOrderCancelled, orderEvents[0].EndState)
	assert.Equal(t, orderHash, orderEvents[0].OrderHash)
}

// TestMissingOrderEventsWithMissingBlocks tests that the orderwatcher will not
// miss block events for orders that were originally validated in a block that no
// longer exists in the database. This covers an edge case where the blockwatcher
// had to catch up a significant number of blocks during a previous call to
// `handleBlockEvents`.
// TODO(jalextowle): De-duplicate the code in this test and the above test
func TestMissingOrderV4EventsWithMissingBlocks(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	// Set up test and orderWatcher
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	ctx := context.Background()
	dbOpts := db.TestOptions()
	dbOpts.MaxMiniHeaders = 1
	database, err := db.New(ctx, dbOpts)
	require.NoError(t, err)

	validator, err := ordervalidator.New(
		&SlowContractCaller{
			caller:            ethRPCClient,
			contractCallDelay: time.Second,
		},
		constants.TestChainID,
		ethereumRPCMaxContentLength,
		ganacheAddresses,
	)
	require.NoError(t, err)

	blockWatcher, orderWatcher := setupOrderWatcherWithValidator(ctx, t, ethRPCClient, database, dbOpts.MaxMiniHeaders, validator)
	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// Create a new order
	signedOrder := scenario.NewSignedTestOrderV4(t, orderopts.SetupMakerState(true))
	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)
	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)

	// Cancel the order
	opts := &bind.TransactOpts{
		From:   signedOrder.Maker,
		Signer: scenario.GetTestSignerFn(signedOrder.Maker),
	}
	trimmedOrder := signedOrder.EthereumAbiLimitOrder()
	txn, err := exchangeV4.CancelLimitOrder(opts, trimmedOrder)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	// Cancel a new order to remove old miniheaders from the database.
	dummyOrder := scenario.NewSignedTestOrderV4(t, orderopts.SetupMakerState(true))
	opts = &bind.TransactOpts{
		From:   dummyOrder.Maker,
		Signer: scenario.GetTestSignerFn(dummyOrder.Maker),
	}
	trimmedOrder = dummyOrder.EthereumAbiLimitOrder()
	txn, err = exchangeV4.CancelLimitOrder(opts, trimmedOrder)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	validationResultsChan := make(chan *ordervalidator.ValidationResults, 1)
	g, innerCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		// NOTE(jalextowle): Sleep to allow the call to ValidateAndStoreValidOrders
		// to begin before syncing to latest block.
		time.Sleep(time.Second)
		err := blockWatcher.SyncToLatestBlock()
		return err
	})
	g.Go(func() error {
		validationResults, err := orderWatcher.ValidateAndStoreValidOrdersV4(innerCtx, []*zeroex.SignedOrderV4{signedOrder}, constants.TestChainID, false, &types.AddOrdersOpts{})
		if err != nil {
			return err
		}
		validationResultsChan <- validationResults
		return nil
	})
	err = g.Wait()
	require.NoError(t, err)

	select {
	case validationResults := <-validationResultsChan:
		require.Equal(t, len(validationResults.Accepted), 1)
		assert.Equal(t, len(validationResults.Rejected), 0)
		orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
		assert.Equal(t, zeroex.ESOrderAdded, orderEvents[0].EndState)
		assert.Equal(t, orderHash, orderEvents[0].OrderHash)
	default:
		t.Fatal("No validation results received")
	}

	// Add new block events and then check to see if the order has been removed from the blockwatcher
	latestBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	nextBlock := &types.MiniHeader{
		Parent:    latestBlock.Hash,
		Hash:      common.HexToHash("0x1"),
		Number:    big.NewInt(0).Add(latestBlock.Number, big.NewInt(1)),
		Timestamp: latestBlock.Timestamp.Add(15 * time.Second),
	}
	newBlockEvents := []*blockwatch.Event{
		{
			Type:        blockwatch.Added,
			BlockHeader: nextBlock,
		},
	}
	orderWatcher.blockEventsChan <- newBlockEvents

	// Await canceled event
	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 10*time.Second)
	assert.Equal(t, zeroex.ESOrderCancelled, orderEvents[0].EndState)
	assert.Equal(t, orderHash, orderEvents[0].OrderHash)
}

func setupOrderWatcherScenarioV4(ctx context.Context, t *testing.T, database *db.DB, signedOrder *zeroex.SignedOrderV4, opts *types.AddOrdersOpts) (*blockwatch.Watcher, chan []*zeroex.OrderEvent) {
	blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)

	// Start watching an order
	watchOrderV4(ctx, t, orderWatcher, blockWatcher, signedOrder, false, opts)

	// Subscribe to OrderWatcher
	orderEventsChan := make(chan []*zeroex.OrderEvent, 10)
	orderWatcher.Subscribe(orderEventsChan)

	return blockWatcher, orderEventsChan
}

func watchOrderV4(ctx context.Context, t *testing.T, orderWatcher *Watcher, blockWatcher *blockwatch.Watcher, signedOrder *zeroex.SignedOrderV4, pinned bool, opts *types.AddOrdersOpts) {
	err := blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	validationResults, err := orderWatcher.ValidateAndStoreValidOrdersV4(ctx, []*zeroex.SignedOrderV4{signedOrder}, constants.TestChainID, pinned, opts)
	require.NoError(t, err)
	if len(validationResults.Rejected) != 0 {
		spew.Dump(validationResults.Rejected)
	}
	require.Len(t, validationResults.Accepted, 1, "Expected order to pass validation and get added to OrderWatcher")
}
