// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package wrappers_v4

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// LibOrderOrder is an auto generated low-level Go binding around an user-defined struct.
type LibOrderOrder struct {
	MakerAddress          common.Address
	TakerAddress          common.Address
	FeeRecipientAddress   common.Address
	SenderAddress         common.Address
	MakerAssetAmount      *big.Int
	TakerAssetAmount      *big.Int
	MakerFee              *big.Int
	TakerFee              *big.Int
	ExpirationTimeSeconds *big.Int
	Salt                  *big.Int
	MakerAssetData        []byte
	TakerAssetData        []byte
	MakerFeeAssetData     []byte
	TakerFeeAssetData     []byte
}

// LibOrderOrderInfo is an auto generated low-level Go binding around an user-defined struct.
type LibOrderOrderInfo struct {
	OrderStatus                 uint8
	OrderHash                   [32]byte
	OrderTakerAssetFilledAmount *big.Int
}

// LibZeroExTransactionZeroExTransaction is an auto generated low-level Go binding around an user-defined struct.
type LibZeroExTransactionZeroExTransaction struct {
	Salt                  *big.Int
	ExpirationTimeSeconds *big.Int
	GasPrice              *big.Int
	SignerAddress         common.Address
	Data                  []byte
}

// WrappersV4ABI is the input ABI used to generate the binding from.
const WrappersV4ABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"exchange_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"chaiBridge_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"dydxBridge_\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[],\"name\":\"EIP712_EXCHANGE_DOMAIN_HASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chaiBridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeAssetProxyId\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC1155AssetData\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenValues\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"callbackData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC20AssetData\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC20BridgeAssetData\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"bridgeAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"bridgeData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC721AssetData\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeMultiAssetData\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"nestedAssetData\",\"type\":\"bytes[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeStaticCallAssetData\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"staticCallTargetAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"staticCallData\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"expectedReturnDataHash\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"transactionData\",\"type\":\"bytes\"}],\"name\":\"decodeZeroExTransactionData\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"functionName\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"dydxBridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenValues\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"callbackData\",\"type\":\"bytes\"}],\"name\":\"encodeERC1155AssetData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"encodeERC20AssetData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"encodeERC721AssetData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"nestedAssetData\",\"type\":\"bytes[]\"}],\"name\":\"encodeMultiAssetData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"staticCallTargetAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"staticCallData\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"expectedReturnDataHash\",\"type\":\"bytes32\"}],\"name\":\"encodeStaticCallAssetData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"erc1155ProxyAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"erc20ProxyAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"erc721ProxyAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"exchangeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getAssetProxyAllowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getBalanceAndAssetProxyAllowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes[]\",\"name\":\"assetData\",\"type\":\"bytes[]\"}],\"name\":\"getBatchAssetProxyAllowances\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"allowances\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes[]\",\"name\":\"assetData\",\"type\":\"bytes[]\"}],\"name\":\"getBatchBalances\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes[]\",\"name\":\"assetData\",\"type\":\"bytes[]\"}],\"name\":\"getBatchBalancesAndAssetProxyAllowances\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"allowances\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"addresses\",\"type\":\"address[]\"}],\"name\":\"getEthBalances\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"exchange\",\"type\":\"address\"}],\"name\":\"getOrderHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"getOrderRelevantState\",\"outputs\":[{\"components\":[{\"internalType\":\"enumLibOrder.OrderStatus\",\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"internalType\":\"structLibOrder.OrderInfo\",\"name\":\"orderInfo\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"fillableTakerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isValidSignature\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"getOrderRelevantStates\",\"outputs\":[{\"components\":[{\"internalType\":\"enumLibOrder.OrderStatus\",\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"internalType\":\"structLibOrder.OrderInfo[]\",\"name\":\"ordersInfo\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"fillableTakerAssetAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bool[]\",\"name\":\"isValidSignature\",\"type\":\"bool[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"}],\"name\":\"getSimulatedOrderMakerTransferResults\",\"outputs\":[{\"internalType\":\"enumLibOrderTransferSimulation.OrderTransferResults\",\"name\":\"orderTransferResults\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"}],\"name\":\"getSimulatedOrderTransferResults\",\"outputs\":[{\"internalType\":\"enumLibOrderTransferSimulation.OrderTransferResults\",\"name\":\"orderTransferResults\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"address[]\",\"name\":\"takerAddresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"}],\"name\":\"getSimulatedOrdersTransferResults\",\"outputs\":[{\"internalType\":\"enumLibOrderTransferSimulation.OrderTransferResults[]\",\"name\":\"orderTransferResults\",\"type\":\"uint8[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLibZeroExTransaction.ZeroExTransaction\",\"name\":\"transaction\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"exchange\",\"type\":\"address\"}],\"name\":\"getTransactionHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"transactionHash\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getTransferableAssetAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"transferableAssetAmount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"revertIfInvalidAssetData\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"staticCallProxyAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// WrappersV4 is an auto generated Go binding around an Ethereum contract.
type WrappersV4 struct {
	WrappersV4Caller     // Read-only binding to the contract
	WrappersV4Transactor // Write-only binding to the contract
	WrappersV4Filterer   // Log filterer for contract events
}

// WrappersV4Caller is an auto generated read-only Go binding around an Ethereum contract.
type WrappersV4Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrappersV4Transactor is an auto generated write-only Go binding around an Ethereum contract.
type WrappersV4Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrappersV4Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WrappersV4Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WrappersV4Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WrappersV4Session struct {
	Contract     *WrappersV4       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WrappersV4CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WrappersV4CallerSession struct {
	Contract *WrappersV4Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// WrappersV4TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WrappersV4TransactorSession struct {
	Contract     *WrappersV4Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// WrappersV4Raw is an auto generated low-level Go binding around an Ethereum contract.
type WrappersV4Raw struct {
	Contract *WrappersV4 // Generic contract binding to access the raw methods on
}

// WrappersV4CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WrappersV4CallerRaw struct {
	Contract *WrappersV4Caller // Generic read-only contract binding to access the raw methods on
}

// WrappersV4TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WrappersV4TransactorRaw struct {
	Contract *WrappersV4Transactor // Generic write-only contract binding to access the raw methods on
}

