package scenario

import (
	"context"
	"math/big"
	"testing"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

// SetupBalancesAndAllowances sets up the proper balance/allowance for the maker/taker of
// ZRX/WETH respectively so that the created orders are fillable.
func SetupBalancesAndAllowances(t *testing.T, makerAddress, takerAddress common.Address, wethAmount *big.Int, zrxAmount *big.Int) {
	// All 1 billion ZRX start in this address
	zrxCoinbase := constants.GanacheAccount0
	if makerAddress == zrxCoinbase {
		t.Errorf("makerAddress cannot be set to the ZRX coinbase address (e.g., the address with the 1 billion ZRX at Genesis)")
	}

	ganacheAddresses := ethereum.NetworkIDToContractAddresses[constants.TestNetworkID]

	ethClient, err := ethclient.Dial(constants.GanacheEndpoint)
	require.NoError(t, err)

	weth9, err := wrappers.NewWETH9(ganacheAddresses.WETH9, ethClient)
	require.NoError(t, err)

	// Convert ETH-WETH
	opts := &bind.TransactOpts{
		From:   takerAddress,
		Value:  wethAmount,
		Signer: getTestSignerFn(takerAddress),
	}
	txn, err := weth9.Deposit(opts)
	require.NoError(t, err)
	receipt, err := bind.WaitMined(context.Background(), ethClient, txn)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, uint64(1))

	zrx, err := wrappers.NewZRXToken(ganacheAddresses.ZRXToken, ethClient)
	require.NoError(t, err)

	// Transfer ZRX to makerAddress
	opts = &bind.TransactOpts{
		From:   zrxCoinbase,
		Signer: getTestSignerFn(zrxCoinbase),
	}
	txn, err = zrx.Transfer(opts, makerAddress, zrxAmount)
	require.NoError(t, err)
	receipt, err = bind.WaitMined(context.Background(), ethClient, txn)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, uint64(1))

	// SET ZRX allowance
	opts = &bind.TransactOpts{
		From:   makerAddress,
		Signer: getTestSignerFn(makerAddress),
	}
	txn, err = zrx.Approve(opts, ganacheAddresses.ERC20Proxy, zrxAmount)
	require.NoError(t, err)
	receipt, err = bind.WaitMined(context.Background(), ethClient, txn)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, uint64(1))

	// SET WETH allowance
	opts = &bind.TransactOpts{
		From:   takerAddress,
		Signer: getTestSignerFn(takerAddress),
	}
	txn, err = weth9.Approve(opts, ganacheAddresses.ERC20Proxy, wethAmount)
	require.NoError(t, err)
	receipt, err = bind.WaitMined(context.Background(), ethClient, txn)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, uint64(1))
}

func getTestSignerFn(signerAddress common.Address) func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
	return func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		testSigner := ethereum.NewTestSigner()
		signature, err := testSigner.(*ethereum.TestSigner).SignTx(signer.Hash(tx).Bytes(), signerAddress)
		if err != nil {
			return nil, err
		}
		return tx.WithSignature(signer, signature)
	}
}
