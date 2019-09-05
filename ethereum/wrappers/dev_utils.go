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
const DevUtilsABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC721AssetData\",\"outputs\":[{\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"name\":\"tokenAddress\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"ownerAddress\",\"type\":\"address\"},{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getBalanceAndAssetProxyAllowance\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"},{\"name\":\"allowance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ERC1155_PROXY_ID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"ownerAddress\",\"type\":\"address\"},{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getTransferableAssetAmount\",\"outputs\":[{\"name\":\"transferableAssetAmount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"ownerAddress\",\"type\":\"address\"},{\"name\":\"assetData\",\"type\":\"bytes[]\"}],\"name\":\"getBatchAssetProxyAllowances\",\"outputs\":[{\"name\":\"allowances\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\"}],\"name\":\"encodeERC20AssetData\",\"outputs\":[{\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"transactionData\",\"type\":\"bytes\"}],\"name\":\"decodeZeroExTransactionData\",\"outputs\":[{\"name\":\"functionName\",\"type\":\"string\"},{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"takerAssetFillAmounts\",\"type\":\"uint256[]\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"ownerAddress\",\"type\":\"address\"},{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getBalance\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"}],\"name\":\"getOrderRelevantStates\",\"outputs\":[{\"components\":[{\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"name\":\"ordersInfo\",\"type\":\"tuple[]\"},{\"name\":\"fillableTakerAssetAmounts\",\"type\":\"uint256[]\"},{\"name\":\"isValidSignature\",\"type\":\"bool[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ERC20_PROXY_ID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC20AssetData\",\"outputs\":[{\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"name\":\"tokenAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"},{\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"getOrderRelevantState\",\"outputs\":[{\"components\":[{\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"name\":\"orderInfo\",\"type\":\"tuple\"},{\"name\":\"fillableTakerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"isValidSignature\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeERC1155AssetData\",\"outputs\":[{\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"name\":\"tokenAddress\",\"type\":\"address\"},{\"name\":\"tokenIds\",\"type\":\"uint256[]\"},{\"name\":\"tokenValues\",\"type\":\"uint256[]\"},{\"name\":\"callbackData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addresses\",\"type\":\"address[]\"}],\"name\":\"getEthBalances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ERC721_PROXY_ID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"encodeERC721AssetData\",\"outputs\":[{\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MULTI_ASSET_PROXY_ID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\"},{\"name\":\"tokenIds\",\"type\":\"uint256[]\"},{\"name\":\"tokenValues\",\"type\":\"uint256[]\"},{\"name\":\"callbackData\",\"type\":\"bytes\"}],\"name\":\"encodeERC1155AssetData\",\"outputs\":[{\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenAddress\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getERC721TokenOwner\",\"outputs\":[{\"name\":\"ownerAddress\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"decodeMultiAssetData\",\"outputs\":[{\"name\":\"assetProxyId\",\"type\":\"bytes4\"},{\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"name\":\"nestedAssetData\",\"type\":\"bytes[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"ownerAddress\",\"type\":\"address\"},{\"name\":\"assetData\",\"type\":\"bytes[]\"}],\"name\":\"getBatchBalances\",\"outputs\":[{\"name\":\"balances\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"ownerAddress\",\"type\":\"address\"},{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getAssetProxyAllowance\",\"outputs\":[{\"name\":\"allowance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"name\":\"nestedAssetData\",\"type\":\"bytes[]\"}],\"name\":\"encodeMultiAssetData\",\"outputs\":[{\"name\":\"assetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"STATIC_CALL_PROXY_ID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"ownerAddress\",\"type\":\"address\"},{\"name\":\"assetData\",\"type\":\"bytes[]\"}],\"name\":\"getBatchBalancesAndAssetProxyAllowances\",\"outputs\":[{\"name\":\"balances\",\"type\":\"uint256[]\"},{\"name\":\"allowances\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_exchange\",\"type\":\"address\"},{\"name\":\"_zrxAssetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

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
func (_DevUtils *DevUtilsCaller) ERC1155PROXYID(opts *bind.CallOpts) ([4]byte, error) {
	var (
		ret0 = new([4]byte)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "ERC1155_PROXY_ID")
	return *ret0, err
}

// ERC1155PROXYID is a free data retrieval call binding the contract method 0x1bd0eb8f.
//
// Solidity: function ERC1155_PROXY_ID() constant returns(bytes4)
func (_DevUtils *DevUtilsSession) ERC1155PROXYID() ([4]byte, error) {
	return _DevUtils.Contract.ERC1155PROXYID(&_DevUtils.CallOpts)
}

// ERC1155PROXYID is a free data retrieval call binding the contract method 0x1bd0eb8f.
//
// Solidity: function ERC1155_PROXY_ID() constant returns(bytes4)
func (_DevUtils *DevUtilsCallerSession) ERC1155PROXYID() ([4]byte, error) {
	return _DevUtils.Contract.ERC1155PROXYID(&_DevUtils.CallOpts)
}

// ERC20PROXYID is a free data retrieval call binding the contract method 0x8ee1a642.
//
// Solidity: function ERC20_PROXY_ID() constant returns(bytes4)
func (_DevUtils *DevUtilsCaller) ERC20PROXYID(opts *bind.CallOpts) ([4]byte, error) {
	var (
		ret0 = new([4]byte)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "ERC20_PROXY_ID")
	return *ret0, err
}

// ERC20PROXYID is a free data retrieval call binding the contract method 0x8ee1a642.
//
// Solidity: function ERC20_PROXY_ID() constant returns(bytes4)
func (_DevUtils *DevUtilsSession) ERC20PROXYID() ([4]byte, error) {
	return _DevUtils.Contract.ERC20PROXYID(&_DevUtils.CallOpts)
}

// ERC20PROXYID is a free data retrieval call binding the contract method 0x8ee1a642.
//
// Solidity: function ERC20_PROXY_ID() constant returns(bytes4)
func (_DevUtils *DevUtilsCallerSession) ERC20PROXYID() ([4]byte, error) {
	return _DevUtils.Contract.ERC20PROXYID(&_DevUtils.CallOpts)
}

// ERC721PROXYID is a free data retrieval call binding the contract method 0xa28fe02e.
//
// Solidity: function ERC721_PROXY_ID() constant returns(bytes4)
func (_DevUtils *DevUtilsCaller) ERC721PROXYID(opts *bind.CallOpts) ([4]byte, error) {
	var (
		ret0 = new([4]byte)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "ERC721_PROXY_ID")
	return *ret0, err
}

// ERC721PROXYID is a free data retrieval call binding the contract method 0xa28fe02e.
//
// Solidity: function ERC721_PROXY_ID() constant returns(bytes4)
func (_DevUtils *DevUtilsSession) ERC721PROXYID() ([4]byte, error) {
	return _DevUtils.Contract.ERC721PROXYID(&_DevUtils.CallOpts)
}

// ERC721PROXYID is a free data retrieval call binding the contract method 0xa28fe02e.
//
// Solidity: function ERC721_PROXY_ID() constant returns(bytes4)
func (_DevUtils *DevUtilsCallerSession) ERC721PROXYID() ([4]byte, error) {
	return _DevUtils.Contract.ERC721PROXYID(&_DevUtils.CallOpts)
}

// MULTIASSETPROXYID is a free data retrieval call binding the contract method 0xb37fda04.
//
// Solidity: function MULTI_ASSET_PROXY_ID() constant returns(bytes4)
func (_DevUtils *DevUtilsCaller) MULTIASSETPROXYID(opts *bind.CallOpts) ([4]byte, error) {
	var (
		ret0 = new([4]byte)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "MULTI_ASSET_PROXY_ID")
	return *ret0, err
}

// MULTIASSETPROXYID is a free data retrieval call binding the contract method 0xb37fda04.
//
// Solidity: function MULTI_ASSET_PROXY_ID() constant returns(bytes4)
func (_DevUtils *DevUtilsSession) MULTIASSETPROXYID() ([4]byte, error) {
	return _DevUtils.Contract.MULTIASSETPROXYID(&_DevUtils.CallOpts)
}

// MULTIASSETPROXYID is a free data retrieval call binding the contract method 0xb37fda04.
//
// Solidity: function MULTI_ASSET_PROXY_ID() constant returns(bytes4)
func (_DevUtils *DevUtilsCallerSession) MULTIASSETPROXYID() ([4]byte, error) {
	return _DevUtils.Contract.MULTIASSETPROXYID(&_DevUtils.CallOpts)
}

// STATICCALLPROXYID is a free data retrieval call binding the contract method 0xd965b998.
//
// Solidity: function STATIC_CALL_PROXY_ID() constant returns(bytes4)
func (_DevUtils *DevUtilsCaller) STATICCALLPROXYID(opts *bind.CallOpts) ([4]byte, error) {
	var (
		ret0 = new([4]byte)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "STATIC_CALL_PROXY_ID")
	return *ret0, err
}

// STATICCALLPROXYID is a free data retrieval call binding the contract method 0xd965b998.
//
// Solidity: function STATIC_CALL_PROXY_ID() constant returns(bytes4)
func (_DevUtils *DevUtilsSession) STATICCALLPROXYID() ([4]byte, error) {
	return _DevUtils.Contract.STATICCALLPROXYID(&_DevUtils.CallOpts)
}

// STATICCALLPROXYID is a free data retrieval call binding the contract method 0xd965b998.
//
// Solidity: function STATIC_CALL_PROXY_ID() constant returns(bytes4)
func (_DevUtils *DevUtilsCallerSession) STATICCALLPROXYID() ([4]byte, error) {
	return _DevUtils.Contract.STATICCALLPROXYID(&_DevUtils.CallOpts)
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

// GetERC721TokenOwner is a free data retrieval call binding the contract method 0xb6988463.
//
// Solidity: function getERC721TokenOwner(address tokenAddress, uint256 tokenId) constant returns(address ownerAddress)
func (_DevUtils *DevUtilsCaller) GetERC721TokenOwner(opts *bind.CallOpts, tokenAddress common.Address, tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DevUtils.contract.Call(opts, out, "getERC721TokenOwner", tokenAddress, tokenId)
	return *ret0, err
}

// GetERC721TokenOwner is a free data retrieval call binding the contract method 0xb6988463.
//
// Solidity: function getERC721TokenOwner(address tokenAddress, uint256 tokenId) constant returns(address ownerAddress)
func (_DevUtils *DevUtilsSession) GetERC721TokenOwner(tokenAddress common.Address, tokenId *big.Int) (common.Address, error) {
	return _DevUtils.Contract.GetERC721TokenOwner(&_DevUtils.CallOpts, tokenAddress, tokenId)
}

// GetERC721TokenOwner is a free data retrieval call binding the contract method 0xb6988463.
//
// Solidity: function getERC721TokenOwner(address tokenAddress, uint256 tokenId) constant returns(address ownerAddress)
func (_DevUtils *DevUtilsCallerSession) GetERC721TokenOwner(tokenAddress common.Address, tokenId *big.Int) (common.Address, error) {
	return _DevUtils.Contract.GetERC721TokenOwner(&_DevUtils.CallOpts, tokenAddress, tokenId)
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

// GetOrderRelevantState is a free data retrieval call binding the contract method 0x8f5afa52.
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

// GetOrderRelevantState is a free data retrieval call binding the contract method 0x8f5afa52.
//
// Solidity: function getOrderRelevantState(Struct0 order, bytes signature) constant returns(Struct1 orderInfo, uint256 fillableTakerAssetAmount, bool isValidSignature)
func (_DevUtils *DevUtilsSession) GetOrderRelevantState(order Struct0, signature []byte) (struct {
	OrderInfo                Struct1
	FillableTakerAssetAmount *big.Int
	IsValidSignature         bool
}, error) {
	return _DevUtils.Contract.GetOrderRelevantState(&_DevUtils.CallOpts, order, signature)
}

// GetOrderRelevantState is a free data retrieval call binding the contract method 0x8f5afa52.
//
// Solidity: function getOrderRelevantState(Struct0 order, bytes signature) constant returns(Struct1 orderInfo, uint256 fillableTakerAssetAmount, bool isValidSignature)
func (_DevUtils *DevUtilsCallerSession) GetOrderRelevantState(order Struct0, signature []byte) (struct {
	OrderInfo                Struct1
	FillableTakerAssetAmount *big.Int
	IsValidSignature         bool
}, error) {
	return _DevUtils.Contract.GetOrderRelevantState(&_DevUtils.CallOpts, order, signature)
}

// GetOrderRelevantStates is a free data retrieval call binding the contract method 0x7f46448d.
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

// GetOrderRelevantStates is a free data retrieval call binding the contract method 0x7f46448d.
//
// Solidity: function getOrderRelevantStates([]Struct0 orders, bytes[] signatures) constant returns([]Struct1 ordersInfo, uint256[] fillableTakerAssetAmounts, bool[] isValidSignature)
func (_DevUtils *DevUtilsSession) GetOrderRelevantStates(orders []Struct0, signatures [][]byte) (struct {
	OrdersInfo                []Struct1
	FillableTakerAssetAmounts []*big.Int
	IsValidSignature          []bool
}, error) {
	return _DevUtils.Contract.GetOrderRelevantStates(&_DevUtils.CallOpts, orders, signatures)
}

// GetOrderRelevantStates is a free data retrieval call binding the contract method 0x7f46448d.
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
