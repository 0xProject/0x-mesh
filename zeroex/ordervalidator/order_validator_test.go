// +build !js

// We currently don't run these tests in WASM because of an issue in Go. See the header of
// eth_watcher_test.go for more details.
package ordervalidator

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/ethrpcclient"
	"github.com/0xProject/0x-mesh/ethereum/ratelimit"
	"github.com/0xProject/0x-mesh/ethereum/signer"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/scenario/orderopts"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const areNewOrders = false

// emptyGetOrderRelevantStatesCallDataByteLength is all the boilerplate ABI encoding required when calling
// `getOrderRelevantStates` that does not include the encoded SignedOrder. By subtracting this amount from the
// calldata length returned from encoding a call to `getOrderRelevantStates` involving a single SignedOrder, we
// get the number of bytes taken up by the SignedOrder alone.
// i.e.: len(`"0x7f46448d0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"`)
const emptyGetOrderRelevantStatesCallDataStringLength = 268

const (
	maxEthRPCRequestsPer24HrUTC = 1000000
	maxEthRPCRequestsPerSeconds = 1000.0
	defaultCheckpointInterval   = 1 * time.Minute
	defaultEthRPCTimeout        = 5 * time.Second
)

var (
	unsupportedAssetData = common.Hex2Bytes("a2cb61b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064")
	malformedAssetData   = []byte("9HJhsAAAAAAAAAAAAAAAAInSSmtMyxtvqiYl")
	malformedSignature   = []byte("9HJhsAAAAAAAAAAAAAAAAInSSmtMyxtvqiYl")
	multiAssetAssetData  = common.Hex2Bytes("94cfcdd7000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000046000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000024f47261b00000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c48000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000204a7cb5fb70000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001800000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000003e90000000000000000000000000000000000000000000000000000000000002711000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000c800000000000000000000000000000000000000000000000000000000000007d10000000000000000000000000000000000000000000000000000000000004e210000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c4800000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
)

// Since these tests must be run sequentially, we don't want them to run as part of
// the normal testing process. They will only be run if the "--serial" flag is used.
var serialTestsEnabled bool

var ganacheAddresses = ethereum.GanacheAddresses

func init() {
	flag.BoolVar(&serialTestsEnabled, "serial", false, "enable serial tests")
	testing.Init()
	flag.Parse()
}

type testCase struct {
	SignedOrder                 *zeroex.SignedOrder
	IsValid                     bool
	ExpectedRejectedOrderStatus RejectedOrderStatus
}

var (
	rpcClient           *ethrpc.Client
	blockchainLifecycle *ethereum.BlockchainLifecycle
	ethClient           *ethclient.Client
	ethRPCClient        ethrpcclient.Client
)

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
	rateLimiter := ratelimit.NewUnlimited()
	rpcClient, err = rpc.Dial(constants.GanacheEndpoint)
	if err != nil {
		panic(err)
	}
	ethClient = ethclient.NewClient(rpcClient)
	ethRPCClient, err = ethrpcclient.New(rpcClient, defaultEthRPCTimeout, rateLimiter)
	if err != nil {
		panic(err)
	}
}

func TestBatchValidateOffChainCases(t *testing.T) {
	var testCases = []testCase{
		testCase{
			SignedOrder:                 scenario.NewSignedTestOrder(t, orderopts.MakerAssetAmount(big.NewInt(0))),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidMakerAssetAmount,
		},
		testCase{
			SignedOrder:                 scenario.NewSignedTestOrder(t, orderopts.TakerAssetAmount(big.NewInt(0))),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidTakerAssetAmount,
		},
		testCase{
			SignedOrder: scenario.NewSignedTestOrder(t, orderopts.MakerAssetData(multiAssetAssetData)),
			IsValid:     true,
		},
		testCase{
			SignedOrder:                 scenario.NewSignedTestOrder(t, orderopts.MakerAssetData(malformedAssetData)),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidMakerAssetData,
		},
		testCase{
			SignedOrder:                 scenario.NewSignedTestOrder(t, orderopts.TakerAssetData(malformedAssetData)),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidTakerAssetData,
		},
		testCase{
			SignedOrder:                 scenario.NewSignedTestOrder(t, orderopts.MakerAssetData(unsupportedAssetData)),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidMakerAssetData,
		},
		testCase{
			SignedOrder:                 scenario.NewSignedTestOrder(t, orderopts.TakerAssetData(unsupportedAssetData)),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidTakerAssetData,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomSignature(t, malformedSignature),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidSignature,
		},
	}

	for _, testCase := range testCases {
		signedOrders := []*zeroex.SignedOrder{
			testCase.SignedOrder,
		}
		orderValidator, err := New(ethClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
		require.NoError(t, err)

		offchainValidOrders, rejectedOrderInfos := orderValidator.BatchOffchainValidation(signedOrders)
		isValid := len(offchainValidOrders) == 1
		assert.Equal(t, testCase.IsValid, isValid, testCase.ExpectedRejectedOrderStatus)
		if !isValid {
			assert.Equal(t, testCase.ExpectedRejectedOrderStatus, rejectedOrderInfos[0].Status)
		}
	}
}

func TestBatchValidateAValidOrder(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	signedOrder := scenario.NewSignedTestOrder(t, orderopts.SetupMakerState(true))
	signedOrders := []*zeroex.SignedOrder{
		signedOrder,
	}

	orderValidator, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	require.NoError(t, err)

	ctx := context.Background()
	latestBlock, err := ethRPCClient.HeaderByNumber(ctx, nil)
	require.NoError(t, err)
	validationResults := orderValidator.BatchValidate(ctx, signedOrders, areNewOrders, latestBlock)
	assert.Len(t, validationResults.Accepted, 1)
	require.Len(t, validationResults.Rejected, 0)
	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, orderHash, validationResults.Accepted[0].OrderHash)
}

