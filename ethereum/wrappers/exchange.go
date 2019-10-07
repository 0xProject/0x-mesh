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

// ExchangeABI is the input ABI used to generate the binding from.
const ExchangeABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"filled\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchFillOrders\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"totalFillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"cancelled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"signerAddress\",\"type\":\"address\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"preSign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"leftOrder\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"rightOrder\",\"type\":\"tuple\"},{\"name\":\"leftSignature\",\"type\":\"bytes\"},{\"name\":\"rightSignature\",\"type\":\"bytes\"}],\"name\":\"matchOrders\",\"outputs\":[{\"components\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"left\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"right\",\"type\":\"tuple\"},{\"name\":\"leftMakerAssetSpreadAmount\",\"type\":\"uint256\"}],\"name\":\"matchedFillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"},{\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"fillOrderNoThrow\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"name\":\"assetProxies\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"}],\"name\":\"batchCancelOrders\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchFillOrKillOrders\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"totalFillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"targetOrderEpoch\",\"type\":\"uint256\"}],\"name\":\"cancelOrdersUpTo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchFillOrdersNoThrow\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"totalFillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"assetProxyId\",\"type\":\"bytes4\"}],\"name\":\"getAssetProxy\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"transactions\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"},{\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"fillOrKillOrder\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"validatorAddress\",\"type\":\"address\"},{\"name\":\"approval\",\"type\":\"bool\"}],\"name\":\"setSignatureValidatorApproval\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedValidators\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketSellOrders\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"totalFillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"}],\"name\":\"getOrdersInfo\",\"outputs\":[{\"components\":[{\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"name\":\"\",\"type\":\"tuple[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"preSigned\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"signerAddress\",\"type\":\"address\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"isValidSignature\",\"outputs\":[{\"name\":\"isValid\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"makerAssetFillAmount\",\"type\":\"uint256\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketBuyOrdersNoThrow\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"totalFillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"},{\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"fillOrder\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"signerAddress\",\"type\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"executeTransaction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"assetProxy\",\"type\":\"address\"}],\"name\":\"registerAssetProxy\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"getOrderInfo\",\"outputs\":[{\"components\":[{\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"name\":\"orderInfo\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"cancelOrder\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"orderEpoch\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ZRX_ASSET_DATA\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketSellOrdersNoThrow\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"totalFillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"EIP712_DOMAIN_HASH\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"makerAssetFillAmount\",\"type\":\"uint256\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketBuyOrders\",\"outputs\":[{\"components\":[{\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"name\":\"takerFeePaid\",\"type\":\"uint256\"}],\"name\":\"totalFillResults\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentContextAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_zrxAssetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"signerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"SignatureValidatorApproval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"takerAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"Fill\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"Cancel\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"orderEpoch\",\"type\":\"uint256\"}],\"name\":\"CancelUpTo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"id\",\"type\":\"bytes4\"},{\"indexed\":false,\"name\":\"assetProxy\",\"type\":\"address\"}],\"name\":\"AssetProxyRegistered\",\"type\":\"event\"}]"

// Exchange is an auto generated Go binding around an Ethereum contract.
type Exchange struct {
	ExchangeCaller     // Read-only binding to the contract
	ExchangeTransactor // Write-only binding to the contract
	ExchangeFilterer   // Log filterer for contract events
}

// ExchangeCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExchangeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExchangeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ExchangeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExchangeSession struct {
	Contract     *Exchange         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ExchangeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExchangeCallerSession struct {
	Contract *ExchangeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ExchangeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExchangeTransactorSession struct {
	Contract     *ExchangeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ExchangeRaw is an auto generated low-level Go binding around an Ethereum contract.
type ExchangeRaw struct {
	Contract *Exchange // Generic contract binding to access the raw methods on
}

// ExchangeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExchangeCallerRaw struct {
	Contract *ExchangeCaller // Generic read-only contract binding to access the raw methods on
}

// ExchangeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExchangeTransactorRaw struct {
	Contract *ExchangeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExchange creates a new instance of Exchange, bound to a specific deployed contract.
func NewExchange(address common.Address, backend bind.ContractBackend) (*Exchange, error) {
	contract, err := bindExchange(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Exchange{ExchangeCaller: ExchangeCaller{contract: contract}, ExchangeTransactor: ExchangeTransactor{contract: contract}, ExchangeFilterer: ExchangeFilterer{contract: contract}}, nil
}

// NewExchangeCaller creates a new read-only instance of Exchange, bound to a specific deployed contract.
func NewExchangeCaller(address common.Address, caller bind.ContractCaller) (*ExchangeCaller, error) {
	contract, err := bindExchange(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangeCaller{contract: contract}, nil
}

// NewExchangeTransactor creates a new write-only instance of Exchange, bound to a specific deployed contract.
func NewExchangeTransactor(address common.Address, transactor bind.ContractTransactor) (*ExchangeTransactor, error) {
	contract, err := bindExchange(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangeTransactor{contract: contract}, nil
}

// NewExchangeFilterer creates a new log filterer instance of Exchange, bound to a specific deployed contract.
func NewExchangeFilterer(address common.Address, filterer bind.ContractFilterer) (*ExchangeFilterer, error) {
	contract, err := bindExchange(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExchangeFilterer{contract: contract}, nil
}

// bindExchange binds a generic wrapper to an already deployed contract.
func bindExchange(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ExchangeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Exchange *ExchangeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Exchange.Contract.ExchangeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Exchange *ExchangeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Exchange.Contract.ExchangeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Exchange *ExchangeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Exchange.Contract.ExchangeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Exchange *ExchangeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Exchange.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Exchange *ExchangeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Exchange.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Exchange *ExchangeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Exchange.Contract.contract.Transact(opts, method, params...)
}

// Struct2 is an auto generated low-level Go binding around an user-defined struct.
type Struct2 struct {
	Left                       Struct3
	Right                      Struct3
	LeftMakerAssetSpreadAmount *big.Int
}

// Struct3 is an auto generated low-level Go binding around an user-defined struct.
type Struct3 struct {
	MakerAssetFilledAmount *big.Int
	TakerAssetFilledAmount *big.Int
	MakerFeePaid           *big.Int
	TakerFeePaid           *big.Int
}

// EIP712DOMAINHASH is a free data retrieval call binding the contract method 0xe306f779.
//
// Solidity: function EIP712_DOMAIN_HASH() constant returns(bytes32)
func (_Exchange *ExchangeCaller) EIP712DOMAINHASH(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "EIP712_DOMAIN_HASH")
	return *ret0, err
}

// EIP712DOMAINHASH is a free data retrieval call binding the contract method 0xe306f779.
//
// Solidity: function EIP712_DOMAIN_HASH() constant returns(bytes32)
func (_Exchange *ExchangeSession) EIP712DOMAINHASH() ([32]byte, error) {
	return _Exchange.Contract.EIP712DOMAINHASH(&_Exchange.CallOpts)
}

// EIP712DOMAINHASH is a free data retrieval call binding the contract method 0xe306f779.
//
// Solidity: function EIP712_DOMAIN_HASH() constant returns(bytes32)
func (_Exchange *ExchangeCallerSession) EIP712DOMAINHASH() ([32]byte, error) {
	return _Exchange.Contract.EIP712DOMAINHASH(&_Exchange.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Exchange *ExchangeCaller) VERSION(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "VERSION")
	return *ret0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Exchange *ExchangeSession) VERSION() (string, error) {
	return _Exchange.Contract.VERSION(&_Exchange.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() constant returns(string)
func (_Exchange *ExchangeCallerSession) VERSION() (string, error) {
	return _Exchange.Contract.VERSION(&_Exchange.CallOpts)
}

// ZRXASSETDATA is a free data retrieval call binding the contract method 0xdb123b1a.
//
// Solidity: function ZRX_ASSET_DATA() constant returns(bytes)
func (_Exchange *ExchangeCaller) ZRXASSETDATA(opts *bind.CallOpts) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "ZRX_ASSET_DATA")
	return *ret0, err
}

// ZRXASSETDATA is a free data retrieval call binding the contract method 0xdb123b1a.
//
// Solidity: function ZRX_ASSET_DATA() constant returns(bytes)
func (_Exchange *ExchangeSession) ZRXASSETDATA() ([]byte, error) {
	return _Exchange.Contract.ZRXASSETDATA(&_Exchange.CallOpts)
}

// ZRXASSETDATA is a free data retrieval call binding the contract method 0xdb123b1a.
//
// Solidity: function ZRX_ASSET_DATA() constant returns(bytes)
func (_Exchange *ExchangeCallerSession) ZRXASSETDATA() ([]byte, error) {
	return _Exchange.Contract.ZRXASSETDATA(&_Exchange.CallOpts)
}

// AllowedValidators is a free data retrieval call binding the contract method 0x7b8e3514.
//
// Solidity: function allowedValidators(address , address ) constant returns(bool)
func (_Exchange *ExchangeCaller) AllowedValidators(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "allowedValidators", arg0, arg1)
	return *ret0, err
}

// AllowedValidators is a free data retrieval call binding the contract method 0x7b8e3514.
//
// Solidity: function allowedValidators(address , address ) constant returns(bool)
func (_Exchange *ExchangeSession) AllowedValidators(arg0 common.Address, arg1 common.Address) (bool, error) {
	return _Exchange.Contract.AllowedValidators(&_Exchange.CallOpts, arg0, arg1)
}

// AllowedValidators is a free data retrieval call binding the contract method 0x7b8e3514.
//
// Solidity: function allowedValidators(address , address ) constant returns(bool)
func (_Exchange *ExchangeCallerSession) AllowedValidators(arg0 common.Address, arg1 common.Address) (bool, error) {
	return _Exchange.Contract.AllowedValidators(&_Exchange.CallOpts, arg0, arg1)
}

// AssetProxies is a free data retrieval call binding the contract method 0x3fd3c997.
//
// Solidity: function assetProxies(bytes4 ) constant returns(address)
func (_Exchange *ExchangeCaller) AssetProxies(opts *bind.CallOpts, arg0 [4]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "assetProxies", arg0)
	return *ret0, err
}

// AssetProxies is a free data retrieval call binding the contract method 0x3fd3c997.
//
// Solidity: function assetProxies(bytes4 ) constant returns(address)
func (_Exchange *ExchangeSession) AssetProxies(arg0 [4]byte) (common.Address, error) {
	return _Exchange.Contract.AssetProxies(&_Exchange.CallOpts, arg0)
}

// AssetProxies is a free data retrieval call binding the contract method 0x3fd3c997.
//
// Solidity: function assetProxies(bytes4 ) constant returns(address)
func (_Exchange *ExchangeCallerSession) AssetProxies(arg0 [4]byte) (common.Address, error) {
	return _Exchange.Contract.AssetProxies(&_Exchange.CallOpts, arg0)
}

// Cancelled is a free data retrieval call binding the contract method 0x2ac12622.
//
// Solidity: function cancelled(bytes32 ) constant returns(bool)
func (_Exchange *ExchangeCaller) Cancelled(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "cancelled", arg0)
	return *ret0, err
}

// Cancelled is a free data retrieval call binding the contract method 0x2ac12622.
//
// Solidity: function cancelled(bytes32 ) constant returns(bool)
func (_Exchange *ExchangeSession) Cancelled(arg0 [32]byte) (bool, error) {
	return _Exchange.Contract.Cancelled(&_Exchange.CallOpts, arg0)
}

// Cancelled is a free data retrieval call binding the contract method 0x2ac12622.
//
// Solidity: function cancelled(bytes32 ) constant returns(bool)
func (_Exchange *ExchangeCallerSession) Cancelled(arg0 [32]byte) (bool, error) {
	return _Exchange.Contract.Cancelled(&_Exchange.CallOpts, arg0)
}

// CurrentContextAddress is a free data retrieval call binding the contract method 0xeea086ba.
//
// Solidity: function currentContextAddress() constant returns(address)
func (_Exchange *ExchangeCaller) CurrentContextAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "currentContextAddress")
	return *ret0, err
}

// CurrentContextAddress is a free data retrieval call binding the contract method 0xeea086ba.
//
// Solidity: function currentContextAddress() constant returns(address)
func (_Exchange *ExchangeSession) CurrentContextAddress() (common.Address, error) {
	return _Exchange.Contract.CurrentContextAddress(&_Exchange.CallOpts)
}

// CurrentContextAddress is a free data retrieval call binding the contract method 0xeea086ba.
//
// Solidity: function currentContextAddress() constant returns(address)
func (_Exchange *ExchangeCallerSession) CurrentContextAddress() (common.Address, error) {
	return _Exchange.Contract.CurrentContextAddress(&_Exchange.CallOpts)
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled(bytes32 ) constant returns(uint256)
func (_Exchange *ExchangeCaller) Filled(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "filled", arg0)
	return *ret0, err
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled(bytes32 ) constant returns(uint256)
func (_Exchange *ExchangeSession) Filled(arg0 [32]byte) (*big.Int, error) {
	return _Exchange.Contract.Filled(&_Exchange.CallOpts, arg0)
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled(bytes32 ) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) Filled(arg0 [32]byte) (*big.Int, error) {
	return _Exchange.Contract.Filled(&_Exchange.CallOpts, arg0)
}

// GetAssetProxy is a free data retrieval call binding the contract method 0x60704108.
//
// Solidity: function getAssetProxy(bytes4 assetProxyId) constant returns(address)
func (_Exchange *ExchangeCaller) GetAssetProxy(opts *bind.CallOpts, assetProxyId [4]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "getAssetProxy", assetProxyId)
	return *ret0, err
}

// GetAssetProxy is a free data retrieval call binding the contract method 0x60704108.
//
// Solidity: function getAssetProxy(bytes4 assetProxyId) constant returns(address)
func (_Exchange *ExchangeSession) GetAssetProxy(assetProxyId [4]byte) (common.Address, error) {
	return _Exchange.Contract.GetAssetProxy(&_Exchange.CallOpts, assetProxyId)
}

// GetAssetProxy is a free data retrieval call binding the contract method 0x60704108.
//
// Solidity: function getAssetProxy(bytes4 assetProxyId) constant returns(address)
func (_Exchange *ExchangeCallerSession) GetAssetProxy(assetProxyId [4]byte) (common.Address, error) {
	return _Exchange.Contract.GetAssetProxy(&_Exchange.CallOpts, assetProxyId)
}