// NewWrappersV4 creates a new instance of WrappersV4, bound to a specific deployed contract.
func NewWrappersV4(address common.Address, backend bind.ContractBackend) (*WrappersV4, error) {
	contract, err := bindWrappersV4(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WrappersV4{WrappersV4Caller: WrappersV4Caller{contract: contract}, WrappersV4Transactor: WrappersV4Transactor{contract: contract}, WrappersV4Filterer: WrappersV4Filterer{contract: contract}}, nil
}

// NewWrappersV4Caller creates a new read-only instance of WrappersV4, bound to a specific deployed contract.
func NewWrappersV4Caller(address common.Address, caller bind.ContractCaller) (*WrappersV4Caller, error) {
	contract, err := bindWrappersV4(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WrappersV4Caller{contract: contract}, nil
}

// NewWrappersV4Transactor creates a new write-only instance of WrappersV4, bound to a specific deployed contract.
func NewWrappersV4Transactor(address common.Address, transactor bind.ContractTransactor) (*WrappersV4Transactor, error) {
	contract, err := bindWrappersV4(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WrappersV4Transactor{contract: contract}, nil
}

// NewWrappersV4Filterer creates a new log filterer instance of WrappersV4, bound to a specific deployed contract.
func NewWrappersV4Filterer(address common.Address, filterer bind.ContractFilterer) (*WrappersV4Filterer, error) {
	contract, err := bindWrappersV4(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WrappersV4Filterer{contract: contract}, nil
}

// bindWrappersV4 binds a generic wrapper to an already deployed contract.
func bindWrappersV4(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(WrappersV4ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WrappersV4 *WrappersV4Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _WrappersV4.Contract.WrappersV4Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WrappersV4 *WrappersV4Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WrappersV4.Contract.WrappersV4Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WrappersV4 *WrappersV4Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WrappersV4.Contract.WrappersV4Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WrappersV4 *WrappersV4CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _WrappersV4.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WrappersV4 *WrappersV4TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WrappersV4.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WrappersV4 *WrappersV4TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WrappersV4.Contract.contract.Transact(opts, method, params...)
}

// EIP712EXCHANGEDOMAINHASH is a free data retrieval call binding the contract method 0xc26cfecd.
//
// Solidity: function EIP712_EXCHANGE_DOMAIN_HASH() view returns(bytes32)
func (_WrappersV4 *WrappersV4Caller) EIP712EXCHANGEDOMAINHASH(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "EIP712_EXCHANGE_DOMAIN_HASH")
	return *ret0, err
}

// EIP712EXCHANGEDOMAINHASH is a free data retrieval call binding the contract method 0xc26cfecd.
//
// Solidity: function EIP712_EXCHANGE_DOMAIN_HASH() view returns(bytes32)
func (_WrappersV4 *WrappersV4Session) EIP712EXCHANGEDOMAINHASH() ([32]byte, error) {
	return _WrappersV4.Contract.EIP712EXCHANGEDOMAINHASH(&_WrappersV4.CallOpts)
}

// EIP712EXCHANGEDOMAINHASH is a free data retrieval call binding the contract method 0xc26cfecd.
//
// Solidity: function EIP712_EXCHANGE_DOMAIN_HASH() view returns(bytes32)
func (_WrappersV4 *WrappersV4CallerSession) EIP712EXCHANGEDOMAINHASH() ([32]byte, error) {
	return _WrappersV4.Contract.EIP712EXCHANGEDOMAINHASH(&_WrappersV4.CallOpts)
}

// ChaiBridgeAddress is a free data retrieval call binding the contract method 0xc82037ef.
//
// Solidity: function chaiBridgeAddress() view returns(address)
func (_WrappersV4 *WrappersV4Caller) ChaiBridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "chaiBridgeAddress")
	return *ret0, err
}

// ChaiBridgeAddress is a free data retrieval call binding the contract method 0xc82037ef.
//
// Solidity: function chaiBridgeAddress() view returns(address)
func (_WrappersV4 *WrappersV4Session) ChaiBridgeAddress() (common.Address, error) {
	return _WrappersV4.Contract.ChaiBridgeAddress(&_WrappersV4.CallOpts)
}

// ChaiBridgeAddress is a free data retrieval call binding the contract method 0xc82037ef.
//
// Solidity: function chaiBridgeAddress() view returns(address)
func (_WrappersV4 *WrappersV4CallerSession) ChaiBridgeAddress() (common.Address, error) {
	return _WrappersV4.Contract.ChaiBridgeAddress(&_WrappersV4.CallOpts)
}

// DecodeAssetProxyId is a free data retrieval call binding the contract method 0xd4695028.
//
// Solidity: function decodeAssetProxyId(bytes assetData) pure returns(bytes4 assetProxyId)
func (_WrappersV4 *WrappersV4Caller) DecodeAssetProxyId(opts *bind.CallOpts, assetData []byte) ([4]byte, error) {
	var (
		ret0 = new([4]byte)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "decodeAssetProxyId", assetData)
	return *ret0, err
}

// DecodeAssetProxyId is a free data retrieval call binding the contract method 0xd4695028.
//
// Solidity: function decodeAssetProxyId(bytes assetData) pure returns(bytes4 assetProxyId)
func (_WrappersV4 *WrappersV4Session) DecodeAssetProxyId(assetData []byte) ([4]byte, error) {
	return _WrappersV4.Contract.DecodeAssetProxyId(&_WrappersV4.CallOpts, assetData)
}

// DecodeAssetProxyId is a free data retrieval call binding the contract method 0xd4695028.
//
// Solidity: function decodeAssetProxyId(bytes assetData) pure returns(bytes4 assetProxyId)
func (_WrappersV4 *WrappersV4CallerSession) DecodeAssetProxyId(assetData []byte) ([4]byte, error) {
	return _WrappersV4.Contract.DecodeAssetProxyId(&_WrappersV4.CallOpts, assetData)
}

// DecodeERC1155AssetData is a free data retrieval call binding the contract method 0x9eadc835.
//
// Solidity: function decodeERC1155AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData)
func (_WrappersV4 *WrappersV4Caller) DecodeERC1155AssetData(opts *bind.CallOpts, assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
	TokenIds     []*big.Int
	TokenValues  []*big.Int
	CallbackData []byte
}, error) {
	ret := new(struct {
		AssetProxyId [4]byte
		TokenAddress common.Address
		TokenIds     []*big.Int
		TokenValues  []*big.Int
		CallbackData []byte
	})
	out := ret
	err := _WrappersV4.contract.Call(opts, out, "decodeERC1155AssetData", assetData)
	return *ret, err
}

// DecodeERC1155AssetData is a free data retrieval call binding the contract method 0x9eadc835.
//
// Solidity: function decodeERC1155AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData)
func (_WrappersV4 *WrappersV4Session) DecodeERC1155AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
	TokenIds     []*big.Int
	TokenValues  []*big.Int
	CallbackData []byte
}, error) {
	return _WrappersV4.Contract.DecodeERC1155AssetData(&_WrappersV4.CallOpts, assetData)
}

// DecodeERC1155AssetData is a free data retrieval call binding the contract method 0x9eadc835.
//
// Solidity: function decodeERC1155AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData)
func (_WrappersV4 *WrappersV4CallerSession) DecodeERC1155AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
	TokenIds     []*big.Int
	TokenValues  []*big.Int
	CallbackData []byte
}, error) {
	return _WrappersV4.Contract.DecodeERC1155AssetData(&_WrappersV4.CallOpts, assetData)
}

// DecodeERC20AssetData is a free data retrieval call binding the contract method 0x8f4ce479.
//
// Solidity: function decodeERC20AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress)
func (_WrappersV4 *WrappersV4Caller) DecodeERC20AssetData(opts *bind.CallOpts, assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
}, error) {
	ret := new(struct {
		AssetProxyId [4]byte
		TokenAddress common.Address
	})
	out := ret
	err := _WrappersV4.contract.Call(opts, out, "decodeERC20AssetData", assetData)
	return *ret, err
}

// DecodeERC20AssetData is a free data retrieval call binding the contract method 0x8f4ce479.
//
// Solidity: function decodeERC20AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress)
func (_WrappersV4 *WrappersV4Session) DecodeERC20AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
}, error) {
	return _WrappersV4.Contract.DecodeERC20AssetData(&_WrappersV4.CallOpts, assetData)
}

// DecodeERC20AssetData is a free data retrieval call binding the contract method 0x8f4ce479.
//
// Solidity: function decodeERC20AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress)
func (_WrappersV4 *WrappersV4CallerSession) DecodeERC20AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
}, error) {
	return _WrappersV4.Contract.DecodeERC20AssetData(&_WrappersV4.CallOpts, assetData)
}

// DecodeERC20BridgeAssetData is a free data retrieval call binding the contract method 0x750bdb30.
//
// Solidity: function decodeERC20BridgeAssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, address bridgeAddress, bytes bridgeData)
func (_WrappersV4 *WrappersV4Caller) DecodeERC20BridgeAssetData(opts *bind.CallOpts, assetData []byte) (struct {
	AssetProxyId  [4]byte
	TokenAddress  common.Address
	BridgeAddress common.Address
	BridgeData    []byte
}, error) {
	ret := new(struct {
		AssetProxyId  [4]byte
		TokenAddress  common.Address
		BridgeAddress common.Address
		BridgeData    []byte
	})
	out := ret
	err := _WrappersV4.contract.Call(opts, out, "decodeERC20BridgeAssetData", assetData)
	return *ret, err
}

