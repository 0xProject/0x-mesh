// +build !js

package orderwatch

import (
	"context"
	"flag"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/blockwatch"
	"github.com/0xProject/0x-mesh/ethereum/ethrpcclient"
	"github.com/0xProject/0x-mesh/ethereum/miniheader"
	"github.com/0xProject/0x-mesh/ethereum/ratelimit"
	"github.com/0xProject/0x-mesh/ethereum/simplestack"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	ethereumRPCRequestTimeout   = 30 * time.Second
	blockWatcherRetentionLimit  = 20
	blockPollingInterval        = 1000 * time.Millisecond
	ethereumRPCMaxContentLength = 524288
	maxEthRPCRequestsPer24HrUTC = 1000000
	maxEthRPCRequestsPerSeconds = 1000.0
)

var (
	makerAddress                = constants.GanacheAccount1
	takerAddress                = constants.GanacheAccount2
	eighteenDecimalsInBaseUnits = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	wethAmount                  = new(big.Int).Mul(big.NewInt(50), eighteenDecimalsInBaseUnits)
	zrxAmount                   = new(big.Int).Mul(big.NewInt(100), eighteenDecimalsInBaseUnits)
	erc1155FungibleAmount       = big.NewInt(100)
	tokenID                     = big.NewInt(1)
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

func init() {
	flag.BoolVar(&serialTestsEnabled, "serial", false, "enable serial tests")
	flag.Parse()
}

func init() {
	var err error
	rpcClient, err = ethrpc.Dial(constants.GanacheEndpoint)
	if err != nil {
		panic(err)
	}
	rateLimiter := ratelimit.NewUnlimited()
	ethRPCClient, err = ethrpcclient.New(constants.GanacheEndpoint, ethereumRPCRequestTimeout, rateLimiter)
	if err != nil {
		panic(err)
	}
	blockchainLifecycle, err = ethereum.NewBlockchainLifecycle(rpcClient)
	if err != nil {
		panic(err)
	}
	ethClient = ethclient.NewClient(rpcClient)
	ganacheAddresses := ethereum.ChainIDToContractAddresses[constants.TestChainID]
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

func TestOrderWatcherUnfundedInsufficientERC20Balance(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	signedOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Transfer makerAsset out of maker address
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	txn, err := zrx.Transfer(opts, constants.GanacheAccount4, zrxAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
	assert.Equal(t, big.NewInt(0), orders[0].FillableTakerAssetAmount)
}

func TestOrderWatcherUnfundedInsufficientERC721Balance(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	signedOrder := scenario.CreateNFTForZRXSignedTestOrder(t, ethClient, makerAddress, takerAddress, tokenID, zrxAmount)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Transfer makerAsset out of maker address
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	txn, err := dummyERC721Token.TransferFrom(opts, makerAddress, constants.GanacheAccount4, tokenID)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
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

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	signedOrder := scenario.CreateNFTForZRXSignedTestOrder(t, ethClient, makerAddress, takerAddress, tokenID, zrxAmount)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Remove Maker's NFT approval to ERC721Proxy
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	ganacheAddresses := ethereum.ChainIDToContractAddresses[constants.TestChainID]
	txn, err := dummyERC721Token.SetApprovalForAll(opts, ganacheAddresses.ERC721Proxy, false)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
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

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	signedOrder := scenario.CreateERC1155ForZRXSignedTestOrder(t, ethClient, makerAddress, takerAddress, tokenID, zrxAmount, erc1155FungibleAmount)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Remove Maker's ERC1155 approval to ERC1155Proxy
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	ganacheAddresses := ethereum.ChainIDToContractAddresses[constants.TestChainID]
	txn, err := erc1155Mintable.SetApprovalForAll(opts, ganacheAddresses.ERC1155Proxy, false)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
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

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	signedOrder := scenario.CreateERC1155ForZRXSignedTestOrder(t, ethClient, makerAddress, takerAddress, tokenID, zrxAmount, erc1155FungibleAmount)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Reduce Maker's ERC1155 balance
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	txn, err := erc1155Mintable.SafeTransferFrom(opts, makerAddress, constants.GanacheAccount4, tokenID, erc1155FungibleAmount, []byte{})
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
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

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	ganacheAddresses := ethereum.ChainIDToContractAddresses[constants.TestChainID]
	signedOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Remove Maker's ZRX approval to ERC20Proxy
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
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

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
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

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	signedOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Transfer makerAsset out of maker address
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	txn, err := zrx.Transfer(opts, constants.GanacheAccount4, zrxAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
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
	txn, err = zrx.Transfer(opts, makerAddress, zrxAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents = <-orderEventsChan
	require.Len(t, orderEvents, 1)
	orderEvent = orderEvents[0]
	assert.Equal(t, zeroex.ESOrderAdded, orderEvent.EndState)

	var newOrders []*meshdb.Order
	err = meshDB.Orders.FindAll(&newOrders)
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

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	signedOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	blockWatcher, _ := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
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
	txn, err := zrx.Transfer(opts, makerAddress, zrxAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	var newOrders []*meshdb.Order
	err = meshDB.Orders.FindAll(&newOrders)
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

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	signedOrder := scenario.CreateWETHForZRXSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Withdraw maker's WETH
	// HACK(fabio): For some reason the txn fails with "out of gas" error with the
	// estimated gas amount
	gasLimit := uint64(50000)
	opts := &bind.TransactOpts{
		From:     makerAddress,
		Signer:   scenario.GetTestSignerFn(makerAddress),
		GasLimit: gasLimit,
	}
	txn, err := weth.Withdraw(opts, wethAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderBecameUnfunded, orderEvent.EndState)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)

	opts = &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
		Value:  wethAmount,
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

	var newOrders []*meshdb.Order
	err = meshDB.Orders.FindAll(&newOrders)
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

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	signedOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Cancel order
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	orderWithoutExchangeAddress := signedOrder.ConvertToOrderWithoutExchangeAddress()
	txn, err := exchange.CancelOrder(opts, orderWithoutExchangeAddress)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderCancelled, orderEvent.EndState)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
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

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	signedOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Cancel order with epoch
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
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

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
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

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	signedOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Fill order
	opts := &bind.TransactOpts{
		From:   takerAddress,
		Signer: scenario.GetTestSignerFn(takerAddress),
	}
	orderWithoutExchangeAddress := signedOrder.ConvertToOrderWithoutExchangeAddress()
	txn, err := exchange.FillOrder(opts, orderWithoutExchangeAddress, wethAmount, signedOrder.Signature)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderFullyFilled, orderEvent.EndState)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
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

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	signedOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Partially fill order
	opts := &bind.TransactOpts{
		From:   takerAddress,
		Signer: scenario.GetTestSignerFn(takerAddress),
	}
	orderWithoutExchangeAddress := signedOrder.ConvertToOrderWithoutExchangeAddress()
	halfAmount := new(big.Int).Div(wethAmount, big.NewInt(2))
	txn, err := exchange.FillOrder(opts, orderWithoutExchangeAddress, halfAmount, signedOrder.Signature)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	orderEvents := waitForOrderEvents(t, orderEventsChan, 1, 4*time.Second)
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderFilled, orderEvent.EndState)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, false, orders[0].IsRemoved)
	assert.Equal(t, halfAmount, orders[0].FillableTakerAssetAmount)
}
func TestOrderWatcherOrderExpiredThenUnexpired(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	// Set up test and orderWatcher
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer func() {
		cancel()
	}()

	// Create and add an order (which will later become expired) to OrderWatcher
	expirationTime := time.Now().Add(24 * time.Hour)
	signedOrder := scenario.CreateSignedTestOrderWithExpirationTime(t, ethClient, makerAddress, takerAddress, expirationTime)
	blockwatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, meshDB)
	watchOrder(ctx, t, orderWatcher, blockwatcher, ethClient, signedOrder)

	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// Simulate a block found with a timestamp past expirationTime
	latestBlock, err := meshDB.FindLatestMiniHeader()
	require.NoError(t, err)
	nextBlock := &miniheader.MiniHeader{
		Parent:    latestBlock.Hash,
		Hash:      common.HexToHash("0x1"),
		Number:    big.NewInt(0).Add(latestBlock.Number, big.NewInt(1)),
		Timestamp: expirationTime.Add(1 * time.Minute),
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
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderExpired, orderEvent.EndState)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
	assert.Equal(t, signedOrder.TakerAssetAmount, orders[0].FillableTakerAssetAmount)

	// Simulate a block re-org
	replacementBlockHash := common.HexToHash("0x2")
	reorgBlockEvents := []*blockwatch.Event{
		&blockwatch.Event{
			Type:        blockwatch.Removed,
			BlockHeader: nextBlock,
		},
		&blockwatch.Event{
			Type: blockwatch.Added,
			BlockHeader: &miniheader.MiniHeader{
				Parent:    nextBlock.Parent,
				Hash:      replacementBlockHash,
				Number:    nextBlock.Number,
				Logs:      []types.Log{},
				Timestamp: expirationTime.Add(-2 * time.Hour),
			},
		},
		&blockwatch.Event{
			Type: blockwatch.Added,
			BlockHeader: &miniheader.MiniHeader{
				Parent:    replacementBlockHash,
				Hash:      common.HexToHash("0x3"),
				Number:    big.NewInt(0).Add(nextBlock.Number, big.NewInt(1)),
				Logs:      []types.Log{},
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

	var newOrders []*meshdb.Order
	err = meshDB.Orders.FindAll(&newOrders)
	require.NoError(t, err)
	require.Len(t, newOrders, 1)
	assert.Equal(t, orderEvent.OrderHash, newOrders[0].Hash)
	assert.Equal(t, false, newOrders[0].IsRemoved)
	assert.Equal(t, signedOrder.TakerAssetAmount, newOrders[0].FillableTakerAssetAmount)
}

func TestOrderWatcherDecreaseExpirationTime(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	// Set up test and orderWatcher. Manually change maxOrders.
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	// Store metadata entry in DB
	metadata := &meshdb.Metadata{
		EthereumChainID:   1337,
		MaxExpirationTime: constants.UnlimitedExpirationTime,
	}
	err = meshDB.SaveMetadata(metadata)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer func() {
		cancel()
	}()
	blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, meshDB)
	orderWatcher.maxOrders = 20

	// create and watch maxOrders orders
	for i := 0; i < orderWatcher.maxOrders; i++ {
		signedOrder := scenario.CreateSignedTestOrderWithExpirationTime(t, ethClient, makerAddress, takerAddress, time.Now().Add(10*time.Minute+time.Duration(i)*time.Minute))
		watchOrder(ctx, t, orderWatcher, blockWatcher, ethClient, signedOrder)
	}

	// We don't care about the order events above for the purposes of this test,
	// so we only subscribe now.
	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// The next order should cause some orders to be removed and the appropriate
	// events to fire.
	signedOrder := scenario.CreateSignedTestOrderWithExpirationTime(t, ethClient, makerAddress, takerAddress, time.Now().Add(10*time.Minute+1*time.Second))
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
	var remainingOrders []*meshdb.Order
	require.NoError(t, meshDB.Orders.FindAll(&remainingOrders))
	require.Len(t, remainingOrders, expectedRemainingOrders)
	for _, order := range remainingOrders {
		assert.True(t, order.SignedOrder.ExpirationTimeSeconds.Cmp(orderWatcher.MaxExpirationTime()) == -1, "remaining order has an expiration time of %s which is *greater than* the maximum of %s", order.SignedOrder.ExpirationTimeSeconds, orderWatcher.MaxExpirationTime())
	}
}

func TestOrderWatcherBatchEmitsAddedEvents(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, meshDB)

	// Subscribe to OrderWatcher
	orderEventsChan := make(chan []*zeroex.OrderEvent, 10)
	orderWatcher.Subscribe(orderEventsChan)

	signedOrders := []*zeroex.SignedOrder{}
	for i := 0; i < 2; i++ {
		signedOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, big.NewInt(1000), big.NewInt(1000))
		signedOrders = append(signedOrders, signedOrder)
	}

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	validationResults, err := orderWatcher.ValidateAndStoreValidOrders(ctx, signedOrders, false, constants.TestChainID)
	require.Len(t, validationResults.Rejected, 0)
	require.NoError(t, err)

	orderEvents := <-orderEventsChan
	require.Len(t, orderEvents, 2)
	for _, orderEvent := range orderEvents {
		assert.Equal(t, zeroex.ESOrderAdded, orderEvent.EndState)
	}

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
	require.NoError(t, err)
	require.Len(t, orders, 2)
}

