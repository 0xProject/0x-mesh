package scenario

import (
	"context"
	"fmt"
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

// CreateZRXForWETHSignedTestOrder creates a valid 0x orders where the maker wishes to trade ZRX for WETH
func CreateZRXForWETHSignedTestOrder(t *testing.T, ethClient *ethclient.Client, makerAddress, takerAddress common.Address, wethAmount *big.Int, zrxAmount *big.Int) *zeroex.SignedOrder {
	// Create order
	testOrder := &zeroex.Order{
		MakerAddress:  makerAddress,
		TakerAddress:  constants.NullAddress,
		SenderAddress: constants.NullAddress,
		// TODO(albrow): We should remove MakerFeeAssetData and TakerFeeAssetData after the DevUtils contract is fixed
		FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		MakerAssetData:        common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		MakerFeeAssetData:     common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		TakerAssetData:        common.Hex2Bytes("f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
		TakerFeeAssetData:     common.Hex2Bytes("f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
		Salt:                  big.NewInt(1548619145450),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      zrxAmount,
		TakerAssetAmount:      wethAmount,
		ExpirationTimeSeconds: big.NewInt(time.Now().Add(24 * time.Hour).Unix()),
		DomainHash:            constants.NetworkIDToDomainHash[constants.TestNetworkID],
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

	weth9, err := wrappers.NewWETH9(ganacheAddresses.WETH9, ethClient)
	require.NoError(t, err)

	// Convert ETH-WETH
	opts := &bind.TransactOpts{
		From:   takerAddress,
		Value:  wethAmount,
		Signer: GetTestSignerFn(takerAddress),
	}
	txn, err := weth9.Deposit(opts)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	zrx, err := wrappers.NewZRXToken(ganacheAddresses.ZRXToken, ethClient)
	require.NoError(t, err)

	// Transfer ZRX to makerAddress
	opts = &bind.TransactOpts{
		From:   zrxCoinbase,
		Signer: GetTestSignerFn(zrxCoinbase),
	}
	txn, err = zrx.Transfer(opts, makerAddress, zrxAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	// SET ZRX allowance
	opts = &bind.TransactOpts{
		From:   makerAddress,
		Signer: GetTestSignerFn(makerAddress),
	}
	txn, err = zrx.Approve(opts, ganacheAddresses.ERC20Proxy, zrxAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	// SET WETH allowance
	opts = &bind.TransactOpts{
		From:   takerAddress,
		Signer: GetTestSignerFn(takerAddress),
	}
	txn, err = weth9.Approve(opts, ganacheAddresses.ERC20Proxy, wethAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	return signedTestOrder
}

// CreateWETHForZRXSignedTestOrder creates a valid 0x orders where the maker wishes to trade WETH for ZRX
func CreateWETHForZRXSignedTestOrder(t *testing.T, ethClient *ethclient.Client, makerAddress, takerAddress common.Address, wethAmount *big.Int, zrxAmount *big.Int) *zeroex.SignedOrder {
	// Create order
	testOrder := &zeroex.Order{
		MakerAddress:        makerAddress,
		TakerAddress:        constants.NullAddress,
		SenderAddress:       constants.NullAddress,
		FeeRecipientAddress: common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		// TODO(albrow): We should remove MakerFeeAssetData and TakerFeeAssetData after the DevUtils contract is fixed
		MakerAssetData:        common.Hex2Bytes("f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
		MakerFeeAssetData:     common.Hex2Bytes("f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
		TakerAssetData:        common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		TakerFeeAssetData:     common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		Salt:                  big.NewInt(1548619145450),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      wethAmount,
		TakerAssetAmount:      zrxAmount,
		ExpirationTimeSeconds: big.NewInt(time.Now().Add(24 * time.Hour).Unix()),
		DomainHash:            constants.NetworkIDToDomainHash[constants.TestNetworkID],
	}

	// Sign Order
	signedTestOrder, err := zeroex.SignTestOrder(testOrder)
	require.NoError(t, err, "could not sign order")

	// Set up balances/allowances

	// All 1 billion ZRX start in this address
	zrxCoinbase := constants.GanacheAccount0
	if takerAddress == zrxCoinbase {
		t.Errorf("takerAddress cannot be set to the ZRX coinbase address (e.g., the address with the 1 billion ZRX at Genesis)")
	}

	ganacheAddresses := ethereum.NetworkIDToContractAddresses[constants.TestNetworkID]

	weth9, err := wrappers.NewWETH9(ganacheAddresses.WETH9, ethClient)
	require.NoError(t, err)

	// Convert ETH-WETH
	opts := &bind.TransactOpts{
		From:   makerAddress,
		Value:  wethAmount,
		Signer: GetTestSignerFn(makerAddress),
	}
	txn, err := weth9.Deposit(opts)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	zrx, err := wrappers.NewZRXToken(ganacheAddresses.ZRXToken, ethClient)
	require.NoError(t, err)

	// Transfer ZRX to takerAddress
	opts = &bind.TransactOpts{
		From:   zrxCoinbase,
		Signer: GetTestSignerFn(zrxCoinbase),
	}
	txn, err = zrx.Transfer(opts, takerAddress, zrxAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	// SET ZRX allowance
	opts = &bind.TransactOpts{
		From:   takerAddress,
		Signer: GetTestSignerFn(takerAddress),
	}
	txn, err = zrx.Approve(opts, ganacheAddresses.ERC20Proxy, zrxAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	// SET WETH allowance
	opts = &bind.TransactOpts{
		From:   makerAddress,
		Signer: GetTestSignerFn(makerAddress),
	}
	txn, err = weth9.Approve(opts, ganacheAddresses.ERC20Proxy, wethAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	return signedTestOrder
}

// CreateNFTForZRXSignedTestOrder creates a valid 0x orders where the maker wishes to trade WETH for ZRX
func CreateNFTForZRXSignedTestOrder(t *testing.T, ethClient *ethclient.Client, makerAddress, takerAddress common.Address, tokenID *big.Int, zrxAmount *big.Int) *zeroex.SignedOrder {
	dummyERC721Token, err := wrappers.NewDummyERC721Token(constants.GanacheDummyERC721TokenAddress, ethClient)
	require.NoError(t, err)

	makerOpts := &bind.TransactOpts{
		From:   makerAddress,
		Signer: GetTestSignerFn(makerAddress),
	}
	txn, err := dummyERC721Token.Mint(makerOpts, makerAddress, tokenID)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	makerAssetDataHex := fmt.Sprintf("%s000000000000000000000000%s000000000000000000000000000000000000000000000000000000000000000%s", zeroex.ERC721AssetDataID, constants.GanacheDummyERC721TokenAddress.Hex()[2:], tokenID)
	makerAssetData := common.Hex2Bytes(
		makerAssetDataHex,
	)

	// Create order
	testOrder := &zeroex.Order{
		MakerAddress:        makerAddress,
		TakerAddress:        constants.NullAddress,
		SenderAddress:       constants.NullAddress,
		FeeRecipientAddress: common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		// TODO(albrow): We should remove MakerFeeAssetData and TakerFeeAssetData after the DevUtils contract is fixed
		MakerAssetData:        makerAssetData,
		MakerFeeAssetData:     makerAssetData,
		TakerAssetData:        common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		TakerFeeAssetData:     common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		Salt:                  big.NewInt(1548619145450),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(1),
		TakerAssetAmount:      zrxAmount,
		ExpirationTimeSeconds: big.NewInt(time.Now().Add(24 * time.Hour).Unix()),
		DomainHash:            constants.NetworkIDToDomainHash[constants.TestNetworkID],
	}

	// Sign Order
	signedTestOrder, err := zeroex.SignTestOrder(testOrder)
	require.NoError(t, err, "could not sign order")

	// Set up balances/allowances

	ganacheAddresses := ethereum.NetworkIDToContractAddresses[constants.TestNetworkID]

	// SET NFT allowance
	txn, err = dummyERC721Token.SetApprovalForAll(makerOpts, ganacheAddresses.ERC721Proxy, true)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	// All 1 billion ZRX start in this address
	zrxCoinbase := constants.GanacheAccount0
	if takerAddress == zrxCoinbase {
		t.Errorf("takerAddress cannot be set to the ZRX coinbase address (e.g., the address with the 1 billion ZRX at Genesis)")
	}

	zrx, err := wrappers.NewZRXToken(ganacheAddresses.ZRXToken, ethClient)
	require.NoError(t, err)

	// Transfer ZRX to takerAddress
	opts := &bind.TransactOpts{
		From:   zrxCoinbase,
		Signer: GetTestSignerFn(zrxCoinbase),
	}
	txn, err = zrx.Transfer(opts, takerAddress, zrxAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	// SET ZRX allowance
	opts = &bind.TransactOpts{
		From:   takerAddress,
		Signer: GetTestSignerFn(takerAddress),
	}
	txn, err = zrx.Approve(opts, ganacheAddresses.ERC20Proxy, zrxAmount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	return signedTestOrder
}

// GetTestSignerFn returns a test signer function that can be used to sign Ethereum transactions
func GetTestSignerFn(signerAddress common.Address) func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
	return func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		testSigner := ethereum.NewTestSigner()
		signature, err := testSigner.(*ethereum.TestSigner).SignTx(signer.Hash(tx).Bytes(), signerAddress)
		if err != nil {
			return nil, err
		}
		return tx.WithSignature(signer, signature)
	}
}

func waitTxnSuccessfullyMined(t *testing.T, ethClient *ethclient.Client, txn *types.Transaction) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancelFn()
	receipt, err := bind.WaitMined(ctx, ethClient, txn)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, uint64(1))
}
