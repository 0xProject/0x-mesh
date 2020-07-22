// +build !js

package orderwatch

import (
	"context"
	"flag"
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/common/types"
	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/db"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/blockwatch"
	"github.com/0xProject/0x-mesh/ethereum/ethrpcclient"
	"github.com/0xProject/0x-mesh/ethereum/ratelimit"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/scenario/orderopts"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
	"github.com/davecgh/go-spew/spew"
	geth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	blockRetentionLimit         = 20
	ethereumRPCRequestTimeout   = 30 * time.Second
	miniHeaderRetentionLimit    = 2
	blockPollingInterval        = 1 * time.Second
	ethereumRPCMaxContentLength = 524288
	maxEthRPCRequestsPer24HrUTC = 1000000
	maxEthRPCRequestsPerSeconds = 1000.0

	// processBlockSleepTime is the amount of time ot wait for order watcher to
	// process block events. If possible, we should listen for order events instead
	// of sleeping, but we need to use this in some places where we don't expect
	// any order events.
	processBlockSleepTime = 350 * time.Millisecond
)

var (
	eighteenDecimalsInBaseUnits = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
)

var (
	rpcClient           *ethrpc.Client
	ethClient           *ethclient.Client
	ethRPCClient        ethrpcclient.Client
	zrx                 *wrappers.ZRXToken
	dummyERC721Token    *wrappers.DummyERC721Token
	erc1155Mintable     *wrappers.ERC1155Mintable
	exchange            *wrappers.Exchange
	weth                *wrappers.WETH9
	blockchainLifecycle *ethereum.BlockchainLifecycle
)

// Since these tests must be run sequentially, we don't want them to run as part of
// the normal testing process. They will only be run if the "--serial" flag is used.
var serialTestsEnabled bool

var ganacheAddresses = ethereum.GanacheAddresses

func init() {
	flag.BoolVar(&serialTestsEnabled, "serial", false, "enable serial tests")
	testing.Init()
	flag.Parse()

	var err error
	rpcClient, err = ethrpc.Dial(constants.GanacheEndpoint)
	if err != nil {
		panic(err)
	}
	rateLimiter := ratelimit.NewUnlimited()
	rpcClient, err := rpc.Dial(constants.GanacheEndpoint)
	if err != nil {
		panic(err)
	}
	ethRPCClient, err = ethrpcclient.New(rpcClient, ethereumRPCRequestTimeout, rateLimiter)
	if err != nil {
		panic(err)
	}
	blockchainLifecycle, err = ethereum.NewBlockchainLifecycle(rpcClient)
	if err != nil {
		panic(err)
	}
	ethClient = ethclient.NewClient(rpcClient)
	zrx, err = wrappers.NewZRXToken(ganacheAddresses.ZRXToken, ethClient)
	if err != nil {
		panic(err)
	}
	dummyERC721Token, err = wrappers.NewDummyERC721Token(constants.GanacheDummyERC721TokenAddress, ethClient)
	if err != nil {
		panic(err)
	}
	erc1155Mintable, err = wrappers.NewERC1155Mintable(constants.GanacheDummyERC1155MintableAddress, ethClient)
	if err != nil {
		panic(err)
	}
	exchange, err = wrappers.NewExchange(ganacheAddresses.Exchange, ethClient)
	if err != nil {
		panic(err)
	}
	weth, err = wrappers.NewWETH9(ganacheAddresses.WETH9, ethClient)
	if err != nil {
		panic(err)
	}
}

func TestOrderWatcherStoresValidOrders(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetData(scenario.ZRXAssetData),
	)
	setupOrderWatcherScenario(ctx, t, ethClient, database, signedOrder)
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)

	latestStoredBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)

	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          false,
		fillableAmount:     signedOrder.TakerAssetAmount,
		lastUpdated:        time.Now(),
		lastValidatedBlock: latestStoredBlock,
	}
	checkOrderState(t, expectedOrderState, orders[0])
}

func TestOrderWatcherUnfundedInsufficientERC20Balance(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetData(scenario.ZRXAssetData),
	)
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, database, signedOrder)

	// Transfer makerAsset out of maker address
	opts := &bind.TransactOpts{
		From:   signedOrder.MakerAddress,
		Signer: scenario.GetTestSignerFn(signedOrder.MakerAddress),
	}
	txn, err := zrx.Transfer(opts, constants.GanacheAccount4, signedOrder.MakerAssetAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState)

	latestStoredBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          true,
		fillableAmount:     big.NewInt(0),
		lastUpdated:        time.Now(),
		lastValidatedBlock: latestStoredBlock,
	}
	checkOrderState(t, expectedOrderState, orders[0])
}

