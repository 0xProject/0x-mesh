package scenario

import (
	"bytes"
	"context"
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

var (
	ganacheAddresses = ethereum.GanacheAddresses
	ZRXAssetData     = common.Hex2Bytes("f47261b0000000000000000000000000871dd7c2b4b25e1aa18728e9d5f2af4c4e431f5c")
	WETHAssetData    = common.Hex2Bytes("f47261b00000000000000000000000000b1ba0af832d7c05fd64161e0db78e85978e8082")
)

func NewTestOrder() *zeroex.Order {
	return &zeroex.Order{
		ChainID:               big.NewInt(constants.TestChainID),
		MakerAddress:          constants.GanacheAccount1,
		TakerAddress:          constants.NullAddress,
		SenderAddress:         constants.NullAddress,
		FeeRecipientAddress:   constants.NullAddress,
		MakerAssetData:        ZRXAssetData,
		MakerFeeAssetData:     constants.NullBytes,
		TakerAssetData:        WETHAssetData,
		TakerFeeAssetData:     constants.NullBytes,
		Salt:                  big.NewInt(int64(time.Now().Nanosecond())),
		MakerFee:              big.NewInt(0),
		TakerFee:              big.NewInt(0),
		MakerAssetAmount:      big.NewInt(100),
		TakerAssetAmount:      big.NewInt(42),
		ExpirationTimeSeconds: big.NewInt(time.Now().Add(24 * time.Hour).Unix()),
		ExchangeAddress:       ganacheAddresses.Exchange,
	}
}

func NewSignedTestOrder(t *testing.T, ethClient *ethclient.Client) *zeroex.SignedOrder {
	order := NewTestOrder()
	signedOrder, err := zeroex.SignTestOrder(order)
	require.NoError(t, err, "could not sign order")
	return signedOrder
}

func NewSignedTestOrderWithState(t *testing.T, ethClient *ethclient.Client) *zeroex.SignedOrder {
	order := NewSignedTestOrder(t, ethClient)
	setupMakerState(t, ethClient, order)
	return order
}

// setupMakerState sets up all the on-chain state in order to make the order fillable. This includes
// setting allowances and transferring the required balances.
func setupMakerState(t *testing.T, ethClient *ethclient.Client, order *zeroex.SignedOrder) {
	assetDataDecoder := zeroex.NewAssetDataDecoder()
	assetDataName, err := assetDataDecoder.GetName(order.MakerAssetData)
	require.NoError(t, err)

	switch assetDataName {
	case "ERC20Token":
		if bytes.Equal(order.MakerAssetData, ZRXAssetData) {
			setZRXBalanceAndAllowance(t, ethClient, order.MakerAddress, order.MakerAssetAmount)
		} else if bytes.Equal(order.MakerAssetData, WETHAssetData) {
			setWETHBalanceAndAllowance(t, ethClient, order.MakerAddress, order.MakerAssetAmount)
		} else {
			t.Errorf("scneario: cannot setup on-chain state for ERC20 assetdata (unsupported token): %s", order.MakerAssetData)
		}
	case "ERC721Token":
		var decodedAssetData zeroex.ERC721AssetData
		require.NoError(t, assetDataDecoder.Decode(order.MakerAssetData, &decodedAssetData))
		if decodedAssetData.Address.Hex() == constants.GanacheDummyERC721TokenAddress.Hex() {
			setDummyERC721BalanceAndAllowance(t, ethClient, order.MakerAddress, decodedAssetData.TokenId)
		} else {
			t.Errorf("scneario: cannot setup on-chain state for ERC721 assetdata (only DummyERC721Token is supported): %s", order.MakerAssetData)
		}
	case "ERC1155Assets":
		var decodedAssetData zeroex.ERC1155AssetData
		require.NoError(t, assetDataDecoder.Decode(order.MakerAssetData, &decodedAssetData))
		if decodedAssetData.Address.Hex() == constants.GanacheDummyERC1155MintableAddress.Hex() {
			setDummyERC1155BalanceAndAllowance(t, ethClient, order.MakerAddress, decodedAssetData.Ids, decodedAssetData.Values)
		} else {
			t.Errorf("scneario: cannot setup on-chain state for ERC1155 assetdata (only DummyERC1155Mintable is supported): %s", order.MakerAssetData)
		}
	default:
		t.Errorf("scenario: cannot setup on-chain state for unsupported assetdata: %s", order.MakerAssetData)
	}

}

// TODO(albrow): Implement setting up taker state.
// // setupTakerState sets up all the on-chain state in order to make the order fillable. This includes
// // setting allowances and transferring the required balances.
// func setupTakerState(t *testing.T, ethClient *ethclient.Client, order *zeroex.SignedOrder, takerAddress common.Address) {
// 	// Set maker allowance
// 	opts = &bind.TransactOpts{
// 		From:   order.MakerAddress,
// 		Signer: GetTestSignerFn(order.MakerAddress),
// 	}
// 	txn, err = zrx.Approve(opts, ganacheAddresses.ERC20Proxy, order.MakerAssetAmount)
// 	require.NoError(t, err)
// 	waitTxnSuccessfullyMined(t, ethClient, txn)
// }

// setWETHBalanceAndAllowance unwraps amount WETH for traderAddress. In other words, the given amount
// will be added to traderAddress's WETH balance.
func setWETHBalanceAndAllowance(t *testing.T, ethClient *ethclient.Client, traderAddress common.Address, amount *big.Int) {
	weth9, err := wrappers.NewWETH9(ganacheAddresses.WETH9, ethClient)
	require.NoError(t, err)

	// Convert ETH to WETH
	opts := &bind.TransactOpts{
		From:   traderAddress,
		Value:  amount,
		Signer: GetTestSignerFn(traderAddress),
	}
	txn, err := weth9.Deposit(opts)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)
}

