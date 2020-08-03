// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package wrappers

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

// DevUtilsABI is the input ABI used to generate the binding from.
const DevUtilsABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"exchange_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"chaiBridge_\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"dydxBridge_\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[],\"name\":\"EIP712_EXCHANGE_DOMAIN_HASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chaiBridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeAssetProxyId\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC1155AssetData\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenValues\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"callbackData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC20AssetData\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC20BridgeAssetData\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"bridgeAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"bridgeData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC721AssetData\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeMultiAssetData\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"nestedAssetData\",\"type\":\"bytes[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeStaticCallAssetData\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"staticCallTargetAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"staticCallData\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"expectedReturnDataHash\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"transactionData\",\"type\":\"bytes\"}],\"name\":\"decodeZeroExTransactionData\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"functionName\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"dydxBridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenValues\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"callbackData\",\"type\":\"bytes\"}],\"name\":\"encodeERC1155AssetData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"encodeERC20AssetData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"encodeERC721AssetData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"nestedAssetData\",\"type\":\"bytes[]\"}],\"name\":\"encodeMultiAssetData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"staticCallTargetAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"staticCallData\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"expectedReturnDataHash\",\"type\":\"bytes32\"}],\"name\":\"encodeStaticCallAssetData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"erc1155ProxyAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"erc20ProxyAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"erc721ProxyAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"exchangeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getAssetProxyAllowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getBalanceAndAssetProxyAllowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes[]\",\"name\":\"assetData\",\"type\":\"bytes[]\"}],\"name\":\"getBatchAssetProxyAllowances\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"allowances\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes[]\",\"name\":\"assetData\",\"type\":\"bytes[]\"}],\"name\":\"getBatchBalances\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes[]\",\"name\":\"assetData\",\"type\":\"bytes[]\"}],\"name\":\"getBatchBalancesAndAssetProxyAllowances\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"allowances\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"addresses\",\"type\":\"address[]\"}],\"name\":\"getEthBalances\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"exchange\",\"type\":\"address\"}],\"name\":\"getOrderHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"getOrderRelevantState\",\"outputs\":[{\"components\":[{\"internalType\":\"enumLibOrder.OrderStatus\",\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"internalType\":\"structLibOrder.OrderInfo\",\"name\":\"orderInfo\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"fillableTakerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isValidSignature\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"getOrderRelevantStates\",\"outputs\":[{\"components\":[{\"internalType\":\"enumLibOrder.OrderStatus\",\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"internalType\":\"structLibOrder.OrderInfo[]\",\"name\":\"ordersInfo\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"fillableTakerAssetAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bool[]\",\"name\":\"isValidSignature\",\"type\":\"bool[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"}],\"name\":\"getSimulatedOrderMakerTransferResults\",\"outputs\":[{\"internalType\":\"enumLibOrderTransferSimulation.OrderTransferResults\",\"name\":\"orderTransferResults\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"}],\"name\":\"getSimulatedOrderTransferResults\",\"outputs\":[{\"internalType\":\"enumLibOrderTransferSimulation.OrderTransferResults\",\"name\":\"orderTransferResults\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"address[]\",\"name\":\"takerAddresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"}],\"name\":\"getSimulatedOrdersTransferResults\",\"outputs\":[{\"internalType\":\"enumLibOrderTransferSimulation.OrderTransferResults[]\",\"name\":\"orderTransferResults\",\"type\":\"uint8[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLibZeroExTransaction.ZeroExTransaction\",\"name\":\"transaction\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"exchange\",\"type\":\"address\"}],\"name\":\"getTransactionHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"transactionHash\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getTransferableAssetAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"transferableAssetAmount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"revertIfInvalidAssetData\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"staticCallProxyAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// DevUtils is an auto generated Go binding around an Ethereum contract.
type DevUtils struct {
	DevUtilsCaller     // Read-only binding to the contract
	DevUtilsTransactor // Write-only binding to the contract
	DevUtilsFilterer   // Log filterer for contract events
}

