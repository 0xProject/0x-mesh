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

// IMetaTransactionsFeatureMetaTransactionData is an auto generated low-level Go binding around an user-defined struct.
type IMetaTransactionsFeatureMetaTransactionData struct {
	Signer                common.Address
	Sender                common.Address
	MinGasPrice           *big.Int
	MaxGasPrice           *big.Int
	ExpirationTimeSeconds *big.Int
	Salt                  *big.Int
	CallData              []byte
	Value                 *big.Int
	FeeToken              common.Address
	FeeAmount             *big.Int
}

// ITransformERC20FeatureTransformation is an auto generated low-level Go binding around an user-defined struct.
type ITransformERC20FeatureTransformation struct {
	DeploymentNonce uint32
	Data            []byte
}

// LibNativeOrderLimitOrder is an auto generated low-level Go binding around an user-defined struct.
type LibNativeOrderLimitOrder struct {
	MakerToken          common.Address
	TakerToken          common.Address
	MakerAmount         *big.Int
	TakerAmount         *big.Int
	TakerTokenFeeAmount *big.Int
	Maker               common.Address
	Taker               common.Address
	Sender              common.Address
	FeeRecipient        common.Address
	Pool                [32]byte
	Expiry              uint64
	Salt                *big.Int
}

// LibNativeOrderOrderInfo is an auto generated low-level Go binding around an user-defined struct.
type LibNativeOrderOrderInfo struct {
	OrderHash              [32]byte
	Status                 uint8
	TakerTokenFilledAmount *big.Int
}

// LibNativeOrderRfqOrder is an auto generated low-level Go binding around an user-defined struct.
type LibNativeOrderRfqOrder struct {
	MakerToken  common.Address
	TakerToken  common.Address
	MakerAmount *big.Int
	TakerAmount *big.Int
	Maker       common.Address
	Taker       common.Address
	TxOrigin    common.Address
	Pool        [32]byte
	Expiry      uint64
	Salt        *big.Int
}

// LibSignatureSignature is an auto generated low-level Go binding around an user-defined struct.
type LibSignatureSignature struct {
	SignatureType uint8
	V             uint8
	R             [32]byte
	S             [32]byte
}

// ExchangeV4ABI is the input ABI used to generate the binding from.
const ExchangeV4ABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"makerToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"takerToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"takerTokenFilledAmount\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"makerTokenFilledAmount\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"takerTokenFeeFilledAmount\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"protocolFeePaid\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"}],\"name\":\"LimitOrderFilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"MetaTransactionExecuted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"migrator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"Migrated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"}],\"name\":\"OrderCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"makerToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"takerToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minValidSalt\",\"type\":\"uint256\"}],\"name\":\"PairCancelledLimitOrders\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"makerToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"takerToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minValidSalt\",\"type\":\"uint256\"}],\"name\":\"PairCancelledRfqOrders\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldImpl\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newImpl\",\"type\":\"address\"}],\"name\":\"ProxyFunctionUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"quoteSigner\",\"type\":\"address\"}],\"name\":\"QuoteSignerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"makerToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"takerToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"takerTokenFilledAmount\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"makerTokenFilledAmount\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"}],\"name\":\"RfqOrderFilled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"origin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"addrs\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"RfqOrderOriginsAllowed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"inputToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"outputToken\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"inputTokenAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outputTokenAmount\",\"type\":\"uint256\"}],\"name\":\"TransformedERC20\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"transformerDeployer\",\"type\":\"address\"}],\"name\":\"TransformerDeployerUpdated\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"makerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFeeAmount\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"internalType\":\"structLibNativeOrder.LimitOrder[]\",\"name\":\"orders\",\"type\":\"tuple[]\"}],\"name\":\"batchCancelLimitOrders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20TokenV06[]\",\"name\":\"makerTokens\",\"type\":\"address[]\"},{\"internalType\":\"contractIERC20TokenV06[]\",\"name\":\"takerTokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"minValidSalts\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelPairLimitOrders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20TokenV06[]\",\"name\":\"makerTokens\",\"type\":\"address[]\"},{\"internalType\":\"contractIERC20TokenV06[]\",\"name\":\"takerTokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"minValidSalts\",\"type\":\"uint256[]\"}],\"name\":\"batchCancelPairRfqOrders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"makerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerAmount\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"txOrigin\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"internalType\":\"structLibNativeOrder.RfqOrder[]\",\"name\":\"orders\",\"type\":\"tuple[]\"}],\"name\":\"batchCancelRfqOrders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeAmount\",\"type\":\"uint256\"}],\"internalType\":\"structIMetaTransactionsFeature.MetaTransactionData[]\",\"name\":\"mtxs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"enumLibSignature.SignatureType\",\"name\":\"signatureType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structLibSignature.Signature[]\",\"name\":\"signatures\",\"type\":\"tuple[]\"}],\"name\":\"batchExecuteMetaTransactions\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"returnResults\",\"type\":\"bytes[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"makerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFeeAmount\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"internalType\":\"structLibNativeOrder.LimitOrder[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"enumLibSignature.SignatureType\",\"name\":\"signatureType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structLibSignature.Signature[]\",\"name\":\"signatures\",\"type\":\"tuple[]\"}],\"name\":\"batchGetLimitOrderRelevantStates\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumLibNativeOrder.OrderStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFilledAmount\",\"type\":\"uint128\"}],\"internalType\":\"structLibNativeOrder.OrderInfo[]\",\"name\":\"orderInfos\",\"type\":\"tuple[]\"},{\"internalType\":\"uint128[]\",\"name\":\"actualFillableTakerTokenAmounts\",\"type\":\"uint128[]\"},{\"internalType\":\"bool[]\",\"name\":\"isSignatureValids\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"makerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerAmount\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"txOrigin\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"internalType\":\"structLibNativeOrder.RfqOrder[]\",\"name\":\"orders\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"enumLibSignature.SignatureType\",\"name\":\"signatureType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structLibSignature.Signature[]\",\"name\":\"signatures\",\"type\":\"tuple[]\"}],\"name\":\"batchGetRfqOrderRelevantStates\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumLibNativeOrder.OrderStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFilledAmount\",\"type\":\"uint128\"}],\"internalType\":\"structLibNativeOrder.OrderInfo[]\",\"name\":\"orderInfos\",\"type\":\"tuple[]\"},{\"internalType\":\"uint128[]\",\"name\":\"actualFillableTakerTokenAmounts\",\"type\":\"uint128[]\"},{\"internalType\":\"bool[]\",\"name\":\"isSignatureValids\",\"type\":\"bool[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"makerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFeeAmount\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"internalType\":\"structLibNativeOrder.LimitOrder\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"cancelLimitOrder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minValidSalt\",\"type\":\"uint256\"}],\"name\":\"cancelPairLimitOrders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minValidSalt\",\"type\":\"uint256\"}],\"name\":\"cancelPairRfqOrders\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"makerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerAmount\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"txOrigin\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"internalType\":\"structLibNativeOrder.RfqOrder\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"cancelRfqOrder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createTransformWallet\",\"outputs\":[{\"internalType\":\"contractIFlashWallet\",\"name\":\"wallet\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeAmount\",\"type\":\"uint256\"}],\"internalType\":\"structIMetaTransactionsFeature.MetaTransactionData\",\"name\":\"mtx\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"enumLibSignature.SignatureType\",\"name\":\"signatureType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structLibSignature.Signature\",\"name\":\"signature\",\"type\":\"tuple\"}],\"name\":\"executeMetaTransaction\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"returnResult\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"impl\",\"type\":\"address\"}],\"name\":\"extend\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"makerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFeeAmount\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"internalType\":\"structLibNativeOrder.LimitOrder\",\"name\":\"order\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"enumLibSignature.SignatureType\",\"name\":\"signatureType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structLibSignature.Signature\",\"name\":\"signature\",\"type\":\"tuple\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFillAmount\",\"type\":\"uint128\"}],\"name\":\"fillLimitOrder\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"takerTokenFilledAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"makerTokenFilledAmount\",\"type\":\"uint128\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"makerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFeeAmount\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"internalType\":\"structLibNativeOrder.LimitOrder\",\"name\":\"order\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"enumLibSignature.SignatureType\",\"name\":\"signatureType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structLibSignature.Signature\",\"name\":\"signature\",\"type\":\"tuple\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFillAmount\",\"type\":\"uint128\"}],\"name\":\"fillOrKillLimitOrder\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"makerTokenFilledAmount\",\"type\":\"uint128\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"makerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerAmount\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"txOrigin\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"internalType\":\"structLibNativeOrder.RfqOrder\",\"name\":\"order\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"enumLibSignature.SignatureType\",\"name\":\"signatureType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structLibSignature.Signature\",\"name\":\"signature\",\"type\":\"tuple\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFillAmount\",\"type\":\"uint128\"}],\"name\":\"fillOrKillRfqOrder\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"makerTokenFilledAmount\",\"type\":\"uint128\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"makerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerAmount\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"txOrigin\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"internalType\":\"structLibNativeOrder.RfqOrder\",\"name\":\"order\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"enumLibSignature.SignatureType\",\"name\":\"signatureType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structLibSignature.Signature\",\"name\":\"signature\",\"type\":\"tuple\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFillAmount\",\"type\":\"uint128\"}],\"name\":\"fillRfqOrder\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"takerTokenFilledAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"makerTokenFilledAmount\",\"type\":\"uint128\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllowanceTarget\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"makerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFeeAmount\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"internalType\":\"structLibNativeOrder.LimitOrder\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"getLimitOrderHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"makerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFeeAmount\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"internalType\":\"structLibNativeOrder.LimitOrder\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"getLimitOrderInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumLibNativeOrder.OrderStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFilledAmount\",\"type\":\"uint128\"}],\"internalType\":\"structLibNativeOrder.OrderInfo\",\"name\":\"orderInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"makerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFeeAmount\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"internalType\":\"structLibNativeOrder.LimitOrder\",\"name\":\"order\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"enumLibSignature.SignatureType\",\"name\":\"signatureType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structLibSignature.Signature\",\"name\":\"signature\",\"type\":\"tuple\"}],\"name\":\"getLimitOrderRelevantState\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumLibNativeOrder.OrderStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFilledAmount\",\"type\":\"uint128\"}],\"internalType\":\"structLibNativeOrder.OrderInfo\",\"name\":\"orderInfo\",\"type\":\"tuple\"},{\"internalType\":\"uint128\",\"name\":\"actualFillableTakerTokenAmount\",\"type\":\"uint128\"},{\"internalType\":\"bool\",\"name\":\"isSignatureValid\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeAmount\",\"type\":\"uint256\"}],\"internalType\":\"structIMetaTransactionsFeature.MetaTransactionData\",\"name\":\"mtx\",\"type\":\"tuple\"}],\"name\":\"getMetaTransactionExecutedBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"addresspayable\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"minGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxGasPrice\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"feeToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"feeAmount\",\"type\":\"uint256\"}],\"internalType\":\"structIMetaTransactionsFeature.MetaTransactionData\",\"name\":\"mtx\",\"type\":\"tuple\"}],\"name\":\"getMetaTransactionHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"mtxHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"mtxHash\",\"type\":\"bytes32\"}],\"name\":\"getMetaTransactionHashExecutedBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProtocolFeeMultiplier\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"multiplier\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getQuoteSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"makerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerAmount\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"txOrigin\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"internalType\":\"structLibNativeOrder.RfqOrder\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"getRfqOrderHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"makerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerAmount\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"txOrigin\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"internalType\":\"structLibNativeOrder.RfqOrder\",\"name\":\"order\",\"type\":\"tuple\"}],\"name\":\"getRfqOrderInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumLibNativeOrder.OrderStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFilledAmount\",\"type\":\"uint128\"}],\"internalType\":\"structLibNativeOrder.OrderInfo\",\"name\":\"orderInfo\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"makerToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"takerToken\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"makerAmount\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"takerAmount\",\"type\":\"uint128\"},{\"internalType\":\"address\",\"name\":\"maker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"taker\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"txOrigin\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"pool\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"expiry\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"salt\",\"type\":\"uint256\"}],\"internalType\":\"structLibNativeOrder.RfqOrder\",\"name\":\"order\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"enumLibSignature.SignatureType\",\"name\":\"signatureType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structLibSignature.Signature\",\"name\":\"signature\",\"type\":\"tuple\"}],\"name\":\"getRfqOrderRelevantState\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumLibNativeOrder.OrderStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint128\",\"name\":\"takerTokenFilledAmount\",\"type\":\"uint128\"}],\"internalType\":\"structLibNativeOrder.OrderInfo\",\"name\":\"orderInfo\",\"type\":\"tuple\"},{\"internalType\":\"uint128\",\"name\":\"actualFillableTakerTokenAmount\",\"type\":\"uint128\"},{\"internalType\":\"bool\",\"name\":\"isSignatureValid\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"},{\"internalType\":\"uint256\",\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"getRollbackEntryAtIndex\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"impl\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"}],\"name\":\"getRollbackLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"rollbackLength\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"getSpendableERC20BalanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransformWallet\",\"outputs\":[{\"internalType\":\"contractIFlashWallet\",\"name\":\"wallet\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTransformerDeployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"isValidHashSignature\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isValid\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"migrate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"ownerAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"origins\",\"type\":\"address[]\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"registerAllowedRfqOrigins\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"selector\",\"type\":\"bytes4\"},{\"internalType\":\"address\",\"name\":\"targetImpl\",\"type\":\"address\"}],\"name\":\"rollback\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"inputToken\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"outputToken\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"provider\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"sellAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBuyAmount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"auxiliaryData\",\"type\":\"bytes\"}],\"name\":\"sellToLiquidityProvider\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"boughtAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20TokenV06[]\",\"name\":\"tokens\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"sellAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minBuyAmount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isSushi\",\"type\":\"bool\"}],\"name\":\"sellToUniswap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"buyAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"quoteSigner\",\"type\":\"address\"}],\"name\":\"setQuoteSigner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"transformerDeployer\",\"type\":\"address\"}],\"name\":\"setTransformerDeployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"poolIds\",\"type\":\"bytes32[]\"}],\"name\":\"transferProtocolFeesForPools\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"inputToken\",\"type\":\"address\"},{\"internalType\":\"contractIERC20TokenV06\",\"name\":\"outputToken\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"inputTokenAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minOutputTokenAmount\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"deploymentNonce\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structITransformERC20Feature.Transformation[]\",\"name\":\"transformations\",\"type\":\"tuple[]\"}],\"name\":\"transformERC20\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"outputTokenAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"validateHashSignature\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// ExchangeV4 is an auto generated Go binding around an Ethereum contract.
type ExchangeV4 struct {
	ExchangeV4Caller     // Read-only binding to the contract
	ExchangeV4Transactor // Write-only binding to the contract
	ExchangeV4Filterer   // Log filterer for contract events
}