// GetOrderInfo is a free data retrieval call binding the contract method 0xc75e0a81.
//
// Solidity: function getOrderInfo(Struct0 order) constant returns(Struct1 orderInfo)
func (_Exchange *ExchangeCaller) GetOrderInfo(opts *bind.CallOpts, order Struct0) (Struct1, error) {
	var (
		ret0 = new(Struct1)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "getOrderInfo", order)
	return *ret0, err
}

// GetOrderInfo is a free data retrieval call binding the contract method 0xc75e0a81.
//
// Solidity: function getOrderInfo(Struct0 order) constant returns(Struct1 orderInfo)
func (_Exchange *ExchangeSession) GetOrderInfo(order Struct0) (Struct1, error) {
	return _Exchange.Contract.GetOrderInfo(&_Exchange.CallOpts, order)
}

// GetOrderInfo is a free data retrieval call binding the contract method 0xc75e0a81.
//
// Solidity: function getOrderInfo(Struct0 order) constant returns(Struct1 orderInfo)
func (_Exchange *ExchangeCallerSession) GetOrderInfo(order Struct0) (Struct1, error) {
	return _Exchange.Contract.GetOrderInfo(&_Exchange.CallOpts, order)
}

// GetOrdersInfo is a free data retrieval call binding the contract method 0x7e9d74dc.
//
// Solidity: function getOrdersInfo([]Struct0 orders) constant returns([]Struct1)
func (_Exchange *ExchangeCaller) GetOrdersInfo(opts *bind.CallOpts, orders []Struct0) ([]Struct1, error) {
	var (
		ret0 = new([]Struct1)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "getOrdersInfo", orders)
	return *ret0, err
}

// GetOrdersInfo is a free data retrieval call binding the contract method 0x7e9d74dc.
//
// Solidity: function getOrdersInfo([]Struct0 orders) constant returns([]Struct1)
func (_Exchange *ExchangeSession) GetOrdersInfo(orders []Struct0) ([]Struct1, error) {
	return _Exchange.Contract.GetOrdersInfo(&_Exchange.CallOpts, orders)
}

// GetOrdersInfo is a free data retrieval call binding the contract method 0x7e9d74dc.
//
// Solidity: function getOrdersInfo([]Struct0 orders) constant returns([]Struct1)
func (_Exchange *ExchangeCallerSession) GetOrdersInfo(orders []Struct0) ([]Struct1, error) {
	return _Exchange.Contract.GetOrdersInfo(&_Exchange.CallOpts, orders)
}

// IsValidSignature is a free data retrieval call binding the contract method 0x93634702.
//
// Solidity: function isValidSignature(bytes32 hash, address signerAddress, bytes signature) constant returns(bool isValid)
func (_Exchange *ExchangeCaller) IsValidSignature(opts *bind.CallOpts, hash [32]byte, signerAddress common.Address, signature []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "isValidSignature", hash, signerAddress, signature)
	return *ret0, err
}

// IsValidSignature is a free data retrieval call binding the contract method 0x93634702.
//
// Solidity: function isValidSignature(bytes32 hash, address signerAddress, bytes signature) constant returns(bool isValid)
func (_Exchange *ExchangeSession) IsValidSignature(hash [32]byte, signerAddress common.Address, signature []byte) (bool, error) {
	return _Exchange.Contract.IsValidSignature(&_Exchange.CallOpts, hash, signerAddress, signature)
}

// IsValidSignature is a free data retrieval call binding the contract method 0x93634702.
//
// Solidity: function isValidSignature(bytes32 hash, address signerAddress, bytes signature) constant returns(bool isValid)
func (_Exchange *ExchangeCallerSession) IsValidSignature(hash [32]byte, signerAddress common.Address, signature []byte) (bool, error) {
	return _Exchange.Contract.IsValidSignature(&_Exchange.CallOpts, hash, signerAddress, signature)
}

// OrderEpoch is a free data retrieval call binding the contract method 0xd9bfa73e.
//
// Solidity: function orderEpoch(address , address ) constant returns(uint256)
func (_Exchange *ExchangeCaller) OrderEpoch(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "orderEpoch", arg0, arg1)
	return *ret0, err
}

// OrderEpoch is a free data retrieval call binding the contract method 0xd9bfa73e.
//
// Solidity: function orderEpoch(address , address ) constant returns(uint256)
func (_Exchange *ExchangeSession) OrderEpoch(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _Exchange.Contract.OrderEpoch(&_Exchange.CallOpts, arg0, arg1)
}

// OrderEpoch is a free data retrieval call binding the contract method 0xd9bfa73e.
//
// Solidity: function orderEpoch(address , address ) constant returns(uint256)
func (_Exchange *ExchangeCallerSession) OrderEpoch(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _Exchange.Contract.OrderEpoch(&_Exchange.CallOpts, arg0, arg1)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Exchange *ExchangeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Exchange *ExchangeSession) Owner() (common.Address, error) {
	return _Exchange.Contract.Owner(&_Exchange.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Exchange *ExchangeCallerSession) Owner() (common.Address, error) {
	return _Exchange.Contract.Owner(&_Exchange.CallOpts)
}

// PreSigned is a free data retrieval call binding the contract method 0x82c174d0.
//
// Solidity: function preSigned(bytes32 , address ) constant returns(bool)
func (_Exchange *ExchangeCaller) PreSigned(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "preSigned", arg0, arg1)
	return *ret0, err
}

// PreSigned is a free data retrieval call binding the contract method 0x82c174d0.
//
// Solidity: function preSigned(bytes32 , address ) constant returns(bool)
func (_Exchange *ExchangeSession) PreSigned(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _Exchange.Contract.PreSigned(&_Exchange.CallOpts, arg0, arg1)
}

// PreSigned is a free data retrieval call binding the contract method 0x82c174d0.
//
// Solidity: function preSigned(bytes32 , address ) constant returns(bool)
func (_Exchange *ExchangeCallerSession) PreSigned(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _Exchange.Contract.PreSigned(&_Exchange.CallOpts, arg0, arg1)
}

// Transactions is a free data retrieval call binding the contract method 0x642f2eaf.
//
// Solidity: function transactions(bytes32 ) constant returns(bool)
func (_Exchange *ExchangeCaller) Transactions(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Exchange.contract.Call(opts, out, "transactions", arg0)
	return *ret0, err
}

// Transactions is a free data retrieval call binding the contract method 0x642f2eaf.
//
// Solidity: function transactions(bytes32 ) constant returns(bool)
func (_Exchange *ExchangeSession) Transactions(arg0 [32]byte) (bool, error) {
	return _Exchange.Contract.Transactions(&_Exchange.CallOpts, arg0)
}