// DevUtilsCaller is an auto generated read-only Go binding around an Ethereum contract.
type DevUtilsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DevUtilsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DevUtilsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DevUtilsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DevUtilsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DevUtilsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DevUtilsSession struct {
	Contract     *DevUtils         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DevUtilsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DevUtilsCallerSession struct {
	Contract *DevUtilsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// DevUtilsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DevUtilsTransactorSession struct {
	Contract     *DevUtilsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// DevUtilsRaw is an auto generated low-level Go binding around an Ethereum contract.
type DevUtilsRaw struct {
	Contract *DevUtils // Generic contract binding to access the raw methods on
}

// DevUtilsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DevUtilsCallerRaw struct {
	Contract *DevUtilsCaller // Generic read-only contract binding to access the raw methods on
}

// DevUtilsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DevUtilsTransactorRaw struct {
	Contract *DevUtilsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDevUtils creates a new instance of DevUtils, bound to a specific deployed contract.
func NewDevUtils(address common.Address, backend bind.ContractBackend) (*DevUtils, error) {
	contract, err := bindDevUtils(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DevUtils{DevUtilsCaller: DevUtilsCaller{contract: contract}, DevUtilsTransactor: DevUtilsTransactor{contract: contract}, DevUtilsFilterer: DevUtilsFilterer{contract: contract}}, nil
}

// NewDevUtilsCaller creates a new read-only instance of DevUtils, bound to a specific deployed contract.
func NewDevUtilsCaller(address common.Address, caller bind.ContractCaller) (*DevUtilsCaller, error) {
	contract, err := bindDevUtils(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DevUtilsCaller{contract: contract}, nil
}

// NewDevUtilsTransactor creates a new write-only instance of DevUtils, bound to a specific deployed contract.
func NewDevUtilsTransactor(address common.Address, transactor bind.ContractTransactor) (*DevUtilsTransactor, error) {
	contract, err := bindDevUtils(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DevUtilsTransactor{contract: contract}, nil
}

// NewDevUtilsFilterer creates a new log filterer instance of DevUtils, bound to a specific deployed contract.
func NewDevUtilsFilterer(address common.Address, filterer bind.ContractFilterer) (*DevUtilsFilterer, error) {
	contract, err := bindDevUtils(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DevUtilsFilterer{contract: contract}, nil
}

// bindDevUtils binds a generic wrapper to an already deployed contract.
func bindDevUtils(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DevUtilsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DevUtils *DevUtilsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DevUtils.Contract.DevUtilsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DevUtils *DevUtilsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DevUtils.Contract.DevUtilsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DevUtils *DevUtilsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DevUtils.Contract.DevUtilsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DevUtils *DevUtilsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DevUtils.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DevUtils *DevUtilsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DevUtils.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DevUtils *DevUtilsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DevUtils.Contract.contract.Transact(opts, method, params...)
}

// EIP712EXCHANGEDOMAINHASH is a free data retrieval call binding the contract method 0xc26cfecd.
//
// Solidity: function EIP712_EXCHANGE_DOMAIN_HASH() view returns(bytes32)
func (_DevUtils *DevUtilsCaller) EIP712EXCHANGEDOMAINHASH(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "EIP712_EXCHANGE_DOMAIN_HASH")
	return *ret0, err
}

// EIP712EXCHANGEDOMAINHASH is a free data retrieval call binding the contract method 0xc26cfecd.
//
// Solidity: function EIP712_EXCHANGE_DOMAIN_HASH() view returns(bytes32)
func (_DevUtils *DevUtilsSession) EIP712EXCHANGEDOMAINHASH() ([32]byte, error) {
	return _DevUtils.Contract.EIP712EXCHANGEDOMAINHASH(&_DevUtils.CallOpts)
}

// EIP712EXCHANGEDOMAINHASH is a free data retrieval call binding the contract method 0xc26cfecd.
//
// Solidity: function EIP712_EXCHANGE_DOMAIN_HASH() view returns(bytes32)
func (_DevUtils *DevUtilsCallerSession) EIP712EXCHANGEDOMAINHASH() ([32]byte, error) {
	return _DevUtils.Contract.EIP712EXCHANGEDOMAINHASH(&_DevUtils.CallOpts)
}

// ChaiBridgeAddress is a free data retrieval call binding the contract method 0xc82037ef.
//
// Solidity: function chaiBridgeAddress() view returns(address)
func (_DevUtils *DevUtilsCaller) ChaiBridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "chaiBridgeAddress")
	return *ret0, err
}

// ChaiBridgeAddress is a free data retrieval call binding the contract method 0xc82037ef.
//
// Solidity: function chaiBridgeAddress() view returns(address)
func (_DevUtils *DevUtilsSession) ChaiBridgeAddress() (common.Address, error) {
	return _DevUtils.Contract.ChaiBridgeAddress(&_DevUtils.CallOpts)
}

// ChaiBridgeAddress is a free data retrieval call binding the contract method 0xc82037ef.
//
// Solidity: function chaiBridgeAddress() view returns(address)
func (_DevUtils *DevUtilsCallerSession) ChaiBridgeAddress() (common.Address, error) {
	return _DevUtils.Contract.ChaiBridgeAddress(&_DevUtils.CallOpts)
}

// DecodeAssetProxyId is a free data retrieval call binding the contract method 0xd4695028.
//
// Solidity: function decodeAssetProxyId(bytes assetData) pure returns(bytes4 assetProxyId)
func (_DevUtils *DevUtilsCaller) DecodeAssetProxyId(opts *bind.CallOpts, assetData []byte) ([4]byte, error) {
	var (
		ret0 = new([4]byte)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "decodeAssetProxyId", assetData)
	return *ret0, err
}

// DecodeAssetProxyId is a free data retrieval call binding the contract method 0xd4695028.
//
// Solidity: function decodeAssetProxyId(bytes assetData) pure returns(bytes4 assetProxyId)
func (_DevUtils *DevUtilsSession) DecodeAssetProxyId(assetData []byte) ([4]byte, error) {
	return _DevUtils.Contract.DecodeAssetProxyId(&_DevUtils.CallOpts, assetData)
}

// DecodeAssetProxyId is a free data retrieval call binding the contract method 0xd4695028.
//
// Solidity: function decodeAssetProxyId(bytes assetData) pure returns(bytes4 assetProxyId)
func (_DevUtils *DevUtilsCallerSession) DecodeAssetProxyId(assetData []byte) ([4]byte, error) {
	return _DevUtils.Contract.DecodeAssetProxyId(&_DevUtils.CallOpts, assetData)
}

// DecodeERC1155AssetData is a free data retrieval call binding the contract method 0x9eadc835.
//
// Solidity: function decodeERC1155AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData)
func (_DevUtils *DevUtilsCaller) DecodeERC1155AssetData(opts *bind.CallOpts, assetData []byte) (struct {
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
	err := _DevUtils.contract.Call(opts, out, "decodeERC1155AssetData", assetData)
	return *ret, err
}

// DecodeERC1155AssetData is a free data retrieval call binding the contract method 0x9eadc835.
//
// Solidity: function decodeERC1155AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData)
func (_DevUtils *DevUtilsSession) DecodeERC1155AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
	TokenIds     []*big.Int
	TokenValues  []*big.Int
	CallbackData []byte
}, error) {
	return _DevUtils.Contract.DecodeERC1155AssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeERC1155AssetData is a free data retrieval call binding the contract method 0x9eadc835.
//
// Solidity: function decodeERC1155AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData)
func (_DevUtils *DevUtilsCallerSession) DecodeERC1155AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
	TokenIds     []*big.Int
	TokenValues  []*big.Int
	CallbackData []byte
}, error) {
	return _DevUtils.Contract.DecodeERC1155AssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeERC20AssetData is a free data retrieval call binding the contract method 0x8f4ce479.
//
// Solidity: function decodeERC20AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress)
func (_DevUtils *DevUtilsCaller) DecodeERC20AssetData(opts *bind.CallOpts, assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
}, error) {
	ret := new(struct {
		AssetProxyId [4]byte
		TokenAddress common.Address
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeERC20AssetData", assetData)
	return *ret, err
}

// DecodeERC20AssetData is a free data retrieval call binding the contract method 0x8f4ce479.
//
// Solidity: function decodeERC20AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress)
func (_DevUtils *DevUtilsSession) DecodeERC20AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
}, error) {
	return _DevUtils.Contract.DecodeERC20AssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeERC20AssetData is a free data retrieval call binding the contract method 0x8f4ce479.
//
// Solidity: function decodeERC20AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress)
func (_DevUtils *DevUtilsCallerSession) DecodeERC20AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
}, error) {
	return _DevUtils.Contract.DecodeERC20AssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeERC20BridgeAssetData is a free data retrieval call binding the contract method 0x750bdb30.
//
// Solidity: function decodeERC20BridgeAssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, address bridgeAddress, bytes bridgeData)
func (_DevUtils *DevUtilsCaller) DecodeERC20BridgeAssetData(opts *bind.CallOpts, assetData []byte) (struct {
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
	err := _DevUtils.contract.Call(opts, out, "decodeERC20BridgeAssetData", assetData)
	return *ret, err
}

// DecodeERC20BridgeAssetData is a free data retrieval call binding the contract method 0x750bdb30.
//
// Solidity: function decodeERC20BridgeAssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, address bridgeAddress, bytes bridgeData)
func (_DevUtils *DevUtilsSession) DecodeERC20BridgeAssetData(assetData []byte) (struct {
	AssetProxyId  [4]byte
	TokenAddress  common.Address
	BridgeAddress common.Address
	BridgeData    []byte
}, error) {
	return _DevUtils.Contract.DecodeERC20BridgeAssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeERC20BridgeAssetData is a free data retrieval call binding the contract method 0x750bdb30.
//
// Solidity: function decodeERC20BridgeAssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, address bridgeAddress, bytes bridgeData)
func (_DevUtils *DevUtilsCallerSession) DecodeERC20BridgeAssetData(assetData []byte) (struct {
	AssetProxyId  [4]byte
	TokenAddress  common.Address
	BridgeAddress common.Address
	BridgeData    []byte
}, error) {
	return _DevUtils.Contract.DecodeERC20BridgeAssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeERC721AssetData is a free data retrieval call binding the contract method 0x04a5618a.
//
// Solidity: function decodeERC721AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, uint256 tokenId)
func (_DevUtils *DevUtilsCaller) DecodeERC721AssetData(opts *bind.CallOpts, assetData []byte) (struct {
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
	err := _DevUtils.contract.Call(opts, out, "decodeERC721AssetData", assetData)
	return *ret, err
}

// DecodeERC721AssetData is a free data retrieval call binding the contract method 0x04a5618a.
//
// Solidity: function decodeERC721AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, uint256 tokenId)
func (_DevUtils *DevUtilsSession) DecodeERC721AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
	TokenId      *big.Int
}, error) {
	return _DevUtils.Contract.DecodeERC721AssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeERC721AssetData is a free data retrieval call binding the contract method 0x04a5618a.
//
// Solidity: function decodeERC721AssetData(bytes assetData) pure returns(bytes4 assetProxyId, address tokenAddress, uint256 tokenId)
func (_DevUtils *DevUtilsCallerSession) DecodeERC721AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
	TokenId      *big.Int
}, error) {
	return _DevUtils.Contract.DecodeERC721AssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeMultiAssetData is a free data retrieval call binding the contract method 0xbbb2dcf6.
//
// Solidity: function decodeMultiAssetData(bytes assetData) pure returns(bytes4 assetProxyId, uint256[] amounts, bytes[] nestedAssetData)
func (_DevUtils *DevUtilsCaller) DecodeMultiAssetData(opts *bind.CallOpts, assetData []byte) (struct {
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
	err := _DevUtils.contract.Call(opts, out, "decodeMultiAssetData", assetData)
	return *ret, err
}

// DecodeMultiAssetData is a free data retrieval call binding the contract method 0xbbb2dcf6.
//
// Solidity: function decodeMultiAssetData(bytes assetData) pure returns(bytes4 assetProxyId, uint256[] amounts, bytes[] nestedAssetData)
func (_DevUtils *DevUtilsSession) DecodeMultiAssetData(assetData []byte) (struct {
	AssetProxyId    [4]byte
	Amounts         []*big.Int
	NestedAssetData [][]byte
}, error) {
	return _DevUtils.Contract.DecodeMultiAssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeMultiAssetData is a free data retrieval call binding the contract method 0xbbb2dcf6.
//
// Solidity: function decodeMultiAssetData(bytes assetData) pure returns(bytes4 assetProxyId, uint256[] amounts, bytes[] nestedAssetData)
func (_DevUtils *DevUtilsCallerSession) DecodeMultiAssetData(assetData []byte) (struct {
	AssetProxyId    [4]byte
	Amounts         []*big.Int
	NestedAssetData [][]byte
}, error) {
	return _DevUtils.Contract.DecodeMultiAssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeStaticCallAssetData is a free data retrieval call binding the contract method 0xca49f47c.
//
// Solidity: function decodeStaticCallAssetData(bytes assetData) pure returns(bytes4 assetProxyId, address staticCallTargetAddress, bytes staticCallData, bytes32 expectedReturnDataHash)
func (_DevUtils *DevUtilsCaller) DecodeStaticCallAssetData(opts *bind.CallOpts, assetData []byte) (struct {
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
	err := _DevUtils.contract.Call(opts, out, "decodeStaticCallAssetData", assetData)
	return *ret, err
}

// DecodeStaticCallAssetData is a free data retrieval call binding the contract method 0xca49f47c.
//
// Solidity: function decodeStaticCallAssetData(bytes assetData) pure returns(bytes4 assetProxyId, address staticCallTargetAddress, bytes staticCallData, bytes32 expectedReturnDataHash)
func (_DevUtils *DevUtilsSession) DecodeStaticCallAssetData(assetData []byte) (struct {
	AssetProxyId            [4]byte
	StaticCallTargetAddress common.Address
	StaticCallData          []byte
	ExpectedReturnDataHash  [32]byte
}, error) {
	return _DevUtils.Contract.DecodeStaticCallAssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeStaticCallAssetData is a free data retrieval call binding the contract method 0xca49f47c.
//
// Solidity: function decodeStaticCallAssetData(bytes assetData) pure returns(bytes4 assetProxyId, address staticCallTargetAddress, bytes staticCallData, bytes32 expectedReturnDataHash)
func (_DevUtils *DevUtilsCallerSession) DecodeStaticCallAssetData(assetData []byte) (struct {
	AssetProxyId            [4]byte
	StaticCallTargetAddress common.Address
	StaticCallData          []byte
	ExpectedReturnDataHash  [32]byte
}, error) {
	return _DevUtils.Contract.DecodeStaticCallAssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeZeroExTransactionData is a free data retrieval call binding the contract method 0x6f83188e.
//
// Solidity: function decodeZeroExTransactionData(bytes transactionData) pure returns(string functionName, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures)
func (_DevUtils *DevUtilsCaller) DecodeZeroExTransactionData(opts *bind.CallOpts, transactionData []byte) (struct {
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
	err := _DevUtils.contract.Call(opts, out, "decodeZeroExTransactionData", transactionData)
	return *ret, err
}

// DecodeZeroExTransactionData is a free data retrieval call binding the contract method 0x6f83188e.
//
// Solidity: function decodeZeroExTransactionData(bytes transactionData) pure returns(string functionName, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures)
func (_DevUtils *DevUtilsSession) DecodeZeroExTransactionData(transactionData []byte) (struct {
	FunctionName          string
	Orders                []LibOrderOrder
	TakerAssetFillAmounts []*big.Int
	Signatures            [][]byte
}, error) {
	return _DevUtils.Contract.DecodeZeroExTransactionData(&_DevUtils.CallOpts, transactionData)
}

// DecodeZeroExTransactionData is a free data retrieval call binding the contract method 0x6f83188e.
//
// Solidity: function decodeZeroExTransactionData(bytes transactionData) pure returns(string functionName, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures)
func (_DevUtils *DevUtilsCallerSession) DecodeZeroExTransactionData(transactionData []byte) (struct {
	FunctionName          string
	Orders                []LibOrderOrder
	TakerAssetFillAmounts []*big.Int
	Signatures            [][]byte
}, error) {
	return _DevUtils.Contract.DecodeZeroExTransactionData(&_DevUtils.CallOpts, transactionData)
}

// DydxBridgeAddress is a free data retrieval call binding the contract method 0xa7530f12.
//
// Solidity: function dydxBridgeAddress() view returns(address)
func (_DevUtils *DevUtilsCaller) DydxBridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "dydxBridgeAddress")
	return *ret0, err
}

// DydxBridgeAddress is a free data retrieval call binding the contract method 0xa7530f12.
//
// Solidity: function dydxBridgeAddress() view returns(address)
func (_DevUtils *DevUtilsSession) DydxBridgeAddress() (common.Address, error) {
	return _DevUtils.Contract.DydxBridgeAddress(&_DevUtils.CallOpts)
}

// DydxBridgeAddress is a free data retrieval call binding the contract method 0xa7530f12.
//
// Solidity: function dydxBridgeAddress() view returns(address)
func (_DevUtils *DevUtilsCallerSession) DydxBridgeAddress() (common.Address, error) {
	return _DevUtils.Contract.DydxBridgeAddress(&_DevUtils.CallOpts)
}

// EncodeERC1155AssetData is a free data retrieval call binding the contract method 0xb43cffe1.
//
// Solidity: function encodeERC1155AssetData(address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData) pure returns(bytes assetData)
func (_DevUtils *DevUtilsCaller) EncodeERC1155AssetData(opts *bind.CallOpts, tokenAddress common.Address, tokenIds []*big.Int, tokenValues []*big.Int, callbackData []byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "encodeERC1155AssetData", tokenAddress, tokenIds, tokenValues, callbackData)
	return *ret0, err
}

// EncodeERC1155AssetData is a free data retrieval call binding the contract method 0xb43cffe1.
//
// Solidity: function encodeERC1155AssetData(address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData) pure returns(bytes assetData)
func (_DevUtils *DevUtilsSession) EncodeERC1155AssetData(tokenAddress common.Address, tokenIds []*big.Int, tokenValues []*big.Int, callbackData []byte) ([]byte, error) {
	return _DevUtils.Contract.EncodeERC1155AssetData(&_DevUtils.CallOpts, tokenAddress, tokenIds, tokenValues, callbackData)
}

// EncodeERC1155AssetData is a free data retrieval call binding the contract method 0xb43cffe1.
//
// Solidity: function encodeERC1155AssetData(address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData) pure returns(bytes assetData)
func (_DevUtils *DevUtilsCallerSession) EncodeERC1155AssetData(tokenAddress common.Address, tokenIds []*big.Int, tokenValues []*big.Int, callbackData []byte) ([]byte, error) {
	return _DevUtils.Contract.EncodeERC1155AssetData(&_DevUtils.CallOpts, tokenAddress, tokenIds, tokenValues, callbackData)
}

// EncodeERC20AssetData is a free data retrieval call binding the contract method 0x590aa875.
//
// Solidity: function encodeERC20AssetData(address tokenAddress) pure returns(bytes assetData)
func (_DevUtils *DevUtilsCaller) EncodeERC20AssetData(opts *bind.CallOpts, tokenAddress common.Address) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "encodeERC20AssetData", tokenAddress)
	return *ret0, err
}

// EncodeERC20AssetData is a free data retrieval call binding the contract method 0x590aa875.
//
// Solidity: function encodeERC20AssetData(address tokenAddress) pure returns(bytes assetData)
func (_DevUtils *DevUtilsSession) EncodeERC20AssetData(tokenAddress common.Address) ([]byte, error) {
	return _DevUtils.Contract.EncodeERC20AssetData(&_DevUtils.CallOpts, tokenAddress)
}

// EncodeERC20AssetData is a free data retrieval call binding the contract method 0x590aa875.
//
// Solidity: function encodeERC20AssetData(address tokenAddress) pure returns(bytes assetData)
func (_DevUtils *DevUtilsCallerSession) EncodeERC20AssetData(tokenAddress common.Address) ([]byte, error) {
	return _DevUtils.Contract.EncodeERC20AssetData(&_DevUtils.CallOpts, tokenAddress)
}

// EncodeERC721AssetData is a free data retrieval call binding the contract method 0xa6627e9f.
//
// Solidity: function encodeERC721AssetData(address tokenAddress, uint256 tokenId) pure returns(bytes assetData)
func (_DevUtils *DevUtilsCaller) EncodeERC721AssetData(opts *bind.CallOpts, tokenAddress common.Address, tokenId *big.Int) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "encodeERC721AssetData", tokenAddress, tokenId)
	return *ret0, err
}

// EncodeERC721AssetData is a free data retrieval call binding the contract method 0xa6627e9f.
//
// Solidity: function encodeERC721AssetData(address tokenAddress, uint256 tokenId) pure returns(bytes assetData)
func (_DevUtils *DevUtilsSession) EncodeERC721AssetData(tokenAddress common.Address, tokenId *big.Int) ([]byte, error) {
	return _DevUtils.Contract.EncodeERC721AssetData(&_DevUtils.CallOpts, tokenAddress, tokenId)
}

// EncodeERC721AssetData is a free data retrieval call binding the contract method 0xa6627e9f.
//
// Solidity: function encodeERC721AssetData(address tokenAddress, uint256 tokenId) pure returns(bytes assetData)
func (_DevUtils *DevUtilsCallerSession) EncodeERC721AssetData(tokenAddress common.Address, tokenId *big.Int) ([]byte, error) {
	return _DevUtils.Contract.EncodeERC721AssetData(&_DevUtils.CallOpts, tokenAddress, tokenId)
}

// EncodeMultiAssetData is a free data retrieval call binding the contract method 0xd3d862d1.
//
// Solidity: function encodeMultiAssetData(uint256[] amounts, bytes[] nestedAssetData) pure returns(bytes assetData)
func (_DevUtils *DevUtilsCaller) EncodeMultiAssetData(opts *bind.CallOpts, amounts []*big.Int, nestedAssetData [][]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "encodeMultiAssetData", amounts, nestedAssetData)
	return *ret0, err
}

// EncodeMultiAssetData is a free data retrieval call binding the contract method 0xd3d862d1.
//
// Solidity: function encodeMultiAssetData(uint256[] amounts, bytes[] nestedAssetData) pure returns(bytes assetData)
func (_DevUtils *DevUtilsSession) EncodeMultiAssetData(amounts []*big.Int, nestedAssetData [][]byte) ([]byte, error) {
	return _DevUtils.Contract.EncodeMultiAssetData(&_DevUtils.CallOpts, amounts, nestedAssetData)
}

// EncodeMultiAssetData is a free data retrieval call binding the contract method 0xd3d862d1.
//
// Solidity: function encodeMultiAssetData(uint256[] amounts, bytes[] nestedAssetData) pure returns(bytes assetData)
func (_DevUtils *DevUtilsCallerSession) EncodeMultiAssetData(amounts []*big.Int, nestedAssetData [][]byte) ([]byte, error) {
	return _DevUtils.Contract.EncodeMultiAssetData(&_DevUtils.CallOpts, amounts, nestedAssetData)
}

// EncodeStaticCallAssetData is a free data retrieval call binding the contract method 0x63eb3992.
//
// Solidity: function encodeStaticCallAssetData(address staticCallTargetAddress, bytes staticCallData, bytes32 expectedReturnDataHash) pure returns(bytes assetData)
func (_DevUtils *DevUtilsCaller) EncodeStaticCallAssetData(opts *bind.CallOpts, staticCallTargetAddress common.Address, staticCallData []byte, expectedReturnDataHash [32]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "encodeStaticCallAssetData", staticCallTargetAddress, staticCallData, expectedReturnDataHash)
	return *ret0, err
}

// EncodeStaticCallAssetData is a free data retrieval call binding the contract method 0x63eb3992.
//
// Solidity: function encodeStaticCallAssetData(address staticCallTargetAddress, bytes staticCallData, bytes32 expectedReturnDataHash) pure returns(bytes assetData)
func (_DevUtils *DevUtilsSession) EncodeStaticCallAssetData(staticCallTargetAddress common.Address, staticCallData []byte, expectedReturnDataHash [32]byte) ([]byte, error) {
	return _DevUtils.Contract.EncodeStaticCallAssetData(&_DevUtils.CallOpts, staticCallTargetAddress, staticCallData, expectedReturnDataHash)
}

// EncodeStaticCallAssetData is a free data retrieval call binding the contract method 0x63eb3992.
//
// Solidity: function encodeStaticCallAssetData(address staticCallTargetAddress, bytes staticCallData, bytes32 expectedReturnDataHash) pure returns(bytes assetData)
func (_DevUtils *DevUtilsCallerSession) EncodeStaticCallAssetData(staticCallTargetAddress common.Address, staticCallData []byte, expectedReturnDataHash [32]byte) ([]byte, error) {
	return _DevUtils.Contract.EncodeStaticCallAssetData(&_DevUtils.CallOpts, staticCallTargetAddress, staticCallData, expectedReturnDataHash)
}

// Erc1155ProxyAddress is a free data retrieval call binding the contract method 0xff84e7cc.
//
// Solidity: function erc1155ProxyAddress() view returns(address)
func (_DevUtils *DevUtilsCaller) Erc1155ProxyAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "erc1155ProxyAddress")
	return *ret0, err
}