// ExchangeV4Caller is an auto generated read-only Go binding around an Ethereum contract.
type ExchangeV4Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeV4Transactor is an auto generated write-only Go binding around an Ethereum contract.
type ExchangeV4Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeV4Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ExchangeV4Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExchangeV4Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExchangeV4Session struct {
	Contract     *ExchangeV4       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ExchangeV4CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExchangeV4CallerSession struct {
	Contract *ExchangeV4Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ExchangeV4TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExchangeV4TransactorSession struct {
	Contract     *ExchangeV4Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ExchangeV4Raw is an auto generated low-level Go binding around an Ethereum contract.
type ExchangeV4Raw struct {
	Contract *ExchangeV4 // Generic contract binding to access the raw methods on
}

// ExchangeV4CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExchangeV4CallerRaw struct {
	Contract *ExchangeV4Caller // Generic read-only contract binding to access the raw methods on
}

// ExchangeV4TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExchangeV4TransactorRaw struct {
	Contract *ExchangeV4Transactor // Generic write-only contract binding to access the raw methods on
}

// NewExchangeV4 creates a new instance of ExchangeV4, bound to a specific deployed contract.
func NewExchangeV4(address common.Address, backend bind.ContractBackend) (*ExchangeV4, error) {
	contract, err := bindExchangeV4(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ExchangeV4{ExchangeV4Caller: ExchangeV4Caller{contract: contract}, ExchangeV4Transactor: ExchangeV4Transactor{contract: contract}, ExchangeV4Filterer: ExchangeV4Filterer{contract: contract}}, nil
}

// NewExchangeV4Caller creates a new read-only instance of ExchangeV4, bound to a specific deployed contract.
func NewExchangeV4Caller(address common.Address, caller bind.ContractCaller) (*ExchangeV4Caller, error) {
	contract, err := bindExchangeV4(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangeV4Caller{contract: contract}, nil
}

// NewExchangeV4Transactor creates a new write-only instance of ExchangeV4, bound to a specific deployed contract.
func NewExchangeV4Transactor(address common.Address, transactor bind.ContractTransactor) (*ExchangeV4Transactor, error) {
	contract, err := bindExchangeV4(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExchangeV4Transactor{contract: contract}, nil
}

// NewExchangeV4Filterer creates a new log filterer instance of ExchangeV4, bound to a specific deployed contract.
func NewExchangeV4Filterer(address common.Address, filterer bind.ContractFilterer) (*ExchangeV4Filterer, error) {
	contract, err := bindExchangeV4(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExchangeV4Filterer{contract: contract}, nil
}

// bindExchangeV4 binds a generic wrapper to an already deployed contract.
func bindExchangeV4(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ExchangeV4ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExchangeV4 *ExchangeV4Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ExchangeV4.Contract.ExchangeV4Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExchangeV4 *ExchangeV4Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExchangeV4.Contract.ExchangeV4Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExchangeV4 *ExchangeV4Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExchangeV4.Contract.ExchangeV4Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExchangeV4 *ExchangeV4CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ExchangeV4.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExchangeV4 *ExchangeV4TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExchangeV4.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExchangeV4 *ExchangeV4TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExchangeV4.Contract.contract.Transact(opts, method, params...)
}

// BatchGetLimitOrderRelevantStates is a free data retrieval call binding the contract method 0xb4658bfb.
//
// Solidity: function batchGetLimitOrderRelevantStates((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256)[] orders, (uint8,uint8,bytes32,bytes32)[] signatures) view returns((bytes32,uint8,uint128)[] orderInfos, uint128[] actualFillableTakerTokenAmounts, bool[] isSignatureValids)
func (_ExchangeV4 *ExchangeV4Caller) BatchGetLimitOrderRelevantStates(opts *bind.CallOpts, orders []LibNativeOrderLimitOrder, signatures []LibSignatureSignature) (struct {
	OrderInfos                      []LibNativeOrderOrderInfo
	ActualFillableTakerTokenAmounts []*big.Int
	IsSignatureValids               []bool
}, error) {
	ret := new(struct {
		OrderInfos                      []LibNativeOrderOrderInfo
		ActualFillableTakerTokenAmounts []*big.Int
		IsSignatureValids               []bool
	})
	out := ret
	err := _ExchangeV4.contract.Call(opts, out, "batchGetLimitOrderRelevantStates", orders, signatures)
	return *ret, err
}

// BatchGetLimitOrderRelevantStates is a free data retrieval call binding the contract method 0xb4658bfb.
//
// Solidity: function batchGetLimitOrderRelevantStates((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256)[] orders, (uint8,uint8,bytes32,bytes32)[] signatures) view returns((bytes32,uint8,uint128)[] orderInfos, uint128[] actualFillableTakerTokenAmounts, bool[] isSignatureValids)
func (_ExchangeV4 *ExchangeV4Session) BatchGetLimitOrderRelevantStates(orders []LibNativeOrderLimitOrder, signatures []LibSignatureSignature) (struct {
	OrderInfos                      []LibNativeOrderOrderInfo
	ActualFillableTakerTokenAmounts []*big.Int
	IsSignatureValids               []bool
}, error) {
	return _ExchangeV4.Contract.BatchGetLimitOrderRelevantStates(&_ExchangeV4.CallOpts, orders, signatures)
}

// BatchGetLimitOrderRelevantStates is a free data retrieval call binding the contract method 0xb4658bfb.
//
// Solidity: function batchGetLimitOrderRelevantStates((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256)[] orders, (uint8,uint8,bytes32,bytes32)[] signatures) view returns((bytes32,uint8,uint128)[] orderInfos, uint128[] actualFillableTakerTokenAmounts, bool[] isSignatureValids)
func (_ExchangeV4 *ExchangeV4CallerSession) BatchGetLimitOrderRelevantStates(orders []LibNativeOrderLimitOrder, signatures []LibSignatureSignature) (struct {
	OrderInfos                      []LibNativeOrderOrderInfo
	ActualFillableTakerTokenAmounts []*big.Int
	IsSignatureValids               []bool
}, error) {
	return _ExchangeV4.Contract.BatchGetLimitOrderRelevantStates(&_ExchangeV4.CallOpts, orders, signatures)
}

// BatchGetRfqOrderRelevantStates is a free data retrieval call binding the contract method 0xad354eeb.
//
// Solidity: function batchGetRfqOrderRelevantStates((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256)[] orders, (uint8,uint8,bytes32,bytes32)[] signatures) view returns((bytes32,uint8,uint128)[] orderInfos, uint128[] actualFillableTakerTokenAmounts, bool[] isSignatureValids)
func (_ExchangeV4 *ExchangeV4Caller) BatchGetRfqOrderRelevantStates(opts *bind.CallOpts, orders []LibNativeOrderRfqOrder, signatures []LibSignatureSignature) (struct {
	OrderInfos                      []LibNativeOrderOrderInfo
	ActualFillableTakerTokenAmounts []*big.Int
	IsSignatureValids               []bool
}, error) {
	ret := new(struct {
		OrderInfos                      []LibNativeOrderOrderInfo
		ActualFillableTakerTokenAmounts []*big.Int
		IsSignatureValids               []bool
	})
	out := ret
	err := _ExchangeV4.contract.Call(opts, out, "batchGetRfqOrderRelevantStates", orders, signatures)
	return *ret, err
}

// BatchGetRfqOrderRelevantStates is a free data retrieval call binding the contract method 0xad354eeb.
//
// Solidity: function batchGetRfqOrderRelevantStates((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256)[] orders, (uint8,uint8,bytes32,bytes32)[] signatures) view returns((bytes32,uint8,uint128)[] orderInfos, uint128[] actualFillableTakerTokenAmounts, bool[] isSignatureValids)
func (_ExchangeV4 *ExchangeV4Session) BatchGetRfqOrderRelevantStates(orders []LibNativeOrderRfqOrder, signatures []LibSignatureSignature) (struct {
	OrderInfos                      []LibNativeOrderOrderInfo
	ActualFillableTakerTokenAmounts []*big.Int
	IsSignatureValids               []bool
}, error) {
	return _ExchangeV4.Contract.BatchGetRfqOrderRelevantStates(&_ExchangeV4.CallOpts, orders, signatures)
}

// BatchGetRfqOrderRelevantStates is a free data retrieval call binding the contract method 0xad354eeb.
//
// Solidity: function batchGetRfqOrderRelevantStates((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256)[] orders, (uint8,uint8,bytes32,bytes32)[] signatures) view returns((bytes32,uint8,uint128)[] orderInfos, uint128[] actualFillableTakerTokenAmounts, bool[] isSignatureValids)
func (_ExchangeV4 *ExchangeV4CallerSession) BatchGetRfqOrderRelevantStates(orders []LibNativeOrderRfqOrder, signatures []LibSignatureSignature) (struct {
	OrderInfos                      []LibNativeOrderOrderInfo
	ActualFillableTakerTokenAmounts []*big.Int
	IsSignatureValids               []bool
}, error) {
	return _ExchangeV4.Contract.BatchGetRfqOrderRelevantStates(&_ExchangeV4.CallOpts, orders, signatures)
}

// GetAllowanceTarget is a free data retrieval call binding the contract method 0xf7c3a33b.
//
// Solidity: function getAllowanceTarget() view returns(address target)
func (_ExchangeV4 *ExchangeV4Caller) GetAllowanceTarget(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "getAllowanceTarget")
	return *ret0, err
}

// GetAllowanceTarget is a free data retrieval call binding the contract method 0xf7c3a33b.
//
// Solidity: function getAllowanceTarget() view returns(address target)
func (_ExchangeV4 *ExchangeV4Session) GetAllowanceTarget() (common.Address, error) {
	return _ExchangeV4.Contract.GetAllowanceTarget(&_ExchangeV4.CallOpts)
}

// GetAllowanceTarget is a free data retrieval call binding the contract method 0xf7c3a33b.
//
// Solidity: function getAllowanceTarget() view returns(address target)
func (_ExchangeV4 *ExchangeV4CallerSession) GetAllowanceTarget() (common.Address, error) {
	return _ExchangeV4.Contract.GetAllowanceTarget(&_ExchangeV4.CallOpts)
}

// GetLimitOrderHash is a free data retrieval call binding the contract method 0xdd11d225.
//
// Solidity: function getLimitOrderHash((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order) view returns(bytes32 orderHash)
func (_ExchangeV4 *ExchangeV4Caller) GetLimitOrderHash(opts *bind.CallOpts, order LibNativeOrderLimitOrder) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "getLimitOrderHash", order)
	return *ret0, err
}

// GetLimitOrderHash is a free data retrieval call binding the contract method 0xdd11d225.
//
// Solidity: function getLimitOrderHash((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order) view returns(bytes32 orderHash)
func (_ExchangeV4 *ExchangeV4Session) GetLimitOrderHash(order LibNativeOrderLimitOrder) ([32]byte, error) {
	return _ExchangeV4.Contract.GetLimitOrderHash(&_ExchangeV4.CallOpts, order)
}

// GetLimitOrderHash is a free data retrieval call binding the contract method 0xdd11d225.
//
// Solidity: function getLimitOrderHash((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order) view returns(bytes32 orderHash)
func (_ExchangeV4 *ExchangeV4CallerSession) GetLimitOrderHash(order LibNativeOrderLimitOrder) ([32]byte, error) {
	return _ExchangeV4.Contract.GetLimitOrderHash(&_ExchangeV4.CallOpts, order)
}

// GetLimitOrderInfo is a free data retrieval call binding the contract method 0x95480889.
//
// Solidity: function getLimitOrderInfo((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order) view returns((bytes32,uint8,uint128) orderInfo)
func (_ExchangeV4 *ExchangeV4Caller) GetLimitOrderInfo(opts *bind.CallOpts, order LibNativeOrderLimitOrder) (LibNativeOrderOrderInfo, error) {
	var (
		ret0 = new(LibNativeOrderOrderInfo)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "getLimitOrderInfo", order)
	return *ret0, err
}

// GetLimitOrderInfo is a free data retrieval call binding the contract method 0x95480889.
//
// Solidity: function getLimitOrderInfo((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order) view returns((bytes32,uint8,uint128) orderInfo)
func (_ExchangeV4 *ExchangeV4Session) GetLimitOrderInfo(order LibNativeOrderLimitOrder) (LibNativeOrderOrderInfo, error) {
	return _ExchangeV4.Contract.GetLimitOrderInfo(&_ExchangeV4.CallOpts, order)
}

// GetLimitOrderInfo is a free data retrieval call binding the contract method 0x95480889.
//
// Solidity: function getLimitOrderInfo((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order) view returns((bytes32,uint8,uint128) orderInfo)
func (_ExchangeV4 *ExchangeV4CallerSession) GetLimitOrderInfo(order LibNativeOrderLimitOrder) (LibNativeOrderOrderInfo, error) {
	return _ExchangeV4.Contract.GetLimitOrderInfo(&_ExchangeV4.CallOpts, order)
}

// GetLimitOrderRelevantState is a free data retrieval call binding the contract method 0x1fb09795.
//
// Solidity: function getLimitOrderRelevantState((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature) view returns((bytes32,uint8,uint128) orderInfo, uint128 actualFillableTakerTokenAmount, bool isSignatureValid)
func (_ExchangeV4 *ExchangeV4Caller) GetLimitOrderRelevantState(opts *bind.CallOpts, order LibNativeOrderLimitOrder, signature LibSignatureSignature) (struct {
	OrderInfo                      LibNativeOrderOrderInfo
	ActualFillableTakerTokenAmount *big.Int
	IsSignatureValid               bool
}, error) {
	ret := new(struct {
		OrderInfo                      LibNativeOrderOrderInfo
		ActualFillableTakerTokenAmount *big.Int
		IsSignatureValid               bool
	})
	out := ret
	err := _ExchangeV4.contract.Call(opts, out, "getLimitOrderRelevantState", order, signature)
	return *ret, err
}

// GetLimitOrderRelevantState is a free data retrieval call binding the contract method 0x1fb09795.
//
// Solidity: function getLimitOrderRelevantState((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature) view returns((bytes32,uint8,uint128) orderInfo, uint128 actualFillableTakerTokenAmount, bool isSignatureValid)
func (_ExchangeV4 *ExchangeV4Session) GetLimitOrderRelevantState(order LibNativeOrderLimitOrder, signature LibSignatureSignature) (struct {
	OrderInfo                      LibNativeOrderOrderInfo
	ActualFillableTakerTokenAmount *big.Int
	IsSignatureValid               bool
}, error) {
	return _ExchangeV4.Contract.GetLimitOrderRelevantState(&_ExchangeV4.CallOpts, order, signature)
}

// GetLimitOrderRelevantState is a free data retrieval call binding the contract method 0x1fb09795.
//
// Solidity: function getLimitOrderRelevantState((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature) view returns((bytes32,uint8,uint128) orderInfo, uint128 actualFillableTakerTokenAmount, bool isSignatureValid)
func (_ExchangeV4 *ExchangeV4CallerSession) GetLimitOrderRelevantState(order LibNativeOrderLimitOrder, signature LibSignatureSignature) (struct {
	OrderInfo                      LibNativeOrderOrderInfo
	ActualFillableTakerTokenAmount *big.Int
	IsSignatureValid               bool
}, error) {
	return _ExchangeV4.Contract.GetLimitOrderRelevantState(&_ExchangeV4.CallOpts, order, signature)
}

// GetMetaTransactionExecutedBlock is a free data retrieval call binding the contract method 0x3fb2da38.
//
// Solidity: function getMetaTransactionExecutedBlock((address,address,uint256,uint256,uint256,uint256,bytes,uint256,address,uint256) mtx) view returns(uint256 blockNumber)
func (_ExchangeV4 *ExchangeV4Caller) GetMetaTransactionExecutedBlock(opts *bind.CallOpts, mtx IMetaTransactionsFeatureMetaTransactionData) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "getMetaTransactionExecutedBlock", mtx)
	return *ret0, err
}

// GetMetaTransactionExecutedBlock is a free data retrieval call binding the contract method 0x3fb2da38.
//
// Solidity: function getMetaTransactionExecutedBlock((address,address,uint256,uint256,uint256,uint256,bytes,uint256,address,uint256) mtx) view returns(uint256 blockNumber)
func (_ExchangeV4 *ExchangeV4Session) GetMetaTransactionExecutedBlock(mtx IMetaTransactionsFeatureMetaTransactionData) (*big.Int, error) {
	return _ExchangeV4.Contract.GetMetaTransactionExecutedBlock(&_ExchangeV4.CallOpts, mtx)
}

// GetMetaTransactionExecutedBlock is a free data retrieval call binding the contract method 0x3fb2da38.
//
// Solidity: function getMetaTransactionExecutedBlock((address,address,uint256,uint256,uint256,uint256,bytes,uint256,address,uint256) mtx) view returns(uint256 blockNumber)
func (_ExchangeV4 *ExchangeV4CallerSession) GetMetaTransactionExecutedBlock(mtx IMetaTransactionsFeatureMetaTransactionData) (*big.Int, error) {
	return _ExchangeV4.Contract.GetMetaTransactionExecutedBlock(&_ExchangeV4.CallOpts, mtx)
}

// GetMetaTransactionHash is a free data retrieval call binding the contract method 0xae550497.
//
// Solidity: function getMetaTransactionHash((address,address,uint256,uint256,uint256,uint256,bytes,uint256,address,uint256) mtx) view returns(bytes32 mtxHash)
func (_ExchangeV4 *ExchangeV4Caller) GetMetaTransactionHash(opts *bind.CallOpts, mtx IMetaTransactionsFeatureMetaTransactionData) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "getMetaTransactionHash", mtx)
	return *ret0, err
}