// Transactions is a free data retrieval call binding the contract method 0x642f2eaf.
//
// Solidity: function transactions(bytes32 ) constant returns(bool)
func (_Exchange *ExchangeCallerSession) Transactions(arg0 [32]byte) (bool, error) {
	return _Exchange.Contract.Transactions(&_Exchange.CallOpts, arg0)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0x4ac14782.
//
// Solidity: function batchCancelOrders([]Struct0 orders) returns()
func (_Exchange *ExchangeTransactor) BatchCancelOrders(opts *bind.TransactOpts, orders []Struct0) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchCancelOrders", orders)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0x4ac14782.
//
// Solidity: function batchCancelOrders([]Struct0 orders) returns()
func (_Exchange *ExchangeSession) BatchCancelOrders(orders []Struct0) (*types.Transaction, error) {
	return _Exchange.Contract.BatchCancelOrders(&_Exchange.TransactOpts, orders)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0x4ac14782.
//
// Solidity: function batchCancelOrders([]Struct0 orders) returns()
func (_Exchange *ExchangeTransactorSession) BatchCancelOrders(orders []Struct0) (*types.Transaction, error) {
	return _Exchange.Contract.BatchCancelOrders(&_Exchange.TransactOpts, orders)
}

// BatchFillOrKillOrders is a paid mutator transaction binding the contract method 0x4d0ae546.
//
// Solidity: function batchFillOrKillOrders([]Struct0 orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeTransactor) BatchFillOrKillOrders(opts *bind.TransactOpts, orders []Struct0, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchFillOrKillOrders", orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrKillOrders is a paid mutator transaction binding the contract method 0x4d0ae546.
//
// Solidity: function batchFillOrKillOrders([]Struct0 orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeSession) BatchFillOrKillOrders(orders []Struct0, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchFillOrKillOrders(&_Exchange.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrKillOrders is a paid mutator transaction binding the contract method 0x4d0ae546.
//
// Solidity: function batchFillOrKillOrders([]Struct0 orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeTransactorSession) BatchFillOrKillOrders(orders []Struct0, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchFillOrKillOrders(&_Exchange.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrders is a paid mutator transaction binding the contract method 0x297bb70b.
//
// Solidity: function batchFillOrders([]Struct0 orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeTransactor) BatchFillOrders(opts *bind.TransactOpts, orders []Struct0, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchFillOrders", orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrders is a paid mutator transaction binding the contract method 0x297bb70b.
//
// Solidity: function batchFillOrders([]Struct0 orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeSession) BatchFillOrders(orders []Struct0, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchFillOrders(&_Exchange.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrders is a paid mutator transaction binding the contract method 0x297bb70b.
//
// Solidity: function batchFillOrders([]Struct0 orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeTransactorSession) BatchFillOrders(orders []Struct0, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchFillOrders(&_Exchange.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrdersNoThrow is a paid mutator transaction binding the contract method 0x50dde190.
//
// Solidity: function batchFillOrdersNoThrow([]Struct0 orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeTransactor) BatchFillOrdersNoThrow(opts *bind.TransactOpts, orders []Struct0, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchFillOrdersNoThrow", orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrdersNoThrow is a paid mutator transaction binding the contract method 0x50dde190.
//
// Solidity: function batchFillOrdersNoThrow([]Struct0 orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeSession) BatchFillOrdersNoThrow(orders []Struct0, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchFillOrdersNoThrow(&_Exchange.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrdersNoThrow is a paid mutator transaction binding the contract method 0x50dde190.
//
// Solidity: function batchFillOrdersNoThrow([]Struct0 orders, uint256[] takerAssetFillAmounts, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeTransactorSession) BatchFillOrdersNoThrow(orders []Struct0, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchFillOrdersNoThrow(&_Exchange.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// CancelOrder is a paid mutator transaction binding the contract method 0xd46b02c3.
//
// Solidity: function cancelOrder(Struct0 order) returns()
func (_Exchange *ExchangeTransactor) CancelOrder(opts *bind.TransactOpts, order Struct0) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "cancelOrder", order)
}

// CancelOrder is a paid mutator transaction binding the contract method 0xd46b02c3.
//
// Solidity: function cancelOrder(Struct0 order) returns()
func (_Exchange *ExchangeSession) CancelOrder(order Struct0) (*types.Transaction, error) {
	return _Exchange.Contract.CancelOrder(&_Exchange.TransactOpts, order)
}

// CancelOrder is a paid mutator transaction binding the contract method 0xd46b02c3.
//
// Solidity: function cancelOrder(Struct0 order) returns()
func (_Exchange *ExchangeTransactorSession) CancelOrder(order Struct0) (*types.Transaction, error) {
	return _Exchange.Contract.CancelOrder(&_Exchange.TransactOpts, order)
}

// CancelOrdersUpTo is a paid mutator transaction binding the contract method 0x4f9559b1.
//
// Solidity: function cancelOrdersUpTo(uint256 targetOrderEpoch) returns()
func (_Exchange *ExchangeTransactor) CancelOrdersUpTo(opts *bind.TransactOpts, targetOrderEpoch *big.Int) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "cancelOrdersUpTo", targetOrderEpoch)
}

// CancelOrdersUpTo is a paid mutator transaction binding the contract method 0x4f9559b1.
//
// Solidity: function cancelOrdersUpTo(uint256 targetOrderEpoch) returns()
func (_Exchange *ExchangeSession) CancelOrdersUpTo(targetOrderEpoch *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.CancelOrdersUpTo(&_Exchange.TransactOpts, targetOrderEpoch)
}

// CancelOrdersUpTo is a paid mutator transaction binding the contract method 0x4f9559b1.
//
// Solidity: function cancelOrdersUpTo(uint256 targetOrderEpoch) returns()
func (_Exchange *ExchangeTransactorSession) CancelOrdersUpTo(targetOrderEpoch *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.CancelOrdersUpTo(&_Exchange.TransactOpts, targetOrderEpoch)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0xbfc8bfce.
//
// Solidity: function executeTransaction(uint256 salt, address signerAddress, bytes data, bytes signature) returns()
func (_Exchange *ExchangeTransactor) ExecuteTransaction(opts *bind.TransactOpts, salt *big.Int, signerAddress common.Address, data []byte, signature []byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "executeTransaction", salt, signerAddress, data, signature)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0xbfc8bfce.
//
// Solidity: function executeTransaction(uint256 salt, address signerAddress, bytes data, bytes signature) returns()
func (_Exchange *ExchangeSession) ExecuteTransaction(salt *big.Int, signerAddress common.Address, data []byte, signature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.ExecuteTransaction(&_Exchange.TransactOpts, salt, signerAddress, data, signature)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0xbfc8bfce.
//
// Solidity: function executeTransaction(uint256 salt, address signerAddress, bytes data, bytes signature) returns()
func (_Exchange *ExchangeTransactorSession) ExecuteTransaction(salt *big.Int, signerAddress common.Address, data []byte, signature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.ExecuteTransaction(&_Exchange.TransactOpts, salt, signerAddress, data, signature)
}

// FillOrKillOrder is a paid mutator transaction binding the contract method 0x64a3bc15.
//
// Solidity: function fillOrKillOrder(Struct0 order, uint256 takerAssetFillAmount, bytes signature) returns(Struct3 fillResults)
func (_Exchange *ExchangeTransactor) FillOrKillOrder(opts *bind.TransactOpts, order Struct0, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "fillOrKillOrder", order, takerAssetFillAmount, signature)
}

// FillOrKillOrder is a paid mutator transaction binding the contract method 0x64a3bc15.
//
// Solidity: function fillOrKillOrder(Struct0 order, uint256 takerAssetFillAmount, bytes signature) returns(Struct3 fillResults)
func (_Exchange *ExchangeSession) FillOrKillOrder(order Struct0, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.FillOrKillOrder(&_Exchange.TransactOpts, order, takerAssetFillAmount, signature)
}

// FillOrKillOrder is a paid mutator transaction binding the contract method 0x64a3bc15.
//
// Solidity: function fillOrKillOrder(Struct0 order, uint256 takerAssetFillAmount, bytes signature) returns(Struct3 fillResults)
func (_Exchange *ExchangeTransactorSession) FillOrKillOrder(order Struct0, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.FillOrKillOrder(&_Exchange.TransactOpts, order, takerAssetFillAmount, signature)
}

// FillOrder is a paid mutator transaction binding the contract method 0xb4be83d5.
//
// Solidity: function fillOrder(Struct0 order, uint256 takerAssetFillAmount, bytes signature) returns(Struct3 fillResults)
func (_Exchange *ExchangeTransactor) FillOrder(opts *bind.TransactOpts, order Struct0, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "fillOrder", order, takerAssetFillAmount, signature)
}