// Erc1155ProxyAddress is a free data retrieval call binding the contract method 0xff84e7cc.
//
// Solidity: function erc1155ProxyAddress() view returns(address)
func (_DevUtils *DevUtilsSession) Erc1155ProxyAddress() (common.Address, error) {
	return _DevUtils.Contract.Erc1155ProxyAddress(&_DevUtils.CallOpts)
}

// Erc1155ProxyAddress is a free data retrieval call binding the contract method 0xff84e7cc.
//
// Solidity: function erc1155ProxyAddress() view returns(address)
func (_DevUtils *DevUtilsCallerSession) Erc1155ProxyAddress() (common.Address, error) {
	return _DevUtils.Contract.Erc1155ProxyAddress(&_DevUtils.CallOpts)
}

// Erc20ProxyAddress is a free data retrieval call binding the contract method 0xee185997.
//
// Solidity: function erc20ProxyAddress() view returns(address)
func (_DevUtils *DevUtilsCaller) Erc20ProxyAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "erc20ProxyAddress")
	return *ret0, err
}

// Erc20ProxyAddress is a free data retrieval call binding the contract method 0xee185997.
//
// Solidity: function erc20ProxyAddress() view returns(address)
func (_DevUtils *DevUtilsSession) Erc20ProxyAddress() (common.Address, error) {
	return _DevUtils.Contract.Erc20ProxyAddress(&_DevUtils.CallOpts)
}