// GetMetaTransactionHash is a free data retrieval call binding the contract method 0xae550497.
//
// Solidity: function getMetaTransactionHash((address,address,uint256,uint256,uint256,uint256,bytes,uint256,address,uint256) mtx) view returns(bytes32 mtxHash)
func (_ExchangeV4 *ExchangeV4Session) GetMetaTransactionHash(mtx IMetaTransactionsFeatureMetaTransactionData) ([32]byte, error) {
	return _ExchangeV4.Contract.GetMetaTransactionHash(&_ExchangeV4.CallOpts, mtx)
}

// GetMetaTransactionHash is a free data retrieval call binding the contract method 0xae550497.
//
// Solidity: function getMetaTransactionHash((address,address,uint256,uint256,uint256,uint256,bytes,uint256,address,uint256) mtx) view returns(bytes32 mtxHash)
func (_ExchangeV4 *ExchangeV4CallerSession) GetMetaTransactionHash(mtx IMetaTransactionsFeatureMetaTransactionData) ([32]byte, error) {
	return _ExchangeV4.Contract.GetMetaTransactionHash(&_ExchangeV4.CallOpts, mtx)
}

// GetMetaTransactionHashExecutedBlock is a free data retrieval call binding the contract method 0x72d17d03.
//
// Solidity: function getMetaTransactionHashExecutedBlock(bytes32 mtxHash) view returns(uint256 blockNumber)
func (_ExchangeV4 *ExchangeV4Caller) GetMetaTransactionHashExecutedBlock(opts *bind.CallOpts, mtxHash [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "getMetaTransactionHashExecutedBlock", mtxHash)
	return *ret0, err
}

// GetMetaTransactionHashExecutedBlock is a free data retrieval call binding the contract method 0x72d17d03.
//
// Solidity: function getMetaTransactionHashExecutedBlock(bytes32 mtxHash) view returns(uint256 blockNumber)
func (_ExchangeV4 *ExchangeV4Session) GetMetaTransactionHashExecutedBlock(mtxHash [32]byte) (*big.Int, error) {
	return _ExchangeV4.Contract.GetMetaTransactionHashExecutedBlock(&_ExchangeV4.CallOpts, mtxHash)
}

// GetMetaTransactionHashExecutedBlock is a free data retrieval call binding the contract method 0x72d17d03.
//
// Solidity: function getMetaTransactionHashExecutedBlock(bytes32 mtxHash) view returns(uint256 blockNumber)
func (_ExchangeV4 *ExchangeV4CallerSession) GetMetaTransactionHashExecutedBlock(mtxHash [32]byte) (*big.Int, error) {
	return _ExchangeV4.Contract.GetMetaTransactionHashExecutedBlock(&_ExchangeV4.CallOpts, mtxHash)
}

// GetProtocolFeeMultiplier is a free data retrieval call binding the contract method 0x487b5c20.
//
// Solidity: function getProtocolFeeMultiplier() view returns(uint32 multiplier)
func (_ExchangeV4 *ExchangeV4Caller) GetProtocolFeeMultiplier(opts *bind.CallOpts) (uint32, error) {
	var (
		ret0 = new(uint32)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "getProtocolFeeMultiplier")
	return *ret0, err
}

// GetProtocolFeeMultiplier is a free data retrieval call binding the contract method 0x487b5c20.
//
// Solidity: function getProtocolFeeMultiplier() view returns(uint32 multiplier)
func (_ExchangeV4 *ExchangeV4Session) GetProtocolFeeMultiplier() (uint32, error) {
	return _ExchangeV4.Contract.GetProtocolFeeMultiplier(&_ExchangeV4.CallOpts)
}

// GetProtocolFeeMultiplier is a free data retrieval call binding the contract method 0x487b5c20.
//
// Solidity: function getProtocolFeeMultiplier() view returns(uint32 multiplier)
func (_ExchangeV4 *ExchangeV4CallerSession) GetProtocolFeeMultiplier() (uint32, error) {
	return _ExchangeV4.Contract.GetProtocolFeeMultiplier(&_ExchangeV4.CallOpts)
}

// GetQuoteSigner is a free data retrieval call binding the contract method 0x9f1ec78b.
//
// Solidity: function getQuoteSigner() view returns(address signer)
func (_ExchangeV4 *ExchangeV4Caller) GetQuoteSigner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "getQuoteSigner")
	return *ret0, err
}

// GetQuoteSigner is a free data retrieval call binding the contract method 0x9f1ec78b.
//
// Solidity: function getQuoteSigner() view returns(address signer)
func (_ExchangeV4 *ExchangeV4Session) GetQuoteSigner() (common.Address, error) {
	return _ExchangeV4.Contract.GetQuoteSigner(&_ExchangeV4.CallOpts)
}

// GetQuoteSigner is a free data retrieval call binding the contract method 0x9f1ec78b.
//
// Solidity: function getQuoteSigner() view returns(address signer)
func (_ExchangeV4 *ExchangeV4CallerSession) GetQuoteSigner() (common.Address, error) {
	return _ExchangeV4.Contract.GetQuoteSigner(&_ExchangeV4.CallOpts)
}

// GetRfqOrderHash is a free data retrieval call binding the contract method 0x016a6d65.
//
// Solidity: function getRfqOrderHash((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order) view returns(bytes32 orderHash)
func (_ExchangeV4 *ExchangeV4Caller) GetRfqOrderHash(opts *bind.CallOpts, order LibNativeOrderRfqOrder) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "getRfqOrderHash", order)
	return *ret0, err
}

// GetRfqOrderHash is a free data retrieval call binding the contract method 0x016a6d65.
//
// Solidity: function getRfqOrderHash((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order) view returns(bytes32 orderHash)
func (_ExchangeV4 *ExchangeV4Session) GetRfqOrderHash(order LibNativeOrderRfqOrder) ([32]byte, error) {
	return _ExchangeV4.Contract.GetRfqOrderHash(&_ExchangeV4.CallOpts, order)
}

// GetRfqOrderHash is a free data retrieval call binding the contract method 0x016a6d65.
//
// Solidity: function getRfqOrderHash((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order) view returns(bytes32 orderHash)
func (_ExchangeV4 *ExchangeV4CallerSession) GetRfqOrderHash(order LibNativeOrderRfqOrder) ([32]byte, error) {
	return _ExchangeV4.Contract.GetRfqOrderHash(&_ExchangeV4.CallOpts, order)
}

// GetRfqOrderInfo is a free data retrieval call binding the contract method 0x346693c5.
//
// Solidity: function getRfqOrderInfo((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order) view returns((bytes32,uint8,uint128) orderInfo)
func (_ExchangeV4 *ExchangeV4Caller) GetRfqOrderInfo(opts *bind.CallOpts, order LibNativeOrderRfqOrder) (LibNativeOrderOrderInfo, error) {
	var (
		ret0 = new(LibNativeOrderOrderInfo)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "getRfqOrderInfo", order)
	return *ret0, err
}

// GetRfqOrderInfo is a free data retrieval call binding the contract method 0x346693c5.
//
// Solidity: function getRfqOrderInfo((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order) view returns((bytes32,uint8,uint128) orderInfo)
func (_ExchangeV4 *ExchangeV4Session) GetRfqOrderInfo(order LibNativeOrderRfqOrder) (LibNativeOrderOrderInfo, error) {
	return _ExchangeV4.Contract.GetRfqOrderInfo(&_ExchangeV4.CallOpts, order)
}

// GetRfqOrderInfo is a free data retrieval call binding the contract method 0x346693c5.
//
// Solidity: function getRfqOrderInfo((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order) view returns((bytes32,uint8,uint128) orderInfo)
func (_ExchangeV4 *ExchangeV4CallerSession) GetRfqOrderInfo(order LibNativeOrderRfqOrder) (LibNativeOrderOrderInfo, error) {
	return _ExchangeV4.Contract.GetRfqOrderInfo(&_ExchangeV4.CallOpts, order)
}

// GetRfqOrderRelevantState is a free data retrieval call binding the contract method 0x37f381d8.
//
// Solidity: function getRfqOrderRelevantState((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature) view returns((bytes32,uint8,uint128) orderInfo, uint128 actualFillableTakerTokenAmount, bool isSignatureValid)
func (_ExchangeV4 *ExchangeV4Caller) GetRfqOrderRelevantState(opts *bind.CallOpts, order LibNativeOrderRfqOrder, signature LibSignatureSignature) (struct {
	OrderInfo                      LibNativeOrderOrderInfo
	ActualFillableTakerTokenAmount *big.Int
	IsSignatureValid               bool
}, error) {
	ret := new(struct {
		OrderInfo                      LibNativeOrderOrderInfo
		ActualFillableTakerTokenAmount *big.Int
		IsSignatureValid               bool
	})
	out := ret
	err := _ExchangeV4.contract.Call(opts, out, "getRfqOrderRelevantState", order, signature)
	return *ret, err
}

// GetRfqOrderRelevantState is a free data retrieval call binding the contract method 0x37f381d8.
//
// Solidity: function getRfqOrderRelevantState((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature) view returns((bytes32,uint8,uint128) orderInfo, uint128 actualFillableTakerTokenAmount, bool isSignatureValid)
func (_ExchangeV4 *ExchangeV4Session) GetRfqOrderRelevantState(order LibNativeOrderRfqOrder, signature LibSignatureSignature) (struct {
	OrderInfo                      LibNativeOrderOrderInfo
	ActualFillableTakerTokenAmount *big.Int
	IsSignatureValid               bool
}, error) {
	return _ExchangeV4.Contract.GetRfqOrderRelevantState(&_ExchangeV4.CallOpts, order, signature)
}

// GetRfqOrderRelevantState is a free data retrieval call binding the contract method 0x37f381d8.
//
// Solidity: function getRfqOrderRelevantState((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature) view returns((bytes32,uint8,uint128) orderInfo, uint128 actualFillableTakerTokenAmount, bool isSignatureValid)
func (_ExchangeV4 *ExchangeV4CallerSession) GetRfqOrderRelevantState(order LibNativeOrderRfqOrder, signature LibSignatureSignature) (struct {
	OrderInfo                      LibNativeOrderOrderInfo
	ActualFillableTakerTokenAmount *big.Int
	IsSignatureValid               bool
}, error) {
	return _ExchangeV4.Contract.GetRfqOrderRelevantState(&_ExchangeV4.CallOpts, order, signature)
}

// GetRollbackEntryAtIndex is a free data retrieval call binding the contract method 0x6ba6bbc2.
//
// Solidity: function getRollbackEntryAtIndex(bytes4 selector, uint256 idx) view returns(address impl)
func (_ExchangeV4 *ExchangeV4Caller) GetRollbackEntryAtIndex(opts *bind.CallOpts, selector [4]byte, idx *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "getRollbackEntryAtIndex", selector, idx)
	return *ret0, err
}

// GetRollbackEntryAtIndex is a free data retrieval call binding the contract method 0x6ba6bbc2.
//
// Solidity: function getRollbackEntryAtIndex(bytes4 selector, uint256 idx) view returns(address impl)
func (_ExchangeV4 *ExchangeV4Session) GetRollbackEntryAtIndex(selector [4]byte, idx *big.Int) (common.Address, error) {
	return _ExchangeV4.Contract.GetRollbackEntryAtIndex(&_ExchangeV4.CallOpts, selector, idx)
}

// GetRollbackEntryAtIndex is a free data retrieval call binding the contract method 0x6ba6bbc2.
//
// Solidity: function getRollbackEntryAtIndex(bytes4 selector, uint256 idx) view returns(address impl)
func (_ExchangeV4 *ExchangeV4CallerSession) GetRollbackEntryAtIndex(selector [4]byte, idx *big.Int) (common.Address, error) {
	return _ExchangeV4.Contract.GetRollbackEntryAtIndex(&_ExchangeV4.CallOpts, selector, idx)
}

// GetRollbackLength is a free data retrieval call binding the contract method 0xdfd00749.
//
// Solidity: function getRollbackLength(bytes4 selector) view returns(uint256 rollbackLength)
func (_ExchangeV4 *ExchangeV4Caller) GetRollbackLength(opts *bind.CallOpts, selector [4]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "getRollbackLength", selector)
	return *ret0, err
}

// GetRollbackLength is a free data retrieval call binding the contract method 0xdfd00749.
//
// Solidity: function getRollbackLength(bytes4 selector) view returns(uint256 rollbackLength)
func (_ExchangeV4 *ExchangeV4Session) GetRollbackLength(selector [4]byte) (*big.Int, error) {
	return _ExchangeV4.Contract.GetRollbackLength(&_ExchangeV4.CallOpts, selector)
}

// GetRollbackLength is a free data retrieval call binding the contract method 0xdfd00749.
//
// Solidity: function getRollbackLength(bytes4 selector) view returns(uint256 rollbackLength)
func (_ExchangeV4 *ExchangeV4CallerSession) GetRollbackLength(selector [4]byte) (*big.Int, error) {
	return _ExchangeV4.Contract.GetRollbackLength(&_ExchangeV4.CallOpts, selector)
}

// GetSpendableERC20BalanceOf is a free data retrieval call binding the contract method 0x496f471e.
//
// Solidity: function getSpendableERC20BalanceOf(address token, address owner) view returns(uint256 amount)
func (_ExchangeV4 *ExchangeV4Caller) GetSpendableERC20BalanceOf(opts *bind.CallOpts, token common.Address, owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "getSpendableERC20BalanceOf", token, owner)
	return *ret0, err
}

// GetSpendableERC20BalanceOf is a free data retrieval call binding the contract method 0x496f471e.
//
// Solidity: function getSpendableERC20BalanceOf(address token, address owner) view returns(uint256 amount)
func (_ExchangeV4 *ExchangeV4Session) GetSpendableERC20BalanceOf(token common.Address, owner common.Address) (*big.Int, error) {
	return _ExchangeV4.Contract.GetSpendableERC20BalanceOf(&_ExchangeV4.CallOpts, token, owner)
}

// GetSpendableERC20BalanceOf is a free data retrieval call binding the contract method 0x496f471e.
//
// Solidity: function getSpendableERC20BalanceOf(address token, address owner) view returns(uint256 amount)
func (_ExchangeV4 *ExchangeV4CallerSession) GetSpendableERC20BalanceOf(token common.Address, owner common.Address) (*big.Int, error) {
	return _ExchangeV4.Contract.GetSpendableERC20BalanceOf(&_ExchangeV4.CallOpts, token, owner)
}

// GetTransformWallet is a free data retrieval call binding the contract method 0xf028e9be.
//
// Solidity: function getTransformWallet() view returns(address wallet)
func (_ExchangeV4 *ExchangeV4Caller) GetTransformWallet(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "getTransformWallet")
	return *ret0, err
}

// GetTransformWallet is a free data retrieval call binding the contract method 0xf028e9be.
//
// Solidity: function getTransformWallet() view returns(address wallet)
func (_ExchangeV4 *ExchangeV4Session) GetTransformWallet() (common.Address, error) {
	return _ExchangeV4.Contract.GetTransformWallet(&_ExchangeV4.CallOpts)
}

// GetTransformWallet is a free data retrieval call binding the contract method 0xf028e9be.
//
// Solidity: function getTransformWallet() view returns(address wallet)
func (_ExchangeV4 *ExchangeV4CallerSession) GetTransformWallet() (common.Address, error) {
	return _ExchangeV4.Contract.GetTransformWallet(&_ExchangeV4.CallOpts)
}

// GetTransformerDeployer is a free data retrieval call binding the contract method 0x4d54cdb6.
//
// Solidity: function getTransformerDeployer() view returns(address deployer)
func (_ExchangeV4 *ExchangeV4Caller) GetTransformerDeployer(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "getTransformerDeployer")
	return *ret0, err
}

// GetTransformerDeployer is a free data retrieval call binding the contract method 0x4d54cdb6.
//
// Solidity: function getTransformerDeployer() view returns(address deployer)
func (_ExchangeV4 *ExchangeV4Session) GetTransformerDeployer() (common.Address, error) {
	return _ExchangeV4.Contract.GetTransformerDeployer(&_ExchangeV4.CallOpts)
}

// GetTransformerDeployer is a free data retrieval call binding the contract method 0x4d54cdb6.
//
// Solidity: function getTransformerDeployer() view returns(address deployer)
func (_ExchangeV4 *ExchangeV4CallerSession) GetTransformerDeployer() (common.Address, error) {
	return _ExchangeV4.Contract.GetTransformerDeployer(&_ExchangeV4.CallOpts)
}

