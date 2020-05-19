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
	ethereumRPCRequestTimeout   = 30 * time.Second
	miniHeaderRetentionLimit    = 2
	blockPollingInterval        = 1000 * time.Millisecond
	ethereumRPCMaxContentLength = 524288
	maxEthRPCRequestsPer24HrUTC = 1000000
	maxEthRPCRequestsPerSeconds = 1000.0
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

// TODO(albrow): Figure out why this test is failing.
func TestOrderWatcherUnfundedInsufficientERC20Balance(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	meshDB, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetData(scenario.ZRXAssetData),
	)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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

	orders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
	assert.Equal(t, big.NewInt(0), orders[0].FillableTakerAssetAmount)
}

func TestOrderWatcherUnfundedInsufficientERC20BalanceForMakerFee(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	meshDB, err := db.New(ctx, db.TestOptions())
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
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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

	orders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
}

func TestOrderWatcherUnfundedInsufficientERC721Balance(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	meshDB, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	tokenID := big.NewInt(1)
	makerAssetData := scenario.GetDummyERC721AssetData(tokenID)
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetAmount(big.NewInt(1)),
		orderopts.MakerAssetData(makerAssetData),
	)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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

	orders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
	assert.Equal(t, big.NewInt(0), orders[0].FillableTakerAssetAmount)
}

func TestOrderWatcherUnfundedInsufficientERC721Allowance(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	meshDB, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	tokenID := big.NewInt(1)
	makerAssetData := scenario.GetDummyERC721AssetData(tokenID)
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetAmount(big.NewInt(1)),
		orderopts.MakerAssetData(makerAssetData),
	)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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

	orders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
	assert.Equal(t, big.NewInt(0), orders[0].FillableTakerAssetAmount)
}

func TestOrderWatcherUnfundedInsufficientERC1155Allowance(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	meshDB, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	makerAssetData := scenario.GetDummyERC1155AssetData(t, []*big.Int{big.NewInt(1)}, []*big.Int{big.NewInt(100)})
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetAmount(big.NewInt(1)),
		orderopts.MakerAssetData(makerAssetData),
	)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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

	orders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
	assert.Equal(t, big.NewInt(0), orders[0].FillableTakerAssetAmount)
}

func TestOrderWatcherUnfundedInsufficientERC1155Balance(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	meshDB, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	tokenID := big.NewInt(1)
	tokenAmount := big.NewInt(100)
	makerAssetData := scenario.GetDummyERC1155AssetData(t, []*big.Int{tokenID}, []*big.Int{tokenAmount})
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetAmount(big.NewInt(1)),
		orderopts.MakerAssetData(makerAssetData),
	)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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

	orders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
	assert.Equal(t, big.NewInt(0), orders[0].FillableTakerAssetAmount)
}

func TestOrderWatcherUnfundedInsufficientERC20Allowance(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	meshDB, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetData(scenario.ZRXAssetData),
	)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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

	orders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
	assert.Equal(t, big.NewInt(0), orders[0].FillableTakerAssetAmount)
}

func TestOrderWatcherUnfundedThenFundedAgain(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	meshDB, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetData(scenario.ZRXAssetData),
		orderopts.TakerAssetData(scenario.WETHAssetData),
	)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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

	orders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
	assert.Equal(t, big.NewInt(0), orders[0].FillableTakerAssetAmount)

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

	newOrders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, newOrders, 1)
	assert.Equal(t, orderEvent.OrderHash, newOrders[0].Hash)
	assert.Equal(t, false, newOrders[0].IsRemoved)
	assert.Equal(t, signedOrder.TakerAssetAmount, newOrders[0].FillableTakerAssetAmount)
}

func TestOrderWatcherNoChange(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	meshDB, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetData(scenario.ZRXAssetData),
		orderopts.TakerAssetData(scenario.WETHAssetData),
	)
	blockWatcher, _ := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	orders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	dbOrder := orders[0]
	assert.Equal(t, false, dbOrder.IsRemoved)

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

	newOrders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, newOrders, 1)
	require.NotEqual(t, dbOrder.LastUpdated, newOrders[0].Hash)
	assert.Equal(t, false, newOrders[0].IsRemoved)
	assert.Equal(t, signedOrder.TakerAssetAmount, orders[0].FillableTakerAssetAmount)
}

