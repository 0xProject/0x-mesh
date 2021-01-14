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

// WrappersV4ABI is the input ABI used to generate the binding from.
const WrappersV4ABI = "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"id\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetProxy\",\"type\":\"address\"}],\"name\":\"AssetProxyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"name\":\"Cancel\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"orderSenderAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"orderEpoch\",\"type\":\"uint256\"}],\"name\":\"CancelUpTo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"name\":\"Fill\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldProtocolFeeCollector\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"updatedProtocolFeeCollector\",\"type\":\"address\"}],\"name\":\"ProtocolFeeCollectorAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldProtocolFeeMultiplier\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedProtocolFeeMultiplier\",\"type\":\"uint256\"}],\"name\":\"ProtocolFeeMultiplier\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isApproved\",\"type\":\"bool\"}],\"name\":\"SignatureValidatorApproval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"transactionHash\",\"type\":\"bytes32\"}],\"name\":\"TransactionExecution\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"EIP1271_MAGIC_VALUE\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"EIP712_EXCHANGE_DOMAIN_HASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedValidators\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"}],\"name\":\"batchCancelOrders\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLibZeroExTransaction.ZeroExTransaction[]\",\"name\":\"transactions\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchExecuteTransactions\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"returnData\",\"type\":\"bytes[]\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchFillOrKillOrders\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"fillResults\",\"type\":\"tuple[]\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchFillOrders\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"fillResults\",\"type\":\"tuple[]\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"batchFillOrdersNoThrow\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"fillResults\",\"type\":\"tuple[]\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"leftOrders\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"rightOrders\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"leftSignatures\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes[]\",\"name\":\"rightSignatures\",\"type\":\"bytes[]\"}],\"name\":\"batchMatchOrders\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"left\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"right\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"profitInLeftMakerAsset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"profitInRightMakerAsset\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.BatchMatchedFillResults\",\"name\":\"batchMatchedFillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"leftOrders\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"rightOrders\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"leftSignatures\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes[]\",\"name\":\"rightSignatures\",\"type\":\"bytes[]\"}],\"name\":\"batchMatchOrdersWithMaximalFill\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"left\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults[]\",\"name\":\"right\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"profitInLeftMakerAsset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"profitInRightMakerAsset\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.BatchMatchedFillResults\",\"name\":\"batchMatchedFillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"cancelOrder\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"targetOrderEpoch\",\"type\":\"uint256\"}],\"name\":\"cancelOrdersUpTo\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"cancelled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentContextAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"detachProtocolFeeCollector\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLibZeroExTransaction.ZeroExTransaction\",\"name\":\"transaction\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"executeTransaction\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"fillOrKillOrder\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"fillOrder\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"filled\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"assetProxyId\",\"type\":\"bytes4\"}],\"name\":\"getAssetProxy\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"assetProxy\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"getOrderInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"enumLibOrder.OrderStatus\",\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"internalType\":\"structLibOrder.OrderInfo\",\"name\":\"orderInfo\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"isValidHashSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isValid\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"isValidOrderSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isValid\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gasPrice\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signerAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLibZeroExTransaction.ZeroExTransaction\",\"name\":\"transaction\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"isValidTransactionSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isValid\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketBuyOrdersFillOrKill\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketBuyOrdersNoThrow\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketSellOrdersFillOrKill\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFillAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes[]\",\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"marketSellOrdersNoThrow\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"fillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"leftOrder\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"rightOrder\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"leftSignature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rightSignature\",\"type\":\"bytes\"}],\"name\":\"matchOrders\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"left\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"right\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"profitInLeftMakerAsset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"profitInRightMakerAsset\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.MatchedFillResults\",\"name\":\"matchedFillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"leftOrder\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"makerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"takerAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"senderAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"makerFeeAssetData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"takerFeeAssetData\",\"type\":\"bytes\"}],\"internalType\":\"structLibOrder.Order\",\"name\":\"rightOrder\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"leftSignature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"rightSignature\",\"type\":\"bytes\"}],\"name\":\"matchOrdersWithMaximalFill\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"left\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"makerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerAssetFilledAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"makerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"takerFeePaid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.FillResults\",\"name\":\"right\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"profitInLeftMakerAsset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"profitInRightMakerAsset\",\"type\":\"uint256\"}],\"internalType\":\"structLibFillResults.MatchedFillResults\",\"name\":\"matchedFillResults\",\"type\":\"tuple\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"orderEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"preSign\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"preSigned\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"protocolFeeCollector\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"protocolFeeMultiplier\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetProxy\",\"type\":\"address\"}],\"name\":\"registerAssetProxy\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"updatedProtocolFeeCollector\",\"type\":\"address\"}],\"name\":\"setProtocolFeeCollectorAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"updatedProtocolFeeMultiplier\",\"type\":\"uint256\"}],\"name\":\"setProtocolFeeMultiplier\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"validatorAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approval\",\"type\":\"bool\"}],\"name\":\"setSignatureValidatorApproval\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"assetData\",\"type\":\"bytes[]\"},{\"internalType\":\"address[]\",\"name\":\"fromAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"toAddresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"simulateDispatchTransferFromCalls\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"transactionsExecuted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

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

// EIP1271MAGICVALUE is a free data retrieval call binding the contract method 0xdd885e2d.
//
// Solidity: function EIP1271_MAGIC_VALUE() view returns(bytes4)
func (_WrappersV4 *WrappersV4Caller) EIP1271MAGICVALUE(opts *bind.CallOpts) ([4]byte, error) {
	var (
		ret0 = new([4]byte)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "EIP1271_MAGIC_VALUE")
	return *ret0, err
}

// EIP1271MAGICVALUE is a free data retrieval call binding the contract method 0xdd885e2d.
//
// Solidity: function EIP1271_MAGIC_VALUE() view returns(bytes4)
func (_WrappersV4 *WrappersV4Session) EIP1271MAGICVALUE() ([4]byte, error) {
	return _WrappersV4.Contract.EIP1271MAGICVALUE(&_WrappersV4.CallOpts)
}

// EIP1271MAGICVALUE is a free data retrieval call binding the contract method 0xdd885e2d.
//
// Solidity: function EIP1271_MAGIC_VALUE() view returns(bytes4)
func (_WrappersV4 *WrappersV4CallerSession) EIP1271MAGICVALUE() ([4]byte, error) {
	return _WrappersV4.Contract.EIP1271MAGICVALUE(&_WrappersV4.CallOpts)
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

// AllowedValidators is a free data retrieval call binding the contract method 0x7b8e3514.
//
// Solidity: function allowedValidators(address , address ) view returns(bool)
func (_WrappersV4 *WrappersV4Caller) AllowedValidators(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "allowedValidators", arg0, arg1)
	return *ret0, err
}

// AllowedValidators is a free data retrieval call binding the contract method 0x7b8e3514.
//
// Solidity: function allowedValidators(address , address ) view returns(bool)
func (_WrappersV4 *WrappersV4Session) AllowedValidators(arg0 common.Address, arg1 common.Address) (bool, error) {
	return _WrappersV4.Contract.AllowedValidators(&_WrappersV4.CallOpts, arg0, arg1)
}

// AllowedValidators is a free data retrieval call binding the contract method 0x7b8e3514.
//
// Solidity: function allowedValidators(address , address ) view returns(bool)
func (_WrappersV4 *WrappersV4CallerSession) AllowedValidators(arg0 common.Address, arg1 common.Address) (bool, error) {
	return _WrappersV4.Contract.AllowedValidators(&_WrappersV4.CallOpts, arg0, arg1)
}

// Cancelled is a free data retrieval call binding the contract method 0x2ac12622.
//
// Solidity: function cancelled(bytes32 ) view returns(bool)
func (_WrappersV4 *WrappersV4Caller) Cancelled(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "cancelled", arg0)
	return *ret0, err
}

// Cancelled is a free data retrieval call binding the contract method 0x2ac12622.
//
// Solidity: function cancelled(bytes32 ) view returns(bool)
func (_WrappersV4 *WrappersV4Session) Cancelled(arg0 [32]byte) (bool, error) {
	return _WrappersV4.Contract.Cancelled(&_WrappersV4.CallOpts, arg0)
}

// Cancelled is a free data retrieval call binding the contract method 0x2ac12622.
//
// Solidity: function cancelled(bytes32 ) view returns(bool)
func (_WrappersV4 *WrappersV4CallerSession) Cancelled(arg0 [32]byte) (bool, error) {
	return _WrappersV4.Contract.Cancelled(&_WrappersV4.CallOpts, arg0)
}

// CurrentContextAddress is a free data retrieval call binding the contract method 0xeea086ba.
//
// Solidity: function currentContextAddress() view returns(address)
func (_WrappersV4 *WrappersV4Caller) CurrentContextAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "currentContextAddress")
	return *ret0, err
}

// CurrentContextAddress is a free data retrieval call binding the contract method 0xeea086ba.
//
// Solidity: function currentContextAddress() view returns(address)
func (_WrappersV4 *WrappersV4Session) CurrentContextAddress() (common.Address, error) {
	return _WrappersV4.Contract.CurrentContextAddress(&_WrappersV4.CallOpts)
}

