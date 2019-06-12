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

// OrderValidationUtilsABI is the input ABI used to generate the binding from.
const OrderValidationUtilsABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC721AssetData\",\"outputs\":[{\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"name\":\"tokenAddress\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"ownerAddress\",\"type\":\"address\"},{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getBalanceAndAssetProxyAllowance\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"},{\"name\":\"allowance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ERC1155_PROXY_ID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"ownerAddress\",\"type\":\"address\"},{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getTransferableAssetAmount\",\"outputs\":[{\"name\":\"transferableAssetAmount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"ownerAddress\",\"type\":\"address\"},{\"name\":\"assetData\",\"type\":\"bytes[]\"}],\"name\":\"getBatchAssetProxyAllowances\",\"outputs\":[{\"name\":\"allowances\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"encodeERC20AssetData\",\"outputs\":[{\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"ownerAddress\",\"type\":\"address\"},{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getBalance\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"getOrderRelevantStates\",\"outputs\":[{\"components\":[{\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"name\":\"ordersInfo\",\"type\":\"tuple[]\"},{\"name\":\"fillableTakerAssetAmounts\",\"type\":\"uint256[]\"},{\"name\":\"isValidSignature\",\"type\":\"bool[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ERC20_PROXY_ID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC20AssetData\",\"outputs\":[{\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"name\":\"tokenAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"getOrderRelevantState\",\"outputs\":[{\"components\":[{\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"name\":\"orderInfo\",\"type\":\"tuple\"},{\"name\":\"fillableTakerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"isValidSignature\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC1155AssetData\",\"outputs\":[{\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"name\":\"tokenAddress\",\"type\":\"address\"},{\"name\":\"tokenIds\",\"type\":\"uint256[]\"},{\"name\":\"tokenValues\",\"type\":\"uint256[]\"},{\"name\":\"callbackData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ERC721_PROXY_ID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"encodeERC721AssetData\",\"outputs\":[{\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MULTI_ASSET_PROXY_ID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\"},{\"name\":\"tokenIds\",\"type\":\"uint256[]\"},{\"name\":\"tokenValues\",\"type\":\"uint256[]\"},{\"name\":\"callbackData\",\"type\":\"bytes\"}],\"name\":\"encodeERC1155AssetData\",\"outputs\":[{\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getERC721TokenOwner\",\"outputs\":[{\"name\":\"ownerAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeMultiAssetData\",\"outputs\":[{\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"name\":\"nestedAssetData\",\"type\":\"bytes[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"ownerAddress\",\"type\":\"address\"},{\"name\":\"assetData\",\"type\":\"bytes[]\"}],\"name\":\"getBatchBalances\",\"outputs\":[{\"name\":\"balances\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"ownerAddress\",\"type\":\"address\"},{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getAssetProxyAllowance\",\"outputs\":[{\"name\":\"allowance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"name\":\"nestedAssetData\",\"type\":\"bytes[]\"}],\"name\":\"encodeMultiAssetData\",\"outputs\":[{\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"ownerAddress\",\"type\":\"address\"},{\"name\":\"assetData\",\"type\":\"bytes[]\"}],\"name\":\"getBatchBalancesAndAssetProxyAllowances\",\"outputs\":[{\"name\":\"balances\",\"type\":\"uint256[]\"},{\"name\":\"allowances\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_exchange\",\"type\":\"address\"},{\"name\":\"_zrxAssetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// OrderValidationUtils is an auto generated Go binding around an Ethereum contract.
type OrderValidationUtils struct {
	OrderValidationUtilsCaller     // Read-only binding to the contract
	OrderValidationUtilsTransactor // Write-only binding to the contract
	OrderValidationUtilsFilterer   // Log filterer for contract events
}

// OrderValidationUtilsCaller is an auto generated read-only Go binding around an Ethereum contract.
type OrderValidationUtilsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OrderValidationUtilsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OrderValidationUtilsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OrderValidationUtilsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OrderValidationUtilsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OrderValidationUtilsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OrderValidationUtilsSession struct {
	Contract     *OrderValidationUtils // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// OrderValidationUtilsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OrderValidationUtilsCallerSession struct {
	Contract *OrderValidationUtilsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// OrderValidationUtilsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OrderValidationUtilsTransactorSession struct {
	Contract     *OrderValidationUtilsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// OrderValidationUtilsRaw is an auto generated low-level Go binding around an Ethereum contract.
type OrderValidationUtilsRaw struct {
	Contract *OrderValidationUtils // Generic contract binding to access the raw methods on
}

// OrderValidationUtilsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OrderValidationUtilsCallerRaw struct {
	Contract *OrderValidationUtilsCaller // Generic read-only contract binding to access the raw methods on
}

// OrderValidationUtilsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OrderValidationUtilsTransactorRaw struct {
	Contract *OrderValidationUtilsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOrderValidationUtils creates a new instance of OrderValidationUtils, bound to a specific deployed contract.
func NewOrderValidationUtils(address common.Address, backend bind.ContractBackend) (*OrderValidationUtils, error) {
	contract, err := bindOrderValidationUtils(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OrderValidationUtils{OrderValidationUtilsCaller: OrderValidationUtilsCaller{contract: contract}, OrderValidationUtilsTransactor: OrderValidationUtilsTransactor{contract: contract}, OrderValidationUtilsFilterer: OrderValidationUtilsFilterer{contract: contract}}, nil
}

// NewOrderValidationUtilsCaller creates a new read-only instance of OrderValidationUtils, bound to a specific deployed contract.
func NewOrderValidationUtilsCaller(address common.Address, caller bind.ContractCaller) (*OrderValidationUtilsCaller, error) {
	contract, err := bindOrderValidationUtils(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OrderValidationUtilsCaller{contract: contract}, nil
}

// NewOrderValidationUtilsTransactor creates a new write-only instance of OrderValidationUtils, bound to a specific deployed contract.
func NewOrderValidationUtilsTransactor(address common.Address, transactor bind.ContractTransactor) (*OrderValidationUtilsTransactor, error) {
	contract, err := bindOrderValidationUtils(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OrderValidationUtilsTransactor{contract: contract}, nil
}

// NewOrderValidationUtilsFilterer creates a new log filterer instance of OrderValidationUtils, bound to a specific deployed contract.
func NewOrderValidationUtilsFilterer(address common.Address, filterer bind.ContractFilterer) (*OrderValidationUtilsFilterer, error) {
	contract, err := bindOrderValidationUtils(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OrderValidationUtilsFilterer{contract: contract}, nil
}

// bindOrderValidationUtils binds a generic wrapper to an already deployed contract.
func bindOrderValidationUtils(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OrderValidationUtilsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OrderValidationUtils *OrderValidationUtilsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OrderValidationUtils.Contract.OrderValidationUtilsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OrderValidationUtils *OrderValidationUtilsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OrderValidationUtils.Contract.OrderValidationUtilsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OrderValidationUtils *OrderValidationUtilsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OrderValidationUtils.Contract.OrderValidationUtilsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OrderValidationUtils *OrderValidationUtilsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OrderValidationUtils.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OrderValidationUtils *OrderValidationUtilsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OrderValidationUtils.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OrderValidationUtils *OrderValidationUtilsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OrderValidationUtils.Contract.contract.Transact(opts, method, params...)
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
}

// Struct1 is an auto generated low-level Go binding around an user-defined struct.
type Struct1 struct {
	OrderStatus                 uint8
	OrderHash                   [32]byte
	OrderTakerAssetFilledAmount *big.Int
}

// ERC1155PROXYID is a free data retrieval call binding the contract method 0x1bd0eb8f.
//
// Solidity: function ERC1155_PROXY_ID() constant returns(bytes4)
func (_OrderValidationUtils *OrderValidationUtilsCaller) ERC1155PROXYID(opts *bind.CallOpts) ([4]byte, error) {
	var (
		ret0 = new([4]byte)
	)
	out := ret0
	err := _OrderValidationUtils.contract.Call(opts, out, "ERC1155_PROXY_ID")
	return *ret0, err
}

// ERC1155PROXYID is a free data retrieval call binding the contract method 0x1bd0eb8f.
//
// Solidity: function ERC1155_PROXY_ID() constant returns(bytes4)
func (_OrderValidationUtils *OrderValidationUtilsSession) ERC1155PROXYID() ([4]byte, error) {
	return _OrderValidationUtils.Contract.ERC1155PROXYID(&_OrderValidationUtils.CallOpts)
}

// ERC1155PROXYID is a free data retrieval call binding the contract method 0x1bd0eb8f.
//
// Solidity: function ERC1155_PROXY_ID() constant returns(bytes4)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) ERC1155PROXYID() ([4]byte, error) {
	return _OrderValidationUtils.Contract.ERC1155PROXYID(&_OrderValidationUtils.CallOpts)
}

// ERC20PROXYID is a free data retrieval call binding the contract method 0x8ee1a642.
//
// Solidity: function ERC20_PROXY_ID() constant returns(bytes4)
func (_OrderValidationUtils *OrderValidationUtilsCaller) ERC20PROXYID(opts *bind.CallOpts) ([4]byte, error) {
	var (
		ret0 = new([4]byte)
	)
	out := ret0
	err := _OrderValidationUtils.contract.Call(opts, out, "ERC20_PROXY_ID")
	return *ret0, err
}

// ERC20PROXYID is a free data retrieval call binding the contract method 0x8ee1a642.
//
// Solidity: function ERC20_PROXY_ID() constant returns(bytes4)
func (_OrderValidationUtils *OrderValidationUtilsSession) ERC20PROXYID() ([4]byte, error) {
	return _OrderValidationUtils.Contract.ERC20PROXYID(&_OrderValidationUtils.CallOpts)
}

// ERC20PROXYID is a free data retrieval call binding the contract method 0x8ee1a642.
//
// Solidity: function ERC20_PROXY_ID() constant returns(bytes4)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) ERC20PROXYID() ([4]byte, error) {
	return _OrderValidationUtils.Contract.ERC20PROXYID(&_OrderValidationUtils.CallOpts)
}

// ERC721PROXYID is a free data retrieval call binding the contract method 0xa28fe02e.
//
// Solidity: function ERC721_PROXY_ID() constant returns(bytes4)
func (_OrderValidationUtils *OrderValidationUtilsCaller) ERC721PROXYID(opts *bind.CallOpts) ([4]byte, error) {
	var (
		ret0 = new([4]byte)
	)
	out := ret0
	err := _OrderValidationUtils.contract.Call(opts, out, "ERC721_PROXY_ID")
	return *ret0, err
}

// ERC721PROXYID is a free data retrieval call binding the contract method 0xa28fe02e.
//
// Solidity: function ERC721_PROXY_ID() constant returns(bytes4)
func (_OrderValidationUtils *OrderValidationUtilsSession) ERC721PROXYID() ([4]byte, error) {
	return _OrderValidationUtils.Contract.ERC721PROXYID(&_OrderValidationUtils.CallOpts)
}

// ERC721PROXYID is a free data retrieval call binding the contract method 0xa28fe02e.
//
// Solidity: function ERC721_PROXY_ID() constant returns(bytes4)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) ERC721PROXYID() ([4]byte, error) {
	return _OrderValidationUtils.Contract.ERC721PROXYID(&_OrderValidationUtils.CallOpts)
}

// MULTIASSETPROXYID is a free data retrieval call binding the contract method 0xb37fda04.
//
// Solidity: function MULTI_ASSET_PROXY_ID() constant returns(bytes4)
func (_OrderValidationUtils *OrderValidationUtilsCaller) MULTIASSETPROXYID(opts *bind.CallOpts) ([4]byte, error) {
	var (
		ret0 = new([4]byte)
	)
	out := ret0
	err := _OrderValidationUtils.contract.Call(opts, out, "MULTI_ASSET_PROXY_ID")
	return *ret0, err
}

// MULTIASSETPROXYID is a free data retrieval call binding the contract method 0xb37fda04.
//
// Solidity: function MULTI_ASSET_PROXY_ID() constant returns(bytes4)
func (_OrderValidationUtils *OrderValidationUtilsSession) MULTIASSETPROXYID() ([4]byte, error) {
	return _OrderValidationUtils.Contract.MULTIASSETPROXYID(&_OrderValidationUtils.CallOpts)
}

// MULTIASSETPROXYID is a free data retrieval call binding the contract method 0xb37fda04.
//
// Solidity: function MULTI_ASSET_PROXY_ID() constant returns(bytes4)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) MULTIASSETPROXYID() ([4]byte, error) {
	return _OrderValidationUtils.Contract.MULTIASSETPROXYID(&_OrderValidationUtils.CallOpts)
}

// DecodeERC1155AssetData is a free data retrieval call binding the contract method 0x9eadc835.
//
// Solidity: function decodeERC1155AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData)
func (_OrderValidationUtils *OrderValidationUtilsCaller) DecodeERC1155AssetData(opts *bind.CallOpts, assetData []byte) (struct {
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
	err := _OrderValidationUtils.contract.Call(opts, out, "decodeERC1155AssetData", assetData)
	return *ret, err
}

// DecodeERC1155AssetData is a free data retrieval call binding the contract method 0x9eadc835.
//
// Solidity: function decodeERC1155AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData)
func (_OrderValidationUtils *OrderValidationUtilsSession) DecodeERC1155AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
	TokenIds     []*big.Int
	TokenValues  []*big.Int
	CallbackData []byte
}, error) {
	return _OrderValidationUtils.Contract.DecodeERC1155AssetData(&_OrderValidationUtils.CallOpts, assetData)
}