// DecodeERC20BridgeAssetData is a free data retrieval call binding the contract method 0x750bdb30.
//
// Solidity: function decodeERC20BridgeAssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, address bridgeAddress, bytes bridgeData)
func (_WrappersV4 *WrappersV4Session) DecodeERC20BridgeAssetData(assetData []byte) (struct {
	AssetProxyId  [4]byte
	TokenAddress  common.Address
	BridgeAddress common.Address
	BridgeData    []byte
}, error) {
	return _WrappersV4.Contract.DecodeERC20BridgeAssetData(&_WrappersV4.CallOpts, assetData)
}

// DecodeERC20BridgeAssetData is a free data retrieval call binding the contract method 0x750bdb30.
//
// Solidity: function decodeERC20BridgeAssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, address bridgeAddress, bytes bridgeData)
func (_WrappersV4 *WrappersV4CallerSession) DecodeERC20BridgeAssetData(assetData []byte) (struct {
	AssetProxyId  [4]byte
	TokenAddress  common.Address
	BridgeAddress common.Address
	BridgeData    []byte
}, error) {
	return _WrappersV4.Contract.DecodeERC20BridgeAssetData(&_WrappersV4.CallOpts, assetData)
}

// DecodeERC721AssetData is a free data retrieval call binding the contract method 0x04a5618a.
//
// Solidity: function decodeERC721AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, uint256 tokenId)
func (_WrappersV4 *WrappersV4Caller) DecodeERC721AssetData(opts *bind.CallOpts, assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
	TokenId      *big.Int
}, error) {
	ret := new(struct {
		AssetProxyId [4]byte
		TokenAddress common.Address
		TokenId      *big.Int
	})
	out := ret
	err := _WrappersV4.contract.Call(opts, out, "decodeERC721AssetData", assetData)
	return *ret, err
}

// DecodeERC721AssetData is a free data retrieval call binding the contract method 0x04a5618a.
//
// Solidity: function decodeERC721AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, uint256 tokenId)
func (_WrappersV4 *WrappersV4Session) DecodeERC721AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
	TokenId      *big.Int
}, error) {
	return _WrappersV4.Contract.DecodeERC721AssetData(&_WrappersV4.CallOpts, assetData)
}

// DecodeERC721AssetData is a free data retrieval call binding the contract method 0x04a5618a.
//
// Solidity: function decodeERC721AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, uint256 tokenId)
func (_WrappersV4 *WrappersV4CallerSession) DecodeERC721AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
	TokenId      *big.Int
}, error) {
	return _WrappersV4.Contract.DecodeERC721AssetData(&_WrappersV4.CallOpts, assetData)
}

// DecodeMultiAssetData is a free data retrieval call binding the contract method 0xbbb2dcf6.
//
// Solidity: function decodeMultiAssetData(bytes assetData) pure returns(bytes4 assetProxyId, uint256[] amounts, bytes[] nestedAssetData)
func (_WrappersV4 *WrappersV4Caller) DecodeMultiAssetData(opts *bind.CallOpts, assetData []byte) (struct {
	AssetProxyId    [4]byte
	Amounts         []*big.Int
	NestedAssetData [][]byte
}, error) {
	ret := new(struct {
		AssetProxyId    [4]byte
		Amounts         []*big.Int
		NestedAssetData [][]byte
	})
	out := ret
	err := _WrappersV4.contract.Call(opts, out, "decodeMultiAssetData", assetData)
	return *ret, err
}

// DecodeMultiAssetData is a free data retrieval call binding the contract method 0xbbb2dcf6.
//
// Solidity: function decodeMultiAssetData(bytes assetData) pure returns(bytes4 assetProxyId, uint256[] amounts, bytes[] nestedAssetData)
func (_WrappersV4 *WrappersV4Session) DecodeMultiAssetData(assetData []byte) (struct {
	AssetProxyId    [4]byte
	Amounts         []*big.Int
	NestedAssetData [][]byte
}, error) {
	return _WrappersV4.Contract.DecodeMultiAssetData(&_WrappersV4.CallOpts, assetData)
}

// DecodeMultiAssetData is a free data retrieval call binding the contract method 0xbbb2dcf6.
//
// Solidity: function decodeMultiAssetData(bytes assetData) pure returns(bytes4 assetProxyId, uint256[] amounts, bytes[] nestedAssetData)
func (_WrappersV4 *WrappersV4CallerSession) DecodeMultiAssetData(assetData []byte) (struct {
	AssetProxyId    [4]byte
	Amounts         []*big.Int
	NestedAssetData [][]byte
}, error) {
	return _WrappersV4.Contract.DecodeMultiAssetData(&_WrappersV4.CallOpts, assetData)
}

// DecodeStaticCallAssetData is a free data retrieval call binding the contract method 0xca49f47c.
//
// Solidity: function decodeStaticCallAssetData(bytes assetData) pure returns(bytes4 assetProxyId, address staticCallTargetAddress, bytes staticCallData, bytes32 expectedReturnDataHash)
func (_WrappersV4 *WrappersV4Caller) DecodeStaticCallAssetData(opts *bind.CallOpts, assetData []byte) (struct {
	AssetProxyId            [4]byte
	StaticCallTargetAddress common.Address
	StaticCallData          []byte
	ExpectedReturnDataHash  [32]byte
}, error) {
	ret := new(struct {
		AssetProxyId            [4]byte
		StaticCallTargetAddress common.Address
		StaticCallData          []byte
		ExpectedReturnDataHash  [32]byte
	})
	out := ret
	err := _WrappersV4.contract.Call(opts, out, "decodeStaticCallAssetData", assetData)
	return *ret, err
}

// DecodeStaticCallAssetData is a free data retrieval call binding the contract method 0xca49f47c.
//
// Solidity: function decodeStaticCallAssetData(bytes assetData) pure returns(bytes4 assetProxyId, address staticCallTargetAddress, bytes staticCallData, bytes32 expectedReturnDataHash)
func (_WrappersV4 *WrappersV4Session) DecodeStaticCallAssetData(assetData []byte) (struct {
	AssetProxyId            [4]byte
	StaticCallTargetAddress common.Address
	StaticCallData          []byte
	ExpectedReturnDataHash  [32]byte
}, error) {
	return _WrappersV4.Contract.DecodeStaticCallAssetData(&_WrappersV4.CallOpts, assetData)
}

// DecodeStaticCallAssetData is a free data retrieval call binding the contract method 0xca49f47c.
//
// Solidity: function decodeStaticCallAssetData(bytes assetData) pure returns(bytes4 assetProxyId, address staticCallTargetAddress, bytes staticCallData, bytes32 expectedReturnDataHash)
func (_WrappersV4 *WrappersV4CallerSession) DecodeStaticCallAssetData(assetData []byte) (struct {
	AssetProxyId            [4]byte
	StaticCallTargetAddress common.Address
	StaticCallData          []byte
	ExpectedReturnDataHash  [32]byte
}, error) {
	return _WrappersV4.Contract.DecodeStaticCallAssetData(&_WrappersV4.CallOpts, assetData)
}

// DecodeZeroExTransactionData is a free data retrieval call binding the contract method 0x6f83188e.
//
// Solidity: function decodeZeroExTransactionData(bytes transactionData) pure returns(string functionName, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures)
func (_WrappersV4 *WrappersV4Caller) DecodeZeroExTransactionData(opts *bind.CallOpts, transactionData []byte) (struct {
	FunctionName          string
	Orders                []LibOrderOrder
	TakerAssetFillAmounts []*big.Int
	Signatures            [][]byte
}, error) {
	ret := new(struct {
		FunctionName          string
		Orders                []LibOrderOrder
		TakerAssetFillAmounts []*big.Int
		Signatures            [][]byte
	})
	out := ret
	err := _WrappersV4.contract.Call(opts, out, "decodeZeroExTransactionData", transactionData)
	return *ret, err
}