// CurrentContextAddress is a free data retrieval call binding the contract method 0xeea086ba.
//
// Solidity: function currentContextAddress() view returns(address)
func (_WrappersV4 *WrappersV4CallerSession) CurrentContextAddress() (common.Address, error) {
	return _WrappersV4.Contract.CurrentContextAddress(&_WrappersV4.CallOpts)
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled(bytes32 ) view returns(uint256)
func (_WrappersV4 *WrappersV4Caller) Filled(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "filled", arg0)
	return *ret0, err
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled(bytes32 ) view returns(uint256)
func (_WrappersV4 *WrappersV4Session) Filled(arg0 [32]byte) (*big.Int, error) {
	return _WrappersV4.Contract.Filled(&_WrappersV4.CallOpts, arg0)
}

// Filled is a free data retrieval call binding the contract method 0x288cdc91.
//
// Solidity: function filled(bytes32 ) view returns(uint256)
func (_WrappersV4 *WrappersV4CallerSession) Filled(arg0 [32]byte) (*big.Int, error) {
	return _WrappersV4.Contract.Filled(&_WrappersV4.CallOpts, arg0)
}

// GetAssetProxy is a free data retrieval call binding the contract method 0x60704108.
//
// Solidity: function getAssetProxy(bytes4 assetProxyId) view returns(address assetProxy)
func (_WrappersV4 *WrappersV4Caller) GetAssetProxy(opts *bind.CallOpts, assetProxyId [4]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "getAssetProxy", assetProxyId)
	return *ret0, err
}

// GetAssetProxy is a free data retrieval call binding the contract method 0x60704108.
//
// Solidity: function getAssetProxy(bytes4 assetProxyId) view returns(address assetProxy)
func (_WrappersV4 *WrappersV4Session) GetAssetProxy(assetProxyId [4]byte) (common.Address, error) {
	return _WrappersV4.Contract.GetAssetProxy(&_WrappersV4.CallOpts, assetProxyId)
}

// GetAssetProxy is a free data retrieval call binding the contract method 0x60704108.
//
// Solidity: function getAssetProxy(bytes4 assetProxyId) view returns(address assetProxy)
func (_WrappersV4 *WrappersV4CallerSession) GetAssetProxy(assetProxyId [4]byte) (common.Address, error) {
	return _WrappersV4.Contract.GetAssetProxy(&_WrappersV4.CallOpts, assetProxyId)
}

// GetOrderInfo is a free data retrieval call binding the contract method 0x9d3fa4b9.
//
// Solidity: function getOrderInfo((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order) view returns((uint8,bytes32,uint256) orderInfo)
func (_WrappersV4 *WrappersV4Caller) GetOrderInfo(opts *bind.CallOpts, order LibOrderOrder) (LibOrderOrderInfo, error) {
	var (
		ret0 = new(LibOrderOrderInfo)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "getOrderInfo", order)
	return *ret0, err
}

// GetOrderInfo is a free data retrieval call binding the contract method 0x9d3fa4b9.
//
// Solidity: function getOrderInfo((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order) view returns((uint8,bytes32,uint256) orderInfo)
func (_WrappersV4 *WrappersV4Session) GetOrderInfo(order LibOrderOrder) (LibOrderOrderInfo, error) {
	return _WrappersV4.Contract.GetOrderInfo(&_WrappersV4.CallOpts, order)
}

// GetOrderInfo is a free data retrieval call binding the contract method 0x9d3fa4b9.
//
// Solidity: function getOrderInfo((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order) view returns((uint8,bytes32,uint256) orderInfo)
func (_WrappersV4 *WrappersV4CallerSession) GetOrderInfo(order LibOrderOrder) (LibOrderOrderInfo, error) {
	return _WrappersV4.Contract.GetOrderInfo(&_WrappersV4.CallOpts, order)
}

// IsValidHashSignature is a free data retrieval call binding the contract method 0x8171c407.
//
// Solidity: function isValidHashSignature(bytes32 hash, address signerAddress, bytes signature) view returns(bool isValid)
func (_WrappersV4 *WrappersV4Caller) IsValidHashSignature(opts *bind.CallOpts, hash [32]byte, signerAddress common.Address, signature []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "isValidHashSignature", hash, signerAddress, signature)
	return *ret0, err
}

// IsValidHashSignature is a free data retrieval call binding the contract method 0x8171c407.
//
// Solidity: function isValidHashSignature(bytes32 hash, address signerAddress, bytes signature) view returns(bool isValid)
func (_WrappersV4 *WrappersV4Session) IsValidHashSignature(hash [32]byte, signerAddress common.Address, signature []byte) (bool, error) {
	return _WrappersV4.Contract.IsValidHashSignature(&_WrappersV4.CallOpts, hash, signerAddress, signature)
}

// IsValidHashSignature is a free data retrieval call binding the contract method 0x8171c407.
//
// Solidity: function isValidHashSignature(bytes32 hash, address signerAddress, bytes signature) view returns(bool isValid)
func (_WrappersV4 *WrappersV4CallerSession) IsValidHashSignature(hash [32]byte, signerAddress common.Address, signature []byte) (bool, error) {
	return _WrappersV4.Contract.IsValidHashSignature(&_WrappersV4.CallOpts, hash, signerAddress, signature)
}

// IsValidOrderSignature is a free data retrieval call binding the contract method 0xa12dcc6f.
//
// Solidity: function isValidOrderSignature((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, bytes signature) view returns(bool isValid)
func (_WrappersV4 *WrappersV4Caller) IsValidOrderSignature(opts *bind.CallOpts, order LibOrderOrder, signature []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "isValidOrderSignature", order, signature)
	return *ret0, err
}

// IsValidOrderSignature is a free data retrieval call binding the contract method 0xa12dcc6f.
//
// Solidity: function isValidOrderSignature((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, bytes signature) view returns(bool isValid)
func (_WrappersV4 *WrappersV4Session) IsValidOrderSignature(order LibOrderOrder, signature []byte) (bool, error) {
	return _WrappersV4.Contract.IsValidOrderSignature(&_WrappersV4.CallOpts, order, signature)
}

// IsValidOrderSignature is a free data retrieval call binding the contract method 0xa12dcc6f.
//
// Solidity: function isValidOrderSignature((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, bytes signature) view returns(bool isValid)
func (_WrappersV4 *WrappersV4CallerSession) IsValidOrderSignature(order LibOrderOrder, signature []byte) (bool, error) {
	return _WrappersV4.Contract.IsValidOrderSignature(&_WrappersV4.CallOpts, order, signature)
}

// IsValidTransactionSignature is a free data retrieval call binding the contract method 0x8d45cd23.
//
// Solidity: function isValidTransactionSignature((uint256,uint256,uint256,address,bytes) transaction, bytes signature) view returns(bool isValid)
func (_WrappersV4 *WrappersV4Caller) IsValidTransactionSignature(opts *bind.CallOpts, transaction LibZeroExTransactionZeroExTransaction, signature []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "isValidTransactionSignature", transaction, signature)
	return *ret0, err
}

// IsValidTransactionSignature is a free data retrieval call binding the contract method 0x8d45cd23.
//
// Solidity: function isValidTransactionSignature((uint256,uint256,uint256,address,bytes) transaction, bytes signature) view returns(bool isValid)
func (_WrappersV4 *WrappersV4Session) IsValidTransactionSignature(transaction LibZeroExTransactionZeroExTransaction, signature []byte) (bool, error) {
	return _WrappersV4.Contract.IsValidTransactionSignature(&_WrappersV4.CallOpts, transaction, signature)
}

// IsValidTransactionSignature is a free data retrieval call binding the contract method 0x8d45cd23.
//
// Solidity: function isValidTransactionSignature((uint256,uint256,uint256,address,bytes) transaction, bytes signature) view returns(bool isValid)
func (_WrappersV4 *WrappersV4CallerSession) IsValidTransactionSignature(transaction LibZeroExTransactionZeroExTransaction, signature []byte) (bool, error) {
	return _WrappersV4.Contract.IsValidTransactionSignature(&_WrappersV4.CallOpts, transaction, signature)
}

// OrderEpoch is a free data retrieval call binding the contract method 0xd9bfa73e.
//
// Solidity: function orderEpoch(address , address ) view returns(uint256)
func (_WrappersV4 *WrappersV4Caller) OrderEpoch(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "orderEpoch", arg0, arg1)
	return *ret0, err
}

// OrderEpoch is a free data retrieval call binding the contract method 0xd9bfa73e.
//
// Solidity: function orderEpoch(address , address ) view returns(uint256)
func (_WrappersV4 *WrappersV4Session) OrderEpoch(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _WrappersV4.Contract.OrderEpoch(&_WrappersV4.CallOpts, arg0, arg1)
}

// OrderEpoch is a free data retrieval call binding the contract method 0xd9bfa73e.
//
// Solidity: function orderEpoch(address , address ) view returns(uint256)
func (_WrappersV4 *WrappersV4CallerSession) OrderEpoch(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _WrappersV4.Contract.OrderEpoch(&_WrappersV4.CallOpts, arg0, arg1)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WrappersV4 *WrappersV4Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WrappersV4 *WrappersV4Session) Owner() (common.Address, error) {
	return _WrappersV4.Contract.Owner(&_WrappersV4.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_WrappersV4 *WrappersV4CallerSession) Owner() (common.Address, error) {
	return _WrappersV4.Contract.Owner(&_WrappersV4.CallOpts)
}

// PreSigned is a free data retrieval call binding the contract method 0x82c174d0.
//
// Solidity: function preSigned(bytes32 , address ) view returns(bool)
func (_WrappersV4 *WrappersV4Caller) PreSigned(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "preSigned", arg0, arg1)
	return *ret0, err
}

// PreSigned is a free data retrieval call binding the contract method 0x82c174d0.
//
// Solidity: function preSigned(bytes32 , address ) view returns(bool)
func (_WrappersV4 *WrappersV4Session) PreSigned(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _WrappersV4.Contract.PreSigned(&_WrappersV4.CallOpts, arg0, arg1)
}

// PreSigned is a free data retrieval call binding the contract method 0x82c174d0.
//
// Solidity: function preSigned(bytes32 , address ) view returns(bool)
func (_WrappersV4 *WrappersV4CallerSession) PreSigned(arg0 [32]byte, arg1 common.Address) (bool, error) {
	return _WrappersV4.Contract.PreSigned(&_WrappersV4.CallOpts, arg0, arg1)
}