// Erc20ProxyAddress is a free data retrieval call binding the contract method 0xee185997.
//
// Solidity: function erc20ProxyAddress() view returns(address)
func (_DevUtils *DevUtilsCallerSession) Erc20ProxyAddress() (common.Address, error) {
	return _DevUtils.Contract.Erc20ProxyAddress(&_DevUtils.CallOpts)
}

// Erc721ProxyAddress is a free data retrieval call binding the contract method 0xef3bb097.
//
// Solidity: function erc721ProxyAddress() view returns(address)
func (_DevUtils *DevUtilsCaller) Erc721ProxyAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "erc721ProxyAddress")
	return *ret0, err
}

// Erc721ProxyAddress is a free data retrieval call binding the contract method 0xef3bb097.
//
// Solidity: function erc721ProxyAddress() view returns(address)
func (_DevUtils *DevUtilsSession) Erc721ProxyAddress() (common.Address, error) {
	return _DevUtils.Contract.Erc721ProxyAddress(&_DevUtils.CallOpts)
}

// Erc721ProxyAddress is a free data retrieval call binding the contract method 0xef3bb097.
//
// Solidity: function erc721ProxyAddress() view returns(address)
func (_DevUtils *DevUtilsCallerSession) Erc721ProxyAddress() (common.Address, error) {
	return _DevUtils.Contract.Erc721ProxyAddress(&_DevUtils.CallOpts)
}