// DecodeZeroExTransactionData is a free data retrieval call binding the contract method 0x6f83188e.
//
// Solidity: function decodeZeroExTransactionData(bytes transactionData) pure returns(string functionName, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures)
func (_WrappersV4 *WrappersV4Session) DecodeZeroExTransactionData(transactionData []byte) (struct {
	FunctionName          string
	Orders                []LibOrderOrder
	TakerAssetFillAmounts []*big.Int
	Signatures            [][]byte
}, error) {
	return _WrappersV4.Contract.DecodeZeroExTransactionData(&_WrappersV4.CallOpts, transactionData)
}

// DecodeZeroExTransactionData is a free data retrieval call binding the contract method 0x6f83188e.
//
// Solidity: function decodeZeroExTransactionData(bytes transactionData) pure returns(string functionName, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures)
func (_WrappersV4 *WrappersV4CallerSession) DecodeZeroExTransactionData(transactionData []byte) (struct {
	FunctionName          string
	Orders                []LibOrderOrder
	TakerAssetFillAmounts []*big.Int
	Signatures            [][]byte
}, error) {
	return _WrappersV4.Contract.DecodeZeroExTransactionData(&_WrappersV4.CallOpts, transactionData)
}

// DydxBridgeAddress is a free data retrieval call binding the contract method 0xa7530f12.
//
// Solidity: function dydxBridgeAddress() view returns(address)
func (_WrappersV4 *WrappersV4Caller) DydxBridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "dydxBridgeAddress")
	return *ret0, err
}

// DydxBridgeAddress is a free data retrieval call binding the contract method 0xa7530f12.
//
// Solidity: function dydxBridgeAddress() view returns(address)
func (_WrappersV4 *WrappersV4Session) DydxBridgeAddress() (common.Address, error) {
	return _WrappersV4.Contract.DydxBridgeAddress(&_WrappersV4.CallOpts)
}

// DydxBridgeAddress is a free data retrieval call binding the contract method 0xa7530f12.
//
// Solidity: function dydxBridgeAddress() view returns(address)
func (_WrappersV4 *WrappersV4CallerSession) DydxBridgeAddress() (common.Address, error) {
	return _WrappersV4.Contract.DydxBridgeAddress(&_WrappersV4.CallOpts)
}

// EncodeERC1155AssetData is a free data retrieval call binding the contract method 0xb43cffe1.
//
// Solidity: function encodeERC1155AssetData(address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData) pure returns(bytes assetData)
func (_WrappersV4 *WrappersV4Caller) EncodeERC1155AssetData(opts *bind.CallOpts, tokenAddress common.Address, tokenIds []*big.Int, tokenValues []*big.Int, callbackData []byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "encodeERC1155AssetData", tokenAddress, tokenIds, tokenValues, callbackData)
	return *ret0, err
}

// EncodeERC1155AssetData is a free data retrieval call binding the contract method 0xb43cffe1.
//
// Solidity: function encodeERC1155AssetData(address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData) pure returns(bytes assetData)
func (_WrappersV4 *WrappersV4Session) EncodeERC1155AssetData(tokenAddress common.Address, tokenIds []*big.Int, tokenValues []*big.Int, callbackData []byte) ([]byte, error) {
	return _WrappersV4.Contract.EncodeERC1155AssetData(&_WrappersV4.CallOpts, tokenAddress, tokenIds, tokenValues, callbackData)
}

// EncodeERC1155AssetData is a free data retrieval call binding the contract method 0xb43cffe1.
//
// Solidity: function encodeERC1155AssetData(address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData) pure returns(bytes assetData)
func (_WrappersV4 *WrappersV4CallerSession) EncodeERC1155AssetData(tokenAddress common.Address, tokenIds []*big.Int, tokenValues []*big.Int, callbackData []byte) ([]byte, error) {
	return _WrappersV4.Contract.EncodeERC1155AssetData(&_WrappersV4.CallOpts, tokenAddress, tokenIds, tokenValues, callbackData)
}

// EncodeERC20AssetData is a free data retrieval call binding the contract method 0x590aa875.
//
// Solidity: function encodeERC20AssetData(address tokenAddress) pure returns(bytes assetData)
func (_WrappersV4 *WrappersV4Caller) EncodeERC20AssetData(opts *bind.CallOpts, tokenAddress common.Address) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "encodeERC20AssetData", tokenAddress)
	return *ret0, err
}

// EncodeERC20AssetData is a free data retrieval call binding the contract method 0x590aa875.
//
// Solidity: function encodeERC20AssetData(address tokenAddress) pure returns(bytes assetData)
func (_WrappersV4 *WrappersV4Session) EncodeERC20AssetData(tokenAddress common.Address) ([]byte, error) {
	return _WrappersV4.Contract.EncodeERC20AssetData(&_WrappersV4.CallOpts, tokenAddress)
}

// EncodeERC20AssetData is a free data retrieval call binding the contract method 0x590aa875.
//
// Solidity: function encodeERC20AssetData(address tokenAddress) pure returns(bytes assetData)
func (_WrappersV4 *WrappersV4CallerSession) EncodeERC20AssetData(tokenAddress common.Address) ([]byte, error) {
	return _WrappersV4.Contract.EncodeERC20AssetData(&_WrappersV4.CallOpts, tokenAddress)
}

// EncodeERC721AssetData is a free data retrieval call binding the contract method 0xa6627e9f.
//
// Solidity: function encodeERC721AssetData(address tokenAddress, uint256 tokenId) pure returns(bytes assetData)
func (_WrappersV4 *WrappersV4Caller) EncodeERC721AssetData(opts *bind.CallOpts, tokenAddress common.Address, tokenId *big.Int) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "encodeERC721AssetData", tokenAddress, tokenId)
	return *ret0, err
}

// EncodeERC721AssetData is a free data retrieval call binding the contract method 0xa6627e9f.
//
// Solidity: function encodeERC721AssetData(address tokenAddress, uint256 tokenId) pure returns(bytes assetData)
func (_WrappersV4 *WrappersV4Session) EncodeERC721AssetData(tokenAddress common.Address, tokenId *big.Int) ([]byte, error) {
	return _WrappersV4.Contract.EncodeERC721AssetData(&_WrappersV4.CallOpts, tokenAddress, tokenId)
}

// EncodeERC721AssetData is a free data retrieval call binding the contract method 0xa6627e9f.
//
// Solidity: function encodeERC721AssetData(address tokenAddress, uint256 tokenId) pure returns(bytes assetData)
func (_WrappersV4 *WrappersV4CallerSession) EncodeERC721AssetData(tokenAddress common.Address, tokenId *big.Int) ([]byte, error) {
	return _WrappersV4.Contract.EncodeERC721AssetData(&_WrappersV4.CallOpts, tokenAddress, tokenId)
}

// EncodeMultiAssetData is a free data retrieval call binding the contract method 0xd3d862d1.
//
// Solidity: function encodeMultiAssetData(uint256[] amounts, bytes[] nestedAssetData) pure returns(bytes assetData)
func (_WrappersV4 *WrappersV4Caller) EncodeMultiAssetData(opts *bind.CallOpts, amounts []*big.Int, nestedAssetData [][]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "encodeMultiAssetData", amounts, nestedAssetData)
	return *ret0, err
}