func TestOrderWatcherUnfundedInsufficientERC20BalanceForMakerFee(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	makerAssetData := scenario.GetDummyERC721AssetData(big.NewInt(1))
	wethFeeAmount := new(big.Int).Mul(big.NewInt(5), eighteenDecimalsInBaseUnits)
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetData(makerAssetData),
		orderopts.MakerAssetAmount(big.NewInt(1)),
		orderopts.MakerFeeAssetData(scenario.WETHAssetData),
		orderopts.MakerFee(wethFeeAmount),
	)
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, database, signedOrder)

	// Transfer makerAsset out of maker address
	opts := &bind.TransactOpts{
		From:   signedOrder.MakerAddress,
		Signer: scenario.GetTestSignerFn(signedOrder.MakerAddress),
	}
	txn, err := weth.Transfer(opts, constants.GanacheAccount4, wethFeeAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState)

	latestStoredBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          true,
		fillableAmount:     big.NewInt(0),
		lastUpdated:        time.Now(),
		lastValidatedBlock: latestStoredBlock,
	}
	checkOrderState(t, expectedOrderState, orders[0])
}

func TestOrderWatcherUnfundedInsufficientERC721Balance(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	tokenID := big.NewInt(1)
	makerAssetData := scenario.GetDummyERC721AssetData(tokenID)
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetAmount(big.NewInt(1)),
		orderopts.MakerAssetData(makerAssetData),
	)
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, database, signedOrder)

	// Transfer makerAsset out of maker address
	opts := &bind.TransactOpts{
		From:   signedOrder.MakerAddress,
		Signer: scenario.GetTestSignerFn(signedOrder.MakerAddress),
	}
	txn, err := dummyERC721Token.TransferFrom(opts, signedOrder.MakerAddress, constants.GanacheAccount4, tokenID)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState)

	latestStoredBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          true,
		fillableAmount:     big.NewInt(0),
		lastUpdated:        time.Now(),
		lastValidatedBlock: latestStoredBlock,
	}
	checkOrderState(t, expectedOrderState, orders[0])

}

func TestOrderWatcherUnfundedInsufficientERC721Allowance(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	tokenID := big.NewInt(1)
	makerAssetData := scenario.GetDummyERC721AssetData(tokenID)
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetAmount(big.NewInt(1)),
		orderopts.MakerAssetData(makerAssetData),
	)
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, database, signedOrder)

	// Remove Maker's NFT approval to ERC721Proxy. We do this by setting the
	// operator/spender to the null address.
	opts := &bind.TransactOpts{
		From:   signedOrder.MakerAddress,
		Signer: scenario.GetTestSignerFn(signedOrder.MakerAddress),
	}
	txn, err := dummyERC721Token.Approve(opts, constants.NullAddress, tokenID)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState)

	latestStoredBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          true,
		fillableAmount:     big.NewInt(0),
		lastUpdated:        time.Now(),
		lastValidatedBlock: latestStoredBlock,
	}
	checkOrderState(t, expectedOrderState, orders[0])
}

func TestOrderWatcherUnfundedInsufficientERC1155Allowance(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	makerAssetData := scenario.GetDummyERC1155AssetData(t, []*big.Int{big.NewInt(1)}, []*big.Int{big.NewInt(100)})
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetAmount(big.NewInt(1)),
		orderopts.MakerAssetData(makerAssetData),
	)
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, database, signedOrder)

	// Remove Maker's ERC1155 approval to ERC1155Proxy
	opts := &bind.TransactOpts{
		From:   signedOrder.MakerAddress,
		Signer: scenario.GetTestSignerFn(signedOrder.MakerAddress),
	}
	txn, err := erc1155Mintable.SetApprovalForAll(opts, ganacheAddresses.ERC1155Proxy, false)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState)

	latestStoredBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          true,
		fillableAmount:     big.NewInt(0),
		lastUpdated:        time.Now(),
		lastValidatedBlock: latestStoredBlock,
	}
	checkOrderState(t, expectedOrderState, orders[0])
}

func TestOrderWatcherUnfundedInsufficientERC1155Balance(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	tokenID := big.NewInt(1)
	tokenAmount := big.NewInt(100)
	makerAssetData := scenario.GetDummyERC1155AssetData(t, []*big.Int{tokenID}, []*big.Int{tokenAmount})
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetAmount(big.NewInt(1)),
		orderopts.MakerAssetData(makerAssetData),
	)
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, database, signedOrder)

	// Reduce Maker's ERC1155 balance
	opts := &bind.TransactOpts{
		From:   signedOrder.MakerAddress,
		Signer: scenario.GetTestSignerFn(signedOrder.MakerAddress),
	}
	txn, err := erc1155Mintable.SafeTransferFrom(opts, signedOrder.MakerAddress, constants.GanacheAccount4, tokenID, tokenAmount, []byte{})
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState)

	latestStoredBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          true,
		fillableAmount:     big.NewInt(0),
		lastUpdated:        time.Now(),
		lastValidatedBlock: latestStoredBlock,
	}
	checkOrderState(t, expectedOrderState, orders[0])
}

