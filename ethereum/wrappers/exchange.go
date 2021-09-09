// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package wrappers

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// LibFillResultsBatchMatchedFillResults is an auto generated low-level Go binding around an user-defined struct.
type LibFillResultsBatchMatchedFillResults struct {
	Left                    []LibFillResultsFillResults
	Right                   []LibFillResultsFillResults
	ProfitInLeftMakerAsset  *big.Int
	ProfitInRightMakerAsset *big.Int
}

// LibFillResultsFillResults is an auto generated low-level Go binding around an user-defined struct.
type LibFillResultsFillResults struct {
	MakerAssetFilledAmount *big.Int
	TakerAssetFilledAmount *big.Int
	MakerFeePaid           *big.Int
	TakerFeePaid           *big.Int
	ProtocolFeePaid        *big.Int
}

// LibFillResultsMatchedFillResults is an auto generated low-level Go binding around an user-defined struct.
type LibFillResultsMatchedFillResults struct {
	Left                    LibFillResultsFillResults
	Right                   LibFillResultsFillResults
	ProfitInLeftMakerAsset  *big.Int
	ProfitInRightMakerAsset *big.Int
}

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

// ExchangeMetaData contains all meta data concerning the Exchange contract.
var ExchangeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"id\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetProxy\",\"type\":\"address\"}],\"name\":\"AssetProxyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"Cancel\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"orderSenderAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"orderEpoch\",\"type\":\"uint256\"}],\"name\":\"CancelUpTo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"name\":\"Fill\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldProtocolFeeCollector\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"updatedProtocolFeeCollector\",\"type\":\"address\"}],\"name\":\"ProtocolFeeCollectorAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldProtocolFeeMultiplier\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedProtocolFeeMultiplier\",\"type\":\"uint256\"}],\"name\":\"ProtocolFeeMultiplier\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isApproved\",\"type\":\"bool\"}],\"name\":\"SignatureValidatorApproval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"transactionHash\",\"type\":\"bytes32\"}],\"name\":\"TransactionExecution\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"EIP1271_MAGIC_VALUE\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"EIP712_EXCHANGE_DOMAIN_HASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedValidators\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"}],\"name\":\"batchCancelOrders\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLibZeroExTransaction.ZeroExTransaction[]\",\"name\":\"transactions\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchExecuteTransactions\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"returnData\",\"type\":\"bytes[]\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchFillOrKillOrders\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"fillResults\",\"type\":\"tuple[]\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchFillOrders\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"fillResults\",\"type\":\"tuple[]\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchFillOrdersNoThrow\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"fillResults\",\"type\":\"tuple[]\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"leftOrders\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"rightOrders\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"leftSignatures\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes[]\",\"name\":\"rightSignatures\",\"type\":\"bytes[]\"}],\"name\":\"batchMatchOrders\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"left\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"right\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"profitInLeftMakerAsset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"profitInRightMakerAsset\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.BatchMatchedFillResults\",\"name\":\"batchMatchedFillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"leftOrders\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"rightOrders\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"leftSignatures\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes[]\",\"name\":\"rightSignatures\",\"type\":\"bytes[]\"}],\"name\":\"batchMatchOrdersWithMaximalFill\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"left\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"right\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"profitInLeftMakerAsset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"profitInRightMakerAsset\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.BatchMatchedFillResults\",\"name\":\"batchMatchedFillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"cancelOrder\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"targetOrderEpoch\",\"type\":\"uint256\"}],\"name\":\"cancelOrdersUpTo\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"cancelled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentContextAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"detachProtocolFeeCollector\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLibZeroExTransaction.ZeroExTransaction\",\"name\":\"transaction\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"executeTransaction\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"fillOrKillOrder\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"fillOrder\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"filled\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"}],\"name\":\"getAssetProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"assetProxy\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"getOrderInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"enumLibOrder.OrderStatus\",\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"internalType\":\"structLibOrder.OrderInfo\",\"name\":\"orderInfo\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"isValidHashSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isValid\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"isValidOrderSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isValid\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLibZeroExTransaction.ZeroExTransaction\",\"name\":\"transaction\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"isValidTransactionSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isValid\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketBuyOrdersFillOrKill\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketBuyOrdersNoThrow\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketSellOrdersFillOrKill\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketSellOrdersNoThrow\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"leftOrder\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"rightOrder\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"leftSignature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rightSignature\",\"type\":\"bytes\"}],\"name\":\"matchOrders\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"left\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"right\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"profitInLeftMakerAsset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"profitInRightMakerAsset\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.MatchedFillResults\",\"name\":\"matchedFillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"leftOrder\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"rightOrder\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"leftSignature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rightSignature\",\"type\":\"bytes\"}],\"name\":\"matchOrdersWithMaximalFill\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"left\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"right\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"profitInLeftMakerAsset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"profitInRightMakerAsset\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.MatchedFillResults\",\"name\":\"matchedFillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"orderEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"preSign\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"preSigned\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"protocolFeeCollector\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"protocolFeeMultiplier\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetProxy\",\"type\":\"address\"}],\"name\":\"registerAssetProxy\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"updatedProtocolFeeCollector\",\"type\":\"address\"}],\"name\":\"setProtocolFeeCollectorAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"updatedProtocolFeeMultiplier\",\"type\":\"uint256\"}],\"name\":\"setProtocolFeeMultiplier\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approval\",\"type\":\"bool\"}],\"name\":\"setSignatureValidatorApproval\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"assetData\",\"type\":\"bytes[]\"},{\"internalType\":\"address[]\",\"name\":\"fromAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"toAddresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"simulateDispatchTransferFromCalls\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"transactionsExecuted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ExchangeABI is the input ABI used to generate the binding from.
// Deprecated: Use ExchangeMetaData.ABI instead.
var ExchangeABI = ExchangeMetaData.ABI

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
func (_Exchange *ExchangeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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
func (_Exchange *ExchangeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
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

// EIP1271MAGICVALUE is a free data retrieval call binding the contract method 0xdd885e2d.
//
// Solidity: function EIP1271_MAGIC_VALUE() view returns(bytes4)
func (_Exchange *ExchangeCaller) EIP1271MAGICVALUE(opts *bind.CallOpts) ([4]byte, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "EIP1271_MAGIC_VALUE")

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

// EIP1271MAGICVALUE is a free data retrieval call binding the contract method 0xdd885e2d.
//
// Solidity: function EIP1271_MAGIC_VALUE() view returns(bytes4)
func (_Exchange *ExchangeSession) EIP1271MAGICVALUE() ([4]byte, error) {
	return _Exchange.Contract.EIP1271MAGICVALUE(&_Exchange.CallOpts)
}

// EIP1271MAGICVALUE is a free data retrieval call binding the contract method 0xdd885e2d.
//
// Solidity: function EIP1271_MAGIC_VALUE() view returns(bytes4)
func (_Exchange *ExchangeCallerSession) EIP1271MAGICVALUE() ([4]byte, error) {
	return _Exchange.Contract.EIP1271MAGICVALUE(&_Exchange.CallOpts)
}

// EIP712EXCHANGEDOMAINHASH is a free data retrieval call binding the contract method 0xc26cfecd.
//
// Solidity: function EIP712_EXCHANGE_DOMAIN_HASH() view returns(bytes32)
func (_Exchange *ExchangeCaller) EIP712EXCHANGEDOMAINHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "EIP712_EXCHANGE_DOMAIN_HASH")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EIP712EXCHANGEDOMAINHASH is a free data retrieval call binding the contract method 0xc26cfecd.
//
// Solidity: function EIP712_EXCHANGE_DOMAIN_HASH() view returns(bytes32)
func (_Exchange *ExchangeSession) EIP712EXCHANGEDOMAINHASH() ([32]byte, error) {
	return _Exchange.Contract.EIP712EXCHANGEDOMAINHASH(&_Exchange.CallOpts)
}

// EIP712EXCHANGEDOMAINHASH is a free data retrieval call binding the contract method 0xc26cfecd.
//
// Solidity: function EIP712_EXCHANGE_DOMAIN_HASH() view returns(bytes32)
func (_Exchange *ExchangeCallerSession) EIP712EXCHANGEDOMAINHASH() ([32]byte, error) {
	return _Exchange.Contract.EIP712EXCHANGEDOMAINHASH(&_Exchange.CallOpts)
}

// AllowedValidators is a free data retrieval call binding the contract method 0x7b8e3514.
//
// Solidity: function allowedValidators(address , address ) view returns(bool)
func (_Exchange *ExchangeCaller) AllowedValidators(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "allowedValidators", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AllowedValidators is a free data retrieval call binding the contract method 0x7b8e3514.
//
// Solidity: function allowedValidators(address , address ) view returns(bool)
func (_Exchange *ExchangeSession) AllowedValidators(arg0 common.Address, arg1 common.Address) (bool, error) {
	return _Exchange.Contract.AllowedValidators(&_Exchange.CallOpts, arg0, arg1)
}

// AllowedValidators is a free data retrieval call binding the contract method 0x7b8e3514.
//
// Solidity: function allowedValidators(address , address ) view returns(bool)
func (_Exchange *ExchangeCallerSession) AllowedValidators(arg0 common.Address, arg1 common.Address) (bool, error) {
	return _Exchange.Contract.AllowedValidators(&_Exchange.CallOpts, arg0, arg1)
}

// Cancelled is a free data retrieval call binding the contract method 0x2ac12622.
//
// Solidity: function cancelled(bytes32 ) view returns(bool)
func (_Exchange *ExchangeCaller) Cancelled(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "cancelled", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Cancelled is a free data retrieval call binding the contract method 0x2ac12622.
//
// Solidity: function cancelled(bytes32 ) view returns(bool)
func (_Exchange *ExchangeSession) Cancelled(arg0 [32]byte) (bool, error) {
	return _Exchange.Contract.Cancelled(&_Exchange.CallOpts, arg0)
}

// Cancelled is a free data retrieval call binding the contract method 0x2ac12622.
//
// Solidity: function cancelled(bytes32 ) view returns(bool)
func (_Exchange *ExchangeCallerSession) Cancelled(arg0 [32]byte) (bool, error) {
	return _Exchange.Contract.Cancelled(&_Exchange.CallOpts, arg0)
}

// CurrentContextAddress is a free data retrieval call binding the contract method 0xeea086ba.
//
// Solidity: function currentContextAddress() view returns(address)
func (_Exchange *ExchangeCaller) CurrentContextAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "currentContextAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CurrentContextAddress is a free data retrieval call binding the contract method 0xeea086ba.
//
// Solidity: function currentContextAddress() view returns(address)
func (_Exchange *ExchangeSession) CurrentContextAddress() (common.Address, error) {
	return _Exchange.Contract.CurrentContextAddress(&_Exchange.CallOpts)
}

// CurrentContextAddress is a free data retrieval call binding the contract method 0xeea086ba.
//
// Solidity: function currentContextAddress() view returns(address)
func (_Exchange *ExchangeCallerSession) CurrentContextAddress() (common.Address, error) {
	return _Exchange.Contract.CurrentContextAddress(&_Exchange.CallOpts)
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled(bytes32 ) view returns(uint256)
func (_Exchange *ExchangeCaller) Filled(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "filled", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled(bytes32 ) view returns(uint256)
func (_Exchange *ExchangeSession) Filled(arg0 [32]byte) (*big.Int, error) {
	return _Exchange.Contract.Filled(&_Exchange.CallOpts, arg0)
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled(bytes32 ) view returns(uint256)
func (_Exchange *ExchangeCallerSession) Filled(arg0 [32]byte) (*big.Int, error) {
	return _Exchange.Contract.Filled(&_Exchange.CallOpts, arg0)
}

// GetAssetProxy is a free data retrieval call binding the contract method 0x60704108.
//
// Solidity: function getAssetProxy(bytes4 assetProxyId) view returns(address assetProxy)
func (_Exchange *ExchangeCaller) GetAssetProxy(opts *bind.CallOpts, assetProxyId [4]byte) (common.Address, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "getAssetProxy", assetProxyId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAssetProxy is a free data retrieval call binding the contract method 0x60704108.
//
// Solidity: function getAssetProxy(bytes4 assetProxyId) view returns(address assetProxy)
func (_Exchange *ExchangeSession) GetAssetProxy(assetProxyId [4]byte) (common.Address, error) {
	return _Exchange.Contract.GetAssetProxy(&_Exchange.CallOpts, assetProxyId)
}

// GetAssetProxy is a free data retrieval call binding the contract method 0x60704108.
//
// Solidity: function getAssetProxy(bytes4 assetProxyId) view returns(address assetProxy)
func (_Exchange *ExchangeCallerSession) GetAssetProxy(assetProxyId [4]byte) (common.Address, error) {
	return _Exchange.Contract.GetAssetProxy(&_Exchange.CallOpts, assetProxyId)
}

// GetOrderInfo is a free data retrieval call binding the contract method 0x9d3fa4b9.
//
// Solidity: function getOrderInfo((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order) view returns((uint8,bytes32,uint256) orderInfo)
func (_Exchange *ExchangeCaller) GetOrderInfo(opts *bind.CallOpts, order LibOrderOrder) (LibOrderOrderInfo, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "getOrderInfo", order)

	if err != nil {
		return *new(LibOrderOrderInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(LibOrderOrderInfo)).(*LibOrderOrderInfo)

	return out0, err

}

// GetOrderInfo is a free data retrieval call binding the contract method 0x9d3fa4b9.
//
// Solidity: function getOrderInfo((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order) view returns((uint8,bytes32,uint256) orderInfo)
func (_Exchange *ExchangeSession) GetOrderInfo(order LibOrderOrder) (LibOrderOrderInfo, error) {
	return _Exchange.Contract.GetOrderInfo(&_Exchange.CallOpts, order)
}

// GetOrderInfo is a free data retrieval call binding the contract method 0x9d3fa4b9.
//
// Solidity: function getOrderInfo((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order) view returns((uint8,bytes32,uint256) orderInfo)
func (_Exchange *ExchangeCallerSession) GetOrderInfo(order LibOrderOrder) (LibOrderOrderInfo, error) {
	return _Exchange.Contract.GetOrderInfo(&_Exchange.CallOpts, order)
}

// IsValidHashSignature is a free data retrieval call binding the contract method 0x8171c407.
//
// Solidity: function isValidHashSignature(bytes32 hash, address signerAddress, bytes signature) view returns(bool isValid)
func (_Exchange *ExchangeCaller) IsValidHashSignature(opts *bind.CallOpts, hash [32]byte, signerAddress common.Address, signature []byte) (bool, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "isValidHashSignature", hash, signerAddress, signature)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidHashSignature is a free data retrieval call binding the contract method 0x8171c407.
//
// Solidity: function isValidHashSignature(bytes32 hash, address signerAddress, bytes signature) view returns(bool isValid)
func (_Exchange *ExchangeSession) IsValidHashSignature(hash [32]byte, signerAddress common.Address, signature []byte) (bool, error) {
	return _Exchange.Contract.IsValidHashSignature(&_Exchange.CallOpts, hash, signerAddress, signature)
}

// IsValidHashSignature is a free data retrieval call binding the contract method 0x8171c407.
//
// Solidity: function isValidHashSignature(bytes32 hash, address signerAddress, bytes signature) view returns(bool isValid)
func (_Exchange *ExchangeCallerSession) IsValidHashSignature(hash [32]byte, signerAddress common.Address, signature []byte) (bool, error) {
	return _Exchange.Contract.IsValidHashSignature(&_Exchange.CallOpts, hash, signerAddress, signature)
}

// IsValidOrderSignature is a free data retrieval call binding the contract method 0xa12dcc6f.
//
// Solidity: function isValidOrderSignature((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, bytes signature) view returns(bool isValid)
func (_Exchange *ExchangeCaller) IsValidOrderSignature(opts *bind.CallOpts, order LibOrderOrder, signature []byte) (bool, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "isValidOrderSignature", order, signature)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidOrderSignature is a free data retrieval call binding the contract method 0xa12dcc6f.
//
// Solidity: function isValidOrderSignature((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, bytes signature) view returns(bool isValid)
func (_Exchange *ExchangeSession) IsValidOrderSignature(order LibOrderOrder, signature []byte) (bool, error) {
	return _Exchange.Contract.IsValidOrderSignature(&_Exchange.CallOpts, order, signature)
}

// IsValidOrderSignature is a free data retrieval call binding the contract method 0xa12dcc6f.
//
// Solidity: function isValidOrderSignature((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, bytes signature) view returns(bool isValid)
func (_Exchange *ExchangeCallerSession) IsValidOrderSignature(order LibOrderOrder, signature []byte) (bool, error) {
	return _Exchange.Contract.IsValidOrderSignature(&_Exchange.CallOpts, order, signature)
}

// IsValidTransactionSignature is a free data retrieval call binding the contract method 0x8d45cd23.
//
// Solidity: function isValidTransactionSignature((uint256,uint256,uint256,address,bytes) transaction, bytes signature) view returns(bool isValid)
func (_Exchange *ExchangeCaller) IsValidTransactionSignature(opts *bind.CallOpts, transaction LibZeroExTransactionZeroExTransaction, signature []byte) (bool, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "isValidTransactionSignature", transaction, signature)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsValidTransactionSignature is a free data retrieval call binding the contract method 0x8d45cd23.
//
// Solidity: function isValidTransactionSignature((uint256,uint256,uint256,address,bytes) transaction, bytes signature) view returns(bool isValid)
func (_Exchange *ExchangeSession) IsValidTransactionSignature(transaction LibZeroExTransactionZeroExTransaction, signature []byte) (bool, error) {
	return _Exchange.Contract.IsValidTransactionSignature(&_Exchange.CallOpts, transaction, signature)
}

// IsValidTransactionSignature is a free data retrieval call binding the contract method 0x8d45cd23.
//
// Solidity: function isValidTransactionSignature((uint256,uint256,uint256,address,bytes) transaction, bytes signature) view returns(bool isValid)
func (_Exchange *ExchangeCallerSession) IsValidTransactionSignature(transaction LibZeroExTransactionZeroExTransaction, signature []byte) (bool, error) {
	return _Exchange.Contract.IsValidTransactionSignature(&_Exchange.CallOpts, transaction, signature)
}

// OrderEpoch is a free data retrieval call binding the contract method 0xd9bfa73e.
//
// Solidity: function orderEpoch(address , address ) view returns(uint256)
func (_Exchange *ExchangeCaller) OrderEpoch(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "orderEpoch", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OrderEpoch is a free data retrieval call binding the contract method 0xd9bfa73e.
//
// Solidity: function orderEpoch(address , address ) view returns(uint256)
func (_Exchange *ExchangeSession) OrderEpoch(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _Exchange.Contract.OrderEpoch(&_Exchange.CallOpts, arg0, arg1)
}

// OrderEpoch is a free data retrieval call binding the contract method 0xd9bfa73e.
//
// Solidity: function orderEpoch(address , address ) view returns(uint256)
func (_Exchange *ExchangeCallerSession) OrderEpoch(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _Exchange.Contract.OrderEpoch(&_Exchange.CallOpts, arg0, arg1)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Exchange *ExchangeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Exchange *ExchangeSession) Owner() (common.Address, error) {
	return _Exchange.Contract.Owner(&_Exchange.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Exchange *ExchangeCallerSession) Owner() (common.Address, error) {
	return _Exchange.Contract.Owner(&_Exchange.CallOpts)
}

// PreSigned is a free data retrieval call binding the contract method 0x82c174d0.
//
// Solidity: function preSigned(bytes32 , address ) view returns(bool)
func (_Exchange *ExchangeCaller) PreSigned(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "preSigned", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// PreSigned is a free data retrieval call binding the contract method 0x82c174d0.
//
// Solidity: function preSigned(bytes32 , address ) view returns(bool)
func (_Exchange *ExchangeSession) PreSigned(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _Exchange.Contract.PreSigned(&_Exchange.CallOpts, arg0, arg1)
}

// PreSigned is a free data retrieval call binding the contract method 0x82c174d0.
//
// Solidity: function preSigned(bytes32 , address ) view returns(bool)
func (_Exchange *ExchangeCallerSession) PreSigned(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _Exchange.Contract.PreSigned(&_Exchange.CallOpts, arg0, arg1)
}

// ProtocolFeeCollector is a free data retrieval call binding the contract method 0x850a1501.
//
// Solidity: function protocolFeeCollector() view returns(address)
func (_Exchange *ExchangeCaller) ProtocolFeeCollector(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "protocolFeeCollector")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProtocolFeeCollector is a free data retrieval call binding the contract method 0x850a1501.
//
// Solidity: function protocolFeeCollector() view returns(address)
func (_Exchange *ExchangeSession) ProtocolFeeCollector() (common.Address, error) {
	return _Exchange.Contract.ProtocolFeeCollector(&_Exchange.CallOpts)
}

// ProtocolFeeCollector is a free data retrieval call binding the contract method 0x850a1501.
//
// Solidity: function protocolFeeCollector() view returns(address)
func (_Exchange *ExchangeCallerSession) ProtocolFeeCollector() (common.Address, error) {
	return _Exchange.Contract.ProtocolFeeCollector(&_Exchange.CallOpts)
}

// ProtocolFeeMultiplier is a free data retrieval call binding the contract method 0x1ce4c78b.
//
// Solidity: function protocolFeeMultiplier() view returns(uint256)
func (_Exchange *ExchangeCaller) ProtocolFeeMultiplier(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "protocolFeeMultiplier")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProtocolFeeMultiplier is a free data retrieval call binding the contract method 0x1ce4c78b.
//
// Solidity: function protocolFeeMultiplier() view returns(uint256)
func (_Exchange *ExchangeSession) ProtocolFeeMultiplier() (*big.Int, error) {
	return _Exchange.Contract.ProtocolFeeMultiplier(&_Exchange.CallOpts)
}

// ProtocolFeeMultiplier is a free data retrieval call binding the contract method 0x1ce4c78b.
//
// Solidity: function protocolFeeMultiplier() view returns(uint256)
func (_Exchange *ExchangeCallerSession) ProtocolFeeMultiplier() (*big.Int, error) {
	return _Exchange.Contract.ProtocolFeeMultiplier(&_Exchange.CallOpts)
}

// TransactionsExecuted is a free data retrieval call binding the contract method 0x0228e168.
//
// Solidity: function transactionsExecuted(bytes32 ) view returns(bool)
func (_Exchange *ExchangeCaller) TransactionsExecuted(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _Exchange.contract.Call(opts, &out, "transactionsExecuted", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// TransactionsExecuted is a free data retrieval call binding the contract method 0x0228e168.
//
// Solidity: function transactionsExecuted(bytes32 ) view returns(bool)
func (_Exchange *ExchangeSession) TransactionsExecuted(arg0 [32]byte) (bool, error) {
	return _Exchange.Contract.TransactionsExecuted(&_Exchange.CallOpts, arg0)
}

// TransactionsExecuted is a free data retrieval call binding the contract method 0x0228e168.
//
// Solidity: function transactionsExecuted(bytes32 ) view returns(bool)
func (_Exchange *ExchangeCallerSession) TransactionsExecuted(arg0 [32]byte) (bool, error) {
	return _Exchange.Contract.TransactionsExecuted(&_Exchange.CallOpts, arg0)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xdedfc1f1.
//
// Solidity: function batchCancelOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders) payable returns()
func (_Exchange *ExchangeTransactor) BatchCancelOrders(opts *bind.TransactOpts, orders []LibOrderOrder) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchCancelOrders", orders)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xdedfc1f1.
//
// Solidity: function batchCancelOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders) payable returns()
func (_Exchange *ExchangeSession) BatchCancelOrders(orders []LibOrderOrder) (*types.Transaction, error) {
	return _Exchange.Contract.BatchCancelOrders(&_Exchange.TransactOpts, orders)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xdedfc1f1.
//
// Solidity: function batchCancelOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders) payable returns()
func (_Exchange *ExchangeTransactorSession) BatchCancelOrders(orders []LibOrderOrder) (*types.Transaction, error) {
	return _Exchange.Contract.BatchCancelOrders(&_Exchange.TransactOpts, orders)
}

// BatchExecuteTransactions is a paid mutator transaction binding the contract method 0xfc74896d.
//
// Solidity: function batchExecuteTransactions((uint256,uint256,uint256,address,bytes)[] transactions, bytes[] signatures) payable returns(bytes[] returnData)
func (_Exchange *ExchangeTransactor) BatchExecuteTransactions(opts *bind.TransactOpts, transactions []LibZeroExTransactionZeroExTransaction, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchExecuteTransactions", transactions, signatures)
}

// BatchExecuteTransactions is a paid mutator transaction binding the contract method 0xfc74896d.
//
// Solidity: function batchExecuteTransactions((uint256,uint256,uint256,address,bytes)[] transactions, bytes[] signatures) payable returns(bytes[] returnData)
func (_Exchange *ExchangeSession) BatchExecuteTransactions(transactions []LibZeroExTransactionZeroExTransaction, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchExecuteTransactions(&_Exchange.TransactOpts, transactions, signatures)
}

// BatchExecuteTransactions is a paid mutator transaction binding the contract method 0xfc74896d.
//
// Solidity: function batchExecuteTransactions((uint256,uint256,uint256,address,bytes)[] transactions, bytes[] signatures) payable returns(bytes[] returnData)
func (_Exchange *ExchangeTransactorSession) BatchExecuteTransactions(transactions []LibZeroExTransactionZeroExTransaction, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchExecuteTransactions(&_Exchange.TransactOpts, transactions, signatures)
}

// BatchFillOrKillOrders is a paid mutator transaction binding the contract method 0xbeee2e14.
//
// Solidity: function batchFillOrKillOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_Exchange *ExchangeTransactor) BatchFillOrKillOrders(opts *bind.TransactOpts, orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchFillOrKillOrders", orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrKillOrders is a paid mutator transaction binding the contract method 0xbeee2e14.
//
// Solidity: function batchFillOrKillOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_Exchange *ExchangeSession) BatchFillOrKillOrders(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchFillOrKillOrders(&_Exchange.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrKillOrders is a paid mutator transaction binding the contract method 0xbeee2e14.
//
// Solidity: function batchFillOrKillOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_Exchange *ExchangeTransactorSession) BatchFillOrKillOrders(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchFillOrKillOrders(&_Exchange.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrders is a paid mutator transaction binding the contract method 0x9694a402.
//
// Solidity: function batchFillOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_Exchange *ExchangeTransactor) BatchFillOrders(opts *bind.TransactOpts, orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchFillOrders", orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrders is a paid mutator transaction binding the contract method 0x9694a402.
//
// Solidity: function batchFillOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_Exchange *ExchangeSession) BatchFillOrders(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchFillOrders(&_Exchange.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrders is a paid mutator transaction binding the contract method 0x9694a402.
//
// Solidity: function batchFillOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_Exchange *ExchangeTransactorSession) BatchFillOrders(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchFillOrders(&_Exchange.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrdersNoThrow is a paid mutator transaction binding the contract method 0x8ea8dfe4.
//
// Solidity: function batchFillOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_Exchange *ExchangeTransactor) BatchFillOrdersNoThrow(opts *bind.TransactOpts, orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchFillOrdersNoThrow", orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrdersNoThrow is a paid mutator transaction binding the contract method 0x8ea8dfe4.
//
// Solidity: function batchFillOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_Exchange *ExchangeSession) BatchFillOrdersNoThrow(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchFillOrdersNoThrow(&_Exchange.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrdersNoThrow is a paid mutator transaction binding the contract method 0x8ea8dfe4.
//
// Solidity: function batchFillOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_Exchange *ExchangeTransactorSession) BatchFillOrdersNoThrow(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchFillOrdersNoThrow(&_Exchange.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchMatchOrders is a paid mutator transaction binding the contract method 0x6fcf3e9e.
//
// Solidity: function batchMatchOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] leftOrders, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) payable returns(((uint256,uint256,uint256,uint256,uint256)[],(uint256,uint256,uint256,uint256,uint256)[],uint256,uint256) batchMatchedFillResults)
func (_Exchange *ExchangeTransactor) BatchMatchOrders(opts *bind.TransactOpts, leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchMatchOrders", leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// BatchMatchOrders is a paid mutator transaction binding the contract method 0x6fcf3e9e.
//
// Solidity: function batchMatchOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] leftOrders, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) payable returns(((uint256,uint256,uint256,uint256,uint256)[],(uint256,uint256,uint256,uint256,uint256)[],uint256,uint256) batchMatchedFillResults)
func (_Exchange *ExchangeSession) BatchMatchOrders(leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchMatchOrders(&_Exchange.TransactOpts, leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// BatchMatchOrders is a paid mutator transaction binding the contract method 0x6fcf3e9e.
//
// Solidity: function batchMatchOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] leftOrders, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) payable returns(((uint256,uint256,uint256,uint256,uint256)[],(uint256,uint256,uint256,uint256,uint256)[],uint256,uint256) batchMatchedFillResults)
func (_Exchange *ExchangeTransactorSession) BatchMatchOrders(leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchMatchOrders(&_Exchange.TransactOpts, leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// BatchMatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0x6a1a80fd.
//
// Solidity: function batchMatchOrdersWithMaximalFill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] leftOrders, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) payable returns(((uint256,uint256,uint256,uint256,uint256)[],(uint256,uint256,uint256,uint256,uint256)[],uint256,uint256) batchMatchedFillResults)
func (_Exchange *ExchangeTransactor) BatchMatchOrdersWithMaximalFill(opts *bind.TransactOpts, leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "batchMatchOrdersWithMaximalFill", leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// BatchMatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0x6a1a80fd.
//
// Solidity: function batchMatchOrdersWithMaximalFill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] leftOrders, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) payable returns(((uint256,uint256,uint256,uint256,uint256)[],(uint256,uint256,uint256,uint256,uint256)[],uint256,uint256) batchMatchedFillResults)
func (_Exchange *ExchangeSession) BatchMatchOrdersWithMaximalFill(leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchMatchOrdersWithMaximalFill(&_Exchange.TransactOpts, leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// BatchMatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0x6a1a80fd.
//
// Solidity: function batchMatchOrdersWithMaximalFill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] leftOrders, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) payable returns(((uint256,uint256,uint256,uint256,uint256)[],(uint256,uint256,uint256,uint256,uint256)[],uint256,uint256) batchMatchedFillResults)
func (_Exchange *ExchangeTransactorSession) BatchMatchOrdersWithMaximalFill(leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.BatchMatchOrdersWithMaximalFill(&_Exchange.TransactOpts, leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x2da62987.
//
// Solidity: function cancelOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order) payable returns()
func (_Exchange *ExchangeTransactor) CancelOrder(opts *bind.TransactOpts, order LibOrderOrder) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "cancelOrder", order)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x2da62987.
//
// Solidity: function cancelOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order) payable returns()
func (_Exchange *ExchangeSession) CancelOrder(order LibOrderOrder) (*types.Transaction, error) {
	return _Exchange.Contract.CancelOrder(&_Exchange.TransactOpts, order)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x2da62987.
//
// Solidity: function cancelOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order) payable returns()
func (_Exchange *ExchangeTransactorSession) CancelOrder(order LibOrderOrder) (*types.Transaction, error) {
	return _Exchange.Contract.CancelOrder(&_Exchange.TransactOpts, order)
}

// CancelOrdersUpTo is a paid mutator transaction binding the contract method 0x4f9559b1.
//
// Solidity: function cancelOrdersUpTo(uint256 targetOrderEpoch) payable returns()
func (_Exchange *ExchangeTransactor) CancelOrdersUpTo(opts *bind.TransactOpts, targetOrderEpoch *big.Int) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "cancelOrdersUpTo", targetOrderEpoch)
}

// CancelOrdersUpTo is a paid mutator transaction binding the contract method 0x4f9559b1.
//
// Solidity: function cancelOrdersUpTo(uint256 targetOrderEpoch) payable returns()
func (_Exchange *ExchangeSession) CancelOrdersUpTo(targetOrderEpoch *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.CancelOrdersUpTo(&_Exchange.TransactOpts, targetOrderEpoch)
}

// CancelOrdersUpTo is a paid mutator transaction binding the contract method 0x4f9559b1.
//
// Solidity: function cancelOrdersUpTo(uint256 targetOrderEpoch) payable returns()
func (_Exchange *ExchangeTransactorSession) CancelOrdersUpTo(targetOrderEpoch *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.CancelOrdersUpTo(&_Exchange.TransactOpts, targetOrderEpoch)
}

// DetachProtocolFeeCollector is a paid mutator transaction binding the contract method 0x0efca185.
//
// Solidity: function detachProtocolFeeCollector() returns()
func (_Exchange *ExchangeTransactor) DetachProtocolFeeCollector(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "detachProtocolFeeCollector")
}

// DetachProtocolFeeCollector is a paid mutator transaction binding the contract method 0x0efca185.
//
// Solidity: function detachProtocolFeeCollector() returns()
func (_Exchange *ExchangeSession) DetachProtocolFeeCollector() (*types.Transaction, error) {
	return _Exchange.Contract.DetachProtocolFeeCollector(&_Exchange.TransactOpts)
}

// DetachProtocolFeeCollector is a paid mutator transaction binding the contract method 0x0efca185.
//
// Solidity: function detachProtocolFeeCollector() returns()
func (_Exchange *ExchangeTransactorSession) DetachProtocolFeeCollector() (*types.Transaction, error) {
	return _Exchange.Contract.DetachProtocolFeeCollector(&_Exchange.TransactOpts)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0x2280c910.
//
// Solidity: function executeTransaction((uint256,uint256,uint256,address,bytes) transaction, bytes signature) payable returns(bytes)
func (_Exchange *ExchangeTransactor) ExecuteTransaction(opts *bind.TransactOpts, transaction LibZeroExTransactionZeroExTransaction, signature []byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "executeTransaction", transaction, signature)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0x2280c910.
//
// Solidity: function executeTransaction((uint256,uint256,uint256,address,bytes) transaction, bytes signature) payable returns(bytes)
func (_Exchange *ExchangeSession) ExecuteTransaction(transaction LibZeroExTransactionZeroExTransaction, signature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.ExecuteTransaction(&_Exchange.TransactOpts, transaction, signature)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0x2280c910.
//
// Solidity: function executeTransaction((uint256,uint256,uint256,address,bytes) transaction, bytes signature) payable returns(bytes)
func (_Exchange *ExchangeTransactorSession) ExecuteTransaction(transaction LibZeroExTransactionZeroExTransaction, signature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.ExecuteTransaction(&_Exchange.TransactOpts, transaction, signature)
}

// FillOrKillOrder is a paid mutator transaction binding the contract method 0xe14b58c4.
//
// Solidity: function fillOrKillOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 takerAssetFillAmount, bytes signature) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeTransactor) FillOrKillOrder(opts *bind.TransactOpts, order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "fillOrKillOrder", order, takerAssetFillAmount, signature)
}

// FillOrKillOrder is a paid mutator transaction binding the contract method 0xe14b58c4.
//
// Solidity: function fillOrKillOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 takerAssetFillAmount, bytes signature) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeSession) FillOrKillOrder(order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.FillOrKillOrder(&_Exchange.TransactOpts, order, takerAssetFillAmount, signature)
}

// FillOrKillOrder is a paid mutator transaction binding the contract method 0xe14b58c4.
//
// Solidity: function fillOrKillOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 takerAssetFillAmount, bytes signature) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeTransactorSession) FillOrKillOrder(order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.FillOrKillOrder(&_Exchange.TransactOpts, order, takerAssetFillAmount, signature)
}

// FillOrder is a paid mutator transaction binding the contract method 0x9b44d556.
//
// Solidity: function fillOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 takerAssetFillAmount, bytes signature) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeTransactor) FillOrder(opts *bind.TransactOpts, order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "fillOrder", order, takerAssetFillAmount, signature)
}

// FillOrder is a paid mutator transaction binding the contract method 0x9b44d556.
//
// Solidity: function fillOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 takerAssetFillAmount, bytes signature) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeSession) FillOrder(order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.FillOrder(&_Exchange.TransactOpts, order, takerAssetFillAmount, signature)
}

// FillOrder is a paid mutator transaction binding the contract method 0x9b44d556.
//
// Solidity: function fillOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 takerAssetFillAmount, bytes signature) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeTransactorSession) FillOrder(order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.FillOrder(&_Exchange.TransactOpts, order, takerAssetFillAmount, signature)
}

// MarketBuyOrdersFillOrKill is a paid mutator transaction binding the contract method 0x8bc8efb3.
//
// Solidity: function marketBuyOrdersFillOrKill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 makerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeTransactor) MarketBuyOrdersFillOrKill(opts *bind.TransactOpts, orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "marketBuyOrdersFillOrKill", orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersFillOrKill is a paid mutator transaction binding the contract method 0x8bc8efb3.
//
// Solidity: function marketBuyOrdersFillOrKill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 makerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeSession) MarketBuyOrdersFillOrKill(orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.MarketBuyOrdersFillOrKill(&_Exchange.TransactOpts, orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersFillOrKill is a paid mutator transaction binding the contract method 0x8bc8efb3.
//
// Solidity: function marketBuyOrdersFillOrKill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 makerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeTransactorSession) MarketBuyOrdersFillOrKill(orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.MarketBuyOrdersFillOrKill(&_Exchange.TransactOpts, orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersNoThrow is a paid mutator transaction binding the contract method 0x78d29ac1.
//
// Solidity: function marketBuyOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 makerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeTransactor) MarketBuyOrdersNoThrow(opts *bind.TransactOpts, orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "marketBuyOrdersNoThrow", orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersNoThrow is a paid mutator transaction binding the contract method 0x78d29ac1.
//
// Solidity: function marketBuyOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 makerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeSession) MarketBuyOrdersNoThrow(orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.MarketBuyOrdersNoThrow(&_Exchange.TransactOpts, orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersNoThrow is a paid mutator transaction binding the contract method 0x78d29ac1.
//
// Solidity: function marketBuyOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 makerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeTransactorSession) MarketBuyOrdersNoThrow(orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.MarketBuyOrdersNoThrow(&_Exchange.TransactOpts, orders, makerAssetFillAmount, signatures)
}

// MarketSellOrdersFillOrKill is a paid mutator transaction binding the contract method 0xa6c3bf33.
//
// Solidity: function marketSellOrdersFillOrKill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 takerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeTransactor) MarketSellOrdersFillOrKill(opts *bind.TransactOpts, orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "marketSellOrdersFillOrKill", orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersFillOrKill is a paid mutator transaction binding the contract method 0xa6c3bf33.
//
// Solidity: function marketSellOrdersFillOrKill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 takerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeSession) MarketSellOrdersFillOrKill(orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.MarketSellOrdersFillOrKill(&_Exchange.TransactOpts, orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersFillOrKill is a paid mutator transaction binding the contract method 0xa6c3bf33.
//
// Solidity: function marketSellOrdersFillOrKill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 takerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeTransactorSession) MarketSellOrdersFillOrKill(orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.MarketSellOrdersFillOrKill(&_Exchange.TransactOpts, orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersNoThrow is a paid mutator transaction binding the contract method 0x369da099.
//
// Solidity: function marketSellOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 takerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeTransactor) MarketSellOrdersNoThrow(opts *bind.TransactOpts, orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "marketSellOrdersNoThrow", orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersNoThrow is a paid mutator transaction binding the contract method 0x369da099.
//
// Solidity: function marketSellOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 takerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeSession) MarketSellOrdersNoThrow(orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.MarketSellOrdersNoThrow(&_Exchange.TransactOpts, orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersNoThrow is a paid mutator transaction binding the contract method 0x369da099.
//
// Solidity: function marketSellOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 takerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_Exchange *ExchangeTransactorSession) MarketSellOrdersNoThrow(orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _Exchange.Contract.MarketSellOrdersNoThrow(&_Exchange.TransactOpts, orders, takerAssetFillAmount, signatures)
}

// MatchOrders is a paid mutator transaction binding the contract method 0x88ec79fb.
//
// Solidity: function matchOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) leftOrder, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) rightOrder, bytes leftSignature, bytes rightSignature) payable returns(((uint256,uint256,uint256,uint256,uint256),(uint256,uint256,uint256,uint256,uint256),uint256,uint256) matchedFillResults)
func (_Exchange *ExchangeTransactor) MatchOrders(opts *bind.TransactOpts, leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "matchOrders", leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrders is a paid mutator transaction binding the contract method 0x88ec79fb.
//
// Solidity: function matchOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) leftOrder, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) rightOrder, bytes leftSignature, bytes rightSignature) payable returns(((uint256,uint256,uint256,uint256,uint256),(uint256,uint256,uint256,uint256,uint256),uint256,uint256) matchedFillResults)
func (_Exchange *ExchangeSession) MatchOrders(leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.MatchOrders(&_Exchange.TransactOpts, leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrders is a paid mutator transaction binding the contract method 0x88ec79fb.
//
// Solidity: function matchOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) leftOrder, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) rightOrder, bytes leftSignature, bytes rightSignature) payable returns(((uint256,uint256,uint256,uint256,uint256),(uint256,uint256,uint256,uint256,uint256),uint256,uint256) matchedFillResults)
func (_Exchange *ExchangeTransactorSession) MatchOrders(leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.MatchOrders(&_Exchange.TransactOpts, leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0xb718e292.
//
// Solidity: function matchOrdersWithMaximalFill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) leftOrder, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) rightOrder, bytes leftSignature, bytes rightSignature) payable returns(((uint256,uint256,uint256,uint256,uint256),(uint256,uint256,uint256,uint256,uint256),uint256,uint256) matchedFillResults)
func (_Exchange *ExchangeTransactor) MatchOrdersWithMaximalFill(opts *bind.TransactOpts, leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "matchOrdersWithMaximalFill", leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0xb718e292.
//
// Solidity: function matchOrdersWithMaximalFill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) leftOrder, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) rightOrder, bytes leftSignature, bytes rightSignature) payable returns(((uint256,uint256,uint256,uint256,uint256),(uint256,uint256,uint256,uint256,uint256),uint256,uint256) matchedFillResults)
func (_Exchange *ExchangeSession) MatchOrdersWithMaximalFill(leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.MatchOrdersWithMaximalFill(&_Exchange.TransactOpts, leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0xb718e292.
//
// Solidity: function matchOrdersWithMaximalFill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) leftOrder, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) rightOrder, bytes leftSignature, bytes rightSignature) payable returns(((uint256,uint256,uint256,uint256,uint256),(uint256,uint256,uint256,uint256,uint256),uint256,uint256) matchedFillResults)
func (_Exchange *ExchangeTransactorSession) MatchOrdersWithMaximalFill(leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _Exchange.Contract.MatchOrdersWithMaximalFill(&_Exchange.TransactOpts, leftOrder, rightOrder, leftSignature, rightSignature)
}

// PreSign is a paid mutator transaction binding the contract method 0x46c02d7a.
//
// Solidity: function preSign(bytes32 hash) payable returns()
func (_Exchange *ExchangeTransactor) PreSign(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "preSign", hash)
}

// PreSign is a paid mutator transaction binding the contract method 0x46c02d7a.
//
// Solidity: function preSign(bytes32 hash) payable returns()
func (_Exchange *ExchangeSession) PreSign(hash [32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.PreSign(&_Exchange.TransactOpts, hash)
}

// PreSign is a paid mutator transaction binding the contract method 0x46c02d7a.
//
// Solidity: function preSign(bytes32 hash) payable returns()
func (_Exchange *ExchangeTransactorSession) PreSign(hash [32]byte) (*types.Transaction, error) {
	return _Exchange.Contract.PreSign(&_Exchange.TransactOpts, hash)
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

// SetProtocolFeeCollectorAddress is a paid mutator transaction binding the contract method 0xc0fa16cc.
//
// Solidity: function setProtocolFeeCollectorAddress(address updatedProtocolFeeCollector) returns()
func (_Exchange *ExchangeTransactor) SetProtocolFeeCollectorAddress(opts *bind.TransactOpts, updatedProtocolFeeCollector common.Address) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "setProtocolFeeCollectorAddress", updatedProtocolFeeCollector)
}

// SetProtocolFeeCollectorAddress is a paid mutator transaction binding the contract method 0xc0fa16cc.
//
// Solidity: function setProtocolFeeCollectorAddress(address updatedProtocolFeeCollector) returns()
func (_Exchange *ExchangeSession) SetProtocolFeeCollectorAddress(updatedProtocolFeeCollector common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.SetProtocolFeeCollectorAddress(&_Exchange.TransactOpts, updatedProtocolFeeCollector)
}

// SetProtocolFeeCollectorAddress is a paid mutator transaction binding the contract method 0xc0fa16cc.
//
// Solidity: function setProtocolFeeCollectorAddress(address updatedProtocolFeeCollector) returns()
func (_Exchange *ExchangeTransactorSession) SetProtocolFeeCollectorAddress(updatedProtocolFeeCollector common.Address) (*types.Transaction, error) {
	return _Exchange.Contract.SetProtocolFeeCollectorAddress(&_Exchange.TransactOpts, updatedProtocolFeeCollector)
}

// SetProtocolFeeMultiplier is a paid mutator transaction binding the contract method 0x9331c742.
//
// Solidity: function setProtocolFeeMultiplier(uint256 updatedProtocolFeeMultiplier) returns()
func (_Exchange *ExchangeTransactor) SetProtocolFeeMultiplier(opts *bind.TransactOpts, updatedProtocolFeeMultiplier *big.Int) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "setProtocolFeeMultiplier", updatedProtocolFeeMultiplier)
}

// SetProtocolFeeMultiplier is a paid mutator transaction binding the contract method 0x9331c742.
//
// Solidity: function setProtocolFeeMultiplier(uint256 updatedProtocolFeeMultiplier) returns()
func (_Exchange *ExchangeSession) SetProtocolFeeMultiplier(updatedProtocolFeeMultiplier *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.SetProtocolFeeMultiplier(&_Exchange.TransactOpts, updatedProtocolFeeMultiplier)
}

// SetProtocolFeeMultiplier is a paid mutator transaction binding the contract method 0x9331c742.
//
// Solidity: function setProtocolFeeMultiplier(uint256 updatedProtocolFeeMultiplier) returns()
func (_Exchange *ExchangeTransactorSession) SetProtocolFeeMultiplier(updatedProtocolFeeMultiplier *big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.SetProtocolFeeMultiplier(&_Exchange.TransactOpts, updatedProtocolFeeMultiplier)
}

// SetSignatureValidatorApproval is a paid mutator transaction binding the contract method 0x77fcce68.
//
// Solidity: function setSignatureValidatorApproval(address validatorAddress, bool approval) payable returns()
func (_Exchange *ExchangeTransactor) SetSignatureValidatorApproval(opts *bind.TransactOpts, validatorAddress common.Address, approval bool) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "setSignatureValidatorApproval", validatorAddress, approval)
}

// SetSignatureValidatorApproval is a paid mutator transaction binding the contract method 0x77fcce68.
//
// Solidity: function setSignatureValidatorApproval(address validatorAddress, bool approval) payable returns()
func (_Exchange *ExchangeSession) SetSignatureValidatorApproval(validatorAddress common.Address, approval bool) (*types.Transaction, error) {
	return _Exchange.Contract.SetSignatureValidatorApproval(&_Exchange.TransactOpts, validatorAddress, approval)
}

// SetSignatureValidatorApproval is a paid mutator transaction binding the contract method 0x77fcce68.
//
// Solidity: function setSignatureValidatorApproval(address validatorAddress, bool approval) payable returns()
func (_Exchange *ExchangeTransactorSession) SetSignatureValidatorApproval(validatorAddress common.Address, approval bool) (*types.Transaction, error) {
	return _Exchange.Contract.SetSignatureValidatorApproval(&_Exchange.TransactOpts, validatorAddress, approval)
}

// SimulateDispatchTransferFromCalls is a paid mutator transaction binding the contract method 0xb04fbddd.
//
// Solidity: function simulateDispatchTransferFromCalls(bytes[] assetData, address[] fromAddresses, address[] toAddresses, uint256[] amounts) returns()
func (_Exchange *ExchangeTransactor) SimulateDispatchTransferFromCalls(opts *bind.TransactOpts, assetData [][]byte, fromAddresses []common.Address, toAddresses []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Exchange.contract.Transact(opts, "simulateDispatchTransferFromCalls", assetData, fromAddresses, toAddresses, amounts)
}

// SimulateDispatchTransferFromCalls is a paid mutator transaction binding the contract method 0xb04fbddd.
//
// Solidity: function simulateDispatchTransferFromCalls(bytes[] assetData, address[] fromAddresses, address[] toAddresses, uint256[] amounts) returns()
func (_Exchange *ExchangeSession) SimulateDispatchTransferFromCalls(assetData [][]byte, fromAddresses []common.Address, toAddresses []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.SimulateDispatchTransferFromCalls(&_Exchange.TransactOpts, assetData, fromAddresses, toAddresses, amounts)
}

// SimulateDispatchTransferFromCalls is a paid mutator transaction binding the contract method 0xb04fbddd.
//
// Solidity: function simulateDispatchTransferFromCalls(bytes[] assetData, address[] fromAddresses, address[] toAddresses, uint256[] amounts) returns()
func (_Exchange *ExchangeTransactorSession) SimulateDispatchTransferFromCalls(assetData [][]byte, fromAddresses []common.Address, toAddresses []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Exchange.Contract.SimulateDispatchTransferFromCalls(&_Exchange.TransactOpts, assetData, fromAddresses, toAddresses, amounts)
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
	event.Raw = log
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
	MakerAssetData      []byte
	TakerAssetData      []byte
	SenderAddress       common.Address
	OrderHash           [32]byte
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterCancel is a free log retrieval operation binding the contract event 0x02c310a9a43963ff31a754a4099cc435ed498049687539d72d7818d9b093415c.
//
// Solidity: event Cancel(address indexed makerAddress, address indexed feeRecipientAddress, bytes makerAssetData, bytes takerAssetData, address senderAddress, bytes32 indexed orderHash)
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

// WatchCancel is a free log subscription operation binding the contract event 0x02c310a9a43963ff31a754a4099cc435ed498049687539d72d7818d9b093415c.
//
// Solidity: event Cancel(address indexed makerAddress, address indexed feeRecipientAddress, bytes makerAssetData, bytes takerAssetData, address senderAddress, bytes32 indexed orderHash)
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

// ParseCancel is a log parse operation binding the contract event 0x02c310a9a43963ff31a754a4099cc435ed498049687539d72d7818d9b093415c.
//
// Solidity: event Cancel(address indexed makerAddress, address indexed feeRecipientAddress, bytes makerAssetData, bytes takerAssetData, address senderAddress, bytes32 indexed orderHash)
func (_Exchange *ExchangeFilterer) ParseCancel(log types.Log) (*ExchangeCancel, error) {
	event := new(ExchangeCancel)
	if err := _Exchange.contract.UnpackLog(event, "Cancel", log); err != nil {
		return nil, err
	}
	event.Raw = log
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
	MakerAddress       common.Address
	OrderSenderAddress common.Address
	OrderEpoch         *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterCancelUpTo is a free log retrieval operation binding the contract event 0x82af639571738f4ebd4268fb0363d8957ebe1bbb9e78dba5ebd69eed39b154f0.
//
// Solidity: event CancelUpTo(address indexed makerAddress, address indexed orderSenderAddress, uint256 orderEpoch)
func (_Exchange *ExchangeFilterer) FilterCancelUpTo(opts *bind.FilterOpts, makerAddress []common.Address, orderSenderAddress []common.Address) (*ExchangeCancelUpToIterator, error) {

	var makerAddressRule []interface{}
	for _, makerAddressItem := range makerAddress {
		makerAddressRule = append(makerAddressRule, makerAddressItem)
	}
	var orderSenderAddressRule []interface{}
	for _, orderSenderAddressItem := range orderSenderAddress {
		orderSenderAddressRule = append(orderSenderAddressRule, orderSenderAddressItem)
	}

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "CancelUpTo", makerAddressRule, orderSenderAddressRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeCancelUpToIterator{contract: _Exchange.contract, event: "CancelUpTo", logs: logs, sub: sub}, nil
}

// WatchCancelUpTo is a free log subscription operation binding the contract event 0x82af639571738f4ebd4268fb0363d8957ebe1bbb9e78dba5ebd69eed39b154f0.
//
// Solidity: event CancelUpTo(address indexed makerAddress, address indexed orderSenderAddress, uint256 orderEpoch)
func (_Exchange *ExchangeFilterer) WatchCancelUpTo(opts *bind.WatchOpts, sink chan<- *ExchangeCancelUpTo, makerAddress []common.Address, orderSenderAddress []common.Address) (event.Subscription, error) {

	var makerAddressRule []interface{}
	for _, makerAddressItem := range makerAddress {
		makerAddressRule = append(makerAddressRule, makerAddressItem)
	}
	var orderSenderAddressRule []interface{}
	for _, orderSenderAddressItem := range orderSenderAddress {
		orderSenderAddressRule = append(orderSenderAddressRule, orderSenderAddressItem)
	}

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "CancelUpTo", makerAddressRule, orderSenderAddressRule)
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
// Solidity: event CancelUpTo(address indexed makerAddress, address indexed orderSenderAddress, uint256 orderEpoch)
func (_Exchange *ExchangeFilterer) ParseCancelUpTo(log types.Log) (*ExchangeCancelUpTo, error) {
	event := new(ExchangeCancelUpTo)
	if err := _Exchange.contract.UnpackLog(event, "CancelUpTo", log); err != nil {
		return nil, err
	}
	event.Raw = log
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
	MakerAssetData         []byte
	TakerAssetData         []byte
	MakerFeeAssetData      []byte
	TakerFeeAssetData      []byte
	OrderHash              [32]byte
	TakerAddress           common.Address
	SenderAddress          common.Address
	MakerAssetFilledAmount *big.Int
	TakerAssetFilledAmount *big.Int
	MakerFeePaid           *big.Int
	TakerFeePaid           *big.Int
	ProtocolFeePaid        *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterFill is a free log retrieval operation binding the contract event 0x6869791f0a34781b29882982cc39e882768cf2c96995c2a110c577c53bc932d5.
//
// Solidity: event Fill(address indexed makerAddress, address indexed feeRecipientAddress, bytes makerAssetData, bytes takerAssetData, bytes makerFeeAssetData, bytes takerFeeAssetData, bytes32 indexed orderHash, address takerAddress, address senderAddress, uint256 makerAssetFilledAmount, uint256 takerAssetFilledAmount, uint256 makerFeePaid, uint256 takerFeePaid, uint256 protocolFeePaid)
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

// WatchFill is a free log subscription operation binding the contract event 0x6869791f0a34781b29882982cc39e882768cf2c96995c2a110c577c53bc932d5.
//
// Solidity: event Fill(address indexed makerAddress, address indexed feeRecipientAddress, bytes makerAssetData, bytes takerAssetData, bytes makerFeeAssetData, bytes takerFeeAssetData, bytes32 indexed orderHash, address takerAddress, address senderAddress, uint256 makerAssetFilledAmount, uint256 takerAssetFilledAmount, uint256 makerFeePaid, uint256 takerFeePaid, uint256 protocolFeePaid)
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

// ParseFill is a log parse operation binding the contract event 0x6869791f0a34781b29882982cc39e882768cf2c96995c2a110c577c53bc932d5.
//
// Solidity: event Fill(address indexed makerAddress, address indexed feeRecipientAddress, bytes makerAssetData, bytes takerAssetData, bytes makerFeeAssetData, bytes takerFeeAssetData, bytes32 indexed orderHash, address takerAddress, address senderAddress, uint256 makerAssetFilledAmount, uint256 takerAssetFilledAmount, uint256 makerFeePaid, uint256 takerFeePaid, uint256 protocolFeePaid)
func (_Exchange *ExchangeFilterer) ParseFill(log types.Log) (*ExchangeFill, error) {
	event := new(ExchangeFill)
	if err := _Exchange.contract.UnpackLog(event, "Fill", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExchangeOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Exchange contract.
type ExchangeOwnershipTransferredIterator struct {
	Event *ExchangeOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ExchangeOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeOwnershipTransferred)
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
		it.Event = new(ExchangeOwnershipTransferred)
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
func (it *ExchangeOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeOwnershipTransferred represents a OwnershipTransferred event raised by the Exchange contract.
type ExchangeOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Exchange *ExchangeFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ExchangeOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeOwnershipTransferredIterator{contract: _Exchange.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Exchange *ExchangeFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ExchangeOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeOwnershipTransferred)
				if err := _Exchange.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Exchange *ExchangeFilterer) ParseOwnershipTransferred(log types.Log) (*ExchangeOwnershipTransferred, error) {
	event := new(ExchangeOwnershipTransferred)
	if err := _Exchange.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExchangeProtocolFeeCollectorAddressIterator is returned from FilterProtocolFeeCollectorAddress and is used to iterate over the raw logs and unpacked data for ProtocolFeeCollectorAddress events raised by the Exchange contract.
type ExchangeProtocolFeeCollectorAddressIterator struct {
	Event *ExchangeProtocolFeeCollectorAddress // Event containing the contract specifics and raw log

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
func (it *ExchangeProtocolFeeCollectorAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeProtocolFeeCollectorAddress)
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
		it.Event = new(ExchangeProtocolFeeCollectorAddress)
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
func (it *ExchangeProtocolFeeCollectorAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeProtocolFeeCollectorAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeProtocolFeeCollectorAddress represents a ProtocolFeeCollectorAddress event raised by the Exchange contract.
type ExchangeProtocolFeeCollectorAddress struct {
	OldProtocolFeeCollector     common.Address
	UpdatedProtocolFeeCollector common.Address
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterProtocolFeeCollectorAddress is a free log retrieval operation binding the contract event 0xe1a5430ebec577336427f40f15822f1f36c5e3509ff209d6db9e6c9e6941cb0b.
//
// Solidity: event ProtocolFeeCollectorAddress(address oldProtocolFeeCollector, address updatedProtocolFeeCollector)
func (_Exchange *ExchangeFilterer) FilterProtocolFeeCollectorAddress(opts *bind.FilterOpts) (*ExchangeProtocolFeeCollectorAddressIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "ProtocolFeeCollectorAddress")
	if err != nil {
		return nil, err
	}
	return &ExchangeProtocolFeeCollectorAddressIterator{contract: _Exchange.contract, event: "ProtocolFeeCollectorAddress", logs: logs, sub: sub}, nil
}

// WatchProtocolFeeCollectorAddress is a free log subscription operation binding the contract event 0xe1a5430ebec577336427f40f15822f1f36c5e3509ff209d6db9e6c9e6941cb0b.
//
// Solidity: event ProtocolFeeCollectorAddress(address oldProtocolFeeCollector, address updatedProtocolFeeCollector)
func (_Exchange *ExchangeFilterer) WatchProtocolFeeCollectorAddress(opts *bind.WatchOpts, sink chan<- *ExchangeProtocolFeeCollectorAddress) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "ProtocolFeeCollectorAddress")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeProtocolFeeCollectorAddress)
				if err := _Exchange.contract.UnpackLog(event, "ProtocolFeeCollectorAddress", log); err != nil {
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

// ParseProtocolFeeCollectorAddress is a log parse operation binding the contract event 0xe1a5430ebec577336427f40f15822f1f36c5e3509ff209d6db9e6c9e6941cb0b.
//
// Solidity: event ProtocolFeeCollectorAddress(address oldProtocolFeeCollector, address updatedProtocolFeeCollector)
func (_Exchange *ExchangeFilterer) ParseProtocolFeeCollectorAddress(log types.Log) (*ExchangeProtocolFeeCollectorAddress, error) {
	event := new(ExchangeProtocolFeeCollectorAddress)
	if err := _Exchange.contract.UnpackLog(event, "ProtocolFeeCollectorAddress", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExchangeProtocolFeeMultiplierIterator is returned from FilterProtocolFeeMultiplier and is used to iterate over the raw logs and unpacked data for ProtocolFeeMultiplier events raised by the Exchange contract.
type ExchangeProtocolFeeMultiplierIterator struct {
	Event *ExchangeProtocolFeeMultiplier // Event containing the contract specifics and raw log

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
func (it *ExchangeProtocolFeeMultiplierIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeProtocolFeeMultiplier)
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
		it.Event = new(ExchangeProtocolFeeMultiplier)
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
func (it *ExchangeProtocolFeeMultiplierIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeProtocolFeeMultiplierIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeProtocolFeeMultiplier represents a ProtocolFeeMultiplier event raised by the Exchange contract.
type ExchangeProtocolFeeMultiplier struct {
	OldProtocolFeeMultiplier     *big.Int
	UpdatedProtocolFeeMultiplier *big.Int
	Raw                          types.Log // Blockchain specific contextual infos
}

// FilterProtocolFeeMultiplier is a free log retrieval operation binding the contract event 0x3a3e76d7a75e198aef1f53137e4f2a8a2ec74e2e9526db8404d08ccc9f1e621d.
//
// Solidity: event ProtocolFeeMultiplier(uint256 oldProtocolFeeMultiplier, uint256 updatedProtocolFeeMultiplier)
func (_Exchange *ExchangeFilterer) FilterProtocolFeeMultiplier(opts *bind.FilterOpts) (*ExchangeProtocolFeeMultiplierIterator, error) {

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "ProtocolFeeMultiplier")
	if err != nil {
		return nil, err
	}
	return &ExchangeProtocolFeeMultiplierIterator{contract: _Exchange.contract, event: "ProtocolFeeMultiplier", logs: logs, sub: sub}, nil
}

// WatchProtocolFeeMultiplier is a free log subscription operation binding the contract event 0x3a3e76d7a75e198aef1f53137e4f2a8a2ec74e2e9526db8404d08ccc9f1e621d.
//
// Solidity: event ProtocolFeeMultiplier(uint256 oldProtocolFeeMultiplier, uint256 updatedProtocolFeeMultiplier)
func (_Exchange *ExchangeFilterer) WatchProtocolFeeMultiplier(opts *bind.WatchOpts, sink chan<- *ExchangeProtocolFeeMultiplier) (event.Subscription, error) {

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "ProtocolFeeMultiplier")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeProtocolFeeMultiplier)
				if err := _Exchange.contract.UnpackLog(event, "ProtocolFeeMultiplier", log); err != nil {
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

// ParseProtocolFeeMultiplier is a log parse operation binding the contract event 0x3a3e76d7a75e198aef1f53137e4f2a8a2ec74e2e9526db8404d08ccc9f1e621d.
//
// Solidity: event ProtocolFeeMultiplier(uint256 oldProtocolFeeMultiplier, uint256 updatedProtocolFeeMultiplier)
func (_Exchange *ExchangeFilterer) ParseProtocolFeeMultiplier(log types.Log) (*ExchangeProtocolFeeMultiplier, error) {
	event := new(ExchangeProtocolFeeMultiplier)
	if err := _Exchange.contract.UnpackLog(event, "ProtocolFeeMultiplier", log); err != nil {
		return nil, err
	}
	event.Raw = log
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
	IsApproved       bool
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSignatureValidatorApproval is a free log retrieval operation binding the contract event 0xa8656e308026eeabce8f0bc18048433252318ab80ac79da0b3d3d8697dfba891.
//
// Solidity: event SignatureValidatorApproval(address indexed signerAddress, address indexed validatorAddress, bool isApproved)
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
// Solidity: event SignatureValidatorApproval(address indexed signerAddress, address indexed validatorAddress, bool isApproved)
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
// Solidity: event SignatureValidatorApproval(address indexed signerAddress, address indexed validatorAddress, bool isApproved)
func (_Exchange *ExchangeFilterer) ParseSignatureValidatorApproval(log types.Log) (*ExchangeSignatureValidatorApproval, error) {
	event := new(ExchangeSignatureValidatorApproval)
	if err := _Exchange.contract.UnpackLog(event, "SignatureValidatorApproval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExchangeTransactionExecutionIterator is returned from FilterTransactionExecution and is used to iterate over the raw logs and unpacked data for TransactionExecution events raised by the Exchange contract.
type ExchangeTransactionExecutionIterator struct {
	Event *ExchangeTransactionExecution // Event containing the contract specifics and raw log

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
func (it *ExchangeTransactionExecutionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeTransactionExecution)
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
		it.Event = new(ExchangeTransactionExecution)
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
func (it *ExchangeTransactionExecutionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeTransactionExecutionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeTransactionExecution represents a TransactionExecution event raised by the Exchange contract.
type ExchangeTransactionExecution struct {
	TransactionHash [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterTransactionExecution is a free log retrieval operation binding the contract event 0xa4a7329f1dd821363067e07d359e347b4af9b1efe4b6cccf13240228af3c800d.
//
// Solidity: event TransactionExecution(bytes32 indexed transactionHash)
func (_Exchange *ExchangeFilterer) FilterTransactionExecution(opts *bind.FilterOpts, transactionHash [][32]byte) (*ExchangeTransactionExecutionIterator, error) {

	var transactionHashRule []interface{}
	for _, transactionHashItem := range transactionHash {
		transactionHashRule = append(transactionHashRule, transactionHashItem)
	}

	logs, sub, err := _Exchange.contract.FilterLogs(opts, "TransactionExecution", transactionHashRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeTransactionExecutionIterator{contract: _Exchange.contract, event: "TransactionExecution", logs: logs, sub: sub}, nil
}

// WatchTransactionExecution is a free log subscription operation binding the contract event 0xa4a7329f1dd821363067e07d359e347b4af9b1efe4b6cccf13240228af3c800d.
//
// Solidity: event TransactionExecution(bytes32 indexed transactionHash)
func (_Exchange *ExchangeFilterer) WatchTransactionExecution(opts *bind.WatchOpts, sink chan<- *ExchangeTransactionExecution, transactionHash [][32]byte) (event.Subscription, error) {

	var transactionHashRule []interface{}
	for _, transactionHashItem := range transactionHash {
		transactionHashRule = append(transactionHashRule, transactionHashItem)
	}

	logs, sub, err := _Exchange.contract.WatchLogs(opts, "TransactionExecution", transactionHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeTransactionExecution)
				if err := _Exchange.contract.UnpackLog(event, "TransactionExecution", log); err != nil {
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

// ParseTransactionExecution is a log parse operation binding the contract event 0xa4a7329f1dd821363067e07d359e347b4af9b1efe4b6cccf13240228af3c800d.
//
// Solidity: event TransactionExecution(bytes32 indexed transactionHash)
func (_Exchange *ExchangeFilterer) ParseTransactionExecution(log types.Log) (*ExchangeTransactionExecution, error) {
	event := new(ExchangeTransactionExecution)
	if err := _Exchange.contract.UnpackLog(event, "TransactionExecution", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
