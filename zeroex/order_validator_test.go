// +build !js

// We currently don't run these tests in WASM because of an issue in Go. See the header of
// eth_watcher_test.go for more details.
package zeroex

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/scenario"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const areNewOrders = false

var makerAddress = constants.GanacheAccount0
var takerAddress = constants.GanacheAccount1
var eighteenDecimalsInBaseUnits = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
var wethAmount = new(big.Int).Mul(big.NewInt(50), eighteenDecimalsInBaseUnits)
var zrxAmount = new(big.Int).Mul(big.NewInt(100), eighteenDecimalsInBaseUnits)

var unsupportedAssetData = common.Hex2Bytes("a2cb61b000000000000000000000000034d402f14d58e001d8efbe6585051bf9706aa064")
var malformedAssetData = []byte("9HJhsAAAAAAAAAAAAAAAAInSSmtMyxtvqiYl")
var malformedSignature = []byte("9HJhsAAAAAAAAAAAAAAAAInSSmtMyxtvqiYl")
var multiAssetAssetData = common.Hex2Bytes("94cfcdd7000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000046000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000024f47261b00000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c48000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000204a7cb5fb70000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001800000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000003e90000000000000000000000000000000000000000000000000000000000002711000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000c800000000000000000000000000000000000000000000000000000000000007d10000000000000000000000000000000000000000000000000000000000004e210000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c4800000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

var testSignedOrder = SignedOrder{
	Order: Order{
		MakerAddress:          makerAddress,
		TakerAddress:          constants.NullAddress,
		SenderAddress:         constants.NullAddress,
		FeeRecipientAddress:   constants.GanacheAccount3,
		MakerAssetData:        common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		TakerAssetData:        common.Hex2Bytes("f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
		Salt:                  big.NewInt(1548619145450),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(1000),
		TakerAssetAmount:      big.NewInt(2000),
		ExpirationTimeSeconds: big.NewInt(time.Now().Add(48 * time.Hour).Unix()),
		ExchangeAddress:       ethereum.NetworkIDToContractAddresses[constants.TestNetworkID].Exchange,
	},
}

type testCase struct {
	SignedOrder                 SignedOrder
	IsValid                     bool
	ExpectedRejectedOrderStatus RejectedOrderStatus
}

var blockchainLifecycle *ethereum.BlockchainLifecycle

func TestSetup(t *testing.T) {
	var err error
	blockchainLifecycle, err = ethereum.NewBlockchainLifecycle(constants.GanacheEndpoint)
	require.NoError(t, err)
}

func TestBatchValidateOffChainCases(t *testing.T) {
	var testCases = []testCase{
		testCase{
			SignedOrder:                 signedOrderWithCustomMakerAssetAmount(t, testSignedOrder, big.NewInt(0)),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidMakerAssetAmount,
		},
		testCase{
			SignedOrder: signedOrderWithCustomMakerAssetAmount(t, testSignedOrder, big.NewInt(1000000)),
			IsValid:     true,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomTakerAssetAmount(t, testSignedOrder, big.NewInt(0)),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidTakerAssetAmount,
		},
		testCase{
			SignedOrder: signedOrderWithCustomTakerAssetAmount(t, testSignedOrder, big.NewInt(1000000)),
			IsValid:     true,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomMakerAssetData(t, testSignedOrder, multiAssetAssetData),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROUnfunded,
		},
		testCase{
			SignedOrder: signedOrderWithCustomTakerAssetData(t, testSignedOrder, multiAssetAssetData),
			IsValid:     true,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomMakerAssetData(t, testSignedOrder, malformedAssetData),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidMakerAssetData,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomTakerAssetData(t, testSignedOrder, malformedAssetData),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidTakerAssetData,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomMakerAssetData(t, testSignedOrder, unsupportedAssetData),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidMakerAssetData,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomTakerAssetData(t, testSignedOrder, unsupportedAssetData),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidTakerAssetData,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomExpirationTimeSeconds(t, testSignedOrder, big.NewInt(time.Now().Add(-5*time.Minute).Unix())),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROExpired,
		},
		testCase{
			SignedOrder:                 signedOrderWithCustomSignature(t, testSignedOrder, malformedSignature),
			IsValid:                     false,
			ExpectedRejectedOrderStatus: ROInvalidSignature,
		},
	}

	for _, testCase := range testCases {
		teardownSubTest := setupSubTest(t)

		ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
		require.NoError(t, err)

		signedOrders := []*SignedOrder{
			&testCase.SignedOrder,
		}

		orderValidator, err := NewOrderValidator(ethClient, constants.TestNetworkID, constants.TestMaxContentLength, 0)
		require.NoError(t, err)

		validationResults := orderValidator.BatchValidate(signedOrders, areNewOrders)
		isValid := len(validationResults.Accepted) == 1
		assert.Equal(t, testCase.IsValid, isValid, testCase.ExpectedRejectedOrderStatus)
		if !isValid {
			assert.Equal(t, testCase.ExpectedRejectedOrderStatus, validationResults.Rejected[0].Status)
		}

		teardownSubTest(t)
	}
}

func TestBatchValidateSignatureInvalid(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	signedOrder := &testSignedOrder
	// Add a correctly formatted signature that does not correspond to this order
	signedOrder.Signature = common.Hex2Bytes("1c3582f06356a1314dbf1c0e534c4d8e92e59b056ee607a7ff5a825f5f2cc5e6151c5cc7fdd420f5608e4d5bef108e42ad90c7a4b408caef32e24374cf387b0d7603")

	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)

	signedOrders := []*SignedOrder{
		signedOrder,
	}

	ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)

	orderValidator, err := NewOrderValidator(ethClient, constants.TestNetworkID, constants.TestMaxContentLength, 0)
	require.NoError(t, err)

	validationResults := orderValidator.BatchValidate(signedOrders, areNewOrders)
	assert.Len(t, validationResults.Accepted, 0)
	require.Len(t, validationResults.Rejected, 1)
	assert.Equal(t, ROInvalidSignature, validationResults.Rejected[0].Status)
	assert.Equal(t, orderHash, validationResults.Rejected[0].OrderHash)
}