func TestOrderWatcherUnfundedInsufficientERC20Allowance(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetData(scenario.ZRXAssetData),
	)
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, database, signedOrder)

	// Remove Maker's ZRX approval to ERC20Proxy
	opts := &bind.TransactOpts{
		From:   signedOrder.MakerAddress,
		Signer: scenario.GetTestSignerFn(signedOrder.MakerAddress),
	}
	txn, err := zrx.Approve(opts, ganacheAddresses.ERC20Proxy, big.NewInt(0))
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState)

	latestStoredBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          true,
		fillableAmount:     big.NewInt(0),
		lastUpdated:        time.Now(),
		lastValidatedBlock: latestStoredBlock,
	}
	checkOrderState(t, expectedOrderState, orders[0])
}

func TestOrderWatcherUnfundedThenFundedAgain(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetData(scenario.ZRXAssetData),
		orderopts.TakerAssetData(scenario.WETHAssetData),
	)
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, database, signedOrder)

	// Transfer makerAsset out of maker address
	opts := &bind.TransactOpts{
		From:   signedOrder.MakerAddress,
		Signer: scenario.GetTestSignerFn(signedOrder.MakerAddress),
	}
	txn, err := zrx.Transfer(opts, constants.GanacheAccount4, signedOrder.MakerAssetAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState)

	latestStoredBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          true,
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
	txn, err = zrx.Transfer(opts, signedOrder.MakerAddress, signedOrder.MakerAssetAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents = <-orderEventsChan
	require.Len(t, orderEvents, 1)
	orderEvent = orderEvents[0]
	assert.Equal(t, zeroex.ESOrderAdded, orderEvent.EndState)

	latestStoredBlock, err = database.GetLatestMiniHeader()
	require.NoError(t, err)
	newOrders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, newOrders, 1)
	expectedOrderState = orderState{
		hash:               expectedOrderHash,
		isRemoved:          false,
		fillableAmount:     signedOrder.TakerAssetAmount,
		lastUpdated:        time.Now(),
		lastValidatedBlock: latestStoredBlock,
	}
	checkOrderState(t, expectedOrderState, newOrders[0])
}

func TestOrderWatcherNoChange(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetData(scenario.ZRXAssetData),
		orderopts.TakerAssetData(scenario.WETHAssetData),
	)
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	blockWatcher, _ := setupOrderWatcherScenario(ctx, t, ethClient, database, signedOrder)

	latestStoredBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          false,
		fillableAmount:     signedOrder.TakerAssetAmount,
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
	txn, err := zrx.Transfer(opts, signedOrder.MakerAddress, signedOrder.MakerAssetAmount)
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
	newOrders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, newOrders, 1)
	expectedOrderState = orderState{
		hash:               expectedOrderHash,
		isRemoved:          false,
		fillableAmount:     signedOrder.TakerAssetAmount,
		lastUpdated:        time.Now(),
		lastValidatedBlock: latestStoredBlock,
	}
	checkOrderState(t, expectedOrderState, newOrders[0])
}

func TestOrderWatcherWETHWithdrawAndDeposit(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetData(scenario.WETHAssetData),
		orderopts.TakerAssetData(scenario.ZRXAssetData),
	)
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, database, signedOrder)

	// Withdraw maker's WETH (i.e. decrease WETH balance)
	// HACK(fabio): For some reason the txn fails with "out of gas" error with the
	// estimated gas amount
	gasLimit := uint64(50000)
	opts := &bind.TransactOpts{
		From:     signedOrder.MakerAddress,
		Signer:   scenario.GetTestSignerFn(signedOrder.MakerAddress),
		GasLimit: gasLimit,
	}
	txn, err := weth.Withdraw(opts, signedOrder.MakerAssetAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState)

	latestStoredBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          true,
		fillableAmount:     big.NewInt(0),
		lastUpdated:        time.Now(),
		lastValidatedBlock: latestStoredBlock,
	}
	checkOrderState(t, expectedOrderState, orders[0])

	// Deposit maker's ETH (i.e. increase WETH balance)
	opts = &bind.TransactOpts{
		From:   signedOrder.MakerAddress,
		Signer: scenario.GetTestSignerFn(signedOrder.MakerAddress),
		Value:  signedOrder.MakerAssetAmount,
	}
	txn, err = weth.Deposit(opts)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents = <-orderEventsChan
	require.Len(t, orderEvents, 1)
	orderEvent = orderEvents[0]
	assert.Equal(t, zeroex.ESOrderAdded, orderEvent.EndState)

	latestStoredBlock, err = database.GetLatestMiniHeader()
	require.NoError(t, err)
	newOrders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, newOrders, 1)
	expectedOrderState = orderState{
		hash:               expectedOrderHash,
		isRemoved:          false,
		fillableAmount:     signedOrder.TakerAssetAmount,
		lastUpdated:        time.Now(),
		lastValidatedBlock: latestStoredBlock,
	}
	checkOrderState(t, expectedOrderState, newOrders[0])
}