func TestBatchOffchainValidateUnsupportedStaticCall(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	// NOTE(jalextowle): This asset data encodes a staticcall to a function called `unsupportedStaticCall`
	makerFeeAssetData := common.Hex2Bytes("c339d10a000000000000000000000000692a70d2e424a56d2c6c27aa97d1a86395877b3a0000000000000000000000000000000000000000000000000000000000000060c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47000000000000000000000000000000000000000000000000000000000000000048b24020700000000000000000000000000000000000000000000000000000000")
	signedTestOrder := scenario.NewSignedTestOrder(t, orderopts.MakerFeeAssetData(makerFeeAssetData))
	signedOrders := []*zeroex.SignedOrder{
		signedTestOrder,
	}

	rateLimiter := ratelimit.NewUnlimited()
	rpcClient, err := rpc.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)
	ethRPCClient, err := ethrpcclient.New(rpcClient, defaultEthRPCTimeout, rateLimiter)
	require.NoError(t, err)

	orderValidator, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	require.NoError(t, err)

	accepted, rejected := orderValidator.BatchOffchainValidation(signedOrders)
	assert.Len(t, accepted, 0)
	require.Len(t, rejected, 1)
	signedTestOrder.ResetHash()
	expectedOrderHash, err := signedTestOrder.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, expectedOrderHash, rejected[0].OrderHash)
	require.Equal(t, ROInvalidMakerFeeAssetData, rejected[0].Status)
}

var checkGasPriceDefaultStaticCallData = common.Hex2Bytes("c339d10a0000000000000000000000002c530e4ecc573f11bd72cf5fdf580d134d25f15f0000000000000000000000000000000000000000000000000000000000000060c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a4700000000000000000000000000000000000000000000000000000000000000004d728f5b700000000000000000000000000000000000000000000000000000000")

var checkGasPriceStaticCallData = common.Hex2Bytes("c339d10a0000000000000000000000002c530e4ecc573f11bd72cf5fdf580d134d25f15f0000000000000000000000000000000000000000000000000000000000000060c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a4700000000000000000000000000000000000000000000000000000000000000024da5b166a000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000")

func TestBatchOffchainValidateMaxGasPriceOrder(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	orderValidator, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	require.NoError(t, err)

	for _, staticCallAssetData := range [][]byte{
		checkGasPriceDefaultStaticCallData,
		checkGasPriceStaticCallData,
	} {
		teardownSubTest := setupSubTest(t)

		// Create the signed order with the staticcall asset data as its MakerFeeAssetData
		signedOrder := scenario.NewSignedTestOrder(t, orderopts.MakerFeeAssetData(staticCallAssetData))
		signedOrders := []*zeroex.SignedOrder{
			signedOrder,
		}

		// Ensure that the order is accepted by offchain validation
		accepted, rejected := orderValidator.BatchOffchainValidation(signedOrders)
		assert.Len(t, accepted, 1)
		require.Len(t, rejected, 0)
		signedOrder.ResetHash()
		expectedOrderHash, err := signedOrder.ComputeOrderHash()
		require.NoError(t, err)
		actualOrderHash, err := accepted[0].ComputeOrderHash()
		require.NoError(t, err)
		assert.Equal(t, expectedOrderHash, actualOrderHash)

		teardownSubTest(t)
	}
}

