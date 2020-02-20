// +build !js

package orderwatch

import (
	"context"
	"flag"
	"math/big"
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
	"github.com/ethereum/go-ethereum/rpc"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
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
	testing.Init()
	flag.Parse()
}

func init() {
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

func TestOrderWatcherUnfundedInsufficientERC20BalanceForMakerFee(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	wethFeeAmount := new(big.Int).Mul(big.NewInt(5), eighteenDecimalsInBaseUnits)
	signedOrder := scenario.CreateNFTForZRXWithWETHMakerFeeSignedTestOrder(t, ethClient, makerAddress, takerAddress, tokenID, zrxAmount, wethFeeAmount)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	blockWatcher, orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Transfer makerAsset out of maker address
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
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

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
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
		Value:  big.NewInt(100000000000000000),
	}
	trimmedOrder := signedOrder.Trim()
	txn, err := exchange.FillOrder(opts, trimmedOrder, wethAmount, signedOrder.Signature)
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
		Value:  big.NewInt(100000000000000000),
	}
	trimmedOrder := signedOrder.Trim()
	halfAmount := new(big.Int).Div(wethAmount, big.NewInt(2))
	txn, err := exchange.FillOrder(opts, trimmedOrder, halfAmount, signedOrder.Signature)
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

func TestOrderWatcherCleanup(t *testing.T) {
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

	// Create and add two orders to OrderWatcher
	amount := big.NewInt(10000)
	signedOrderOne := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, amount, amount)
	watchOrder(ctx, t, orderWatcher, blockWatcher, ethClient, signedOrderOne)
	signedOrderTwo := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, amount, amount)
	watchOrder(ctx, t, orderWatcher, blockWatcher, ethClient, signedOrderTwo)
	signedOrderOneHash, err := signedOrderTwo.ComputeOrderHash()
	require.NoError(t, err)

	// Set lastUpdate for signedOrderOne to more than defaultLastUpdatedBuffer so that signedOrderOne
	// does not get re-validated by the cleanup job
	signedOrderOneDB := &meshdb.Order{}
	err = meshDB.Orders.FindByID(signedOrderOneHash.Bytes(), signedOrderOneDB)
	require.NoError(t, err)
	signedOrderOneDB.LastUpdated = time.Now().Add(-defaultLastUpdatedBuffer - 1*time.Minute)
	err = meshDB.Orders.Update(signedOrderOneDB)
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