// EncodeMultiAssetData is a free data retrieval call binding the contract method 0xd3d862d1.
//
// Solidity: function encodeMultiAssetData(uint256[] amounts, bytes[] nestedAssetData) pure returns(bytes assetData)
func (_WrappersV4 *WrappersV4Session) EncodeMultiAssetData(amounts []*big.Int, nestedAssetData [][]byte) ([]byte, error) {
	return _WrappersV4.Contract.EncodeMultiAssetData(&_WrappersV4.CallOpts, amounts, nestedAssetData)
}

// EncodeMultiAssetData is a free data retrieval call binding the contract method 0xd3d862d1.
//
// Solidity: function encodeMultiAssetData(uint256[] amounts, bytes[] nestedAssetData) pure returns(bytes assetData)
func (_WrappersV4 *WrappersV4CallerSession) EncodeMultiAssetData(amounts []*big.Int, nestedAssetData [][]byte) ([]byte, error) {
	return _WrappersV4.Contract.EncodeMultiAssetData(&_WrappersV4.CallOpts, amounts, nestedAssetData)
}

// EncodeStaticCallAssetData is a free data retrieval call binding the contract method 0x63eb3992.
//
// Solidity: function encodeStaticCallAssetData(address staticCallTargetAddress, bytes staticCallData, bytes32 expectedReturnDataHash) pure returns(bytes assetData)
func (_WrappersV4 *WrappersV4Caller) EncodeStaticCallAssetData(opts *bind.CallOpts, staticCallTargetAddress common.Address, staticCallData []byte, expectedReturnDataHash [32]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "encodeStaticCallAssetData", staticCallTargetAddress, staticCallData, expectedReturnDataHash)
	return *ret0, err
}

// EncodeStaticCallAssetData is a free data retrieval call binding the contract method 0x63eb3992.
//
// Solidity: function encodeStaticCallAssetData(address staticCallTargetAddress, bytes staticCallData, bytes32 expectedReturnDataHash) pure returns(bytes assetData)
func (_WrappersV4 *WrappersV4Session) EncodeStaticCallAssetData(staticCallTargetAddress common.Address, staticCallData []byte, expectedReturnDataHash [32]byte) ([]byte, error) {
	return _WrappersV4.Contract.EncodeStaticCallAssetData(&_WrappersV4.CallOpts, staticCallTargetAddress, staticCallData, expectedReturnDataHash)
}

// EncodeStaticCallAssetData is a free data retrieval call binding the contract method 0x63eb3992.
//
// Solidity: function encodeStaticCallAssetData(address staticCallTargetAddress, bytes staticCallData, bytes32 expectedReturnDataHash) pure returns(bytes assetData)
func (_WrappersV4 *WrappersV4CallerSession) EncodeStaticCallAssetData(staticCallTargetAddress common.Address, staticCallData []byte, expectedReturnDataHash [32]byte) ([]byte, error) {
	return _WrappersV4.Contract.EncodeStaticCallAssetData(&_WrappersV4.CallOpts, staticCallTargetAddress, staticCallData, expectedReturnDataHash)
}

// Erc1155ProxyAddress is a free data retrieval call binding the contract method 0xff84e7cc.
//
// Solidity: function erc1155ProxyAddress() view returns(address)
func (_WrappersV4 *WrappersV4Caller) Erc1155ProxyAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "erc1155ProxyAddress")
	return *ret0, err
}

// Erc1155ProxyAddress is a free data retrieval call binding the contract method 0xff84e7cc.
//
// Solidity: function erc1155ProxyAddress() view returns(address)
func (_WrappersV4 *WrappersV4Session) Erc1155ProxyAddress() (common.Address, error) {
	return _WrappersV4.Contract.Erc1155ProxyAddress(&_WrappersV4.CallOpts)
}

// Erc1155ProxyAddress is a free data retrieval call binding the contract method 0xff84e7cc.
//
// Solidity: function erc1155ProxyAddress() view returns(address)
func (_WrappersV4 *WrappersV4CallerSession) Erc1155ProxyAddress() (common.Address, error) {
	return _WrappersV4.Contract.Erc1155ProxyAddress(&_WrappersV4.CallOpts)
}

// Erc20ProxyAddress is a free data retrieval call binding the contract method 0xee185997.
//
// Solidity: function erc20ProxyAddress() view returns(address)
func (_WrappersV4 *WrappersV4Caller) Erc20ProxyAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "erc20ProxyAddress")
	return *ret0, err
}

// Erc20ProxyAddress is a free data retrieval call binding the contract method 0xee185997.
//
// Solidity: function erc20ProxyAddress() view returns(address)
func (_WrappersV4 *WrappersV4Session) Erc20ProxyAddress() (common.Address, error) {
	return _WrappersV4.Contract.Erc20ProxyAddress(&_WrappersV4.CallOpts)
}

// Erc20ProxyAddress is a free data retrieval call binding the contract method 0xee185997.
//
// Solidity: function erc20ProxyAddress() view returns(address)
func (_WrappersV4 *WrappersV4CallerSession) Erc20ProxyAddress() (common.Address, error) {
	return _WrappersV4.Contract.Erc20ProxyAddress(&_WrappersV4.CallOpts)
}

// Erc721ProxyAddress is a free data retrieval call binding the contract method 0xef3bb097.
//
// Solidity: function erc721ProxyAddress() view returns(address)
func (_WrappersV4 *WrappersV4Caller) Erc721ProxyAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "erc721ProxyAddress")
	return *ret0, err
}

// Erc721ProxyAddress is a free data retrieval call binding the contract method 0xef3bb097.
//
// Solidity: function erc721ProxyAddress() view returns(address)
func (_WrappersV4 *WrappersV4Session) Erc721ProxyAddress() (common.Address, error) {
	return _WrappersV4.Contract.Erc721ProxyAddress(&_WrappersV4.CallOpts)
}

// Erc721ProxyAddress is a free data retrieval call binding the contract method 0xef3bb097.
//
// Solidity: function erc721ProxyAddress() view returns(address)
func (_WrappersV4 *WrappersV4CallerSession) Erc721ProxyAddress() (common.Address, error) {
	return _WrappersV4.Contract.Erc721ProxyAddress(&_WrappersV4.CallOpts)
}

// ExchangeAddress is a free data retrieval call binding the contract method 0x9cd01605.
//
// Solidity: function exchangeAddress() view returns(address)
func (_WrappersV4 *WrappersV4Caller) ExchangeAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "exchangeAddress")
	return *ret0, err
}

// ExchangeAddress is a free data retrieval call binding the contract method 0x9cd01605.
//
// Solidity: function exchangeAddress() view returns(address)
func (_WrappersV4 *WrappersV4Session) ExchangeAddress() (common.Address, error) {
	return _WrappersV4.Contract.ExchangeAddress(&_WrappersV4.CallOpts)
}

// ExchangeAddress is a free data retrieval call binding the contract method 0x9cd01605.
//
// Solidity: function exchangeAddress() view returns(address)
func (_WrappersV4 *WrappersV4CallerSession) ExchangeAddress() (common.Address, error) {
	return _WrappersV4.Contract.ExchangeAddress(&_WrappersV4.CallOpts)
}

// GetEthBalances is a free data retrieval call binding the contract method 0xa0901e51.
//
// Solidity: function getEthBalances(address[] addresses) view returns(uint256[])
func (_WrappersV4 *WrappersV4Caller) GetEthBalances(opts *bind.CallOpts, addresses []common.Address) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "getEthBalances", addresses)
	return *ret0, err
}

// GetEthBalances is a free data retrieval call binding the contract method 0xa0901e51.
//
// Solidity: function getEthBalances(address[] addresses) view returns(uint256[])
func (_WrappersV4 *WrappersV4Session) GetEthBalances(addresses []common.Address) ([]*big.Int, error) {
	return _WrappersV4.Contract.GetEthBalances(&_WrappersV4.CallOpts, addresses)
}