// ExchangeAddress is a free data retrieval call binding the contract method 0x9cd01605.
//
// Solidity: function exchangeAddress() view returns(address)
func (_DevUtils *DevUtilsCaller) ExchangeAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "exchangeAddress")
	return *ret0, err
}

// ExchangeAddress is a free data retrieval call binding the contract method 0x9cd01605.
//
// Solidity: function exchangeAddress() view returns(address)
func (_DevUtils *DevUtilsSession) ExchangeAddress() (common.Address, error) {
	return _DevUtils.Contract.ExchangeAddress(&_DevUtils.CallOpts)
}

// ExchangeAddress is a free data retrieval call binding the contract method 0x9cd01605.
//
// Solidity: function exchangeAddress() view returns(address)
func (_DevUtils *DevUtilsCallerSession) ExchangeAddress() (common.Address, error) {
	return _DevUtils.Contract.ExchangeAddress(&_DevUtils.CallOpts)
}

// GetEthBalances is a free data retrieval call binding the contract method 0xa0901e51.
//
// Solidity: function getEthBalances(address[] addresses) view returns(uint256[])
func (_DevUtils *DevUtilsCaller) GetEthBalances(opts *bind.CallOpts, addresses []common.Address) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "getEthBalances", addresses)
	return *ret0, err
}

// GetEthBalances is a free data retrieval call binding the contract method 0xa0901e51.
//
// Solidity: function getEthBalances(address[] addresses) view returns(uint256[])
func (_DevUtils *DevUtilsSession) GetEthBalances(addresses []common.Address) ([]*big.Int, error) {
	return _DevUtils.Contract.GetEthBalances(&_DevUtils.CallOpts, addresses)
}

// GetEthBalances is a free data retrieval call binding the contract method 0xa0901e51.
//
// Solidity: function getEthBalances(address[] addresses) view returns(uint256[])
func (_DevUtils *DevUtilsCallerSession) GetEthBalances(addresses []common.Address) ([]*big.Int, error) {
	return _DevUtils.Contract.GetEthBalances(&_DevUtils.CallOpts, addresses)
}

// GetOrderHash is a free data retrieval call binding the contract method 0xa070cac8.
//
// Solidity: function getOrderHash((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 chainId, address exchange) pure returns(bytes32 orderHash)
func (_DevUtils *DevUtilsCaller) GetOrderHash(opts *bind.CallOpts, order LibOrderOrder, chainId *big.Int, exchange common.Address) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "getOrderHash", order, chainId, exchange)
	return *ret0, err
}

// GetOrderHash is a free data retrieval call binding the contract method 0xa070cac8.
//
// Solidity: function getOrderHash((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 chainId, address exchange) pure returns(bytes32 orderHash)
func (_DevUtils *DevUtilsSession) GetOrderHash(order LibOrderOrder, chainId *big.Int, exchange common.Address) ([32]byte, error) {
	return _DevUtils.Contract.GetOrderHash(&_DevUtils.CallOpts, order, chainId, exchange)
}

// GetOrderHash is a free data retrieval call binding the contract method 0xa070cac8.
//
// Solidity: function getOrderHash((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 chainId, address exchange) pure returns(bytes32 orderHash)
func (_DevUtils *DevUtilsCallerSession) GetOrderHash(order LibOrderOrder, chainId *big.Int, exchange common.Address) ([32]byte, error) {
	return _DevUtils.Contract.GetOrderHash(&_DevUtils.CallOpts, order, chainId, exchange)
}

