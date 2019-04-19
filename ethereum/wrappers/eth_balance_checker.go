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

// EthBalanceCheckerABI is the input ABI used to generate the binding from.
const EthBalanceCheckerABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"addresses\",\"type\":\"address[]\"}],\"name\":\"getEthBalances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// EthBalanceChecker is an auto generated Go binding around an Ethereum contract.
type EthBalanceChecker struct {
	EthBalanceCheckerCaller     // Read-only binding to the contract
	EthBalanceCheckerTransactor // Write-only binding to the contract
	EthBalanceCheckerFilterer   // Log filterer for contract events
}

// EthBalanceCheckerCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthBalanceCheckerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthBalanceCheckerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthBalanceCheckerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthBalanceCheckerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthBalanceCheckerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthBalanceCheckerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthBalanceCheckerSession struct {
	Contract     *EthBalanceChecker // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// EthBalanceCheckerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthBalanceCheckerCallerSession struct {
	Contract *EthBalanceCheckerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// EthBalanceCheckerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthBalanceCheckerTransactorSession struct {
	Contract     *EthBalanceCheckerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// EthBalanceCheckerRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthBalanceCheckerRaw struct {
	Contract *EthBalanceChecker // Generic contract binding to access the raw methods on
}

// EthBalanceCheckerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthBalanceCheckerCallerRaw struct {
	Contract *EthBalanceCheckerCaller // Generic read-only contract binding to access the raw methods on
}

// EthBalanceCheckerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthBalanceCheckerTransactorRaw struct {
	Contract *EthBalanceCheckerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEthBalanceChecker creates a new instance of EthBalanceChecker, bound to a specific deployed contract.
func NewEthBalanceChecker(address common.Address, backend bind.ContractBackend) (*EthBalanceChecker, error) {
	contract, err := bindEthBalanceChecker(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EthBalanceChecker{EthBalanceCheckerCaller: EthBalanceCheckerCaller{contract: contract}, EthBalanceCheckerTransactor: EthBalanceCheckerTransactor{contract: contract}, EthBalanceCheckerFilterer: EthBalanceCheckerFilterer{contract: contract}}, nil
}

// NewEthBalanceCheckerCaller creates a new read-only instance of EthBalanceChecker, bound to a specific deployed contract.
func NewEthBalanceCheckerCaller(address common.Address, caller bind.ContractCaller) (*EthBalanceCheckerCaller, error) {
	contract, err := bindEthBalanceChecker(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthBalanceCheckerCaller{contract: contract}, nil
}

// NewEthBalanceCheckerTransactor creates a new write-only instance of EthBalanceChecker, bound to a specific deployed contract.
func NewEthBalanceCheckerTransactor(address common.Address, transactor bind.ContractTransactor) (*EthBalanceCheckerTransactor, error) {
	contract, err := bindEthBalanceChecker(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthBalanceCheckerTransactor{contract: contract}, nil
}

// NewEthBalanceCheckerFilterer creates a new log filterer instance of EthBalanceChecker, bound to a specific deployed contract.
func NewEthBalanceCheckerFilterer(address common.Address, filterer bind.ContractFilterer) (*EthBalanceCheckerFilterer, error) {
	contract, err := bindEthBalanceChecker(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthBalanceCheckerFilterer{contract: contract}, nil
}

// bindEthBalanceChecker binds a generic wrapper to an already deployed contract.
func bindEthBalanceChecker(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EthBalanceCheckerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthBalanceChecker *EthBalanceCheckerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _EthBalanceChecker.Contract.EthBalanceCheckerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthBalanceChecker *EthBalanceCheckerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthBalanceChecker.Contract.EthBalanceCheckerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthBalanceChecker *EthBalanceCheckerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthBalanceChecker.Contract.EthBalanceCheckerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EthBalanceChecker *EthBalanceCheckerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _EthBalanceChecker.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EthBalanceChecker *EthBalanceCheckerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EthBalanceChecker.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EthBalanceChecker *EthBalanceCheckerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EthBalanceChecker.Contract.contract.Transact(opts, method, params...)
}

// GetEthBalances is a free data retrieval call binding the contract method 0xa0901e51.
//
// Solidity: function getEthBalances(address[] addresses) constant returns(uint256[])
func (_EthBalanceChecker *EthBalanceCheckerCaller) GetEthBalances(opts *bind.CallOpts, addresses []common.Address) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _EthBalanceChecker.contract.Call(opts, out, "getEthBalances", addresses)
	return *ret0, err
}

// GetEthBalances is a free data retrieval call binding the contract method 0xa0901e51.
//
// Solidity: function getEthBalances(address[] addresses) constant returns(uint256[])
func (_EthBalanceChecker *EthBalanceCheckerSession) GetEthBalances(addresses []common.Address) ([]*big.Int, error) {
	return _EthBalanceChecker.Contract.GetEthBalances(&_EthBalanceChecker.CallOpts, addresses)
}

// GetEthBalances is a free data retrieval call binding the contract method 0xa0901e51.
//
// Solidity: function getEthBalances(address[] addresses) constant returns(uint256[])
func (_EthBalanceChecker *EthBalanceCheckerCallerSession) GetEthBalances(addresses []common.Address) ([]*big.Int, error) {
	return _EthBalanceChecker.Contract.GetEthBalances(&_EthBalanceChecker.CallOpts, addresses)
}