// DecodeERC1155AssetData is a free data retrieval call binding the contract method 0x9eadc835.
//
// Solidity: function decodeERC1155AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) DecodeERC1155AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
	TokenIds     []*big.Int
	TokenValues  []*big.Int
	CallbackData []byte
}, error) {
	return _OrderValidationUtils.Contract.DecodeERC1155AssetData(&_OrderValidationUtils.CallOpts, assetData)
}

// DecodeERC20AssetData is a free data retrieval call binding the contract method 0x8f4ce479.
//
// Solidity: function decodeERC20AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress)
func (_OrderValidationUtils *OrderValidationUtilsCaller) DecodeERC20AssetData(opts *bind.CallOpts, assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
}, error) {
	ret := new(struct {
		AssetProxyId [4]byte
		TokenAddress common.Address
	})
	out := ret
	err := _OrderValidationUtils.contract.Call(opts, out, "decodeERC20AssetData", assetData)
	return *ret, err
}

// DecodeERC20AssetData is a free data retrieval call binding the contract method 0x8f4ce479.
//
// Solidity: function decodeERC20AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress)
func (_OrderValidationUtils *OrderValidationUtilsSession) DecodeERC20AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
}, error) {
	return _OrderValidationUtils.Contract.DecodeERC20AssetData(&_OrderValidationUtils.CallOpts, assetData)
}