func TestOrderWatcherCanceled(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	signedOrder := scenario.NewSignedTestOrder(t, orderopts.SetupMakerState(true))
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, database, signedOrder)

	// Cancel order
	opts := &bind.TransactOpts{
		From:   signedOrder.MakerAddress,
		Signer: scenario.GetTestSignerFn(signedOrder.MakerAddress),
	}
	trimmedOrder := signedOrder.Trim()
	txn, err := exchange.CancelOrder(opts, trimmedOrder)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderCancelled, orderEvent.EndState)

	latestStoredBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          true,
		fillableAmount:     big.NewInt(0),
		lastUpdated:        time.Now(),
		lastValidatedBlock: latestStoredBlock,
	}
	checkOrderState(t, expectedOrderState, orders[0])
}

func TestOrderWatcherCancelUpTo(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	signedOrder := scenario.NewSignedTestOrder(t, orderopts.SetupMakerState(true))
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, database, signedOrder)

	// Cancel order with epoch
	opts := &bind.TransactOpts{
		From:   signedOrder.MakerAddress,
		Signer: scenario.GetTestSignerFn(signedOrder.MakerAddress),
	}
	targetOrderEpoch := signedOrder.Salt
	txn, err := exchange.CancelOrdersUpTo(opts, targetOrderEpoch)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderCancelled, orderEvent.EndState)

	latestStoredBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          true,
		fillableAmount:     big.NewInt(0),
		lastUpdated:        time.Now(),
		lastValidatedBlock: latestStoredBlock,
	}
	checkOrderState(t, expectedOrderState, orders[0])
}

func TestOrderWatcherERC20Filled(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	takerAddress := constants.GanacheAccount3
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.SetupTakerAddress(takerAddress),
	)
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, database, signedOrder)

	// Fill order
	opts := &bind.TransactOpts{
		From:   takerAddress,
		Signer: scenario.GetTestSignerFn(takerAddress),
		Value:  big.NewInt(100000000000000000),
	}
	trimmedOrder := signedOrder.Trim()
	txn, err := exchange.FillOrder(opts, trimmedOrder, signedOrder.TakerAssetAmount, signedOrder.Signature)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderFullyFilled, orderEvent.EndState)

	latestStoredBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          true,
		fillableAmount:     big.NewInt(0),
		lastUpdated:        time.Now(),
		lastValidatedBlock: latestStoredBlock,
	}
	checkOrderState(t, expectedOrderState, orders[0])
}

func TestOrderWatcherERC20PartiallyFilled(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	takerAddress := constants.GanacheAccount3
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.SetupTakerAddress(takerAddress),
	)
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, database, signedOrder)

	// Partially fill order
	opts := &bind.TransactOpts{
		From:   takerAddress,
		Signer: scenario.GetTestSignerFn(takerAddress),
		Value:  big.NewInt(100000000000000000),
	}
	trimmedOrder := signedOrder.Trim()
	halfAmount := new(big.Int).Div(signedOrder.TakerAssetAmount, big.NewInt(2))
	txn, err := exchange.FillOrder(opts, trimmedOrder, halfAmount, signedOrder.Signature)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderFilled, orderEvent.EndState)

	latestStoredBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          false,
		fillableAmount:     halfAmount,
		lastUpdated:        time.Now(),
		lastValidatedBlock: latestStoredBlock,
	}
	checkOrderState(t, expectedOrderState, orders[0])
}