// Scenario: An order is sent to OrderWatcher for validation and storage. The order validation `eth_call`
// takes a long time, during which a fill for this order was made. Once the order validation completes,
// we expect to process blocks missed during validation and for the fill to get emitted.
func TestOrderWatcherValidateAndStoreValidOrdersHighLatencyEventsCatchup(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	// Set up test and orderWatcher
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer func() {
		cancel()
	}()
	blockingChan := make(chan struct{})
	blockingCallEthRPCClient, err := NewBlockingCallEthRPCClient(ethRPCClient, blockingChan)
	require.NoError(t, err)
	blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, blockingCallEthRPCClient, meshDB)

	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	signedOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)

	// Kick off adding order to OrderWatcher within a separate go-routine so that we can block on the validation
	// `eth_call`, simulating network latency, and while it's blocked, fill the order and process the block containing
	// the fill
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		watchOrder(ctx, t, orderWatcher, blockWatcher, ethClient, signedOrder)
	}()

	// Partially fill order
	opts := &bind.TransactOpts{
		From:   takerAddress,
		Signer: scenario.GetTestSignerFn(takerAddress),
	}
	orderWithoutExchangeAddress := signedOrder.ConvertToOrderWithoutExchangeAddress()
	partialFillAmount := big.NewInt(0).Div(wethAmount, big.NewInt(2))
	txn, err := exchange.FillOrder(opts, orderWithoutExchangeAddress, partialFillAmount, signedOrder.Signature)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)

	// Make sure no order events were emitted
	select {
	case _ = <-orderEventsChan:
		t.Fatalf("No order events expected since this order isn't being watched yet")
	case <-time.After(250 * time.Millisecond):
		// noop
	}

	// Unblock validation `eth_call` twice. The first validation of the order, and then the
	// subsequent re-validation of the order necessary because of the missed block events
	blockingChan <- struct{}{}
	blockingChan <- struct{}{}

	// Wait for adding order to OrderWatcher to complete
	wg.Wait()

	// Expect added AND fill event for order to be emitted, due to the OrderWatcher
	// backfilling missed block events past while it was validating the order
	orderEvents := waitForOrderEvents(t, orderEventsChan, 2, 4*time.Second)
	firstOrderEvent := orderEvents[0]
	assert.Equal(t, zeroex.ESOrderAdded, firstOrderEvent.EndState)
	secondOrderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, secondOrderHash, firstOrderEvent.OrderHash)

	secondOrderEvent := orderEvents[1]
	assert.Equal(t, zeroex.ESOrderFilled, secondOrderEvent.EndState)
	secondOrderHash, err = signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, secondOrderHash, secondOrderEvent.OrderHash)
}

