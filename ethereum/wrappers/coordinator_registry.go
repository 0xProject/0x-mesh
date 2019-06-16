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

// CoordinatorRegistryABI is the input ABI used to generate the binding from.
const CoordinatorRegistryABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"coordinatorEndpoint\",\"type\":\"string\"}],\"name\":\"setCoordinatorEndpoint\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"coordinatorOperator\",\"type\":\"address\"}],\"name\":\"getCoordinatorEndpoint\",\"outputs\":[{\"name\":\"coordinatorEndpoint\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"coordinatorOperator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"coordinatorEndpoint\",\"type\":\"string\"}],\"name\":\"CoordinatorEndpointSet\",\"type\":\"event\"}]"

// CoordinatorRegistry is an auto generated Go binding around an Ethereum contract.
type CoordinatorRegistry struct {
	CoordinatorRegistryCaller     // Read-only binding to the contract
	CoordinatorRegistryTransactor // Write-only binding to the contract
	CoordinatorRegistryFilterer   // Log filterer for contract events
}

// CoordinatorRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type CoordinatorRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CoordinatorRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CoordinatorRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CoordinatorRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CoordinatorRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CoordinatorRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CoordinatorRegistrySession struct {
	Contract     *CoordinatorRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// CoordinatorRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CoordinatorRegistryCallerSession struct {
	Contract *CoordinatorRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// CoordinatorRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CoordinatorRegistryTransactorSession struct {
	Contract     *CoordinatorRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// CoordinatorRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type CoordinatorRegistryRaw struct {
	Contract *CoordinatorRegistry // Generic contract binding to access the raw methods on
}

// CoordinatorRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CoordinatorRegistryCallerRaw struct {
	Contract *CoordinatorRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// CoordinatorRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CoordinatorRegistryTransactorRaw struct {
	Contract *CoordinatorRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCoordinatorRegistry creates a new instance of CoordinatorRegistry, bound to a specific deployed contract.
func NewCoordinatorRegistry(address common.Address, backend bind.ContractBackend) (*CoordinatorRegistry, error) {
	contract, err := bindCoordinatorRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CoordinatorRegistry{CoordinatorRegistryCaller: CoordinatorRegistryCaller{contract: contract}, CoordinatorRegistryTransactor: CoordinatorRegistryTransactor{contract: contract}, CoordinatorRegistryFilterer: CoordinatorRegistryFilterer{contract: contract}}, nil
}

// NewCoordinatorRegistryCaller creates a new read-only instance of CoordinatorRegistry, bound to a specific deployed contract.
func NewCoordinatorRegistryCaller(address common.Address, caller bind.ContractCaller) (*CoordinatorRegistryCaller, error) {
	contract, err := bindCoordinatorRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CoordinatorRegistryCaller{contract: contract}, nil
}

// NewCoordinatorRegistryTransactor creates a new write-only instance of CoordinatorRegistry, bound to a specific deployed contract.
func NewCoordinatorRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*CoordinatorRegistryTransactor, error) {
	contract, err := bindCoordinatorRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CoordinatorRegistryTransactor{contract: contract}, nil
}

// NewCoordinatorRegistryFilterer creates a new log filterer instance of CoordinatorRegistry, bound to a specific deployed contract.
func NewCoordinatorRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*CoordinatorRegistryFilterer, error) {
	contract, err := bindCoordinatorRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CoordinatorRegistryFilterer{contract: contract}, nil
}

// bindCoordinatorRegistry binds a generic wrapper to an already deployed contract.
func bindCoordinatorRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CoordinatorRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CoordinatorRegistry *CoordinatorRegistryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CoordinatorRegistry.Contract.CoordinatorRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CoordinatorRegistry *CoordinatorRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CoordinatorRegistry.Contract.CoordinatorRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CoordinatorRegistry *CoordinatorRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CoordinatorRegistry.Contract.CoordinatorRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CoordinatorRegistry *CoordinatorRegistryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CoordinatorRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CoordinatorRegistry *CoordinatorRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CoordinatorRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CoordinatorRegistry *CoordinatorRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CoordinatorRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetCoordinatorEndpoint is a free data retrieval call binding the contract method 0x6c90fedb.
//
// Solidity: function getCoordinatorEndpoint(address coordinatorOperator) constant returns(string coordinatorEndpoint)
func (_CoordinatorRegistry *CoordinatorRegistryCaller) GetCoordinatorEndpoint(opts *bind.CallOpts, coordinatorOperator common.Address) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _CoordinatorRegistry.contract.Call(opts, out, "getCoordinatorEndpoint", coordinatorOperator)
	return *ret0, err
}

// GetCoordinatorEndpoint is a free data retrieval call binding the contract method 0x6c90fedb.
//
// Solidity: function getCoordinatorEndpoint(address coordinatorOperator) constant returns(string coordinatorEndpoint)
func (_CoordinatorRegistry *CoordinatorRegistrySession) GetCoordinatorEndpoint(coordinatorOperator common.Address) (string, error) {
	return _CoordinatorRegistry.Contract.GetCoordinatorEndpoint(&_CoordinatorRegistry.CallOpts, coordinatorOperator)
}