// DecodeERC20AssetData is a free data retrieval call binding the contract method 0x8f4ce479.
//
// Solidity: function decodeERC20AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) DecodeERC20AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
}, error) {
	return _OrderValidationUtils.Contract.DecodeERC20AssetData(&_OrderValidationUtils.CallOpts, assetData)
}

// DecodeERC721AssetData is a free data retrieval call binding the contract method 0x04a5618a.
//
// Solidity: function decodeERC721AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress, uint256 tokenId)
func (_OrderValidationUtils *OrderValidationUtilsCaller) DecodeERC721AssetData(opts *bind.CallOpts, assetData []byte) (struct {
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
	err := _OrderValidationUtils.contract.Call(opts, out, "decodeERC721AssetData", assetData)
	return *ret, err
}

// DecodeERC721AssetData is a free data retrieval call binding the contract method 0x04a5618a.
//
// Solidity: function decodeERC721AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress, uint256 tokenId)
func (_OrderValidationUtils *OrderValidationUtilsSession) DecodeERC721AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
	TokenId      *big.Int
}, error) {
	return _OrderValidationUtils.Contract.DecodeERC721AssetData(&_OrderValidationUtils.CallOpts, assetData)
}

// DecodeERC721AssetData is a free data retrieval call binding the contract method 0x04a5618a.
//
// Solidity: function decodeERC721AssetData(bytes assetData) constant returns(bytes4 assetProxyId, address tokenAddress, uint256 tokenId)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) DecodeERC721AssetData(assetData []byte) (struct {
	AssetProxyId [4]byte
	TokenAddress common.Address
	TokenId      *big.Int
}, error) {
	return _OrderValidationUtils.Contract.DecodeERC721AssetData(&_OrderValidationUtils.CallOpts, assetData)
}