func TestOrderWatcherWETHWithdrawAndDeposit(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	meshDB, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.MakerAssetData(scenario.WETHAssetData),
		orderopts.TakerAssetData(scenario.ZRXAssetData),
	)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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

	orders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)

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

	newOrders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, newOrders, 1)
	assert.Equal(t, orderEvent.OrderHash, newOrders[0].Hash)
	assert.Equal(t, false, newOrders[0].IsRemoved)
	assert.Equal(t, big.NewInt(0), orders[0].FillableTakerAssetAmount)
}

func TestOrderWatcherCanceled(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	meshDB, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	signedOrder := scenario.NewSignedTestOrder(t, orderopts.SetupMakerState(true))
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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

	orders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
	assert.Equal(t, big.NewInt(0), orders[0].FillableTakerAssetAmount)
}

func TestOrderWatcherCancelUpTo(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	meshDB, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	signedOrder := scenario.NewSignedTestOrder(t, orderopts.SetupMakerState(true))
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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

	orders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
	assert.Equal(t, big.NewInt(0), orders[0].FillableTakerAssetAmount)
}

func TestOrderWatcherERC20Filled(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	meshDB, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	takerAddress := constants.GanacheAccount3
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.SetupTakerAddress(takerAddress),
	)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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

	orders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
	assert.Equal(t, big.NewInt(0), orders[0].FillableTakerAssetAmount)
}

func TestOrderWatcherERC20PartiallyFilled(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	meshDB, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	takerAddress := constants.GanacheAccount3
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.SetupTakerAddress(takerAddress),
	)
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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

	orders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, false, orders[0].IsRemoved)
	assert.Equal(t, halfAmount, orders[0].FillableTakerAssetAmount)
}

// TODO(albrow): Needs more MiniHeader methods.
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
	blockwatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, database)
	watchOrder(ctx, t, orderWatcher, blockwatcher, ethClient, signedOrder)

	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// Simulate a block found with a timestamp past expirationTime
	latestBlocks, err := database.FindMiniHeaders(&db.MiniHeaderQuery{
		Limit: 1,
		Sort: []db.MiniHeaderSort{
			{
				Field:     db.MFNumber,
				Direction: db.Descending,
			},
		},
	})
	require.NoError(t, err)
	if len(latestBlocks) == 0 {
		t.Error("No miniHeaders stored in database")
	}
	latestBlock := latestBlocks[0]
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
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
	assert.Equal(t, signedOrder.TakerAssetAmount, orders[0].FillableTakerAssetAmount)

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
	assert.Equal(t, orderEvent.OrderHash, newOrders[0].Hash)
	assert.Equal(t, false, newOrders[0].IsRemoved)
	assert.Equal(t, signedOrder.TakerAssetAmount, newOrders[0].FillableTakerAssetAmount)
}