// IsValidHashSignature is a free data retrieval call binding the contract method 0x8171c407.
//
// Solidity: function isValidHashSignature(bytes32 hash, address signer, bytes signature) view returns(bool isValid)
func (_ExchangeV4 *ExchangeV4Caller) IsValidHashSignature(opts *bind.CallOpts, hash [32]byte, signer common.Address, signature []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "isValidHashSignature", hash, signer, signature)
	return *ret0, err
}

// IsValidHashSignature is a free data retrieval call binding the contract method 0x8171c407.
//
// Solidity: function isValidHashSignature(bytes32 hash, address signer, bytes signature) view returns(bool isValid)
func (_ExchangeV4 *ExchangeV4Session) IsValidHashSignature(hash [32]byte, signer common.Address, signature []byte) (bool, error) {
	return _ExchangeV4.Contract.IsValidHashSignature(&_ExchangeV4.CallOpts, hash, signer, signature)
}

// IsValidHashSignature is a free data retrieval call binding the contract method 0x8171c407.
//
// Solidity: function isValidHashSignature(bytes32 hash, address signer, bytes signature) view returns(bool isValid)
func (_ExchangeV4 *ExchangeV4CallerSession) IsValidHashSignature(hash [32]byte, signer common.Address, signature []byte) (bool, error) {
	return _ExchangeV4.Contract.IsValidHashSignature(&_ExchangeV4.CallOpts, hash, signer, signature)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address ownerAddress)
func (_ExchangeV4 *ExchangeV4Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ExchangeV4.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address ownerAddress)
func (_ExchangeV4 *ExchangeV4Session) Owner() (common.Address, error) {
	return _ExchangeV4.Contract.Owner(&_ExchangeV4.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address ownerAddress)
func (_ExchangeV4 *ExchangeV4CallerSession) Owner() (common.Address, error) {
	return _ExchangeV4.Contract.Owner(&_ExchangeV4.CallOpts)
}

// ValidateHashSignature is a free data retrieval call binding the contract method 0xf68fd38d.
//
// Solidity: function validateHashSignature(bytes32 hash, address signer, bytes signature) view returns()
func (_ExchangeV4 *ExchangeV4Caller) ValidateHashSignature(opts *bind.CallOpts, hash [32]byte, signer common.Address, signature []byte) error {
	var ()
	out := &[]interface{}{}
	err := _ExchangeV4.contract.Call(opts, out, "validateHashSignature", hash, signer, signature)
	return err
}

// ValidateHashSignature is a free data retrieval call binding the contract method 0xf68fd38d.
//
// Solidity: function validateHashSignature(bytes32 hash, address signer, bytes signature) view returns()
func (_ExchangeV4 *ExchangeV4Session) ValidateHashSignature(hash [32]byte, signer common.Address, signature []byte) error {
	return _ExchangeV4.Contract.ValidateHashSignature(&_ExchangeV4.CallOpts, hash, signer, signature)
}

// ValidateHashSignature is a free data retrieval call binding the contract method 0xf68fd38d.
//
// Solidity: function validateHashSignature(bytes32 hash, address signer, bytes signature) view returns()
func (_ExchangeV4 *ExchangeV4CallerSession) ValidateHashSignature(hash [32]byte, signer common.Address, signature []byte) error {
	return _ExchangeV4.Contract.ValidateHashSignature(&_ExchangeV4.CallOpts, hash, signer, signature)
}

// BatchCancelLimitOrders is a paid mutator transaction binding the contract method 0x9baa45a8.
//
// Solidity: function batchCancelLimitOrders((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256)[] orders) returns()
func (_ExchangeV4 *ExchangeV4Transactor) BatchCancelLimitOrders(opts *bind.TransactOpts, orders []LibNativeOrderLimitOrder) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "batchCancelLimitOrders", orders)
}

// BatchCancelLimitOrders is a paid mutator transaction binding the contract method 0x9baa45a8.
//
// Solidity: function batchCancelLimitOrders((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256)[] orders) returns()
func (_ExchangeV4 *ExchangeV4Session) BatchCancelLimitOrders(orders []LibNativeOrderLimitOrder) (*types.Transaction, error) {
	return _ExchangeV4.Contract.BatchCancelLimitOrders(&_ExchangeV4.TransactOpts, orders)
}

// BatchCancelLimitOrders is a paid mutator transaction binding the contract method 0x9baa45a8.
//
// Solidity: function batchCancelLimitOrders((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256)[] orders) returns()
func (_ExchangeV4 *ExchangeV4TransactorSession) BatchCancelLimitOrders(orders []LibNativeOrderLimitOrder) (*types.Transaction, error) {
	return _ExchangeV4.Contract.BatchCancelLimitOrders(&_ExchangeV4.TransactOpts, orders)
}

// BatchCancelPairLimitOrders is a paid mutator transaction binding the contract method 0x86a0c8d7.
//
// Solidity: function batchCancelPairLimitOrders(address[] makerTokens, address[] takerTokens, uint256[] minValidSalts) returns()
func (_ExchangeV4 *ExchangeV4Transactor) BatchCancelPairLimitOrders(opts *bind.TransactOpts, makerTokens []common.Address, takerTokens []common.Address, minValidSalts []*big.Int) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "batchCancelPairLimitOrders", makerTokens, takerTokens, minValidSalts)
}

// BatchCancelPairLimitOrders is a paid mutator transaction binding the contract method 0x86a0c8d7.
//
// Solidity: function batchCancelPairLimitOrders(address[] makerTokens, address[] takerTokens, uint256[] minValidSalts) returns()
func (_ExchangeV4 *ExchangeV4Session) BatchCancelPairLimitOrders(makerTokens []common.Address, takerTokens []common.Address, minValidSalts []*big.Int) (*types.Transaction, error) {
	return _ExchangeV4.Contract.BatchCancelPairLimitOrders(&_ExchangeV4.TransactOpts, makerTokens, takerTokens, minValidSalts)
}

// BatchCancelPairLimitOrders is a paid mutator transaction binding the contract method 0x86a0c8d7.
//
// Solidity: function batchCancelPairLimitOrders(address[] makerTokens, address[] takerTokens, uint256[] minValidSalts) returns()
func (_ExchangeV4 *ExchangeV4TransactorSession) BatchCancelPairLimitOrders(makerTokens []common.Address, takerTokens []common.Address, minValidSalts []*big.Int) (*types.Transaction, error) {
	return _ExchangeV4.Contract.BatchCancelPairLimitOrders(&_ExchangeV4.TransactOpts, makerTokens, takerTokens, minValidSalts)
}

// BatchCancelPairRfqOrders is a paid mutator transaction binding the contract method 0x0f0e8cf7.
//
// Solidity: function batchCancelPairRfqOrders(address[] makerTokens, address[] takerTokens, uint256[] minValidSalts) returns()
func (_ExchangeV4 *ExchangeV4Transactor) BatchCancelPairRfqOrders(opts *bind.TransactOpts, makerTokens []common.Address, takerTokens []common.Address, minValidSalts []*big.Int) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "batchCancelPairRfqOrders", makerTokens, takerTokens, minValidSalts)
}

// BatchCancelPairRfqOrders is a paid mutator transaction binding the contract method 0x0f0e8cf7.
//
// Solidity: function batchCancelPairRfqOrders(address[] makerTokens, address[] takerTokens, uint256[] minValidSalts) returns()
func (_ExchangeV4 *ExchangeV4Session) BatchCancelPairRfqOrders(makerTokens []common.Address, takerTokens []common.Address, minValidSalts []*big.Int) (*types.Transaction, error) {
	return _ExchangeV4.Contract.BatchCancelPairRfqOrders(&_ExchangeV4.TransactOpts, makerTokens, takerTokens, minValidSalts)
}

// BatchCancelPairRfqOrders is a paid mutator transaction binding the contract method 0x0f0e8cf7.
//
// Solidity: function batchCancelPairRfqOrders(address[] makerTokens, address[] takerTokens, uint256[] minValidSalts) returns()
func (_ExchangeV4 *ExchangeV4TransactorSession) BatchCancelPairRfqOrders(makerTokens []common.Address, takerTokens []common.Address, minValidSalts []*big.Int) (*types.Transaction, error) {
	return _ExchangeV4.Contract.BatchCancelPairRfqOrders(&_ExchangeV4.TransactOpts, makerTokens, takerTokens, minValidSalts)
}

// BatchCancelRfqOrders is a paid mutator transaction binding the contract method 0xf6e0f6a5.
//
// Solidity: function batchCancelRfqOrders((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256)[] orders) returns()
func (_ExchangeV4 *ExchangeV4Transactor) BatchCancelRfqOrders(opts *bind.TransactOpts, orders []LibNativeOrderRfqOrder) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "batchCancelRfqOrders", orders)
}

// BatchCancelRfqOrders is a paid mutator transaction binding the contract method 0xf6e0f6a5.
//
// Solidity: function batchCancelRfqOrders((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256)[] orders) returns()
func (_ExchangeV4 *ExchangeV4Session) BatchCancelRfqOrders(orders []LibNativeOrderRfqOrder) (*types.Transaction, error) {
	return _ExchangeV4.Contract.BatchCancelRfqOrders(&_ExchangeV4.TransactOpts, orders)
}

// BatchCancelRfqOrders is a paid mutator transaction binding the contract method 0xf6e0f6a5.
//
// Solidity: function batchCancelRfqOrders((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256)[] orders) returns()
func (_ExchangeV4 *ExchangeV4TransactorSession) BatchCancelRfqOrders(orders []LibNativeOrderRfqOrder) (*types.Transaction, error) {
	return _ExchangeV4.Contract.BatchCancelRfqOrders(&_ExchangeV4.TransactOpts, orders)
}

// BatchExecuteMetaTransactions is a paid mutator transaction binding the contract method 0xc5579ec8.
//
// Solidity: function batchExecuteMetaTransactions((address,address,uint256,uint256,uint256,uint256,bytes,uint256,address,uint256)[] mtxs, (uint8,uint8,bytes32,bytes32)[] signatures) payable returns(bytes[] returnResults)
func (_ExchangeV4 *ExchangeV4Transactor) BatchExecuteMetaTransactions(opts *bind.TransactOpts, mtxs []IMetaTransactionsFeatureMetaTransactionData, signatures []LibSignatureSignature) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "batchExecuteMetaTransactions", mtxs, signatures)
}

// BatchExecuteMetaTransactions is a paid mutator transaction binding the contract method 0xc5579ec8.
//
// Solidity: function batchExecuteMetaTransactions((address,address,uint256,uint256,uint256,uint256,bytes,uint256,address,uint256)[] mtxs, (uint8,uint8,bytes32,bytes32)[] signatures) payable returns(bytes[] returnResults)
func (_ExchangeV4 *ExchangeV4Session) BatchExecuteMetaTransactions(mtxs []IMetaTransactionsFeatureMetaTransactionData, signatures []LibSignatureSignature) (*types.Transaction, error) {
	return _ExchangeV4.Contract.BatchExecuteMetaTransactions(&_ExchangeV4.TransactOpts, mtxs, signatures)
}

// BatchExecuteMetaTransactions is a paid mutator transaction binding the contract method 0xc5579ec8.
//
// Solidity: function batchExecuteMetaTransactions((address,address,uint256,uint256,uint256,uint256,bytes,uint256,address,uint256)[] mtxs, (uint8,uint8,bytes32,bytes32)[] signatures) payable returns(bytes[] returnResults)
func (_ExchangeV4 *ExchangeV4TransactorSession) BatchExecuteMetaTransactions(mtxs []IMetaTransactionsFeatureMetaTransactionData, signatures []LibSignatureSignature) (*types.Transaction, error) {
	return _ExchangeV4.Contract.BatchExecuteMetaTransactions(&_ExchangeV4.TransactOpts, mtxs, signatures)
}

// CancelLimitOrder is a paid mutator transaction binding the contract method 0x7d49ec1a.
//
// Solidity: function cancelLimitOrder((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order) returns()
func (_ExchangeV4 *ExchangeV4Transactor) CancelLimitOrder(opts *bind.TransactOpts, order LibNativeOrderLimitOrder) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "cancelLimitOrder", order)
}

// CancelLimitOrder is a paid mutator transaction binding the contract method 0x7d49ec1a.
//
// Solidity: function cancelLimitOrder((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order) returns()
func (_ExchangeV4 *ExchangeV4Session) CancelLimitOrder(order LibNativeOrderLimitOrder) (*types.Transaction, error) {
	return _ExchangeV4.Contract.CancelLimitOrder(&_ExchangeV4.TransactOpts, order)
}

// CancelLimitOrder is a paid mutator transaction binding the contract method 0x7d49ec1a.
//
// Solidity: function cancelLimitOrder((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order) returns()
func (_ExchangeV4 *ExchangeV4TransactorSession) CancelLimitOrder(order LibNativeOrderLimitOrder) (*types.Transaction, error) {
	return _ExchangeV4.Contract.CancelLimitOrder(&_ExchangeV4.TransactOpts, order)
}

// CancelPairLimitOrders is a paid mutator transaction binding the contract method 0xd0a55fb0.
//
// Solidity: function cancelPairLimitOrders(address makerToken, address takerToken, uint256 minValidSalt) returns()
func (_ExchangeV4 *ExchangeV4Transactor) CancelPairLimitOrders(opts *bind.TransactOpts, makerToken common.Address, takerToken common.Address, minValidSalt *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "cancelPairLimitOrders", makerToken, takerToken, minValidSalt)
}

// CancelPairLimitOrders is a paid mutator transaction binding the contract method 0xd0a55fb0.
//
// Solidity: function cancelPairLimitOrders(address makerToken, address takerToken, uint256 minValidSalt) returns()
func (_ExchangeV4 *ExchangeV4Session) CancelPairLimitOrders(makerToken common.Address, takerToken common.Address, minValidSalt *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.Contract.CancelPairLimitOrders(&_ExchangeV4.TransactOpts, makerToken, takerToken, minValidSalt)
}

// CancelPairLimitOrders is a paid mutator transaction binding the contract method 0xd0a55fb0.
//
// Solidity: function cancelPairLimitOrders(address makerToken, address takerToken, uint256 minValidSalt) returns()
func (_ExchangeV4 *ExchangeV4TransactorSession) CancelPairLimitOrders(makerToken common.Address, takerToken common.Address, minValidSalt *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.Contract.CancelPairLimitOrders(&_ExchangeV4.TransactOpts, makerToken, takerToken, minValidSalt)
}

// CancelPairRfqOrders is a paid mutator transaction binding the contract method 0x9a4f809c.
//
// Solidity: function cancelPairRfqOrders(address makerToken, address takerToken, uint256 minValidSalt) returns()
func (_ExchangeV4 *ExchangeV4Transactor) CancelPairRfqOrders(opts *bind.TransactOpts, makerToken common.Address, takerToken common.Address, minValidSalt *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "cancelPairRfqOrders", makerToken, takerToken, minValidSalt)
}

// CancelPairRfqOrders is a paid mutator transaction binding the contract method 0x9a4f809c.
//
// Solidity: function cancelPairRfqOrders(address makerToken, address takerToken, uint256 minValidSalt) returns()
func (_ExchangeV4 *ExchangeV4Session) CancelPairRfqOrders(makerToken common.Address, takerToken common.Address, minValidSalt *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.Contract.CancelPairRfqOrders(&_ExchangeV4.TransactOpts, makerToken, takerToken, minValidSalt)
}

// CancelPairRfqOrders is a paid mutator transaction binding the contract method 0x9a4f809c.
//
// Solidity: function cancelPairRfqOrders(address makerToken, address takerToken, uint256 minValidSalt) returns()
func (_ExchangeV4 *ExchangeV4TransactorSession) CancelPairRfqOrders(makerToken common.Address, takerToken common.Address, minValidSalt *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.Contract.CancelPairRfqOrders(&_ExchangeV4.TransactOpts, makerToken, takerToken, minValidSalt)
}

// CancelRfqOrder is a paid mutator transaction binding the contract method 0xfe55a3ef.
//
// Solidity: function cancelRfqOrder((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order) returns()
func (_ExchangeV4 *ExchangeV4Transactor) CancelRfqOrder(opts *bind.TransactOpts, order LibNativeOrderRfqOrder) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "cancelRfqOrder", order)
}