// ProtocolFeeCollector is a free data retrieval call binding the contract method 0x850a1501.
//
// Solidity: function protocolFeeCollector() view returns(address)
func (_WrappersV4 *WrappersV4Caller) ProtocolFeeCollector(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "protocolFeeCollector")
	return *ret0, err
}

// ProtocolFeeCollector is a free data retrieval call binding the contract method 0x850a1501.
//
// Solidity: function protocolFeeCollector() view returns(address)
func (_WrappersV4 *WrappersV4Session) ProtocolFeeCollector() (common.Address, error) {
	return _WrappersV4.Contract.ProtocolFeeCollector(&_WrappersV4.CallOpts)
}

// ProtocolFeeCollector is a free data retrieval call binding the contract method 0x850a1501.
//
// Solidity: function protocolFeeCollector() view returns(address)
func (_WrappersV4 *WrappersV4CallerSession) ProtocolFeeCollector() (common.Address, error) {
	return _WrappersV4.Contract.ProtocolFeeCollector(&_WrappersV4.CallOpts)
}

// ProtocolFeeMultiplier is a free data retrieval call binding the contract method 0x1ce4c78b.
//
// Solidity: function protocolFeeMultiplier() view returns(uint256)
func (_WrappersV4 *WrappersV4Caller) ProtocolFeeMultiplier(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "protocolFeeMultiplier")
	return *ret0, err
}

// ProtocolFeeMultiplier is a free data retrieval call binding the contract method 0x1ce4c78b.
//
// Solidity: function protocolFeeMultiplier() view returns(uint256)
func (_WrappersV4 *WrappersV4Session) ProtocolFeeMultiplier() (*big.Int, error) {
	return _WrappersV4.Contract.ProtocolFeeMultiplier(&_WrappersV4.CallOpts)
}

// ProtocolFeeMultiplier is a free data retrieval call binding the contract method 0x1ce4c78b.
//
// Solidity: function protocolFeeMultiplier() view returns(uint256)
func (_WrappersV4 *WrappersV4CallerSession) ProtocolFeeMultiplier() (*big.Int, error) {
	return _WrappersV4.Contract.ProtocolFeeMultiplier(&_WrappersV4.CallOpts)
}

// TransactionsExecuted is a free data retrieval call binding the contract method 0x0228e168.
//
// Solidity: function transactionsExecuted(bytes32 ) view returns(bool)
func (_WrappersV4 *WrappersV4Caller) TransactionsExecuted(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _WrappersV4.contract.Call(opts, out, "transactionsExecuted", arg0)
	return *ret0, err
}

// TransactionsExecuted is a free data retrieval call binding the contract method 0x0228e168.
//
// Solidity: function transactionsExecuted(bytes32 ) view returns(bool)
func (_WrappersV4 *WrappersV4Session) TransactionsExecuted(arg0 [32]byte) (bool, error) {
	return _WrappersV4.Contract.TransactionsExecuted(&_WrappersV4.CallOpts, arg0)
}

// TransactionsExecuted is a free data retrieval call binding the contract method 0x0228e168.
//
// Solidity: function transactionsExecuted(bytes32 ) view returns(bool)
func (_WrappersV4 *WrappersV4CallerSession) TransactionsExecuted(arg0 [32]byte) (bool, error) {
	return _WrappersV4.Contract.TransactionsExecuted(&_WrappersV4.CallOpts, arg0)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xdedfc1f1.
//
// Solidity: function batchCancelOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders) payable returns()
func (_WrappersV4 *WrappersV4Transactor) BatchCancelOrders(opts *bind.TransactOpts, orders []LibOrderOrder) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "batchCancelOrders", orders)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xdedfc1f1.
//
// Solidity: function batchCancelOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders) payable returns()
func (_WrappersV4 *WrappersV4Session) BatchCancelOrders(orders []LibOrderOrder) (*types.Transaction, error) {
	return _WrappersV4.Contract.BatchCancelOrders(&_WrappersV4.TransactOpts, orders)
}

// BatchCancelOrders is a paid mutator transaction binding the contract method 0xdedfc1f1.
//
// Solidity: function batchCancelOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders) payable returns()
func (_WrappersV4 *WrappersV4TransactorSession) BatchCancelOrders(orders []LibOrderOrder) (*types.Transaction, error) {
	return _WrappersV4.Contract.BatchCancelOrders(&_WrappersV4.TransactOpts, orders)
}

// BatchExecuteTransactions is a paid mutator transaction binding the contract method 0xfc74896d.
//
// Solidity: function batchExecuteTransactions((uint256,uint256,uint256,address,bytes)[] transactions, bytes[] signatures) payable returns(bytes[] returnData)
func (_WrappersV4 *WrappersV4Transactor) BatchExecuteTransactions(opts *bind.TransactOpts, transactions []LibZeroExTransactionZeroExTransaction, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "batchExecuteTransactions", transactions, signatures)
}

// BatchExecuteTransactions is a paid mutator transaction binding the contract method 0xfc74896d.
//
// Solidity: function batchExecuteTransactions((uint256,uint256,uint256,address,bytes)[] transactions, bytes[] signatures) payable returns(bytes[] returnData)
func (_WrappersV4 *WrappersV4Session) BatchExecuteTransactions(transactions []LibZeroExTransactionZeroExTransaction, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.BatchExecuteTransactions(&_WrappersV4.TransactOpts, transactions, signatures)
}

// BatchExecuteTransactions is a paid mutator transaction binding the contract method 0xfc74896d.
//
// Solidity: function batchExecuteTransactions((uint256,uint256,uint256,address,bytes)[] transactions, bytes[] signatures) payable returns(bytes[] returnData)
func (_WrappersV4 *WrappersV4TransactorSession) BatchExecuteTransactions(transactions []LibZeroExTransactionZeroExTransaction, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.BatchExecuteTransactions(&_WrappersV4.TransactOpts, transactions, signatures)
}

// BatchFillOrKillOrders is a paid mutator transaction binding the contract method 0xbeee2e14.
//
// Solidity: function batchFillOrKillOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_WrappersV4 *WrappersV4Transactor) BatchFillOrKillOrders(opts *bind.TransactOpts, orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "batchFillOrKillOrders", orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrKillOrders is a paid mutator transaction binding the contract method 0xbeee2e14.
//
// Solidity: function batchFillOrKillOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_WrappersV4 *WrappersV4Session) BatchFillOrKillOrders(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.BatchFillOrKillOrders(&_WrappersV4.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrKillOrders is a paid mutator transaction binding the contract method 0xbeee2e14.
//
// Solidity: function batchFillOrKillOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_WrappersV4 *WrappersV4TransactorSession) BatchFillOrKillOrders(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.BatchFillOrKillOrders(&_WrappersV4.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrders is a paid mutator transaction binding the contract method 0x9694a402.
//
// Solidity: function batchFillOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_WrappersV4 *WrappersV4Transactor) BatchFillOrders(opts *bind.TransactOpts, orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "batchFillOrders", orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrders is a paid mutator transaction binding the contract method 0x9694a402.
//
// Solidity: function batchFillOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_WrappersV4 *WrappersV4Session) BatchFillOrders(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.BatchFillOrders(&_WrappersV4.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrders is a paid mutator transaction binding the contract method 0x9694a402.
//
// Solidity: function batchFillOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_WrappersV4 *WrappersV4TransactorSession) BatchFillOrders(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.BatchFillOrders(&_WrappersV4.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrdersNoThrow is a paid mutator transaction binding the contract method 0x8ea8dfe4.
//
// Solidity: function batchFillOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_WrappersV4 *WrappersV4Transactor) BatchFillOrdersNoThrow(opts *bind.TransactOpts, orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "batchFillOrdersNoThrow", orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrdersNoThrow is a paid mutator transaction binding the contract method 0x8ea8dfe4.
//
// Solidity: function batchFillOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_WrappersV4 *WrappersV4Session) BatchFillOrdersNoThrow(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.BatchFillOrdersNoThrow(&_WrappersV4.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchFillOrdersNoThrow is a paid mutator transaction binding the contract method 0x8ea8dfe4.
//
// Solidity: function batchFillOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256[] takerAssetFillAmounts, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256)[] fillResults)
func (_WrappersV4 *WrappersV4TransactorSession) BatchFillOrdersNoThrow(orders []LibOrderOrder, takerAssetFillAmounts []*big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.BatchFillOrdersNoThrow(&_WrappersV4.TransactOpts, orders, takerAssetFillAmounts, signatures)
}

