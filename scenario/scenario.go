// Package scenario allows creating orders for testing purposes with a variety of options.
// It also supports setting up the necessary on-chain state for both the taker and maker.
package scenario

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/0xProject/0x-mesh/ethereum"
	"github.com/0xProject/0x-mesh/ethereum/signer"
	"github.com/0xProject/0x-mesh/ethereum/wrappers"
	"github.com/0xProject/0x-mesh/scenario/orderopts"
	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

var (
	ethClient        *ethclient.Client
	ganacheAddresses = ethereum.GanacheAddresses
	ZRXAssetData     = constants.ZRXAssetData
	WETHAssetData    = constants.WETHAssetData
)

func init() {
	rpcClient, err := rpc.Dial(constants.GanacheEndpoint)
	if err != nil {
		panic(err)
	}
	ethClient = ethclient.NewClient(rpcClient)
}

func defaultTestOrder() *zeroex.Order {
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

func defaultConfig() *orderopts.Config {
	return &orderopts.Config{
		Order:             defaultTestOrder(),
		OrderV4:           defaultTestOrderV4(),
		SetupMakerState:   false,
		SetupTakerAddress: constants.NullAddress,
	}
}

func NewTestOrder(t *testing.T, opts ...orderopts.Option) *zeroex.Order {
	cfg := defaultConfig()
	require.NoError(t, cfg.Apply(opts...))
	return newTestOrder(cfg)
}

func newTestOrder(cfg *orderopts.Config) *zeroex.Order {
	return cfg.Order
}

func NewSignedTestOrder(t *testing.T, opts ...orderopts.Option) *zeroex.SignedOrder {
	cfg := defaultConfig()
	require.NoError(t, cfg.Apply(opts...))

	order := newTestOrder(cfg)
	signedOrder, err := zeroex.SignTestOrder(order)
	require.NoError(t, err, "could not sign order")

	if cfg.SetupMakerState {
		setupMakerState(t, signedOrder)
	}
	if cfg.SetupTakerAddress != constants.NullAddress {
		setupTakerState(t, signedOrder, cfg.SetupTakerAddress)
	}

	return signedOrder
}

// NewSignedTestOrdersBatch efficiently creates numOrders orders with independent options.
// If the options require setting up maker or taker state, that state will be set up efficiently
// with one transaction per address.
//
// optionsForIndex is a function which returns the options for creating the order at a specific
// index (between 0 and numOrders). For example, you can create ERC721 orders which each have a unique
// token ID. optionsForIndex can be nil to always use the default options. It can return nil to
// use the default options for an order at a specific index.
func NewSignedTestOrdersBatch(t *testing.T, numOrders int, optionsForIndex func(index int) []orderopts.Option) []*zeroex.SignedOrder {
	allRequiredBalances := map[common.Address]*tokenBalances{}

	allOrders := make([]*zeroex.SignedOrder, numOrders)
	for i := 0; i < numOrders; i++ {
		// Apply the options (if any) for the order we will create at this index.
		cfg := defaultConfig()
		if optionsForIndex != nil {
			opts := optionsForIndex(i)
			if opts != nil {
				require.NoError(t, cfg.Apply(opts...))
			}
		}

		// Create the order based on the cfg.
		order := newTestOrder(cfg)
		signedOrder, err := zeroex.SignTestOrder(order)
		require.NoError(t, err, "could not sign order")
		allOrders[i] = signedOrder

		// Add maker and taker balances as needed to the set of required balances.
		if cfg.SetupMakerState {
			makerBalancesForThisOrder := requiredMakerBalances(t, signedOrder)
			makerBalances, found := allRequiredBalances[signedOrder.MakerAddress]
			if !found {
				allRequiredBalances[order.MakerAddress] = makerBalancesForThisOrder
			} else {
				makerBalances.add(makerBalancesForThisOrder)
			}
		}
		if cfg.SetupTakerAddress != constants.NullAddress {
			takerBalancesForThisOrder := requiredTakerBalances(t, signedOrder)
			takerBalances, found := allRequiredBalances[cfg.SetupTakerAddress]
			if !found {
				allRequiredBalances[cfg.SetupTakerAddress] = takerBalancesForThisOrder
			} else {
				takerBalances.add(takerBalancesForThisOrder)
			}
		}
	}

	// Setup all the required balances.
	for traderAddress, requiredBalances := range allRequiredBalances {
		setupBalanceAndAllowance(t, traderAddress, requiredBalances)
	}

	return allOrders
}

// OptionsForAll is a convenience function which can be used in combination with NewSignedTestOrdersBatch
// when you want all orders to be created with the same options. It returns a function which can be used
// as optionsForIndex which always returns the given options, regardless of the index.
func OptionsForAll(opts ...orderopts.Option) func(_ int) []orderopts.Option {
	return func(_ int) []orderopts.Option {
		return opts
	}
}

type tokenBalances struct {
	zrx           *big.Int
	weth          *big.Int
	erc721Tokens  []*big.Int
	erc1155Tokens []erc1155TokenAmount
}

type erc1155TokenAmount struct {
	tokenID *big.Int
	amount  *big.Int
}

func newTokenBalances() *tokenBalances {
	return &tokenBalances{
		zrx:           big.NewInt(0),
		weth:          big.NewInt(0),
		erc721Tokens:  []*big.Int{},
		erc1155Tokens: []erc1155TokenAmount{},
	}
}

func (x *tokenBalances) add(y *tokenBalances) {
	x.zrx.Add(x.zrx, y.zrx)
	x.weth.Add(x.weth, y.weth)
	x.erc721Tokens = append(x.erc721Tokens, y.erc721Tokens...)
	for _, yToken := range y.erc1155Tokens {
		found := false
		for xIndex, xToken := range x.erc1155Tokens {
			if xToken.tokenID.Cmp(yToken.tokenID) == 0 {
				found = true
				x.erc1155Tokens[xIndex] = erc1155TokenAmount{
					tokenID: xToken.tokenID,
					amount:  xToken.amount.Add(xToken.amount, yToken.amount),
				}
			}
		}
		if !found {
			x.erc1155Tokens = append(x.erc1155Tokens, yToken)
		}
	}
}

func (x *tokenBalances) addTokenAmount(t *testing.T, token common.Address, amount *big.Int) {
	if token == ganacheAddresses.ZRXToken {
		x.zrx.Add(x.zrx, amount)
	} else if token == ganacheAddresses.WETH9 {
		x.weth.Add(x.weth, amount)
	} else {
		t.Fatalf("scenario: cannot setup on-chain state for ERC20 token (unsupported token): %s", token.Hex())
	}
}

var zero = big.NewInt(0)

func isZero(x *big.Int) bool {
	return x.Cmp(zero) == 0
}

// setupMakerState sets up all the on-chain state in order to make the order fillable. This includes
// setting allowances and transferring the required balances.
func setupMakerState(t *testing.T, order *zeroex.SignedOrder) {
	requiredMakerBalances := requiredMakerBalances(t, order)
	setupBalanceAndAllowance(t, order.MakerAddress, requiredMakerBalances)
}

// setupTakerState sets up all the on-chain state needed by taker in order to fill the order.
// This includes setting allowances and transferring the required balances.
func setupTakerState(t *testing.T, order *zeroex.SignedOrder, taker common.Address) {
	requiredTakerBalances := requiredTakerBalances(t, order)
	setupBalanceAndAllowance(t, taker, requiredTakerBalances)
}

func setupBalanceAndAllowance(t *testing.T, traderAddress common.Address, requiredBalances *tokenBalances) {
	if !isZero(requiredBalances.zrx) {
		setZRXBalanceAndAllowance(t, traderAddress, requiredBalances.zrx)
	}
	if !isZero(requiredBalances.weth) {
		setWETHBalanceAndAllowance(t, traderAddress, requiredBalances.weth)
	}
	if len(requiredBalances.erc721Tokens) != 0 {
		for _, tokenId := range requiredBalances.erc721Tokens {
			setDummyERC721BalanceAndAllowance(t, traderAddress, tokenId)
		}
	}
	if len(requiredBalances.erc1155Tokens) != 0 {
		setDummyERC1155BalanceAndAllowance(t, traderAddress, requiredBalances.erc1155Tokens)
	}
}

func requiredMakerBalances(t *testing.T, order *zeroex.SignedOrder) *tokenBalances {
	balances := newTokenBalances()
	balances.add(requiredBalancesForAssetData(t, order.MakerAssetData, order.MakerAssetAmount))
	if len(order.MakerFeeAssetData) != 0 && !isZero(order.MakerFee) {
		balances.add(requiredBalancesForAssetData(t, order.MakerFeeAssetData, order.MakerFee))
	}
	return balances
}

func requiredTakerBalances(t *testing.T, order *zeroex.SignedOrder) *tokenBalances {
	balances := newTokenBalances()
	balances.add(requiredBalancesForAssetData(t, order.TakerAssetData, order.TakerAssetAmount))
	if len(order.TakerFeeAssetData) != 0 && !isZero(order.TakerFee) {
		balances.add(requiredBalancesForAssetData(t, order.TakerFeeAssetData, order.TakerFee))
	}
	return balances
}

func requiredBalancesForAssetData(t *testing.T, assetData []byte, assetAmount *big.Int) *tokenBalances {
	balances := newTokenBalances()
	assetDataDecoder := zeroex.NewAssetDataDecoder()
	assetDataName, err := assetDataDecoder.GetName(assetData)
	require.NoError(t, err)
	switch assetDataName {
	case "ERC20Token":
		if bytes.Equal(assetData, ZRXAssetData) {
			balances.zrx = assetAmount
			return balances
		} else if bytes.Equal(assetData, WETHAssetData) {
			balances.weth = assetAmount
			return balances
		} else {
			t.Fatalf("scenario: cannot setup on-chain state for ERC20 assetdata (unsupported token): %s", common.Bytes2Hex(assetData))
		}
	case "ERC721Token":
		var decodedAssetData zeroex.ERC721AssetData
		require.NoError(t, assetDataDecoder.Decode(assetData, &decodedAssetData))
		if decodedAssetData.Address.Hex() == constants.GanacheDummyERC721TokenAddress.Hex() {
			balances.erc721Tokens = []*big.Int{decodedAssetData.TokenId}
			return balances
		} else {
			t.Fatalf("scneario: cannot setup on-chain state for ERC721 assetdata (only DummyERC721Token is supported): %s", common.Bytes2Hex(assetData))
		}
	case "ERC1155Assets":
		var decodedAssetData zeroex.ERC1155AssetData
		require.NoError(t, assetDataDecoder.Decode(assetData, &decodedAssetData))

		if len(decodedAssetData.Ids) != len(decodedAssetData.Values) {
			t.Fatalf("scenario: tokenIDs and amounts are not the same length (%d and %d respectively)", len(decodedAssetData.Ids), len(decodedAssetData.Values))
		}

		if decodedAssetData.Address.Hex() == constants.GanacheDummyERC1155MintableAddress.Hex() {
			balances.erc1155Tokens = make([]erc1155TokenAmount, len(decodedAssetData.Ids))
			for i, tokenID := range decodedAssetData.Ids {
				totalAmount := big.NewInt(0).Mul(decodedAssetData.Values[i], assetAmount)
				balances.erc1155Tokens[i] = erc1155TokenAmount{
					tokenID: tokenID,
					amount:  totalAmount,
				}
			}
			return balances
		} else {
			t.Fatalf("scneario: cannot setup on-chain state for ERC1155 assetdata (only DummyERC1155Mintable is supported): %s", common.Bytes2Hex(assetData))
		}
	case "StaticCall":
		var decodedAssetData zeroex.StaticCallAssetData
		require.NoError(t, assetDataDecoder.Decode(assetData, &decodedAssetData))
		staticCallDataName, err := assetDataDecoder.GetName(decodedAssetData.StaticCallData)
		require.NoError(t, err)
		if staticCallDataName != "checkGasPrice" {
			t.Fatalf("scneario: cannot setup on-chain state for StaticCall assetdata (only checkGasPrice is supported): (%s) %s", staticCallDataName, common.Bytes2Hex(assetData))
		}
		// Note(albrow): So far there is no additional state required for the types of StaticCall asset data that we support.
		return balances
	}

	// Note(albrow): We don't currently support setting balances and allowances for MAP orders. If needed in
	// the future, we an support MAP orders by recursively calling requiredBalancesForAssetData and adding the
	// result to balances.

	t.Fatalf("scenario: cannot setup on-chain state for unsupported assetdata: (%s) %s", assetDataName, common.Bytes2Hex(assetData))
	return nil
}

// setWETHBalanceAndAllowance unwraps amount WETH for traderAddress. In other words, the given amount
// will be added to traderAddress's WETH balance.
func setWETHBalanceAndAllowance(t *testing.T, traderAddress common.Address, amount *big.Int) {
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
	waitTxnSuccessfullyMined(t, txn)

	// Set WETH allowance
	opts = &bind.TransactOpts{
		From:   traderAddress,
		Signer: GetTestSignerFn(traderAddress),
	}
	// V3
	_, err = weth9.Approve(opts, ganacheAddresses.ERC20Proxy, amount)
	require.NoError(t, err)
	// V4
	_, err = weth9.Approve(opts, ganacheAddresses.ExchangeProxy, amount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, txn)
}

// setZRXBalanceAndAllowance transfers amount ZRX to traderAddress and sets the appropriate allowance.
func setZRXBalanceAndAllowance(t *testing.T, traderAddress common.Address, amount *big.Int) {
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
	waitTxnSuccessfullyMined(t, txn)

	// Set ZRX allowance
	opts = &bind.TransactOpts{
		From:   traderAddress,
		Signer: GetTestSignerFn(traderAddress),
	}
	// V3
	_, err = zrx.Approve(opts, ganacheAddresses.ERC20Proxy, amount)
	require.NoError(t, err)
	// V4
	txn, err = zrx.Approve(opts, ganacheAddresses.ExchangeProxy, amount)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, txn)
}

