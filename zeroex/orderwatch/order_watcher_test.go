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

var makerAddress = constants.GanacheAccount1
var takerAddress = constants.GanacheAccount2
var eighteenDecimalsInBaseUnits = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
var wethAmount = new(big.Int).Mul(big.NewInt(50), eighteenDecimalsInBaseUnits)
var zrxAmount = new(big.Int).Mul(big.NewInt(100), eighteenDecimalsInBaseUnits)
var tokenID = big.NewInt(1)

// Since these tests must be run sequentially, we don't want them to run as part of
// the normal testing process. They will only be run if the "--serial" flag is used.
var serialTestsEnabled bool

func init() {
	flag.BoolVar(&serialTestsEnabled, "serial", false, "enable serial tests")
	flag.Parse()
}

var rpcClient *ethrpc.Client
var ethClient *ethclient.Client
var zrx *wrappers.ZRXToken
var dummyERC721Token *wrappers.DummyERC721Token
var exchange *wrappers.Exchange
var weth *wrappers.WETH9

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
	orderEventChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Transfer makerAsset out of maker address
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	txn, err := zrx.Transfer(opts, constants.GanacheAccount4, zrxAmount)
	require.NoError(t, err)
	ctx, cancelFn = context.WithTimeout(context.Background(), 4*time.Second)
	receipt, err := bind.WaitMined(ctx, ethClient, txn)
	cancelFn()
	require.NoError(t, err)
	assert.Equal(t, receipt.Status, uint64(1))

	var orderEvents []*zeroex.OrderEvent
	select {
	case orderEvents = <-orderEventChan:
		// continue with the test
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for order event")
	}
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.EKOrderBecameUnfunded, orderEvent.Kind)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
	require.NoError(t, err)
	assert.Len(t, orders, 1)
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
	orderEventChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Transfer makerAsset out of maker address
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	txn, err := dummyERC721Token.TransferFrom(opts, makerAddress, constants.GanacheAccount4, tokenID)
	require.NoError(t, err)
	ctx, cancelFn = context.WithTimeout(context.Background(), 4*time.Second)
	receipt, err := bind.WaitMined(ctx, ethClient, txn)
	cancelFn()
	require.NoError(t, err)
	assert.Equal(t, receipt.Status, uint64(1))

	var orderEvents []*zeroex.OrderEvent
	select {
	case orderEvents = <-orderEventChan:
		// continue with the test
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for order event")
	}
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.EKOrderBecameUnfunded, orderEvent.Kind)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
	require.NoError(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
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
	orderEventChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Remove Maker's NFT approval to ERC721Proxy
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	ganacheAddresses := ethereum.NetworkIDToContractAddresses[constants.TestNetworkID]
	txn, err := dummyERC721Token.SetApprovalForAll(opts, ganacheAddresses.ERC721Proxy, false)
	require.NoError(t, err)
	ctx, cancelFn = context.WithTimeout(context.Background(), 4*time.Second)
	receipt, err := bind.WaitMined(ctx, ethClient, txn)
	cancelFn()
	require.NoError(t, err)
	assert.Equal(t, receipt.Status, uint64(1))

	var orderEvents []*zeroex.OrderEvent
	select {
	case orderEvents = <-orderEventChan:
		// continue with the test
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for order event")
	}
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.EKOrderBecameUnfunded, orderEvent.Kind)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
	require.NoError(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
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
	orderEventChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Remove Maker's ZRX approval to ERC20Proxy
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	txn, err := zrx.Approve(opts, ganacheAddresses.ERC20Proxy, big.NewInt(0))
	require.NoError(t, err)
	ctx, cancelFn = context.WithTimeout(context.Background(), 4*time.Second)
	receipt, err := bind.WaitMined(ctx, ethClient, txn)
	cancelFn()
	require.NoError(t, err)
	assert.Equal(t, receipt.Status, uint64(1))

	var orderEvents []*zeroex.OrderEvent
	select {
	case orderEvents = <-orderEventChan:
		// continue with the test
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for order event")
	}
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.EKOrderBecameUnfunded, orderEvent.Kind)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
	require.NoError(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
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
	orderEventChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Transfer makerAsset out of maker address
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	txn, err := zrx.Transfer(opts, constants.GanacheAccount4, zrxAmount)
	require.NoError(t, err)
	ctx, cancelFn = context.WithTimeout(context.Background(), 4*time.Second)
	receipt, err := bind.WaitMined(ctx, ethClient, txn)
	cancelFn()
	require.NoError(t, err)
	assert.Equal(t, receipt.Status, uint64(1))

	var orderEvents []*zeroex.OrderEvent
	select {
	case orderEvents = <-orderEventChan:
		// continue with the test
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for order event")
	}
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.EKOrderBecameUnfunded, orderEvent.Kind)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
	require.NoError(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)

	// Transfer makerAsset back to maker address
	zrxCoinbase := constants.GanacheAccount0
	opts = &bind.TransactOpts{
		From:   zrxCoinbase,
		Signer: scenario.GetTestSignerFn(zrxCoinbase),
	}
	txn, err = zrx.Transfer(opts, makerAddress, zrxAmount)
	require.NoError(t, err)
	ctx, cancelFn = context.WithTimeout(context.Background(), 4*time.Second)
	receipt, err = bind.WaitMined(ctx, ethClient, txn)
	cancelFn()
	require.NoError(t, err)
	assert.Equal(t, receipt.Status, uint64(1))

	orderEvents = <-orderEventChan
	require.Len(t, orderEvents, 1)
	orderEvent = orderEvents[0]
	assert.Equal(t, zeroex.EKOrderAdded, orderEvent.Kind)

	var newOrders []*meshdb.Order
	err = meshDB.Orders.FindAll(&newOrders)
	require.NoError(t, err)
	require.Len(t, newOrders, 1)
	assert.Equal(t, orderEvent.OrderHash, newOrders[0].Hash)
	assert.Equal(t, false, newOrders[0].IsRemoved)
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
	assert.Len(t, orders, 1)
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
	ctx, cancelFn = context.WithTimeout(context.Background(), 4*time.Second)
	receipt, err := bind.WaitMined(ctx, ethClient, txn)
	cancelFn()
	require.NoError(t, err)
	assert.Equal(t, receipt.Status, uint64(1))

	var newOrders []*meshdb.Order
	err = meshDB.Orders.FindAll(&newOrders)
	require.NoError(t, err)
	require.Len(t, newOrders, 1)
	require.NotEqual(t, dbOrder.LastUpdated, newOrders[0].Hash)
	assert.Equal(t, false, newOrders[0].IsRemoved)
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
	orderEventChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

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
	ctx, cancelFn = context.WithTimeout(context.Background(), 4*time.Second)
	receipt, err := bind.WaitMined(ctx, ethClient, txn)
	cancelFn()
	require.NoError(t, err)
	assert.Equal(t, receipt.Status, uint64(1))

	var orderEvents []*zeroex.OrderEvent
	select {
	case orderEvents = <-orderEventChan:
		// continue with the test
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for order event")
	}
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.EKOrderBecameUnfunded, orderEvent.Kind)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
	require.NoError(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)

	opts = &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
		Value:  wethAmount,
	}
	txn, err = weth.Deposit(opts)
	require.NoError(t, err)
	ctx, cancelFn = context.WithTimeout(context.Background(), 4*time.Second)
	receipt, err = bind.WaitMined(ctx, ethClient, txn)
	cancelFn()
	require.NoError(t, err)
	assert.Equal(t, receipt.Status, uint64(1))

	orderEvents = <-orderEventChan
	require.Len(t, orderEvents, 1)
	orderEvent = orderEvents[0]
	assert.Equal(t, zeroex.EKOrderAdded, orderEvent.Kind)

	var newOrders []*meshdb.Order
	err = meshDB.Orders.FindAll(&newOrders)
	require.NoError(t, err)
	require.Len(t, newOrders, 1)
	assert.Equal(t, orderEvent.OrderHash, newOrders[0].Hash)
	assert.Equal(t, false, newOrders[0].IsRemoved)
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
	orderEventChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Cancel order
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	orderWithoutExchangeAddress := signedOrder.ConvertToOrderWithoutExchangeAddress()
	txn, err := exchange.CancelOrder(opts, orderWithoutExchangeAddress)
	require.NoError(t, err)
	ctx, cancelFn = context.WithTimeout(context.Background(), 4*time.Second)
	receipt, err := bind.WaitMined(ctx, ethClient, txn)
	cancelFn()
	require.NoError(t, err)
	assert.Equal(t, receipt.Status, uint64(1))

	var orderEvents []*zeroex.OrderEvent
	select {
	case orderEvents = <-orderEventChan:
		// continue with the test
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for order event")
	}
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.EKOrderCancelled, orderEvent.Kind)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
	require.NoError(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
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
	orderEventChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Cancel order with epoch
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: scenario.GetTestSignerFn(makerAddress),
	}
	targetOrderEpoch := signedOrder.Salt
	txn, err := exchange.CancelOrdersUpTo(opts, targetOrderEpoch)
	require.NoError(t, err)
	ctx, cancelFn = context.WithTimeout(context.Background(), 4*time.Second)
	receipt, err := bind.WaitMined(ctx, ethClient, txn)
	cancelFn()
	require.NoError(t, err)
	assert.Equal(t, receipt.Status, uint64(1))

	var orderEvents []*zeroex.OrderEvent
	select {
	case orderEvents = <-orderEventChan:
		// continue with the test
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for order event")
	}
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.EKOrderCancelled, orderEvent.Kind)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
	require.NoError(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
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
	orderEventChan := setupOrderWatcherScenario(ctx, t, ethClient, meshDB, signedOrder)

	// Fill order
	opts := &bind.TransactOpts{
		From:   takerAddress,
		Signer: scenario.GetTestSignerFn(takerAddress),
	}
	orderWithoutExchangeAddress := signedOrder.ConvertToOrderWithoutExchangeAddress()
	txn, err := exchange.FillOrder(opts, orderWithoutExchangeAddress, wethAmount, signedOrder.Signature)
	require.NoError(t, err)
	ctx, cancelFn = context.WithTimeout(context.Background(), 4*time.Second)
	receipt, err := bind.WaitMined(ctx, ethClient, txn)
	cancelFn()
	require.NoError(t, err)
	assert.Equal(t, receipt.Status, uint64(1))

	var orderEvents []*zeroex.OrderEvent
	select {
	case orderEvents = <-orderEventChan:
		// continue with the test
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for order event")
	}
	require.Len(t, orderEvents, 1)
	orderEvent := orderEvents[0]
	assert.Equal(t, zeroex.EKOrderFullyFilled, orderEvent.Kind)

	var orders []*meshdb.Order
	err = meshDB.Orders.FindAll(&orders)
	require.NoError(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, orderEvent.OrderHash, orders[0].Hash)
	assert.Equal(t, true, orders[0].IsRemoved)
}

func setupOrderWatcherScenario(ctx context.Context, t *testing.T, ethClient *ethclient.Client, meshDB *meshdb.MeshDB, signedOrder *zeroex.SignedOrder) chan []*zeroex.OrderEvent {
	orderWatcher := setupOrderWatcher(ctx, t, ethClient, meshDB)

	// Start watching an order
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

	// Subscribe to OrderWatcher
	orderEventChan := make(chan []*zeroex.OrderEvent, 10)
	orderWatcher.Subscribe(orderEventChan)

	return orderEventChan
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
	orderWatcher, err := New(meshDB, blockWatcher, orderValidator, constants.TestNetworkID, 0)
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

var blockchainLifecycle *ethereum.BlockchainLifecycle

func setupSubTest(t *testing.T) func(t *testing.T) {
	blockchainLifecycle.Start(t)
	return func(t *testing.T) {
		blockchainLifecycle.Revert(t)
	}
}