// TODO(albrow): Re-enable this test or move it.
func TestOrderWatcherDecreaseExpirationTime(t *testing.T) {
	t.Skip("Decreasing expiratin time is not yet implemented")
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	// Set up test and orderWatcher. Manually change maxOrders.
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	meshDB, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	// Store metadata entry in DB
	metadata := &types.Metadata{
		EthereumChainID:   1337,
		MaxExpirationTime: constants.UnlimitedExpirationTime,
	}
	err = meshDB.SaveMetadata(metadata)
	require.NoError(t, err)

	blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, meshDB)
	orderWatcher.maxOrders = 20

	// Create and watch maxOrders orders. Each order has a different expiration time.
	optionsForIndex := func(index int) []orderopts.Option {
		expirationTime := time.Now().Add(10*time.Minute + time.Duration(index)*time.Minute)
		expirationTimeSeconds := big.NewInt(expirationTime.Unix())
		return []orderopts.Option{
			orderopts.SetupMakerState(true),
			orderopts.ExpirationTimeSeconds(expirationTimeSeconds),
		}
	}
	signedOrders := scenario.NewSignedTestOrdersBatch(t, orderWatcher.maxOrders, optionsForIndex)
	for _, signedOrder := range signedOrders {
		watchOrder(ctx, t, orderWatcher, blockWatcher, ethClient, signedOrder)
	}

	// We don't care about the order events above for the purposes of this test,
	// so we only subscribe now.
	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// The next order should cause some orders to be removed and the appropriate
	// events to fire.
	expirationTime := time.Now().Add(10*time.Minute + 1*time.Second)
	expirationTimeSeconds := big.NewInt(expirationTime.Unix())
	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SetupMakerState(true),
		orderopts.ExpirationTimeSeconds(expirationTimeSeconds),
	)
	watchOrder(ctx, t, orderWatcher, blockWatcher, ethClient, signedOrder)
	expectedOrderEvents := int(float64(orderWatcher.maxOrders)*(1-maxOrdersTrimRatio)) + 1
	orderEvents := waitForOrderEvents(t, orderEventsChan, expectedOrderEvents, 4*time.Second)
	require.Len(t, orderEvents, expectedOrderEvents, "wrong number of order events were fired")
	for i, orderEvent := range orderEvents {
		// Last event should be ADDED. The other events should be STOPPED_WATCHING.
		if i == expectedOrderEvents-1 {
			assert.Equal(t, zeroex.ESOrderAdded, orderEvent.EndState, "order event %d had wrong EndState", i)
		} else {
			// For STOPPED_WATCHING events, we also make sure that the expiration time is after
			// the current max expiration time.
			assert.Equal(t, zeroex.ESStoppedWatching, orderEvent.EndState, "order event %d had wrong EndState", i)
			orderExpirationTime := orderEvent.SignedOrder.ExpirationTimeSeconds
			assert.True(t, orderExpirationTime.Cmp(orderWatcher.MaxExpirationTime()) != -1, "remaining order has an expiration time of %s which is *less than* the maximum of %s", orderExpirationTime, orderWatcher.MaxExpirationTime())
		}
	}

	// Now we check that the correct number of orders remain and that all
	// remaining orders have an expiration time less than the current max.
	expectedRemainingOrders := int(float64(orderWatcher.maxOrders)*maxOrdersTrimRatio) + 1
	remainingOrders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, remainingOrders, expectedRemainingOrders)
	for _, order := range remainingOrders {
		assert.True(t, order.ExpirationTimeSeconds.Cmp(orderWatcher.MaxExpirationTime()) == -1, "remaining order has an expiration time of %s which is *greater than* the maximum of %s", order.ExpirationTimeSeconds, orderWatcher.MaxExpirationTime())
	}
}

func TestOrderWatcherBatchEmitsAddedEvents(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	meshDB, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)

	blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, meshDB)

	// Subscribe to OrderWatcher
	orderEventsChan := make(chan []*zeroex.OrderEvent, 10)
	orderWatcher.Subscribe(orderEventsChan)

	// Create numOrders test orders in a batch.
	numOrders := 2
	orderOptions := scenario.OptionsForAll(orderopts.SetupMakerState(true))
	signedOrders := scenario.NewSignedTestOrdersBatch(t, numOrders, orderOptions)

	// Creating a valid order involves transferring sufficient funds to the maker, and setting their allowance for
	// the maker asset. These transactions must be mined and Mesh's BlockWatcher poller must process these blocks
	// in order for the order validation run at order submission to occur at a block number equal or higher then
	// the one where these state changes were included. With the BlockWatcher poller configured to run every 200ms,
	// we wait 500ms here to give it ample time to run before submitting the above order to the Mesh node.
	time.Sleep(500 * time.Millisecond)

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

	orders, err := meshDB.FindOrders(nil)
	require.NoError(t, err)
	require.Len(t, orders, numOrders)
}