// DecodeMultiAssetData is a free data retrieval call binding the contract method 0xbbb2dcf6.
//
// Solidity: function decodeMultiAssetData(bytes assetData) constant returns(bytes4 assetProxyId, uint256[] amounts, bytes[] nestedAssetData)
func (_OrderValidationUtils *OrderValidationUtilsCaller) DecodeMultiAssetData(opts *bind.CallOpts, assetData []byte) (struct {
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
	err := _OrderValidationUtils.contract.Call(opts, out, "decodeMultiAssetData", assetData)
	return *ret, err
}

// DecodeMultiAssetData is a free data retrieval call binding the contract method 0xbbb2dcf6.
//
// Solidity: function decodeMultiAssetData(bytes assetData) constant returns(bytes4 assetProxyId, uint256[] amounts, bytes[] nestedAssetData)
func (_OrderValidationUtils *OrderValidationUtilsSession) DecodeMultiAssetData(assetData []byte) (struct {
	AssetProxyId    [4]byte
	Amounts         []*big.Int
	NestedAssetData [][]byte
}, error) {
	return _OrderValidationUtils.Contract.DecodeMultiAssetData(&_OrderValidationUtils.CallOpts, assetData)
}

// DecodeMultiAssetData is a free data retrieval call binding the contract method 0xbbb2dcf6.
//
// Solidity: function decodeMultiAssetData(bytes assetData) constant returns(bytes4 assetProxyId, uint256[] amounts, bytes[] nestedAssetData)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) DecodeMultiAssetData(assetData []byte) (struct {
	AssetProxyId    [4]byte
	Amounts         []*big.Int
	NestedAssetData [][]byte
}, error) {
	return _OrderValidationUtils.Contract.DecodeMultiAssetData(&_OrderValidationUtils.CallOpts, assetData)
}

// EncodeERC1155AssetData is a free data retrieval call binding the contract method 0xb43cffe1.
//
// Solidity: function encodeERC1155AssetData(address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData) constant returns(bytes assetData)
func (_OrderValidationUtils *OrderValidationUtilsCaller) EncodeERC1155AssetData(opts *bind.CallOpts, tokenAddress common.Address, tokenIds []*big.Int, tokenValues []*big.Int, callbackData []byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _OrderValidationUtils.contract.Call(opts, out, "encodeERC1155AssetData", tokenAddress, tokenIds, tokenValues, callbackData)
	return *ret0, err
}

// EncodeERC1155AssetData is a free data retrieval call binding the contract method 0xb43cffe1.
//
// Solidity: function encodeERC1155AssetData(address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData) constant returns(bytes assetData)
func (_OrderValidationUtils *OrderValidationUtilsSession) EncodeERC1155AssetData(tokenAddress common.Address, tokenIds []*big.Int, tokenValues []*big.Int, callbackData []byte) ([]byte, error) {
	return _OrderValidationUtils.Contract.EncodeERC1155AssetData(&_OrderValidationUtils.CallOpts, tokenAddress, tokenIds, tokenValues, callbackData)
}

// EncodeERC1155AssetData is a free data retrieval call binding the contract method 0xb43cffe1.
//
// Solidity: function encodeERC1155AssetData(address tokenAddress, uint256[] tokenIds, uint256[] tokenValues, bytes callbackData) constant returns(bytes assetData)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) EncodeERC1155AssetData(tokenAddress common.Address, tokenIds []*big.Int, tokenValues []*big.Int, callbackData []byte) ([]byte, error) {
	return _OrderValidationUtils.Contract.EncodeERC1155AssetData(&_OrderValidationUtils.CallOpts, tokenAddress, tokenIds, tokenValues, callbackData)
}

// EncodeERC20AssetData is a free data retrieval call binding the contract method 0x590aa875.
//
// Solidity: function encodeERC20AssetData(address tokenAddress) constant returns(bytes assetData)
func (_OrderValidationUtils *OrderValidationUtilsCaller) EncodeERC20AssetData(opts *bind.CallOpts, tokenAddress common.Address) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _OrderValidationUtils.contract.Call(opts, out, "encodeERC20AssetData", tokenAddress)
	return *ret0, err
}

