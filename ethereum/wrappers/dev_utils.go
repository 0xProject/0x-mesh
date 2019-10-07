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
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// DevUtilsABI is the input ABI used to generate the binding from.
const DevUtilsABI = "[{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"decodeOrderStatusError\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumLibOrder.OrderStatus\",\"name\":\"orderStatus\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC721AssetData\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getBalanceAndAssetProxyAllowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"decodeIncompleteFillError\",\"outputs\":[{\"internalType\":\"enumLibExchangeRichErrors.IncompleteFillErrorCode\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"expectedAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"actualAssetFillAmount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getTransferableAssetAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"transferableAssetAmount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"decodeAssetProxyTransferError\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"errorData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"decodeNegativeSpreadError\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"leftOrderHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"rightOrderHash\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"decodeAssetProxyDispatchError\",\"outputs\":[{\"internalType\":\"enumLibExchangeRichErrors.AssetProxyDispatchErrorCodes\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"decodeSignatureWalletError\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"errorData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"decodeFillError\",\"outputs\":[{\"internalType\":\"enumLibExchangeRichErrors.FillErrorCodes\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes[]\",\"name\":\"assetData\",\"type\":\"bytes[]\"}],\"name\":\"getBatchAssetProxyAllowances\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"allowances\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"encodeERC20AssetData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"decodeOrderEpochError\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"orderSenderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"currentEpoch\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"transactionData\",\"type\":\"bytes\"}],\"name\":\"decodeZeroExTransactionData\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"functionName\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"decodeAssetProxyExistsError\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"assetProxyAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"decodeSignatureValidatorNotApprovedError\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC20AssetData\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"decodeSignatureError\",\"outputs\":[{\"internalType\":\"enumLibExchangeRichErrors.SignatureErrorCodes\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC1155AssetData\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenValues\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"callbackData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"addresses\",\"type\":\"address[]\"}],\"name\":\"getEthBalances\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"address[]\",\"name\":\"takerAddresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"}],\"name\":\"getSimulatedOrdersTransferResults\",\"outputs\":[{\"internalType\":\"enumOrderTransferSimulationUtils.OrderTransferResults[]\",\"name\":\"orderTransferResults\",\"type\":\"uint8[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"encodeERC721AssetData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"decodeEIP1271SignatureError\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"verifyingContractAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"errorData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenIds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenValues\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"callbackData\",\"type\":\"bytes\"}],\"name\":\"encodeERC1155AssetData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeMultiAssetData\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"nestedAssetData\",\"type\":\"bytes[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"decodeTransactionExecutionError\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"transactionHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"errorData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"decodeTransactionError\",\"outputs\":[{\"internalType\":\"enumLibExchangeRichErrors.TransactionErrorCodes\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"transactionHash\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes[]\",\"name\":\"assetData\",\"type\":\"bytes[]\"}],\"name\":\"getBatchBalances\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getAssetProxyAllowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"}],\"name\":\"getSimulatedOrderTransferResults\",\"outputs\":[{\"internalType\":\"enumOrderTransferSimulationUtils.OrderTransferResults\",\"name\":\"orderTransferResults\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"nestedAssetData\",\"type\":\"bytes[]\"}],\"name\":\"encodeMultiAssetData\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"getOrderRelevantStates\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"internalType\":\"structLibOrder.OrderInfo[]\",\"name\":\"ordersInfo\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"fillableTakerAssetAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bool[]\",\"name\":\"isValidSignature\",\"type\":\"bool[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes[]\",\"name\":\"assetData\",\"type\":\"bytes[]\"}],\"name\":\"getBatchBalancesAndAssetProxyAllowances\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"balances\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"allowances\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"getOrderRelevantState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"internalType\":\"structLibOrder.OrderInfo\",\"name\":\"orderInfo\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"fillableTakerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isValidSignature\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"encoded\",\"type\":\"bytes\"}],\"name\":\"decodeExchangeInvalidContextError\",\"outputs\":[{\"internalType\":\"enumLibExchangeRichErrors.ExchangeContextErrorCodes\",\"name\":\"errorCode\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"contextAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_exchange\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

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