func TestOrderWatcherUpdateBlockHeadersStoredInDBHeaderExists(t *testing.T) {
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)

	headerOne := &miniheader.MiniHeader{
		Number:    big.NewInt(5),
		Hash:      common.HexToHash("0x293b9ea024055a3e9eddbf9b9383dc7731744111894af6aa038594dc1b61f87f"),
		Parent:    common.HexToHash("0x26b13ac89500f7fcdd141b7d1b30f3a82178431eca325d1cf10998f9d68ff5ba"),
		Timestamp: time.Now().UTC(),
	}

	testCases := []struct {
		events              []*blockwatch.Event
		startMiniHeaders    []*miniheader.MiniHeader
		expectedMiniHeaders []*miniheader.MiniHeader
	}{
		// Scenario 1: Header 1 exists in DB. Get's removed and then re-added.
		{
			events: []*blockwatch.Event{
				&blockwatch.Event{
					Type:        blockwatch.Removed,
					BlockHeader: headerOne,
				},
				&blockwatch.Event{
					Type:        blockwatch.Added,
					BlockHeader: headerOne,
				},
			},
			startMiniHeaders: []*miniheader.MiniHeader{
				headerOne,
			},
			expectedMiniHeaders: []*miniheader.MiniHeader{
				headerOne,
			},
		},
		// Scenario 2: Header doesn't exist, get's added and then removed
		{
			events: []*blockwatch.Event{
				&blockwatch.Event{
					Type:        blockwatch.Added,
					BlockHeader: headerOne,
				},
				&blockwatch.Event{
					Type:        blockwatch.Removed,
					BlockHeader: headerOne,
				},
			},
			startMiniHeaders:    []*miniheader.MiniHeader{},
			expectedMiniHeaders: []*miniheader.MiniHeader{},
		},
		// Scenario 3: Header added, removed then re-added
		{
			events: []*blockwatch.Event{
				&blockwatch.Event{
					Type:        blockwatch.Added,
					BlockHeader: headerOne,
				},
				&blockwatch.Event{
					Type:        blockwatch.Removed,
					BlockHeader: headerOne,
				},
				&blockwatch.Event{
					Type:        blockwatch.Added,
					BlockHeader: headerOne,
				},
			},
			startMiniHeaders: []*miniheader.MiniHeader{},
			expectedMiniHeaders: []*miniheader.MiniHeader{
				headerOne,
			},
		},
		// Scenario 4: Header removed, added then removed again
		{
			events: []*blockwatch.Event{
				&blockwatch.Event{
					Type:        blockwatch.Removed,
					BlockHeader: headerOne,
				},
				&blockwatch.Event{
					Type:        blockwatch.Added,
					BlockHeader: headerOne,
				},
				&blockwatch.Event{
					Type:        blockwatch.Removed,
					BlockHeader: headerOne,
				},
			},
			startMiniHeaders: []*miniheader.MiniHeader{
				headerOne,
			},
			expectedMiniHeaders: []*miniheader.MiniHeader{},
		},
		// Scenario 5: Call added twice for the same block
		{
			events: []*blockwatch.Event{
				&blockwatch.Event{
					Type:        blockwatch.Added,
					BlockHeader: headerOne,
				},
				&blockwatch.Event{
					Type:        blockwatch.Added,
					BlockHeader: headerOne,
				},
			},
			startMiniHeaders: []*miniheader.MiniHeader{},
			expectedMiniHeaders: []*miniheader.MiniHeader{
				headerOne,
			},
		},
		// Scenario 6: Call removed twice for the same block
		{
			events: []*blockwatch.Event{
				&blockwatch.Event{
					Type:        blockwatch.Removed,
					BlockHeader: headerOne,
				},
				&blockwatch.Event{
					Type:        blockwatch.Removed,
					BlockHeader: headerOne,
				},
			},
			startMiniHeaders: []*miniheader.MiniHeader{
				headerOne,
			},
			expectedMiniHeaders: []*miniheader.MiniHeader{},
		},
	}

	for _, testCase := range testCases {
		for _, startMiniHeader := range testCase.startMiniHeaders {
			err = meshDB.MiniHeaders.Insert(startMiniHeader)
			require.NoError(t, err)
		}

		miniHeadersColTxn := meshDB.MiniHeaders.OpenTransaction()
		defer func() {
			_ = miniHeadersColTxn.Discard()
		}()

		err = updateBlockHeadersStoredInDB(miniHeadersColTxn, testCase.events)
		require.NoError(t, err)

		err = miniHeadersColTxn.Commit()
		require.NoError(t, err)

		miniHeaders := []*miniheader.MiniHeader{}
		err = meshDB.MiniHeaders.FindAll(&miniHeaders)
		require.NoError(t, err)
		assert.Equal(t, testCase.expectedMiniHeaders, miniHeaders)

		err := meshDB.ClearAllMiniHeaders()
		require.NoError(t, err)
	}
}