// FillOrder is a paid mutator transaction binding the contract method 0xb4be83d5.
//
// Solidity: function fillOrder(Struct0 order, uint256 takerAssetFillAmount, bytes signature) returns(Struct3 fillResults)
func (_Exchange *ExchangeSession) FillOrder(order Struct0, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.FillOrder(&_Exchange.TransactOpts, order, takerAssetFillAmount, signature)
}

// FillOrder is a paid mutator transaction binding the contract method 0xb4be83d5.
//
// Solidity: function fillOrder(Struct0 order, uint256 takerAssetFillAmount, bytes signature) returns(Struct3 fillResults)
func (_Exchange *ExchangeTransactorSession) FillOrder(order Struct0, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.FillOrder(&_Exchange.TransactOpts, order, takerAssetFillAmount, signature)
}

// FillOrderNoThrow is a paid mutator transaction binding the contract method 0x3e228bae.
//
// Solidity: function fillOrderNoThrow(Struct0 order, uint256 takerAssetFillAmount, bytes signature) returns(Struct3 fillResults)
func (_Exchange *ExchangeTransactor) FillOrderNoThrow(opts *bind.TransactOpts, order Struct0, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "fillOrderNoThrow", order, takerAssetFillAmount, signature)
}

// FillOrderNoThrow is a paid mutator transaction binding the contract method 0x3e228bae.
//
// Solidity: function fillOrderNoThrow(Struct0 order, uint256 takerAssetFillAmount, bytes signature) returns(Struct3 fillResults)
func (_Exchange *ExchangeSession) FillOrderNoThrow(order Struct0, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.FillOrderNoThrow(&_Exchange.TransactOpts, order, takerAssetFillAmount, signature)
}

// FillOrderNoThrow is a paid mutator transaction binding the contract method 0x3e228bae.
//
// Solidity: function fillOrderNoThrow(Struct0 order, uint256 takerAssetFillAmount, bytes signature) returns(Struct3 fillResults)
func (_Exchange *ExchangeTransactorSession) FillOrderNoThrow(order Struct0, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.FillOrderNoThrow(&_Exchange.TransactOpts, order, takerAssetFillAmount, signature)
}

// MarketBuyOrders is a paid mutator transaction binding the contract method 0xe5fa431b.
//
// Solidity: function marketBuyOrders([]Struct0 orders, uint256 makerAssetFillAmount, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeTransactor) MarketBuyOrders(opts *bind.TransactOpts, orders []Struct0, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "marketBuyOrders", orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrders is a paid mutator transaction binding the contract method 0xe5fa431b.
//
// Solidity: function marketBuyOrders([]Struct0 orders, uint256 makerAssetFillAmount, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeSession) MarketBuyOrders(orders []Struct0, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.MarketBuyOrders(&_Exchange.TransactOpts, orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrders is a paid mutator transaction binding the contract method 0xe5fa431b.
//
// Solidity: function marketBuyOrders([]Struct0 orders, uint256 makerAssetFillAmount, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeTransactorSession) MarketBuyOrders(orders []Struct0, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.MarketBuyOrders(&_Exchange.TransactOpts, orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersNoThrow is a paid mutator transaction binding the contract method 0xa3e20380.
//
// Solidity: function marketBuyOrdersNoThrow([]Struct0 orders, uint256 makerAssetFillAmount, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeTransactor) MarketBuyOrdersNoThrow(opts *bind.TransactOpts, orders []Struct0, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "marketBuyOrdersNoThrow", orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersNoThrow is a paid mutator transaction binding the contract method 0xa3e20380.
//
// Solidity: function marketBuyOrdersNoThrow([]Struct0 orders, uint256 makerAssetFillAmount, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeSession) MarketBuyOrdersNoThrow(orders []Struct0, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.MarketBuyOrdersNoThrow(&_Exchange.TransactOpts, orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersNoThrow is a paid mutator transaction binding the contract method 0xa3e20380.
//
// Solidity: function marketBuyOrdersNoThrow([]Struct0 orders, uint256 makerAssetFillAmount, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeTransactorSession) MarketBuyOrdersNoThrow(orders []Struct0, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.MarketBuyOrdersNoThrow(&_Exchange.TransactOpts, orders, makerAssetFillAmount, signatures)
}

// MarketSellOrders is a paid mutator transaction binding the contract method 0x7e1d9808.
//
// Solidity: function marketSellOrders([]Struct0 orders, uint256 takerAssetFillAmount, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeTransactor) MarketSellOrders(opts *bind.TransactOpts, orders []Struct0, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "marketSellOrders", orders, takerAssetFillAmount, signatures)
}

// MarketSellOrders is a paid mutator transaction binding the contract method 0x7e1d9808.
//
// Solidity: function marketSellOrders([]Struct0 orders, uint256 takerAssetFillAmount, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeSession) MarketSellOrders(orders []Struct0, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.MarketSellOrders(&_Exchange.TransactOpts, orders, takerAssetFillAmount, signatures)
}

// MarketSellOrders is a paid mutator transaction binding the contract method 0x7e1d9808.
//
// Solidity: function marketSellOrders([]Struct0 orders, uint256 takerAssetFillAmount, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeTransactorSession) MarketSellOrders(orders []Struct0, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.MarketSellOrders(&_Exchange.TransactOpts, orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersNoThrow is a paid mutator transaction binding the contract method 0xdd1c7d18.
//
// Solidity: function marketSellOrdersNoThrow([]Struct0 orders, uint256 takerAssetFillAmount, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeTransactor) MarketSellOrdersNoThrow(opts *bind.TransactOpts, orders []Struct0, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "marketSellOrdersNoThrow", orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersNoThrow is a paid mutator transaction binding the contract method 0xdd1c7d18.
//
// Solidity: function marketSellOrdersNoThrow([]Struct0 orders, uint256 takerAssetFillAmount, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeSession) MarketSellOrdersNoThrow(orders []Struct0, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.MarketSellOrdersNoThrow(&_Exchange.TransactOpts, orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersNoThrow is a paid mutator transaction binding the contract method 0xdd1c7d18.
//
// Solidity: function marketSellOrdersNoThrow([]Struct0 orders, uint256 takerAssetFillAmount, bytes[] signatures) returns(Struct3 totalFillResults)
func (_Exchange *ExchangeTransactorSession) MarketSellOrdersNoThrow(orders []Struct0, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.MarketSellOrdersNoThrow(&_Exchange.TransactOpts, orders, takerAssetFillAmount, signatures)
}

