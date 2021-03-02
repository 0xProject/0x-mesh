// +build !js

package integrationtests

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/scenario"
	"github.com/0xProject/0x-mesh/scenario/orderopts"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddOrdersSuccessV4(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	wg := &sync.WaitGroup{}
	client, _ := buildAndStartGraphQLServer(t, ctx, wg)

	// Create a new valid order.
	signedTestOrder := scenario.NewSignedTestOrderV4(t, orderopts.SetupMakerState(true))
	time.Sleep(blockProcessingWaitTime)

	fmt.Printf("%+v\n", signedTestOrder.Sender.Hex())

	_, err := client.GetStats(ctx)
	require.NoError(t, err)

	// Send the "AddOrders" request to the GraphQL server.
	validationResponse, err := client.AddOrdersV4(ctx, []*zeroex.SignedOrderV4{signedTestOrder})
	require.NoError(t, err)
	fmt.Printf("validation response: %+v\n", validationResponse)

	fmt.Println(validationResponse.Rejected)
	fmt.Println(validationResponse.Accepted)

	// Ensure that the validation results contain only the order that was
	// sent to the GraphQL server and that the order was marked as valid.
	require.Len(t, validationResponse.Accepted, 1)
	assert.Len(t, validationResponse.Rejected, 0)
	expectedFillableTakerAssetAmount := signedTestOrder.TakerAmount
	expectedOrderHash, err := signedTestOrder.ComputeOrderHash()
	require.NoError(t, err, "could not compute order hash for standalone order")
	fmt.Printf("the hash is: %+v and fillable is: %+v\n", expectedOrderHash, expectedFillableTakerAssetAmount)
	// expectedAcceptedOrder := &gqlclient.OrderWithMetadata{
	// 	ChainID:                  signedTestOrder.ChainID,
	// 	ExchangeAddress:          signedTestOrder.ExchangeAddress,
	// 	MakerAddress:             signedTestOrder.MakerAddress,
	// 	MakerAssetData:           signedTestOrder.MakerAssetData,
	// 	MakerAssetAmount:         signedTestOrder.MakerAssetAmount,
	// 	MakerFeeAssetData:        signedTestOrder.MakerFeeAssetData,
	// 	MakerFee:                 signedTestOrder.MakerFee,
	// 	TakerAddress:             signedTestOrder.TakerAddress,
	// 	TakerAssetData:           signedTestOrder.TakerAssetData,
	// 	TakerAssetAmount:         signedTestOrder.TakerAssetAmount,
	// 	TakerFeeAssetData:        signedTestOrder.TakerFeeAssetData,
	// 	TakerFee:                 signedTestOrder.TakerFee,
	// 	SenderAddress:            signedTestOrder.SenderAddress,
	// 	FeeRecipientAddress:      signedTestOrder.FeeRecipientAddress,
	// 	ExpirationTimeSeconds:    signedTestOrder.ExpirationTimeSeconds,
	// 	Salt:                     signedTestOrder.Salt,
	// 	Signature:                signedTestOrder.Signature,
	// 	Hash:                     expectedOrderHash,
	// 	FillableTakerAssetAmount: expectedFillableTakerAssetAmount,
	// }
	// assert.Equal(t, expectedAcceptedOrder, accepted.Order, "accepted.Order")
	// assert.Equal(t, true, accepted.IsNew, "accepted.IsNew")
	findResults, err := client.FindOrdersV4(ctx)
	require.NoError(t, err)
	fmt.Println(len(findResults))

	cancel()
	wg.Wait()
}

func TestGetOrderV4(t *testing.T) {
	teardownSubTest := setupSubTest(t)
	defer teardownSubTest(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	wg := &sync.WaitGroup{}
	client, _ := buildAndStartGraphQLServer(t, ctx, wg)

	orderOptions := orderopts.SetupMakerState(true)
	signedTestOrder := scenario.NewSignedTestOrderV4(t, orderOptions)
	time.Sleep(blockProcessingWaitTime)

	validationResponse, err := client.AddOrdersV4(ctx, []*zeroex.SignedOrderV4{signedTestOrder})
	require.NoError(t, err)
	assert.Len(t, validationResponse.Accepted, 1)
	assert.Len(t, validationResponse.Rejected, 0)

	acceptedOne := validationResponse.Accepted[0].Order.Hash
	fmt.Printf("the hash is: %+v\n", acceptedOne)

	expectedHash, err := signedTestOrder.ComputeOrderHash()
	require.NoError(t, err)
	actualOrder, err := client.GetOrderV4(ctx, expectedHash)
	require.NoError(t, err)

	fmt.Printf("%+v\n", actualOrder)

	cancel()
	wg.Wait()
}