// CancelRfqOrder is a paid mutator transaction binding the contract method 0xfe55a3ef.
//
// Solidity: function cancelRfqOrder((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order) returns()
func (_ExchangeV4 *ExchangeV4Session) CancelRfqOrder(order LibNativeOrderRfqOrder) (*types.Transaction, error) {
	return _ExchangeV4.Contract.CancelRfqOrder(&_ExchangeV4.TransactOpts, order)
}

// CancelRfqOrder is a paid mutator transaction binding the contract method 0xfe55a3ef.
//
// Solidity: function cancelRfqOrder((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order) returns()
func (_ExchangeV4 *ExchangeV4TransactorSession) CancelRfqOrder(order LibNativeOrderRfqOrder) (*types.Transaction, error) {
	return _ExchangeV4.Contract.CancelRfqOrder(&_ExchangeV4.TransactOpts, order)
}

// CreateTransformWallet is a paid mutator transaction binding the contract method 0x287b071b.
//
// Solidity: function createTransformWallet() returns(address wallet)
func (_ExchangeV4 *ExchangeV4Transactor) CreateTransformWallet(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "createTransformWallet")
}

// CreateTransformWallet is a paid mutator transaction binding the contract method 0x287b071b.
//
// Solidity: function createTransformWallet() returns(address wallet)
func (_ExchangeV4 *ExchangeV4Session) CreateTransformWallet() (*types.Transaction, error) {
	return _ExchangeV4.Contract.CreateTransformWallet(&_ExchangeV4.TransactOpts)
}

// CreateTransformWallet is a paid mutator transaction binding the contract method 0x287b071b.
//
// Solidity: function createTransformWallet() returns(address wallet)
func (_ExchangeV4 *ExchangeV4TransactorSession) CreateTransformWallet() (*types.Transaction, error) {
	return _ExchangeV4.Contract.CreateTransformWallet(&_ExchangeV4.TransactOpts)
}

// ExecuteMetaTransaction is a paid mutator transaction binding the contract method 0x3d61ed3e.
//
// Solidity: function executeMetaTransaction((address,address,uint256,uint256,uint256,uint256,bytes,uint256,address,uint256) mtx, (uint8,uint8,bytes32,bytes32) signature) payable returns(bytes returnResult)
func (_ExchangeV4 *ExchangeV4Transactor) ExecuteMetaTransaction(opts *bind.TransactOpts, mtx IMetaTransactionsFeatureMetaTransactionData, signature LibSignatureSignature) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "executeMetaTransaction", mtx, signature)
}

// ExecuteMetaTransaction is a paid mutator transaction binding the contract method 0x3d61ed3e.
//
// Solidity: function executeMetaTransaction((address,address,uint256,uint256,uint256,uint256,bytes,uint256,address,uint256) mtx, (uint8,uint8,bytes32,bytes32) signature) payable returns(bytes returnResult)
func (_ExchangeV4 *ExchangeV4Session) ExecuteMetaTransaction(mtx IMetaTransactionsFeatureMetaTransactionData, signature LibSignatureSignature) (*types.Transaction, error) {
	return _ExchangeV4.Contract.ExecuteMetaTransaction(&_ExchangeV4.TransactOpts, mtx, signature)
}

// ExecuteMetaTransaction is a paid mutator transaction binding the contract method 0x3d61ed3e.
//
// Solidity: function executeMetaTransaction((address,address,uint256,uint256,uint256,uint256,bytes,uint256,address,uint256) mtx, (uint8,uint8,bytes32,bytes32) signature) payable returns(bytes returnResult)
func (_ExchangeV4 *ExchangeV4TransactorSession) ExecuteMetaTransaction(mtx IMetaTransactionsFeatureMetaTransactionData, signature LibSignatureSignature) (*types.Transaction, error) {
	return _ExchangeV4.Contract.ExecuteMetaTransaction(&_ExchangeV4.TransactOpts, mtx, signature)
}

// Extend is a paid mutator transaction binding the contract method 0x6eb224cb.
//
// Solidity: function extend(bytes4 selector, address impl) returns()
func (_ExchangeV4 *ExchangeV4Transactor) Extend(opts *bind.TransactOpts, selector [4]byte, impl common.Address) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "extend", selector, impl)
}

// Extend is a paid mutator transaction binding the contract method 0x6eb224cb.
//
// Solidity: function extend(bytes4 selector, address impl) returns()
func (_ExchangeV4 *ExchangeV4Session) Extend(selector [4]byte, impl common.Address) (*types.Transaction, error) {
	return _ExchangeV4.Contract.Extend(&_ExchangeV4.TransactOpts, selector, impl)
}

// Extend is a paid mutator transaction binding the contract method 0x6eb224cb.
//
// Solidity: function extend(bytes4 selector, address impl) returns()
func (_ExchangeV4 *ExchangeV4TransactorSession) Extend(selector [4]byte, impl common.Address) (*types.Transaction, error) {
	return _ExchangeV4.Contract.Extend(&_ExchangeV4.TransactOpts, selector, impl)
}

// FillLimitOrder is a paid mutator transaction binding the contract method 0xf6274f66.
//
// Solidity: function fillLimitOrder((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature, uint128 takerTokenFillAmount) payable returns(uint128 takerTokenFilledAmount, uint128 makerTokenFilledAmount)
func (_ExchangeV4 *ExchangeV4Transactor) FillLimitOrder(opts *bind.TransactOpts, order LibNativeOrderLimitOrder, signature LibSignatureSignature, takerTokenFillAmount *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "fillLimitOrder", order, signature, takerTokenFillAmount)
}

// FillLimitOrder is a paid mutator transaction binding the contract method 0xf6274f66.
//
// Solidity: function fillLimitOrder((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature, uint128 takerTokenFillAmount) payable returns(uint128 takerTokenFilledAmount, uint128 makerTokenFilledAmount)
func (_ExchangeV4 *ExchangeV4Session) FillLimitOrder(order LibNativeOrderLimitOrder, signature LibSignatureSignature, takerTokenFillAmount *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.Contract.FillLimitOrder(&_ExchangeV4.TransactOpts, order, signature, takerTokenFillAmount)
}

// FillLimitOrder is a paid mutator transaction binding the contract method 0xf6274f66.
//
// Solidity: function fillLimitOrder((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature, uint128 takerTokenFillAmount) payable returns(uint128 takerTokenFilledAmount, uint128 makerTokenFilledAmount)
func (_ExchangeV4 *ExchangeV4TransactorSession) FillLimitOrder(order LibNativeOrderLimitOrder, signature LibSignatureSignature, takerTokenFillAmount *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.Contract.FillLimitOrder(&_ExchangeV4.TransactOpts, order, signature, takerTokenFillAmount)
}

// FillOrKillLimitOrder is a paid mutator transaction binding the contract method 0x9240529c.
//
// Solidity: function fillOrKillLimitOrder((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature, uint128 takerTokenFillAmount) payable returns(uint128 makerTokenFilledAmount)
func (_ExchangeV4 *ExchangeV4Transactor) FillOrKillLimitOrder(opts *bind.TransactOpts, order LibNativeOrderLimitOrder, signature LibSignatureSignature, takerTokenFillAmount *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "fillOrKillLimitOrder", order, signature, takerTokenFillAmount)
}

// FillOrKillLimitOrder is a paid mutator transaction binding the contract method 0x9240529c.
//
// Solidity: function fillOrKillLimitOrder((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature, uint128 takerTokenFillAmount) payable returns(uint128 makerTokenFilledAmount)
func (_ExchangeV4 *ExchangeV4Session) FillOrKillLimitOrder(order LibNativeOrderLimitOrder, signature LibSignatureSignature, takerTokenFillAmount *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.Contract.FillOrKillLimitOrder(&_ExchangeV4.TransactOpts, order, signature, takerTokenFillAmount)
}

// FillOrKillLimitOrder is a paid mutator transaction binding the contract method 0x9240529c.
//
// Solidity: function fillOrKillLimitOrder((address,address,uint128,uint128,uint128,address,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature, uint128 takerTokenFillAmount) payable returns(uint128 makerTokenFilledAmount)
func (_ExchangeV4 *ExchangeV4TransactorSession) FillOrKillLimitOrder(order LibNativeOrderLimitOrder, signature LibSignatureSignature, takerTokenFillAmount *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.Contract.FillOrKillLimitOrder(&_ExchangeV4.TransactOpts, order, signature, takerTokenFillAmount)
}

// FillOrKillRfqOrder is a paid mutator transaction binding the contract method 0x438cdfc5.
//
// Solidity: function fillOrKillRfqOrder((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature, uint128 takerTokenFillAmount) returns(uint128 makerTokenFilledAmount)
func (_ExchangeV4 *ExchangeV4Transactor) FillOrKillRfqOrder(opts *bind.TransactOpts, order LibNativeOrderRfqOrder, signature LibSignatureSignature, takerTokenFillAmount *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "fillOrKillRfqOrder", order, signature, takerTokenFillAmount)
}

// FillOrKillRfqOrder is a paid mutator transaction binding the contract method 0x438cdfc5.
//
// Solidity: function fillOrKillRfqOrder((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature, uint128 takerTokenFillAmount) returns(uint128 makerTokenFilledAmount)
func (_ExchangeV4 *ExchangeV4Session) FillOrKillRfqOrder(order LibNativeOrderRfqOrder, signature LibSignatureSignature, takerTokenFillAmount *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.Contract.FillOrKillRfqOrder(&_ExchangeV4.TransactOpts, order, signature, takerTokenFillAmount)
}

// FillOrKillRfqOrder is a paid mutator transaction binding the contract method 0x438cdfc5.
//
// Solidity: function fillOrKillRfqOrder((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature, uint128 takerTokenFillAmount) returns(uint128 makerTokenFilledAmount)
func (_ExchangeV4 *ExchangeV4TransactorSession) FillOrKillRfqOrder(order LibNativeOrderRfqOrder, signature LibSignatureSignature, takerTokenFillAmount *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.Contract.FillOrKillRfqOrder(&_ExchangeV4.TransactOpts, order, signature, takerTokenFillAmount)
}

// FillRfqOrder is a paid mutator transaction binding the contract method 0xaa77476c.
//
// Solidity: function fillRfqOrder((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature, uint128 takerTokenFillAmount) returns(uint128 takerTokenFilledAmount, uint128 makerTokenFilledAmount)
func (_ExchangeV4 *ExchangeV4Transactor) FillRfqOrder(opts *bind.TransactOpts, order LibNativeOrderRfqOrder, signature LibSignatureSignature, takerTokenFillAmount *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "fillRfqOrder", order, signature, takerTokenFillAmount)
}

// FillRfqOrder is a paid mutator transaction binding the contract method 0xaa77476c.
//
// Solidity: function fillRfqOrder((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature, uint128 takerTokenFillAmount) returns(uint128 takerTokenFilledAmount, uint128 makerTokenFilledAmount)
func (_ExchangeV4 *ExchangeV4Session) FillRfqOrder(order LibNativeOrderRfqOrder, signature LibSignatureSignature, takerTokenFillAmount *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.Contract.FillRfqOrder(&_ExchangeV4.TransactOpts, order, signature, takerTokenFillAmount)
}

// FillRfqOrder is a paid mutator transaction binding the contract method 0xaa77476c.
//
// Solidity: function fillRfqOrder((address,address,uint128,uint128,address,address,address,bytes32,uint64,uint256) order, (uint8,uint8,bytes32,bytes32) signature, uint128 takerTokenFillAmount) returns(uint128 takerTokenFilledAmount, uint128 makerTokenFilledAmount)
func (_ExchangeV4 *ExchangeV4TransactorSession) FillRfqOrder(order LibNativeOrderRfqOrder, signature LibSignatureSignature, takerTokenFillAmount *big.Int) (*types.Transaction, error) {
	return _ExchangeV4.Contract.FillRfqOrder(&_ExchangeV4.TransactOpts, order, signature, takerTokenFillAmount)
}

// Migrate is a paid mutator transaction binding the contract method 0x261fe679.
//
// Solidity: function migrate(address target, bytes data, address newOwner) returns()
func (_ExchangeV4 *ExchangeV4Transactor) Migrate(opts *bind.TransactOpts, target common.Address, data []byte, newOwner common.Address) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "migrate", target, data, newOwner)
}

// Migrate is a paid mutator transaction binding the contract method 0x261fe679.
//
// Solidity: function migrate(address target, bytes data, address newOwner) returns()
func (_ExchangeV4 *ExchangeV4Session) Migrate(target common.Address, data []byte, newOwner common.Address) (*types.Transaction, error) {
	return _ExchangeV4.Contract.Migrate(&_ExchangeV4.TransactOpts, target, data, newOwner)
}

// Migrate is a paid mutator transaction binding the contract method 0x261fe679.
//
// Solidity: function migrate(address target, bytes data, address newOwner) returns()
func (_ExchangeV4 *ExchangeV4TransactorSession) Migrate(target common.Address, data []byte, newOwner common.Address) (*types.Transaction, error) {
	return _ExchangeV4.Contract.Migrate(&_ExchangeV4.TransactOpts, target, data, newOwner)
}

// RegisterAllowedRfqOrigins is a paid mutator transaction binding the contract method 0xb09f1fb1.
//
// Solidity: function registerAllowedRfqOrigins(address[] origins, bool allowed) returns()
func (_ExchangeV4 *ExchangeV4Transactor) RegisterAllowedRfqOrigins(opts *bind.TransactOpts, origins []common.Address, allowed bool) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "registerAllowedRfqOrigins", origins, allowed)
}

// RegisterAllowedRfqOrigins is a paid mutator transaction binding the contract method 0xb09f1fb1.
//
// Solidity: function registerAllowedRfqOrigins(address[] origins, bool allowed) returns()
func (_ExchangeV4 *ExchangeV4Session) RegisterAllowedRfqOrigins(origins []common.Address, allowed bool) (*types.Transaction, error) {
	return _ExchangeV4.Contract.RegisterAllowedRfqOrigins(&_ExchangeV4.TransactOpts, origins, allowed)
}

// RegisterAllowedRfqOrigins is a paid mutator transaction binding the contract method 0xb09f1fb1.
//
// Solidity: function registerAllowedRfqOrigins(address[] origins, bool allowed) returns()
func (_ExchangeV4 *ExchangeV4TransactorSession) RegisterAllowedRfqOrigins(origins []common.Address, allowed bool) (*types.Transaction, error) {
	return _ExchangeV4.Contract.RegisterAllowedRfqOrigins(&_ExchangeV4.TransactOpts, origins, allowed)
}

// Rollback is a paid mutator transaction binding the contract method 0x9db64a40.
//
// Solidity: function rollback(bytes4 selector, address targetImpl) returns()
func (_ExchangeV4 *ExchangeV4Transactor) Rollback(opts *bind.TransactOpts, selector [4]byte, targetImpl common.Address) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "rollback", selector, targetImpl)
}

// Rollback is a paid mutator transaction binding the contract method 0x9db64a40.
//
// Solidity: function rollback(bytes4 selector, address targetImpl) returns()
func (_ExchangeV4 *ExchangeV4Session) Rollback(selector [4]byte, targetImpl common.Address) (*types.Transaction, error) {
	return _ExchangeV4.Contract.Rollback(&_ExchangeV4.TransactOpts, selector, targetImpl)
}

// Rollback is a paid mutator transaction binding the contract method 0x9db64a40.
//
// Solidity: function rollback(bytes4 selector, address targetImpl) returns()
func (_ExchangeV4 *ExchangeV4TransactorSession) Rollback(selector [4]byte, targetImpl common.Address) (*types.Transaction, error) {
	return _ExchangeV4.Contract.Rollback(&_ExchangeV4.TransactOpts, selector, targetImpl)
}

// SellToLiquidityProvider is a paid mutator transaction binding the contract method 0xf7fcd384.
//
// Solidity: function sellToLiquidityProvider(address inputToken, address outputToken, address provider, address recipient, uint256 sellAmount, uint256 minBuyAmount, bytes auxiliaryData) payable returns(uint256 boughtAmount)
func (_ExchangeV4 *ExchangeV4Transactor) SellToLiquidityProvider(opts *bind.TransactOpts, inputToken common.Address, outputToken common.Address, provider common.Address, recipient common.Address, sellAmount *big.Int, minBuyAmount *big.Int, auxiliaryData []byte) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "sellToLiquidityProvider", inputToken, outputToken, provider, recipient, sellAmount, minBuyAmount, auxiliaryData)
}