// HACK(jalextowle): This function is technically non-payable, but we mark it as "view" in the
// abi so that we can get it to generate a function that can use "ethcall"
// GetOrderRelevantState is a free data retrieval call binding the contract method 0xe77286eb.
//
// Solidity: function getOrderRelevantState((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, bytes signature) view returns((uint8,bytes32,uint256) orderInfo, uint256 fillableTakerAssetAmount, bool isValidSignature)
func (_DevUtils *DevUtilsCaller) GetOrderRelevantState(opts *bind.CallOpts, order LibOrderOrder, signature []byte) (struct {
	OrderInfo                LibOrderOrderInfo
	FillableTakerAssetAmount *big.Int
	IsValidSignature         bool
}, error) {
	ret := new(struct {
		OrderInfo                LibOrderOrderInfo
		FillableTakerAssetAmount *big.Int
		IsValidSignature         bool
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "getOrderRelevantState", order, signature)
	return *ret, err
}

// GetOrderRelevantState is a free data retrieval call binding the contract method 0xe77286eb.
//
// Solidity: function getOrderRelevantState((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, bytes signature) view returns((uint8,bytes32,uint256) orderInfo, uint256 fillableTakerAssetAmount, bool isValidSignature)
func (_DevUtils *DevUtilsSession) GetOrderRelevantState(order LibOrderOrder, signature []byte) (struct {
	OrderInfo                LibOrderOrderInfo
	FillableTakerAssetAmount *big.Int
	IsValidSignature         bool
}, error) {
	return _DevUtils.Contract.GetOrderRelevantState(&_DevUtils.CallOpts, order, signature)
}

// GetOrderRelevantState is a free data retrieval call binding the contract method 0xe77286eb.
//
// Solidity: function getOrderRelevantState((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, bytes signature) view returns((uint8,bytes32,uint256) orderInfo, uint256 fillableTakerAssetAmount, bool isValidSignature)
func (_DevUtils *DevUtilsCallerSession) GetOrderRelevantState(order LibOrderOrder, signature []byte) (struct {
	OrderInfo                LibOrderOrderInfo
	FillableTakerAssetAmount *big.Int
	IsValidSignature         bool
}, error) {
	return _DevUtils.Contract.GetOrderRelevantState(&_DevUtils.CallOpts, order, signature)
}

// GetOrderRelevantStates is a free data retrieval call binding the contract method 0xe25cabf7.
//
// Solidity: function getOrderRelevantStates((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, bytes[] signatures) view returns((uint8,bytes32,uint256)[] ordersInfo, uint256[] fillableTakerAssetAmounts, bool[] isValidSignature)
func (_DevUtils *DevUtilsCaller) GetOrderRelevantStates(opts *bind.CallOpts, orders []LibOrderOrder, signatures [][]byte) (struct {
	OrdersInfo                []LibOrderOrderInfo
	FillableTakerAssetAmounts []*big.Int
	IsValidSignature          []bool
}, error) {
	ret := new(struct {
		OrdersInfo                []LibOrderOrderInfo
		FillableTakerAssetAmounts []*big.Int
		IsValidSignature          []bool
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "getOrderRelevantStates", orders, signatures)
	return *ret, err
}

// GetOrderRelevantStates is a free data retrieval call binding the contract method 0xe25cabf7.
//
// Solidity: function getOrderRelevantStates((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, bytes[] signatures) view returns((uint8,bytes32,uint256)[] ordersInfo, uint256[] fillableTakerAssetAmounts, bool[] isValidSignature)
func (_DevUtils *DevUtilsSession) GetOrderRelevantStates(orders []LibOrderOrder, signatures [][]byte) (struct {
	OrdersInfo                []LibOrderOrderInfo
	FillableTakerAssetAmounts []*big.Int
	IsValidSignature          []bool
}, error) {
	return _DevUtils.Contract.GetOrderRelevantStates(&_DevUtils.CallOpts, orders, signatures)
}

// GetOrderRelevantStates is a free data retrieval call binding the contract method 0xe25cabf7.
//
// Solidity: function getOrderRelevantStates((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, bytes[] signatures) view returns((uint8,bytes32,uint256)[] ordersInfo, uint256[] fillableTakerAssetAmounts, bool[] isValidSignature)
func (_DevUtils *DevUtilsCallerSession) GetOrderRelevantStates(orders []LibOrderOrder, signatures [][]byte) (struct {
	OrdersInfo                []LibOrderOrderInfo
	FillableTakerAssetAmounts []*big.Int
	IsValidSignature          []bool
}, error) {
	return _DevUtils.Contract.GetOrderRelevantStates(&_DevUtils.CallOpts, orders, signatures)
}

// GetTransactionHash is a free data retrieval call binding the contract method 0x79c9c426.
//
// Solidity: function getTransactionHash((uint256,uint256,uint256,address,bytes) transaction, uint256 chainId, address exchange) pure returns(bytes32 transactionHash)
func (_DevUtils *DevUtilsCaller) GetTransactionHash(opts *bind.CallOpts, transaction LibZeroExTransactionZeroExTransaction, chainId *big.Int, exchange common.Address) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "getTransactionHash", transaction, chainId, exchange)
	return *ret0, err
}

// GetTransactionHash is a free data retrieval call binding the contract method 0x79c9c426.
//
// Solidity: function getTransactionHash((uint256,uint256,uint256,address,bytes) transaction, uint256 chainId, address exchange) pure returns(bytes32 transactionHash)
func (_DevUtils *DevUtilsSession) GetTransactionHash(transaction LibZeroExTransactionZeroExTransaction, chainId *big.Int, exchange common.Address) ([32]byte, error) {
	return _DevUtils.Contract.GetTransactionHash(&_DevUtils.CallOpts, transaction, chainId, exchange)
}

// GetTransactionHash is a free data retrieval call binding the contract method 0x79c9c426.
//
// Solidity: function getTransactionHash((uint256,uint256,uint256,address,bytes) transaction, uint256 chainId, address exchange) pure returns(bytes32 transactionHash)
func (_DevUtils *DevUtilsCallerSession) GetTransactionHash(transaction LibZeroExTransactionZeroExTransaction, chainId *big.Int, exchange common.Address) ([32]byte, error) {
	return _DevUtils.Contract.GetTransactionHash(&_DevUtils.CallOpts, transaction, chainId, exchange)
}

// RevertIfInvalidAssetData is a free data retrieval call binding the contract method 0x46eb65cb.
//
// Solidity: function revertIfInvalidAssetData(bytes assetData) pure returns()
func (_DevUtils *DevUtilsCaller) RevertIfInvalidAssetData(opts *bind.CallOpts, assetData []byte) error {
	var ()
	out := &[]interface{}{}
	err := _DevUtils.contract.Call(opts, out, "revertIfInvalidAssetData", assetData)
	return err
}

// RevertIfInvalidAssetData is a free data retrieval call binding the contract method 0x46eb65cb.
//
// Solidity: function revertIfInvalidAssetData(bytes assetData) pure returns()
func (_DevUtils *DevUtilsSession) RevertIfInvalidAssetData(assetData []byte) error {
	return _DevUtils.Contract.RevertIfInvalidAssetData(&_DevUtils.CallOpts, assetData)
}

// RevertIfInvalidAssetData is a free data retrieval call binding the contract method 0x46eb65cb.
//
// Solidity: function revertIfInvalidAssetData(bytes assetData) pure returns()
func (_DevUtils *DevUtilsCallerSession) RevertIfInvalidAssetData(assetData []byte) error {
	return _DevUtils.Contract.RevertIfInvalidAssetData(&_DevUtils.CallOpts, assetData)
}

// StaticCallProxyAddress is a free data retrieval call binding the contract method 0x9baf2705.
//
// Solidity: function staticCallProxyAddress() view returns(address)
func (_DevUtils *DevUtilsCaller) StaticCallProxyAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "staticCallProxyAddress")
	return *ret0, err
}

// StaticCallProxyAddress is a free data retrieval call binding the contract method 0x9baf2705.
//
// Solidity: function staticCallProxyAddress() view returns(address)
func (_DevUtils *DevUtilsSession) StaticCallProxyAddress() (common.Address, error) {
	return _DevUtils.Contract.StaticCallProxyAddress(&_DevUtils.CallOpts)
}

// StaticCallProxyAddress is a free data retrieval call binding the contract method 0x9baf2705.
//
// Solidity: function staticCallProxyAddress() view returns(address)
func (_DevUtils *DevUtilsCallerSession) StaticCallProxyAddress() (common.Address, error) {
	return _DevUtils.Contract.StaticCallProxyAddress(&_DevUtils.CallOpts)
}

// GetAssetProxyAllowance is a paid mutator transaction binding the contract method 0xd186037f.
//
// Solidity: function getAssetProxyAllowance(address ownerAddress, bytes assetData) returns(uint256 allowance)
func (_DevUtils *DevUtilsTransactor) GetAssetProxyAllowance(opts *bind.TransactOpts, ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _DevUtils.contract.Transact(opts, "getAssetProxyAllowance", ownerAddress, assetData)
}

