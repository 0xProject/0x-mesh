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
	"strings"
	"testing"
	"time"

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

const areNewOrderV3s = false

// emptyGetOrderV3RelevantStatesCallDataByteLength is all the boilerplate ABI encoding required when calling
// `getOrde3RelevantStates` that does not include the encoded SignedOrderV3. By subtracting this amount from the
// calldata length returned from encoding a call to `getOrderRelevantStates` involving a single SignedOrderV3, we
// get the number of bytes taken up by the SignedOrderV3 alone.
// i.e.: len(`"0x7f46448d0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"`)
const emptyGetOrderV3RelevantStatesCallDataStringLength = 268

const (
	defaultEthRPCTimeout = 5 * time.Second
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
	SignedOrderV3               *zeroex.SignedOrderV3
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
		{
			SignedOrderV3:               scenario.NewSignedTestOrderV3(t, orderopts.MakerAssetAmount(big.NewInt(0))),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidMakerAssetAmount,
		},
		{
			SignedOrderV3:               scenario.NewSignedTestOrderV3(t, orderopts.TakerAssetAmount(big.NewInt(0))),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidTakerAssetAmount,
		},
		{
			SignedOrderV3: scenario.NewSignedTestOrderV3(t, orderopts.MakerAssetData(multiAssetAssetData)),
			IsValid:       true,
		},
		{
			SignedOrderV3:               scenario.NewSignedTestOrderV3(t, orderopts.MakerAssetData(malformedAssetData)),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidMakerAssetData,
		},
		{
			SignedOrderV3:               scenario.NewSignedTestOrderV3(t, orderopts.TakerAssetData(malformedAssetData)),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidTakerAssetData,
		},
		{
			SignedOrderV3:               scenario.NewSignedTestOrderV3(t, orderopts.MakerAssetData(unsupportedAssetData)),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidMakerAssetData,
		},
		{
			SignedOrderV3:               scenario.NewSignedTestOrderV3(t, orderopts.TakerAssetData(unsupportedAssetData)),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidTakerAssetData,
		},
		{
			SignedOrderV3:               signedOrderV3WithCustomSignature(t, malformedSignature),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidSignature,
		},
	}

	for _, testCase := range testCases {
		signedOrderV3s := []*zeroex.SignedOrderV3{
			testCase.SignedOrderV3,
		}
		orderValidator, err := New(ethClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
		require.NoError(t, err)

		offchainValidOrderV3s, rejectedOrderV3Infos := orderValidator.BatchOffchainValidation(signedOrderV3s)
		isValid := len(offchainValidOrderV3s) == 1
		assert.Equal(t, testCase.IsValid, isValid, testCase.ExpectedRejectedOrderStatus)
		if !isValid {
			assert.Equal(t, testCase.ExpectedRejectedOrderStatus, rejectedOrderV3Infos[0].Status)
		}
	}
}

func TestBatchValidateAValidOrderV3(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	signedOrderV3 := scenario.NewSignedTestOrderV3(t, orderopts.SetupMakerState(true))
	signedOrderV3s := []*zeroex.SignedOrderV3{
		signedOrderV3,
	}

	orderValidator, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	require.NoError(t, err)

	ctx := context.Background()
	latestBlock, err := ethRPCClient.HeaderByNumber(ctx, nil)
	require.NoError(t, err)
	validationResults := orderValidator.BatchValidate(ctx, signedOrderV3s, areNewOrderV3s, latestBlock)
	assert.Len(t, validationResults.Accepted, 1)
	require.Len(t, validationResults.Rejected, 0)
	orderHash, err := signedOrderV3.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, orderHash, validationResults.Accepted[0].OrderHash)
}

func TestBatchOffchainValidateZeroFeeAmount(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	makerFeeAssetData := common.Hex2Bytes("deadbeef")
	signedTestOrderV3 := scenario.NewSignedTestOrderV3(t, orderopts.MakerFeeAssetData(makerFeeAssetData))
	signedOrderV3s := []*zeroex.SignedOrderV3{
		signedTestOrderV3,
	}

	rateLimiter := ratelimit.NewUnlimited()
	rpcClient, err := rpc.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)
	ethRPCClient, err := ethrpcclient.New(rpcClient, defaultEthRPCTimeout, rateLimiter)
	require.NoError(t, err)

	orderValidator, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	require.NoError(t, err)

	accepted, rejected := orderValidator.BatchOffchainValidation(signedOrderV3s)
	assert.Len(t, accepted, 1)
	require.Len(t, rejected, 0)
	signedTestOrderV3.ResetHash()
	expectedOrderV3Hash, err := signedTestOrderV3.ComputeOrderHash()
	require.NoError(t, err)
	actualOrderV3Hash, err := accepted[0].ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, expectedOrderV3Hash, actualOrderV3Hash)
}