// GetCoordinatorEndpoint is a free data retrieval call binding the contract method 0x6c90fedb.
//
// Solidity: function getCoordinatorEndpoint(address coordinatorOperator) constant returns(string coordinatorEndpoint)
func (_CoordinatorRegistry *CoordinatorRegistryCallerSession) GetCoordinatorEndpoint(coordinatorOperator common.Address) (string, error) {
	return _CoordinatorRegistry.Contract.GetCoordinatorEndpoint(&_CoordinatorRegistry.CallOpts, coordinatorOperator)
}

// SetCoordinatorEndpoint is a paid mutator transaction binding the contract method 0x5b2388be.
//
// Solidity: function setCoordinatorEndpoint(string coordinatorEndpoint) returns()
func (_CoordinatorRegistry *CoordinatorRegistryTransactor) SetCoordinatorEndpoint(opts *bind.TransactOpts, coordinatorEndpoint string) (*types.Transaction, error) {
	return _CoordinatorRegistry.contract.Transact(opts, "setCoordinatorEndpoint", coordinatorEndpoint)
}

// SetCoordinatorEndpoint is a paid mutator transaction binding the contract method 0x5b2388be.
//
// Solidity: function setCoordinatorEndpoint(string coordinatorEndpoint) returns()
func (_CoordinatorRegistry *CoordinatorRegistrySession) SetCoordinatorEndpoint(coordinatorEndpoint string) (*types.Transaction, error) {
	return _CoordinatorRegistry.Contract.SetCoordinatorEndpoint(&_CoordinatorRegistry.TransactOpts, coordinatorEndpoint)
}

// SetCoordinatorEndpoint is a paid mutator transaction binding the contract method 0x5b2388be.
//
// Solidity: function setCoordinatorEndpoint(string coordinatorEndpoint) returns()
func (_CoordinatorRegistry *CoordinatorRegistryTransactorSession) SetCoordinatorEndpoint(coordinatorEndpoint string) (*types.Transaction, error) {
	return _CoordinatorRegistry.Contract.SetCoordinatorEndpoint(&_CoordinatorRegistry.TransactOpts, coordinatorEndpoint)
}

// CoordinatorRegistryCoordinatorEndpointSetIterator is returned from FilterCoordinatorEndpointSet and is used to iterate over the raw logs and unpacked data for CoordinatorEndpointSet events raised by the CoordinatorRegistry contract.
type CoordinatorRegistryCoordinatorEndpointSetIterator struct {
	Event *CoordinatorRegistryCoordinatorEndpointSet // Event containing the contract specifics and raw log

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
func (it *CoordinatorRegistryCoordinatorEndpointSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CoordinatorRegistryCoordinatorEndpointSet)
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
		it.Event = new(CoordinatorRegistryCoordinatorEndpointSet)
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
func (it *CoordinatorRegistryCoordinatorEndpointSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CoordinatorRegistryCoordinatorEndpointSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CoordinatorRegistryCoordinatorEndpointSet represents a CoordinatorEndpointSet event raised by the CoordinatorRegistry contract.
type CoordinatorRegistryCoordinatorEndpointSet struct {
	CoordinatorOperator common.Address
	CoordinatorEndpoint string
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterCoordinatorEndpointSet is a free log retrieval operation binding the contract event 0xd060052768902f3eecb84b8eae9d3a2608a1a9e60811a33968b46b8d552f266e.
//
// Solidity: event CoordinatorEndpointSet(address coordinatorOperator, string coordinatorEndpoint)
func (_CoordinatorRegistry *CoordinatorRegistryFilterer) FilterCoordinatorEndpointSet(opts *bind.FilterOpts) (*CoordinatorRegistryCoordinatorEndpointSetIterator, error) {

	logs, sub, err := _CoordinatorRegistry.contract.FilterLogs(opts, "CoordinatorEndpointSet")
	if err != nil {
		return nil, err
	}
	return &CoordinatorRegistryCoordinatorEndpointSetIterator{contract: _CoordinatorRegistry.contract, event: "CoordinatorEndpointSet", logs: logs, sub: sub}, nil
}

// WatchCoordinatorEndpointSet is a free log subscription operation binding the contract event 0xd060052768902f3eecb84b8eae9d3a2608a1a9e60811a33968b46b8d552f266e.
//
// Solidity: event CoordinatorEndpointSet(address coordinatorOperator, string coordinatorEndpoint)
func (_CoordinatorRegistry *CoordinatorRegistryFilterer) WatchCoordinatorEndpointSet(opts *bind.WatchOpts, sink chan<- *CoordinatorRegistryCoordinatorEndpointSet) (event.Subscription, error) {

	logs, sub, err := _CoordinatorRegistry.contract.WatchLogs(opts, "CoordinatorEndpointSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CoordinatorRegistryCoordinatorEndpointSet)
				if err := _CoordinatorRegistry.contract.UnpackLog(event, "CoordinatorEndpointSet", log); err != nil {
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