// Scenario: While we validate orders at a specific block height, the block at that height gets re-orged out.
// We are no longer certain if we validated the orders against the right block so we reject the orders.
func TestOrderWatcherValidateAndStoreValidOrdersHighLatencyValidationBlockReorgedOut(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	// Set up test and orderWatcher
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer func() {
		cancel()
	}()
	blockingChan := make(chan struct{})
	blockingCallEthRPCClient, err := NewBlockingCallEthRPCClient(ethRPCClient, blockingChan)
	require.NoError(t, err)
	blockWatcher, orderWatcher := setupOrderWatcher(ctx, t, blockingCallEthRPCClient, meshDB)

	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	signedOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	err = blockWatcher.SyncToLatestBlock()
	require.NoError(t, err)
	// We wait 200ms to allow time for the newly mined blocks to get processed by the OrderWatcher, so that when
	// we try to validateAndStore the order below, the validation is performed at the latest block number
	time.Sleep(200 * time.Millisecond)

	// Kick off adding order to OrderWatcher within a separate go-routine so that we can block on the validation
	// `eth_call`, simulating network latency, and while it's blocked, fill the order and process the block containing
	// the fill
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		validationResults, err := orderWatcher.ValidateAndStoreValidOrders(ctx, []*zeroex.SignedOrder{signedOrder}, false, constants.TestChainID)
		require.NoError(t, err)
		require.Len(t, validationResults.Rejected, 1)
		require.Len(t, validationResults.Accepted, 0)
		assert.Equal(t, ordervalidator.ROEthRPCRequestFailed, validationResults.Rejected[0].Status)
	}()

	// Simulate a block re-org by replacing the latest block with a new one, and adding one other order
	// ontop of it
	storedBlocks, err := meshDB.FindAllMiniHeadersSortedByNumber()
	require.NoError(t, err)
	latestBlock := storedBlocks[len(storedBlocks)-1]
	replacementBlockHash := common.HexToHash("0x123456789")
	replacementBlock := &miniheader.MiniHeader{
		Parent:    storedBlocks[len(storedBlocks)-2].Hash,
		Hash:      replacementBlockHash,
		Number:    latestBlock.Number,
		Logs:      []types.Log{},
		Timestamp: time.Now().Add(-12 * time.Second),
	}
	newLatestBlock := &miniheader.MiniHeader{
		Parent:    replacementBlockHash,
		Hash:      common.HexToHash("0x987654321"),
		Number:    big.NewInt(0).Add(latestBlock.Number, big.NewInt(1)),
		Logs:      []types.Log{},
		Timestamp: time.Now().Add(-1 * time.Second),
	}
	err = meshDB.MiniHeaders.Delete(latestBlock.ID())
	require.NoError(t, err)
	err = meshDB.MiniHeaders.Insert(replacementBlock)
	require.NoError(t, err)
	err = meshDB.MiniHeaders.Insert(newLatestBlock)
	require.NoError(t, err)

	// Unblock validation `eth_call`
	blockingChan <- struct{}{}

	// Wait for adding order to OrderWatcher to complete
	wg.Wait()
}