func TestBatchValidateUnregisteredCoordinatorSoftCancels(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	signedOrder := &testSignedOrder
	signedOrder.SenderAddress = ethereum.NetworkIDToContractAddresses[constants.TestNetworkID].Coordinator
	// Address for which there is no entry in the Coordinator registry
	signedOrder.FeeRecipientAddress = constants.GanacheAccount4

	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)

	signedOrders := []*SignedOrder{
		signedOrder,
	}

	ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)

	orderValidator, err := NewOrderValidator(ethClient, constants.TestNetworkID, constants.TestMaxContentLength, 0)
	require.NoError(t, err)

	validationResults := orderValidator.BatchValidate(signedOrders, areNewOrders)
	assert.Len(t, validationResults.Accepted, 0)
	require.Len(t, validationResults.Rejected, 1)
	assert.Equal(t, ROCoordinatorEndpointNotFound, validationResults.Rejected[0].Status)
	assert.Equal(t, orderHash, validationResults.Rejected[0].OrderHash)
}

func TestBatchValidateCoordinatorSoftCancels(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	signedOrder := &testSignedOrder
	signedOrder.SenderAddress = ethereum.NetworkIDToContractAddresses[constants.TestNetworkID].Coordinator
	orderHash, err := signedOrder.ComputeOrderHash()
	require.NoError(t, err)
	signedOrders := []*SignedOrder{
		signedOrder,
	}

	ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)
	orderValidator, err := NewOrderValidator(ethClient, constants.TestNetworkID, constants.TestMaxContentLength, 0)
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
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			testSigner := ethereum.NewTestSigner()
			signature, err := testSigner.(*ethereum.TestSigner).SignTx(signer.Hash(tx).Bytes(), address)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}
	_, err = orderValidator.coordinatorRegistry.SetCoordinatorEndpoint(opts, testServer.URL)
	require.NoError(t, err)

	validationResults := orderValidator.BatchValidate(signedOrders, areNewOrders)
	assert.Len(t, validationResults.Accepted, 0)
	require.Len(t, validationResults.Rejected, 1)
	assert.Equal(t, ROCoordinatorSoftCancelled, validationResults.Rejected[0].Status)
	assert.Equal(t, orderHash, validationResults.Rejected[0].OrderHash)
}

const singleOrderPayloadSize = 1980

func TestComputeOptimalChunkSizesMaxContentLengthTooLow(t *testing.T) {
	signedOrder, err := SignTestOrder(&testSignedOrder.Order)
	require.NoError(t, err)

	ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)

	maxContentLength := singleOrderPayloadSize - 10
	orderValidator, err := NewOrderValidator(ethClient, constants.TestNetworkID, maxContentLength, 0)
	require.NoError(t, err)

	signedOrders := []*SignedOrder{signedOrder}
	assert.Panics(t, func() {
		orderValidator.computeOptimalChunkSizes(signedOrders)
	})
}

func TestComputeOptimalChunkSizes(t *testing.T) {
	signedOrder, err := SignTestOrder(&testSignedOrder.Order)
	require.NoError(t, err)

	ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)

	maxContentLength := singleOrderPayloadSize * 3
	orderValidator, err := NewOrderValidator(ethClient, constants.TestNetworkID, maxContentLength, 0)
	require.NoError(t, err)

	signedOrders := []*SignedOrder{signedOrder, signedOrder, signedOrder, signedOrder}
	chunkSizes := orderValidator.computeOptimalChunkSizes(signedOrders)
	expectedChunkSizes := []int{3, 1}
	assert.Equal(t, expectedChunkSizes, chunkSizes)
}

