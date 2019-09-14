package scenario

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func CreateZRXForWETHSignedTestOrder(t *testing.T, makerAddress, takerAddress common.Address, wethAmount *big.Int, zrxAmount *big.Int) *zeroex.SignedOrder {
	// Create order
	testOrder := &zeroex.Order{
		MakerAddress:          makerAddress,
		TakerAddress:          constants.NullAddress,
		SenderAddress:         constants.NullAddress,
		FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		MakerAssetData:        common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		TakerAssetData:        common.Hex2Bytes("f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
		Salt:                  big.NewInt(1548619145450),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(1000),
		TakerAssetAmount:      big.NewInt(2000),
		ExpirationTimeSeconds: big.NewInt(time.Now().Add(24 * time.Hour).Unix()),
		ExchangeAddress:       ethereum.NetworkIDToContractAddresses[constants.TestNetworkID].Exchange,
	}

	// Sign Order
	signedTestOrder, err := zeroex.SignTestOrder(testOrder)
	require.NoError(t, err, "could not sign order")

	// Set up balances/allowances

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

	return signedTestOrder
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