// GetEthBalances is a free data retrieval call binding the contract method 0xa0901e51.
//
// Solidity: function getEthBalances(address[] addresses) view returns(uint256[])
func (_WrappersV4 *WrappersV4CallerSession) GetEthBalances(addresses []common.Address) ([]*big.Int, error) {
	return _WrappersV4.Contract.GetEthBalances(&_WrappersV4.CallOpts, addresses)
}

// GetOrderHash is a free data retrieval call binding the contract method 0xa070cac8.
//
// Solidity: function getOrderHash((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 chainId, address exchange) pure returns(bytes32 orderHash)
func (_WrappersV4 *WrappersV4Caller) GetOrderHash(opts *bind.CallOpts, order LibOrderOrder, chainId *big.Int, exchange common.Address) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "getOrderHash", order, chainId, exchange)
	return *ret0, err
}

// GetOrderHash is a free data retrieval call binding the contract method 0xa070cac8.
//
// Solidity: function getOrderHash((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 chainId, address exchange) pure returns(bytes32 orderHash)
func (_WrappersV4 *WrappersV4Session) GetOrderHash(order LibOrderOrder, chainId *big.Int, exchange common.Address) ([32]byte, error) {
	return _WrappersV4.Contract.GetOrderHash(&_WrappersV4.CallOpts, order, chainId, exchange)
}

// GetOrderHash is a free data retrieval call binding the contract method 0xa070cac8.
//
// Solidity: function getOrderHash((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 chainId, address exchange) pure returns(bytes32 orderHash)
func (_WrappersV4 *WrappersV4CallerSession) GetOrderHash(order LibOrderOrder, chainId *big.Int, exchange common.Address) ([32]byte, error) {
	return _WrappersV4.Contract.GetOrderHash(&_WrappersV4.CallOpts, order, chainId, exchange)
}

// GetTransactionHash is a free data retrieval call binding the contract method 0x79c9c426.
//
// Solidity: function getTransactionHash((uint256,uint256,uint256,address,bytes) transaction, uint256 chainId, address exchange) pure returns(bytes32 transactionHash)
func (_WrappersV4 *WrappersV4Caller) GetTransactionHash(opts *bind.CallOpts, transaction LibZeroExTransactionZeroExTransaction, chainId *big.Int, exchange common.Address) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "getTransactionHash", transaction, chainId, exchange)
	return *ret0, err
}

// GetTransactionHash is a free data retrieval call binding the contract method 0x79c9c426.
//
// Solidity: function getTransactionHash((uint256,uint256,uint256,address,bytes) transaction, uint256 chainId, address exchange) pure returns(bytes32 transactionHash)
func (_WrappersV4 *WrappersV4Session) GetTransactionHash(transaction LibZeroExTransactionZeroExTransaction, chainId *big.Int, exchange common.Address) ([32]byte, error) {
	return _WrappersV4.Contract.GetTransactionHash(&_WrappersV4.CallOpts, transaction, chainId, exchange)
}

// GetTransactionHash is a free data retrieval call binding the contract method 0x79c9c426.
//
// Solidity: function getTransactionHash((uint256,uint256,uint256,address,bytes) transaction, uint256 chainId, address exchange) pure returns(bytes32 transactionHash)
func (_WrappersV4 *WrappersV4CallerSession) GetTransactionHash(transaction LibZeroExTransactionZeroExTransaction, chainId *big.Int, exchange common.Address) ([32]byte, error) {
	return _WrappersV4.Contract.GetTransactionHash(&_WrappersV4.CallOpts, transaction, chainId, exchange)
}

// RevertIfInvalidAssetData is a free data retrieval call binding the contract method 0x46eb65cb.
//
// Solidity: function revertIfInvalidAssetData(bytes assetData) pure returns()
func (_WrappersV4 *WrappersV4Caller) RevertIfInvalidAssetData(opts *bind.CallOpts, assetData []byte) error {
	var ()
	out := &[]interface{}{}
	err := _WrappersV4.contract.Call(opts, out, "revertIfInvalidAssetData", assetData)
	return err
}

// RevertIfInvalidAssetData is a free data retrieval call binding the contract method 0x46eb65cb.
//
// Solidity: function revertIfInvalidAssetData(bytes assetData) pure returns()
func (_WrappersV4 *WrappersV4Session) RevertIfInvalidAssetData(assetData []byte) error {
	return _WrappersV4.Contract.RevertIfInvalidAssetData(&_WrappersV4.CallOpts, assetData)
}

// RevertIfInvalidAssetData is a free data retrieval call binding the contract method 0x46eb65cb.
//
// Solidity: function revertIfInvalidAssetData(bytes assetData) pure returns()
func (_WrappersV4 *WrappersV4CallerSession) RevertIfInvalidAssetData(assetData []byte) error {
	return _WrappersV4.Contract.RevertIfInvalidAssetData(&_WrappersV4.CallOpts, assetData)
}

// StaticCallProxyAddress is a free data retrieval call binding the contract method 0x9baf2705.
//
// Solidity: function staticCallProxyAddress() view returns(address)
func (_WrappersV4 *WrappersV4Caller) StaticCallProxyAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "staticCallProxyAddress")
	return *ret0, err
}

// StaticCallProxyAddress is a free data retrieval call binding the contract method 0x9baf2705.
//
// Solidity: function staticCallProxyAddress() view returns(address)
func (_WrappersV4 *WrappersV4Session) StaticCallProxyAddress() (common.Address, error) {
	return _WrappersV4.Contract.StaticCallProxyAddress(&_WrappersV4.CallOpts)
}

// StaticCallProxyAddress is a free data retrieval call binding the contract method 0x9baf2705.
//
// Solidity: function staticCallProxyAddress() view returns(address)
func (_WrappersV4 *WrappersV4CallerSession) StaticCallProxyAddress() (common.Address, error) {
	return _WrappersV4.Contract.StaticCallProxyAddress(&_WrappersV4.CallOpts)
}

// GetAssetProxyAllowance is a paid mutator transaction binding the contract method 0xd186037f.
//
// Solidity: function getAssetProxyAllowance(address ownerAddress, bytes assetData) returns(uint256 allowance)
func (_WrappersV4 *WrappersV4Transactor) GetAssetProxyAllowance(opts *bind.TransactOpts, ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "getAssetProxyAllowance", ownerAddress, assetData)
}

// GetAssetProxyAllowance is a paid mutator transaction binding the contract method 0xd186037f.
//
// Solidity: function getAssetProxyAllowance(address ownerAddress, bytes assetData) returns(uint256 allowance)
func (_WrappersV4 *WrappersV4Session) GetAssetProxyAllowance(ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetAssetProxyAllowance(&_WrappersV4.TransactOpts, ownerAddress, assetData)
}

// GetAssetProxyAllowance is a paid mutator transaction binding the contract method 0xd186037f.
//
// Solidity: function getAssetProxyAllowance(address ownerAddress, bytes assetData) returns(uint256 allowance)
func (_WrappersV4 *WrappersV4TransactorSession) GetAssetProxyAllowance(ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetAssetProxyAllowance(&_WrappersV4.TransactOpts, ownerAddress, assetData)
}

// GetBalance is a paid mutator transaction binding the contract method 0x7d727512.
//
// Solidity: function getBalance(address ownerAddress, bytes assetData) returns(uint256 balance)
func (_WrappersV4 *WrappersV4Transactor) GetBalance(opts *bind.TransactOpts, ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "getBalance", ownerAddress, assetData)
}

// GetBalance is a paid mutator transaction binding the contract method 0x7d727512.
//
// Solidity: function getBalance(address ownerAddress, bytes assetData) returns(uint256 balance)
func (_WrappersV4 *WrappersV4Session) GetBalance(ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetBalance(&_WrappersV4.TransactOpts, ownerAddress, assetData)
}

