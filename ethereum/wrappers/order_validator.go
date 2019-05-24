// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.
// Fabio: This file was generated with a fork of `go-ethereum` that added
// ABIV2 support: https://github.com/ethereum/go-ethereum/pull/18491.

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

// OrderValidatorABI is the input ABI used to generate the binding from.
const OrderValidatorABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"target\",\"type\":\"address\"},{\"name\":\"assetData\",\"type\":\"bytes\"}],\"name\":\"getBalanceAndAllowance\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"},{\"name\":\"allowance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"takerAddresses\",\"type\":\"address[]\"}],\"name\":\"getTradersInfo\",\"outputs\":[{\"components\":[{\"name\":\"makerBalance\",\"type\":\"uint256\"},{\"name\":\"makerAllowance\",\"type\":\"uint256\"},{\"name\":\"takerBalance\",\"type\":\"uint256\"},{\"name\":\"takerAllowance\",\"type\":\"uint256\"},{\"name\":\"makerZrxBalance\",\"type\":\"uint256\"},{\"name\":\"makerZrxAllowance\",\"type\":\"uint256\"},{\"name\":\"takerZrxBalance\",\"type\":\"uint256\"},{\"name\":\"takerZrxAllowance\",\"type\":\"uint256\"}],\"name\":\"\",\"type\":\"tuple[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getERC721TokenOwner\",\"outputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"},{\"name\":\"signature\",\"type\":\"bytes\"},{\"name\":\"takerAddress\",\"type\":\"address\"}],\"name\":\"getOrderAndTraderInfo\",\"outputs\":[{\"components\":[{\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"name\":\"orderInfo\",\"type\":\"tuple\"},{\"components\":[{\"name\":\"makerBalance\",\"type\":\"uint256\"},{\"name\":\"makerAllowance\",\"type\":\"uint256\"},{\"name\":\"takerBalance\",\"type\":\"uint256\"},{\"name\":\"takerAllowance\",\"type\":\"uint256\"},{\"name\":\"makerZrxBalance\",\"type\":\"uint256\"},{\"name\":\"makerZrxAllowance\",\"type\":\"uint256\"},{\"name\":\"takerZrxBalance\",\"type\":\"uint256\"},{\"name\":\"takerZrxAllowance\",\"type\":\"uint256\"}],\"name\":\"traderInfo\",\"type\":\"tuple\"},{\"name\":\"isValidSignature\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"target\",\"type\":\"address\"},{\"name\":\"assetData\",\"type\":\"bytes[]\"}],\"name\":\"getBalancesAndAllowances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"},{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"order\",\"type\":\"tuple\"},{\"name\":\"takerAddress\",\"type\":\"address\"}],\"name\":\"getTraderInfo\",\"outputs\":[{\"components\":[{\"name\":\"makerBalance\",\"type\":\"uint256\"},{\"name\":\"makerAllowance\",\"type\":\"uint256\"},{\"name\":\"takerBalance\",\"type\":\"uint256\"},{\"name\":\"takerAllowance\",\"type\":\"uint256\"},{\"name\":\"makerZrxBalance\",\"type\":\"uint256\"},{\"name\":\"makerZrxAllowance\",\"type\":\"uint256\"},{\"name\":\"takerZrxBalance\",\"type\":\"uint256\"},{\"name\":\"takerZrxAllowance\",\"type\":\"uint256\"}],\"name\":\"traderInfo\",\"type\":\"tuple\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"components\":[{\"name\":\"makerAddress\",\"type\":\"address\"},{\"name\":\"takerAddress\",\"type\":\"address\"},{\"name\":\"feeRecipientAddress\",\"type\":\"address\"},{\"name\":\"senderAddress\",\"type\":\"address\"},{\"name\":\"makerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"takerAssetAmount\",\"type\":\"uint256\"},{\"name\":\"makerFee\",\"type\":\"uint256\"},{\"name\":\"takerFee\",\"type\":\"uint256\"},{\"name\":\"expirationTimeSeconds\",\"type\":\"uint256\"},{\"name\":\"salt\",\"type\":\"uint256\"},{\"name\":\"makerAssetData\",\"type\":\"bytes\"},{\"name\":\"takerAssetData\",\"type\":\"bytes\"}],\"name\":\"orders\",\"type\":\"tuple[]\"},{\"name\":\"signatures\",\"type\":\"bytes[]\"},{\"name\":\"takerAddresses\",\"type\":\"address[]\"}],\"name\":\"getOrdersAndTradersInfo\",\"outputs\":[{\"components\":[{\"name\":\"orderStatus\",\"type\":\"uint8\"},{\"name\":\"orderHash\",\"type\":\"bytes32\"},{\"name\":\"orderTakerAssetFilledAmount\",\"type\":\"uint256\"}],\"name\":\"ordersInfo\",\"type\":\"tuple[]\"},{\"components\":[{\"name\":\"makerBalance\",\"type\":\"uint256\"},{\"name\":\"makerAllowance\",\"type\":\"uint256\"},{\"name\":\"takerBalance\",\"type\":\"uint256\"},{\"name\":\"takerAllowance\",\"type\":\"uint256\"},{\"name\":\"makerZrxBalance\",\"type\":\"uint256\"},{\"name\":\"makerZrxAllowance\",\"type\":\"uint256\"},{\"name\":\"takerZrxBalance\",\"type\":\"uint256\"},{\"name\":\"takerZrxAllowance\",\"type\":\"uint256\"}],\"name\":\"tradersInfo\",\"type\":\"tuple[]\"},{\"name\":\"isValidSignature\",\"type\":\"bool[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_exchange\",\"type\":\"address\"},{\"name\":\"_zrxAssetData\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// OrderValidator is an auto generated Go binding around an Ethereum contract.
type OrderValidator struct {
	OrderValidatorCaller     // Read-only binding to the contract
	OrderValidatorTransactor // Write-only binding to the contract
	OrderValidatorFilterer   // Log filterer for contract events
}

// OrderValidatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type OrderValidatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OrderValidatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OrderValidatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OrderValidatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OrderValidatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OrderValidatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OrderValidatorSession struct {
	Contract     *OrderValidator   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OrderValidatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OrderValidatorCallerSession struct {
	Contract *OrderValidatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// OrderValidatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OrderValidatorTransactorSession struct {
	Contract     *OrderValidatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// OrderValidatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type OrderValidatorRaw struct {
	Contract *OrderValidator // Generic contract binding to access the raw methods on
}

// OrderValidatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OrderValidatorCallerRaw struct {
	Contract *OrderValidatorCaller // Generic read-only contract binding to access the raw methods on
}

// OrderValidatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OrderValidatorTransactorRaw struct {
	Contract *OrderValidatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOrderValidator creates a new instance of OrderValidator, bound to a specific deployed contract.
func NewOrderValidator(address common.Address, backend bind.ContractBackend) (*OrderValidator, error) {
	contract, err := bindOrderValidator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OrderValidator{OrderValidatorCaller: OrderValidatorCaller{contract: contract}, OrderValidatorTransactor: OrderValidatorTransactor{contract: contract}, OrderValidatorFilterer: OrderValidatorFilterer{contract: contract}}, nil
}

// NewOrderValidatorCaller creates a new read-only instance of OrderValidator, bound to a specific deployed contract.
func NewOrderValidatorCaller(address common.Address, caller bind.ContractCaller) (*OrderValidatorCaller, error) {
	contract, err := bindOrderValidator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OrderValidatorCaller{contract: contract}, nil
}

// NewOrderValidatorTransactor creates a new write-only instance of OrderValidator, bound to a specific deployed contract.
func NewOrderValidatorTransactor(address common.Address, transactor bind.ContractTransactor) (*OrderValidatorTransactor, error) {
	contract, err := bindOrderValidator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OrderValidatorTransactor{contract: contract}, nil
}

// NewOrderValidatorFilterer creates a new log filterer instance of OrderValidator, bound to a specific deployed contract.
func NewOrderValidatorFilterer(address common.Address, filterer bind.ContractFilterer) (*OrderValidatorFilterer, error) {
	contract, err := bindOrderValidator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OrderValidatorFilterer{contract: contract}, nil
}

// bindOrderValidator binds a generic wrapper to an already deployed contract.
func bindOrderValidator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OrderValidatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OrderValidator *OrderValidatorRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OrderValidator.Contract.OrderValidatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OrderValidator *OrderValidatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OrderValidator.Contract.OrderValidatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OrderValidator *OrderValidatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OrderValidator.Contract.OrderValidatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OrderValidator *OrderValidatorCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OrderValidator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OrderValidator *OrderValidatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OrderValidator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OrderValidator *OrderValidatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OrderValidator.Contract.contract.Transact(opts, method, params...)
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
	MakerBalance      *big.Int
	MakerAllowance    *big.Int
	TakerBalance      *big.Int
	TakerAllowance    *big.Int
	MakerZrxBalance   *big.Int
	MakerZrxAllowance *big.Int
	TakerZrxBalance   *big.Int
	TakerZrxAllowance *big.Int
}