func TestBatchOffchainValidateUnsupportedStaticCall(t *testing.T) {
	if !serialTestsEnabled {
		t.Skip("Serial tests (tests which cannot run in parallel) are disabled. You can enable them with the --serial flag")
	}

	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)
	// NOTE(jalextowle): This asset data encodes a staticcall to a function called `unsupportedStaticCall`
	makerFeeAssetData := common.Hex2Bytes("c339d10a000000000000000000000000692a70d2e424a56d2c6c27aa97d1a86395877b3a0000000000000000000000000000000000000000000000000000000000000060c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47000000000000000000000000000000000000000000000000000000000000000048b24020700000000000000000000000000000000000000000000000000000000")
	signedTestOrderV3 := scenario.NewSignedTestOrderV3(
		t,
		orderopts.MakerFeeAssetData(makerFeeAssetData),
		orderopts.MakerFee(big.NewInt(1)),
	)
	signedOrderV3s := []*zeroex.SignedOrderV3{
		signedTestOrderV3,
	}

	rateLimiter := ratelimit.NewUnlimited()
	rpcClient, err := rpc.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)
	ethRPCClient, err := ethrpcclient.New(rpcClient, defaultEthRPCTimeout, rateLimiter)
	require.NoError(t, err)

	orderValidator, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	require.NoError(t, err)

	accepted, rejected := orderValidator.BatchOffchainValidation(signedOrderV3s)
	assert.Len(t, accepted, 0)
	require.Len(t, rejected, 1)
	signedTestOrderV3.ResetHash()
	expectedOrderV3Hash, err := signedTestOrderV3.ComputeOrderHash()
	require.NoError(t, err)
	assert.Equal(t, expectedOrderV3Hash, rejected[0].OrderHash)
	require.Equal(t, ROInvalidMakerFeeAssetData, rejected[0].Status)
}

var checkGasPriceDefaultStaticCallData = common.Hex2Bytes("c339d10a0000000000000000000000002c530e4ecc573f11bd72cf5fdf580d134d25f15f0000000000000000000000000000000000000000000000000000000000000060c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a4700000000000000000000000000000000000000000000000000000000000000004d728f5b700000000000000000000000000000000000000000000000000000000")

var checkGasPriceStaticCallData = common.Hex2Bytes("c339d10a0000000000000000000000002c530e4ecc573f11bd72cf5fdf580d134d25f15f0000000000000000000000000000000000000000000000000000000000000060c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a4700000000000000000000000000000000000000000000000000000000000000024da5b166a000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000")

func TestBatchOffchainValidateMaxGasPriceOrderV3(t *testing.T) {
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
		signedOrderV3 := scenario.NewSignedTestOrderV3(t, orderopts.MakerFeeAssetData(staticCallAssetData))
		signedOrderV3s := []*zeroex.SignedOrderV3{
			signedOrderV3,
		}

		// Ensure that the order is accepted by offchain validation
		accepted, rejected := orderValidator.BatchOffchainValidation(signedOrderV3s)
		assert.Len(t, accepted, 1)
		require.Len(t, rejected, 0)
		signedOrderV3.ResetHash()
		expectedOrderV3Hash, err := signedOrderV3.ComputeOrderHash()
		require.NoError(t, err)
		actualOrderV3Hash, err := accepted[0].ComputeOrderHash()
		require.NoError(t, err)
		assert.Equal(t, expectedOrderV3Hash, actualOrderV3Hash)

		teardownSubTest(t)
	}
}

func TestBatchValidateMaxGasPriceOrderV3(t *testing.T) {
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
		signedOrderV3 := scenario.NewSignedTestOrderV3(t, orderopts.SetupMakerState(true), orderopts.MakerFeeAssetData(staticCallAssetData))
		signedOrderV3s := []*zeroex.SignedOrderV3{
			signedOrderV3,
		}

		// Ensure that the order is accepted by offchain validation
		ctx := context.Background()
		latestBlock, err := ethRPCClient.HeaderByNumber(ctx, nil)
		require.NoError(t, err)
		validationResults := orderValidator.BatchValidate(ctx, signedOrderV3s, areNewOrderV3s, latestBlock)
		assert.Len(t, validationResults.Accepted, 1)
		require.Len(t, validationResults.Rejected, 0)
		expectedOrderV3Hash, err := signedOrderV3.ComputeOrderHash()
		require.NoError(t, err)
		assert.Equal(t, expectedOrderV3Hash, validationResults.Accepted[0].OrderHash)

		teardownSubTest(t)
	}
}

func TestBatchValidateSignatureInvalid(t *testing.T) {
	signedOrderV3 := signedOrderV3WithCustomSignature(t, malformedSignature)
	signedOrderV3s := []*zeroex.SignedOrderV3{
		signedOrderV3,
	}
	orderHash, err := signedOrderV3.ComputeOrderHash()
	require.NoError(t, err)

	orderValidator, err := New(ethRPCClient, constants.TestChainID, constants.TestMaxContentLength, ganacheAddresses)
	require.NoError(t, err)

	ctx := context.Background()
	latestBlock, err := ethRPCClient.HeaderByNumber(ctx, nil)
	require.NoError(t, err)
	validationResults := orderValidator.BatchValidate(ctx, signedOrderV3s, areNewOrderV3s, latestBlock)
	assert.Len(t, validationResults.Accepted, 0)
	require.Len(t, validationResults.Rejected, 1)
	assert.Equal(t, ROInvalidSignature, validationResults.Rejected[0].Status)
	assert.Equal(t, orderHash, validationResults.Rejected[0].OrderHash)
}