func TestOrderWatcherHandleOrderExpirationsExpired(t *testing.T) {
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
	signedOrderOne := scenario.CreateSignedTestOrderWithExpirationTime(t, ethClient, makerAddress, takerAddress, expirationTime)
	signedOrderTwo := scenario.CreateSignedTestOrderWithExpirationTime(t, ethClient, makerAddress, takerAddress, expirationTime)
	blockwatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, meshDB)
	watchOrder(ctx, t, orderWatcher, blockwatcher, ethClient, signedOrderOne)
	watchOrder(ctx, t, orderWatcher, blockwatcher, ethClient, signedOrderTwo)

	signedOrderOneHash, err := signedOrderOne.ComputeOrderHash()
	require.NoError(t, err)
	var orderOne meshdb.Order
	err = meshDB.Orders.FindByID(signedOrderOneHash.Bytes(), &orderOne)
	require.NoError(t, err)
	// Since we flag SignedOrderOne for revalidation, we expect `handleOrderExpirations` not to return an
	// expiry event for it.
	ordersToRevalidate := map[common.Hash]*meshdb.Order{
		signedOrderOneHash: &orderOne,
	}

	ordersColTxn := meshDB.Orders.OpenTransaction()
	defer func() {
		_ = ordersColTxn.Discard()
	}()

	previousLatestBlockTimestamp := expirationTime.Add(-1 * time.Minute)
	latestBlockTimestamp := expirationTime.Add(1 * time.Second)
	orderEvents, err := orderWatcher.handleOrderExpirations(ordersColTxn, latestBlockTimestamp, previousLatestBlockTimestamp, ordersToRevalidate)
	require.NoError(t, err)

	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	signedOrderTwoHash, err := signedOrderTwo.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, signedOrderTwoHash, orderEvent.OrderHash)
	assert.Equal(t, zeroex.ESOrderExpired, orderEvent.EndState)
	assert.Equal(t, big.NewInt(0), orderEvent.FillableTakerAssetAmount)
	assert.Len(t, orderEvent.ContractEvents, 0)

	err = ordersColTxn.Commit()
	require.NoError(t, err)

	var orderTwo meshdb.Order
	err = meshDB.Orders.FindByID(signedOrderTwoHash.Bytes(), &orderTwo)
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
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer func() {
		cancel()
	}()

	// Create and add an order (which will later become expired) to OrderWatcher
	expirationTime := time.Now().Add(24 * time.Hour)
	signedOrderOne := scenario.CreateSignedTestOrderWithExpirationTime(t, ethClient, makerAddress, takerAddress, expirationTime)
	signedOrderTwo := scenario.CreateSignedTestOrderWithExpirationTime(t, ethClient, makerAddress, takerAddress, expirationTime)
	blockwatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, meshDB)
	watchOrder(ctx, t, orderWatcher, blockwatcher, ethClient, signedOrderOne)
	watchOrder(ctx, t, orderWatcher, blockwatcher, ethClient, signedOrderTwo)

	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// Simulate a block found with a timestamp past expirationTime
	latestBlock, err := meshDB.FindLatestMiniHeader()
	require.NoError(t, err)
	blockTimestamp := expirationTime.Add(1 * time.Minute)
	nextBlock := &miniheader.MiniHeader{
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
	orderEvents := waitForOrderEvents(t, orderEventsChan, 2, 4*time.Second)
	require.Len(t, orderEvents, 2)
	for _, orderEvent := range orderEvents {
		assert.Equal(t, zeroex.ESOrderExpired, orderEvent.EndState)
	}

	signedOrderOneHash, err := signedOrderOne.ComputeOrderHash()
	require.NoError(t, err)
	var orderOne meshdb.Order
	err = meshDB.Orders.FindByID(signedOrderOneHash.Bytes(), &orderOne)
	require.NoError(t, err)
	// Since we flag SignedOrderOne for revalidation, we expect `handleOrderExpirations` not to return an
	// unexpiry event for it.
	ordersToRevalidate := map[common.Hash]*meshdb.Order{
		signedOrderOneHash: &orderOne,
	}

	ordersColTxn := meshDB.Orders.OpenTransaction()
	defer func() {
		_ = ordersColTxn.Discard()
	}()

	// LatestBlockTimestamp is earlier than previous latest simulating block-reorg where new latest block
	// has an earlier timestamp than the last
	previousLatestBlockTimestamp := blockTimestamp
	latestBlockTimestamp := expirationTime.Add(-1 * time.Minute)
	orderEvents, err = orderWatcher.handleOrderExpirations(ordersColTxn, latestBlockTimestamp, previousLatestBlockTimestamp, ordersToRevalidate)
	require.NoError(t, err)

	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	signedOrderTwoHash, err := signedOrderTwo.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, signedOrderTwoHash, orderEvent.OrderHash)
	assert.Equal(t, zeroex.ESOrderUnexpired, orderEvent.EndState)
	assert.Equal(t, signedOrderTwo.TakerAssetAmount, orderEvent.FillableTakerAssetAmount)
	assert.Len(t, orderEvent.ContractEvents, 0)

	err = ordersColTxn.Commit()
	require.NoError(t, err)

	var orderTwo meshdb.Order
	err = meshDB.Orders.FindByID(signedOrderTwoHash.Bytes(), &orderTwo)
	require.NoError(t, err)
	assert.Equal(t, false, orderTwo.IsRemoved)
}