// MatchOrders is a paid mutator transaction binding the contract method 0x3c28d861.
//
// Solidity: function matchOrders(Struct0 leftOrder, Struct0 rightOrder, bytes leftSignature, bytes rightSignature) returns(Struct2 matchedFillResults)
func (_Exchange *ExchangeTransactor) MatchOrders(opts *bind.TransactOpts, leftOrder Struct0, rightOrder Struct0, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "matchOrders", leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrders is a paid mutator transaction binding the contract method 0x3c28d861.
//
// Solidity: function matchOrders(Struct0 leftOrder, Struct0 rightOrder, bytes leftSignature, bytes rightSignature) returns(Struct2 matchedFillResults)
func (_Exchange *ExchangeSession) MatchOrders(leftOrder Struct0, rightOrder Struct0, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.MatchOrders(&_Exchange.TransactOpts, leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrders is a paid mutator transaction binding the contract method 0x3c28d861.
//
// Solidity: function matchOrders(Struct0 leftOrder, Struct0 rightOrder, bytes leftSignature, bytes rightSignature) returns(Struct2 matchedFillResults)
func (_Exchange *ExchangeTransactorSession) MatchOrders(leftOrder Struct0, rightOrder Struct0, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.MatchOrders(&_Exchange.TransactOpts, leftOrder, rightOrder, leftSignature, rightSignature)
}

// PreSign is a paid mutator transaction binding the contract method 0x3683ef8e.
//
// Solidity: function preSign(bytes32 hash, address signerAddress, bytes signature) returns()
func (_Exchange *ExchangeTransactor) PreSign(opts *bind.TransactOpts, hash [32]byte, signerAddress common.Address, signature []byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "preSign", hash, signerAddress, signature)
}

// PreSign is a paid mutator transaction binding the contract method 0x3683ef8e.
//
// Solidity: function preSign(bytes32 hash, address signerAddress, bytes signature) returns()
func (_Exchange *ExchangeSession) PreSign(hash [32]byte, signerAddress common.Address, signature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.PreSign(&_Exchange.TransactOpts, hash, signerAddress, signature)
}

// PreSign is a paid mutator transaction binding the contract method 0x3683ef8e.
//
// Solidity: function preSign(bytes32 hash, address signerAddress, bytes signature) returns()
func (_Exchange *ExchangeTransactorSession) PreSign(hash [32]byte, signerAddress common.Address, signature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.PreSign(&_Exchange.TransactOpts, hash, signerAddress, signature)
}

// RegisterAssetProxy is a paid mutator transaction binding the contract method 0xc585bb93.
//
// Solidity: function registerAssetProxy(address assetProxy) returns()
func (_Exchange *ExchangeTransactor) RegisterAssetProxy(opts *bind.TransactOpts, assetProxy common.Address) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "registerAssetProxy", assetProxy)
}

// RegisterAssetProxy is a paid mutator transaction binding the contract method 0xc585bb93.
//
// Solidity: function registerAssetProxy(address assetProxy) returns()
func (_Exchange *ExchangeSession) RegisterAssetProxy(assetProxy common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.RegisterAssetProxy(&_Exchange.TransactOpts, assetProxy)
}

// RegisterAssetProxy is a paid mutator transaction binding the contract method 0xc585bb93.
//
// Solidity: function registerAssetProxy(address assetProxy) returns()
func (_Exchange *ExchangeTransactorSession) RegisterAssetProxy(assetProxy common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.RegisterAssetProxy(&_Exchange.TransactOpts, assetProxy)
}

// SetSignatureValidatorApproval is a paid mutator transaction binding the contract method 0x77fcce68.
//
// Solidity: function setSignatureValidatorApproval(address validatorAddress, bool approval) returns()
func (_Exchange *ExchangeTransactor) SetSignatureValidatorApproval(opts *bind.TransactOpts, validatorAddress common.Address, approval bool) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "setSignatureValidatorApproval", validatorAddress, approval)
}

// SetSignatureValidatorApproval is a paid mutator transaction binding the contract method 0x77fcce68.
//
// Solidity: function setSignatureValidatorApproval(address validatorAddress, bool approval) returns()
func (_Exchange *ExchangeSession) SetSignatureValidatorApproval(validatorAddress common.Address, approval bool) (*types.Transaction, error) {
	return _Exchange.Contract.SetSignatureValidatorApproval(&_Exchange.TransactOpts, validatorAddress, approval)
}