func TestOrderWatcherOrderExpiredThenUnexpired(t *testing.T) {
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

	// Create and add an order (which will later become expired) to OrderWatcher
	expirationTime := time.Now().Add(24 * time.Hour)
	expirationTimeSeconds := big.NewInt(expirationTime.Unix())
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.ExpirationTimeSeconds(expirationTimeSeconds),
	)
	expectedOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	blockwatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)
	watchOrder(ctx, t, orderWatcher, blockwatcher, signedOrder, false)

	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// Simulate a block found with a timestamp past expirationTime
	latestBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
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
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderExpired, orderEvent.EndState)

	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	expectedOrderState := orderState{
		hash:               expectedOrderHash,
		isRemoved:          true,
		fillableAmount:     signedOrder.TakerAssetAmount,
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
	require.Len(t, orderEvents, 1)
	orderEvent = orderEvents[0]
	assert.Equal(t, zeroex.ESOrderUnexpired, orderEvent.EndState)

	newOrders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, newOrders, 1)
	expectedOrderState = orderState{
		hash:               expectedOrderHash,
		isRemoved:          false,
		fillableAmount:     signedOrder.TakerAssetAmount,
		lastUpdated:        time.Now(),
		lastValidatedBlock: reorgBlockEvents[len(reorgBlockEvents)-1].BlockHeader,
	}
	checkOrderState(t, expectedOrderState, newOrders[0])
}

func TestOrderWatcherDecreaseExpirationTime(t *testing.T) {
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
	signedOrders := scenario.NewSignedTestOrdersBatch(t, maxOrders, optionsForIndex)
	for _, signedOrder := range signedOrders {
		watchOrder(ctx, t, orderWatcher, blockWatcher, signedOrder, false)
	}

	// We don't care about the order events above for the purposes of this test,
	// so we only subscribe now.
	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// The next order should cause some orders to be removed and the appropriate
	// events to fire.
	expirationTime := time.Now().Add(10*time.Minute + 1*time.Second)
	expirationTimeSeconds := big.NewInt(expirationTime.Unix())
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.ExpirationTimeSeconds(expirationTimeSeconds),
	)
	watchOrder(ctx, t, orderWatcher, blockWatcher, signedOrder, false)
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
			orderExpirationTime := orderEvent.SignedOrder.ExpirationTimeSeconds
			assert.True(t, orderExpirationTime.Cmp(storedMaxExpirationTime) == -1, "ADDED order has an expiration time of %s which is *greater than* the maximum of %s", orderExpirationTime, storedMaxExpirationTime)
		case zeroex.ESStoppedWatching:
			numStoppedWatching += 1
			orderExpirationTime := orderEvent.SignedOrder.ExpirationTimeSeconds
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
	remainingOrders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, remainingOrders, expectedRemainingOrders)
	for _, order := range remainingOrders {
		assert.True(t, order.ExpirationTimeSeconds.Cmp(storedMaxExpirationTime) != 1, "remaining order has an expiration time of %s which is *greater than* the maximum of %s", order.ExpirationTimeSeconds, storedMaxExpirationTime)
	}

	// Confirm that a pinned order will be accepted even if its expiration
	// is greater than the current max.
	pinnedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.ExpirationTimeSeconds(big.NewInt(0).Add(storedMaxExpirationTime, big.NewInt(10))),
	)
	pinnedOrderHash, err := pinnedOrder.ComputeOrderHash()
	require.NoError(t, err)
	watchOrder(ctx, t, orderWatcher, blockWatcher, pinnedOrder, true)

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

func TestOrderWatcherBatchEmitsAddedEvents(t *testing.T) {
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
	signedOrders := scenario.NewSignedTestOrdersBatch(t, numOrders, orderOptions)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	validationResults, err := orderWatcher.ValidateAndStoreValidOrders(ctx, signedOrders, false, constants.TestChainID)
	require.Len(t, validationResults.Rejected, 0)
	require.NoError(t, err)

	orderEvents := <-orderEventsChan
	require.Len(t, orderEvents, numOrders)
	for _, orderEvent := range orderEvents {
		assert.Equal(t, zeroex.ESOrderAdded, orderEvent.EndState)
	}

	orders, err := database.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, numOrders)
}

func TestOrderWatcherCleanup(t *testing.T) {
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
	signedOrders := scenario.NewSignedTestOrdersBatch(t, 2, orderOptions)
	signedOrderOne := signedOrders[0]
	watchOrder(ctx, t, orderWatcher, blockWatcher, signedOrderOne, false)
	signedOrderTwo := signedOrders[1]
	watchOrder(ctx, t, orderWatcher, blockWatcher, signedOrderTwo, false)
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
	case _ = <-orderEventsChan:
		t.Error("Expected no orderEvents to fire after calling Cleanup()")
	case <-time.After(100 * time.Millisecond):
		// Noop
	}
}