// GetAssetProxyAllowance is a paid mutator transaction binding the contract method 0xd186037f.
//
// Solidity: function getAssetProxyAllowance(address ownerAddress, bytes assetData) returns(uint256 allowance)
func (_DevUtils *DevUtilsSession) GetAssetProxyAllowance(ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _DevUtils.Contract.GetAssetProxyAllowance(&_DevUtils.TransactOpts, ownerAddress, assetData)
}

// GetAssetProxyAllowance is a paid mutator transaction binding the contract method 0xd186037f.
//
// Solidity: function getAssetProxyAllowance(address ownerAddress, bytes assetData) returns(uint256 allowance)
func (_DevUtils *DevUtilsTransactorSession) GetAssetProxyAllowance(ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _DevUtils.Contract.GetAssetProxyAllowance(&_DevUtils.TransactOpts, ownerAddress, assetData)
}

// GetBalance is a paid mutator transaction binding the contract method 0x7d727512.
//
// Solidity: function getBalance(address ownerAddress, bytes assetData) returns(uint256 balance)
func (_DevUtils *DevUtilsTransactor) GetBalance(opts *bind.TransactOpts, ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _DevUtils.contract.Transact(opts, "getBalance", ownerAddress, assetData)
}

// GetBalance is a paid mutator transaction binding the contract method 0x7d727512.
//
// Solidity: function getBalance(address ownerAddress, bytes assetData) returns(uint256 balance)
func (_DevUtils *DevUtilsSession) GetBalance(ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _DevUtils.Contract.GetBalance(&_DevUtils.TransactOpts, ownerAddress, assetData)
}

// GetBalance is a paid mutator transaction binding the contract method 0x7d727512.
//
// Solidity: function getBalance(address ownerAddress, bytes assetData) returns(uint256 balance)
func (_DevUtils *DevUtilsTransactorSession) GetBalance(ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _DevUtils.Contract.GetBalance(&_DevUtils.TransactOpts, ownerAddress, assetData)
}

// GetBalanceAndAssetProxyAllowance is a paid mutator transaction binding the contract method 0x0d7b7d76.
//
// Solidity: function getBalanceAndAssetProxyAllowance(address ownerAddress, bytes assetData) returns(uint256 balance, uint256 allowance)
func (_DevUtils *DevUtilsTransactor) GetBalanceAndAssetProxyAllowance(opts *bind.TransactOpts, ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _DevUtils.contract.Transact(opts, "getBalanceAndAssetProxyAllowance", ownerAddress, assetData)
}

// GetBalanceAndAssetProxyAllowance is a paid mutator transaction binding the contract method 0x0d7b7d76.
//
// Solidity: function getBalanceAndAssetProxyAllowance(address ownerAddress, bytes assetData) returns(uint256 balance, uint256 allowance)
func (_DevUtils *DevUtilsSession) GetBalanceAndAssetProxyAllowance(ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _DevUtils.Contract.GetBalanceAndAssetProxyAllowance(&_DevUtils.TransactOpts, ownerAddress, assetData)
}

// GetBalanceAndAssetProxyAllowance is a paid mutator transaction binding the contract method 0x0d7b7d76.
//
// Solidity: function getBalanceAndAssetProxyAllowance(address ownerAddress, bytes assetData) returns(uint256 balance, uint256 allowance)
func (_DevUtils *DevUtilsTransactorSession) GetBalanceAndAssetProxyAllowance(ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _DevUtils.Contract.GetBalanceAndAssetProxyAllowance(&_DevUtils.TransactOpts, ownerAddress, assetData)
}

// GetBatchAssetProxyAllowances is a paid mutator transaction binding the contract method 0x4dfdac20.
//
// Solidity: function getBatchAssetProxyAllowances(address ownerAddress, bytes[] assetData) returns(uint256[] allowances)
func (_DevUtils *DevUtilsTransactor) GetBatchAssetProxyAllowances(opts *bind.TransactOpts, ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _DevUtils.contract.Transact(opts, "getBatchAssetProxyAllowances", ownerAddress, assetData)
}

// GetBatchAssetProxyAllowances is a paid mutator transaction binding the contract method 0x4dfdac20.
//
// Solidity: function getBatchAssetProxyAllowances(address ownerAddress, bytes[] assetData) returns(uint256[] allowances)
func (_DevUtils *DevUtilsSession) GetBatchAssetProxyAllowances(ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _DevUtils.Contract.GetBatchAssetProxyAllowances(&_DevUtils.TransactOpts, ownerAddress, assetData)
}

// GetBatchAssetProxyAllowances is a paid mutator transaction binding the contract method 0x4dfdac20.
//
// Solidity: function getBatchAssetProxyAllowances(address ownerAddress, bytes[] assetData) returns(uint256[] allowances)
func (_DevUtils *DevUtilsTransactorSession) GetBatchAssetProxyAllowances(ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _DevUtils.Contract.GetBatchAssetProxyAllowances(&_DevUtils.TransactOpts, ownerAddress, assetData)
}

// GetBatchBalances is a paid mutator transaction binding the contract method 0xd001c5dc.
//
// Solidity: function getBatchBalances(address ownerAddress, bytes[] assetData) returns(uint256[] balances)
func (_DevUtils *DevUtilsTransactor) GetBatchBalances(opts *bind.TransactOpts, ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _DevUtils.contract.Transact(opts, "getBatchBalances", ownerAddress, assetData)
}

// GetBatchBalances is a paid mutator transaction binding the contract method 0xd001c5dc.
//
// Solidity: function getBatchBalances(address ownerAddress, bytes[] assetData) returns(uint256[] balances)
func (_DevUtils *DevUtilsSession) GetBatchBalances(ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _DevUtils.Contract.GetBatchBalances(&_DevUtils.TransactOpts, ownerAddress, assetData)
}

// GetBatchBalances is a paid mutator transaction binding the contract method 0xd001c5dc.
//
// Solidity: function getBatchBalances(address ownerAddress, bytes[] assetData) returns(uint256[] balances)
func (_DevUtils *DevUtilsTransactorSession) GetBatchBalances(ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _DevUtils.Contract.GetBatchBalances(&_DevUtils.TransactOpts, ownerAddress, assetData)
}

// GetBatchBalancesAndAssetProxyAllowances is a paid mutator transaction binding the contract method 0xe4e6e7da.
//
// Solidity: function getBatchBalancesAndAssetProxyAllowances(address ownerAddress, bytes[] assetData) returns(uint256[] balances, uint256[] allowances)
func (_DevUtils *DevUtilsTransactor) GetBatchBalancesAndAssetProxyAllowances(opts *bind.TransactOpts, ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _DevUtils.contract.Transact(opts, "getBatchBalancesAndAssetProxyAllowances", ownerAddress, assetData)
}

// GetBatchBalancesAndAssetProxyAllowances is a paid mutator transaction binding the contract method 0xe4e6e7da.
//
// Solidity: function getBatchBalancesAndAssetProxyAllowances(address ownerAddress, bytes[] assetData) returns(uint256[] balances, uint256[] allowances)
func (_DevUtils *DevUtilsSession) GetBatchBalancesAndAssetProxyAllowances(ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _DevUtils.Contract.GetBatchBalancesAndAssetProxyAllowances(&_DevUtils.TransactOpts, ownerAddress, assetData)
}

// GetBatchBalancesAndAssetProxyAllowances is a paid mutator transaction binding the contract method 0xe4e6e7da.
//
// Solidity: function getBatchBalancesAndAssetProxyAllowances(address ownerAddress, bytes[] assetData) returns(uint256[] balances, uint256[] allowances)
func (_DevUtils *DevUtilsTransactorSession) GetBatchBalancesAndAssetProxyAllowances(ownerAddress common.Address, assetData [][]byte) (*types.Transaction, error) {
	return _DevUtils.Contract.GetBatchBalancesAndAssetProxyAllowances(&_DevUtils.TransactOpts, ownerAddress, assetData)
}