// EncodeERC20AssetData is a free data retrieval call binding the contract method 0x590aa875.
//
// Solidity: function encodeERC20AssetData(address tokenAddress) constant returns(bytes assetData)
func (_OrderValidationUtils *OrderValidationUtilsSession) EncodeERC20AssetData(tokenAddress common.Address) ([]byte, error) {
	return _OrderValidationUtils.Contract.EncodeERC20AssetData(&_OrderValidationUtils.CallOpts, tokenAddress)
}

// EncodeERC20AssetData is a free data retrieval call binding the contract method 0x590aa875.
//
// Solidity: function encodeERC20AssetData(address tokenAddress) constant returns(bytes assetData)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) EncodeERC20AssetData(tokenAddress common.Address) ([]byte, error) {
	return _OrderValidationUtils.Contract.EncodeERC20AssetData(&_OrderValidationUtils.CallOpts, tokenAddress)
}

// EncodeERC721AssetData is a free data retrieval call binding the contract method 0xa6627e9f.
//
// Solidity: function encodeERC721AssetData(address tokenAddress, uint256 tokenId) constant returns(bytes assetData)
func (_OrderValidationUtils *OrderValidationUtilsCaller) EncodeERC721AssetData(opts *bind.CallOpts, tokenAddress common.Address, tokenId *big.Int) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _OrderValidationUtils.contract.Call(opts, out, "encodeERC721AssetData", tokenAddress, tokenId)
	return *ret0, err
}

// EncodeERC721AssetData is a free data retrieval call binding the contract method 0xa6627e9f.
//
// Solidity: function encodeERC721AssetData(address tokenAddress, uint256 tokenId) constant returns(bytes assetData)
func (_OrderValidationUtils *OrderValidationUtilsSession) EncodeERC721AssetData(tokenAddress common.Address, tokenId *big.Int) ([]byte, error) {
	return _OrderValidationUtils.Contract.EncodeERC721AssetData(&_OrderValidationUtils.CallOpts, tokenAddress, tokenId)
}

// EncodeERC721AssetData is a free data retrieval call binding the contract method 0xa6627e9f.
//
// Solidity: function encodeERC721AssetData(address tokenAddress, uint256 tokenId) constant returns(bytes assetData)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) EncodeERC721AssetData(tokenAddress common.Address, tokenId *big.Int) ([]byte, error) {
	return _OrderValidationUtils.Contract.EncodeERC721AssetData(&_OrderValidationUtils.CallOpts, tokenAddress, tokenId)
}

// EncodeMultiAssetData is a free data retrieval call binding the contract method 0xd3d862d1.
//
// Solidity: function encodeMultiAssetData(uint256[] amounts, bytes[] nestedAssetData) constant returns(bytes assetData)
func (_OrderValidationUtils *OrderValidationUtilsCaller) EncodeMultiAssetData(opts *bind.CallOpts, amounts []*big.Int, nestedAssetData [][]byte) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _OrderValidationUtils.contract.Call(opts, out, "encodeMultiAssetData", amounts, nestedAssetData)
	return *ret0, err
}

// EncodeMultiAssetData is a free data retrieval call binding the contract method 0xd3d862d1.
//
// Solidity: function encodeMultiAssetData(uint256[] amounts, bytes[] nestedAssetData) constant returns(bytes assetData)
func (_OrderValidationUtils *OrderValidationUtilsSession) EncodeMultiAssetData(amounts []*big.Int, nestedAssetData [][]byte) ([]byte, error) {
	return _OrderValidationUtils.Contract.EncodeMultiAssetData(&_OrderValidationUtils.CallOpts, amounts, nestedAssetData)
}

// EncodeMultiAssetData is a free data retrieval call binding the contract method 0xd3d862d1.
//
// Solidity: function encodeMultiAssetData(uint256[] amounts, bytes[] nestedAssetData) constant returns(bytes assetData)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) EncodeMultiAssetData(amounts []*big.Int, nestedAssetData [][]byte) ([]byte, error) {
	return _OrderValidationUtils.Contract.EncodeMultiAssetData(&_OrderValidationUtils.CallOpts, amounts, nestedAssetData)
}

// GetAssetProxyAllowance is a free data retrieval call binding the contract method 0xd186037f.
//
// Solidity: function getAssetProxyAllowance(address ownerAddress, bytes assetData) constant returns(uint256 allowance)
func (_OrderValidationUtils *OrderValidationUtilsCaller) GetAssetProxyAllowance(opts *bind.CallOpts, ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OrderValidationUtils.contract.Call(opts, out, "getAssetProxyAllowance", ownerAddress, assetData)
	return *ret0, err
}

// GetAssetProxyAllowance is a free data retrieval call binding the contract method 0xd186037f.
//
// Solidity: function getAssetProxyAllowance(address ownerAddress, bytes assetData) constant returns(uint256 allowance)
func (_OrderValidationUtils *OrderValidationUtilsSession) GetAssetProxyAllowance(ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	return _OrderValidationUtils.Contract.GetAssetProxyAllowance(&_OrderValidationUtils.CallOpts, ownerAddress, assetData)
}

// GetAssetProxyAllowance is a free data retrieval call binding the contract method 0xd186037f.
//
// Solidity: function getAssetProxyAllowance(address ownerAddress, bytes assetData) constant returns(uint256 allowance)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) GetAssetProxyAllowance(ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	return _OrderValidationUtils.Contract.GetAssetProxyAllowance(&_OrderValidationUtils.CallOpts, ownerAddress, assetData)
}