// Struct2 is an auto generated low-level Go binding around an user-defined struct.
type Struct2 struct {
	OrderStatus                 uint8
	OrderHash                   [32]byte
	OrderTakerAssetFilledAmount *big.Int
}

// GetBalanceAndAllowance is a free data retrieval call binding the contract method 0x2cd0fc73.
//
// Solidity: function getBalanceAndAllowance(address target, bytes assetData) constant returns(uint256 balance, uint256 allowance)
func (_OrderValidator *OrderValidatorCaller) GetBalanceAndAllowance(opts *bind.CallOpts, target common.Address, assetData []byte) (struct {
	Balance   *big.Int
	Allowance *big.Int
}, error) {
	ret := new(struct {
		Balance   *big.Int
		Allowance *big.Int
	})
	out := ret
	err := _OrderValidator.contract.Call(opts, out, "getBalanceAndAllowance", target, assetData)
	return *ret, err
}

// GetBalanceAndAllowance is a free data retrieval call binding the contract method 0x2cd0fc73.
//
// Solidity: function getBalanceAndAllowance(address target, bytes assetData) constant returns(uint256 balance, uint256 allowance)
func (_OrderValidator *OrderValidatorSession) GetBalanceAndAllowance(target common.Address, assetData []byte) (struct {
	Balance   *big.Int
	Allowance *big.Int
}, error) {
	return _OrderValidator.Contract.GetBalanceAndAllowance(&_OrderValidator.CallOpts, target, assetData)
}

// GetBalanceAndAllowance is a free data retrieval call binding the contract method 0x2cd0fc73.
//
// Solidity: function getBalanceAndAllowance(address target, bytes assetData) constant returns(uint256 balance, uint256 allowance)
func (_OrderValidator *OrderValidatorCallerSession) GetBalanceAndAllowance(target common.Address, assetData []byte) (struct {
	Balance   *big.Int
	Allowance *big.Int
}, error) {
	return _OrderValidator.Contract.GetBalanceAndAllowance(&_OrderValidator.CallOpts, target, assetData)
}