// BatchMatchOrders is a paid mutator transaction binding the contract method 0x6fcf3e9e.
//
// Solidity: function batchMatchOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] leftOrders, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) payable returns(((uint256,uint256,uint256,uint256,uint256)[],(uint256,uint256,uint256,uint256,uint256)[],uint256,uint256) batchMatchedFillResults)
func (_WrappersV4 *WrappersV4Transactor) BatchMatchOrders(opts *bind.TransactOpts, leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "batchMatchOrders", leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// BatchMatchOrders is a paid mutator transaction binding the contract method 0x6fcf3e9e.
//
// Solidity: function batchMatchOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] leftOrders, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) payable returns(((uint256,uint256,uint256,uint256,uint256)[],(uint256,uint256,uint256,uint256,uint256)[],uint256,uint256) batchMatchedFillResults)
func (_WrappersV4 *WrappersV4Session) BatchMatchOrders(leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.BatchMatchOrders(&_WrappersV4.TransactOpts, leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// BatchMatchOrders is a paid mutator transaction binding the contract method 0x6fcf3e9e.
//
// Solidity: function batchMatchOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] leftOrders, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) payable returns(((uint256,uint256,uint256,uint256,uint256)[],(uint256,uint256,uint256,uint256,uint256)[],uint256,uint256) batchMatchedFillResults)
func (_WrappersV4 *WrappersV4TransactorSession) BatchMatchOrders(leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.BatchMatchOrders(&_WrappersV4.TransactOpts, leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// BatchMatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0x6a1a80fd.
//
// Solidity: function batchMatchOrdersWithMaximalFill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] leftOrders, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) payable returns(((uint256,uint256,uint256,uint256,uint256)[],(uint256,uint256,uint256,uint256,uint256)[],uint256,uint256) batchMatchedFillResults)
func (_WrappersV4 *WrappersV4Transactor) BatchMatchOrdersWithMaximalFill(opts *bind.TransactOpts, leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "batchMatchOrdersWithMaximalFill", leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// BatchMatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0x6a1a80fd.
//
// Solidity: function batchMatchOrdersWithMaximalFill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] leftOrders, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) payable returns(((uint256,uint256,uint256,uint256,uint256)[],(uint256,uint256,uint256,uint256,uint256)[],uint256,uint256) batchMatchedFillResults)
func (_WrappersV4 *WrappersV4Session) BatchMatchOrdersWithMaximalFill(leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.BatchMatchOrdersWithMaximalFill(&_WrappersV4.TransactOpts, leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// BatchMatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0x6a1a80fd.
//
// Solidity: function batchMatchOrdersWithMaximalFill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] leftOrders, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] rightOrders, bytes[] leftSignatures, bytes[] rightSignatures) payable returns(((uint256,uint256,uint256,uint256,uint256)[],(uint256,uint256,uint256,uint256,uint256)[],uint256,uint256) batchMatchedFillResults)
func (_WrappersV4 *WrappersV4TransactorSession) BatchMatchOrdersWithMaximalFill(leftOrders []LibOrderOrder, rightOrders []LibOrderOrder, leftSignatures [][]byte, rightSignatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.BatchMatchOrdersWithMaximalFill(&_WrappersV4.TransactOpts, leftOrders, rightOrders, leftSignatures, rightSignatures)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x2da62987.
//
// Solidity: function cancelOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order) payable returns()
func (_WrappersV4 *WrappersV4Transactor) CancelOrder(opts *bind.TransactOpts, order LibOrderOrder) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "cancelOrder", order)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x2da62987.
//
// Solidity: function cancelOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order) payable returns()
func (_WrappersV4 *WrappersV4Session) CancelOrder(order LibOrderOrder) (*types.Transaction, error) {
	return _WrappersV4.Contract.CancelOrder(&_WrappersV4.TransactOpts, order)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x2da62987.
//
// Solidity: function cancelOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order) payable returns()
func (_WrappersV4 *WrappersV4TransactorSession) CancelOrder(order LibOrderOrder) (*types.Transaction, error) {
	return _WrappersV4.Contract.CancelOrder(&_WrappersV4.TransactOpts, order)
}

// CancelOrdersUpTo is a paid mutator transaction binding the contract method 0x4f9559b1.
//
// Solidity: function cancelOrdersUpTo(uint256 targetOrderEpoch) payable returns()
func (_WrappersV4 *WrappersV4Transactor) CancelOrdersUpTo(opts *bind.TransactOpts, targetOrderEpoch *big.Int) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "cancelOrdersUpTo", targetOrderEpoch)
}

// CancelOrdersUpTo is a paid mutator transaction binding the contract method 0x4f9559b1.
//
// Solidity: function cancelOrdersUpTo(uint256 targetOrderEpoch) payable returns()
func (_WrappersV4 *WrappersV4Session) CancelOrdersUpTo(targetOrderEpoch *big.Int) (*types.Transaction, error) {
	return _WrappersV4.Contract.CancelOrdersUpTo(&_WrappersV4.TransactOpts, targetOrderEpoch)
}

// CancelOrdersUpTo is a paid mutator transaction binding the contract method 0x4f9559b1.
//
// Solidity: function cancelOrdersUpTo(uint256 targetOrderEpoch) payable returns()
func (_WrappersV4 *WrappersV4TransactorSession) CancelOrdersUpTo(targetOrderEpoch *big.Int) (*types.Transaction, error) {
	return _WrappersV4.Contract.CancelOrdersUpTo(&_WrappersV4.TransactOpts, targetOrderEpoch)
}

// DetachProtocolFeeCollector is a paid mutator transaction binding the contract method 0x0efca185.
//
// Solidity: function detachProtocolFeeCollector() returns()
func (_WrappersV4 *WrappersV4Transactor) DetachProtocolFeeCollector(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "detachProtocolFeeCollector")
}

// DetachProtocolFeeCollector is a paid mutator transaction binding the contract method 0x0efca185.
//
// Solidity: function detachProtocolFeeCollector() returns()
func (_WrappersV4 *WrappersV4Session) DetachProtocolFeeCollector() (*types.Transaction, error) {
	return _WrappersV4.Contract.DetachProtocolFeeCollector(&_WrappersV4.TransactOpts)
}

// DetachProtocolFeeCollector is a paid mutator transaction binding the contract method 0x0efca185.
//
// Solidity: function detachProtocolFeeCollector() returns()
func (_WrappersV4 *WrappersV4TransactorSession) DetachProtocolFeeCollector() (*types.Transaction, error) {
	return _WrappersV4.Contract.DetachProtocolFeeCollector(&_WrappersV4.TransactOpts)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0x2280c910.
//
// Solidity: function executeTransaction((uint256,uint256,uint256,address,bytes) transaction, bytes signature) payable returns(bytes)
func (_WrappersV4 *WrappersV4Transactor) ExecuteTransaction(opts *bind.TransactOpts, transaction LibZeroExTransactionZeroExTransaction, signature []byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "executeTransaction", transaction, signature)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0x2280c910.
//
// Solidity: function executeTransaction((uint256,uint256,uint256,address,bytes) transaction, bytes signature) payable returns(bytes)
func (_WrappersV4 *WrappersV4Session) ExecuteTransaction(transaction LibZeroExTransactionZeroExTransaction, signature []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.ExecuteTransaction(&_WrappersV4.TransactOpts, transaction, signature)
}

// ExecuteTransaction is a paid mutator transaction binding the contract method 0x2280c910.
//
// Solidity: function executeTransaction((uint256,uint256,uint256,address,bytes) transaction, bytes signature) payable returns(bytes)
func (_WrappersV4 *WrappersV4TransactorSession) ExecuteTransaction(transaction LibZeroExTransactionZeroExTransaction, signature []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.ExecuteTransaction(&_WrappersV4.TransactOpts, transaction, signature)
}

// FillOrKillOrder is a paid mutator transaction binding the contract method 0xe14b58c4.
//
// Solidity: function fillOrKillOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 takerAssetFillAmount, bytes signature) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4Transactor) FillOrKillOrder(opts *bind.TransactOpts, order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "fillOrKillOrder", order, takerAssetFillAmount, signature)
}

// FillOrKillOrder is a paid mutator transaction binding the contract method 0xe14b58c4.
//
// Solidity: function fillOrKillOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 takerAssetFillAmount, bytes signature) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4Session) FillOrKillOrder(order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.FillOrKillOrder(&_WrappersV4.TransactOpts, order, takerAssetFillAmount, signature)
}

// FillOrKillOrder is a paid mutator transaction binding the contract method 0xe14b58c4.
//
// Solidity: function fillOrKillOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 takerAssetFillAmount, bytes signature) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4TransactorSession) FillOrKillOrder(order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.FillOrKillOrder(&_WrappersV4.TransactOpts, order, takerAssetFillAmount, signature)
}

// FillOrder is a paid mutator transaction binding the contract method 0x9b44d556.
//
// Solidity: function fillOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 takerAssetFillAmount, bytes signature) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4Transactor) FillOrder(opts *bind.TransactOpts, order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "fillOrder", order, takerAssetFillAmount, signature)
}

// FillOrder is a paid mutator transaction binding the contract method 0x9b44d556.
//
// Solidity: function fillOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 takerAssetFillAmount, bytes signature) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4Session) FillOrder(order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.FillOrder(&_WrappersV4.TransactOpts, order, takerAssetFillAmount, signature)
}

// FillOrder is a paid mutator transaction binding the contract method 0x9b44d556.
//
// Solidity: function fillOrder((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) order, uint256 takerAssetFillAmount, bytes signature) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4TransactorSession) FillOrder(order LibOrderOrder, takerAssetFillAmount *big.Int, signature []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.FillOrder(&_WrappersV4.TransactOpts, order, takerAssetFillAmount, signature)
}

// MarketBuyOrdersFillOrKill is a paid mutator transaction binding the contract method 0x8bc8efb3.
//
// Solidity: function marketBuyOrdersFillOrKill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 makerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4Transactor) MarketBuyOrdersFillOrKill(opts *bind.TransactOpts, orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "marketBuyOrdersFillOrKill", orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersFillOrKill is a paid mutator transaction binding the contract method 0x8bc8efb3.
//
// Solidity: function marketBuyOrdersFillOrKill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 makerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4Session) MarketBuyOrdersFillOrKill(orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.MarketBuyOrdersFillOrKill(&_WrappersV4.TransactOpts, orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersFillOrKill is a paid mutator transaction binding the contract method 0x8bc8efb3.
//
// Solidity: function marketBuyOrdersFillOrKill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 makerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4TransactorSession) MarketBuyOrdersFillOrKill(orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.MarketBuyOrdersFillOrKill(&_WrappersV4.TransactOpts, orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersNoThrow is a paid mutator transaction binding the contract method 0x78d29ac1.
//
// Solidity: function marketBuyOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 makerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4Transactor) MarketBuyOrdersNoThrow(opts *bind.TransactOpts, orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "marketBuyOrdersNoThrow", orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersNoThrow is a paid mutator transaction binding the contract method 0x78d29ac1.
//
// Solidity: function marketBuyOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 makerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4Session) MarketBuyOrdersNoThrow(orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.MarketBuyOrdersNoThrow(&_WrappersV4.TransactOpts, orders, makerAssetFillAmount, signatures)
}