func TestOrderWatcherCleanup(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	meshDB, err := db.New(ctx, db.TestOptions())
	require.NoError(t, err)
	blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, meshDB)

	// Create and add two orders to OrderWatcher
	orderOptions := scenario.OptionsForAll(orderopts.SetupMakerState(true))
	signedOrders := scenario.NewSignedTestOrdersBatch(t, 2, orderOptions)
	signedOrderOne := signedOrders[0]
	watchOrder(ctx, t, orderWatcher, blockWatcher, ethClient, signedOrderOne)
	signedOrderTwo := signedOrders[1]
	watchOrder(ctx, t, orderWatcher, blockWatcher, ethClient, signedOrderTwo)
	signedOrderOneHash, err := signedOrderTwo.ComputeOrderHash()
	require.NoError(t, err)

	// Set lastUpdate for signedOrderOne to more than defaultLastUpdatedBuffer so that signedOrderOne
	// does not get re-validated by the cleanup job
	err = meshDB.UpdateOrder(signedOrderOneHash, func(orderToUpdate *types.OrderWithMetadata) (*types.OrderWithMetadata, error) {
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
	meshDB, err := db.New(ctx, db.TestOptions())
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
	blockwatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, meshDB)
	watchOrder(ctx, t, orderWatcher, blockwatcher, ethClient, signedOrderOne)
	watchOrder(ctx, t, orderWatcher, blockwatcher, ethClient, signedOrderTwo)

	signedOrderOneHash, err := signedOrderOne.ComputeOrderHash()
	require.NoError(t, err)
	orderOne, err := meshDB.GetOrder(signedOrderOneHash)
	require.NoError(t, err)
	// Since we flag SignedOrderOne for revalidation, we expect `handleOrderExpirations` not to return an
	// expiry event for it.
	ordersToRevalidate := map[common.Hash]*types.OrderWithMetadata{
		signedOrderOneHash: orderOne,
	}

	// previousLatestBlockTimestamp := expirationTime.Add(-1 * time.Minute)
	latestBlockTimestamp := expirationTime.Add(1 * time.Second)
	orderEvents, err := orderWatcher.handleOrderExpirations(latestBlockTimestamp, ordersToRevalidate)
	require.NoError(t, err)

	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	signedOrderTwoHash, err := signedOrderTwo.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, signedOrderTwoHash, orderEvent.OrderHash)
	assert.Equal(t, zeroex.ESOrderExpired, orderEvent.EndState)
	assert.Equal(t, big.NewInt(0), orderEvent.FillableTakerAssetAmount)
	assert.Len(t, orderEvent.ContractEvents, 0)

	orderTwo, err := meshDB.GetOrder(signedOrderTwoHash)
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
	watchOrder(ctx, t, orderWatcher, blockwatcher, ethClient, signedOrderOne)
	watchOrder(ctx, t, orderWatcher, blockwatcher, ethClient, signedOrderTwo)

	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// Simulate a block found with a timestamp past expirationTime
	latestBlocks, err := database.FindMiniHeaders(&db.MiniHeaderQuery{
		Limit: 1,
		Sort: []db.MiniHeaderSort{
			{
				Field:     db.MFNumber,
				Direction: db.Descending,
			},
		},
	})
	require.NoError(t, err)
	if len(latestBlocks) == 0 {
		t.Error("No miniHeaders stored in database")
	}
	latestBlock := latestBlocks[0]
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

	// LatestBlockTimestamp is earlier than previous latest simulating block-reorg where new latest block
	// has an earlier timestamp than the last
	latestBlockTimestamp := expirationTime.Add(-1 * time.Minute)
	orderEvents, err = orderWatcher.handleOrderExpirations(latestBlockTimestamp, ordersToRevalidate)
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
	watchOrder(ctx, t, orderWatcher, blockwatcher, ethClient, signedOrder)

	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// Simulate a block found with a timestamp past expirationTime. This will mark the order as removed
	// and will remove it from the expiration watcher.
	latestBlocks, err := database.FindMiniHeaders(&db.MiniHeaderQuery{
		Limit: 1,
		Sort: []db.MiniHeaderSort{
			{
				Field:     db.MFNumber,
				Direction: db.Descending,
			},
		},
	})
	require.NoError(t, err)
	if len(latestBlocks) == 0 {
		t.Error("No miniHeaders stored in database")
	}
	latestBlock := latestBlocks[0]
	blockTimestamp := expirationTime.Add(1 * time.Minute)
	nextBlock := &types.MiniHeader{
		Parent:    latestBlock.Hash,
		Hash:      common.HexToHash("0x1"),
		Number:    big.NewInt(0).Add(latestBlock.Number, big.NewInt(1)),
		Timestamp: blockTimestamp,
	}
	expiringBlockEvents := []*blockwatch.Event{
		&blockwatch.Event{
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
			&ordervalidator.AcceptedOrderInfo{
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
	validationBlockTimestamp := expirationTime.Add(-1 * time.Minute)
	orderEvents, err = orderWatcher.convertValidationResultsIntoOrderEvents(&validationResults, orderHashToDBOrder, orderHashToEvents, validationBlockTimestamp)
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

	// var existingOrder meshdb.Order
	// err = database.Orders.FindByID(orderHash.Bytes(), &existingOrder)
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

func setupOrderWatcherScenario(ctx context.Context, t *testing.T, ethClient *ethclient.Client, meshDB *db.DB, signedOrder *zeroex.SignedOrder) (*blockwatch.Watcher, chan []*zeroex.OrderEvent) {
	blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, meshDB)

	// Start watching an order
	watchOrder(ctx, t, orderWatcher, blockWatcher, ethClient, signedOrder)

	// Subscribe to OrderWatcher
	orderEventsChan := make(chan []*zeroex.OrderEvent, 10)
	orderWatcher.Subscribe(orderEventsChan)

	return blockWatcher, orderEventsChan
}

func watchOrder(ctx context.Context, t *testing.T, orderWatcher *Watcher, blockWatcher *blockwatch.Watcher, ethClient *ethclient.Client, signedOrder *zeroex.SignedOrder) {
	err := blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	validationResults, err := orderWatcher.ValidateAndStoreValidOrders(ctx, []*zeroex.SignedOrder{signedOrder}, false, constants.TestChainID)
	require.NoError(t, err)
	if len(validationResults.Rejected) != 0 {
		spew.Dump(validationResults.Rejected)
	}
	require.Len(t, validationResults.Accepted, 1, "Expected order to pass validation and get added to OrderWatcher")
}

func setupOrderWatcher(ctx context.Context, t *testing.T, ethRPCClient ethrpcclient.Client, meshDB *db.DB) (*blockwatch.Watcher, *Watcher) {
	blockWatcherClient, err := blockwatch.NewRpcClient(ethRPCClient)
	require.NoError(t, err)
	topics := GetRelevantTopics()
	// TODO(albrow): May need to be updated.
	blockWatcherConfig := blockwatch.Config{
		DB:              meshDB,
		PollingInterval: blockPollingInterval,
		WithLogs:        true,
		Topics:          topics,
		Client:          blockWatcherClient,
	}
	blockWatcher := blockwatch.New(blockWatcherConfig)
	orderValidator, err := ordervalidator.New(ethRPCClient, constants.TestChainID, ethereumRPCMaxContentLength, ganacheAddresses)
	require.NoError(t, err)
	orderWatcher, err := New(Config{
		DB:                meshDB,
		BlockWatcher:      blockWatcher,
		OrderValidator:    orderValidator,
		ChainID:           constants.TestChainID,
		ContractAddresses: ganacheAddresses,
		MaxExpirationTime: constants.UnlimitedExpirationTime,
		MaxOrders:         1000,
	})
	require.NoError(t, err)

	// Start OrderWatcher
	go func() {
		err := orderWatcher.Watch(ctx)
		require.NoError(t, err)
	}()

	// Ensure at least one block has been processed and is stored in the DB
	// before tests run
	// storedBlocks, err := meshDB.FindAllMiniHeadersSortedByNumber()
	storedBlocks, err := meshDB.FindMiniHeaders(nil)
	require.NoError(t, err)
	if len(storedBlocks) == 0 {
		err := blockWatcher.SyncToLatestBlock()
		require.NoError(t, err)
	}

	err = orderWatcher.WaitForAtLeastOneBlockToBeProcessed(ctx)
	require.NoError(t, err)

	return blockWatcher, orderWatcher
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
	ctx, cancelFn := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancelFn()
	receipt, err := bind.WaitMined(ctx, ethClient, txn)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, uint64(1))
}