// SellToLiquidityProvider is a paid mutator transaction binding the contract method 0xf7fcd384.
//
// Solidity: function sellToLiquidityProvider(address inputToken, address outputToken, address provider, address recipient, uint256 sellAmount, uint256 minBuyAmount, bytes auxiliaryData) payable returns(uint256 boughtAmount)
func (_ExchangeV4 *ExchangeV4Session) SellToLiquidityProvider(inputToken common.Address, outputToken common.Address, provider common.Address, recipient common.Address, sellAmount *big.Int, minBuyAmount *big.Int, auxiliaryData []byte) (*types.Transaction, error) {
	return _ExchangeV4.Contract.SellToLiquidityProvider(&_ExchangeV4.TransactOpts, inputToken, outputToken, provider, recipient, sellAmount, minBuyAmount, auxiliaryData)
}

// SellToLiquidityProvider is a paid mutator transaction binding the contract method 0xf7fcd384.
//
// Solidity: function sellToLiquidityProvider(address inputToken, address outputToken, address provider, address recipient, uint256 sellAmount, uint256 minBuyAmount, bytes auxiliaryData) payable returns(uint256 boughtAmount)
func (_ExchangeV4 *ExchangeV4TransactorSession) SellToLiquidityProvider(inputToken common.Address, outputToken common.Address, provider common.Address, recipient common.Address, sellAmount *big.Int, minBuyAmount *big.Int, auxiliaryData []byte) (*types.Transaction, error) {
	return _ExchangeV4.Contract.SellToLiquidityProvider(&_ExchangeV4.TransactOpts, inputToken, outputToken, provider, recipient, sellAmount, minBuyAmount, auxiliaryData)
}

// SellToUniswap is a paid mutator transaction binding the contract method 0xd9627aa4.
//
// Solidity: function sellToUniswap(address[] tokens, uint256 sellAmount, uint256 minBuyAmount, bool isSushi) payable returns(uint256 buyAmount)
func (_ExchangeV4 *ExchangeV4Transactor) SellToUniswap(opts *bind.TransactOpts, tokens []common.Address, sellAmount *big.Int, minBuyAmount *big.Int, isSushi bool) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "sellToUniswap", tokens, sellAmount, minBuyAmount, isSushi)
}

// SellToUniswap is a paid mutator transaction binding the contract method 0xd9627aa4.
//
// Solidity: function sellToUniswap(address[] tokens, uint256 sellAmount, uint256 minBuyAmount, bool isSushi) payable returns(uint256 buyAmount)
func (_ExchangeV4 *ExchangeV4Session) SellToUniswap(tokens []common.Address, sellAmount *big.Int, minBuyAmount *big.Int, isSushi bool) (*types.Transaction, error) {
	return _ExchangeV4.Contract.SellToUniswap(&_ExchangeV4.TransactOpts, tokens, sellAmount, minBuyAmount, isSushi)
}

// SellToUniswap is a paid mutator transaction binding the contract method 0xd9627aa4.
//
// Solidity: function sellToUniswap(address[] tokens, uint256 sellAmount, uint256 minBuyAmount, bool isSushi) payable returns(uint256 buyAmount)
func (_ExchangeV4 *ExchangeV4TransactorSession) SellToUniswap(tokens []common.Address, sellAmount *big.Int, minBuyAmount *big.Int, isSushi bool) (*types.Transaction, error) {
	return _ExchangeV4.Contract.SellToUniswap(&_ExchangeV4.TransactOpts, tokens, sellAmount, minBuyAmount, isSushi)
}

// SetQuoteSigner is a paid mutator transaction binding the contract method 0x56ce180a.
//
// Solidity: function setQuoteSigner(address quoteSigner) returns()
func (_ExchangeV4 *ExchangeV4Transactor) SetQuoteSigner(opts *bind.TransactOpts, quoteSigner common.Address) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "setQuoteSigner", quoteSigner)
}

// SetQuoteSigner is a paid mutator transaction binding the contract method 0x56ce180a.
//
// Solidity: function setQuoteSigner(address quoteSigner) returns()
func (_ExchangeV4 *ExchangeV4Session) SetQuoteSigner(quoteSigner common.Address) (*types.Transaction, error) {
	return _ExchangeV4.Contract.SetQuoteSigner(&_ExchangeV4.TransactOpts, quoteSigner)
}

// SetQuoteSigner is a paid mutator transaction binding the contract method 0x56ce180a.
//
// Solidity: function setQuoteSigner(address quoteSigner) returns()
func (_ExchangeV4 *ExchangeV4TransactorSession) SetQuoteSigner(quoteSigner common.Address) (*types.Transaction, error) {
	return _ExchangeV4.Contract.SetQuoteSigner(&_ExchangeV4.TransactOpts, quoteSigner)
}

// SetTransformerDeployer is a paid mutator transaction binding the contract method 0x87c96419.
//
// Solidity: function setTransformerDeployer(address transformerDeployer) returns()
func (_ExchangeV4 *ExchangeV4Transactor) SetTransformerDeployer(opts *bind.TransactOpts, transformerDeployer common.Address) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "setTransformerDeployer", transformerDeployer)
}

// SetTransformerDeployer is a paid mutator transaction binding the contract method 0x87c96419.
//
// Solidity: function setTransformerDeployer(address transformerDeployer) returns()
func (_ExchangeV4 *ExchangeV4Session) SetTransformerDeployer(transformerDeployer common.Address) (*types.Transaction, error) {
	return _ExchangeV4.Contract.SetTransformerDeployer(&_ExchangeV4.TransactOpts, transformerDeployer)
}

// SetTransformerDeployer is a paid mutator transaction binding the contract method 0x87c96419.
//
// Solidity: function setTransformerDeployer(address transformerDeployer) returns()
func (_ExchangeV4 *ExchangeV4TransactorSession) SetTransformerDeployer(transformerDeployer common.Address) (*types.Transaction, error) {
	return _ExchangeV4.Contract.SetTransformerDeployer(&_ExchangeV4.TransactOpts, transformerDeployer)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ExchangeV4 *ExchangeV4Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ExchangeV4 *ExchangeV4Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ExchangeV4.Contract.TransferOwnership(&_ExchangeV4.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ExchangeV4 *ExchangeV4TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ExchangeV4.Contract.TransferOwnership(&_ExchangeV4.TransactOpts, newOwner)
}

// TransferProtocolFeesForPools is a paid mutator transaction binding the contract method 0x3cd2f026.
//
// Solidity: function transferProtocolFeesForPools(bytes32[] poolIds) returns()
func (_ExchangeV4 *ExchangeV4Transactor) TransferProtocolFeesForPools(opts *bind.TransactOpts, poolIds [][32]byte) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "transferProtocolFeesForPools", poolIds)
}

// TransferProtocolFeesForPools is a paid mutator transaction binding the contract method 0x3cd2f026.
//
// Solidity: function transferProtocolFeesForPools(bytes32[] poolIds) returns()
func (_ExchangeV4 *ExchangeV4Session) TransferProtocolFeesForPools(poolIds [][32]byte) (*types.Transaction, error) {
	return _ExchangeV4.Contract.TransferProtocolFeesForPools(&_ExchangeV4.TransactOpts, poolIds)
}

// TransferProtocolFeesForPools is a paid mutator transaction binding the contract method 0x3cd2f026.
//
// Solidity: function transferProtocolFeesForPools(bytes32[] poolIds) returns()
func (_ExchangeV4 *ExchangeV4TransactorSession) TransferProtocolFeesForPools(poolIds [][32]byte) (*types.Transaction, error) {
	return _ExchangeV4.Contract.TransferProtocolFeesForPools(&_ExchangeV4.TransactOpts, poolIds)
}

// TransformERC20 is a paid mutator transaction binding the contract method 0x415565b0.
//
// Solidity: function transformERC20(address inputToken, address outputToken, uint256 inputTokenAmount, uint256 minOutputTokenAmount, (uint32,bytes)[] transformations) payable returns(uint256 outputTokenAmount)
func (_ExchangeV4 *ExchangeV4Transactor) TransformERC20(opts *bind.TransactOpts, inputToken common.Address, outputToken common.Address, inputTokenAmount *big.Int, minOutputTokenAmount *big.Int, transformations []ITransformERC20FeatureTransformation) (*types.Transaction, error) {
	return _ExchangeV4.contract.Transact(opts, "transformERC20", inputToken, outputToken, inputTokenAmount, minOutputTokenAmount, transformations)
}

// TransformERC20 is a paid mutator transaction binding the contract method 0x415565b0.
//
// Solidity: function transformERC20(address inputToken, address outputToken, uint256 inputTokenAmount, uint256 minOutputTokenAmount, (uint32,bytes)[] transformations) payable returns(uint256 outputTokenAmount)
func (_ExchangeV4 *ExchangeV4Session) TransformERC20(inputToken common.Address, outputToken common.Address, inputTokenAmount *big.Int, minOutputTokenAmount *big.Int, transformations []ITransformERC20FeatureTransformation) (*types.Transaction, error) {
	return _ExchangeV4.Contract.TransformERC20(&_ExchangeV4.TransactOpts, inputToken, outputToken, inputTokenAmount, minOutputTokenAmount, transformations)
}

// TransformERC20 is a paid mutator transaction binding the contract method 0x415565b0.
//
// Solidity: function transformERC20(address inputToken, address outputToken, uint256 inputTokenAmount, uint256 minOutputTokenAmount, (uint32,bytes)[] transformations) payable returns(uint256 outputTokenAmount)
func (_ExchangeV4 *ExchangeV4TransactorSession) TransformERC20(inputToken common.Address, outputToken common.Address, inputTokenAmount *big.Int, minOutputTokenAmount *big.Int, transformations []ITransformERC20FeatureTransformation) (*types.Transaction, error) {
	return _ExchangeV4.Contract.TransformERC20(&_ExchangeV4.TransactOpts, inputToken, outputToken, inputTokenAmount, minOutputTokenAmount, transformations)
}