// GetBalance is a free data retrieval call binding the contract method 0x7d727512.
//
// Solidity: function getBalance(address ownerAddress, bytes assetData) constant returns(uint256 balance)
func (_OrderValidationUtils *OrderValidationUtilsCaller) GetBalance(opts *bind.CallOpts, ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OrderValidationUtils.contract.Call(opts, out, "getBalance", ownerAddress, assetData)
	return *ret0, err
}

// GetBalance is a free data retrieval call binding the contract method 0x7d727512.
//
// Solidity: function getBalance(address ownerAddress, bytes assetData) constant returns(uint256 balance)
func (_OrderValidationUtils *OrderValidationUtilsSession) GetBalance(ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	return _OrderValidationUtils.Contract.GetBalance(&_OrderValidationUtils.CallOpts, ownerAddress, assetData)
}

// GetBalance is a free data retrieval call binding the contract method 0x7d727512.
//
// Solidity: function getBalance(address ownerAddress, bytes assetData) constant returns(uint256 balance)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) GetBalance(ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	return _OrderValidationUtils.Contract.GetBalance(&_OrderValidationUtils.CallOpts, ownerAddress, assetData)
}

// GetBalanceAndAssetProxyAllowance is a free data retrieval call binding the contract method 0x0d7b7d76.
//
// Solidity: function getBalanceAndAssetProxyAllowance(address ownerAddress, bytes assetData) constant returns(uint256 balance, uint256 allowance)
func (_OrderValidationUtils *OrderValidationUtilsCaller) GetBalanceAndAssetProxyAllowance(opts *bind.CallOpts, ownerAddress common.Address, assetData []byte) (struct {
	Balance   *big.Int
	Allowance *big.Int
}, error) {
	ret := new(struct {
		Balance   *big.Int
		Allowance *big.Int
	})
	out := ret
	err := _OrderValidationUtils.contract.Call(opts, out, "getBalanceAndAssetProxyAllowance", ownerAddress, assetData)
	return *ret, err
}

// GetBalanceAndAssetProxyAllowance is a free data retrieval call binding the contract method 0x0d7b7d76.
//
// Solidity: function getBalanceAndAssetProxyAllowance(address ownerAddress, bytes assetData) constant returns(uint256 balance, uint256 allowance)
func (_OrderValidationUtils *OrderValidationUtilsSession) GetBalanceAndAssetProxyAllowance(ownerAddress common.Address, assetData []byte) (struct {
	Balance   *big.Int
	Allowance *big.Int
}, error) {
	return _OrderValidationUtils.Contract.GetBalanceAndAssetProxyAllowance(&_OrderValidationUtils.CallOpts, ownerAddress, assetData)
}

// GetBalanceAndAssetProxyAllowance is a free data retrieval call binding the contract method 0x0d7b7d76.
//
// Solidity: function getBalanceAndAssetProxyAllowance(address ownerAddress, bytes assetData) constant returns(uint256 balance, uint256 allowance)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) GetBalanceAndAssetProxyAllowance(ownerAddress common.Address, assetData []byte) (struct {
	Balance   *big.Int
	Allowance *big.Int
}, error) {
	return _OrderValidationUtils.Contract.GetBalanceAndAssetProxyAllowance(&_OrderValidationUtils.CallOpts, ownerAddress, assetData)
}

// GetBatchAssetProxyAllowances is a free data retrieval call binding the contract method 0x4dfdac20.
//
// Solidity: function getBatchAssetProxyAllowances(address ownerAddress, bytes[] assetData) constant returns(uint256[] allowances)
func (_OrderValidationUtils *OrderValidationUtilsCaller) GetBatchAssetProxyAllowances(opts *bind.CallOpts, ownerAddress common.Address, assetData [][]byte) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _OrderValidationUtils.contract.Call(opts, out, "getBatchAssetProxyAllowances", ownerAddress, assetData)
	return *ret0, err
}

// GetBatchAssetProxyAllowances is a free data retrieval call binding the contract method 0x4dfdac20.
//
// Solidity: function getBatchAssetProxyAllowances(address ownerAddress, bytes[] assetData) constant returns(uint256[] allowances)
func (_OrderValidationUtils *OrderValidationUtilsSession) GetBatchAssetProxyAllowances(ownerAddress common.Address, assetData [][]byte) ([]*big.Int, error) {
	return _OrderValidationUtils.Contract.GetBatchAssetProxyAllowances(&_OrderValidationUtils.CallOpts, ownerAddress, assetData)
}

// GetBatchAssetProxyAllowances is a free data retrieval call binding the contract method 0x4dfdac20.
//
// Solidity: function getBatchAssetProxyAllowances(address ownerAddress, bytes[] assetData) constant returns(uint256[] allowances)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) GetBatchAssetProxyAllowances(ownerAddress common.Address, assetData [][]byte) ([]*big.Int, error) {
	return _OrderValidationUtils.Contract.GetBatchAssetProxyAllowances(&_OrderValidationUtils.CallOpts, ownerAddress, assetData)
}

// GetBatchBalances is a free data retrieval call binding the contract method 0xd001c5dc.
//
// Solidity: function getBatchBalances(address ownerAddress, bytes[] assetData) constant returns(uint256[] balances)
func (_OrderValidationUtils *OrderValidationUtilsCaller) GetBatchBalances(opts *bind.CallOpts, ownerAddress common.Address, assetData [][]byte) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _OrderValidationUtils.contract.Call(opts, out, "getBatchBalances", ownerAddress, assetData)
	return *ret0, err
}