// MarketBuyOrdersNoThrow is a paid mutator transaction binding the contract method 0x78d29ac1.
//
// Solidity: function marketBuyOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 makerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4TransactorSession) MarketBuyOrdersNoThrow(orders []LibOrderOrder, makerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.MarketBuyOrdersNoThrow(&_WrappersV4.TransactOpts, orders, makerAssetFillAmount, signatures)
}

// MarketSellOrdersFillOrKill is a paid mutator transaction binding the contract method 0xa6c3bf33.
//
// Solidity: function marketSellOrdersFillOrKill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 takerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4Transactor) MarketSellOrdersFillOrKill(opts *bind.TransactOpts, orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "marketSellOrdersFillOrKill", orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersFillOrKill is a paid mutator transaction binding the contract method 0xa6c3bf33.
//
// Solidity: function marketSellOrdersFillOrKill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 takerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4Session) MarketSellOrdersFillOrKill(orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.MarketSellOrdersFillOrKill(&_WrappersV4.TransactOpts, orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersFillOrKill is a paid mutator transaction binding the contract method 0xa6c3bf33.
//
// Solidity: function marketSellOrdersFillOrKill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 takerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4TransactorSession) MarketSellOrdersFillOrKill(orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.MarketSellOrdersFillOrKill(&_WrappersV4.TransactOpts, orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersNoThrow is a paid mutator transaction binding the contract method 0x369da099.
//
// Solidity: function marketSellOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 takerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4Transactor) MarketSellOrdersNoThrow(opts *bind.TransactOpts, orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "marketSellOrdersNoThrow", orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersNoThrow is a paid mutator transaction binding the contract method 0x369da099.
//
// Solidity: function marketSellOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 takerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4Session) MarketSellOrdersNoThrow(orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.MarketSellOrdersNoThrow(&_WrappersV4.TransactOpts, orders, takerAssetFillAmount, signatures)
}

// MarketSellOrdersNoThrow is a paid mutator transaction binding the contract method 0x369da099.
//
// Solidity: function marketSellOrdersNoThrow((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes)[] orders, uint256 takerAssetFillAmount, bytes[] signatures) payable returns((uint256,uint256,uint256,uint256,uint256) fillResults)
func (_WrappersV4 *WrappersV4TransactorSession) MarketSellOrdersNoThrow(orders []LibOrderOrder, takerAssetFillAmount *big.Int, signatures [][]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.MarketSellOrdersNoThrow(&_WrappersV4.TransactOpts, orders, takerAssetFillAmount, signatures)
}

// MatchOrders is a paid mutator transaction binding the contract method 0x88ec79fb.
//
// Solidity: function matchOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) leftOrder, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) rightOrder, bytes leftSignature, bytes rightSignature) payable returns(((uint256,uint256,uint256,uint256,uint256),(uint256,uint256,uint256,uint256,uint256),uint256,uint256) matchedFillResults)
func (_WrappersV4 *WrappersV4Transactor) MatchOrders(opts *bind.TransactOpts, leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "matchOrders", leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrders is a paid mutator transaction binding the contract method 0x88ec79fb.
//
// Solidity: function matchOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) leftOrder, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) rightOrder, bytes leftSignature, bytes rightSignature) payable returns(((uint256,uint256,uint256,uint256,uint256),(uint256,uint256,uint256,uint256,uint256),uint256,uint256) matchedFillResults)
func (_WrappersV4 *WrappersV4Session) MatchOrders(leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.MatchOrders(&_WrappersV4.TransactOpts, leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrders is a paid mutator transaction binding the contract method 0x88ec79fb.
//
// Solidity: function matchOrders((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) leftOrder, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) rightOrder, bytes leftSignature, bytes rightSignature) payable returns(((uint256,uint256,uint256,uint256,uint256),(uint256,uint256,uint256,uint256,uint256),uint256,uint256) matchedFillResults)
func (_WrappersV4 *WrappersV4TransactorSession) MatchOrders(leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.MatchOrders(&_WrappersV4.TransactOpts, leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0xb718e292.
//
// Solidity: function matchOrdersWithMaximalFill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) leftOrder, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) rightOrder, bytes leftSignature, bytes rightSignature) payable returns(((uint256,uint256,uint256,uint256,uint256),(uint256,uint256,uint256,uint256,uint256),uint256,uint256) matchedFillResults)
func (_WrappersV4 *WrappersV4Transactor) MatchOrdersWithMaximalFill(opts *bind.TransactOpts, leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "matchOrdersWithMaximalFill", leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0xb718e292.
//
// Solidity: function matchOrdersWithMaximalFill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) leftOrder, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) rightOrder, bytes leftSignature, bytes rightSignature) payable returns(((uint256,uint256,uint256,uint256,uint256),(uint256,uint256,uint256,uint256,uint256),uint256,uint256) matchedFillResults)
func (_WrappersV4 *WrappersV4Session) MatchOrdersWithMaximalFill(leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.MatchOrdersWithMaximalFill(&_WrappersV4.TransactOpts, leftOrder, rightOrder, leftSignature, rightSignature)
}

// MatchOrdersWithMaximalFill is a paid mutator transaction binding the contract method 0xb718e292.
//
// Solidity: function matchOrdersWithMaximalFill((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) leftOrder, (address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes,bytes,bytes) rightOrder, bytes leftSignature, bytes rightSignature) payable returns(((uint256,uint256,uint256,uint256,uint256),(uint256,uint256,uint256,uint256,uint256),uint256,uint256) matchedFillResults)
func (_WrappersV4 *WrappersV4TransactorSession) MatchOrdersWithMaximalFill(leftOrder LibOrderOrder, rightOrder LibOrderOrder, leftSignature []byte, rightSignature []byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.MatchOrdersWithMaximalFill(&_WrappersV4.TransactOpts, leftOrder, rightOrder, leftSignature, rightSignature)
}

// PreSign is a paid mutator transaction binding the contract method 0x46c02d7a.
//
// Solidity: function preSign(bytes32 hash) payable returns()
func (_WrappersV4 *WrappersV4Transactor) PreSign(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "preSign", hash)
}

// PreSign is a paid mutator transaction binding the contract method 0x46c02d7a.
//
// Solidity: function preSign(bytes32 hash) payable returns()
func (_WrappersV4 *WrappersV4Session) PreSign(hash [32]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.PreSign(&_WrappersV4.TransactOpts, hash)
}

// PreSign is a paid mutator transaction binding the contract method 0x46c02d7a.
//
// Solidity: function preSign(bytes32 hash) payable returns()
func (_WrappersV4 *WrappersV4TransactorSession) PreSign(hash [32]byte) (*types.Transaction, error) {
	return _WrappersV4.Contract.PreSign(&_WrappersV4.TransactOpts, hash)
}

// RegisterAssetProxy is a paid mutator transaction binding the contract method 0xc585bb93.
//
// Solidity: function registerAssetProxy(address assetProxy) returns()
func (_WrappersV4 *WrappersV4Transactor) RegisterAssetProxy(opts *bind.TransactOpts, assetProxy common.Address) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "registerAssetProxy", assetProxy)
}

// RegisterAssetProxy is a paid mutator transaction binding the contract method 0xc585bb93.
//
// Solidity: function registerAssetProxy(address assetProxy) returns()
func (_WrappersV4 *WrappersV4Session) RegisterAssetProxy(assetProxy common.Address) (*types.Transaction, error) {
	return _WrappersV4.Contract.RegisterAssetProxy(&_WrappersV4.TransactOpts, assetProxy)
}

// RegisterAssetProxy is a paid mutator transaction binding the contract method 0xc585bb93.
//
// Solidity: function registerAssetProxy(address assetProxy) returns()
func (_WrappersV4 *WrappersV4TransactorSession) RegisterAssetProxy(assetProxy common.Address) (*types.Transaction, error) {
	return _WrappersV4.Contract.RegisterAssetProxy(&_WrappersV4.TransactOpts, assetProxy)
}

// SetProtocolFeeCollectorAddress is a paid mutator transaction binding the contract method 0xc0fa16cc.
//
// Solidity: function setProtocolFeeCollectorAddress(address updatedProtocolFeeCollector) returns()
func (_WrappersV4 *WrappersV4Transactor) SetProtocolFeeCollectorAddress(opts *bind.TransactOpts, updatedProtocolFeeCollector common.Address) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "setProtocolFeeCollectorAddress", updatedProtocolFeeCollector)
}

// SetProtocolFeeCollectorAddress is a paid mutator transaction binding the contract method 0xc0fa16cc.
//
// Solidity: function setProtocolFeeCollectorAddress(address updatedProtocolFeeCollector) returns()
func (_WrappersV4 *WrappersV4Session) SetProtocolFeeCollectorAddress(updatedProtocolFeeCollector common.Address) (*types.Transaction, error) {
	return _WrappersV4.Contract.SetProtocolFeeCollectorAddress(&_WrappersV4.TransactOpts, updatedProtocolFeeCollector)
}

// SetProtocolFeeCollectorAddress is a paid mutator transaction binding the contract method 0xc0fa16cc.
//
// Solidity: function setProtocolFeeCollectorAddress(address updatedProtocolFeeCollector) returns()
func (_WrappersV4 *WrappersV4TransactorSession) SetProtocolFeeCollectorAddress(updatedProtocolFeeCollector common.Address) (*types.Transaction, error) {
	return _WrappersV4.Contract.SetProtocolFeeCollectorAddress(&_WrappersV4.TransactOpts, updatedProtocolFeeCollector)
}

// SetProtocolFeeMultiplier is a paid mutator transaction binding the contract method 0x9331c742.
//
// Solidity: function setProtocolFeeMultiplier(uint256 updatedProtocolFeeMultiplier) returns()
func (_WrappersV4 *WrappersV4Transactor) SetProtocolFeeMultiplier(opts *bind.TransactOpts, updatedProtocolFeeMultiplier *big.Int) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "setProtocolFeeMultiplier", updatedProtocolFeeMultiplier)
}

