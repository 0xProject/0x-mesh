// Package scenario allows creating orders for testing purposes with a variety of options.
// It also supports setting up the necessary on-chain state for both the taker and maker.
package scenario

import (
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum/signer"
	"github.com/0xProject/0x-mesh/scenario/orderopts"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/stretchr/testify/require"
)

func defaultTestOrderV4() *zeroex.OrderV4 {
	return &zeroex.OrderV4{
		ChainID:         big.NewInt(constants.TestChainID),
		ExchangeAddress: ganacheAddresses.ExchangeProxy,
		// FIXME: Ganache snapshot currently requires valid token addresses or it will revert
		MakerToken:          ganacheAddresses.WETH9,
		TakerToken:          ganacheAddresses.ZRXToken,
		MakerAmount:         big.NewInt(100),
		TakerAmount:         big.NewInt(42),
		TakerTokenFeeAmount: big.NewInt(0),
		Maker:               constants.GanacheAccount1,
		Taker:               constants.NullAddress,
		Sender:              constants.NullAddress,
		FeeRecipient:        constants.NullAddress,
		Pool:                zeroex.HexToBytes32("0000000000000000000000000000000000000000000000000000000000000000"),
		Expiry:              big.NewInt(time.Now().Add(24 * time.Hour).Unix()),
		Salt:                big.NewInt(int64(time.Now().Nanosecond())),
	}
}

func NewSignedTestOrderV4(t *testing.T, opts ...orderopts.Option) *zeroex.SignedOrderV4 {
	// Generate v4 order with options applied
	cfg := defaultConfig()
	require.NoError(t, cfg.Apply(opts...))
	order := cfg.OrderV4

	// Sign order
	testSigner := signer.NewTestSigner()
	signedOrder, err := zeroex.SignOrderV4(testSigner, order)
	require.NoError(t, err, "could not sign order")

	if cfg.SetupMakerState {
		// sets up all the on-chain state in order to make the order fillable for maker
		requiredMakerBalances := newTokenBalances()
		requiredMakerBalances.addTokenAmount(t, order.MakerToken, order.MakerAmount)
		setupBalanceAndAllowance(t, order.Maker, requiredMakerBalances)
	}
	if cfg.SetupTakerAddress != constants.NullAddress {
		// sets up all the on-chain state in order to make the order fillable taker
		requiredMakerBalances := newTokenBalances()
		requiredMakerBalances.addTokenAmount(t, order.TakerToken, order.TakerAmount)
		requiredMakerBalances.addTokenAmount(t, order.TakerToken, order.TakerTokenFeeAmount)
		setupBalanceAndAllowance(t, order.Taker, requiredMakerBalances)
	}

	return signedOrder
}

// NewSignedTestOrdersBatchV4 creates numOrders orders with independent options.
//
// Unlike the V3 version it currently does not efficiently set balances.
func NewSignedTestOrdersBatchV4(t *testing.T, numOrders int, optionsForIndex func(index int) []orderopts.Option) []*zeroex.SignedOrderV4 {
	allOrders := make([]*zeroex.SignedOrderV4, numOrders)
	for i := 0; i < numOrders; i++ {
		if optionsForIndex != nil {
			allOrders[i] = NewSignedTestOrderV4(t, optionsForIndex(i)...)
		} else {
			allOrders[i] = NewSignedTestOrderV4(t)
		}
	}
	return allOrders
}