// GetBatchBalances is a free data retrieval call binding the contract method 0xd001c5dc.
//
// Solidity: function getBatchBalances(address ownerAddress, bytes[] assetData) constant returns(uint256[] balances)
func (_OrderValidationUtils *OrderValidationUtilsSession) GetBatchBalances(ownerAddress common.Address, assetData [][]byte) ([]*big.Int, error) {
	return _OrderValidationUtils.Contract.GetBatchBalances(&_OrderValidationUtils.CallOpts, ownerAddress, assetData)
}

// GetBatchBalances is a free data retrieval call binding the contract method 0xd001c5dc.
//
// Solidity: function getBatchBalances(address ownerAddress, bytes[] assetData) constant returns(uint256[] balances)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) GetBatchBalances(ownerAddress common.Address, assetData [][]byte) ([]*big.Int, error) {
	return _OrderValidationUtils.Contract.GetBatchBalances(&_OrderValidationUtils.CallOpts, ownerAddress, assetData)
}

// GetBatchBalancesAndAssetProxyAllowances is a free data retrieval call binding the contract method 0xe4e6e7da.
//
// Solidity: function getBatchBalancesAndAssetProxyAllowances(address ownerAddress, bytes[] assetData) constant returns(uint256[] balances, uint256[] allowances)
func (_OrderValidationUtils *OrderValidationUtilsCaller) GetBatchBalancesAndAssetProxyAllowances(opts *bind.CallOpts, ownerAddress common.Address, assetData [][]byte) (struct {
	Balances   []*big.Int
	Allowances []*big.Int
}, error) {
	ret := new(struct {
		Balances   []*big.Int
		Allowances []*big.Int
	})
	out := ret
	err := _OrderValidationUtils.contract.Call(opts, out, "getBatchBalancesAndAssetProxyAllowances", ownerAddress, assetData)
	return *ret, err
}

// GetBatchBalancesAndAssetProxyAllowances is a free data retrieval call binding the contract method 0xe4e6e7da.
//
// Solidity: function getBatchBalancesAndAssetProxyAllowances(address ownerAddress, bytes[] assetData) constant returns(uint256[] balances, uint256[] allowances)
func (_OrderValidationUtils *OrderValidationUtilsSession) GetBatchBalancesAndAssetProxyAllowances(ownerAddress common.Address, assetData [][]byte) (struct {
	Balances   []*big.Int
	Allowances []*big.Int
}, error) {
	return _OrderValidationUtils.Contract.GetBatchBalancesAndAssetProxyAllowances(&_OrderValidationUtils.CallOpts, ownerAddress, assetData)
}

// GetBatchBalancesAndAssetProxyAllowances is a free data retrieval call binding the contract method 0xe4e6e7da.
//
// Solidity: function getBatchBalancesAndAssetProxyAllowances(address ownerAddress, bytes[] assetData) constant returns(uint256[] balances, uint256[] allowances)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) GetBatchBalancesAndAssetProxyAllowances(ownerAddress common.Address, assetData [][]byte) (struct {
	Balances   []*big.Int
	Allowances []*big.Int
}, error) {
	return _OrderValidationUtils.Contract.GetBatchBalancesAndAssetProxyAllowances(&_OrderValidationUtils.CallOpts, ownerAddress, assetData)
}

// GetERC721TokenOwner is a free data retrieval call binding the contract method 0xb6988463.
//
// Solidity: function getERC721TokenOwner(address tokenAddress, uint256 tokenId) constant returns(address ownerAddress)
func (_OrderValidationUtils *OrderValidationUtilsCaller) GetERC721TokenOwner(opts *bind.CallOpts, tokenAddress common.Address, tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _OrderValidationUtils.contract.Call(opts, out, "getERC721TokenOwner", tokenAddress, tokenId)
	return *ret0, err
}

// GetERC721TokenOwner is a free data retrieval call binding the contract method 0xb6988463.
//
// Solidity: function getERC721TokenOwner(address tokenAddress, uint256 tokenId) constant returns(address ownerAddress)
func (_OrderValidationUtils *OrderValidationUtilsSession) GetERC721TokenOwner(tokenAddress common.Address, tokenId *big.Int) (common.Address, error) {
	return _OrderValidationUtils.Contract.GetERC721TokenOwner(&_OrderValidationUtils.CallOpts, tokenAddress, tokenId)
}

// GetERC721TokenOwner is a free data retrieval call binding the contract method 0xb6988463.
//
// Solidity: function getERC721TokenOwner(address tokenAddress, uint256 tokenId) constant returns(address ownerAddress)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) GetERC721TokenOwner(tokenAddress common.Address, tokenId *big.Int) (common.Address, error) {
	return _OrderValidationUtils.Contract.GetERC721TokenOwner(&_OrderValidationUtils.CallOpts, tokenAddress, tokenId)
}