func TestOrderWatcherMaintainMiniHeaderRetentionLimit(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	// Set up test and orderWatcher
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	err = meshDB.UpdateMiniHeaderRetentionLimit(miniHeaderRetentionLimit)
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer func() {
		cancel()
	}()
	_, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, meshDB)

	latestMiniHeader, err := meshDB.FindLatestMiniHeader()
	require.NoError(t, err)

	headerOne := &miniheader.MiniHeader{
		Number:    big.NewInt(0).Add(latestMiniHeader.Number, big.NewInt(1)),
		Hash:      common.HexToHash("0x293b9ea024055a3e9eddbf9b9383dc7731744111894af6aa038594dc1b61f87f"),
		Parent:    common.HexToHash("0x26b13ac89500f7fcdd141b7d1b30f3a82178431eca325d1cf10998f9d68ff5ba"),
		Timestamp: time.Now().UTC(),
	}
	headerTwo := &miniheader.MiniHeader{
		Number:    big.NewInt(0).Add(headerOne.Number, big.NewInt(1)),
		Hash:      common.HexToHash("0x72ca9481b09b8c00b2c38575e5652f2de1077f1676c6b868cf575229fcb06a96"),
		Parent:    common.HexToHash("0x293b9ea024055a3e9eddbf9b9383dc7731744111894af6aa038594dc1b61f87f"),
		Timestamp: time.Now().UTC(),
	}
	headerThree := &miniheader.MiniHeader{
		Number:    big.NewInt(0).Add(headerTwo.Number, big.NewInt(1)),
		Hash:      common.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
		Parent:    common.HexToHash("0x72ca9481b09b8c00b2c38575e5652f2de1077f1676c6b868cf575229fcb06a96"),
		Timestamp: time.Now().UTC(),
	}

	blockEvents := []*blockwatch.Event{
		&blockwatch.Event{
			Type:        blockwatch.Added,
			BlockHeader: headerOne,
		},
		&blockwatch.Event{
			Type:        blockwatch.Added,
			BlockHeader: headerTwo,
		},
		&blockwatch.Event{
			Type:        blockwatch.Added,
			BlockHeader: headerThree,
		},
	}
	err = orderWatcher.handleBlockEvents(ctx, blockEvents)
	require.NoError(t, err)

	latestMiniHeader, err = meshDB.FindLatestMiniHeader()
	require.NoError(t, err)
	assert.Equal(t, headerThree.Hash, latestMiniHeader.Hash)

	totalMiniHeaders, err := meshDB.MiniHeaders.Count()
	require.NoError(t, err)
	assert.Equal(t, meshDB.MiniHeaderRetentionLimit, totalMiniHeaders)
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
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer func() {
		cancel()
	}()

	// Create and add an order (which will later become expired) to OrderWatcher
	expirationTime := time.Now().Add(24 * time.Hour)
	signedOrderOne := scenario.CreateSignedTestOrderWithExpirationTime(t, ethClient, makerAddress, takerAddress, expirationTime)
	blockwatcher, orderWatcher := setupOrderWatcher(ctx, t, ethRPCClient, meshDB)
	watchOrder(ctx, t, orderWatcher, blockwatcher, ethClient, signedOrderOne)

	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// Simulate a block found with a timestamp past expirationTime. This will mark the order as removed
	// and will remove it from the expiration watcher.
	latestBlock, err := meshDB.FindLatestMiniHeader()
	require.NoError(t, err)
	blockTimestamp := expirationTime.Add(1 * time.Minute)
	nextBlock := &miniheader.MiniHeader{
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

	signedOrderOneHash, err := signedOrderOne.ComputeOrderHash()
	require.NoError(t, err)
	var orderOne meshdb.Order
	err = meshDB.Orders.FindByID(signedOrderOneHash.Bytes(), &orderOne)
	require.NoError(t, err)

	ordersColTxn := meshDB.Orders.OpenTransaction()
	defer func() {
		_ = ordersColTxn.Discard()
	}()

	validationResults := ordervalidator.ValidationResults{
		Accepted: []*ordervalidator.AcceptedOrderInfo{
			&ordervalidator.AcceptedOrderInfo{
				OrderHash:                signedOrderOneHash,
				SignedOrder:              signedOrderOne,
				FillableTakerAssetAmount: big.NewInt(1).Div(signedOrderOne.TakerAssetAmount, big.NewInt(2)),
				IsNew:                    false,
			},
		},
		Rejected: []*ordervalidator.RejectedOrderInfo{},
	}
	orderHashToDBOrder := map[common.Hash]*meshdb.Order{
		signedOrderOneHash: &orderOne,
	}
	exchangeFillEvent := "ExchangeFillEvent"
	orderHashToEvents := map[common.Hash][]*zeroex.ContractEvent{
		signedOrderOneHash: []*zeroex.ContractEvent{
			&zeroex.ContractEvent{
				Kind: exchangeFillEvent,
			},
		},
	}
	validationBlockTimestamp := expirationTime.Add(-1 * time.Minute)
	orderEvents, err = orderWatcher.convertValidationResultsIntoOrderEvents(ordersColTxn, &validationResults, orderHashToDBOrder, orderHashToEvents, validationBlockTimestamp)
	require.NoError(t, err)

	require.Len(t, orderEvents, 2)
	orderEventTwo := orderEvents[0]
	assert.Equal(t, signedOrderOneHash, orderEventTwo.OrderHash)
	assert.Equal(t, zeroex.ESOrderUnexpired, orderEventTwo.EndState)
	assert.Len(t, orderEventTwo.ContractEvents, 0)
	orderEventOne := orderEvents[1]
	assert.Equal(t, signedOrderOneHash, orderEventOne.OrderHash)
	assert.Equal(t, zeroex.ESOrderFilled, orderEventOne.EndState)
	assert.Len(t, orderEventOne.ContractEvents, 1)
	assert.Equal(t, orderEventOne.ContractEvents[0].Kind, exchangeFillEvent)

	err = ordersColTxn.Commit()
	require.NoError(t, err)

	var orderTwo meshdb.Order
	err = meshDB.Orders.FindByID(signedOrderOneHash.Bytes(), &orderTwo)
	require.NoError(t, err)
	assert.Equal(t, false, orderTwo.IsRemoved)
}

func TestDrainAllBlockEventsChan(t *testing.T) {
	blockEventsChan := make(chan []*blockwatch.Event, 100)
	ts := time.Now().Add(1 * time.Hour)
	blockEventsOne := []*blockwatch.Event{
		&blockwatch.Event{
			Type: blockwatch.Added,
			BlockHeader: &miniheader.MiniHeader{
				Parent:    common.HexToHash("0x0"),
				Hash:      common.HexToHash("0x1"),
				Number:    big.NewInt(1),
				Timestamp: ts,
			},
		},
	}
	blockEventsChan <- blockEventsOne

	blockEventsTwo := []*blockwatch.Event{
		&blockwatch.Event{
			Type: blockwatch.Added,
			BlockHeader: &miniheader.MiniHeader{
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
	stack := simplestack.New(meshDB.MiniHeaderRetentionLimit, []*miniheader.MiniHeader{})
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