// GetSimulatedOrderMakerTransferResults is a paid mutator transaction binding the contract method 0x7982653e.
//
// Solidity: function getSimulatedOrderMakerTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, address takerAddress, uint256 takerAssetFillAmount) returns(uint8 orderTransferResults)
func (_DevUtils *DevUtilsTransactor) GetSimulatedOrderMakerTransferResults(opts *bind.TransactOpts, order LibOrderOrder, takerAddress common.Address, takerAssetFillAmount *big.Int) (*types.Transaction, error) {
	return _DevUtils.contract.Transact(opts, "getSimulatedOrderMakerTransferResults", order, takerAddress, takerAssetFillAmount)
}

// GetSimulatedOrderMakerTransferResults is a paid mutator transaction binding the contract method 0x7982653e.
//
// Solidity: function getSimulatedOrderMakerTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, address takerAddress, uint256 takerAssetFillAmount) returns(uint8 orderTransferResults)
func (_DevUtils *DevUtilsSession) GetSimulatedOrderMakerTransferResults(order LibOrderOrder, takerAddress common.Address, takerAssetFillAmount *big.Int) (*types.Transaction, error) {
	return _DevUtils.Contract.GetSimulatedOrderMakerTransferResults(&_DevUtils.TransactOpts, order, takerAddress, takerAssetFillAmount)
}

// GetSimulatedOrderMakerTransferResults is a paid mutator transaction binding the contract method 0x7982653e.
//
// Solidity: function getSimulatedOrderMakerTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, address takerAddress, uint256 takerAssetFillAmount) returns(uint8 orderTransferResults)
func (_DevUtils *DevUtilsTransactorSession) GetSimulatedOrderMakerTransferResults(order LibOrderOrder, takerAddress common.Address, takerAssetFillAmount *big.Int) (*types.Transaction, error) {
	return _DevUtils.Contract.GetSimulatedOrderMakerTransferResults(&_DevUtils.TransactOpts, order, takerAddress, takerAssetFillAmount)
}

// GetSimulatedOrderTransferResults is a paid mutator transaction binding the contract method 0xd3637905.
//
// Solidity: function getSimulatedOrderTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, address takerAddress, uint256 takerAssetFillAmount) returns(uint8 orderTransferResults)
func (_DevUtils *DevUtilsTransactor) GetSimulatedOrderTransferResults(opts *bind.TransactOpts, order LibOrderOrder, takerAddress common.Address, takerAssetFillAmount *big.Int) (*types.Transaction, error) {
	return _DevUtils.contract.Transact(opts, "getSimulatedOrderTransferResults", order, takerAddress, takerAssetFillAmount)
}

// GetSimulatedOrderTransferResults is a paid mutator transaction binding the contract method 0xd3637905.
//
// Solidity: function getSimulatedOrderTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, address takerAddress, uint256 takerAssetFillAmount) returns(uint8 orderTransferResults)
func (_DevUtils *DevUtilsSession) GetSimulatedOrderTransferResults(order LibOrderOrder, takerAddress common.Address, takerAssetFillAmount *big.Int) (*types.Transaction, error) {
	return _DevUtils.Contract.GetSimulatedOrderTransferResults(&_DevUtils.TransactOpts, order, takerAddress, takerAssetFillAmount)
}

// GetSimulatedOrderTransferResults is a paid mutator transaction binding the contract method 0xd3637905.
//
// Solidity: function getSimulatedOrderTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, address takerAddress, uint256 takerAssetFillAmount) returns(uint8 orderTransferResults)
func (_DevUtils *DevUtilsTransactorSession) GetSimulatedOrderTransferResults(order LibOrderOrder, takerAddress common.Address, takerAssetFillAmount *big.Int) (*types.Transaction, error) {
	return _DevUtils.Contract.GetSimulatedOrderTransferResults(&_DevUtils.TransactOpts, order, takerAddress, takerAssetFillAmount)
}

// GetSimulatedOrdersTransferResults is a paid mutator transaction binding the contract method 0xa5cd62ba.
//
// Solidity: function getSimulatedOrdersTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, address[] takerAddresses, uint256[] takerAssetFillAmounts) returns(uint8[] orderTransferResults)
func (_DevUtils *DevUtilsTransactor) GetSimulatedOrdersTransferResults(opts *bind.TransactOpts, orders []LibOrderOrder, takerAddresses []common.Address, takerAssetFillAmounts []*big.Int) (*types.Transaction, error) {
	return _DevUtils.contract.Transact(opts, "getSimulatedOrdersTransferResults", orders, takerAddresses, takerAssetFillAmounts)
}

// GetSimulatedOrdersTransferResults is a paid mutator transaction binding the contract method 0xa5cd62ba.
//
// Solidity: function getSimulatedOrdersTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, address[] takerAddresses, uint256[] takerAssetFillAmounts) returns(uint8[] orderTransferResults)
func (_DevUtils *DevUtilsSession) GetSimulatedOrdersTransferResults(orders []LibOrderOrder, takerAddresses []common.Address, takerAssetFillAmounts []*big.Int) (*types.Transaction, error) {
	return _DevUtils.Contract.GetSimulatedOrdersTransferResults(&_DevUtils.TransactOpts, orders, takerAddresses, takerAssetFillAmounts)
}

// GetSimulatedOrdersTransferResults is a paid mutator transaction binding the contract method 0xa5cd62ba.
//
// Solidity: function getSimulatedOrdersTransferResults((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, address[] takerAddresses, uint256[] takerAssetFillAmounts) returns(uint8[] orderTransferResults)
func (_DevUtils *DevUtilsTransactorSession) GetSimulatedOrdersTransferResults(orders []LibOrderOrder, takerAddresses []common.Address, takerAssetFillAmounts []*big.Int) (*types.Transaction, error) {
	return _DevUtils.Contract.GetSimulatedOrdersTransferResults(&_DevUtils.TransactOpts, orders, takerAddresses, takerAssetFillAmounts)
}

// GetTransferableAssetAmount is a paid mutator transaction binding the contract method 0x2322cf76.
//
// Solidity: function getTransferableAssetAmount(address ownerAddress, bytes assetData) returns(uint256 transferableAssetAmount)
func (_DevUtils *DevUtilsTransactor) GetTransferableAssetAmount(opts *bind.TransactOpts, ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _DevUtils.contract.Transact(opts, "getTransferableAssetAmount", ownerAddress, assetData)
}

// GetTransferableAssetAmount is a paid mutator transaction binding the contract method 0x2322cf76.
//
// Solidity: function getTransferableAssetAmount(address ownerAddress, bytes assetData) returns(uint256 transferableAssetAmount)
func (_DevUtils *DevUtilsSession) GetTransferableAssetAmount(ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _DevUtils.Contract.GetTransferableAssetAmount(&_DevUtils.TransactOpts, ownerAddress, assetData)
}

// GetTransferableAssetAmount is a paid mutator transaction binding the contract method 0x2322cf76.
//
// Solidity: function getTransferableAssetAmount(address ownerAddress, bytes assetData) returns(uint256 transferableAssetAmount)
func (_DevUtils *DevUtilsTransactorSession) GetTransferableAssetAmount(ownerAddress common.Address, assetData []byte) (*types.Transaction, error) {
	return _DevUtils.Contract.GetTransferableAssetAmount(&_DevUtils.TransactOpts, ownerAddress, assetData)
}
