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
	"github.com/0xProject/0x-mesh/ethereum/dbstack"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/0xProject/0x-mesh/meshdb"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/0xProject/0x-mesh/zeroex/ordervalidator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	blockchainLifecycle, err = ethereum.NewBlockchainLifecycle(rpcClient)
	if err != nil {
		panic(err)
	}
	ethClient = ethclient.NewClient(rpcClient)
	ganacheAddresses := ethereum.NetworkIDToContractAddresses[constants.TestNetworkID]
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
	orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Transfer makerAsset out of maker address
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	txn, err := zrx.Transfer(opts, constants.GanacheAccount4, zrxAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

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
	orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Transfer makerAsset out of maker address
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	txn, err := dummyERC721Token.TransferFrom(opts, makerAddress, constants.GanacheAccount4, tokenID)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

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
	orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Remove Maker's NFT approval to ERC721Proxy
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	ganacheAddresses := ethereum.NetworkIDToContractAddresses[constants.TestNetworkID]
	txn, err := dummyERC721Token.SetApprovalForAll(opts, ganacheAddresses.ERC721Proxy, false)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

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
	orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Remove Maker's ERC1155 approval to ERC1155Proxy
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	ganacheAddresses := ethereum.NetworkIDToContractAddresses[constants.TestNetworkID]
	txn, err := erc1155Mintable.SetApprovalForAll(opts, ganacheAddresses.ERC1155Proxy, false)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

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
	orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Reduce Maker's ERC1155 balance
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	txn, err := erc1155Mintable.SafeTransferFrom(opts, makerAddress, constants.GanacheAccount4, tokenID, erc1155FungibleAmount, []byte{})
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

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

	ganacheAddresses := ethereum.NetworkIDToContractAddresses[constants.TestNetworkID]
	signedOrder := scenario.CreateZRXForWETHSignedTestOrder(t, ethClient, makerAddress, takerAddress, wethAmount, zrxAmount)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Remove Maker's ZRX approval to ERC20Proxy
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	txn, err := zrx.Approve(opts, ganacheAddresses.ERC20Proxy, big.NewInt(0))
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

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
	orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Transfer makerAsset out of maker address
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	txn, err := zrx.Transfer(opts, constants.GanacheAccount4, zrxAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

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
	setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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
	orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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
	orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Cancel order
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	orderWithoutExchangeAddress := signedOrder.ConvertToOrderWithoutExchangeAddress()
	txn, err := exchange.CancelOrder(opts, orderWithoutExchangeAddress)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

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
	orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Cancel order with epoch
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	targetOrderEpoch := signedOrder.Salt
	txn, err := exchange.CancelOrdersUpTo(opts, targetOrderEpoch)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

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
	orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Fill order
	opts := &bind.TransactOpts{
		From:   takerAddress,
		Signer: scenario.GetTestSignerFn(takerAddress),
	}
	orderWithoutExchangeAddress := signedOrder.ConvertToOrderWithoutExchangeAddress()
	txn, err := exchange.FillOrder(opts, orderWithoutExchangeAddress, wethAmount, signedOrder.Signature)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

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
	orderEventsChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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

func TestOrderWatcherDecreaseExpirationTime(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	// Set up test and orderWatcher. Manually change maxOrders.
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	meshDB, err := meshdb.New("/tmp/leveldb_testing/" + uuid.New().String())
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer func() {
		cancel()
	}()
	orderWatcher := setupOrderWatcher(ctx, t, ethClient, meshDB)
	orderWatcher.maxOrders = 20

	// create and watch maxOrders orders
	for i := 0; i < orderWatcher.maxOrders; i++ {
		signedOrder := scenario.CreateSignedTestOrderWithExpirationTime(t, ethClient, makerAddress, takerAddress, time.Now().Add(10*time.Minute+time.Duration(i)*time.Minute))
		watchOrder(t, orderWatcher, signedOrder)
	}

	// We don't care about the order events above for the purposes of this test,
	// so we only subscribe now.
	orderEventsChan := make(chan []*zeroex.OrderEvent, 2*orderWatcher.maxOrders)
	orderWatcher.Subscribe(orderEventsChan)

	// The next order should cause some orders to be removed and the appropriate
	// events to fire.
	signedOrder := scenario.CreateSignedTestOrderWithExpirationTime(t, ethClient, makerAddress, takerAddress, time.Now().Add(10*time.Minute+1*time.Second))
	watchOrder(t, orderWatcher, signedOrder)
	expectedOrderEvents := int(float64(orderWatcher.maxOrders)*(1-maxOrdersTrimRatio)) + 1
	orderEvents := waitForOrderEvents(t, orderEventsChan, expectedOrderEvents, 4*time.Second)
	require.Len(t, orderEvents, expectedOrderEvents, "wrong number of order events were fired")
	for i, orderEvent := range orderEvents {
		// Last event should be ADDED. The other events should be REMOVED.
		if i == expectedOrderEvents-1 {
			assert.Equal(t, zeroex.ESOrderAdded, orderEvent.EndState, "order event %d had wrong EndState", i)
		} else {
			// For REMOVED events, we also make sure that the expiration time is after
			// the current max expiration time.
			assert.Equal(t, zeroex.ESOrderRemoved, orderEvent.EndState, "order event %d had wrong EndState", i)
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

func setupOrderWatcherScenario(ctx context.Context, t *testing.T, ethClient *ethclient.Client, meshDB *meshdb.MeshDB, signedOrder *zeroex.SignedOrder) chan []*zeroex.OrderEvent {
	orderWatcher := setupOrderWatcher(ctx, t, ethClient, meshDB)

	// Start watching an order
	watchOrder(t, orderWatcher, signedOrder)

	// Subscribe to OrderWatcher
	orderEventsChan := make(chan []*zeroex.OrderEvent, 10)
	orderWatcher.Subscribe(orderEventsChan)

	return orderEventsChan
}

func watchOrder(t *testing.T, orderWatcher *Watcher, signedOrder *zeroex.SignedOrder) {
	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	orderInfo := &ordervalidator.AcceptedOrderInfo{
		SignedOrder:              signedOrder,
		OrderHash:                orderHash,
		FillableTakerAssetAmount: signedOrder.TakerAssetAmount,
		IsNew:                    true,
	}
	err = orderWatcher.Add(orderInfo)
	require.NoError(t, err)
}

func setupOrderWatcher(ctx context.Context, t *testing.T, ethClient *ethclient.Client, meshDB *meshdb.MeshDB) *Watcher {
	// Init OrderWatcher
	blockWatcherClient, err := blockwatch.NewRpcClient(constants.GanacheEndpoint, ethereumRPCRequestTimeout)
	require.NoError(t, err)
	topics := GetRelevantTopics()
	stack := dbstack.New(meshDB, blockWatcherRetentionLimit)
	blockWatcherConfig := blockwatch.Config{
		Stack:           stack,
		PollingInterval: blockPollingInterval,
		StartBlockDepth: ethrpc.LatestBlockNumber,
		WithLogs:        true,
		Topics:          topics,
		Client:          blockWatcherClient,
	}
	blockWatcher := blockwatch.New(blockWatcherConfig)
	orderValidator, err := ordervalidator.New(ethClient, constants.TestNetworkID, ethereumRPCMaxContentLength, 0)
	require.NoError(t, err)
	orderWatcher, err := New(Config{
		MeshDB:           meshDB,
		BlockWatcher:     blockWatcher,
		OrderValidator:   orderValidator,
		NetworkID:        constants.TestNetworkID,
		ExpirationBuffer: 0,
	})
	require.NoError(t, err)

	// Start the block watcher.
	go func() {
		err := blockWatcher.Watch(ctx)
		require.NoError(t, err)
	}()

	// Start OrderWatcher
	go func() {
		err := orderWatcher.Watch(ctx)
		require.NoError(t, err)
	}()

	return orderWatcher
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