// ExchangeV4LimitOrderFilledIterator is returned from FilterLimitOrderFilled and is used to iterate over the raw logs and unpacked data for LimitOrderFilled events raised by the ExchangeV4 contract.
type ExchangeV4LimitOrderFilledIterator struct {
	Event *ExchangeV4LimitOrderFilled // Event containing the contract specifics and raw log

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
func (it *ExchangeV4LimitOrderFilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeV4LimitOrderFilled)
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
		it.Event = new(ExchangeV4LimitOrderFilled)
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
func (it *ExchangeV4LimitOrderFilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeV4LimitOrderFilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeV4LimitOrderFilled represents a LimitOrderFilled event raised by the ExchangeV4 contract.
type ExchangeV4LimitOrderFilled struct {
	OrderHash                 [32]byte
	Maker                     common.Address
	Taker                     common.Address
	FeeRecipient              common.Address
	MakerToken                common.Address
	TakerToken                common.Address
	TakerTokenFilledAmount    *big.Int
	MakerTokenFilledAmount    *big.Int
	TakerTokenFeeFilledAmount *big.Int
	ProtocolFeePaid           *big.Int
	Pool                      [32]byte
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterLimitOrderFilled is a free log retrieval operation binding the contract event 0xab614d2b738543c0ea21f56347cf696a3a0c42a7cbec3212a5ca22a4dcff2124.
//
// Solidity: event LimitOrderFilled(bytes32 orderHash, address maker, address taker, address feeRecipient, address makerToken, address takerToken, uint128 takerTokenFilledAmount, uint128 makerTokenFilledAmount, uint128 takerTokenFeeFilledAmount, uint256 protocolFeePaid, bytes32 pool)
func (_ExchangeV4 *ExchangeV4Filterer) FilterLimitOrderFilled(opts *bind.FilterOpts) (*ExchangeV4LimitOrderFilledIterator, error) {

	logs, sub, err := _ExchangeV4.contract.FilterLogs(opts, "LimitOrderFilled")
	if err != nil {
		return nil, err
	}
	return &ExchangeV4LimitOrderFilledIterator{contract: _ExchangeV4.contract, event: "LimitOrderFilled", logs: logs, sub: sub}, nil
}

// WatchLimitOrderFilled is a free log subscription operation binding the contract event 0xab614d2b738543c0ea21f56347cf696a3a0c42a7cbec3212a5ca22a4dcff2124.
//
// Solidity: event LimitOrderFilled(bytes32 orderHash, address maker, address taker, address feeRecipient, address makerToken, address takerToken, uint128 takerTokenFilledAmount, uint128 makerTokenFilledAmount, uint128 takerTokenFeeFilledAmount, uint256 protocolFeePaid, bytes32 pool)
func (_ExchangeV4 *ExchangeV4Filterer) WatchLimitOrderFilled(opts *bind.WatchOpts, sink chan<- *ExchangeV4LimitOrderFilled) (event.Subscription, error) {

	logs, sub, err := _ExchangeV4.contract.WatchLogs(opts, "LimitOrderFilled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeV4LimitOrderFilled)
				if err := _ExchangeV4.contract.UnpackLog(event, "LimitOrderFilled", log); err != nil {
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

// ParseLimitOrderFilled is a log parse operation binding the contract event 0xab614d2b738543c0ea21f56347cf696a3a0c42a7cbec3212a5ca22a4dcff2124.
//
// Solidity: event LimitOrderFilled(bytes32 orderHash, address maker, address taker, address feeRecipient, address makerToken, address takerToken, uint128 takerTokenFilledAmount, uint128 makerTokenFilledAmount, uint128 takerTokenFeeFilledAmount, uint256 protocolFeePaid, bytes32 pool)
func (_ExchangeV4 *ExchangeV4Filterer) ParseLimitOrderFilled(log types.Log) (*ExchangeV4LimitOrderFilled, error) {
	event := new(ExchangeV4LimitOrderFilled)
	if err := _ExchangeV4.contract.UnpackLog(event, "LimitOrderFilled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangeV4MetaTransactionExecutedIterator is returned from FilterMetaTransactionExecuted and is used to iterate over the raw logs and unpacked data for MetaTransactionExecuted events raised by the ExchangeV4 contract.
type ExchangeV4MetaTransactionExecutedIterator struct {
	Event *ExchangeV4MetaTransactionExecuted // Event containing the contract specifics and raw log

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
func (it *ExchangeV4MetaTransactionExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeV4MetaTransactionExecuted)
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
		it.Event = new(ExchangeV4MetaTransactionExecuted)
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
func (it *ExchangeV4MetaTransactionExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeV4MetaTransactionExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeV4MetaTransactionExecuted represents a MetaTransactionExecuted event raised by the ExchangeV4 contract.
type ExchangeV4MetaTransactionExecuted struct {
	Hash     [32]byte
	Selector [4]byte
	Signer   common.Address
	Sender   common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMetaTransactionExecuted is a free log retrieval operation binding the contract event 0x7f4fe3ff8ae440e1570c558da08440b26f89fb1c1f2910cd91ca6452955f121a.
//
// Solidity: event MetaTransactionExecuted(bytes32 hash, bytes4 indexed selector, address signer, address sender)
func (_ExchangeV4 *ExchangeV4Filterer) FilterMetaTransactionExecuted(opts *bind.FilterOpts, selector [][4]byte) (*ExchangeV4MetaTransactionExecutedIterator, error) {

	var selectorRule []interface{}
	for _, selectorItem := range selector {
		selectorRule = append(selectorRule, selectorItem)
	}

	logs, sub, err := _ExchangeV4.contract.FilterLogs(opts, "MetaTransactionExecuted", selectorRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeV4MetaTransactionExecutedIterator{contract: _ExchangeV4.contract, event: "MetaTransactionExecuted", logs: logs, sub: sub}, nil
}

// WatchMetaTransactionExecuted is a free log subscription operation binding the contract event 0x7f4fe3ff8ae440e1570c558da08440b26f89fb1c1f2910cd91ca6452955f121a.
//
// Solidity: event MetaTransactionExecuted(bytes32 hash, bytes4 indexed selector, address signer, address sender)
func (_ExchangeV4 *ExchangeV4Filterer) WatchMetaTransactionExecuted(opts *bind.WatchOpts, sink chan<- *ExchangeV4MetaTransactionExecuted, selector [][4]byte) (event.Subscription, error) {

	var selectorRule []interface{}
	for _, selectorItem := range selector {
		selectorRule = append(selectorRule, selectorItem)
	}

	logs, sub, err := _ExchangeV4.contract.WatchLogs(opts, "MetaTransactionExecuted", selectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeV4MetaTransactionExecuted)
				if err := _ExchangeV4.contract.UnpackLog(event, "MetaTransactionExecuted", log); err != nil {
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

// ParseMetaTransactionExecuted is a log parse operation binding the contract event 0x7f4fe3ff8ae440e1570c558da08440b26f89fb1c1f2910cd91ca6452955f121a.
//
// Solidity: event MetaTransactionExecuted(bytes32 hash, bytes4 indexed selector, address signer, address sender)
func (_ExchangeV4 *ExchangeV4Filterer) ParseMetaTransactionExecuted(log types.Log) (*ExchangeV4MetaTransactionExecuted, error) {
	event := new(ExchangeV4MetaTransactionExecuted)
	if err := _ExchangeV4.contract.UnpackLog(event, "MetaTransactionExecuted", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangeV4MigratedIterator is returned from FilterMigrated and is used to iterate over the raw logs and unpacked data for Migrated events raised by the ExchangeV4 contract.
type ExchangeV4MigratedIterator struct {
	Event *ExchangeV4Migrated // Event containing the contract specifics and raw log

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
func (it *ExchangeV4MigratedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeV4Migrated)
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
		it.Event = new(ExchangeV4Migrated)
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
func (it *ExchangeV4MigratedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeV4MigratedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeV4Migrated represents a Migrated event raised by the ExchangeV4 contract.
type ExchangeV4Migrated struct {
	Caller   common.Address
	Migrator common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterMigrated is a free log retrieval operation binding the contract event 0xe1b831b0e6f3aa16b4b1a6bd526b5cdeab4940744ca6e0251f5fe5f8caf1c81a.
//
// Solidity: event Migrated(address caller, address migrator, address newOwner)
func (_ExchangeV4 *ExchangeV4Filterer) FilterMigrated(opts *bind.FilterOpts) (*ExchangeV4MigratedIterator, error) {

	logs, sub, err := _ExchangeV4.contract.FilterLogs(opts, "Migrated")
	if err != nil {
		return nil, err
	}
	return &ExchangeV4MigratedIterator{contract: _ExchangeV4.contract, event: "Migrated", logs: logs, sub: sub}, nil
}

// WatchMigrated is a free log subscription operation binding the contract event 0xe1b831b0e6f3aa16b4b1a6bd526b5cdeab4940744ca6e0251f5fe5f8caf1c81a.
//
// Solidity: event Migrated(address caller, address migrator, address newOwner)
func (_ExchangeV4 *ExchangeV4Filterer) WatchMigrated(opts *bind.WatchOpts, sink chan<- *ExchangeV4Migrated) (event.Subscription, error) {

	logs, sub, err := _ExchangeV4.contract.WatchLogs(opts, "Migrated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeV4Migrated)
				if err := _ExchangeV4.contract.UnpackLog(event, "Migrated", log); err != nil {
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

// ParseMigrated is a log parse operation binding the contract event 0xe1b831b0e6f3aa16b4b1a6bd526b5cdeab4940744ca6e0251f5fe5f8caf1c81a.
//
// Solidity: event Migrated(address caller, address migrator, address newOwner)
func (_ExchangeV4 *ExchangeV4Filterer) ParseMigrated(log types.Log) (*ExchangeV4Migrated, error) {
	event := new(ExchangeV4Migrated)
	if err := _ExchangeV4.contract.UnpackLog(event, "Migrated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangeV4OrderCancelledIterator is returned from FilterOrderCancelled and is used to iterate over the raw logs and unpacked data for OrderCancelled events raised by the ExchangeV4 contract.
type ExchangeV4OrderCancelledIterator struct {
	Event *ExchangeV4OrderCancelled // Event containing the contract specifics and raw log

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
func (it *ExchangeV4OrderCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeV4OrderCancelled)
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
		it.Event = new(ExchangeV4OrderCancelled)
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
func (it *ExchangeV4OrderCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeV4OrderCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeV4OrderCancelled represents a OrderCancelled event raised by the ExchangeV4 contract.
type ExchangeV4OrderCancelled struct {
	OrderHash [32]byte
	Maker     common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOrderCancelled is a free log retrieval operation binding the contract event 0xa6eb7cdc219e1518ced964e9a34e61d68a94e4f1569db3e84256ba981ba52753.
//
// Solidity: event OrderCancelled(bytes32 orderHash, address maker)
func (_ExchangeV4 *ExchangeV4Filterer) FilterOrderCancelled(opts *bind.FilterOpts) (*ExchangeV4OrderCancelledIterator, error) {

	logs, sub, err := _ExchangeV4.contract.FilterLogs(opts, "OrderCancelled")
	if err != nil {
		return nil, err
	}
	return &ExchangeV4OrderCancelledIterator{contract: _ExchangeV4.contract, event: "OrderCancelled", logs: logs, sub: sub}, nil
}

// WatchOrderCancelled is a free log subscription operation binding the contract event 0xa6eb7cdc219e1518ced964e9a34e61d68a94e4f1569db3e84256ba981ba52753.
//
// Solidity: event OrderCancelled(bytes32 orderHash, address maker)
func (_ExchangeV4 *ExchangeV4Filterer) WatchOrderCancelled(opts *bind.WatchOpts, sink chan<- *ExchangeV4OrderCancelled) (event.Subscription, error) {

	logs, sub, err := _ExchangeV4.contract.WatchLogs(opts, "OrderCancelled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeV4OrderCancelled)
				if err := _ExchangeV4.contract.UnpackLog(event, "OrderCancelled", log); err != nil {
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

// ParseOrderCancelled is a log parse operation binding the contract event 0xa6eb7cdc219e1518ced964e9a34e61d68a94e4f1569db3e84256ba981ba52753.
//
// Solidity: event OrderCancelled(bytes32 orderHash, address maker)
func (_ExchangeV4 *ExchangeV4Filterer) ParseOrderCancelled(log types.Log) (*ExchangeV4OrderCancelled, error) {
	event := new(ExchangeV4OrderCancelled)
	if err := _ExchangeV4.contract.UnpackLog(event, "OrderCancelled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangeV4OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ExchangeV4 contract.
type ExchangeV4OwnershipTransferredIterator struct {
	Event *ExchangeV4OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ExchangeV4OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeV4OwnershipTransferred)
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
		it.Event = new(ExchangeV4OwnershipTransferred)
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
func (it *ExchangeV4OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeV4OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeV4OwnershipTransferred represents a OwnershipTransferred event raised by the ExchangeV4 contract.
type ExchangeV4OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ExchangeV4 *ExchangeV4Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ExchangeV4OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ExchangeV4.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeV4OwnershipTransferredIterator{contract: _ExchangeV4.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ExchangeV4 *ExchangeV4Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ExchangeV4OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ExchangeV4.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeV4OwnershipTransferred)
				if err := _ExchangeV4.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ExchangeV4 *ExchangeV4Filterer) ParseOwnershipTransferred(log types.Log) (*ExchangeV4OwnershipTransferred, error) {
	event := new(ExchangeV4OwnershipTransferred)
	if err := _ExchangeV4.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangeV4PairCancelledLimitOrdersIterator is returned from FilterPairCancelledLimitOrders and is used to iterate over the raw logs and unpacked data for PairCancelledLimitOrders events raised by the ExchangeV4 contract.
type ExchangeV4PairCancelledLimitOrdersIterator struct {
	Event *ExchangeV4PairCancelledLimitOrders // Event containing the contract specifics and raw log

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
func (it *ExchangeV4PairCancelledLimitOrdersIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeV4PairCancelledLimitOrders)
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
		it.Event = new(ExchangeV4PairCancelledLimitOrders)
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
func (it *ExchangeV4PairCancelledLimitOrdersIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeV4PairCancelledLimitOrdersIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeV4PairCancelledLimitOrders represents a PairCancelledLimitOrders event raised by the ExchangeV4 contract.
type ExchangeV4PairCancelledLimitOrders struct {
	Maker        common.Address
	MakerToken   common.Address
	TakerToken   common.Address
	MinValidSalt *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterPairCancelledLimitOrders is a free log retrieval operation binding the contract event 0xa91fe7ae62fce669df2c7f880f8c14d178531aae72515558e5c948e37c32a572.
//
// Solidity: event PairCancelledLimitOrders(address maker, address makerToken, address takerToken, uint256 minValidSalt)
func (_ExchangeV4 *ExchangeV4Filterer) FilterPairCancelledLimitOrders(opts *bind.FilterOpts) (*ExchangeV4PairCancelledLimitOrdersIterator, error) {

	logs, sub, err := _ExchangeV4.contract.FilterLogs(opts, "PairCancelledLimitOrders")
	if err != nil {
		return nil, err
	}
	return &ExchangeV4PairCancelledLimitOrdersIterator{contract: _ExchangeV4.contract, event: "PairCancelledLimitOrders", logs: logs, sub: sub}, nil
}

// WatchPairCancelledLimitOrders is a free log subscription operation binding the contract event 0xa91fe7ae62fce669df2c7f880f8c14d178531aae72515558e5c948e37c32a572.
//
// Solidity: event PairCancelledLimitOrders(address maker, address makerToken, address takerToken, uint256 minValidSalt)
func (_ExchangeV4 *ExchangeV4Filterer) WatchPairCancelledLimitOrders(opts *bind.WatchOpts, sink chan<- *ExchangeV4PairCancelledLimitOrders) (event.Subscription, error) {

	logs, sub, err := _ExchangeV4.contract.WatchLogs(opts, "PairCancelledLimitOrders")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeV4PairCancelledLimitOrders)
				if err := _ExchangeV4.contract.UnpackLog(event, "PairCancelledLimitOrders", log); err != nil {
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

// ParsePairCancelledLimitOrders is a log parse operation binding the contract event 0xa91fe7ae62fce669df2c7f880f8c14d178531aae72515558e5c948e37c32a572.
//
// Solidity: event PairCancelledLimitOrders(address maker, address makerToken, address takerToken, uint256 minValidSalt)
func (_ExchangeV4 *ExchangeV4Filterer) ParsePairCancelledLimitOrders(log types.Log) (*ExchangeV4PairCancelledLimitOrders, error) {
	event := new(ExchangeV4PairCancelledLimitOrders)
	if err := _ExchangeV4.contract.UnpackLog(event, "PairCancelledLimitOrders", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangeV4PairCancelledRfqOrdersIterator is returned from FilterPairCancelledRfqOrders and is used to iterate over the raw logs and unpacked data for PairCancelledRfqOrders events raised by the ExchangeV4 contract.
type ExchangeV4PairCancelledRfqOrdersIterator struct {
	Event *ExchangeV4PairCancelledRfqOrders // Event containing the contract specifics and raw log

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
func (it *ExchangeV4PairCancelledRfqOrdersIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeV4PairCancelledRfqOrders)
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
		it.Event = new(ExchangeV4PairCancelledRfqOrders)
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
func (it *ExchangeV4PairCancelledRfqOrdersIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeV4PairCancelledRfqOrdersIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeV4PairCancelledRfqOrders represents a PairCancelledRfqOrders event raised by the ExchangeV4 contract.
type ExchangeV4PairCancelledRfqOrders struct {
	Maker        common.Address
	MakerToken   common.Address
	TakerToken   common.Address
	MinValidSalt *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterPairCancelledRfqOrders is a free log retrieval operation binding the contract event 0xfe7ffb1edfe79f4df716cb2dcad21cf2f31b104d816a7976ba1e6e4653c1efb1.
//
// Solidity: event PairCancelledRfqOrders(address maker, address makerToken, address takerToken, uint256 minValidSalt)
func (_ExchangeV4 *ExchangeV4Filterer) FilterPairCancelledRfqOrders(opts *bind.FilterOpts) (*ExchangeV4PairCancelledRfqOrdersIterator, error) {

	logs, sub, err := _ExchangeV4.contract.FilterLogs(opts, "PairCancelledRfqOrders")
	if err != nil {
		return nil, err
	}
	return &ExchangeV4PairCancelledRfqOrdersIterator{contract: _ExchangeV4.contract, event: "PairCancelledRfqOrders", logs: logs, sub: sub}, nil
}

// WatchPairCancelledRfqOrders is a free log subscription operation binding the contract event 0xfe7ffb1edfe79f4df716cb2dcad21cf2f31b104d816a7976ba1e6e4653c1efb1.
//
// Solidity: event PairCancelledRfqOrders(address maker, address makerToken, address takerToken, uint256 minValidSalt)
func (_ExchangeV4 *ExchangeV4Filterer) WatchPairCancelledRfqOrders(opts *bind.WatchOpts, sink chan<- *ExchangeV4PairCancelledRfqOrders) (event.Subscription, error) {

	logs, sub, err := _ExchangeV4.contract.WatchLogs(opts, "PairCancelledRfqOrders")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeV4PairCancelledRfqOrders)
				if err := _ExchangeV4.contract.UnpackLog(event, "PairCancelledRfqOrders", log); err != nil {
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

// ParsePairCancelledRfqOrders is a log parse operation binding the contract event 0xfe7ffb1edfe79f4df716cb2dcad21cf2f31b104d816a7976ba1e6e4653c1efb1.
//
// Solidity: event PairCancelledRfqOrders(address maker, address makerToken, address takerToken, uint256 minValidSalt)
func (_ExchangeV4 *ExchangeV4Filterer) ParsePairCancelledRfqOrders(log types.Log) (*ExchangeV4PairCancelledRfqOrders, error) {
	event := new(ExchangeV4PairCancelledRfqOrders)
	if err := _ExchangeV4.contract.UnpackLog(event, "PairCancelledRfqOrders", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangeV4ProxyFunctionUpdatedIterator is returned from FilterProxyFunctionUpdated and is used to iterate over the raw logs and unpacked data for ProxyFunctionUpdated events raised by the ExchangeV4 contract.
type ExchangeV4ProxyFunctionUpdatedIterator struct {
	Event *ExchangeV4ProxyFunctionUpdated // Event containing the contract specifics and raw log

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
func (it *ExchangeV4ProxyFunctionUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeV4ProxyFunctionUpdated)
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
		it.Event = new(ExchangeV4ProxyFunctionUpdated)
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
func (it *ExchangeV4ProxyFunctionUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeV4ProxyFunctionUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeV4ProxyFunctionUpdated represents a ProxyFunctionUpdated event raised by the ExchangeV4 contract.
type ExchangeV4ProxyFunctionUpdated struct {
	Selector [4]byte
	OldImpl  common.Address
	NewImpl  common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterProxyFunctionUpdated is a free log retrieval operation binding the contract event 0x2ae221083467de52078b0096696ab88d8d53a7ecb44bb65b56a2bab687598367.
//
// Solidity: event ProxyFunctionUpdated(bytes4 indexed selector, address oldImpl, address newImpl)
func (_ExchangeV4 *ExchangeV4Filterer) FilterProxyFunctionUpdated(opts *bind.FilterOpts, selector [][4]byte) (*ExchangeV4ProxyFunctionUpdatedIterator, error) {

	var selectorRule []interface{}
	for _, selectorItem := range selector {
		selectorRule = append(selectorRule, selectorItem)
	}

	logs, sub, err := _ExchangeV4.contract.FilterLogs(opts, "ProxyFunctionUpdated", selectorRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeV4ProxyFunctionUpdatedIterator{contract: _ExchangeV4.contract, event: "ProxyFunctionUpdated", logs: logs, sub: sub}, nil
}

// WatchProxyFunctionUpdated is a free log subscription operation binding the contract event 0x2ae221083467de52078b0096696ab88d8d53a7ecb44bb65b56a2bab687598367.
//
// Solidity: event ProxyFunctionUpdated(bytes4 indexed selector, address oldImpl, address newImpl)
func (_ExchangeV4 *ExchangeV4Filterer) WatchProxyFunctionUpdated(opts *bind.WatchOpts, sink chan<- *ExchangeV4ProxyFunctionUpdated, selector [][4]byte) (event.Subscription, error) {

	var selectorRule []interface{}
	for _, selectorItem := range selector {
		selectorRule = append(selectorRule, selectorItem)
	}

	logs, sub, err := _ExchangeV4.contract.WatchLogs(opts, "ProxyFunctionUpdated", selectorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeV4ProxyFunctionUpdated)
				if err := _ExchangeV4.contract.UnpackLog(event, "ProxyFunctionUpdated", log); err != nil {
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

// ParseProxyFunctionUpdated is a log parse operation binding the contract event 0x2ae221083467de52078b0096696ab88d8d53a7ecb44bb65b56a2bab687598367.
//
// Solidity: event ProxyFunctionUpdated(bytes4 indexed selector, address oldImpl, address newImpl)
func (_ExchangeV4 *ExchangeV4Filterer) ParseProxyFunctionUpdated(log types.Log) (*ExchangeV4ProxyFunctionUpdated, error) {
	event := new(ExchangeV4ProxyFunctionUpdated)
	if err := _ExchangeV4.contract.UnpackLog(event, "ProxyFunctionUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangeV4QuoteSignerUpdatedIterator is returned from FilterQuoteSignerUpdated and is used to iterate over the raw logs and unpacked data for QuoteSignerUpdated events raised by the ExchangeV4 contract.
type ExchangeV4QuoteSignerUpdatedIterator struct {
	Event *ExchangeV4QuoteSignerUpdated // Event containing the contract specifics and raw log

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
func (it *ExchangeV4QuoteSignerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeV4QuoteSignerUpdated)
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
		it.Event = new(ExchangeV4QuoteSignerUpdated)
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
func (it *ExchangeV4QuoteSignerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeV4QuoteSignerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeV4QuoteSignerUpdated represents a QuoteSignerUpdated event raised by the ExchangeV4 contract.
type ExchangeV4QuoteSignerUpdated struct {
	QuoteSigner common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterQuoteSignerUpdated is a free log retrieval operation binding the contract event 0xf5550c5eea19b48ac6eb5f03abdc4f59c0a60697abb3d973cd68669703b5c8b9.
//
// Solidity: event QuoteSignerUpdated(address quoteSigner)
func (_ExchangeV4 *ExchangeV4Filterer) FilterQuoteSignerUpdated(opts *bind.FilterOpts) (*ExchangeV4QuoteSignerUpdatedIterator, error) {

	logs, sub, err := _ExchangeV4.contract.FilterLogs(opts, "QuoteSignerUpdated")
	if err != nil {
		return nil, err
	}
	return &ExchangeV4QuoteSignerUpdatedIterator{contract: _ExchangeV4.contract, event: "QuoteSignerUpdated", logs: logs, sub: sub}, nil
}

// WatchQuoteSignerUpdated is a free log subscription operation binding the contract event 0xf5550c5eea19b48ac6eb5f03abdc4f59c0a60697abb3d973cd68669703b5c8b9.
//
// Solidity: event QuoteSignerUpdated(address quoteSigner)
func (_ExchangeV4 *ExchangeV4Filterer) WatchQuoteSignerUpdated(opts *bind.WatchOpts, sink chan<- *ExchangeV4QuoteSignerUpdated) (event.Subscription, error) {

	logs, sub, err := _ExchangeV4.contract.WatchLogs(opts, "QuoteSignerUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeV4QuoteSignerUpdated)
				if err := _ExchangeV4.contract.UnpackLog(event, "QuoteSignerUpdated", log); err != nil {
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

// ParseQuoteSignerUpdated is a log parse operation binding the contract event 0xf5550c5eea19b48ac6eb5f03abdc4f59c0a60697abb3d973cd68669703b5c8b9.
//
// Solidity: event QuoteSignerUpdated(address quoteSigner)
func (_ExchangeV4 *ExchangeV4Filterer) ParseQuoteSignerUpdated(log types.Log) (*ExchangeV4QuoteSignerUpdated, error) {
	event := new(ExchangeV4QuoteSignerUpdated)
	if err := _ExchangeV4.contract.UnpackLog(event, "QuoteSignerUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangeV4RfqOrderFilledIterator is returned from FilterRfqOrderFilled and is used to iterate over the raw logs and unpacked data for RfqOrderFilled events raised by the ExchangeV4 contract.
type ExchangeV4RfqOrderFilledIterator struct {
	Event *ExchangeV4RfqOrderFilled // Event containing the contract specifics and raw log

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
func (it *ExchangeV4RfqOrderFilledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeV4RfqOrderFilled)
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
		it.Event = new(ExchangeV4RfqOrderFilled)
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
func (it *ExchangeV4RfqOrderFilledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeV4RfqOrderFilledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeV4RfqOrderFilled represents a RfqOrderFilled event raised by the ExchangeV4 contract.
type ExchangeV4RfqOrderFilled struct {
	OrderHash              [32]byte
	Maker                  common.Address
	Taker                  common.Address
	MakerToken             common.Address
	TakerToken             common.Address
	TakerTokenFilledAmount *big.Int
	MakerTokenFilledAmount *big.Int
	Pool                   [32]byte
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterRfqOrderFilled is a free log retrieval operation binding the contract event 0x829fa99d94dc4636925b38632e625736a614c154d55006b7ab6bea979c210c32.
//
// Solidity: event RfqOrderFilled(bytes32 orderHash, address maker, address taker, address makerToken, address takerToken, uint128 takerTokenFilledAmount, uint128 makerTokenFilledAmount, bytes32 pool)
func (_ExchangeV4 *ExchangeV4Filterer) FilterRfqOrderFilled(opts *bind.FilterOpts) (*ExchangeV4RfqOrderFilledIterator, error) {

	logs, sub, err := _ExchangeV4.contract.FilterLogs(opts, "RfqOrderFilled")
	if err != nil {
		return nil, err
	}
	return &ExchangeV4RfqOrderFilledIterator{contract: _ExchangeV4.contract, event: "RfqOrderFilled", logs: logs, sub: sub}, nil
}

// WatchRfqOrderFilled is a free log subscription operation binding the contract event 0x829fa99d94dc4636925b38632e625736a614c154d55006b7ab6bea979c210c32.
//
// Solidity: event RfqOrderFilled(bytes32 orderHash, address maker, address taker, address makerToken, address takerToken, uint128 takerTokenFilledAmount, uint128 makerTokenFilledAmount, bytes32 pool)
func (_ExchangeV4 *ExchangeV4Filterer) WatchRfqOrderFilled(opts *bind.WatchOpts, sink chan<- *ExchangeV4RfqOrderFilled) (event.Subscription, error) {

	logs, sub, err := _ExchangeV4.contract.WatchLogs(opts, "RfqOrderFilled")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeV4RfqOrderFilled)
				if err := _ExchangeV4.contract.UnpackLog(event, "RfqOrderFilled", log); err != nil {
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

// ParseRfqOrderFilled is a log parse operation binding the contract event 0x829fa99d94dc4636925b38632e625736a614c154d55006b7ab6bea979c210c32.
//
// Solidity: event RfqOrderFilled(bytes32 orderHash, address maker, address taker, address makerToken, address takerToken, uint128 takerTokenFilledAmount, uint128 makerTokenFilledAmount, bytes32 pool)
func (_ExchangeV4 *ExchangeV4Filterer) ParseRfqOrderFilled(log types.Log) (*ExchangeV4RfqOrderFilled, error) {
	event := new(ExchangeV4RfqOrderFilled)
	if err := _ExchangeV4.contract.UnpackLog(event, "RfqOrderFilled", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangeV4RfqOrderOriginsAllowedIterator is returned from FilterRfqOrderOriginsAllowed and is used to iterate over the raw logs and unpacked data for RfqOrderOriginsAllowed events raised by the ExchangeV4 contract.
type ExchangeV4RfqOrderOriginsAllowedIterator struct {
	Event *ExchangeV4RfqOrderOriginsAllowed // Event containing the contract specifics and raw log

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
func (it *ExchangeV4RfqOrderOriginsAllowedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeV4RfqOrderOriginsAllowed)
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
		it.Event = new(ExchangeV4RfqOrderOriginsAllowed)
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
func (it *ExchangeV4RfqOrderOriginsAllowedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeV4RfqOrderOriginsAllowedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeV4RfqOrderOriginsAllowed represents a RfqOrderOriginsAllowed event raised by the ExchangeV4 contract.
type ExchangeV4RfqOrderOriginsAllowed struct {
	Origin  common.Address
	Addrs   []common.Address
	Allowed bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRfqOrderOriginsAllowed is a free log retrieval operation binding the contract event 0x02dfead5eb769b298e82dd9650b31c40559a3d42701dbf53c931bc2682847c31.
//
// Solidity: event RfqOrderOriginsAllowed(address origin, address[] addrs, bool allowed)
func (_ExchangeV4 *ExchangeV4Filterer) FilterRfqOrderOriginsAllowed(opts *bind.FilterOpts) (*ExchangeV4RfqOrderOriginsAllowedIterator, error) {

	logs, sub, err := _ExchangeV4.contract.FilterLogs(opts, "RfqOrderOriginsAllowed")
	if err != nil {
		return nil, err
	}
	return &ExchangeV4RfqOrderOriginsAllowedIterator{contract: _ExchangeV4.contract, event: "RfqOrderOriginsAllowed", logs: logs, sub: sub}, nil
}

// WatchRfqOrderOriginsAllowed is a free log subscription operation binding the contract event 0x02dfead5eb769b298e82dd9650b31c40559a3d42701dbf53c931bc2682847c31.
//
// Solidity: event RfqOrderOriginsAllowed(address origin, address[] addrs, bool allowed)
func (_ExchangeV4 *ExchangeV4Filterer) WatchRfqOrderOriginsAllowed(opts *bind.WatchOpts, sink chan<- *ExchangeV4RfqOrderOriginsAllowed) (event.Subscription, error) {

	logs, sub, err := _ExchangeV4.contract.WatchLogs(opts, "RfqOrderOriginsAllowed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeV4RfqOrderOriginsAllowed)
				if err := _ExchangeV4.contract.UnpackLog(event, "RfqOrderOriginsAllowed", log); err != nil {
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

// ParseRfqOrderOriginsAllowed is a log parse operation binding the contract event 0x02dfead5eb769b298e82dd9650b31c40559a3d42701dbf53c931bc2682847c31.
//
// Solidity: event RfqOrderOriginsAllowed(address origin, address[] addrs, bool allowed)
func (_ExchangeV4 *ExchangeV4Filterer) ParseRfqOrderOriginsAllowed(log types.Log) (*ExchangeV4RfqOrderOriginsAllowed, error) {
	event := new(ExchangeV4RfqOrderOriginsAllowed)
	if err := _ExchangeV4.contract.UnpackLog(event, "RfqOrderOriginsAllowed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangeV4TransformedERC20Iterator is returned from FilterTransformedERC20 and is used to iterate over the raw logs and unpacked data for TransformedERC20 events raised by the ExchangeV4 contract.
type ExchangeV4TransformedERC20Iterator struct {
	Event *ExchangeV4TransformedERC20 // Event containing the contract specifics and raw log

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
func (it *ExchangeV4TransformedERC20Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeV4TransformedERC20)
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
		it.Event = new(ExchangeV4TransformedERC20)
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
func (it *ExchangeV4TransformedERC20Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeV4TransformedERC20Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeV4TransformedERC20 represents a TransformedERC20 event raised by the ExchangeV4 contract.
type ExchangeV4TransformedERC20 struct {
	Taker             common.Address
	InputToken        common.Address
	OutputToken       common.Address
	InputTokenAmount  *big.Int
	OutputTokenAmount *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterTransformedERC20 is a free log retrieval operation binding the contract event 0x0f6672f78a59ba8e5e5b5d38df3ebc67f3c792e2c9259b8d97d7f00dd78ba1b3.
//
// Solidity: event TransformedERC20(address indexed taker, address inputToken, address outputToken, uint256 inputTokenAmount, uint256 outputTokenAmount)
func (_ExchangeV4 *ExchangeV4Filterer) FilterTransformedERC20(opts *bind.FilterOpts, taker []common.Address) (*ExchangeV4TransformedERC20Iterator, error) {

	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	logs, sub, err := _ExchangeV4.contract.FilterLogs(opts, "TransformedERC20", takerRule)
	if err != nil {
		return nil, err
	}
	return &ExchangeV4TransformedERC20Iterator{contract: _ExchangeV4.contract, event: "TransformedERC20", logs: logs, sub: sub}, nil
}

// WatchTransformedERC20 is a free log subscription operation binding the contract event 0x0f6672f78a59ba8e5e5b5d38df3ebc67f3c792e2c9259b8d97d7f00dd78ba1b3.
//
// Solidity: event TransformedERC20(address indexed taker, address inputToken, address outputToken, uint256 inputTokenAmount, uint256 outputTokenAmount)
func (_ExchangeV4 *ExchangeV4Filterer) WatchTransformedERC20(opts *bind.WatchOpts, sink chan<- *ExchangeV4TransformedERC20, taker []common.Address) (event.Subscription, error) {

	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	logs, sub, err := _ExchangeV4.contract.WatchLogs(opts, "TransformedERC20", takerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeV4TransformedERC20)
				if err := _ExchangeV4.contract.UnpackLog(event, "TransformedERC20", log); err != nil {
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

// ParseTransformedERC20 is a log parse operation binding the contract event 0x0f6672f78a59ba8e5e5b5d38df3ebc67f3c792e2c9259b8d97d7f00dd78ba1b3.
//
// Solidity: event TransformedERC20(address indexed taker, address inputToken, address outputToken, uint256 inputTokenAmount, uint256 outputTokenAmount)
func (_ExchangeV4 *ExchangeV4Filterer) ParseTransformedERC20(log types.Log) (*ExchangeV4TransformedERC20, error) {
	event := new(ExchangeV4TransformedERC20)
	if err := _ExchangeV4.contract.UnpackLog(event, "TransformedERC20", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ExchangeV4TransformerDeployerUpdatedIterator is returned from FilterTransformerDeployerUpdated and is used to iterate over the raw logs and unpacked data for TransformerDeployerUpdated events raised by the ExchangeV4 contract.
type ExchangeV4TransformerDeployerUpdatedIterator struct {
	Event *ExchangeV4TransformerDeployerUpdated // Event containing the contract specifics and raw log

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
func (it *ExchangeV4TransformerDeployerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExchangeV4TransformerDeployerUpdated)
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
		it.Event = new(ExchangeV4TransformerDeployerUpdated)
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
func (it *ExchangeV4TransformerDeployerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExchangeV4TransformerDeployerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExchangeV4TransformerDeployerUpdated represents a TransformerDeployerUpdated event raised by the ExchangeV4 contract.
type ExchangeV4TransformerDeployerUpdated struct {
	TransformerDeployer common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterTransformerDeployerUpdated is a free log retrieval operation binding the contract event 0xfd45604abad79c16e23348a137ed8292661be1b8eba6e4806ebed6833b1c046a.
//
// Solidity: event TransformerDeployerUpdated(address transformerDeployer)
func (_ExchangeV4 *ExchangeV4Filterer) FilterTransformerDeployerUpdated(opts *bind.FilterOpts) (*ExchangeV4TransformerDeployerUpdatedIterator, error) {

	logs, sub, err := _ExchangeV4.contract.FilterLogs(opts, "TransformerDeployerUpdated")
	if err != nil {
		return nil, err
	}
	return &ExchangeV4TransformerDeployerUpdatedIterator{contract: _ExchangeV4.contract, event: "TransformerDeployerUpdated", logs: logs, sub: sub}, nil
}

// WatchTransformerDeployerUpdated is a free log subscription operation binding the contract event 0xfd45604abad79c16e23348a137ed8292661be1b8eba6e4806ebed6833b1c046a.
//
// Solidity: event TransformerDeployerUpdated(address transformerDeployer)
func (_ExchangeV4 *ExchangeV4Filterer) WatchTransformerDeployerUpdated(opts *bind.WatchOpts, sink chan<- *ExchangeV4TransformerDeployerUpdated) (event.Subscription, error) {

	logs, sub, err := _ExchangeV4.contract.WatchLogs(opts, "TransformerDeployerUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExchangeV4TransformerDeployerUpdated)
				if err := _ExchangeV4.contract.UnpackLog(event, "TransformerDeployerUpdated", log); err != nil {
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

// ParseTransformerDeployerUpdated is a log parse operation binding the contract event 0xfd45604abad79c16e23348a137ed8292661be1b8eba6e4806ebed6833b1c046a.
//
// Solidity: event TransformerDeployerUpdated(address transformerDeployer)
func (_ExchangeV4 *ExchangeV4Filterer) ParseTransformerDeployerUpdated(log types.Log) (*ExchangeV4TransformerDeployerUpdated, error) {
	event := new(ExchangeV4TransformerDeployerUpdated)
	if err := _ExchangeV4.contract.UnpackLog(event, "TransformerDeployerUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}