// GetOrderRelevantState is a free data retrieval call binding the contract method 0x8f5afa52.
//
// Solidity: function getOrderRelevantState(Struct0 order, bytes signature) constant returns(Struct1 orderInfo, uint256 fillableTakerAssetAmount, bool isValidSignature)
func (_OrderValidationUtils *OrderValidationUtilsCaller) GetOrderRelevantState(opts *bind.CallOpts, order Struct0, signature []byte) (struct {
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
	err := _OrderValidationUtils.contract.Call(opts, out, "getOrderRelevantState", order, signature)
	return *ret, err
}

// GetOrderRelevantState is a free data retrieval call binding the contract method 0x8f5afa52.
//
// Solidity: function getOrderRelevantState(Struct0 order, bytes signature) constant returns(Struct1 orderInfo, uint256 fillableTakerAssetAmount, bool isValidSignature)
func (_OrderValidationUtils *OrderValidationUtilsSession) GetOrderRelevantState(order Struct0, signature []byte) (struct {
	OrderInfo                Struct1
	FillableTakerAssetAmount *big.Int
	IsValidSignature         bool
}, error) {
	return _OrderValidationUtils.Contract.GetOrderRelevantState(&_OrderValidationUtils.CallOpts, order, signature)
}

// GetOrderRelevantState is a free data retrieval call binding the contract method 0x8f5afa52.
//
// Solidity: function getOrderRelevantState(Struct0 order, bytes signature) constant returns(Struct1 orderInfo, uint256 fillableTakerAssetAmount, bool isValidSignature)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) GetOrderRelevantState(order Struct0, signature []byte) (struct {
	OrderInfo                Struct1
	FillableTakerAssetAmount *big.Int
	IsValidSignature         bool
}, error) {
	return _OrderValidationUtils.Contract.GetOrderRelevantState(&_OrderValidationUtils.CallOpts, order, signature)
}

// GetOrderRelevantStates is a free data retrieval call binding the contract method 0x7f46448d.
//
// Solidity: function getOrderRelevantStates((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes)[] orders, bytes[] signatures) constant returns((uint8,bytes32,uint256)[] ordersInfo, uint256[] fillableTakerAssetAmounts, bool[] isValidSignature)
func (_OrderValidationUtils *OrderValidationUtilsCaller) GetOrderRelevantStates(opts *bind.CallOpts, orders []Struct0, signatures [][]byte) (struct {
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
	err := _OrderValidationUtils.contract.Call(opts, out, "getOrderRelevantStates", orders, signatures)
	return *ret, err
}

// GetOrderRelevantStates is a free data retrieval call binding the contract method 0x7f46448d.
//
// Solidity: function getOrderRelevantStates((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes)[] orders, bytes[] signatures) constant returns((uint8,bytes32,uint256)[] ordersInfo, uint256[] fillableTakerAssetAmounts, bool[] isValidSignature)
func (_OrderValidationUtils *OrderValidationUtilsSession) GetOrderRelevantStates(orders []Struct0, signatures [][]byte) (struct {
	OrdersInfo                []Struct1
	FillableTakerAssetAmounts []*big.Int
	IsValidSignature          []bool
}, error) {
	return _OrderValidationUtils.Contract.GetOrderRelevantStates(&_OrderValidationUtils.CallOpts, orders, signatures)
}

// GetOrderRelevantStates is a free data retrieval call binding the contract method 0x7f46448d.
//
// Solidity: function getOrderRelevantStates((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes)[] orders, bytes[] signatures) constant returns((uint8,bytes32,uint256)[] ordersInfo, uint256[] fillableTakerAssetAmounts, bool[] isValidSignature)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) GetOrderRelevantStates(orders []Struct0, signatures [][]byte) (struct {
	OrdersInfo                []Struct1
	FillableTakerAssetAmounts []*big.Int
	IsValidSignature          []bool
}, error) {
	return _OrderValidationUtils.Contract.GetOrderRelevantStates(&_OrderValidationUtils.CallOpts, orders, signatures)
}

// GetTransferableAssetAmount is a free data retrieval call binding the contract method 0x2322cf76.
//
// Solidity: function getTransferableAssetAmount(address ownerAddress, bytes assetData) constant returns(uint256 transferableAssetAmount)
func (_OrderValidationUtils *OrderValidationUtilsCaller) GetTransferableAssetAmount(opts *bind.CallOpts, ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OrderValidationUtils.contract.Call(opts, out, "getTransferableAssetAmount", ownerAddress, assetData)
	return *ret0, err
}

// GetTransferableAssetAmount is a free data retrieval call binding the contract method 0x2322cf76.
//
// Solidity: function getTransferableAssetAmount(address ownerAddress, bytes assetData) constant returns(uint256 transferableAssetAmount)
func (_OrderValidationUtils *OrderValidationUtilsSession) GetTransferableAssetAmount(ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	return _OrderValidationUtils.Contract.GetTransferableAssetAmount(&_OrderValidationUtils.CallOpts, ownerAddress, assetData)
}

// GetTransferableAssetAmount is a free data retrieval call binding the contract method 0x2322cf76.
//
// Solidity: function getTransferableAssetAmount(address ownerAddress, bytes assetData) constant returns(uint256 transferableAssetAmount)
func (_OrderValidationUtils *OrderValidationUtilsCallerSession) GetTransferableAssetAmount(ownerAddress common.Address, assetData []byte) (*big.Int, error) {
	return _OrderValidationUtils.Contract.GetTransferableAssetAmount(&_OrderValidationUtils.CallOpts, ownerAddress, assetData)
}