// SetProtocolFeeMultiplier is a paid mutator transaction binding the contract method 0x9331c742.
//
// Solidity: function setProtocolFeeMultiplier(uint256 updatedProtocolFeeMultiplier) returns()
func (_WrappersV4 *WrappersV4Session) SetProtocolFeeMultiplier(updatedProtocolFeeMultiplier *big.Int) (*types.Transaction, error) {
	return _WrappersV4.Contract.SetProtocolFeeMultiplier(&_WrappersV4.TransactOpts, updatedProtocolFeeMultiplier)
}

// SetProtocolFeeMultiplier is a paid mutator transaction binding the contract method 0x9331c742.
//
// Solidity: function setProtocolFeeMultiplier(uint256 updatedProtocolFeeMultiplier) returns()
func (_WrappersV4 *WrappersV4TransactorSession) SetProtocolFeeMultiplier(updatedProtocolFeeMultiplier *big.Int) (*types.Transaction, error) {
	return _WrappersV4.Contract.SetProtocolFeeMultiplier(&_WrappersV4.TransactOpts, updatedProtocolFeeMultiplier)
}

// SetSignatureValidatorApproval is a paid mutator transaction binding the contract method 0x77fcce68.
//
// Solidity: function setSignatureValidatorApproval(address validatorAddress, bool approval) payable returns()
func (_WrappersV4 *WrappersV4Transactor) SetSignatureValidatorApproval(opts *bind.TransactOpts, validatorAddress common.Address, approval bool) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "setSignatureValidatorApproval", validatorAddress, approval)
}

// SetSignatureValidatorApproval is a paid mutator transaction binding the contract method 0x77fcce68.
//
// Solidity: function setSignatureValidatorApproval(address validatorAddress, bool approval) payable returns()
func (_WrappersV4 *WrappersV4Session) SetSignatureValidatorApproval(validatorAddress common.Address, approval bool) (*types.Transaction, error) {
	return _WrappersV4.Contract.SetSignatureValidatorApproval(&_WrappersV4.TransactOpts, validatorAddress, approval)
}

// SetSignatureValidatorApproval is a paid mutator transaction binding the contract method 0x77fcce68.
//
// Solidity: function setSignatureValidatorApproval(address validatorAddress, bool approval) payable returns()
func (_WrappersV4 *WrappersV4TransactorSession) SetSignatureValidatorApproval(validatorAddress common.Address, approval bool) (*types.Transaction, error) {
	return _WrappersV4.Contract.SetSignatureValidatorApproval(&_WrappersV4.TransactOpts, validatorAddress, approval)
}

// SimulateDispatchTransferFromCalls is a paid mutator transaction binding the contract method 0xb04fbddd.
//
// Solidity: function simulateDispatchTransferFromCalls(bytes[] assetData, address[] fromAddresses, address[] toAddresses, uint256[] amounts) returns()
func (_WrappersV4 *WrappersV4Transactor) SimulateDispatchTransferFromCalls(opts *bind.TransactOpts, assetData [][]byte, fromAddresses []common.Address, toAddresses []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "simulateDispatchTransferFromCalls", assetData, fromAddresses, toAddresses, amounts)
}

// SimulateDispatchTransferFromCalls is a paid mutator transaction binding the contract method 0xb04fbddd.
//
// Solidity: function simulateDispatchTransferFromCalls(bytes[] assetData, address[] fromAddresses, address[] toAddresses, uint256[] amounts) returns()
func (_WrappersV4 *WrappersV4Session) SimulateDispatchTransferFromCalls(assetData [][]byte, fromAddresses []common.Address, toAddresses []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _WrappersV4.Contract.SimulateDispatchTransferFromCalls(&_WrappersV4.TransactOpts, assetData, fromAddresses, toAddresses, amounts)
}