func TestOrderWatcherHandleOrderExpirationsExpired(t *testing.T) {
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
	signedOrders := scenario.NewSignedTestOrdersBatch(t, 2, orderOptions)
	signedOrderOne := signedOrders[0]
	signedOrderTwo := signedOrders[1]
	blockwatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)
	watchOrder(ctx, t, orderWatcher, blockwatcher, signedOrderOne, false)
	watchOrder(ctx, t, orderWatcher, blockwatcher, signedOrderTwo, false)

	signedOrderOneHash, err := signedOrderOne.ComputeOrderHash()
	require.NoError(t, err)
	orderOne, err := database.GetOrder(signedOrderOneHash)
	require.NoError(t, err)
	// Since we flag SignedOrderOne for revalidation, we expect `handleOrderExpirations` not to return an
	// expiry event for it.
	ordersToRevalidate := map[common.Hash]*types.OrderWithMetadata{
		signedOrderOneHash: orderOne,
	}

	// Make a "fake" block with a timestamp 1 second after expirationTime.
	latestBlock, err := database.GetLatestMiniHeader()
	require.NoError(t, err)
	latestBlock.Timestamp = expirationTime.Add(1 * time.Second)
	orderEvents, err := orderWatcher.handleOrderExpirations(latestBlock, ordersToRevalidate)
	require.NoError(t, err)

	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	signedOrderTwoHash, err := signedOrderTwo.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, signedOrderTwoHash, orderEvent.OrderHash)
	assert.Equal(t, zeroex.ESOrderExpired, orderEvent.EndState)
	assert.Equal(t, big.NewInt(0), orderEvent.FillableTakerAssetAmount)
	assert.Len(t, orderEvent.ContractEvents, 0)

	orderTwo, err := database.GetOrder(signedOrderTwoHash)
	require.NoError(t, err)
	assert.Equal(t, true, orderTwo.IsRemoved)
}

func TestOrderWatcherHandleOrderExpirationsUnexpired(t *testing.T) {
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
	signedOrders := scenario.NewSignedTestOrdersBatch(t, 2, orderOptions)
	signedOrderOne := signedOrders[0]
	signedOrderTwo := signedOrders[1]
	blockwatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)
	watchOrder(ctx, t, orderWatcher, blockwatcher, signedOrderOne, false)
	watchOrder(ctx, t, orderWatcher, blockwatcher, signedOrderTwo, false)

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
	orderEvents, err = orderWatcher.handleOrderExpirations(latestBlock, ordersToRevalidate)
	require.NoError(t, err)

	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	signedOrderTwoHash, err := signedOrderTwo.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, signedOrderTwoHash, orderEvent.OrderHash)
	assert.Equal(t, zeroex.ESOrderUnexpired, orderEvent.EndState)
	assert.Equal(t, signedOrderTwo.TakerAssetAmount, orderEvent.FillableTakerAssetAmount)
	assert.Len(t, orderEvent.ContractEvents, 0)

	orderTwo, err := database.GetOrder(signedOrderTwoHash)
	require.NoError(t, err)
	assert.Equal(t, false, orderTwo.IsRemoved)
}

// Scenario: Order has become unexpired and filled in the same block events processed. We test this case using
// `convertValidationResultsIntoOrderEvents` since we cannot properly time-travel using Ganache.
// Source: https://github.com/trufflesuite/ganache-cli/issues/708
func TestConvertValidationResultsIntoOrderEventsUnexpired(t *testing.T) {
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
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.ExpirationTimeSeconds(expirationTimeSeconds),
	)
	blockwatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)
	watchOrder(ctx, t, orderWatcher, blockwatcher, signedOrder, false)

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
				SignedOrder:              signedOrder,
				FillableTakerAssetAmount: big.NewInt(1).Div(signedOrder.TakerAssetAmount, big.NewInt(2)),
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
	orderEvents, err = orderWatcher.convertValidationResultsIntoOrderEvents(&validationResults, orderHashToDBOrder, orderHashToEvents, validationBlock)
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
}

func TestDrainAllBlockEventsChan(t *testing.T) {
	blockEventsChan := make(chan []*blockwatch.Event, 100)
	ts := time.Now().Add(1 * time.Hour)
	blockEventsOne := []*blockwatch.Event{
		{
			Type: blockwatch.Added,
			BlockHeader: &types.MiniHeader{
				Parent:    common.HexToHash("0x0"),
				Hash:      common.HexToHash("0x1"),
				Number:    big.NewInt(1),
				Timestamp: ts,
			},
		},
	}
	blockEventsChan <- blockEventsOne

	blockEventsTwo := []*blockwatch.Event{
		{
			Type: blockwatch.Added,
			BlockHeader: &types.MiniHeader{
				Parent:    common.HexToHash("0x1"),
				Hash:      common.HexToHash("0x2"),
				Number:    big.NewInt(2),
				Timestamp: ts.Add(1 * time.Second),
			},
		},
	}
	blockEventsChan <- blockEventsTwo

	max := 2 // enough
	allEvents := drainBlockEventsChan(blockEventsChan, max)
	assert.Len(t, allEvents, 2, "Two events should be drained from the channel")
	require.Equal(t, allEvents[0], blockEventsOne[0])
	require.Equal(t, allEvents[1], blockEventsTwo[0])

	// Test case where more than max events in channel
	blockEventsChan <- blockEventsOne
	blockEventsChan <- blockEventsTwo

	max = 1
	allEvents = drainBlockEventsChan(blockEventsChan, max)
	assert.Len(t, allEvents, 1, "Only max number of events should be drained from the channel, even if more than max events are present")
	require.Equal(t, allEvents[0], blockEventsOne[0])
}

