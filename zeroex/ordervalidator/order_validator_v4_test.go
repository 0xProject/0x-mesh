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
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	invalidSignedOrder.Signature.SignatureType = zeroex.InvalidSignatureV4

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

	signedOrder := scenario.NewSignedTestOrderV4(t, orderopts.SetupMakerState(true))
	signedOrders := []*zeroex.SignedOrderV4{
		signedOrder,
	}

	orderValidator, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	require.NoError(t, err)

	ctx := context.Background()
	latestBlock, err := ethRPCClient.HeaderByNumber(ctx, nil)
	require.NoError(t, err)
	validationResults := orderValidator.BatchValidateV4(ctx, signedOrders, areNewOrders, latestBlock)
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
	signedTestOrder := scenario.NewSignedTestOrderV4(t, orderopts.MakerFeeAssetData(makerFeeAssetData))
	signedOrders := []*zeroex.SignedOrderV4{
		signedTestOrder,
	}

	rateLimiter := ratelimit.NewUnlimited()
	rpcClient, err := rpc.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)
	ethRPCClient, err := ethrpcclient.New(rpcClient, defaultEthRPCTimeout, rateLimiter)
	require.NoError(t, err)

	orderValidator, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	require.NoError(t, err)

	accepted, rejected := orderValidator.BatchOffchainValidationV4(signedOrders)
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
	signedTestOrder := scenario.NewSignedTestOrderV4(
		t,
		orderopts.MakerFee(big.NewInt(1)),
	)
	signedOrders := []*zeroex.SignedOrderV4{
		signedTestOrder,
	}

	rateLimiter := ratelimit.NewUnlimited()
	rpcClient, err := rpc.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)
	ethRPCClient, err := ethrpcclient.New(rpcClient, defaultEthRPCTimeout, rateLimiter)
	require.NoError(t, err)

	orderValidator, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	require.NoError(t, err)

	accepted, rejected := orderValidator.BatchOffchainValidationV4(signedOrders)
	assert.Len(t, accepted, 1)
	require.Len(t, rejected, 0)
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
		signedOrder := scenario.NewSignedTestOrderV4(t, orderopts.MakerFeeAssetData(staticCallAssetData))
		signedOrders := []*zeroex.SignedOrderV4{
			signedOrder,
		}

		// Ensure that the order is accepted by offchain validation
		accepted, rejected := orderValidator.BatchOffchainValidationV4(signedOrders)
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
		signedOrder := scenario.NewSignedTestOrderV4(t, orderopts.SetupMakerState(true), orderopts.MakerFeeAssetData(staticCallAssetData))
		signedOrders := []*zeroex.SignedOrderV4{
			signedOrder,
		}

		// Ensure that the order is accepted by offchain validation
		ctx := context.Background()
		latestBlock, err := ethRPCClient.HeaderByNumber(ctx, nil)
		require.NoError(t, err)
		validationResults := orderValidator.BatchValidateV4(ctx, signedOrders, areNewOrders, latestBlock)
		assert.Len(t, validationResults.Accepted, 1)
		require.Len(t, validationResults.Rejected, 0)
		expectedOrderHash, err := signedOrder.ComputeOrderHash()
		require.NoError(t, err)
		assert.Equal(t, expectedOrderHash, validationResults.Accepted[0].OrderHash)

		teardownSubTest(t)
	}
}

func makeEthClient() (*OrderValidator, *bind.CallOpts, error) {
	eth, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	if err != nil {
		return nil, nil, err
	}
	ctx := context.Background()
	latestBlock, err := ethRPCClient.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	opts := &bind.CallOpts{
		// HACK(albrow): From field should not be required for eth_call but
		// including it here is a workaround for a bug in Ganache. Removing
		// this line causes Ganache to crash.
		From:    constants.GanacheDummyERC721TokenAddress,
		Pending: false,
		Context: ctx,
	}
	opts.BlockNumber = latestBlock.Number
	return eth, opts, nil
}

func TestOrderHashV4(t *testing.T) {
	signedOrder := scenario.NewSignedTestOrderV4(t)
	order := signedOrder.OrderV4

	localHash, err := order.ComputeOrderHash()
	require.NoError(t, err)

	eth, opts, err := makeEthClient()
	require.NoError(t, err)

	ethHash, err := eth.exchangeV4.GetLimitOrderHash(opts, order.EthereumAbiLimitOrder())
	require.NoError(t, err)
	assert.Equal(t, common.BytesToHash(ethHash[:]), localHash)
}