// Struct0 is an auto generated low-level Go binding around an user-defined struct.
type Struct0 struct {
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

// Struct1 is an auto generated low-level Go binding around an user-defined struct.
type Struct1 struct {
	OrderStatus                 uint8
	OrderHash                   [32]byte
	OrderTakerAssetFilledAmount *big.Int
}

// DecodeAssetProxyDispatchError is a free data retrieval call binding the contract method 0x32aae3ad.
//
// Solidity: function decodeAssetProxyDispatchError(bytes encoded) constant returns(uint8 errorCode, bytes32 orderHash, bytes assetData)
func (_DevUtils *DevUtilsCaller) DecodeAssetProxyDispatchError(opts *bind.CallOpts, encoded []byte) (struct {
	ErrorCode uint8
	OrderHash [32]byte
	AssetData []byte
}, error) {
	ret := new(struct {
		ErrorCode uint8
		OrderHash [32]byte
		AssetData []byte
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeAssetProxyDispatchError", encoded)
	return *ret, err
}

// DecodeAssetProxyDispatchError is a free data retrieval call binding the contract method 0x32aae3ad.
//
// Solidity: function decodeAssetProxyDispatchError(bytes encoded) constant returns(uint8 errorCode, bytes32 orderHash, bytes assetData)
func (_DevUtils *DevUtilsSession) DecodeAssetProxyDispatchError(encoded []byte) (struct {
	ErrorCode uint8
	OrderHash [32]byte
	AssetData []byte
}, error) {
	return _DevUtils.Contract.DecodeAssetProxyDispatchError(&_DevUtils.CallOpts, encoded)
}

// DecodeAssetProxyDispatchError is a free data retrieval call binding the contract method 0x32aae3ad.
//
// Solidity: function decodeAssetProxyDispatchError(bytes encoded) constant returns(uint8 errorCode, bytes32 orderHash, bytes assetData)
func (_DevUtils *DevUtilsCallerSession) DecodeAssetProxyDispatchError(encoded []byte) (struct {
	ErrorCode uint8
	OrderHash [32]byte
	AssetData []byte
}, error) {
	return _DevUtils.Contract.DecodeAssetProxyDispatchError(&_DevUtils.CallOpts, encoded)
}

// DecodeAssetProxyExistsError is a free data retrieval call binding the contract method 0x7914b2ec.
//
// Solidity: function decodeAssetProxyExistsError(bytes encoded) constant returns(bytes4 assetProxyId, address assetProxyAddress)
func (_DevUtils *DevUtilsCaller) DecodeAssetProxyExistsError(opts *bind.CallOpts, encoded []byte) (struct {
	AssetProxyId      [4]byte
	AssetProxyAddress common.Address
}, error) {
	ret := new(struct {
		AssetProxyId      [4]byte
		AssetProxyAddress common.Address
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeAssetProxyExistsError", encoded)
	return *ret, err
}

// DecodeAssetProxyExistsError is a free data retrieval call binding the contract method 0x7914b2ec.
//
// Solidity: function decodeAssetProxyExistsError(bytes encoded) constant returns(bytes4 assetProxyId, address assetProxyAddress)
func (_DevUtils *DevUtilsSession) DecodeAssetProxyExistsError(encoded []byte) (struct {
	AssetProxyId      [4]byte
	AssetProxyAddress common.Address
}, error) {
	return _DevUtils.Contract.DecodeAssetProxyExistsError(&_DevUtils.CallOpts, encoded)
}

// DecodeAssetProxyExistsError is a free data retrieval call binding the contract method 0x7914b2ec.
//
// Solidity: function decodeAssetProxyExistsError(bytes encoded) constant returns(bytes4 assetProxyId, address assetProxyAddress)
func (_DevUtils *DevUtilsCallerSession) DecodeAssetProxyExistsError(encoded []byte) (struct {
	AssetProxyId      [4]byte
	AssetProxyAddress common.Address
}, error) {
	return _DevUtils.Contract.DecodeAssetProxyExistsError(&_DevUtils.CallOpts, encoded)
}

// DecodeAssetProxyTransferError is a free data retrieval call binding the contract method 0x314853ff.
//
// Solidity: function decodeAssetProxyTransferError(bytes encoded) constant returns(bytes32 orderHash, bytes assetData, bytes errorData)
func (_DevUtils *DevUtilsCaller) DecodeAssetProxyTransferError(opts *bind.CallOpts, encoded []byte) (struct {
	OrderHash [32]byte
	AssetData []byte
	ErrorData []byte
}, error) {
	ret := new(struct {
		OrderHash [32]byte
		AssetData []byte
		ErrorData []byte
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeAssetProxyTransferError", encoded)
	return *ret, err
}

// DecodeAssetProxyTransferError is a free data retrieval call binding the contract method 0x314853ff.
//
// Solidity: function decodeAssetProxyTransferError(bytes encoded) constant returns(bytes32 orderHash, bytes assetData, bytes errorData)
func (_DevUtils *DevUtilsSession) DecodeAssetProxyTransferError(encoded []byte) (struct {
	OrderHash [32]byte
	AssetData []byte
	ErrorData []byte
}, error) {
	return _DevUtils.Contract.DecodeAssetProxyTransferError(&_DevUtils.CallOpts, encoded)
}

// DecodeAssetProxyTransferError is a free data retrieval call binding the contract method 0x314853ff.
//
// Solidity: function decodeAssetProxyTransferError(bytes encoded) constant returns(bytes32 orderHash, bytes assetData, bytes errorData)
func (_DevUtils *DevUtilsCallerSession) DecodeAssetProxyTransferError(encoded []byte) (struct {
	OrderHash [32]byte
	AssetData []byte
	ErrorData []byte
}, error) {
	return _DevUtils.Contract.DecodeAssetProxyTransferError(&_DevUtils.CallOpts, encoded)
}

// DecodeEIP1271SignatureError is a free data retrieval call binding the contract method 0xacaedc74.
//
// Solidity: function decodeEIP1271SignatureError(bytes encoded) constant returns(address verifyingContractAddress, bytes data, bytes signature, bytes errorData)
func (_DevUtils *DevUtilsCaller) DecodeEIP1271SignatureError(opts *bind.CallOpts, encoded []byte) (struct {
	VerifyingContractAddress common.Address
	Data                     []byte
	Signature                []byte
	ErrorData                []byte
}, error) {
	ret := new(struct {
		VerifyingContractAddress common.Address
		Data                     []byte
		Signature                []byte
		ErrorData                []byte
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeEIP1271SignatureError", encoded)
	return *ret, err
}

// DecodeEIP1271SignatureError is a free data retrieval call binding the contract method 0xacaedc74.
//
// Solidity: function decodeEIP1271SignatureError(bytes encoded) constant returns(address verifyingContractAddress, bytes data, bytes signature, bytes errorData)
func (_DevUtils *DevUtilsSession) DecodeEIP1271SignatureError(encoded []byte) (struct {
	VerifyingContractAddress common.Address
	Data                     []byte
	Signature                []byte
	ErrorData                []byte
}, error) {
	return _DevUtils.Contract.DecodeEIP1271SignatureError(&_DevUtils.CallOpts, encoded)
}

// DecodeEIP1271SignatureError is a free data retrieval call binding the contract method 0xacaedc74.
//
// Solidity: function decodeEIP1271SignatureError(bytes encoded) constant returns(address verifyingContractAddress, bytes data, bytes signature, bytes errorData)
func (_DevUtils *DevUtilsCallerSession) DecodeEIP1271SignatureError(encoded []byte) (struct {
	VerifyingContractAddress common.Address
	Data                     []byte
	Signature                []byte
	ErrorData                []byte
}, error) {
	return _DevUtils.Contract.DecodeEIP1271SignatureError(&_DevUtils.CallOpts, encoded)
}

// DecodeERC1155AssetData is a free data retrieval call binding the contract method 0x9eadc835.
//
// Solidity: function decodeERC1155AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData)
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
// Solidity: function decodeERC1155AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData)
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
// Solidity: function decodeERC1155AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData)
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
// Solidity: function decodeERC20AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress)
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
// Solidity: function decodeERC20AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress)
func (_DevUtils *DevUtilsSession) DecodeERC20AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
}, error) {
	return _DevUtils.Contract.DecodeERC20AssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeERC20AssetData is a free data retrieval call binding the contract method 0x8f4ce479.
//
// Solidity: function decodeERC20AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress)
func (_DevUtils *DevUtilsCallerSession) DecodeERC20AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
}, error) {
	return _DevUtils.Contract.DecodeERC20AssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeERC721AssetData is a free data retrieval call binding the contract method 0x04a5618a.
//
// Solidity: function decodeERC721AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress, uint256 tokenId)
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
// Solidity: function decodeERC721AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress, uint256 tokenId)
func (_DevUtils *DevUtilsSession) DecodeERC721AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
	TokenId      *big.Int
}, error) {
	return _DevUtils.Contract.DecodeERC721AssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeERC721AssetData is a free data retrieval call binding the contract method 0x04a5618a.
//
// Solidity: function decodeERC721AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress, uint256 tokenId)
func (_DevUtils *DevUtilsCallerSession) DecodeERC721AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
	TokenId      *big.Int
}, error) {
	return _DevUtils.Contract.DecodeERC721AssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeExchangeInvalidContextError is a free data retrieval call binding the contract method 0xee4f5a94.
//
// Solidity: function decodeExchangeInvalidContextError(bytes encoded) constant returns(uint8 errorCode, bytes32 orderHash, address contextAddress)
func (_DevUtils *DevUtilsCaller) DecodeExchangeInvalidContextError(opts *bind.CallOpts, encoded []byte) (struct {
	ErrorCode      uint8
	OrderHash      [32]byte
	ContextAddress common.Address
}, error) {
	ret := new(struct {
		ErrorCode      uint8
		OrderHash      [32]byte
		ContextAddress common.Address
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeExchangeInvalidContextError", encoded)
	return *ret, err
}

// DecodeExchangeInvalidContextError is a free data retrieval call binding the contract method 0xee4f5a94.
//
// Solidity: function decodeExchangeInvalidContextError(bytes encoded) constant returns(uint8 errorCode, bytes32 orderHash, address contextAddress)
func (_DevUtils *DevUtilsSession) DecodeExchangeInvalidContextError(encoded []byte) (struct {
	ErrorCode      uint8
	OrderHash      [32]byte
	ContextAddress common.Address
}, error) {
	return _DevUtils.Contract.DecodeExchangeInvalidContextError(&_DevUtils.CallOpts, encoded)
}

// DecodeExchangeInvalidContextError is a free data retrieval call binding the contract method 0xee4f5a94.
//
// Solidity: function decodeExchangeInvalidContextError(bytes encoded) constant returns(uint8 errorCode, bytes32 orderHash, address contextAddress)
func (_DevUtils *DevUtilsCallerSession) DecodeExchangeInvalidContextError(encoded []byte) (struct {
	ErrorCode      uint8
	OrderHash      [32]byte
	ContextAddress common.Address
}, error) {
	return _DevUtils.Contract.DecodeExchangeInvalidContextError(&_DevUtils.CallOpts, encoded)
}

// DecodeFillError is a free data retrieval call binding the contract method 0x459be5e2.
//
// Solidity: function decodeFillError(bytes encoded) constant returns(uint8 errorCode, bytes32 orderHash)
func (_DevUtils *DevUtilsCaller) DecodeFillError(opts *bind.CallOpts, encoded []byte) (struct {
	ErrorCode uint8
	OrderHash [32]byte
}, error) {
	ret := new(struct {
		ErrorCode uint8
		OrderHash [32]byte
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeFillError", encoded)
	return *ret, err
}

// DecodeFillError is a free data retrieval call binding the contract method 0x459be5e2.
//
// Solidity: function decodeFillError(bytes encoded) constant returns(uint8 errorCode, bytes32 orderHash)
func (_DevUtils *DevUtilsSession) DecodeFillError(encoded []byte) (struct {
	ErrorCode uint8
	OrderHash [32]byte
}, error) {
	return _DevUtils.Contract.DecodeFillError(&_DevUtils.CallOpts, encoded)
}

// DecodeFillError is a free data retrieval call binding the contract method 0x459be5e2.
//
// Solidity: function decodeFillError(bytes encoded) constant returns(uint8 errorCode, bytes32 orderHash)
func (_DevUtils *DevUtilsCallerSession) DecodeFillError(encoded []byte) (struct {
	ErrorCode uint8
	OrderHash [32]byte
}, error) {
	return _DevUtils.Contract.DecodeFillError(&_DevUtils.CallOpts, encoded)
}

// DecodeIncompleteFillError is a free data retrieval call binding the contract method 0x165979e1.
//
// Solidity: function decodeIncompleteFillError(bytes encoded) constant returns(uint8 errorCode, uint256 expectedAssetFillAmount, uint256 actualAssetFillAmount)
func (_DevUtils *DevUtilsCaller) DecodeIncompleteFillError(opts *bind.CallOpts, encoded []byte) (struct {
	ErrorCode               uint8
	ExpectedAssetFillAmount *big.Int
	ActualAssetFillAmount   *big.Int
}, error) {
	ret := new(struct {
		ErrorCode               uint8
		ExpectedAssetFillAmount *big.Int
		ActualAssetFillAmount   *big.Int
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeIncompleteFillError", encoded)
	return *ret, err
}

// DecodeIncompleteFillError is a free data retrieval call binding the contract method 0x165979e1.
//
// Solidity: function decodeIncompleteFillError(bytes encoded) constant returns(uint8 errorCode, uint256 expectedAssetFillAmount, uint256 actualAssetFillAmount)
func (_DevUtils *DevUtilsSession) DecodeIncompleteFillError(encoded []byte) (struct {
	ErrorCode               uint8
	ExpectedAssetFillAmount *big.Int
	ActualAssetFillAmount   *big.Int
}, error) {
	return _DevUtils.Contract.DecodeIncompleteFillError(&_DevUtils.CallOpts, encoded)
}

// DecodeIncompleteFillError is a free data retrieval call binding the contract method 0x165979e1.
//
// Solidity: function decodeIncompleteFillError(bytes encoded) constant returns(uint8 errorCode, uint256 expectedAssetFillAmount, uint256 actualAssetFillAmount)
func (_DevUtils *DevUtilsCallerSession) DecodeIncompleteFillError(encoded []byte) (struct {
	ErrorCode               uint8
	ExpectedAssetFillAmount *big.Int
	ActualAssetFillAmount   *big.Int
}, error) {
	return _DevUtils.Contract.DecodeIncompleteFillError(&_DevUtils.CallOpts, encoded)
}

// DecodeMultiAssetData is a free data retrieval call binding the contract method 0xbbb2dcf6.
//
// Solidity: function decodeMultiAssetData(bytes assetData) constant returns(bytes4 assetProxyId, uint256[] amounts, bytes[] nestedAssetData)
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
// Solidity: function decodeMultiAssetData(bytes assetData) constant returns(bytes4 assetProxyId, uint256[] amounts, bytes[] nestedAssetData)
func (_DevUtils *DevUtilsSession) DecodeMultiAssetData(assetData []byte) (struct {
	AssetProxyId    [4]byte
	Amounts         []*big.Int
	NestedAssetData [][]byte
}, error) {
	return _DevUtils.Contract.DecodeMultiAssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeMultiAssetData is a free data retrieval call binding the contract method 0xbbb2dcf6.
//
// Solidity: function decodeMultiAssetData(bytes assetData) constant returns(bytes4 assetProxyId, uint256[] amounts, bytes[] nestedAssetData)
func (_DevUtils *DevUtilsCallerSession) DecodeMultiAssetData(assetData []byte) (struct {
	AssetProxyId    [4]byte
	Amounts         []*big.Int
	NestedAssetData [][]byte
}, error) {
	return _DevUtils.Contract.DecodeMultiAssetData(&_DevUtils.CallOpts, assetData)
}

// DecodeNegativeSpreadError is a free data retrieval call binding the contract method 0x327d3054.
//
// Solidity: function decodeNegativeSpreadError(bytes encoded) constant returns(bytes32 leftOrderHash, bytes32 rightOrderHash)
func (_DevUtils *DevUtilsCaller) DecodeNegativeSpreadError(opts *bind.CallOpts, encoded []byte) (struct {
	LeftOrderHash  [32]byte
	RightOrderHash [32]byte
}, error) {
	ret := new(struct {
		LeftOrderHash  [32]byte
		RightOrderHash [32]byte
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeNegativeSpreadError", encoded)
	return *ret, err
}

// DecodeNegativeSpreadError is a free data retrieval call binding the contract method 0x327d3054.
//
// Solidity: function decodeNegativeSpreadError(bytes encoded) constant returns(bytes32 leftOrderHash, bytes32 rightOrderHash)
func (_DevUtils *DevUtilsSession) DecodeNegativeSpreadError(encoded []byte) (struct {
	LeftOrderHash  [32]byte
	RightOrderHash [32]byte
}, error) {
	return _DevUtils.Contract.DecodeNegativeSpreadError(&_DevUtils.CallOpts, encoded)
}

// DecodeNegativeSpreadError is a free data retrieval call binding the contract method 0x327d3054.
//
// Solidity: function decodeNegativeSpreadError(bytes encoded) constant returns(bytes32 leftOrderHash, bytes32 rightOrderHash)
func (_DevUtils *DevUtilsCallerSession) DecodeNegativeSpreadError(encoded []byte) (struct {
	LeftOrderHash  [32]byte
	RightOrderHash [32]byte
}, error) {
	return _DevUtils.Contract.DecodeNegativeSpreadError(&_DevUtils.CallOpts, encoded)
}

// DecodeOrderEpochError is a free data retrieval call binding the contract method 0x65129042.
//
// Solidity: function decodeOrderEpochError(bytes encoded) constant returns(address makerAddress, address orderSenderAddress, uint256 currentEpoch)
func (_DevUtils *DevUtilsCaller) DecodeOrderEpochError(opts *bind.CallOpts, encoded []byte) (struct {
	MakerAddress       common.Address
	OrderSenderAddress common.Address
	CurrentEpoch       *big.Int
}, error) {
	ret := new(struct {
		MakerAddress       common.Address
		OrderSenderAddress common.Address
		CurrentEpoch       *big.Int
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeOrderEpochError", encoded)
	return *ret, err
}

// DecodeOrderEpochError is a free data retrieval call binding the contract method 0x65129042.
//
// Solidity: function decodeOrderEpochError(bytes encoded) constant returns(address makerAddress, address orderSenderAddress, uint256 currentEpoch)
func (_DevUtils *DevUtilsSession) DecodeOrderEpochError(encoded []byte) (struct {
	MakerAddress       common.Address
	OrderSenderAddress common.Address
	CurrentEpoch       *big.Int
}, error) {
	return _DevUtils.Contract.DecodeOrderEpochError(&_DevUtils.CallOpts, encoded)
}

// DecodeOrderEpochError is a free data retrieval call binding the contract method 0x65129042.
//
// Solidity: function decodeOrderEpochError(bytes encoded) constant returns(address makerAddress, address orderSenderAddress, uint256 currentEpoch)
func (_DevUtils *DevUtilsCallerSession) DecodeOrderEpochError(encoded []byte) (struct {
	MakerAddress       common.Address
	OrderSenderAddress common.Address
	CurrentEpoch       *big.Int
}, error) {
	return _DevUtils.Contract.DecodeOrderEpochError(&_DevUtils.CallOpts, encoded)
}

// DecodeOrderStatusError is a free data retrieval call binding the contract method 0x02d0aec3.
//
// Solidity: function decodeOrderStatusError(bytes encoded) constant returns(bytes32 orderHash, uint8 orderStatus)
func (_DevUtils *DevUtilsCaller) DecodeOrderStatusError(opts *bind.CallOpts, encoded []byte) (struct {
	OrderHash   [32]byte
	OrderStatus uint8
}, error) {
	ret := new(struct {
		OrderHash   [32]byte
		OrderStatus uint8
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeOrderStatusError", encoded)
	return *ret, err
}

// DecodeOrderStatusError is a free data retrieval call binding the contract method 0x02d0aec3.
//
// Solidity: function decodeOrderStatusError(bytes encoded) constant returns(bytes32 orderHash, uint8 orderStatus)
func (_DevUtils *DevUtilsSession) DecodeOrderStatusError(encoded []byte) (struct {
	OrderHash   [32]byte
	OrderStatus uint8
}, error) {
	return _DevUtils.Contract.DecodeOrderStatusError(&_DevUtils.CallOpts, encoded)
}

// DecodeOrderStatusError is a free data retrieval call binding the contract method 0x02d0aec3.
//
// Solidity: function decodeOrderStatusError(bytes encoded) constant returns(bytes32 orderHash, uint8 orderStatus)
func (_DevUtils *DevUtilsCallerSession) DecodeOrderStatusError(encoded []byte) (struct {
	OrderHash   [32]byte
	OrderStatus uint8
}, error) {
	return _DevUtils.Contract.DecodeOrderStatusError(&_DevUtils.CallOpts, encoded)
}

// DecodeSignatureError is a free data retrieval call binding the contract method 0x9a7e7526.
//
// Solidity: function decodeSignatureError(bytes encoded) constant returns(uint8 errorCode, bytes32 hash, address signerAddress, bytes signature)
func (_DevUtils *DevUtilsCaller) DecodeSignatureError(opts *bind.CallOpts, encoded []byte) (struct {
	ErrorCode     uint8
	Hash          [32]byte
	SignerAddress common.Address
	Signature     []byte
}, error) {
	ret := new(struct {
		ErrorCode     uint8
		Hash          [32]byte
		SignerAddress common.Address
		Signature     []byte
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeSignatureError", encoded)
	return *ret, err
}

// DecodeSignatureError is a free data retrieval call binding the contract method 0x9a7e7526.
//
// Solidity: function decodeSignatureError(bytes encoded) constant returns(uint8 errorCode, bytes32 hash, address signerAddress, bytes signature)
func (_DevUtils *DevUtilsSession) DecodeSignatureError(encoded []byte) (struct {
	ErrorCode     uint8
	Hash          [32]byte
	SignerAddress common.Address
	Signature     []byte
}, error) {
	return _DevUtils.Contract.DecodeSignatureError(&_DevUtils.CallOpts, encoded)
}

// DecodeSignatureError is a free data retrieval call binding the contract method 0x9a7e7526.
//
// Solidity: function decodeSignatureError(bytes encoded) constant returns(uint8 errorCode, bytes32 hash, address signerAddress, bytes signature)
func (_DevUtils *DevUtilsCallerSession) DecodeSignatureError(encoded []byte) (struct {
	ErrorCode     uint8
	Hash          [32]byte
	SignerAddress common.Address
	Signature     []byte
}, error) {
	return _DevUtils.Contract.DecodeSignatureError(&_DevUtils.CallOpts, encoded)
}

// DecodeSignatureValidatorNotApprovedError is a free data retrieval call binding the contract method 0x7b66ad34.
//
// Solidity: function decodeSignatureValidatorNotApprovedError(bytes encoded) constant returns(address signerAddress, address validatorAddress)
func (_DevUtils *DevUtilsCaller) DecodeSignatureValidatorNotApprovedError(opts *bind.CallOpts, encoded []byte) (struct {
	SignerAddress    common.Address
	ValidatorAddress common.Address
}, error) {
	ret := new(struct {
		SignerAddress    common.Address
		ValidatorAddress common.Address
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeSignatureValidatorNotApprovedError", encoded)
	return *ret, err
}

// DecodeSignatureValidatorNotApprovedError is a free data retrieval call binding the contract method 0x7b66ad34.
//
// Solidity: function decodeSignatureValidatorNotApprovedError(bytes encoded) constant returns(address signerAddress, address validatorAddress)
func (_DevUtils *DevUtilsSession) DecodeSignatureValidatorNotApprovedError(encoded []byte) (struct {
	SignerAddress    common.Address
	ValidatorAddress common.Address
}, error) {
	return _DevUtils.Contract.DecodeSignatureValidatorNotApprovedError(&_DevUtils.CallOpts, encoded)
}

// DecodeSignatureValidatorNotApprovedError is a free data retrieval call binding the contract method 0x7b66ad34.
//
// Solidity: function decodeSignatureValidatorNotApprovedError(bytes encoded) constant returns(address signerAddress, address validatorAddress)
func (_DevUtils *DevUtilsCallerSession) DecodeSignatureValidatorNotApprovedError(encoded []byte) (struct {
	SignerAddress    common.Address
	ValidatorAddress common.Address
}, error) {
	return _DevUtils.Contract.DecodeSignatureValidatorNotApprovedError(&_DevUtils.CallOpts, encoded)
}

// DecodeSignatureWalletError is a free data retrieval call binding the contract method 0x3db6dc61.
//
// Solidity: function decodeSignatureWalletError(bytes encoded) constant returns(bytes32 hash, address signerAddress, bytes signature, bytes errorData)
func (_DevUtils *DevUtilsCaller) DecodeSignatureWalletError(opts *bind.CallOpts, encoded []byte) (struct {
	Hash          [32]byte
	SignerAddress common.Address
	Signature     []byte
	ErrorData     []byte
}, error) {
	ret := new(struct {
		Hash          [32]byte
		SignerAddress common.Address
		Signature     []byte
		ErrorData     []byte
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeSignatureWalletError", encoded)
	return *ret, err
}

// DecodeSignatureWalletError is a free data retrieval call binding the contract method 0x3db6dc61.
//
// Solidity: function decodeSignatureWalletError(bytes encoded) constant returns(bytes32 hash, address signerAddress, bytes signature, bytes errorData)
func (_DevUtils *DevUtilsSession) DecodeSignatureWalletError(encoded []byte) (struct {
	Hash          [32]byte
	SignerAddress common.Address
	Signature     []byte
	ErrorData     []byte
}, error) {
	return _DevUtils.Contract.DecodeSignatureWalletError(&_DevUtils.CallOpts, encoded)
}

// DecodeSignatureWalletError is a free data retrieval call binding the contract method 0x3db6dc61.
//
// Solidity: function decodeSignatureWalletError(bytes encoded) constant returns(bytes32 hash, address signerAddress, bytes signature, bytes errorData)
func (_DevUtils *DevUtilsCallerSession) DecodeSignatureWalletError(encoded []byte) (struct {
	Hash          [32]byte
	SignerAddress common.Address
	Signature     []byte
	ErrorData     []byte
}, error) {
	return _DevUtils.Contract.DecodeSignatureWalletError(&_DevUtils.CallOpts, encoded)
}

// DecodeTransactionError is a free data retrieval call binding the contract method 0xcafd3a07.
//
// Solidity: function decodeTransactionError(bytes encoded) constant returns(uint8 errorCode, bytes32 transactionHash)
func (_DevUtils *DevUtilsCaller) DecodeTransactionError(opts *bind.CallOpts, encoded []byte) (struct {
	ErrorCode       uint8
	TransactionHash [32]byte
}, error) {
	ret := new(struct {
		ErrorCode       uint8
		TransactionHash [32]byte
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeTransactionError", encoded)
	return *ret, err
}

// DecodeTransactionError is a free data retrieval call binding the contract method 0xcafd3a07.
//
// Solidity: function decodeTransactionError(bytes encoded) constant returns(uint8 errorCode, bytes32 transactionHash)
func (_DevUtils *DevUtilsSession) DecodeTransactionError(encoded []byte) (struct {
	ErrorCode       uint8
	TransactionHash [32]byte
}, error) {
	return _DevUtils.Contract.DecodeTransactionError(&_DevUtils.CallOpts, encoded)
}

// DecodeTransactionError is a free data retrieval call binding the contract method 0xcafd3a07.
//
// Solidity: function decodeTransactionError(bytes encoded) constant returns(uint8 errorCode, bytes32 transactionHash)
func (_DevUtils *DevUtilsCallerSession) DecodeTransactionError(encoded []byte) (struct {
	ErrorCode       uint8
	TransactionHash [32]byte
}, error) {
	return _DevUtils.Contract.DecodeTransactionError(&_DevUtils.CallOpts, encoded)
}

// DecodeTransactionExecutionError is a free data retrieval call binding the contract method 0xbc03f964.
//
// Solidity: function decodeTransactionExecutionError(bytes encoded) constant returns(bytes32 transactionHash, bytes errorData)
func (_DevUtils *DevUtilsCaller) DecodeTransactionExecutionError(opts *bind.CallOpts, encoded []byte) (struct {
	TransactionHash [32]byte
	ErrorData       []byte
}, error) {
	ret := new(struct {
		TransactionHash [32]byte
		ErrorData       []byte
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeTransactionExecutionError", encoded)
	return *ret, err
}

// DecodeTransactionExecutionError is a free data retrieval call binding the contract method 0xbc03f964.
//
// Solidity: function decodeTransactionExecutionError(bytes encoded) constant returns(bytes32 transactionHash, bytes errorData)
func (_DevUtils *DevUtilsSession) DecodeTransactionExecutionError(encoded []byte) (struct {
	TransactionHash [32]byte
	ErrorData       []byte
}, error) {
	return _DevUtils.Contract.DecodeTransactionExecutionError(&_DevUtils.CallOpts, encoded)
}

// DecodeTransactionExecutionError is a free data retrieval call binding the contract method 0xbc03f964.
//
// Solidity: function decodeTransactionExecutionError(bytes encoded) constant returns(bytes32 transactionHash, bytes errorData)
func (_DevUtils *DevUtilsCallerSession) DecodeTransactionExecutionError(encoded []byte) (struct {
	TransactionHash [32]byte
	ErrorData       []byte
}, error) {
	return _DevUtils.Contract.DecodeTransactionExecutionError(&_DevUtils.CallOpts, encoded)
}

// DecodeZeroExTransactionData is a free data retrieval call binding the contract method 0x6f83188e.
//
// Solidity: function decodeZeroExTransactionData(bytes transactionData) constant returns(string functionName, []Struct0 orders, uint256[] takerAssetFillAmounts, bytes[] signatures)
func (_DevUtils *DevUtilsCaller) DecodeZeroExTransactionData(opts *bind.CallOpts, transactionData []byte) (struct {
	FunctionName          string
	Orders                []Struct0
	TakerAssetFillAmounts []*big.Int
	Signatures            [][]byte
}, error) {
	ret := new(struct {
		FunctionName          string
		Orders                []Struct0
		TakerAssetFillAmounts []*big.Int
		Signatures            [][]byte
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "decodeZeroExTransactionData", transactionData)
	return *ret, err
}

// DecodeZeroExTransactionData is a free data retrieval call binding the contract method 0x6f83188e.
//
// Solidity: function decodeZeroExTransactionData(bytes transactionData) constant returns(string functionName, []Struct0 orders, uint256[] takerAssetFillAmounts, bytes[] signatures)
func (_DevUtils *DevUtilsSession) DecodeZeroExTransactionData(transactionData []byte) (struct {
	FunctionName          string
	Orders                []Struct0
	TakerAssetFillAmounts []*big.Int
	Signatures            [][]byte
}, error) {
	return _DevUtils.Contract.DecodeZeroExTransactionData(&_DevUtils.CallOpts, transactionData)
}

// DecodeZeroExTransactionData is a free data retrieval call binding the contract method 0x6f83188e.
//
// Solidity: function decodeZeroExTransactionData(bytes transactionData) constant returns(string functionName, []Struct0 orders, uint256[] takerAssetFillAmounts, bytes[] signatures)
func (_DevUtils *DevUtilsCallerSession) DecodeZeroExTransactionData(transactionData []byte) (struct {
	FunctionName          string
	Orders                []Struct0
	TakerAssetFillAmounts []*big.Int
	Signatures            [][]byte
}, error) {
	return _DevUtils.Contract.DecodeZeroExTransactionData(&_DevUtils.CallOpts, transactionData)
}

// EncodeERC1155AssetData is a free data retrieval call binding the contract method 0xb43cffe1.
//
// Solidity: function encodeERC1155AssetData(address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData) constant returns(bytes assetData)
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
// Solidity: function encodeERC1155AssetData(address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData) constant returns(bytes assetData)
func (_DevUtils *DevUtilsSession) EncodeERC1155AssetData(tokenAddress common.Address, tokenIds []*big.Int, tokenValues []*big.Int, callbackData []byte) ([]byte, error) {
	return _DevUtils.Contract.EncodeERC1155AssetData(&_DevUtils.CallOpts, tokenAddress, tokenIds, tokenValues, callbackData)
}

// EncodeERC1155AssetData is a free data retrieval call binding the contract method 0xb43cffe1.
//
// Solidity: function encodeERC1155AssetData(address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData) constant returns(bytes assetData)
func (_DevUtils *DevUtilsCallerSession) EncodeERC1155AssetData(tokenAddress common.Address, tokenIds []*big.Int, tokenValues []*big.Int, callbackData []byte) ([]byte, error) {
	return _DevUtils.Contract.EncodeERC1155AssetData(&_DevUtils.CallOpts, tokenAddress, tokenIds, tokenValues, callbackData)
}

// EncodeERC20AssetData is a free data retrieval call binding the contract method 0x590aa875.
//
// Solidity: function encodeERC20AssetData(address tokenAddress) constant returns(bytes assetData)
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
// Solidity: function encodeERC20AssetData(address tokenAddress) constant returns(bytes assetData)
func (_DevUtils *DevUtilsSession) EncodeERC20AssetData(tokenAddress common.Address) ([]byte, error) {
	return _DevUtils.Contract.EncodeERC20AssetData(&_DevUtils.CallOpts, tokenAddress)
}

// EncodeERC20AssetData is a free data retrieval call binding the contract method 0x590aa875.
//
// Solidity: function encodeERC20AssetData(address tokenAddress) constant returns(bytes assetData)
func (_DevUtils *DevUtilsCallerSession) EncodeERC20AssetData(tokenAddress common.Address) ([]byte, error) {
	return _DevUtils.Contract.EncodeERC20AssetData(&_DevUtils.CallOpts, tokenAddress)
}

// EncodeERC721AssetData is a free data retrieval call binding the contract method 0xa6627e9f.
//
// Solidity: function encodeERC721AssetData(address tokenAddress, uint256 tokenId) constant returns(bytes assetData)
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
// Solidity: function encodeERC721AssetData(address tokenAddress, uint256 tokenId) constant returns(bytes assetData)
func (_DevUtils *DevUtilsSession) EncodeERC721AssetData(tokenAddress common.Address, tokenId *big.Int) ([]byte, error) {
	return _DevUtils.Contract.EncodeERC721AssetData(&_DevUtils.CallOpts, tokenAddress, tokenId)
}

// EncodeERC721AssetData is a free data retrieval call binding the contract method 0xa6627e9f.
//
// Solidity: function encodeERC721AssetData(address tokenAddress, uint256 tokenId) constant returns(bytes assetData)
func (_DevUtils *DevUtilsCallerSession) EncodeERC721AssetData(tokenAddress common.Address, tokenId *big.Int) ([]byte, error) {
	return _DevUtils.Contract.EncodeERC721AssetData(&_DevUtils.CallOpts, tokenAddress, tokenId)
}

// EncodeMultiAssetData is a free data retrieval call binding the contract method 0xd3d862d1.
//
// Solidity: function encodeMultiAssetData(uint256[] amounts, bytes[] nestedAssetData) constant returns(bytes assetData)
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
// Solidity: function encodeMultiAssetData(uint256[] amounts, bytes[] nestedAssetData) constant returns(bytes assetData)
func (_DevUtils *DevUtilsSession) EncodeMultiAssetData(amounts []*big.Int, nestedAssetData [][]byte) ([]byte, error) {
	return _DevUtils.Contract.EncodeMultiAssetData(&_DevUtils.CallOpts, amounts, nestedAssetData)
}

// EncodeMultiAssetData is a free data retrieval call binding the contract method 0xd3d862d1.
//
// Solidity: function encodeMultiAssetData(uint256[] amounts, bytes[] nestedAssetData) constant returns(bytes assetData)
func (_DevUtils *DevUtilsCallerSession) EncodeMultiAssetData(amounts []*big.Int, nestedAssetData [][]byte) ([]byte, error) {
	return _DevUtils.Contract.EncodeMultiAssetData(&_DevUtils.CallOpts, amounts, nestedAssetData)
}

// GetAssetProxyAllowance is a free data retrieval call binding the contract method 0xd186037f.
//
// Solidity: function getAssetProxyAllowance(address ownerAddress, bytes assetData) constant returns(uint256 allowance)
func (_DevUtils *DevUtilsCaller) GetAssetProxyAllowance(opts *bind.CallOpts, ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "getAssetProxyAllowance", ownerAddress, assetData)
	return *ret0, err
}

// GetAssetProxyAllowance is a free data retrieval call binding the contract method 0xd186037f.
//
// Solidity: function getAssetProxyAllowance(address ownerAddress, bytes assetData) constant returns(uint256 allowance)
func (_DevUtils *DevUtilsSession) GetAssetProxyAllowance(ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	return _DevUtils.Contract.GetAssetProxyAllowance(&_DevUtils.CallOpts, ownerAddress, assetData)
}

// GetAssetProxyAllowance is a free data retrieval call binding the contract method 0xd186037f.
//
// Solidity: function getAssetProxyAllowance(address ownerAddress, bytes assetData) constant returns(uint256 allowance)
func (_DevUtils *DevUtilsCallerSession) GetAssetProxyAllowance(ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	return _DevUtils.Contract.GetAssetProxyAllowance(&_DevUtils.CallOpts, ownerAddress, assetData)
}

// GetBalance is a free data retrieval call binding the contract method 0x7d727512.
//
// Solidity: function getBalance(address ownerAddress, bytes assetData) constant returns(uint256 balance)
func (_DevUtils *DevUtilsCaller) GetBalance(opts *bind.CallOpts, ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "getBalance", ownerAddress, assetData)
	return *ret0, err
}

// GetBalance is a free data retrieval call binding the contract method 0x7d727512.
//
// Solidity: function getBalance(address ownerAddress, bytes assetData) constant returns(uint256 balance)
func (_DevUtils *DevUtilsSession) GetBalance(ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	return _DevUtils.Contract.GetBalance(&_DevUtils.CallOpts, ownerAddress, assetData)
}

// GetBalance is a free data retrieval call binding the contract method 0x7d727512.
//
// Solidity: function getBalance(address ownerAddress, bytes assetData) constant returns(uint256 balance)
func (_DevUtils *DevUtilsCallerSession) GetBalance(ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	return _DevUtils.Contract.GetBalance(&_DevUtils.CallOpts, ownerAddress, assetData)
}

// GetBalanceAndAssetProxyAllowance is a free data retrieval call binding the contract method 0x0d7b7d76.
//
// Solidity: function getBalanceAndAssetProxyAllowance(address ownerAddress, bytes assetData) constant returns(uint256 balance, uint256 allowance)
func (_DevUtils *DevUtilsCaller) GetBalanceAndAssetProxyAllowance(opts *bind.CallOpts, ownerAddress common.Address, assetData []byte) (struct {
	Balance   *big.Int
	Allowance *big.Int
}, error) {
	ret := new(struct {
		Balance   *big.Int
		Allowance *big.Int
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "getBalanceAndAssetProxyAllowance", ownerAddress, assetData)
	return *ret, err
}

// GetBalanceAndAssetProxyAllowance is a free data retrieval call binding the contract method 0x0d7b7d76.
//
// Solidity: function getBalanceAndAssetProxyAllowance(address ownerAddress, bytes assetData) constant returns(uint256 balance, uint256 allowance)
func (_DevUtils *DevUtilsSession) GetBalanceAndAssetProxyAllowance(ownerAddress common.Address, assetData []byte) (struct {
	Balance   *big.Int
	Allowance *big.Int
}, error) {
	return _DevUtils.Contract.GetBalanceAndAssetProxyAllowance(&_DevUtils.CallOpts, ownerAddress, assetData)
}

// GetBalanceAndAssetProxyAllowance is a free data retrieval call binding the contract method 0x0d7b7d76.
//
// Solidity: function getBalanceAndAssetProxyAllowance(address ownerAddress, bytes assetData) constant returns(uint256 balance, uint256 allowance)
func (_DevUtils *DevUtilsCallerSession) GetBalanceAndAssetProxyAllowance(ownerAddress common.Address, assetData []byte) (struct {
	Balance   *big.Int
	Allowance *big.Int
}, error) {
	return _DevUtils.Contract.GetBalanceAndAssetProxyAllowance(&_DevUtils.CallOpts, ownerAddress, assetData)
}

// GetBatchAssetProxyAllowances is a free data retrieval call binding the contract method 0x4dfdac20.
//
// Solidity: function getBatchAssetProxyAllowances(address ownerAddress, bytes[] assetData) constant returns(uint256[] allowances)
func (_DevUtils *DevUtilsCaller) GetBatchAssetProxyAllowances(opts *bind.CallOpts, ownerAddress common.Address, assetData [][]byte) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "getBatchAssetProxyAllowances", ownerAddress, assetData)
	return *ret0, err
}

// GetBatchAssetProxyAllowances is a free data retrieval call binding the contract method 0x4dfdac20.
//
// Solidity: function getBatchAssetProxyAllowances(address ownerAddress, bytes[] assetData) constant returns(uint256[] allowances)
func (_DevUtils *DevUtilsSession) GetBatchAssetProxyAllowances(ownerAddress common.Address, assetData [][]byte) ([]*big.Int, error) {
	return _DevUtils.Contract.GetBatchAssetProxyAllowances(&_DevUtils.CallOpts, ownerAddress, assetData)
}

// GetBatchAssetProxyAllowances is a free data retrieval call binding the contract method 0x4dfdac20.
//
// Solidity: function getBatchAssetProxyAllowances(address ownerAddress, bytes[] assetData) constant returns(uint256[] allowances)
func (_DevUtils *DevUtilsCallerSession) GetBatchAssetProxyAllowances(ownerAddress common.Address, assetData [][]byte) ([]*big.Int, error) {
	return _DevUtils.Contract.GetBatchAssetProxyAllowances(&_DevUtils.CallOpts, ownerAddress, assetData)
}

// GetBatchBalances is a free data retrieval call binding the contract method 0xd001c5dc.
//
// Solidity: function getBatchBalances(address ownerAddress, bytes[] assetData) constant returns(uint256[] balances)
func (_DevUtils *DevUtilsCaller) GetBatchBalances(opts *bind.CallOpts, ownerAddress common.Address, assetData [][]byte) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "getBatchBalances", ownerAddress, assetData)
	return *ret0, err
}

// GetBatchBalances is a free data retrieval call binding the contract method 0xd001c5dc.
//
// Solidity: function getBatchBalances(address ownerAddress, bytes[] assetData) constant returns(uint256[] balances)
func (_DevUtils *DevUtilsSession) GetBatchBalances(ownerAddress common.Address, assetData [][]byte) ([]*big.Int, error) {
	return _DevUtils.Contract.GetBatchBalances(&_DevUtils.CallOpts, ownerAddress, assetData)
}

// GetBatchBalances is a free data retrieval call binding the contract method 0xd001c5dc.
//
// Solidity: function getBatchBalances(address ownerAddress, bytes[] assetData) constant returns(uint256[] balances)
func (_DevUtils *DevUtilsCallerSession) GetBatchBalances(ownerAddress common.Address, assetData [][]byte) ([]*big.Int, error) {
	return _DevUtils.Contract.GetBatchBalances(&_DevUtils.CallOpts, ownerAddress, assetData)
}

// GetBatchBalancesAndAssetProxyAllowances is a free data retrieval call binding the contract method 0xe4e6e7da.
//
// Solidity: function getBatchBalancesAndAssetProxyAllowances(address ownerAddress, bytes[] assetData) constant returns(uint256[] balances, uint256[] allowances)
func (_DevUtils *DevUtilsCaller) GetBatchBalancesAndAssetProxyAllowances(opts *bind.CallOpts, ownerAddress common.Address, assetData [][]byte) (struct {
	Balances   []*big.Int
	Allowances []*big.Int
}, error) {
	ret := new(struct {
		Balances   []*big.Int
		Allowances []*big.Int
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "getBatchBalancesAndAssetProxyAllowances", ownerAddress, assetData)
	return *ret, err
}

// GetBatchBalancesAndAssetProxyAllowances is a free data retrieval call binding the contract method 0xe4e6e7da.
//
// Solidity: function getBatchBalancesAndAssetProxyAllowances(address ownerAddress, bytes[] assetData) constant returns(uint256[] balances, uint256[] allowances)
func (_DevUtils *DevUtilsSession) GetBatchBalancesAndAssetProxyAllowances(ownerAddress common.Address, assetData [][]byte) (struct {
	Balances   []*big.Int
	Allowances []*big.Int
}, error) {
	return _DevUtils.Contract.GetBatchBalancesAndAssetProxyAllowances(&_DevUtils.CallOpts, ownerAddress, assetData)
}

// GetBatchBalancesAndAssetProxyAllowances is a free data retrieval call binding the contract method 0xe4e6e7da.
//
// Solidity: function getBatchBalancesAndAssetProxyAllowances(address ownerAddress, bytes[] assetData) constant returns(uint256[] balances, uint256[] allowances)
func (_DevUtils *DevUtilsCallerSession) GetBatchBalancesAndAssetProxyAllowances(ownerAddress common.Address, assetData [][]byte) (struct {
	Balances   []*big.Int
	Allowances []*big.Int
}, error) {
	return _DevUtils.Contract.GetBatchBalancesAndAssetProxyAllowances(&_DevUtils.CallOpts, ownerAddress, assetData)
}

// GetEthBalances is a free data retrieval call binding the contract method 0xa0901e51.
//
// Solidity: function getEthBalances(address[] addresses) constant returns(uint256[])
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
// Solidity: function getEthBalances(address[] addresses) constant returns(uint256[])
func (_DevUtils *DevUtilsSession) GetEthBalances(addresses []common.Address) ([]*big.Int, error) {
	return _DevUtils.Contract.GetEthBalances(&_DevUtils.CallOpts, addresses)
}

// GetEthBalances is a free data retrieval call binding the contract method 0xa0901e51.
//
// Solidity: function getEthBalances(address[] addresses) constant returns(uint256[])
func (_DevUtils *DevUtilsCallerSession) GetEthBalances(addresses []common.Address) ([]*big.Int, error) {
	return _DevUtils.Contract.GetEthBalances(&_DevUtils.CallOpts, addresses)
}

// GetOrderRelevantState is a free data retrieval call binding the contract method 0xe77286eb.
//
// Solidity: function getOrderRelevantState(Struct0 order, bytes signature) constant returns(Struct1 orderInfo, uint256 fillableTakerAssetAmount, bool isValidSignature)
func (_DevUtils *DevUtilsCaller) GetOrderRelevantState(opts *bind.CallOpts, order Struct0, signature []byte) (struct {
	OrderInfo                Struct1
	FillableTakerAssetAmount *big.Int
	IsValidSignature         bool
}, error) {
	ret := new(struct {
		OrderInfo                Struct1
		FillableTakerAssetAmount *big.Int
		IsValidSignature         bool
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "getOrderRelevantState", order, signature)
	return *ret, err
}

// GetOrderRelevantState is a free data retrieval call binding the contract method 0xe77286eb.
//
// Solidity: function getOrderRelevantState(Struct0 order, bytes signature) constant returns(Struct1 orderInfo, uint256 fillableTakerAssetAmount, bool isValidSignature)
func (_DevUtils *DevUtilsSession) GetOrderRelevantState(order Struct0, signature []byte) (struct {
	OrderInfo                Struct1
	FillableTakerAssetAmount *big.Int
	IsValidSignature         bool
}, error) {
	return _DevUtils.Contract.GetOrderRelevantState(&_DevUtils.CallOpts, order, signature)
}

// GetOrderRelevantState is a free data retrieval call binding the contract method 0xe77286eb.
//
// Solidity: function getOrderRelevantState(Struct0 order, bytes signature) constant returns(Struct1 orderInfo, uint256 fillableTakerAssetAmount, bool isValidSignature)
func (_DevUtils *DevUtilsCallerSession) GetOrderRelevantState(order Struct0, signature []byte) (struct {
	OrderInfo                Struct1
	FillableTakerAssetAmount *big.Int
	IsValidSignature         bool
}, error) {
	return _DevUtils.Contract.GetOrderRelevantState(&_DevUtils.CallOpts, order, signature)
}

// GetOrderRelevantStates is a free data retrieval call binding the contract method 0xe25cabf7.
//
// Solidity: function getOrderRelevantStates([]Struct0 orders, bytes[] signatures) constant returns([]Struct1 ordersInfo, uint256[] fillableTakerAssetAmounts, bool[] isValidSignature)
func (_DevUtils *DevUtilsCaller) GetOrderRelevantStates(opts *bind.CallOpts, orders []Struct0, signatures [][]byte) (struct {
	OrdersInfo                []Struct1
	FillableTakerAssetAmounts []*big.Int
	IsValidSignature          []bool
}, error) {
	ret := new(struct {
		OrdersInfo                []Struct1
		FillableTakerAssetAmounts []*big.Int
		IsValidSignature          []bool
	})
	out := ret
	err := _DevUtils.contract.Call(opts, out, "getOrderRelevantStates", orders, signatures)
	return *ret, err
}

// GetOrderRelevantStates is a free data retrieval call binding the contract method 0xe25cabf7.
//
// Solidity: function getOrderRelevantStates([]Struct0 orders, bytes[] signatures) constant returns([]Struct1 ordersInfo, uint256[] fillableTakerAssetAmounts, bool[] isValidSignature)
func (_DevUtils *DevUtilsSession) GetOrderRelevantStates(orders []Struct0, signatures [][]byte) (struct {
	OrdersInfo                []Struct1
	FillableTakerAssetAmounts []*big.Int
	IsValidSignature          []bool
}, error) {
	return _DevUtils.Contract.GetOrderRelevantStates(&_DevUtils.CallOpts, orders, signatures)
}

// GetOrderRelevantStates is a free data retrieval call binding the contract method 0xe25cabf7.
//
// Solidity: function getOrderRelevantStates([]Struct0 orders, bytes[] signatures) constant returns([]Struct1 ordersInfo, uint256[] fillableTakerAssetAmounts, bool[] isValidSignature)
func (_DevUtils *DevUtilsCallerSession) GetOrderRelevantStates(orders []Struct0, signatures [][]byte) (struct {
	OrdersInfo                []Struct1
	FillableTakerAssetAmounts []*big.Int
	IsValidSignature          []bool
}, error) {
	return _DevUtils.Contract.GetOrderRelevantStates(&_DevUtils.CallOpts, orders, signatures)
}

// GetTransferableAssetAmount is a free data retrieval call binding the contract method 0x2322cf76.
//
// Solidity: function getTransferableAssetAmount(address ownerAddress, bytes assetData) constant returns(uint256 transferableAssetAmount)
func (_DevUtils *DevUtilsCaller) GetTransferableAssetAmount(opts *bind.CallOpts, ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "getTransferableAssetAmount", ownerAddress, assetData)
	return *ret0, err
}

// GetTransferableAssetAmount is a free data retrieval call binding the contract method 0x2322cf76.
//
// Solidity: function getTransferableAssetAmount(address ownerAddress, bytes assetData) constant returns(uint256 transferableAssetAmount)
func (_DevUtils *DevUtilsSession) GetTransferableAssetAmount(ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	return _DevUtils.Contract.GetTransferableAssetAmount(&_DevUtils.CallOpts, ownerAddress, assetData)
}

// GetTransferableAssetAmount is a free data retrieval call binding the contract method 0x2322cf76.
//
// Solidity: function getTransferableAssetAmount(address ownerAddress, bytes assetData) constant returns(uint256 transferableAssetAmount)
func (_DevUtils *DevUtilsCallerSession) GetTransferableAssetAmount(ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	return _DevUtils.Contract.GetTransferableAssetAmount(&_DevUtils.CallOpts, ownerAddress, assetData)
}

// GetSimulatedOrderTransferResults is a paid mutator transaction binding the contract method 0xd3637905.
//
// Solidity: function getSimulatedOrderTransferResults(Struct0 order, address takerAddress, uint256 takerAssetFillAmount) returns(uint8 orderTransferResults)
func (_DevUtils *DevUtilsTransactor) GetSimulatedOrderTransferResults(opts *bind.TransactOpts, order Struct0, takerAddress common.Address, takerAssetFillAmount *big.Int) (*types.Transaction, error) {
	return _DevUtils.contract.Transact(opts, "getSimulatedOrderTransferResults", order, takerAddress, takerAssetFillAmount)
}

// GetSimulatedOrderTransferResults is a paid mutator transaction binding the contract method 0xd3637905.
//
// Solidity: function getSimulatedOrderTransferResults(Struct0 order, address takerAddress, uint256 takerAssetFillAmount) returns(uint8 orderTransferResults)
func (_DevUtils *DevUtilsSession) GetSimulatedOrderTransferResults(order Struct0, takerAddress common.Address, takerAssetFillAmount *big.Int) (*types.Transaction, error) {
	return _DevUtils.Contract.GetSimulatedOrderTransferResults(&_DevUtils.TransactOpts, order, takerAddress, takerAssetFillAmount)
}

// GetSimulatedOrderTransferResults is a paid mutator transaction binding the contract method 0xd3637905.
//
// Solidity: function getSimulatedOrderTransferResults(Struct0 order, address takerAddress, uint256 takerAssetFillAmount) returns(uint8 orderTransferResults)
func (_DevUtils *DevUtilsTransactorSession) GetSimulatedOrderTransferResults(order Struct0, takerAddress common.Address, takerAssetFillAmount *big.Int) (*types.Transaction, error) {
	return _DevUtils.Contract.GetSimulatedOrderTransferResults(&_DevUtils.TransactOpts, order, takerAddress, takerAssetFillAmount)
}

// GetSimulatedOrdersTransferResults is a paid mutator transaction binding the contract method 0xa5cd62ba.
//
// Solidity: function getSimulatedOrdersTransferResults([]Struct0 orders, address[] takerAddresses, uint256[] takerAssetFillAmounts) returns(uint8[] orderTransferResults)
func (_DevUtils *DevUtilsTransactor) GetSimulatedOrdersTransferResults(opts *bind.TransactOpts, orders []Struct0, takerAddresses []common.Address, takerAssetFillAmounts []*big.Int) (*types.Transaction, error) {
	return _DevUtils.contract.Transact(opts, "getSimulatedOrdersTransferResults", orders, takerAddresses, takerAssetFillAmounts)
}

// GetSimulatedOrdersTransferResults is a paid mutator transaction binding the contract method 0xa5cd62ba.
//
// Solidity: function getSimulatedOrdersTransferResults([]Struct0 orders, address[] takerAddresses, uint256[] takerAssetFillAmounts) returns(uint8[] orderTransferResults)
func (_DevUtils *DevUtilsSession) GetSimulatedOrdersTransferResults(orders []Struct0, takerAddresses []common.Address, takerAssetFillAmounts []*big.Int) (*types.Transaction, error) {
	return _DevUtils.Contract.GetSimulatedOrdersTransferResults(&_DevUtils.TransactOpts, orders, takerAddresses, takerAssetFillAmounts)
}

// GetSimulatedOrdersTransferResults is a paid mutator transaction binding the contract method 0xa5cd62ba.
//
// Solidity: function getSimulatedOrdersTransferResults([]Struct0 orders, address[] takerAddresses, uint256[] takerAssetFillAmounts) returns(uint8[] orderTransferResults)
func (_DevUtils *DevUtilsTransactorSession) GetSimulatedOrdersTransferResults(orders []Struct0, takerAddresses []common.Address, takerAssetFillAmounts []*big.Int) (*types.Transaction, error) {
	return _DevUtils.Contract.GetSimulatedOrdersTransferResults(&_DevUtils.TransactOpts, orders, takerAddresses, takerAssetFillAmounts)
}