func TestMissingOrderEvents(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	// Set up test and orderWatcher
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	// FIXME(jalextowle): Add timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	database, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	validator, err := ordervalidator.New(
		&SlowContractCaller{
			caller:            ethRPCClient,
			contractCallDelay: 5 * time.Second,
		},
		constants.TestChainID,
		ethereumRPCMaxContentLength,
		ganacheAddresses,
	)
	require.NoError(t, err)

	blockWatcher, orderWatcher := setupOrderWatcherWithValidator(ctx, t, ethRPCClient, database, validator)
	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// Create a new order
	signedOrder := scenario.NewSignedTestOrder(t, orderopts.SetupMakerState(true))
	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	// Cancel the order
	opts := &bind.TransactOpts{
		From:   signedOrder.MakerAddress,
		Signer: scenario.GetTestSignerFn(signedOrder.MakerAddress),
	}
	trimmedOrder := signedOrder.Trim()
	txn, err := exchange.CancelOrder(opts, trimmedOrder)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	syncErrChan := make(chan error)
	validationErrChan := make(chan error)
	validationResultsChan := make(chan *ordervalidator.ValidationResults)

	go func() {
		err = blockWatcher.SyncToLatestBlock()
		syncErrChan <- err
	}()

	go func() {
		validationResults, err := orderWatcher.ValidateAndStoreValidOrders(ctx, []*zeroex.SignedOrder{signedOrder}, false, constants.TestChainID)
		if err != nil {
			validationErrChan <- err
			return
		}
		validationResultsChan <- validationResults
	}()

	err = <-syncErrChan
	require.NoError(t, err)

	select {
	case err := <-validationErrChan:
		t.Error(err)
	case validationResults := <-validationResultsChan:
		require.Equal(t, len(validationResults.Accepted), 1)
		assert.Equal(t, len(validationResults.Rejected), 0)
		orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
		assert.Equal(t, zeroex.ESOrderAdded, orderEvents[0].EndState)
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
	orderEvents := waitForOrderEvents(t, orderEventsChan, 2, 10*time.Second)
	assert.Equal(t, zeroex.ESOrderCancelled, orderEvents[0].EndState)
	// TODO(jalextowle): This event probably shouldn't be fired
	assert.Equal(t, zeroex.ESOrderUnexpired, orderEvents[1].EndState)
}

func setupOrderWatcherScenario(ctx context.Context, t *testing.T, ethClient *ethclient.Client, database *db.DB, signedOrder *zeroex.SignedOrder) (*blockwatch.Watcher, chan []*zeroex.OrderEvent) {
	blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)

	// Start watching an order
	watchOrder(ctx, t, orderWatcher, blockWatcher, signedOrder, false)

	// Subscribe to OrderWatcher
	orderEventsChan := make(chan []*zeroex.OrderEvent, 10)
	orderWatcher.Subscribe(orderEventsChan)

	return blockWatcher, orderEventsChan
}

func watchOrder(ctx context.Context, t *testing.T, orderWatcher *Watcher, blockWatcher *blockwatch.Watcher, signedOrder *zeroex.SignedOrder, pinned bool) {
	err := blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	validationResults, err := orderWatcher.ValidateAndStoreValidOrders(ctx, []*zeroex.SignedOrder{signedOrder}, pinned, constants.TestChainID)
	require.NoError(t, err)
	if len(validationResults.Rejected) != 0 {
		spew.Dump(validationResults.Rejected)
	}
	require.Len(t, validationResults.Accepted, 1, "Expected order to pass validation and get added to OrderWatcher")
}

func setupOrderWatcher(ctx context.Context, t *testing.T, ethRPCClient ethrpcclient.Client, database *db.DB) (*blockwatch.Watcher, *Watcher) {
	orderValidator, err := ordervalidator.New(ethRPCClient, constants.TestChainID, ethereumRPCMaxContentLength, ganacheAddresses)
	require.NoError(t, err)
	return setupOrderWatcherWithValidator(ctx, t, ethRPCClient, database, orderValidator)
}

