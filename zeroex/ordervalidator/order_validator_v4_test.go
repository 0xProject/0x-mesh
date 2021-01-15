// +build !js

// We currently don't run these tests in WASM because of an issue in Go. See the header of
// eth_watcher_test.go for more details.
package ordervalidator

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/ethrpcclient"
	"github.com/0xProject/0x-mesh/ethereum/ratelimit"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/scenario/orderopts"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// emptyBatchGetLimitOrderRelevantStatesCallDataStringLength is all the boilerplate ABI encoding
// required when calling `batchGetLimitOrderRelevantStates` that does not include the encoded
// SignedOrderV4. By subtracting this amount from the calldata length returned from encoding a
// call to `batchGetLimitOrderRelevantStates` involving a single SignedOrder, we get the number of
// bytes taken up by the SignedOrderV4 alone in hex encoding including prefix and string quotation marks. i.e.: len(`"0x[...]"`)
const emptyBatchGetLimitOrderRelevantStatesCallDataStringLength = 4 + 8 + 64*4

type testCaseV4 struct {
	SignedOrder                 *zeroex.SignedOrderV4
	IsValid                     bool
	ExpectedRejectedOrderStatus RejectedOrderStatus
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

func TestBatchValidateOffChainCasesV4(t *testing.T) {
	invalidSignedOrder := scenario.NewSignedTestOrderV4(t)
	invalidSignedOrder.R = zeroex.HexToBytes32("a2cb61b6585051bf9706585051bf97034d402f14d58e001d8efbe6585051bf97") // Random

	var testCases = []testCaseV4{
		{
			SignedOrder: scenario.NewSignedTestOrderV4(t),
			IsValid:     true,
		},
		{
			SignedOrder:                 scenario.NewSignedTestOrderV4(t, orderopts.MakerAssetAmount(big.NewInt(0))),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidMakerAssetAmount,
		},
		{
			SignedOrder:                 scenario.NewSignedTestOrderV4(t, orderopts.TakerAssetAmount(big.NewInt(0))),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidTakerAssetAmount,
		},
		{
			SignedOrder:                 invalidSignedOrder,
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidSignature,
		},
	}

	for _, testCase := range testCases {
		signedOrders := []*zeroex.SignedOrderV4{
			testCase.SignedOrder,
		}
		orderValidator, err := New(ethClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
		require.NoError(t, err)

		offchainValidOrders, rejectedOrderInfos := orderValidator.BatchOffchainValidationV4(signedOrders)
		isValid := len(offchainValidOrders) == 1
		assert.Equal(t, testCase.IsValid, isValid, testCase.ExpectedRejectedOrderStatus)
		if !isValid {
			assert.Equal(t, testCase.ExpectedRejectedOrderStatus, rejectedOrderInfos[0].Status)
		}
	}
}

func TestBatchValidateAValidOrderV4(t *testing.T) {
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

func TestBatchOffchainValidateZeroFeeAmountV4(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	makerFeeAssetData := common.Hex2Bytes("deadbeef")
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
	assert.Len(t, accepted, 1)
	require.Len(t, rejected, 0)
	signedTestOrder.ResetHash()
	expectedOrderHash, err := signedTestOrder.ComputeOrderHash()
	require.NoError(t, err)
	actualOrderHash, err := accepted[0].ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, expectedOrderHash, actualOrderHash)
}

func TestBatchOffchainValidateUnsupportedStaticCallV4(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	// NOTE(jalextowle): This asset data encodes a staticcall to a function called `unsupportedStaticCall`
	makerFeeAssetData := common.Hex2Bytes("c339d10a000000000000000000000000692a70d2e424a56d2c6c27aa97d1a86395877b3a0000000000000000000000000000000000000000000000000000000000000060c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47000000000000000000000000000000000000000000000000000000000000000048b24020700000000000000000000000000000000000000000000000000000000")
	signedTestOrder := scenario.NewSignedTestOrder(
		t,
		orderopts.MakerFeeAssetData(makerFeeAssetData),
		orderopts.MakerFee(big.NewInt(1)),
	)
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

func TestBatchOffchainValidateMaxGasPriceOrderV4(t *testing.T) {
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

func TestBatchValidateMaxGasPriceOrderV4(t *testing.T) {
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

func TestBatchValidateSignatureInvalidV4(t *testing.T) {
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

func TestComputeOptimalChunkSizesMaxContentLengthTooLowV4(t *testing.T) {
	signedOrder := scenario.NewSignedTestOrder(t)
	maxContentLength := singleOrderPayloadSize - 10
	orderValidator, err := New(ethRPCClient, constants.TestChainID, maxContentLength, ganacheAddresses)
	require.NoError(t, err)

	signedOrders := []*zeroex.SignedOrder{signedOrder}
	assert.Panics(t, func() {
		orderValidator.computeOptimalChunkSizes(signedOrders)
	})
}

func TestComputeOptimalChunkSizesV4(t *testing.T) {
	signedOrder := scenario.NewSignedTestOrder(t)
	maxContentLength := singleOrderPayloadSize * 3
	orderValidator, err := New(ethRPCClient, constants.TestChainID, maxContentLength, ganacheAddresses)
	require.NoError(t, err)

	signedOrders := []*zeroex.SignedOrder{signedOrder, signedOrder, signedOrder, signedOrder}
	chunkSizes := orderValidator.computeOptimalChunkSizes(signedOrders)
	expectedChunkSizes := []int{3, 1}
	assert.Equal(t, expectedChunkSizes, chunkSizes)
}

func abiEncodeV4(signedOrder *zeroex.SignedOrderV4) ([]byte, error) {
	abiOrder := signedOrder.EthereumAbiLimitOrder()
	abiSignature := signedOrder.EthereumAbiSignature()

	exchangeV4ABI, err := abi.JSON(strings.NewReader(wrappers.ExchangeV4ABI))
	if err != nil {
		return []byte{}, err
	}

	data, err := exchangeV4ABI.Pack(
		"batchGetLimitOrderRelevantStates",
		[]wrappers.LibNativeOrderLimitOrder{abiOrder},
		[]wrappers.LibSignatureSignature{abiSignature},
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

func TestComputeABIEncodedSignedOrderStringByteLengthV4(t *testing.T) {
	testOrder := scenario.NewSignedTestOrderV4(t)
	testCases := []*zeroex.SignedOrderV4{testOrder}

	for _, signedOrder := range testCases {
		label := fmt.Sprintf("test order: %v", signedOrder)

		encoded, err := abiEncodeV4(signedOrder)
		require.NoError(t, err)
		t.Logf("abiEncoded = %s\n", encoded)
		expectedLength := len(encoded) - emptyBatchGetLimitOrderRelevantStatesCallDataStringLength

		assert.Equal(t, expectedLength, signedOrderV4AbiHexLength, label)
	}
}