// GetBalance is a paid mutator transaction binding the contract method 0x7d727512.
//
// Solidity: function getBalance(address ownerAddress, bytes assetData) returns(uint256 balance)
func (_WrappersV4 *WrappersV4TransactorSession) GetBalance(ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetBalance(&_WrappersV4.TransactOpts, ownerAddress, assetData)
}

// GetBalanceAndAssetProxyAllowance is a paid mutator transaction binding the contract method 0x0d7b7d76.
//
// Solidity: function getBalanceAndAssetProxyAllowance(address ownerAddress, bytes assetData) returns(uint256 balance, uint256 allowance)
func (_WrappersV4 *WrappersV4Transactor) GetBalanceAndAssetProxyAllowance(opts *bind.TransactOpts, ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "getBalanceAndAssetProxyAllowance", ownerAddress, assetData)
}

// GetBalanceAndAssetProxyAllowance is a paid mutator transaction binding the contract method 0x0d7b7d76.
//
// Solidity: function getBalanceAndAssetProxyAllowance(address ownerAddress, bytes assetData) returns(uint256 balance, uint256 allowance)
func (_WrappersV4 *WrappersV4Session) GetBalanceAndAssetProxyAllowance(ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetBalanceAndAssetProxyAllowance(&_WrappersV4.TransactOpts, ownerAddress, assetData)
}

// GetBalanceAndAssetProxyAllowance is a paid mutator transaction binding the contract method 0x0d7b7d76.
//
// Solidity: function getBalanceAndAssetProxyAllowance(address ownerAddress, bytes assetData) returns(uint256 balance, uint256 allowance)
func (_WrappersV4 *WrappersV4TransactorSession) GetBalanceAndAssetProxyAllowance(ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetBalanceAndAssetProxyAllowance(&_WrappersV4.TransactOpts, ownerAddress, assetData)
}

// GetBatchAssetProxyAllowances is a paid mutator transaction binding the contract method 0x4dfdac20.
//
// Solidity: function getBatchAssetProxyAllowances(address ownerAddress, bytes[] assetData) returns(uint256[] allowances)
func (_WrappersV4 *WrappersV4Transactor) GetBatchAssetProxyAllowances(opts *bind.TransactOpts, ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "getBatchAssetProxyAllowances", ownerAddress, assetData)
}

// GetBatchAssetProxyAllowances is a paid mutator transaction binding the contract method 0x4dfdac20.
//
// Solidity: function getBatchAssetProxyAllowances(address ownerAddress, bytes[] assetData) returns(uint256[] allowances)
func (_WrappersV4 *WrappersV4Session) GetBatchAssetProxyAllowances(ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetBatchAssetProxyAllowances(&_WrappersV4.TransactOpts, ownerAddress, assetData)
}

// GetBatchAssetProxyAllowances is a paid mutator transaction binding the contract method 0x4dfdac20.
//
// Solidity: function getBatchAssetProxyAllowances(address ownerAddress, bytes[] assetData) returns(uint256[] allowances)
func (_WrappersV4 *WrappersV4TransactorSession) GetBatchAssetProxyAllowances(ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetBatchAssetProxyAllowances(&_WrappersV4.TransactOpts, ownerAddress, assetData)
}

// GetBatchBalances is a paid mutator transaction binding the contract method 0xd001c5dc.
//
// Solidity: function getBatchBalances(address ownerAddress, bytes[] assetData) returns(uint256[] balances)
func (_WrappersV4 *WrappersV4Transactor) GetBatchBalances(opts *bind.TransactOpts, ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "getBatchBalances", ownerAddress, assetData)
}

// GetBatchBalances is a paid mutator transaction binding the contract method 0xd001c5dc.
//
// Solidity: function getBatchBalances(address ownerAddress, bytes[] assetData) returns(uint256[] balances)
func (_WrappersV4 *WrappersV4Session) GetBatchBalances(ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetBatchBalances(&_WrappersV4.TransactOpts, ownerAddress, assetData)
}

// GetBatchBalances is a paid mutator transaction binding the contract method 0xd001c5dc.
//
// Solidity: function getBatchBalances(address ownerAddress, bytes[] assetData) returns(uint256[] balances)
func (_WrappersV4 *WrappersV4TransactorSession) GetBatchBalances(ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetBatchBalances(&_WrappersV4.TransactOpts, ownerAddress, assetData)
}

// GetBatchBalancesAndAssetProxyAllowances is a paid mutator transaction binding the contract method 0xe4e6e7da.
//
// Solidity: function getBatchBalancesAndAssetProxyAllowances(address ownerAddress, bytes[] assetData) returns(uint256[] balances, uint256[] allowances)
func (_WrappersV4 *WrappersV4Transactor) GetBatchBalancesAndAssetProxyAllowances(opts *bind.TransactOpts, ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "getBatchBalancesAndAssetProxyAllowances", ownerAddress, assetData)
}

// GetBatchBalancesAndAssetProxyAllowances is a paid mutator transaction binding the contract method 0xe4e6e7da.
//
// Solidity: function getBatchBalancesAndAssetProxyAllowances(address ownerAddress, bytes[] assetData) returns(uint256[] balances, uint256[] allowances)
func (_WrappersV4 *WrappersV4Session) GetBatchBalancesAndAssetProxyAllowances(ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetBatchBalancesAndAssetProxyAllowances(&_WrappersV4.TransactOpts, ownerAddress, assetData)
}

// GetBatchBalancesAndAssetProxyAllowances is a paid mutator transaction binding the contract method 0xe4e6e7da.
//
// Solidity: function getBatchBalancesAndAssetProxyAllowances(address ownerAddress, bytes[] assetData) returns(uint256[] balances, uint256[] allowances)
func (_WrappersV4 *WrappersV4TransactorSession) GetBatchBalancesAndAssetProxyAllowances(ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetBatchBalancesAndAssetProxyAllowances(&_WrappersV4.TransactOpts, ownerAddress, assetData)
}

// GetOrderRelevantState is a paid mutator transaction binding the contract method 0xe77286eb.
//
// Solidity: function getOrderRelevantState((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, bytes signature) returns((uint8,bytes32,uint256) orderInfo, uint256 fillableTakerAssetAmount, bool isValidSignature)
func (_WrappersV4 *WrappersV4Transactor) GetOrderRelevantState(opts *bind.TransactOpts, order LibOrderOrder, signature []byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "getOrderRelevantState", order, signature)
}

// GetOrderRelevantState is a paid mutator transaction binding the contract method 0xe77286eb.
//
// Solidity: function getOrderRelevantState((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, bytes signature) returns((uint8,bytes32,uint256) orderInfo, uint256 fillableTakerAssetAmount, bool isValidSignature)
func (_WrappersV4 *WrappersV4Session) GetOrderRelevantState(order LibOrderOrder, signature []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetOrderRelevantState(&_WrappersV4.TransactOpts, order, signature)
}

// GetOrderRelevantState is a paid mutator transaction binding the contract method 0xe77286eb.
//
// Solidity: function getOrderRelevantState((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, bytes signature) returns((uint8,bytes32,uint256) orderInfo, uint256 fillableTakerAssetAmount, bool isValidSignature)
func (_WrappersV4 *WrappersV4TransactorSession) GetOrderRelevantState(order LibOrderOrder, signature []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetOrderRelevantState(&_WrappersV4.TransactOpts, order, signature)
}

// GetOrderRelevantStates is a paid mutator transaction binding the contract method 0xe25cabf7.
//
// Solidity: function getOrderRelevantStates((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, bytes[] signatures) returns((uint8,bytes32,uint256)[] ordersInfo, uint256[] fillableTakerAssetAmounts, bool[] isValidSignature)
func (_WrappersV4 *WrappersV4Transactor) GetOrderRelevantStates(opts *bind.TransactOpts, orders []LibOrderOrder, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "getOrderRelevantStates", orders, signatures)
}