func TestOrderStateV4(t *testing.T) {
	order := &zeroex.OrderV4{
		ChainID:         big.NewInt(1337),
		ExchangeAddress: ganacheAddresses.ExchangeProxy,

		// TODO: Invalid token addresses currently make the call fail. This should be fixed soon, but for now make sure we use valid tokens.
		MakerToken:          ganacheAddresses.WETH9,
		TakerToken:          ganacheAddresses.ZRXToken,
		MakerAmount:         big.NewInt(1234),
		TakerAmount:         big.NewInt(5678),
		TakerTokenFeeAmount: big.NewInt(9101112),
		Maker:               constants.NullAddress, // Will be set by signer
		Taker:               common.HexToAddress("0x615312fb74c31303eab07dea520019bb23f4c6c2"),
		Sender:              common.HexToAddress("0x70f2d6c7acd257a6700d745b76c602ceefeb8e20"),
		FeeRecipient:        common.HexToAddress("0xcc3c7ea403427154ec908203ba6c418bd699f7ce"),
		Pool:                zeroex.HexToBytes32("0x0bbff69b85a87da39511aefc3211cb9aff00e1a1779dc35b8f3635d8b5ea2680"),
		Expiry:              big.NewInt(9223372036854775807),
		Salt:                big.NewInt(2001),
	}
	signed := order.TestSign(t)
	orderHash, err := signed.OrderV4.ComputeOrderHash()
	require.NoError(t, err)

	eth, opts, err := makeEthClient()
	require.NoError(t, err)

	ethState, err := eth.exchangeV4.GetLimitOrderRelevantState(opts, signed.OrderV4.EthereumAbiLimitOrder(), signed.EthereumAbiSignature())
	require.NoError(t, err)
	assert.Equal(t, orderHash, common.BytesToHash(ethState.OrderInfo.OrderHash[:]))
	assert.Equal(t, zeroex.OSInvalidMakerAssetAmount, zeroex.OrderStatus(ethState.OrderInfo.Status))
	assert.Equal(t, "0", ethState.OrderInfo.TakerTokenFilledAmount.String())
	assert.Equal(t, "0", ethState.ActualFillableTakerTokenAmount.String())
	assert.Equal(t, true, ethState.IsSignatureValid)
}

func TestBatchValidateV4(t *testing.T) {
	signedOrder := scenario.NewSignedTestOrderV4(t, orderopts.SetupMakerState(true))
	orderHash, err := signedOrder.OrderV4.ComputeOrderHash()
	require.NoError(t, err)

	orderValidator, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	require.NoError(t, err)

	ctx := context.Background()
	latestBlock, err := ethRPCClient.HeaderByNumber(ctx, nil)
	require.NoError(t, err)
	signedOrders := []*zeroex.SignedOrderV4{signedOrder, signedOrder, signedOrder, signedOrder}
	validationResults := orderValidator.BatchValidateV4(ctx, signedOrders, areNewOrders, latestBlock)
	assert.Len(t, validationResults.Accepted, 4)
	assert.Len(t, validationResults.Rejected, 0)
	assert.Equal(t, orderHash, validationResults.Accepted[0].OrderHash)
	assert.Equal(t, orderHash, validationResults.Accepted[1].OrderHash)
	assert.Equal(t, orderHash, validationResults.Accepted[2].OrderHash)
	assert.Equal(t, orderHash, validationResults.Accepted[3].OrderHash)
}

func TestBatchValidateSignatureInvalidV4(t *testing.T) {
	t.Skip("FIME: Invalid signatures cause batchGetLimitOrderRelevantStates to revert.")

	signedOrder := scenario.NewSignedTestOrderV4(t)
	signedOrder.Signature.R = zeroex.HexToBytes32("0000000000000000000000000000000000000000000000000000000000000000")

	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	t.Logf("hash = %+v\n", orderHash.Hex())

	orderValidator, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	require.NoError(t, err)

	ctx := context.Background()
	latestBlock, err := ethRPCClient.HeaderByNumber(ctx, nil)
	require.NoError(t, err)
	signedOrders := []*zeroex.SignedOrderV4{signedOrder}
	validationResults := orderValidator.BatchValidateV4(ctx, signedOrders, areNewOrders, latestBlock)
	assert.Len(t, validationResults.Accepted, 0)
	require.Len(t, validationResults.Rejected, 1)
	assert.Equal(t, ROInvalidSignature, validationResults.Rejected[0].Status)
	assert.Equal(t, orderHash, validationResults.Rejected[0].OrderHash)
}

func TestComputeOptimalChunkSizesMaxContentLengthTooLowV4(t *testing.T) {
	signedOrder := scenario.NewSignedTestOrderV4(t)
	maxContentLength := signedOrderV4AbiHexLength - 10
	orderValidator, err := New(ethRPCClient, constants.TestChainID, maxContentLength, ganacheAddresses)
	require.NoError(t, err)

	signedOrders := []*zeroex.SignedOrderV4{signedOrder}
	assert.Panics(t, func() {
		orderValidator.computeOptimalChunkSizesV4(signedOrders)
	})
}

func TestComputeOptimalChunkSizesV4(t *testing.T) {
	signedOrder := scenario.NewSignedTestOrderV4(t)
	maxContentLength := jsonRPCPayloadByteLength + signedOrderV4AbiHexLength*3
	orderValidator, err := New(ethRPCClient, constants.TestChainID, maxContentLength, ganacheAddresses)
	require.NoError(t, err)

	signedOrders := []*zeroex.SignedOrderV4{signedOrder, signedOrder, signedOrder, signedOrder}
	chunkSizes := orderValidator.computeOptimalChunkSizesV4(signedOrders)
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
		// t.Logf("abiEncoded = %s\n", encoded)
		expectedLength := len(encoded) - emptyBatchGetLimitOrderRelevantStatesCallDataStringLength

		assert.Equal(t, expectedLength, signedOrderV4AbiHexLength, label)
	}
}