// GetBalancesAndAllowances is a free data retrieval call binding the contract method 0xc6b7f4ee.
//
// Solidity: function getBalancesAndAllowances(address target, bytes[] assetData) constant returns(uint256[], uint256[])
func (_OrderValidator *OrderValidatorCaller) GetBalancesAndAllowances(opts *bind.CallOpts, target common.Address, assetData [][]byte) ([]*big.Int, []*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
		ret1 = new([]*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _OrderValidator.contract.Call(opts, out, "getBalancesAndAllowances", target, assetData)
	return *ret0, *ret1, err
}

// GetBalancesAndAllowances is a free data retrieval call binding the contract method 0xc6b7f4ee.
//
// Solidity: function getBalancesAndAllowances(address target, bytes[] assetData) constant returns(uint256[], uint256[])
func (_OrderValidator *OrderValidatorSession) GetBalancesAndAllowances(target common.Address, assetData [][]byte) ([]*big.Int, []*big.Int, error) {
	return _OrderValidator.Contract.GetBalancesAndAllowances(&_OrderValidator.CallOpts, target, assetData)
}

// GetBalancesAndAllowances is a free data retrieval call binding the contract method 0xc6b7f4ee.
//
// Solidity: function getBalancesAndAllowances(address target, bytes[] assetData) constant returns(uint256[], uint256[])
func (_OrderValidator *OrderValidatorCallerSession) GetBalancesAndAllowances(target common.Address, assetData [][]byte) ([]*big.Int, []*big.Int, error) {
	return _OrderValidator.Contract.GetBalancesAndAllowances(&_OrderValidator.CallOpts, target, assetData)
}

// GetERC721TokenOwner is a free data retrieval call binding the contract method 0xb6988463.
//
// Solidity: function getERC721TokenOwner(address token, uint256 tokenId) constant returns(address owner)
func (_OrderValidator *OrderValidatorCaller) GetERC721TokenOwner(opts *bind.CallOpts, token common.Address, tokenId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _OrderValidator.contract.Call(opts, out, "getERC721TokenOwner", token, tokenId)
	return *ret0, err
}

// GetERC721TokenOwner is a free data retrieval call binding the contract method 0xb6988463.
//
// Solidity: function getERC721TokenOwner(address token, uint256 tokenId) constant returns(address owner)
func (_OrderValidator *OrderValidatorSession) GetERC721TokenOwner(token common.Address, tokenId *big.Int) (common.Address, error) {
	return _OrderValidator.Contract.GetERC721TokenOwner(&_OrderValidator.CallOpts, token, tokenId)
}

// GetERC721TokenOwner is a free data retrieval call binding the contract method 0xb6988463.
//
// Solidity: function getERC721TokenOwner(address token, uint256 tokenId) constant returns(address owner)
func (_OrderValidator *OrderValidatorCallerSession) GetERC721TokenOwner(token common.Address, tokenId *big.Int) (common.Address, error) {
	return _OrderValidator.Contract.GetERC721TokenOwner(&_OrderValidator.CallOpts, token, tokenId)
}

// GetOrderAndTraderInfo is a free data retrieval call binding the contract method 0xc6319fc5.
//
// Solidity: function getOrderAndTraderInfo(Struct0 order, bytes signature, address takerAddress) constant returns(Struct2 orderInfo, Struct1 traderInfo, bool isValidSignature)
func (_OrderValidator *OrderValidatorCaller) GetOrderAndTraderInfo(opts *bind.CallOpts, order Struct0, signature []byte, takerAddress common.Address) (struct {
	OrderInfo        Struct2
	TraderInfo       Struct1
	IsValidSignature bool
}, error) {
	ret := new(struct {
		OrderInfo        Struct2
		TraderInfo       Struct1
		IsValidSignature bool
	})
	out := ret
	err := _OrderValidator.contract.Call(opts, out, "getOrderAndTraderInfo", order, signature, takerAddress)
	return *ret, err
}

// GetOrderAndTraderInfo is a free data retrieval call binding the contract method 0xc6319fc5.
//
// Solidity: function getOrderAndTraderInfo(Struct0 order, bytes signature, address takerAddress) constant returns(Struct2 orderInfo, Struct1 traderInfo, bool isValidSignature)
func (_OrderValidator *OrderValidatorSession) GetOrderAndTraderInfo(order Struct0, signature []byte, takerAddress common.Address) (struct {
	OrderInfo        Struct2
	TraderInfo       Struct1
	IsValidSignature bool
}, error) {
	return _OrderValidator.Contract.GetOrderAndTraderInfo(&_OrderValidator.CallOpts, order, signature, takerAddress)
}

// GetOrderAndTraderInfo is a free data retrieval call binding the contract method 0xc6319fc5.
//
// Solidity: function getOrderAndTraderInfo(Struct0 order, bytes signature, address takerAddress) constant returns(Struct2 orderInfo, Struct1 traderInfo, bool isValidSignature)
func (_OrderValidator *OrderValidatorCallerSession) GetOrderAndTraderInfo(order Struct0, signature []byte, takerAddress common.Address) (struct {
	OrderInfo        Struct2
	TraderInfo       Struct1
	IsValidSignature bool
}, error) {
	return _OrderValidator.Contract.GetOrderAndTraderInfo(&_OrderValidator.CallOpts, order, signature, takerAddress)
}

// GetOrdersAndTradersInfo is a free data retrieval call binding the contract method 0xf7c1c03f.
//
// Solidity: function getOrdersAndTradersInfo((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes)[] orders, bytes[] signatures, address[] takerAddresses) constant returns((uint8,bytes32,uint256)[] ordersInfo, (uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[] tradersInfo, bool[] isValidSignature)
func (_OrderValidator *OrderValidatorCaller) GetOrdersAndTradersInfo(opts *bind.CallOpts, orders []Struct0, signatures [][]byte, takerAddresses []common.Address) (struct {
	OrdersInfo       []Struct2
	TradersInfo      []Struct1
	IsValidSignature []bool
}, error) {
	ret := new(struct {
		OrdersInfo       []Struct2
		TradersInfo      []Struct1
		IsValidSignature []bool
	})
	out := ret
	err := _OrderValidator.contract.Call(opts, out, "getOrdersAndTradersInfo", orders, signatures, takerAddresses)
	return *ret, err
}

// GetOrdersAndTradersInfo is a free data retrieval call binding the contract method 0xf7c1c03f.
//
// Solidity: function getOrdersAndTradersInfo((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes)[] orders, bytes[] signatures, address[] takerAddresses) constant returns((uint8,bytes32,uint256)[] ordersInfo, (uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[] tradersInfo, bool[] isValidSignature)
func (_OrderValidator *OrderValidatorSession) GetOrdersAndTradersInfo(orders []Struct0, signatures [][]byte, takerAddresses []common.Address) (struct {
	OrdersInfo       []Struct2
	TradersInfo      []Struct1
	IsValidSignature []bool
}, error) {
	return _OrderValidator.Contract.GetOrdersAndTradersInfo(&_OrderValidator.CallOpts, orders, signatures, takerAddresses)
}

// GetOrdersAndTradersInfo is a free data retrieval call binding the contract method 0xf7c1c03f.
//
// Solidity: function getOrdersAndTradersInfo((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes)[] orders, bytes[] signatures, address[] takerAddresses) constant returns((uint8,bytes32,uint256)[] ordersInfo, (uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[] tradersInfo, bool[] isValidSignature)
func (_OrderValidator *OrderValidatorCallerSession) GetOrdersAndTradersInfo(orders []Struct0, signatures [][]byte, takerAddresses []common.Address) (struct {
	OrdersInfo       []Struct2
	TradersInfo      []Struct1
	IsValidSignature []bool
}, error) {
	return _OrderValidator.Contract.GetOrdersAndTradersInfo(&_OrderValidator.CallOpts, orders, signatures, takerAddresses)
}

// GetTraderInfo is a free data retrieval call binding the contract method 0xf241ffb0.
//
// Solidity: function getTraderInfo(Struct0 order, address takerAddress) constant returns(Struct1 traderInfo)
func (_OrderValidator *OrderValidatorCaller) GetTraderInfo(opts *bind.CallOpts, order Struct0, takerAddress common.Address) (Struct1, error) {
	var (
		ret0 = new(Struct1)
	)
	out := ret0
	err := _OrderValidator.contract.Call(opts, out, "getTraderInfo", order, takerAddress)
	return *ret0, err
}

// GetTraderInfo is a free data retrieval call binding the contract method 0xf241ffb0.
//
// Solidity: function getTraderInfo(Struct0 order, address takerAddress) constant returns(Struct1 traderInfo)
func (_OrderValidator *OrderValidatorSession) GetTraderInfo(order Struct0, takerAddress common.Address) (Struct1, error) {
	return _OrderValidator.Contract.GetTraderInfo(&_OrderValidator.CallOpts, order, takerAddress)
}

// GetTraderInfo is a free data retrieval call binding the contract method 0xf241ffb0.
//
// Solidity: function getTraderInfo(Struct0 order, address takerAddress) constant returns(Struct1 traderInfo)
func (_OrderValidator *OrderValidatorCallerSession) GetTraderInfo(order Struct0, takerAddress common.Address) (Struct1, error) {
	return _OrderValidator.Contract.GetTraderInfo(&_OrderValidator.CallOpts, order, takerAddress)
}

// GetTradersInfo is a free data retrieval call binding the contract method 0x690d3114.
//
// Solidity: function getTradersInfo((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes)[] orders, address[] takerAddresses) constant returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_OrderValidator *OrderValidatorCaller) GetTradersInfo(opts *bind.CallOpts, orders []Struct0, takerAddresses []common.Address) ([]Struct1, error) {
	var (
		ret0 = new([]Struct1)
	)
	out := ret0
	err := _OrderValidator.contract.Call(opts, out, "getTradersInfo", orders, takerAddresses)
	return *ret0, err
}

// GetTradersInfo is a free data retrieval call binding the contract method 0x690d3114.
//
// Solidity: function getTradersInfo((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes)[] orders, address[] takerAddresses) constant returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_OrderValidator *OrderValidatorSession) GetTradersInfo(orders []Struct0, takerAddresses []common.Address) ([]Struct1, error) {
	return _OrderValidator.Contract.GetTradersInfo(&_OrderValidator.CallOpts, orders, takerAddresses)
}

// GetTradersInfo is a free data retrieval call binding the contract method 0x690d3114.
//
// Solidity: function getTradersInfo((address,address,address,address,uint256,uint256,uint256,uint256,uint256,uint256,bytes,bytes)[] orders, address[] takerAddresses) constant returns((uint256,uint256,uint256,uint256,uint256,uint256,uint256,uint256)[])
func (_OrderValidator *OrderValidatorCallerSession) GetTradersInfo(orders []Struct0, takerAddresses []common.Address) ([]Struct1, error) {
	return _OrderValidator.Contract.GetTradersInfo(&_OrderValidator.CallOpts, orders, takerAddresses)
}