func setDummyERC721BalanceAndAllowance(t *testing.T, traderAddress common.Address, tokenID *big.Int) {
	// Transfer NFT to traderAddress
	dummyERC721Token, err := wrappers.NewDummyERC721Token(constants.GanacheDummyERC721TokenAddress, ethClient)
	require.NoError(t, err)

	opts := &bind.TransactOpts{
		From:   traderAddress,
		Signer: GetTestSignerFn(traderAddress),
	}
	txn, err := dummyERC721Token.Mint(opts, traderAddress, tokenID)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, txn)

	// Set allowance
	txn, err = dummyERC721Token.Approve(opts, ganacheAddresses.ERC721Proxy, tokenID)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, txn)
}

func setDummyERC1155BalanceAndAllowance(t *testing.T, traderAddress common.Address, tokenAmounts []erc1155TokenAmount) {
	// Mint the necessary ERC1155 tokens
	erc1155Mintable, err := wrappers.NewERC1155Mintable(constants.GanacheDummyERC1155MintableAddress, ethClient)
	require.NoError(t, err)

	// HACK(fabio): For some reason the txn fails with "out of gas" error with the
	// estimated gas amount
	gasLimit := uint64(50000)
	opts := &bind.TransactOpts{
		From:     traderAddress,
		Signer:   GetTestSignerFn(traderAddress),
		GasLimit: gasLimit,
	}

	for _, tokenAmount := range tokenAmounts {
		uri := ""
		txn, err := erc1155Mintable.CreateWithType(opts, tokenAmount.tokenID, uri)
		require.NoError(t, err)
		waitTxnSuccessfullyMined(t, txn)

		txn, err = erc1155Mintable.MintFungible(opts, tokenAmount.tokenID, []common.Address{traderAddress}, []*big.Int{tokenAmount.amount})
		require.NoError(t, err)
		waitTxnSuccessfullyMined(t, txn)
	}

	// Set ERC1155 allowance
	// HACK(albrow): erc1155Mintable does not allow setting allowance per token id.
	txn, err := erc1155Mintable.SetApprovalForAll(opts, ganacheAddresses.ERC1155Proxy, true)
	require.NoError(t, err)
	waitTxnSuccessfullyMined(t, txn)
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

func GetDummyERC721AssetData(tokenID *big.Int) []byte {
	makerAssetDataHex := fmt.Sprintf("%s000000000000000000000000%s000000000000000000000000000000000000000000000000000000000000000%s", zeroex.ERC721AssetDataID, constants.GanacheDummyERC721TokenAddress.Hex()[2:], tokenID)
	return common.Hex2Bytes(makerAssetDataHex)
}

func GetDummyERC1155AssetData(t *testing.T, tokenIDs []*big.Int, amounts []*big.Int) []byte {
	devUtils, err := wrappers.NewDevUtils(ganacheAddresses.DevUtils, ethClient)
	require.NoError(t, err)

	callOpts := &bind.CallOpts{
		From: constants.GanacheAccount1,
	}
	assetData, err := devUtils.EncodeERC1155AssetData(
		callOpts,
		constants.GanacheDummyERC1155MintableAddress,
		tokenIDs,
		amounts,
		[]byte{},
	)
	require.NoError(t, err)
	return assetData
}

func waitTxnSuccessfullyMined(t *testing.T, txn *types.Transaction) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancelFn()
	receipt, err := bind.WaitMined(ctx, ethClient, txn)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, uint64(1))
}