// SetSignatureValidatorApproval is a paid mutator transaction binding the contract method 0x77fcce68.
//
// Solidity: function setSignatureValidatorApproval(address validatorAddress, bool approval) returns()
func (_Exchange *ExchangeTransactorSession) SetSignatureValidatorApproval(validatorAddress common.Address, approval bool) (*types.Transaction, error) {
	return _Exchange.Contract.SetSignatureValidatorApproval(&_Exchange.TransactOpts, validatorAddress, approval)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Exchange *ExchangeTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Exchange *ExchangeSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.TransferOwnership(&_Exchange.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Exchange *ExchangeTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.TransferOwnership(&_Exchange.TransactOpts, newOwner)
}

// ExchangeAssetProxyRegisteredIterator is returned from FilterAssetProxyRegistered and is used to iterate over the raw logs and unpacked data for AssetProxyRegistered events raised by the Exchange contract.
type ExchangeAssetProxyRegisteredIterator struct {
	Event *ExchangeAssetProxyRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangeAssetProxyRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeAssetProxyRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangeAssetProxyRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangeAssetProxyRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeAssetProxyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeAssetProxyRegistered represents a AssetProxyRegistered event raised by the Exchange contract.
type ExchangeAssetProxyRegistered struct {
	Id         [4]byte
	AssetProxy common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAssetProxyRegistered is a free log retrieval operation binding the contract event 0xd2c6b762299c609bdb96520b58a49bfb80186934d4f71a86a367571a15c03194.
//
// Solidity: event AssetProxyRegistered(bytes4 id, address assetProxy)
func (_Exchange *ExchangeFilterer) FilterAssetProxyRegistered(opts *bind.FilterOpts) (*ExchangeAssetProxyRegisteredIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "AssetProxyRegistered")
	if err != nil {
		return nil, err
	}
	return &ExchangeAssetProxyRegisteredIterator{contract: _Exchange.contract, event: "AssetProxyRegistered", logs: logs, sub: sub}, nil
}

// WatchAssetProxyRegistered is a free log subscription operation binding the contract event 0xd2c6b762299c609bdb96520b58a49bfb80186934d4f71a86a367571a15c03194.
//
// Solidity: event AssetProxyRegistered(bytes4 id, address assetProxy)
func (_Exchange *ExchangeFilterer) WatchAssetProxyRegistered(opts *bind.WatchOpts, sink chan<- *ExchangeAssetProxyRegistered) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "AssetProxyRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeAssetProxyRegistered)
				if err := _Exchange.contract.UnpackLog(event, "AssetProxyRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAssetProxyRegistered is a log parse operation binding the contract event 0xd2c6b762299c609bdb96520b58a49bfb80186934d4f71a86a367571a15c03194.
//
// Solidity: event AssetProxyRegistered(bytes4 id, address assetProxy)
func (_Exchange *ExchangeFilterer) ParseAssetProxyRegistered(log types.Log) (*ExchangeAssetProxyRegistered, error) {
	event := new(ExchangeAssetProxyRegistered)
	if err := _Exchange.contract.UnpackLog(event, "AssetProxyRegistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangeCancelIterator is returned from FilterCancel and is used to iterate over the raw logs and unpacked data for Cancel events raised by the Exchange contract.
type ExchangeCancelIterator struct {
	Event *ExchangeCancel // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangeCancelIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeCancel)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangeCancel)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangeCancelIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeCancelIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeCancel represents a Cancel event raised by the Exchange contract.
type ExchangeCancel struct {
	MakerAddress        common.Address
	FeeRecipientAddress common.Address
	SenderAddress       common.Address
	OrderHash           [32]byte
	MakerAssetData      []byte
	TakerAssetData      []byte
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterCancel is a free log retrieval operation binding the contract event 0xdc47b3613d9fe400085f6dbdc99453462279057e6207385042827ed6b1a62cf7.
//
// Solidity: event Cancel(address indexed makerAddress, address indexed feeRecipientAddress, address senderAddress, bytes32 indexed orderHash, bytes makerAssetData, bytes takerAssetData)
func (_Exchange *ExchangeFilterer) FilterCancel(opts *bind.FilterOpts, makerAddress []common.Address, feeRecipientAddress []common.Address, orderHash [][32]byte) (*ExchangeCancelIterator, error) {

	var makerAddressRule []interface{}
	for _, makerAddressItem := range makerAddress {
		makerAddressRule = append(makerAddressRule, makerAddressItem)
	}
	var feeRecipientAddressRule []interface{}
	for _, feeRecipientAddressItem := range feeRecipientAddress {
		feeRecipientAddressRule = append(feeRecipientAddressRule, feeRecipientAddressItem)
	}

	var orderHashRule []interface{}
	for _, orderHashItem := range orderHash {
		orderHashRule = append(orderHashRule, orderHashItem)
	}

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "Cancel", makerAddressRule, feeRecipientAddressRule, orderHashRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeCancelIterator{contract: _Exchange.contract, event: "Cancel", logs: logs, sub: sub}, nil
}

// WatchCancel is a free log subscription operation binding the contract event 0xdc47b3613d9fe400085f6dbdc99453462279057e6207385042827ed6b1a62cf7.
//
// Solidity: event Cancel(address indexed makerAddress, address indexed feeRecipientAddress, address senderAddress, bytes32 indexed orderHash, bytes makerAssetData, bytes takerAssetData)
func (_Exchange *ExchangeFilterer) WatchCancel(opts *bind.WatchOpts, sink chan<- *ExchangeCancel, makerAddress []common.Address, feeRecipientAddress []common.Address, orderHash [][32]byte) (event.Subscription, error) {

	var makerAddressRule []interface{}
	for _, makerAddressItem := range makerAddress {
		makerAddressRule = append(makerAddressRule, makerAddressItem)
	}
	var feeRecipientAddressRule []interface{}
	for _, feeRecipientAddressItem := range feeRecipientAddress {
		feeRecipientAddressRule = append(feeRecipientAddressRule, feeRecipientAddressItem)
	}

	var orderHashRule []interface{}
	for _, orderHashItem := range orderHash {
		orderHashRule = append(orderHashRule, orderHashItem)
	}

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "Cancel", makerAddressRule, feeRecipientAddressRule, orderHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeCancel)
				if err := _Exchange.contract.UnpackLog(event, "Cancel", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCancel is a log parse operation binding the contract event 0xdc47b3613d9fe400085f6dbdc99453462279057e6207385042827ed6b1a62cf7.
//
// Solidity: event Cancel(address indexed makerAddress, address indexed feeRecipientAddress, address senderAddress, bytes32 indexed orderHash, bytes makerAssetData, bytes takerAssetData)
func (_Exchange *ExchangeFilterer) ParseCancel(log types.Log) (*ExchangeCancel, error) {
	event := new(ExchangeCancel)
	if err := _Exchange.contract.UnpackLog(event, "Cancel", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangeCancelUpToIterator is returned from FilterCancelUpTo and is used to iterate over the raw logs and unpacked data for CancelUpTo events raised by the Exchange contract.
type ExchangeCancelUpToIterator struct {
	Event *ExchangeCancelUpTo // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangeCancelUpToIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeCancelUpTo)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangeCancelUpTo)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangeCancelUpToIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeCancelUpToIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeCancelUpTo represents a CancelUpTo event raised by the Exchange contract.
type ExchangeCancelUpTo struct {
	MakerAddress  common.Address
	SenderAddress common.Address
	OrderEpoch    *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterCancelUpTo is a free log retrieval operation binding the contract event 0x82af639571738f4ebd4268fb0363d8957ebe1bbb9e78dba5ebd69eed39b154f0.
//
// Solidity: event CancelUpTo(address indexed makerAddress, address indexed senderAddress, uint256 orderEpoch)
func (_Exchange *ExchangeFilterer) FilterCancelUpTo(opts *bind.FilterOpts, makerAddress []common.Address, senderAddress []common.Address) (*ExchangeCancelUpToIterator, error) {

	var makerAddressRule []interface{}
	for _, makerAddressItem := range makerAddress {
		makerAddressRule = append(makerAddressRule, makerAddressItem)
	}
	var senderAddressRule []interface{}
	for _, senderAddressItem := range senderAddress {
		senderAddressRule = append(senderAddressRule, senderAddressItem)
	}

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "CancelUpTo", makerAddressRule, senderAddressRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeCancelUpToIterator{contract: _Exchange.contract, event: "CancelUpTo", logs: logs, sub: sub}, nil
}

// WatchCancelUpTo is a free log subscription operation binding the contract event 0x82af639571738f4ebd4268fb0363d8957ebe1bbb9e78dba5ebd69eed39b154f0.
//
// Solidity: event CancelUpTo(address indexed makerAddress, address indexed senderAddress, uint256 orderEpoch)
func (_Exchange *ExchangeFilterer) WatchCancelUpTo(opts *bind.WatchOpts, sink chan<- *ExchangeCancelUpTo, makerAddress []common.Address, senderAddress []common.Address) (event.Subscription, error) {

	var makerAddressRule []interface{}
	for _, makerAddressItem := range makerAddress {
		makerAddressRule = append(makerAddressRule, makerAddressItem)
	}
	var senderAddressRule []interface{}
	for _, senderAddressItem := range senderAddress {
		senderAddressRule = append(senderAddressRule, senderAddressItem)
	}

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "CancelUpTo", makerAddressRule, senderAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeCancelUpTo)
				if err := _Exchange.contract.UnpackLog(event, "CancelUpTo", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCancelUpTo is a log parse operation binding the contract event 0x82af639571738f4ebd4268fb0363d8957ebe1bbb9e78dba5ebd69eed39b154f0.
//
// Solidity: event CancelUpTo(address indexed makerAddress, address indexed senderAddress, uint256 orderEpoch)
func (_Exchange *ExchangeFilterer) ParseCancelUpTo(log types.Log) (*ExchangeCancelUpTo, error) {
	event := new(ExchangeCancelUpTo)
	if err := _Exchange.contract.UnpackLog(event, "CancelUpTo", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangeFillIterator is returned from FilterFill and is used to iterate over the raw logs and unpacked data for Fill events raised by the Exchange contract.
type ExchangeFillIterator struct {
	Event *ExchangeFill // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangeFillIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeFill)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangeFill)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangeFillIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeFillIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeFill represents a Fill event raised by the Exchange contract.
type ExchangeFill struct {
	MakerAddress           common.Address
	FeeRecipientAddress    common.Address
	TakerAddress           common.Address
	SenderAddress          common.Address
	MakerAssetFilledAmount *big.Int
	TakerAssetFilledAmount *big.Int
	MakerFeePaid           *big.Int
	TakerFeePaid           *big.Int
	OrderHash              [32]byte
	MakerAssetData         []byte
	TakerAssetData         []byte
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterFill is a free log retrieval operation binding the contract event 0x0bcc4c97732e47d9946f229edb95f5b6323f601300e4690de719993f3c371129.
//
// Solidity: event Fill(address indexed makerAddress, address indexed feeRecipientAddress, address takerAddress, address senderAddress, uint256 makerAssetFilledAmount, uint256 takerAssetFilledAmount, uint256 makerFeePaid, uint256 takerFeePaid, bytes32 indexed orderHash, bytes makerAssetData, bytes takerAssetData)
func (_Exchange *ExchangeFilterer) FilterFill(opts *bind.FilterOpts, makerAddress []common.Address, feeRecipientAddress []common.Address, orderHash [][32]byte) (*ExchangeFillIterator, error) {

	var makerAddressRule []interface{}
	for _, makerAddressItem := range makerAddress {
		makerAddressRule = append(makerAddressRule, makerAddressItem)
	}
	var feeRecipientAddressRule []interface{}
	for _, feeRecipientAddressItem := range feeRecipientAddress {
		feeRecipientAddressRule = append(feeRecipientAddressRule, feeRecipientAddressItem)
	}

	var orderHashRule []interface{}
	for _, orderHashItem := range orderHash {
		orderHashRule = append(orderHashRule, orderHashItem)
	}

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "Fill", makerAddressRule, feeRecipientAddressRule, orderHashRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeFillIterator{contract: _Exchange.contract, event: "Fill", logs: logs, sub: sub}, nil
}

// WatchFill is a free log subscription operation binding the contract event 0x0bcc4c97732e47d9946f229edb95f5b6323f601300e4690de719993f3c371129.
//
// Solidity: event Fill(address indexed makerAddress, address indexed feeRecipientAddress, address takerAddress, address senderAddress, uint256 makerAssetFilledAmount, uint256 takerAssetFilledAmount, uint256 makerFeePaid, uint256 takerFeePaid, bytes32 indexed orderHash, bytes makerAssetData, bytes takerAssetData)
func (_Exchange *ExchangeFilterer) WatchFill(opts *bind.WatchOpts, sink chan<- *ExchangeFill, makerAddress []common.Address, feeRecipientAddress []common.Address, orderHash [][32]byte) (event.Subscription, error) {

	var makerAddressRule []interface{}
	for _, makerAddressItem := range makerAddress {
		makerAddressRule = append(makerAddressRule, makerAddressItem)
	}
	var feeRecipientAddressRule []interface{}
	for _, feeRecipientAddressItem := range feeRecipientAddress {
		feeRecipientAddressRule = append(feeRecipientAddressRule, feeRecipientAddressItem)
	}

	var orderHashRule []interface{}
	for _, orderHashItem := range orderHash {
		orderHashRule = append(orderHashRule, orderHashItem)
	}

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "Fill", makerAddressRule, feeRecipientAddressRule, orderHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeFill)
				if err := _Exchange.contract.UnpackLog(event, "Fill", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFill is a log parse operation binding the contract event 0x0bcc4c97732e47d9946f229edb95f5b6323f601300e4690de719993f3c371129.
//
// Solidity: event Fill(address indexed makerAddress, address indexed feeRecipientAddress, address takerAddress, address senderAddress, uint256 makerAssetFilledAmount, uint256 takerAssetFilledAmount, uint256 makerFeePaid, uint256 takerFeePaid, bytes32 indexed orderHash, bytes makerAssetData, bytes takerAssetData)
func (_Exchange *ExchangeFilterer) ParseFill(log types.Log) (*ExchangeFill, error) {
	event := new(ExchangeFill)
	if err := _Exchange.contract.UnpackLog(event, "Fill", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangeSignatureValidatorApprovalIterator is returned from FilterSignatureValidatorApproval and is used to iterate over the raw logs and unpacked data for SignatureValidatorApproval events raised by the Exchange contract.
type ExchangeSignatureValidatorApprovalIterator struct {
	Event *ExchangeSignatureValidatorApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ExchangeSignatureValidatorApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeSignatureValidatorApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ExchangeSignatureValidatorApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ExchangeSignatureValidatorApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeSignatureValidatorApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeSignatureValidatorApproval represents a SignatureValidatorApproval event raised by the Exchange contract.
type ExchangeSignatureValidatorApproval struct {
	SignerAddress    common.Address
	ValidatorAddress common.Address
	Approved         bool
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSignatureValidatorApproval is a free log retrieval operation binding the contract event 0xa8656e308026eeabce8f0bc18048433252318ab80ac79da0b3d3d8697dfba891.
//
// Solidity: event SignatureValidatorApproval(address indexed signerAddress, address indexed validatorAddress, bool approved)
func (_Exchange *ExchangeFilterer) FilterSignatureValidatorApproval(opts *bind.FilterOpts, signerAddress []common.Address, validatorAddress []common.Address) (*ExchangeSignatureValidatorApprovalIterator, error) {

	var signerAddressRule []interface{}
	for _, signerAddressItem := range signerAddress {
		signerAddressRule = append(signerAddressRule, signerAddressItem)
	}
	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "SignatureValidatorApproval", signerAddressRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeSignatureValidatorApprovalIterator{contract: _Exchange.contract, event: "SignatureValidatorApproval", logs: logs, sub: sub}, nil
}

// WatchSignatureValidatorApproval is a free log subscription operation binding the contract event 0xa8656e308026eeabce8f0bc18048433252318ab80ac79da0b3d3d8697dfba891.
//
// Solidity: event SignatureValidatorApproval(address indexed signerAddress, address indexed validatorAddress, bool approved)
func (_Exchange *ExchangeFilterer) WatchSignatureValidatorApproval(opts *bind.WatchOpts, sink chan<- *ExchangeSignatureValidatorApproval, signerAddress []common.Address, validatorAddress []common.Address) (event.Subscription, error) {

	var signerAddressRule []interface{}
	for _, signerAddressItem := range signerAddress {
		signerAddressRule = append(signerAddressRule, signerAddressItem)
	}
	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "SignatureValidatorApproval", signerAddressRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeSignatureValidatorApproval)
				if err := _Exchange.contract.UnpackLog(event, "SignatureValidatorApproval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSignatureValidatorApproval is a log parse operation binding the contract event 0xa8656e308026eeabce8f0bc18048433252318ab80ac79da0b3d3d8697dfba891.
//
// Solidity: event SignatureValidatorApproval(address indexed signerAddress, address indexed validatorAddress, bool approved)
func (_Exchange *ExchangeFilterer) ParseSignatureValidatorApproval(log types.Log) (*ExchangeSignatureValidatorApproval, error) {
	event := new(ExchangeSignatureValidatorApproval)
	if err := _Exchange.contract.UnpackLog(event, "SignatureValidatorApproval", log); err != nil {
		return nil, err
	}
	return event, nil
}
