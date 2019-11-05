package scenario

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/signer"
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
		MakerAddress:          makerAddress,
		TakerAddress:          constants.NullAddress,
		SenderAddress:         constants.NullAddress,
		FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		MakerAssetData:        common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		TakerAssetData:        common.Hex2Bytes("f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
		Salt:                  big.NewInt(1548619145450),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      zrxAmount,
		TakerAssetAmount:      wethAmount,
		ExpirationTimeSeconds: big.NewInt(time.Now().Add(24 * time.Hour).Unix()),
		ExchangeAddress:       ethereum.ChainIDToContractAddresses[constants.TestChainID].Exchange,
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

	ganacheAddresses := ethereum.ChainIDToContractAddresses[constants.TestChainID]

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

// CreateSignedTestOrderWithExpirationTime creates a valid 0x orders where the maker wishes to trade ZRX for WETH
func CreateSignedTestOrderWithExpirationTime(t *testing.T, ethClient *ethclient.Client, makerAddress, takerAddress common.Address, expirationTime time.Time) *zeroex.SignedOrder {
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
		TakerAssetAmount:      big.NewInt(1000),
		ExpirationTimeSeconds: big.NewInt(expirationTime.Unix()),
		ExchangeAddress:       ethereum.ChainIDToContractAddresses[constants.TestChainID].Exchange,
	}

	// Sign Order
	signedTestOrder, err := zeroex.SignTestOrder(testOrder)
	require.NoError(t, err, "could not sign order")

	return signedTestOrder
}

// CreateWETHForZRXSignedTestOrder creates a valid 0x orders where the maker wishes to trade WETH for ZRX
func CreateWETHForZRXSignedTestOrder(t *testing.T, ethClient *ethclient.Client, makerAddress, takerAddress common.Address, wethAmount *big.Int, zrxAmount *big.Int) *zeroex.SignedOrder {
	// Create order
	testOrder := &zeroex.Order{
		MakerAddress:          makerAddress,
		TakerAddress:          constants.NullAddress,
		SenderAddress:         constants.NullAddress,
		FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		MakerAssetData:        common.Hex2Bytes("f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082"),
		TakerAssetData:        common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		Salt:                  big.NewInt(1548619145450),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      wethAmount,
		TakerAssetAmount:      zrxAmount,
		ExpirationTimeSeconds: big.NewInt(time.Now().Add(24 * time.Hour).Unix()),
		ExchangeAddress:       ethereum.ChainIDToContractAddresses[constants.TestChainID].Exchange,
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

	ganacheAddresses := ethereum.ChainIDToContractAddresses[constants.TestChainID]

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

// CreateNFTForZRXSignedTestOrder creates a valid 0x orders where the maker wishes to trade an NFT for ZRX
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
		MakerAddress:          makerAddress,
		TakerAddress:          constants.NullAddress,
		SenderAddress:         constants.NullAddress,
		FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		MakerAssetData:        makerAssetData,
		TakerAssetData:        common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		Salt:                  big.NewInt(1548619145450),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(1),
		TakerAssetAmount:      zrxAmount,
		ExpirationTimeSeconds: big.NewInt(time.Now().Add(24 * time.Hour).Unix()),
		ExchangeAddress:       ethereum.ChainIDToContractAddresses[constants.TestChainID].Exchange,
	}

	// Sign Order
	signedTestOrder, err := zeroex.SignTestOrder(testOrder)
	require.NoError(t, err, "could not sign order")

	// Set up balances/allowances

	ganacheAddresses := ethereum.ChainIDToContractAddresses[constants.TestChainID]

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

// CreateERC1155ForZRXSignedTestOrder creates a valid 0x orders where the maker wishes to trade an ERC1155 for ZRX
func CreateERC1155ForZRXSignedTestOrder(t *testing.T, ethClient *ethclient.Client, makerAddress, takerAddress common.Address, tokenID *big.Int, zrxAmount, erc1155FungibleAmount *big.Int) *zeroex.SignedOrder {
	erc1155Mintable, err := wrappers.NewERC1155Mintable(constants.GanacheDummyERC1155MintableAddress, ethClient)
	require.NoError(t, err)

	// Withdraw maker's WETH
	// HACK(fabio): For some reason the txn fails with "out of gas" error with the
	// estimated gas amount
	gasLimit := uint64(50000)
	makerOpts := &bind.TransactOpts{
		From:     makerAddress,
		Signer:   GetTestSignerFn(makerAddress),
		GasLimit: gasLimit,
	}
	uri := ""
	txn, err := erc1155Mintable.CreateWithType(makerOpts, tokenID, uri)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	txn, err = erc1155Mintable.MintFungible(makerOpts, tokenID, []common.Address{makerAddress}, []*big.Int{erc1155FungibleAmount})
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	ganacheAddresses := ethereum.ChainIDToContractAddresses[constants.TestChainID]

	devUtils, err := wrappers.NewDevUtils(ganacheAddresses.DevUtils, ethClient)
	require.NoError(t, err)

	callOpts := &bind.CallOpts{
		From: makerAddress,
	}
	erc1155AssetData, err := devUtils.EncodeERC1155AssetData(
		callOpts,
		constants.GanacheDummyERC1155MintableAddress,
		[]*big.Int{tokenID},
		[]*big.Int{erc1155FungibleAmount},
		[]byte{},
	)
	require.NoError(t, err)

	// Create order
	testOrder := &zeroex.Order{
		MakerAddress:          makerAddress,
		TakerAddress:          constants.NullAddress,
		SenderAddress:         constants.NullAddress,
		FeeRecipientAddress:   common.HexToAddress("0xa258b39954cef5cb142fd567a46cddb31a670124"),
		MakerAssetData:        erc1155AssetData,
		TakerAssetData:        common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c"),
		Salt:                  big.NewInt(1548619145450),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(1),
		TakerAssetAmount:      zrxAmount,
		ExpirationTimeSeconds: big.NewInt(time.Now().Add(24 * time.Hour).Unix()),
		ExchangeAddress:       ethereum.ChainIDToContractAddresses[constants.TestChainID].Exchange,
	}

	// Sign Order
	signedTestOrder, err := zeroex.SignTestOrder(testOrder)
	require.NoError(t, err, "could not sign order")

	// Set up balances/allowances

	// SET ERC1155 allowance
	txn, err = erc1155Mintable.SetApprovalForAll(makerOpts, ganacheAddresses.ERC1155Proxy, true)
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
	return func(s types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
		testSigner := signer.NewTestSigner()
		signature, err := testSigner.(*signer.TestSigner).SignTx(s.Hash(tx).Bytes(), signerAddress)
		if err != nil {
			return nil, err
		}
		return tx.WithSignature(s, signature)
	}
}

func waitTxnSuccessfullyMined(t *testing.T, ethClient *ethclient.Client, txn *types.Transaction) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancelFn()
	receipt, err := bind.WaitMined(ctx, ethClient, txn)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, uint64(1))
}