const singleOrderV3PayloadSize = 2236

func TestComputeOptimalChunkSizesMaxContentLengthTooLow(t *testing.T) {
	signedOrderV3 := scenario.NewSignedTestOrderV3(t)
	maxContentLength := singleOrderV3PayloadSize - 10
	orderValidator, err := New(ethRPCClient, constants.TestChainID, maxContentLength, ganacheAddresses)
	require.NoError(t, err)

	signedOrderV3s := []*zeroex.SignedOrderV3{signedOrderV3}
	assert.Panics(t, func() {
		orderValidator.computeOptimalChunkSizes(signedOrderV3s)
	})
}

func TestComputeOptimalChunkSizes(t *testing.T) {
	signedOrderV3 := scenario.NewSignedTestOrderV3(t)
	maxContentLength := singleOrderV3PayloadSize * 3
	orderValidator, err := New(ethRPCClient, constants.TestChainID, maxContentLength, ganacheAddresses)
	require.NoError(t, err)

	signedOrderV3s := []*zeroex.SignedOrderV3{signedOrderV3, signedOrderV3, signedOrderV3, signedOrderV3}
	chunkSizes := orderValidator.computeOptimalChunkSizes(signedOrderV3s)
	expectedChunkSizes := []int{3, 1}
	assert.Equal(t, expectedChunkSizes, chunkSizes)
}

func TestComputeOptimalChunkSizesMultiAssetOrderV3(t *testing.T) {
	signedOrderV3 := scenario.NewSignedTestOrderV3(t)
	signedMultiAssetOrderV3 := scenario.NewSignedTestOrderV3(t, orderopts.MakerAssetData(multiAssetAssetData))

	maxContentLength := singleOrderV3PayloadSize * 3
	orderValidator, err := New(ethRPCClient, constants.TestChainID, maxContentLength, ganacheAddresses)
	require.NoError(t, err)

	signedOrderV3s := []*zeroex.SignedOrderV3{signedMultiAssetOrderV3, signedOrderV3, signedOrderV3, signedOrderV3, signedOrderV3}
	chunkSizes := orderValidator.computeOptimalChunkSizes(signedOrderV3s)
	expectedChunkSizes := []int{2, 3} // MultiAsset order is larger so can only fit two orders in first chunk
	assert.Equal(t, expectedChunkSizes, chunkSizes)
}

func abiEncode(signedOrderV3 *zeroex.SignedOrderV3) ([]byte, error) {
	trimmedOrderV3 := signedOrderV3.Trim()

	devUtilsABI, err := abi.JSON(strings.NewReader(wrappers.DevUtilsABI))
	if err != nil {
		return []byte{}, err
	}

	data, err := devUtilsABI.Pack(
		"getOrderRelevantStates",
		[]wrappers.LibOrderOrder{trimmedOrderV3},
		[][]byte{signedOrderV3.Signature},
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

func TestComputeABIEncodedSignedOrderV3StringByteLength(t *testing.T) {
	testOrderV3 := scenario.NewSignedTestOrderV3(t)

	testMultiAssetOrderV3 := scenario.NewSignedTestOrderV3(t)
	testMultiAssetOrderV3.OrderV3.MakerAssetData = common.Hex2Bytes("123412304102340120350120340123041023401234102341234234523452345234")
	testMultiAssetOrderV3.OrderV3.MakerAssetData = common.Hex2Bytes("132519348523094582039457283452")
	testMultiAssetOrderV3.OrderV3.MakerAssetData = multiAssetAssetData
	testMultiAssetOrderV3.OrderV3.MakerAssetData = common.Hex2Bytes("324857203942034562893723452345246529837")

	testCases := []*zeroex.SignedOrderV3{testOrderV3, testMultiAssetOrderV3}

	for _, signedOrderV3 := range testCases {
		label := fmt.Sprintf("test order: %v", signedOrderV3)

		encoded, err := abiEncode(signedOrderV3)
		require.NoError(t, err)
		expectedLength := len(encoded) - emptyGetOrderV3RelevantStatesCallDataStringLength

		length := computeABIEncodedSignedOrderV3StringLength(signedOrderV3)

		assert.Equal(t, expectedLength, length, label)
	}
}

func setupSubTest(t *testing.T) func(t *testing.T) {
	blockchainLifecycle.Start(t)
	return func(t *testing.T) {
		blockchainLifecycle.Revert(t)
	}
}

func signedOrderV3WithCustomSignature(t *testing.T, signature []byte) *zeroex.SignedOrderV3 {
	signedOrderV3 := scenario.NewSignedTestOrderV3(t)
	signedOrderV3.Signature = signature
	return signedOrderV3
}