func TestBatchValidateMaxGasPriceOrder(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	orderValidator, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	require.NoError(t, err)

	for _, staticCallAssetData := range [][]byte{
		checkGasPriceDefaultStaticCallData,
		checkGasPriceStaticCallData,
	} {

		teardownSubTest := setupSubTest(t)

		// Create the signed order with the staticcall asset data as its MakerFeeAssetData
		signedOrder := scenario.NewSignedTestOrder(t, orderopts.SetupMakerState(true), orderopts.MakerFeeAssetData(staticCallAssetData))
		signedOrders := []*zeroex.SignedOrder{
			signedOrder,
		}

		// Ensure that the order is accepted by offchain validation
		ctx := context.Background()
		latestBlock, err := ethRPCClient.HeaderByNumber(ctx, nil)
		require.NoError(t, err)
		validationResults := orderValidator.BatchValidate(ctx, signedOrders, areNewOrders, latestBlock)
		assert.Len(t, validationResults.Accepted, 1)
		require.Len(t, validationResults.Rejected, 0)
		expectedOrderHash, err := signedOrder.ComputeOrderHash()
		require.NoError(t, err)
		assert.Equal(t, expectedOrderHash, validationResults.Accepted[0].OrderHash)

		teardownSubTest(t)
	}
}

func TestBatchValidateSignatureInvalid(t *testing.T) {
	signedOrder := signedOrderWithCustomSignature(t, malformedSignature)
	signedOrders := []*zeroex.SignedOrder{
		signedOrder,
	}
	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)

	orderValidator, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	require.NoError(t, err)

	ctx := context.Background()
	latestBlock, err := ethRPCClient.HeaderByNumber(ctx, nil)
	require.NoError(t, err)
	validationResults := orderValidator.BatchValidate(ctx, signedOrders, areNewOrders, latestBlock)
	assert.Len(t, validationResults.Accepted, 0)
	require.Len(t, validationResults.Rejected, 1)
	assert.Equal(t, ROInvalidSignature, validationResults.Rejected[0].Status)
	assert.Equal(t, orderHash, validationResults.Rejected[0].OrderHash)
}

func TestBatchValidateUnregisteredCoordinator(t *testing.T) {
	// FeeRecipientAddress is an address for which there is no entry in the Coordinator registry
	signedOrder := scenario.NewSignedTestOrder(t, orderopts.SenderAddress(ganacheAddresses.Coordinator), orderopts.FeeRecipientAddress(constants.GanacheAccount4))
	signedOrder.FeeRecipientAddress = constants.GanacheAccount4
	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	signedOrders := []*zeroex.SignedOrder{
		signedOrder,
	}

	orderValidator, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	require.NoError(t, err)

	ctx := context.Background()
	latestBlock, err := ethRPCClient.HeaderByNumber(ctx, nil)
	require.NoError(t, err)
	validationResults := orderValidator.BatchValidate(ctx, signedOrders, areNewOrders, latestBlock)
	assert.Len(t, validationResults.Accepted, 0)
	require.Len(t, validationResults.Rejected, 1)
	assert.Equal(t, ROCoordinatorEndpointNotFound, validationResults.Rejected[0].Status)
	assert.Equal(t, orderHash, validationResults.Rejected[0].OrderHash)
}

func TestBatchValidateCoordinatorSoftCancels(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	signedOrder := scenario.NewSignedTestOrder(t,
		orderopts.SenderAddress(ganacheAddresses.Coordinator),
		orderopts.SetupMakerState(true),
		orderopts.FeeRecipientAddress(constants.GanacheAccount3),
	)
	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	signedOrders := []*zeroex.SignedOrder{
		signedOrder,
	}

	orderValidator, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	require.NoError(t, err)

	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		_, err := res.Write([]byte(fmt.Sprintf(`{"orderHashes": ["%s"]}`, orderHash.Hex())))
		require.NoError(t, err)
	}))
	defer testServer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	opts := &bind.TransactOpts{
		From:    signedOrder.FeeRecipientAddress,
		Context: ctx,
		Signer: func(s types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			testSigner := signer.NewTestSigner()
			signature, err := testSigner.(*signer.TestSigner).SignTx(s.Hash(tx).Bytes(), address)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(s, signature)
		},
	}

	coordinatorRegistryAddress := ganacheAddresses.CoordinatorRegistry
	coordinatorRegistry, err := wrappers.NewCoordinatorRegistry(coordinatorRegistryAddress, ethClient)
	require.NoError(t, err)
	_, err = coordinatorRegistry.SetCoordinatorEndpoint(opts, testServer.URL)
	require.NoError(t, err)

	ctx = context.Background()
	latestBlock, err := ethRPCClient.HeaderByNumber(ctx, nil)
	require.NoError(t, err)
	validationResults := orderValidator.BatchValidate(ctx, signedOrders, areNewOrders, latestBlock)
	assert.Len(t, validationResults.Accepted, 0)
	require.Len(t, validationResults.Rejected, 1)
	assert.Equal(t, ROCoordinatorSoftCancelled, validationResults.Rejected[0].Status)
	assert.Equal(t, orderHash, validationResults.Rejected[0].OrderHash)
}