func setupOrderWatcherScenario(ctx context.Context, t *testing.T, ethClient *ethclient.Client, meshDB *meshdb.MeshDB, signedOrder *zeroex.SignedOrder) (*blockwatch.Watcher, chan []*zeroex.OrderEvent) {
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
	require.Len(t, validationResults.Accepted, 1, "Expected order to pass validation and get added to OrderWatcher")
}

func setupOrderWatcher(ctx context.Context, t *testing.T, ethRPCClient ethrpcclient.Client, meshDB *meshdb.MeshDB) (*blockwatch.Watcher, *Watcher) {
	blockWatcherClient, err := blockwatch.NewRpcClient(ethRPCClient)
	require.NoError(t, err)
	topics := GetRelevantTopics()
	stack := simplestack.New(blockWatcherRetentionLimit, []*miniheader.MiniHeader{})
	blockWatcherConfig := blockwatch.Config{
		Stack:           stack,
		PollingInterval: blockPollingInterval,
		WithLogs:        true,
		Topics:          topics,
		Client:          blockWatcherClient,
	}
	blockWatcher := blockwatch.New(blockWatcherConfig)
	orderValidator, err := ordervalidator.New(ethRPCClient, constants.TestChainID, ethereumRPCMaxContentLength)
	require.NoError(t, err)
	orderWatcher, err := New(Config{
		MeshDB:            meshDB,
		BlockWatcher:      blockWatcher,
		OrderValidator:    orderValidator,
		ChainID:           constants.TestChainID,
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
	storedBlocks, err := meshDB.FindAllMiniHeadersSortedByNumber()
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

func waitTxnSuccessfullyMined(t *testing.T, ethClient *ethclient.Client, txn *types.Transaction) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancelFn()
	receipt, err := bind.WaitMined(ctx, ethClient, txn)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, uint64(1))
}