// SimulateDispatchTransferFromCalls is a paid mutator transaction binding the contract method 0xb04fbddd.
//
// Solidity: function simulateDispatchTransferFromCalls(bytes[] assetData, address[] fromAddresses, address[] toAddresses, uint256[] amounts) returns()
func (_WrappersV4 *WrappersV4TransactorSession) SimulateDispatchTransferFromCalls(assetData [][]byte, fromAddresses []common.Address, toAddresses []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _WrappersV4.Contract.SimulateDispatchTransferFromCalls(&_WrappersV4.TransactOpts, assetData, fromAddresses, toAddresses, amounts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WrappersV4 *WrappersV4Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _WrappersV4.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WrappersV4 *WrappersV4Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WrappersV4.Contract.TransferOwnership(&_WrappersV4.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_WrappersV4 *WrappersV4TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _WrappersV4.Contract.TransferOwnership(&_WrappersV4.TransactOpts, newOwner)
}

// WrappersV4AssetProxyRegisteredIterator is returned from FilterAssetProxyRegistered and is used to iterate over the raw logs and unpacked data for AssetProxyRegistered events raised by the WrappersV4 contract.
type WrappersV4AssetProxyRegisteredIterator struct {
	Event *WrappersV4AssetProxyRegistered // Event containing the contract specifics and raw log

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
func (it *WrappersV4AssetProxyRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappersV4AssetProxyRegistered)
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
		it.Event = new(WrappersV4AssetProxyRegistered)
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
func (it *WrappersV4AssetProxyRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappersV4AssetProxyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappersV4AssetProxyRegistered represents a AssetProxyRegistered event raised by the WrappersV4 contract.
type WrappersV4AssetProxyRegistered struct {
	Id         [4]byte
	AssetProxy common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAssetProxyRegistered is a free log retrieval operation binding the contract event 0xd2c6b762299c609bdb96520b58a49bfb80186934d4f71a86a367571a15c03194.
//
// Solidity: event AssetProxyRegistered(bytes4 id, address assetProxy)
func (_WrappersV4 *WrappersV4Filterer) FilterAssetProxyRegistered(opts *bind.FilterOpts) (*WrappersV4AssetProxyRegisteredIterator, error) {

	logs, sub, err := _WrappersV4.contract.FilterLogs(opts, "AssetProxyRegistered")
	if err != nil {
		return nil, err
	}
	return &WrappersV4AssetProxyRegisteredIterator{contract: _WrappersV4.contract, event: "AssetProxyRegistered", logs: logs, sub: sub}, nil
}

// WatchAssetProxyRegistered is a free log subscription operation binding the contract event 0xd2c6b762299c609bdb96520b58a49bfb80186934d4f71a86a367571a15c03194.
//
// Solidity: event AssetProxyRegistered(bytes4 id, address assetProxy)
func (_WrappersV4 *WrappersV4Filterer) WatchAssetProxyRegistered(opts *bind.WatchOpts, sink chan<- *WrappersV4AssetProxyRegistered) (event.Subscription, error) {

	logs, sub, err := _WrappersV4.contract.WatchLogs(opts, "AssetProxyRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappersV4AssetProxyRegistered)
				if err := _WrappersV4.contract.UnpackLog(event, "AssetProxyRegistered", log); err != nil {
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
func (_WrappersV4 *WrappersV4Filterer) ParseAssetProxyRegistered(log types.Log) (*WrappersV4AssetProxyRegistered, error) {
	event := new(WrappersV4AssetProxyRegistered)
	if err := _WrappersV4.contract.UnpackLog(event, "AssetProxyRegistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WrappersV4CancelIterator is returned from FilterCancel and is used to iterate over the raw logs and unpacked data for Cancel events raised by the WrappersV4 contract.
type WrappersV4CancelIterator struct {
	Event *WrappersV4Cancel // Event containing the contract specifics and raw log

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
func (it *WrappersV4CancelIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappersV4Cancel)
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
		it.Event = new(WrappersV4Cancel)
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
func (it *WrappersV4CancelIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappersV4CancelIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappersV4Cancel represents a Cancel event raised by the WrappersV4 contract.
type WrappersV4Cancel struct {
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
func (_WrappersV4 *WrappersV4Filterer) FilterCancel(opts *bind.FilterOpts, makerAddress []common.Address, feeRecipientAddress []common.Address, orderHash [][32]byte) (*WrappersV4CancelIterator, error) {

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

	logs, sub, err := _WrappersV4.contract.FilterLogs(opts, "Cancel", makerAddressRule, feeRecipientAddressRule, orderHashRule)
	if err != nil {
		return nil, err
	}
	return &WrappersV4CancelIterator{contract: _WrappersV4.contract, event: "Cancel", logs: logs, sub: sub}, nil
}

// WatchCancel is a free log subscription operation binding the contract event 0x02c310a9a43963ff31a754a4099cc435ed498049687539d72d7818d9b093415c.
//
// Solidity: event Cancel(address indexed makerAddress, address indexed feeRecipientAddress, bytes makerAssetData, bytes takerAssetData, address senderAddress, bytes32 indexed orderHash)
func (_WrappersV4 *WrappersV4Filterer) WatchCancel(opts *bind.WatchOpts, sink chan<- *WrappersV4Cancel, makerAddress []common.Address, feeRecipientAddress []common.Address, orderHash [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _WrappersV4.contract.WatchLogs(opts, "Cancel", makerAddressRule, feeRecipientAddressRule, orderHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappersV4Cancel)
				if err := _WrappersV4.contract.UnpackLog(event, "Cancel", log); err != nil {
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
func (_WrappersV4 *WrappersV4Filterer) ParseCancel(log types.Log) (*WrappersV4Cancel, error) {
	event := new(WrappersV4Cancel)
	if err := _WrappersV4.contract.UnpackLog(event, "Cancel", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WrappersV4CancelUpToIterator is returned from FilterCancelUpTo and is used to iterate over the raw logs and unpacked data for CancelUpTo events raised by the WrappersV4 contract.
type WrappersV4CancelUpToIterator struct {
	Event *WrappersV4CancelUpTo // Event containing the contract specifics and raw log

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
func (it *WrappersV4CancelUpToIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappersV4CancelUpTo)
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
		it.Event = new(WrappersV4CancelUpTo)
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
func (it *WrappersV4CancelUpToIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappersV4CancelUpToIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappersV4CancelUpTo represents a CancelUpTo event raised by the WrappersV4 contract.
type WrappersV4CancelUpTo struct {
	MakerAddress       common.Address
	OrderSenderAddress common.Address
	OrderEpoch         *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterCancelUpTo is a free log retrieval operation binding the contract event 0x82af639571738f4ebd4268fb0363d8957ebe1bbb9e78dba5ebd69eed39b154f0.
//
// Solidity: event CancelUpTo(address indexed makerAddress, address indexed orderSenderAddress, uint256 orderEpoch)
func (_WrappersV4 *WrappersV4Filterer) FilterCancelUpTo(opts *bind.FilterOpts, makerAddress []common.Address, orderSenderAddress []common.Address) (*WrappersV4CancelUpToIterator, error) {

	var makerAddressRule []interface{}
	for _, makerAddressItem := range makerAddress {
		makerAddressRule = append(makerAddressRule, makerAddressItem)
	}
	var orderSenderAddressRule []interface{}
	for _, orderSenderAddressItem := range orderSenderAddress {
		orderSenderAddressRule = append(orderSenderAddressRule, orderSenderAddressItem)
	}

	logs, sub, err := _WrappersV4.contract.FilterLogs(opts, "CancelUpTo", makerAddressRule, orderSenderAddressRule)
	if err != nil {
		return nil, err
	}
	return &WrappersV4CancelUpToIterator{contract: _WrappersV4.contract, event: "CancelUpTo", logs: logs, sub: sub}, nil
}

// WatchCancelUpTo is a free log subscription operation binding the contract event 0x82af639571738f4ebd4268fb0363d8957ebe1bbb9e78dba5ebd69eed39b154f0.
//
// Solidity: event CancelUpTo(address indexed makerAddress, address indexed orderSenderAddress, uint256 orderEpoch)
func (_WrappersV4 *WrappersV4Filterer) WatchCancelUpTo(opts *bind.WatchOpts, sink chan<- *WrappersV4CancelUpTo, makerAddress []common.Address, orderSenderAddress []common.Address) (event.Subscription, error) {

	var makerAddressRule []interface{}
	for _, makerAddressItem := range makerAddress {
		makerAddressRule = append(makerAddressRule, makerAddressItem)
	}
	var orderSenderAddressRule []interface{}
	for _, orderSenderAddressItem := range orderSenderAddress {
		orderSenderAddressRule = append(orderSenderAddressRule, orderSenderAddressItem)
	}

	logs, sub, err := _WrappersV4.contract.WatchLogs(opts, "CancelUpTo", makerAddressRule, orderSenderAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappersV4CancelUpTo)
				if err := _WrappersV4.contract.UnpackLog(event, "CancelUpTo", log); err != nil {
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
func (_WrappersV4 *WrappersV4Filterer) ParseCancelUpTo(log types.Log) (*WrappersV4CancelUpTo, error) {
	event := new(WrappersV4CancelUpTo)
	if err := _WrappersV4.contract.UnpackLog(event, "CancelUpTo", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WrappersV4FillIterator is returned from FilterFill and is used to iterate over the raw logs and unpacked data for Fill events raised by the WrappersV4 contract.
type WrappersV4FillIterator struct {
	Event *WrappersV4Fill // Event containing the contract specifics and raw log

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
func (it *WrappersV4FillIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappersV4Fill)
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
		it.Event = new(WrappersV4Fill)
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
func (it *WrappersV4FillIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappersV4FillIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappersV4Fill represents a Fill event raised by the WrappersV4 contract.
type WrappersV4Fill struct {
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
func (_WrappersV4 *WrappersV4Filterer) FilterFill(opts *bind.FilterOpts, makerAddress []common.Address, feeRecipientAddress []common.Address, orderHash [][32]byte) (*WrappersV4FillIterator, error) {

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

	logs, sub, err := _WrappersV4.contract.FilterLogs(opts, "Fill", makerAddressRule, feeRecipientAddressRule, orderHashRule)
	if err != nil {
		return nil, err
	}
	return &WrappersV4FillIterator{contract: _WrappersV4.contract, event: "Fill", logs: logs, sub: sub}, nil
}

// WatchFill is a free log subscription operation binding the contract event 0x6869791f0a34781b29882982cc39e882768cf2c96995c2a110c577c53bc932d5.
//
// Solidity: event Fill(address indexed makerAddress, address indexed feeRecipientAddress, bytes makerAssetData, bytes takerAssetData, bytes makerFeeAssetData, bytes takerFeeAssetData, bytes32 indexed orderHash, address takerAddress, address senderAddress, uint256 makerAssetFilledAmount, uint256 takerAssetFilledAmount, uint256 makerFeePaid, uint256 takerFeePaid, uint256 protocolFeePaid)
func (_WrappersV4 *WrappersV4Filterer) WatchFill(opts *bind.WatchOpts, sink chan<- *WrappersV4Fill, makerAddress []common.Address, feeRecipientAddress []common.Address, orderHash [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _WrappersV4.contract.WatchLogs(opts, "Fill", makerAddressRule, feeRecipientAddressRule, orderHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappersV4Fill)
				if err := _WrappersV4.contract.UnpackLog(event, "Fill", log); err != nil {
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
func (_WrappersV4 *WrappersV4Filterer) ParseFill(log types.Log) (*WrappersV4Fill, error) {
	event := new(WrappersV4Fill)
	if err := _WrappersV4.contract.UnpackLog(event, "Fill", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WrappersV4OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the WrappersV4 contract.
type WrappersV4OwnershipTransferredIterator struct {
	Event *WrappersV4OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *WrappersV4OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappersV4OwnershipTransferred)
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
		it.Event = new(WrappersV4OwnershipTransferred)
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
func (it *WrappersV4OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappersV4OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappersV4OwnershipTransferred represents a OwnershipTransferred event raised by the WrappersV4 contract.
type WrappersV4OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WrappersV4 *WrappersV4Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*WrappersV4OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _WrappersV4.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &WrappersV4OwnershipTransferredIterator{contract: _WrappersV4.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_WrappersV4 *WrappersV4Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *WrappersV4OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _WrappersV4.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappersV4OwnershipTransferred)
				if err := _WrappersV4.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_WrappersV4 *WrappersV4Filterer) ParseOwnershipTransferred(log types.Log) (*WrappersV4OwnershipTransferred, error) {
	event := new(WrappersV4OwnershipTransferred)
	if err := _WrappersV4.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WrappersV4ProtocolFeeCollectorAddressIterator is returned from FilterProtocolFeeCollectorAddress and is used to iterate over the raw logs and unpacked data for ProtocolFeeCollectorAddress events raised by the WrappersV4 contract.
type WrappersV4ProtocolFeeCollectorAddressIterator struct {
	Event *WrappersV4ProtocolFeeCollectorAddress // Event containing the contract specifics and raw log

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
func (it *WrappersV4ProtocolFeeCollectorAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappersV4ProtocolFeeCollectorAddress)
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
		it.Event = new(WrappersV4ProtocolFeeCollectorAddress)
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
func (it *WrappersV4ProtocolFeeCollectorAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappersV4ProtocolFeeCollectorAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappersV4ProtocolFeeCollectorAddress represents a ProtocolFeeCollectorAddress event raised by the WrappersV4 contract.
type WrappersV4ProtocolFeeCollectorAddress struct {
	OldProtocolFeeCollector     common.Address
	UpdatedProtocolFeeCollector common.Address
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterProtocolFeeCollectorAddress is a free log retrieval operation binding the contract event 0xe1a5430ebec577336427f40f15822f1f36c5e3509ff209d6db9e6c9e6941cb0b.
//
// Solidity: event ProtocolFeeCollectorAddress(address oldProtocolFeeCollector, address updatedProtocolFeeCollector)
func (_WrappersV4 *WrappersV4Filterer) FilterProtocolFeeCollectorAddress(opts *bind.FilterOpts) (*WrappersV4ProtocolFeeCollectorAddressIterator, error) {

	logs, sub, err := _WrappersV4.contract.FilterLogs(opts, "ProtocolFeeCollectorAddress")
	if err != nil {
		return nil, err
	}
	return &WrappersV4ProtocolFeeCollectorAddressIterator{contract: _WrappersV4.contract, event: "ProtocolFeeCollectorAddress", logs: logs, sub: sub}, nil
}

// WatchProtocolFeeCollectorAddress is a free log subscription operation binding the contract event 0xe1a5430ebec577336427f40f15822f1f36c5e3509ff209d6db9e6c9e6941cb0b.
//
// Solidity: event ProtocolFeeCollectorAddress(address oldProtocolFeeCollector, address updatedProtocolFeeCollector)
func (_WrappersV4 *WrappersV4Filterer) WatchProtocolFeeCollectorAddress(opts *bind.WatchOpts, sink chan<- *WrappersV4ProtocolFeeCollectorAddress) (event.Subscription, error) {

	logs, sub, err := _WrappersV4.contract.WatchLogs(opts, "ProtocolFeeCollectorAddress")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappersV4ProtocolFeeCollectorAddress)
				if err := _WrappersV4.contract.UnpackLog(event, "ProtocolFeeCollectorAddress", log); err != nil {
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
func (_WrappersV4 *WrappersV4Filterer) ParseProtocolFeeCollectorAddress(log types.Log) (*WrappersV4ProtocolFeeCollectorAddress, error) {
	event := new(WrappersV4ProtocolFeeCollectorAddress)
	if err := _WrappersV4.contract.UnpackLog(event, "ProtocolFeeCollectorAddress", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WrappersV4ProtocolFeeMultiplierIterator is returned from FilterProtocolFeeMultiplier and is used to iterate over the raw logs and unpacked data for ProtocolFeeMultiplier events raised by the WrappersV4 contract.
type WrappersV4ProtocolFeeMultiplierIterator struct {
	Event *WrappersV4ProtocolFeeMultiplier // Event containing the contract specifics and raw log

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
func (it *WrappersV4ProtocolFeeMultiplierIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappersV4ProtocolFeeMultiplier)
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
		it.Event = new(WrappersV4ProtocolFeeMultiplier)
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
func (it *WrappersV4ProtocolFeeMultiplierIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappersV4ProtocolFeeMultiplierIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappersV4ProtocolFeeMultiplier represents a ProtocolFeeMultiplier event raised by the WrappersV4 contract.
type WrappersV4ProtocolFeeMultiplier struct {
	OldProtocolFeeMultiplier     *big.Int
	UpdatedProtocolFeeMultiplier *big.Int
	Raw                          types.Log // Blockchain specific contextual infos
}

// FilterProtocolFeeMultiplier is a free log retrieval operation binding the contract event 0x3a3e76d7a75e198aef1f53137e4f2a8a2ec74e2e9526db8404d08ccc9f1e621d.
//
// Solidity: event ProtocolFeeMultiplier(uint256 oldProtocolFeeMultiplier, uint256 updatedProtocolFeeMultiplier)
func (_WrappersV4 *WrappersV4Filterer) FilterProtocolFeeMultiplier(opts *bind.FilterOpts) (*WrappersV4ProtocolFeeMultiplierIterator, error) {

	logs, sub, err := _WrappersV4.contract.FilterLogs(opts, "ProtocolFeeMultiplier")
	if err != nil {
		return nil, err
	}
	return &WrappersV4ProtocolFeeMultiplierIterator{contract: _WrappersV4.contract, event: "ProtocolFeeMultiplier", logs: logs, sub: sub}, nil
}

// WatchProtocolFeeMultiplier is a free log subscription operation binding the contract event 0x3a3e76d7a75e198aef1f53137e4f2a8a2ec74e2e9526db8404d08ccc9f1e621d.
//
// Solidity: event ProtocolFeeMultiplier(uint256 oldProtocolFeeMultiplier, uint256 updatedProtocolFeeMultiplier)
func (_WrappersV4 *WrappersV4Filterer) WatchProtocolFeeMultiplier(opts *bind.WatchOpts, sink chan<- *WrappersV4ProtocolFeeMultiplier) (event.Subscription, error) {

	logs, sub, err := _WrappersV4.contract.WatchLogs(opts, "ProtocolFeeMultiplier")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappersV4ProtocolFeeMultiplier)
				if err := _WrappersV4.contract.UnpackLog(event, "ProtocolFeeMultiplier", log); err != nil {
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
func (_WrappersV4 *WrappersV4Filterer) ParseProtocolFeeMultiplier(log types.Log) (*WrappersV4ProtocolFeeMultiplier, error) {
	event := new(WrappersV4ProtocolFeeMultiplier)
	if err := _WrappersV4.contract.UnpackLog(event, "ProtocolFeeMultiplier", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WrappersV4SignatureValidatorApprovalIterator is returned from FilterSignatureValidatorApproval and is used to iterate over the raw logs and unpacked data for SignatureValidatorApproval events raised by the WrappersV4 contract.
type WrappersV4SignatureValidatorApprovalIterator struct {
	Event *WrappersV4SignatureValidatorApproval // Event containing the contract specifics and raw log

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
func (it *WrappersV4SignatureValidatorApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappersV4SignatureValidatorApproval)
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
		it.Event = new(WrappersV4SignatureValidatorApproval)
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
func (it *WrappersV4SignatureValidatorApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappersV4SignatureValidatorApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappersV4SignatureValidatorApproval represents a SignatureValidatorApproval event raised by the WrappersV4 contract.
type WrappersV4SignatureValidatorApproval struct {
	SignerAddress    common.Address
	ValidatorAddress common.Address
	IsApproved       bool
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSignatureValidatorApproval is a free log retrieval operation binding the contract event 0xa8656e308026eeabce8f0bc18048433252318ab80ac79da0b3d3d8697dfba891.
//
// Solidity: event SignatureValidatorApproval(address indexed signerAddress, address indexed validatorAddress, bool isApproved)
func (_WrappersV4 *WrappersV4Filterer) FilterSignatureValidatorApproval(opts *bind.FilterOpts, signerAddress []common.Address, validatorAddress []common.Address) (*WrappersV4SignatureValidatorApprovalIterator, error) {

	var signerAddressRule []interface{}
	for _, signerAddressItem := range signerAddress {
		signerAddressRule = append(signerAddressRule, signerAddressItem)
	}
	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _WrappersV4.contract.FilterLogs(opts, "SignatureValidatorApproval", signerAddressRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return &WrappersV4SignatureValidatorApprovalIterator{contract: _WrappersV4.contract, event: "SignatureValidatorApproval", logs: logs, sub: sub}, nil
}

// WatchSignatureValidatorApproval is a free log subscription operation binding the contract event 0xa8656e308026eeabce8f0bc18048433252318ab80ac79da0b3d3d8697dfba891.
//
// Solidity: event SignatureValidatorApproval(address indexed signerAddress, address indexed validatorAddress, bool isApproved)
func (_WrappersV4 *WrappersV4Filterer) WatchSignatureValidatorApproval(opts *bind.WatchOpts, sink chan<- *WrappersV4SignatureValidatorApproval, signerAddress []common.Address, validatorAddress []common.Address) (event.Subscription, error) {

	var signerAddressRule []interface{}
	for _, signerAddressItem := range signerAddress {
		signerAddressRule = append(signerAddressRule, signerAddressItem)
	}
	var validatorAddressRule []interface{}
	for _, validatorAddressItem := range validatorAddress {
		validatorAddressRule = append(validatorAddressRule, validatorAddressItem)
	}

	logs, sub, err := _WrappersV4.contract.WatchLogs(opts, "SignatureValidatorApproval", signerAddressRule, validatorAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappersV4SignatureValidatorApproval)
				if err := _WrappersV4.contract.UnpackLog(event, "SignatureValidatorApproval", log); err != nil {
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
func (_WrappersV4 *WrappersV4Filterer) ParseSignatureValidatorApproval(log types.Log) (*WrappersV4SignatureValidatorApproval, error) {
	event := new(WrappersV4SignatureValidatorApproval)
	if err := _WrappersV4.contract.UnpackLog(event, "SignatureValidatorApproval", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WrappersV4TransactionExecutionIterator is returned from FilterTransactionExecution and is used to iterate over the raw logs and unpacked data for TransactionExecution events raised by the WrappersV4 contract.
type WrappersV4TransactionExecutionIterator struct {
	Event *WrappersV4TransactionExecution // Event containing the contract specifics and raw log

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
func (it *WrappersV4TransactionExecutionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WrappersV4TransactionExecution)
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
		it.Event = new(WrappersV4TransactionExecution)
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
func (it *WrappersV4TransactionExecutionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WrappersV4TransactionExecutionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WrappersV4TransactionExecution represents a TransactionExecution event raised by the WrappersV4 contract.
type WrappersV4TransactionExecution struct {
	TransactionHash [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterTransactionExecution is a free log retrieval operation binding the contract event 0xa4a7329f1dd821363067e07d359e347b4af9b1efe4b6cccf13240228af3c800d.
//
// Solidity: event TransactionExecution(bytes32 indexed transactionHash)
func (_WrappersV4 *WrappersV4Filterer) FilterTransactionExecution(opts *bind.FilterOpts, transactionHash [][32]byte) (*WrappersV4TransactionExecutionIterator, error) {

	var transactionHashRule []interface{}
	for _, transactionHashItem := range transactionHash {
		transactionHashRule = append(transactionHashRule, transactionHashItem)
	}

	logs, sub, err := _WrappersV4.contract.FilterLogs(opts, "TransactionExecution", transactionHashRule)
	if err != nil {
		return nil, err
	}
	return &WrappersV4TransactionExecutionIterator{contract: _WrappersV4.contract, event: "TransactionExecution", logs: logs, sub: sub}, nil
}

// WatchTransactionExecution is a free log subscription operation binding the contract event 0xa4a7329f1dd821363067e07d359e347b4af9b1efe4b6cccf13240228af3c800d.
//
// Solidity: event TransactionExecution(bytes32 indexed transactionHash)
func (_WrappersV4 *WrappersV4Filterer) WatchTransactionExecution(opts *bind.WatchOpts, sink chan<- *WrappersV4TransactionExecution, transactionHash [][32]byte) (event.Subscription, error) {

	var transactionHashRule []interface{}
	for _, transactionHashItem := range transactionHash {
		transactionHashRule = append(transactionHashRule, transactionHashItem)
	}

	logs, sub, err := _WrappersV4.contract.WatchLogs(opts, "TransactionExecution", transactionHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WrappersV4TransactionExecution)
				if err := _WrappersV4.contract.UnpackLog(event, "TransactionExecution", log); err != nil {
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
func (_WrappersV4 *WrappersV4Filterer) ParseTransactionExecution(log types.Log) (*WrappersV4TransactionExecution, error) {
	event := new(WrappersV4TransactionExecution)
	if err := _WrappersV4.contract.UnpackLog(event, "TransactionExecution", log); err != nil {
		return nil, err
	}
	return event, nil
}