func setupOrderWatcherWithValidator(ctx context.Context, t *testing.T, ethRPCClient ethrpcclient.Client, database *db.DB, v *ordervalidator.OrderValidator) (*blockwatch.Watcher, *Watcher) {
	blockWatcherClient, err := blockwatch.NewRpcClient(ethRPCClient)
	require.NoError(t, err)
	topics := GetRelevantTopics()
	blockWatcherConfig := blockwatch.Config{
		DB:              database,
		PollingInterval: blockPollingInterval,
		WithLogs:        true,
		Topics:          topics,
		Client:          blockWatcherClient,
	}
	blockWatcher := blockwatch.New(blockRetentionLimit, blockWatcherConfig)
	require.NoError(t, err)
	orderWatcher, err := New(Config{
		DB:                database,
		BlockWatcher:      blockWatcher,
		OrderValidator:    v,
		ChainID:           constants.TestChainID,
		ContractAddresses: ganacheAddresses,
		MaxOrders:         1000,
	})
	require.NoError(t, err)

	// Start OrderWatcher
	go func() {
		err := orderWatcher.Watch(ctx)
		if err != nil {
			panic(err)
		}
	}()

	// Ensure at least one block has been processed and is stored in the DB
	// before tests run
	storedBlocks, err := database.FindMiniHeaders(nil)
	require.NoError(t, err)
	if len(storedBlocks) == 0 {
		err := blockWatcher.SyncToLatestBlock()
		require.NoError(t, err)
	}

	err = orderWatcher.WaitForAtLeastOneBlockToBeProcessed(ctx)
	require.NoError(t, err)

	return blockWatcher, orderWatcher
}

var _ bind.ContractCaller = &SlowContractCaller{}

// SlowContractCaller satisfies the bind.ContractCall interface by wrapping another
// contract caller and adding delays before the contract call.
type SlowContractCaller struct {
	caller            bind.ContractCaller
	contractCallDelay time.Duration
	codeAtDelay       time.Duration
}

func (s *SlowContractCaller) CallContract(ctx context.Context, call geth.CallMsg, blockNumber *big.Int) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(s.contractCallDelay):
	}
	return s.caller.CallContract(ctx, call, blockNumber)
}

func (s *SlowContractCaller) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(s.codeAtDelay):
	}
	return s.caller.CodeAt(ctx, contract, blockNumber)
}

func setupSubTest(t *testing.T) func(t *testing.T) {
	blockchainLifecycle.Start(t)
	return func(t *testing.T) {
		blockchainLifecycle.Revert(t)
	}
}

func waitForOrderEvents(t *testing.T, orderEventsChan <-chan []*zeroex.OrderEvent, expectedNumberOfEvents int, timeout time.Duration) []*zeroex.OrderEvent {
	allOrderEvents := []*zeroex.OrderEvent{}
	for {
		select {
		case orderEvents := <-orderEventsChan:
			allOrderEvents = append(allOrderEvents, orderEvents...)
			if len(allOrderEvents) >= expectedNumberOfEvents {
				return allOrderEvents
			}
			continue
		case <-time.After(timeout):
			t.Fatalf("timed out waiting for %d order events (received %d events)", expectedNumberOfEvents, len(allOrderEvents))
		}
	}
}

func waitTxnSuccessfullyMined(t *testing.T, ethClient *ethclient.Client, txn *ethtypes.Transaction) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, ethClient, txn)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, uint64(1))
}

// orderState contains the order hash and some metadata fields for an order that are updated
// by/controlled by order watcher.
type orderState struct {
	hash               common.Hash
	isRemoved          bool
	fillableAmount     *big.Int
	lastUpdated        time.Time
	lastValidatedBlock *types.MiniHeader
}

func checkOrderState(t *testing.T, expectedState orderState, order *types.OrderWithMetadata) {
	assert.Equal(t, expectedState.hash, order.Hash, "Hash")
	assert.Equal(t, expectedState.isRemoved, order.IsRemoved, "IsRemoved")
	assert.Equal(t, expectedState.fillableAmount, order.FillableTakerAssetAmount, "FillableTakerAssetAmount")
	assert.WithinDuration(t, expectedState.lastUpdated, order.LastUpdated, 4*time.Second, "LastUpdated")
	assert.Equal(t, expectedState.lastValidatedBlock.Number, order.LastValidatedBlockNumber, "LastValidatedBlockNumber")
	assert.Equal(t, expectedState.lastValidatedBlock.Hash, order.LastValidatedBlockHash, "LastValidatedBlockHash")
}