// GetOrderRelevantStates is a paid mutator transaction binding the contract method 0xe25cabf7.
//
// Solidity: function getOrderRelevantStates((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, bytes[] signatures) returns((uint8,bytes32,uint256)[] ordersInfo, uint256[] fillableTakerAssetAmounts, bool[] isValidSignature)
func (_WrappersV4 *WrappersV4Session) GetOrderRelevantStates(orders []LibOrderOrder, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetOrderRelevantStates(&_WrappersV4.TransactOpts, orders, signatures)
}

// GetOrderRelevantStates is a paid mutator transaction binding the contract method 0xe25cabf7.
//
// Solidity: function getOrderRelevantStates((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, bytes[] signatures) returns((uint8,bytes32,uint256)[] ordersInfo, uint256[] fillableTakerAssetAmounts, bool[] isValidSignature)
func (_WrappersV4 *WrappersV4TransactorSession) GetOrderRelevantStates(orders []LibOrderOrder, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetOrderRelevantStates(&_WrappersV4.TransactOpts, orders, signatures)
}

// GetSimulatedOrderMakerTransferResults is a paid mutator transaction binding the contract method 0x7982653e.
//
// Solidity: function getSimulatedOrderMakerTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, address takerAddress, uint256 takerAssetFillAmount) returns(uint8 orderTransferResults)
func (_WrappersV4 *WrappersV4Transactor) GetSimulatedOrderMakerTransferResults(opts *bind.TransactOpts, order LibOrderOrder, takerAddress common.Address, takerAssetFillAmount *big.Int) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "getSimulatedOrderMakerTransferResults", order, takerAddress, takerAssetFillAmount)
}

// GetSimulatedOrderMakerTransferResults is a paid mutator transaction binding the contract method 0x7982653e.
//
// Solidity: function getSimulatedOrderMakerTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, address takerAddress, uint256 takerAssetFillAmount) returns(uint8 orderTransferResults)
func (_WrappersV4 *WrappersV4Session) GetSimulatedOrderMakerTransferResults(order LibOrderOrder, takerAddress common.Address, takerAssetFillAmount *big.Int) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetSimulatedOrderMakerTransferResults(&_WrappersV4.TransactOpts, order, takerAddress, takerAssetFillAmount)
}

// GetSimulatedOrderMakerTransferResults is a paid mutator transaction binding the contract method 0x7982653e.
//
// Solidity: function getSimulatedOrderMakerTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, address takerAddress, uint256 takerAssetFillAmount) returns(uint8 orderTransferResults)
func (_WrappersV4 *WrappersV4TransactorSession) GetSimulatedOrderMakerTransferResults(order LibOrderOrder, takerAddress common.Address, takerAssetFillAmount *big.Int) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetSimulatedOrderMakerTransferResults(&_WrappersV4.TransactOpts, order, takerAddress, takerAssetFillAmount)
}

// GetSimulatedOrderTransferResults is a paid mutator transaction binding the contract method 0xd3637905.
//
// Solidity: function getSimulatedOrderTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, address takerAddress, uint256 takerAssetFillAmount) returns(uint8 orderTransferResults)
func (_WrappersV4 *WrappersV4Transactor) GetSimulatedOrderTransferResults(opts *bind.TransactOpts, order LibOrderOrder, takerAddress common.Address, takerAssetFillAmount *big.Int) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "getSimulatedOrderTransferResults", order, takerAddress, takerAssetFillAmount)
}

// GetSimulatedOrderTransferResults is a paid mutator transaction binding the contract method 0xd3637905.
//
// Solidity: function getSimulatedOrderTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, address takerAddress, uint256 takerAssetFillAmount) returns(uint8 orderTransferResults)
func (_WrappersV4 *WrappersV4Session) GetSimulatedOrderTransferResults(order LibOrderOrder, takerAddress common.Address, takerAssetFillAmount *big.Int) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetSimulatedOrderTransferResults(&_WrappersV4.TransactOpts, order, takerAddress, takerAssetFillAmount)
}

// GetSimulatedOrderTransferResults is a paid mutator transaction binding the contract method 0xd3637905.
//
// Solidity: function getSimulatedOrderTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, address takerAddress, uint256 takerAssetFillAmount) returns(uint8 orderTransferResults)
func (_WrappersV4 *WrappersV4TransactorSession) GetSimulatedOrderTransferResults(order LibOrderOrder, takerAddress common.Address, takerAssetFillAmount *big.Int) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetSimulatedOrderTransferResults(&_WrappersV4.TransactOpts, order, takerAddress, takerAssetFillAmount)
}

// GetSimulatedOrdersTransferResults is a paid mutator transaction binding the contract method 0xa5cd62ba.
//
// Solidity: function getSimulatedOrdersTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, address[] takerAddresses, uint256[] takerAssetFillAmounts) returns(uint8[] orderTransferResults)
func (_WrappersV4 *WrappersV4Transactor) GetSimulatedOrdersTransferResults(opts *bind.TransactOpts, orders []LibOrderOrder, takerAddresses []common.Address, takerAssetFillAmounts []*big.Int) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "getSimulatedOrdersTransferResults", orders, takerAddresses, takerAssetFillAmounts)
}

// GetSimulatedOrdersTransferResults is a paid mutator transaction binding the contract method 0xa5cd62ba.
//
// Solidity: function getSimulatedOrdersTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, address[] takerAddresses, uint256[] takerAssetFillAmounts) returns(uint8[] orderTransferResults)
func (_WrappersV4 *WrappersV4Session) GetSimulatedOrdersTransferResults(orders []LibOrderOrder, takerAddresses []common.Address, takerAssetFillAmounts []*big.Int) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetSimulatedOrdersTransferResults(&_WrappersV4.TransactOpts, orders, takerAddresses, takerAssetFillAmounts)
}

// GetSimulatedOrdersTransferResults is a paid mutator transaction binding the contract method 0xa5cd62ba.
//
// Solidity: function getSimulatedOrdersTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, address[] takerAddresses, uint256[] takerAssetFillAmounts) returns(uint8[] orderTransferResults)
func (_WrappersV4 *WrappersV4TransactorSession) GetSimulatedOrdersTransferResults(orders []LibOrderOrder, takerAddresses []common.Address, takerAssetFillAmounts []*big.Int) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetSimulatedOrdersTransferResults(&_WrappersV4.TransactOpts, orders, takerAddresses, takerAssetFillAmounts)
}

// GetTransferableAssetAmount is a paid mutator transaction binding the contract method 0x2322cf76.
//
// Solidity: function getTransferableAssetAmount(address ownerAddress, bytes assetData) returns(uint256 transferableAssetAmount)
func (_WrappersV4 *WrappersV4Transactor) GetTransferableAssetAmount(opts *bind.TransactOpts, ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "getTransferableAssetAmount", ownerAddress, assetData)
}

// GetTransferableAssetAmount is a paid mutator transaction binding the contract method 0x2322cf76.
//
// Solidity: function getTransferableAssetAmount(address ownerAddress, bytes assetData) returns(uint256 transferableAssetAmount)
func (_WrappersV4 *WrappersV4Session) GetTransferableAssetAmount(ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetTransferableAssetAmount(&_WrappersV4.TransactOpts, ownerAddress, assetData)
}

// GetTransferableAssetAmount is a paid mutator transaction binding the contract method 0x2322cf76.
//
// Solidity: function getTransferableAssetAmount(address ownerAddress, bytes assetData) returns(uint256 transferableAssetAmount)
func (_WrappersV4 *WrappersV4TransactorSession) GetTransferableAssetAmount(ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.GetTransferableAssetAmount(&_WrappersV4.TransactOpts, ownerAddress, assetData)
}