var testMultiAssetSignedOrder = SignedOrder{
	Order: Order{
		MakerAddress:          constants.GanacheAccount0,
		TakerAddress:          constants.NullAddress,
		SenderAddress:         constants.NullAddress,
		FeeRecipientAddress:   common.HexToAddress("0x6ecbe1db9ef729cbe972c83fb886247691fb6beb"),
		MakerAssetData:        common.Hex2Bytes("94cfcdd7000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000046000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000024f47261b00000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c48000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000204a7cb5fb70000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c480000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001800000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000003e90000000000000000000000000000000000000000000000000000000000002711000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000c800000000000000000000000000000000000000000000000000000000000007d10000000000000000000000000000000000000000000000000000000000004e210000000000000000000000000000000000000000000000000000000000000044025717920000000000000000000000001dc4c1cefef38a777b15aa20260a54e584b16c4800000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		TakerAssetData:        common.Hex2Bytes("f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
		Salt:                  big.NewInt(1548619145450),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(1000),
		TakerAssetAmount:      big.NewInt(2000),
		ExpirationTimeSeconds: big.NewInt(time.Now().Add(48 * time.Hour).Unix()),
		ExchangeAddress:       ethereum.NetworkIDToContractAddresses[constants.TestNetworkID].Exchange,
	},
}

func TestComputeOptimalChunkSizesMultiAssetOrder(t *testing.T) {
	signedOrder, err := SignTestOrder(&testSignedOrder.Order)
	require.NoError(t, err)
	signedMultiAssetOrder, err := SignTestOrder(&testMultiAssetSignedOrder.Order)
	require.NoError(t, err)

	ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)

	maxContentLength := singleOrderPayloadSize * 3
	orderValidator, err := NewOrderValidator(ethClient, constants.TestNetworkID, maxContentLength, 0)
	require.NoError(t, err)

	signedOrders := []*SignedOrder{signedMultiAssetOrder, signedOrder, signedOrder, signedOrder, signedOrder}
	chunkSizes := orderValidator.computeOptimalChunkSizes(signedOrders)
	expectedChunkSizes := []int{2, 3} // MultiAsset order is larger so can only fit two orders in first chunk
	assert.Equal(t, expectedChunkSizes, chunkSizes)
}

func setupSubTest(t *testing.T) func(t *testing.T) {
	blockchainLifecycle.Start(t)
	scenario.SetupBalancesAndAllowances(t, makerAddress, takerAddress, wethAmount, zrxAmount)
	return func(t *testing.T) {
		blockchainLifecycle.Revert(t)
	}
}

func signedOrderWithCustomMakerAssetAmount(t *testing.T, signedOrder SignedOrder, makerAssetAmount *big.Int) SignedOrder {
	signedOrder.MakerAssetAmount = makerAssetAmount
	signedOrderWithSignature, err := SignTestOrder(&signedOrder.Order)
	require.NoError(t, err)
	return *signedOrderWithSignature
}

func signedOrderWithCustomTakerAssetAmount(t *testing.T, signedOrder SignedOrder, takerAssetAmount *big.Int) SignedOrder {
	signedOrder.TakerAssetAmount = takerAssetAmount
	signedOrderWithSignature, err := SignTestOrder(&signedOrder.Order)
	require.NoError(t, err)
	return *signedOrderWithSignature
}

func signedOrderWithCustomMakerAssetData(t *testing.T, signedOrder SignedOrder, makerAssetData []byte) SignedOrder {
	signedOrder.MakerAssetData = makerAssetData
	signedOrderWithSignature, err := SignTestOrder(&signedOrder.Order)
	require.NoError(t, err)
	return *signedOrderWithSignature
}

func signedOrderWithCustomTakerAssetData(t *testing.T, signedOrder SignedOrder, takerAssetData []byte) SignedOrder {
	signedOrder.TakerAssetData = takerAssetData
	signedOrderWithSignature, err := SignTestOrder(&signedOrder.Order)
	require.NoError(t, err)
	return *signedOrderWithSignature
}

func signedOrderWithCustomExpirationTimeSeconds(t *testing.T, signedOrder SignedOrder, expirationTimeSeconds *big.Int) SignedOrder {
	signedOrder.ExpirationTimeSeconds = expirationTimeSeconds
	signedOrderWithSignature, err := SignTestOrder(&signedOrder.Order)
	require.NoError(t, err)
	return *signedOrderWithSignature
}

func signedOrderWithCustomSignature(t *testing.T, signedOrder SignedOrder, signature []byte) SignedOrder {
	signedOrder.Signature = signature
	return signedOrder
}

func copyOrder(order Order) Order {
	return order
}