// setZRXBalanceAndAllowance transfers amount ZRX to traderAddress and sets the appropriate allowance.
func setZRXBalanceAndAllowance(t *testing.T, ethClient *ethclient.Client, traderAddress common.Address, amount *big.Int) {
	zrx, err := wrappers.NewZRXToken(ganacheAddresses.ZRXToken, ethClient)
	require.NoError(t, err)

	// Transfer ZRX to traderAddress
	zrxCoinbase := constants.GanacheAccount0
	opts := &bind.TransactOpts{
		From:   zrxCoinbase,
		Signer: GetTestSignerFn(zrxCoinbase),
	}
	txn, err := zrx.Transfer(opts, traderAddress, amount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	// Set ZRX allowance
	opts = &bind.TransactOpts{
		From:   traderAddress,
		Signer: GetTestSignerFn(traderAddress),
	}
	txn, err = zrx.Approve(opts, ganacheAddresses.ERC20Proxy, amount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)
}

func setDummyERC721BalanceAndAllowance(t *testing.T, ethClient *ethclient.Client, traderAddress common.Address, tokenID *big.Int) {
	// Transfer NFT to traderAddress
	dummyERC721Token, err := wrappers.NewDummyERC721Token(constants.GanacheDummyERC721TokenAddress, ethClient)
	require.NoError(t, err)

	opts := &bind.TransactOpts{
		From:   traderAddress,
		Signer: GetTestSignerFn(traderAddress),
	}
	txn, err := dummyERC721Token.Mint(opts, traderAddress, tokenID)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)

	// Set allowance
	txn, err = dummyERC721Token.Approve(opts, ganacheAddresses.ERC721Proxy, tokenID)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)
}

func setDummyERC1155BalanceAndAllowance(t *testing.T, ethClient *ethclient.Client, traderAddress common.Address, tokenIDs []*big.Int, amounts []*big.Int) {
	// Mint the necessary ERC1155 tokens
	erc1155Mintable, err := wrappers.NewERC1155Mintable(constants.GanacheDummyERC1155MintableAddress, ethClient)
	require.NoError(t, err)

	if len(tokenIDs) != len(amounts) {
		t.Errorf("scenario: tokenIDs and amounts are not the same length (%d and %d respectively)", len(tokenIDs), len(amounts))
	}

	// HACK(fabio): For some reason the txn fails with "out of gas" error with the
	// estimated gas amount
	gasLimit := uint64(50000)
	opts := &bind.TransactOpts{
		From:     traderAddress,
		Signer:   GetTestSignerFn(traderAddress),
		GasLimit: gasLimit,
	}

	for i, tokenID := range tokenIDs {
		amount := amounts[i]

		uri := ""
		txn, err := erc1155Mintable.CreateWithType(opts, tokenID, uri)
		require.NoError(t, err)
		waitTxnSuccessfullyMined(t, ethClient, txn)

		txn, err = erc1155Mintable.MintFungible(opts, tokenID, []common.Address{traderAddress}, []*big.Int{amount})
		require.NoError(t, err)
		waitTxnSuccessfullyMined(t, ethClient, txn)
	}

	// Set ERC1155 allowance
	// HACK(albrow): erc1155Mintable does not allow setting allowance per token id.
	txn, err := erc1155Mintable.SetApprovalForAll(opts, ganacheAddresses.ERC1155Proxy, true)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, ethClient, txn)
}

// TODO(albrow): Create a way to supply the following options:
//
//    - makerAsset
//    - makerAssetAmount
//    - takerAsset
//    - takerAssetAmount
//    - expirationTime
//    - makerFeeAsset and makerFeeAssetAmount
//
// Maybe with chaining like in libp2p?
//

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