const singleOrderPayloadSize = 2236

func TestComputeOptimalChunkSizesMaxContentLengthTooLow(t *testing.T) {
	signedOrder := scenario.NewSignedTestOrder(t)
	maxContentLength := singleOrderPayloadSize - 10
	orderValidator, err := New(ethRPCClient, constants.TestChainID, maxContentLength, ganacheAddresses)
	require.NoError(t, err)

	signedOrders := []*zeroex.SignedOrder{signedOrder}
	assert.Panics(t, func() {
		orderValidator.computeOptimalChunkSizes(signedOrders)
	})
}

func TestComputeOptimalChunkSizes(t *testing.T) {
	signedOrder := scenario.NewSignedTestOrder(t)
	maxContentLength := singleOrderPayloadSize * 3
	orderValidator, err := New(ethRPCClient, constants.TestChainID, maxContentLength, ganacheAddresses)
	require.NoError(t, err)

	signedOrders := []*zeroex.SignedOrder{signedOrder, signedOrder, signedOrder, signedOrder}
	chunkSizes := orderValidator.computeOptimalChunkSizes(signedOrders)
	expectedChunkSizes := []int{3, 1}
	assert.Equal(t, expectedChunkSizes, chunkSizes)
}

func TestComputeOptimalChunkSizesMultiAssetOrder(t *testing.T) {
	signedOrder := scenario.NewSignedTestOrder(t)
	signedMultiAssetOrder := scenario.NewSignedTestOrder(t, orderopts.MakerAssetData(multiAssetAssetData))

	maxContentLength := singleOrderPayloadSize * 3
	orderValidator, err := New(ethRPCClient, constants.TestChainID, maxContentLength, ganacheAddresses)
	require.NoError(t, err)

	signedOrders := []*zeroex.SignedOrder{signedMultiAssetOrder, signedOrder, signedOrder, signedOrder, signedOrder}
	chunkSizes := orderValidator.computeOptimalChunkSizes(signedOrders)
	expectedChunkSizes := []int{2, 3} // MultiAsset order is larger so can only fit two orders in first chunk
	assert.Equal(t, expectedChunkSizes, chunkSizes)
}

func abiEncode(signedOrder *zeroex.SignedOrder) ([]byte, error) {
	trimmedOrder := signedOrder.Trim()

	devUtilsABI, err := abi.JSON(strings.NewReader(wrappers.DevUtilsABI))
	if err != nil {
		return []byte{}, err
	}

	data, err := devUtilsABI.Pack(
		"getOrderRelevantStates",
		[]wrappers.LibOrderOrder{trimmedOrder},
		[][]byte{signedOrder.Signature},
	)
	if err != nil {
		return []byte{}, err
	}

	dataBytes := hexutil.Bytes(data)
	encodedData, err := json.Marshal(dataBytes)
	if err != nil {
		return []byte{}, err
	}

	return encodedData, nil
}

func TestComputeABIEncodedSignedOrderStringByteLength(t *testing.T) {
	testOrder := scenario.NewSignedTestOrder(t)

	testMultiAssetOrder := scenario.NewSignedTestOrder(t)
	testMultiAssetOrder.Order.MakerAssetData = common.Hex2Bytes("123412304102340120350120340123041023401234102341234234523452345234")
	testMultiAssetOrder.Order.MakerAssetData = common.Hex2Bytes("132519348523094582039457283452")
	testMultiAssetOrder.Order.MakerAssetData = multiAssetAssetData
	testMultiAssetOrder.Order.MakerAssetData = common.Hex2Bytes("324857203942034562893723452345246529837")

	testCases := []*zeroex.SignedOrder{testOrder, testMultiAssetOrder}

	for _, signedOrder := range testCases {
		label := fmt.Sprintf("test order: %v", signedOrder)

		encoded, err := abiEncode(signedOrder)
		require.NoError(t, err)
		expectedLength := len(encoded) - emptyGetOrderRelevantStatesCallDataStringLength

		length := computeABIEncodedSignedOrderStringLength(signedOrder)

		assert.Equal(t, expectedLength, length, label)
	}
}

func setupSubTest(t *testing.T) func(t *testing.T) {
	blockchainLifecycle.Start(t)
	return func(t *testing.T) {
		blockchainLifecycle.Revert(t)
	}
}

func signedOrderWithCustomSignature(t *testing.T, signature []byte) *zeroex.SignedOrder {
	signedOrder := scenario.NewSignedTestOrder(t)
	signedOrder.Signature = signature
	return signedOrder
}

func copyOrder(order zeroex.Order) zeroex.Order {
	return order
}
