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

// ERC1155MintableMetaData contains all meta data concerning the ERC1155Mintable contract.
var ERC1155MintableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"TransferBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"URI\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"ERC1155_BATCH_RECEIVED\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ERC1155_RECEIVED\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"owners\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"name\":\"balanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"balances_\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"isNF\",\"type\":\"bool\"}],\"name\":\"create\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"type_\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"type_\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"}],\"name\":\"createWithType\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"creators\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getNonFungibleBaseType\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getNonFungibleIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"isFungible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"isNonFungible\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"isNonFungibleBaseType\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"isNonFungibleItem\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"maxIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"to\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"quantities\",\"type\":\"uint256[]\"}],\"name\":\"mintFungible\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"type_\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"to\",\"type\":\"address[]\"}],\"name\":\"mintNonFungible\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeBatchTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ERC1155MintableABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC1155MintableMetaData.ABI instead.
var ERC1155MintableABI = ERC1155MintableMetaData.ABI

// ERC1155Mintable is an auto generated Go binding around an Ethereum contract.
type ERC1155Mintable struct {
	ERC1155MintableCaller     // Read-only binding to the contract
	ERC1155MintableTransactor // Write-only binding to the contract
	ERC1155MintableFilterer   // Log filterer for contract events
}

// ERC1155MintableCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC1155MintableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1155MintableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC1155MintableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1155MintableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC1155MintableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1155MintableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC1155MintableSession struct {
	Contract     *ERC1155Mintable  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC1155MintableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC1155MintableCallerSession struct {
	Contract *ERC1155MintableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// ERC1155MintableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC1155MintableTransactorSession struct {
	Contract     *ERC1155MintableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ERC1155MintableRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC1155MintableRaw struct {
	Contract *ERC1155Mintable // Generic contract binding to access the raw methods on
}

// ERC1155MintableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC1155MintableCallerRaw struct {
	Contract *ERC1155MintableCaller // Generic read-only contract binding to access the raw methods on
}

// ERC1155MintableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC1155MintableTransactorRaw struct {
	Contract *ERC1155MintableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC1155Mintable creates a new instance of ERC1155Mintable, bound to a specific deployed contract.
func NewERC1155Mintable(address common.Address, backend bind.ContractBackend) (*ERC1155Mintable, error) {
	contract, err := bindERC1155Mintable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC1155Mintable{ERC1155MintableCaller: ERC1155MintableCaller{contract: contract}, ERC1155MintableTransactor: ERC1155MintableTransactor{contract: contract}, ERC1155MintableFilterer: ERC1155MintableFilterer{contract: contract}}, nil
}

// NewERC1155MintableCaller creates a new read-only instance of ERC1155Mintable, bound to a specific deployed contract.
func NewERC1155MintableCaller(address common.Address, caller bind.ContractCaller) (*ERC1155MintableCaller, error) {
	contract, err := bindERC1155Mintable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1155MintableCaller{contract: contract}, nil
}

// NewERC1155MintableTransactor creates a new write-only instance of ERC1155Mintable, bound to a specific deployed contract.
func NewERC1155MintableTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC1155MintableTransactor, error) {
	contract, err := bindERC1155Mintable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1155MintableTransactor{contract: contract}, nil
}

// NewERC1155MintableFilterer creates a new log filterer instance of ERC1155Mintable, bound to a specific deployed contract.
func NewERC1155MintableFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC1155MintableFilterer, error) {
	contract, err := bindERC1155Mintable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC1155MintableFilterer{contract: contract}, nil
}

// bindERC1155Mintable binds a generic wrapper to an already deployed contract.
func bindERC1155Mintable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC1155MintableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1155Mintable *ERC1155MintableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1155Mintable.Contract.ERC1155MintableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1155Mintable *ERC1155MintableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.ERC1155MintableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1155Mintable *ERC1155MintableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.ERC1155MintableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1155Mintable *ERC1155MintableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1155Mintable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1155Mintable *ERC1155MintableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1155Mintable *ERC1155MintableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.contract.Transact(opts, method, params...)
}

// ERC1155BATCHRECEIVED is a free data retrieval call binding the contract method 0xfc67bf1c.
//
// Solidity: function ERC1155_BATCH_RECEIVED() view returns(bytes4)
func (_ERC1155Mintable *ERC1155MintableCaller) ERC1155BATCHRECEIVED(opts *bind.CallOpts) ([4]byte, error) {
	var out []interface{}
	err := _ERC1155Mintable.contract.Call(opts, &out, "ERC1155_BATCH_RECEIVED")

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

// ERC1155BATCHRECEIVED is a free data retrieval call binding the contract method 0xfc67bf1c.
//
// Solidity: function ERC1155_BATCH_RECEIVED() view returns(bytes4)
func (_ERC1155Mintable *ERC1155MintableSession) ERC1155BATCHRECEIVED() ([4]byte, error) {
	return _ERC1155Mintable.Contract.ERC1155BATCHRECEIVED(&_ERC1155Mintable.CallOpts)
}

// ERC1155BATCHRECEIVED is a free data retrieval call binding the contract method 0xfc67bf1c.
//
// Solidity: function ERC1155_BATCH_RECEIVED() view returns(bytes4)
func (_ERC1155Mintable *ERC1155MintableCallerSession) ERC1155BATCHRECEIVED() ([4]byte, error) {
	return _ERC1155Mintable.Contract.ERC1155BATCHRECEIVED(&_ERC1155Mintable.CallOpts)
}

// ERC1155RECEIVED is a free data retrieval call binding the contract method 0xe0a5c949.
//
// Solidity: function ERC1155_RECEIVED() view returns(bytes4)
func (_ERC1155Mintable *ERC1155MintableCaller) ERC1155RECEIVED(opts *bind.CallOpts) ([4]byte, error) {
	var out []interface{}
	err := _ERC1155Mintable.contract.Call(opts, &out, "ERC1155_RECEIVED")

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

// ERC1155RECEIVED is a free data retrieval call binding the contract method 0xe0a5c949.
//
// Solidity: function ERC1155_RECEIVED() view returns(bytes4)
func (_ERC1155Mintable *ERC1155MintableSession) ERC1155RECEIVED() ([4]byte, error) {
	return _ERC1155Mintable.Contract.ERC1155RECEIVED(&_ERC1155Mintable.CallOpts)
}

// ERC1155RECEIVED is a free data retrieval call binding the contract method 0xe0a5c949.
//
// Solidity: function ERC1155_RECEIVED() view returns(bytes4)
func (_ERC1155Mintable *ERC1155MintableCallerSession) ERC1155RECEIVED() ([4]byte, error) {
	return _ERC1155Mintable.Contract.ERC1155RECEIVED(&_ERC1155Mintable.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address owner, uint256 id) view returns(uint256)
func (_ERC1155Mintable *ERC1155MintableCaller) BalanceOf(opts *bind.CallOpts, owner common.Address, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ERC1155Mintable.contract.Call(opts, &out, "balanceOf", owner, id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address owner, uint256 id) view returns(uint256)
func (_ERC1155Mintable *ERC1155MintableSession) BalanceOf(owner common.Address, id *big.Int) (*big.Int, error) {
	return _ERC1155Mintable.Contract.BalanceOf(&_ERC1155Mintable.CallOpts, owner, id)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address owner, uint256 id) view returns(uint256)
func (_ERC1155Mintable *ERC1155MintableCallerSession) BalanceOf(owner common.Address, id *big.Int) (*big.Int, error) {
	return _ERC1155Mintable.Contract.BalanceOf(&_ERC1155Mintable.CallOpts, owner, id)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] owners, uint256[] ids) view returns(uint256[] balances_)
func (_ERC1155Mintable *ERC1155MintableCaller) BalanceOfBatch(opts *bind.CallOpts, owners []common.Address, ids []*big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _ERC1155Mintable.contract.Call(opts, &out, "balanceOfBatch", owners, ids)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] owners, uint256[] ids) view returns(uint256[] balances_)
func (_ERC1155Mintable *ERC1155MintableSession) BalanceOfBatch(owners []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _ERC1155Mintable.Contract.BalanceOfBatch(&_ERC1155Mintable.CallOpts, owners, ids)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] owners, uint256[] ids) view returns(uint256[] balances_)
func (_ERC1155Mintable *ERC1155MintableCallerSession) BalanceOfBatch(owners []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _ERC1155Mintable.Contract.BalanceOfBatch(&_ERC1155Mintable.CallOpts, owners, ids)
}

// Creators is a free data retrieval call binding the contract method 0xcd53d08e.
//
// Solidity: function creators(uint256 ) view returns(address)
func (_ERC1155Mintable *ERC1155MintableCaller) Creators(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ERC1155Mintable.contract.Call(opts, &out, "creators", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Creators is a free data retrieval call binding the contract method 0xcd53d08e.
//
// Solidity: function creators(uint256 ) view returns(address)
func (_ERC1155Mintable *ERC1155MintableSession) Creators(arg0 *big.Int) (common.Address, error) {
	return _ERC1155Mintable.Contract.Creators(&_ERC1155Mintable.CallOpts, arg0)
}

// Creators is a free data retrieval call binding the contract method 0xcd53d08e.
//
// Solidity: function creators(uint256 ) view returns(address)
func (_ERC1155Mintable *ERC1155MintableCallerSession) Creators(arg0 *big.Int) (common.Address, error) {
	return _ERC1155Mintable.Contract.Creators(&_ERC1155Mintable.CallOpts, arg0)
}

// GetNonFungibleBaseType is a free data retrieval call binding the contract method 0x6f969c2d.
//
// Solidity: function getNonFungibleBaseType(uint256 id) pure returns(uint256)
func (_ERC1155Mintable *ERC1155MintableCaller) GetNonFungibleBaseType(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ERC1155Mintable.contract.Call(opts, &out, "getNonFungibleBaseType", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNonFungibleBaseType is a free data retrieval call binding the contract method 0x6f969c2d.
//
// Solidity: function getNonFungibleBaseType(uint256 id) pure returns(uint256)
func (_ERC1155Mintable *ERC1155MintableSession) GetNonFungibleBaseType(id *big.Int) (*big.Int, error) {
	return _ERC1155Mintable.Contract.GetNonFungibleBaseType(&_ERC1155Mintable.CallOpts, id)
}

// GetNonFungibleBaseType is a free data retrieval call binding the contract method 0x6f969c2d.
//
// Solidity: function getNonFungibleBaseType(uint256 id) pure returns(uint256)
func (_ERC1155Mintable *ERC1155MintableCallerSession) GetNonFungibleBaseType(id *big.Int) (*big.Int, error) {
	return _ERC1155Mintable.Contract.GetNonFungibleBaseType(&_ERC1155Mintable.CallOpts, id)
}

// GetNonFungibleIndex is a free data retrieval call binding the contract method 0x9cca1c64.
//
// Solidity: function getNonFungibleIndex(uint256 id) pure returns(uint256)
func (_ERC1155Mintable *ERC1155MintableCaller) GetNonFungibleIndex(opts *bind.CallOpts, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ERC1155Mintable.contract.Call(opts, &out, "getNonFungibleIndex", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNonFungibleIndex is a free data retrieval call binding the contract method 0x9cca1c64.
//
// Solidity: function getNonFungibleIndex(uint256 id) pure returns(uint256)
func (_ERC1155Mintable *ERC1155MintableSession) GetNonFungibleIndex(id *big.Int) (*big.Int, error) {
	return _ERC1155Mintable.Contract.GetNonFungibleIndex(&_ERC1155Mintable.CallOpts, id)
}

// GetNonFungibleIndex is a free data retrieval call binding the contract method 0x9cca1c64.
//
// Solidity: function getNonFungibleIndex(uint256 id) pure returns(uint256)
func (_ERC1155Mintable *ERC1155MintableCallerSession) GetNonFungibleIndex(id *big.Int) (*big.Int, error) {
	return _ERC1155Mintable.Contract.GetNonFungibleIndex(&_ERC1155Mintable.CallOpts, id)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ERC1155Mintable *ERC1155MintableCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _ERC1155Mintable.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ERC1155Mintable *ERC1155MintableSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ERC1155Mintable.Contract.IsApprovedForAll(&_ERC1155Mintable.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ERC1155Mintable *ERC1155MintableCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ERC1155Mintable.Contract.IsApprovedForAll(&_ERC1155Mintable.CallOpts, owner, operator)
}

// IsFungible is a free data retrieval call binding the contract method 0xadebf6f2.
//
// Solidity: function isFungible(uint256 id) pure returns(bool)
func (_ERC1155Mintable *ERC1155MintableCaller) IsFungible(opts *bind.CallOpts, id *big.Int) (bool, error) {
	var out []interface{}
	err := _ERC1155Mintable.contract.Call(opts, &out, "isFungible", id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsFungible is a free data retrieval call binding the contract method 0xadebf6f2.
//
// Solidity: function isFungible(uint256 id) pure returns(bool)
func (_ERC1155Mintable *ERC1155MintableSession) IsFungible(id *big.Int) (bool, error) {
	return _ERC1155Mintable.Contract.IsFungible(&_ERC1155Mintable.CallOpts, id)
}

// IsFungible is a free data retrieval call binding the contract method 0xadebf6f2.
//
// Solidity: function isFungible(uint256 id) pure returns(bool)
func (_ERC1155Mintable *ERC1155MintableCallerSession) IsFungible(id *big.Int) (bool, error) {
	return _ERC1155Mintable.Contract.IsFungible(&_ERC1155Mintable.CallOpts, id)
}

// IsNonFungible is a free data retrieval call binding the contract method 0xe44591f0.
//
// Solidity: function isNonFungible(uint256 id) pure returns(bool)
func (_ERC1155Mintable *ERC1155MintableCaller) IsNonFungible(opts *bind.CallOpts, id *big.Int) (bool, error) {
	var out []interface{}
	err := _ERC1155Mintable.contract.Call(opts, &out, "isNonFungible", id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsNonFungible is a free data retrieval call binding the contract method 0xe44591f0.
//
// Solidity: function isNonFungible(uint256 id) pure returns(bool)
func (_ERC1155Mintable *ERC1155MintableSession) IsNonFungible(id *big.Int) (bool, error) {
	return _ERC1155Mintable.Contract.IsNonFungible(&_ERC1155Mintable.CallOpts, id)
}

// IsNonFungible is a free data retrieval call binding the contract method 0xe44591f0.
//
// Solidity: function isNonFungible(uint256 id) pure returns(bool)
func (_ERC1155Mintable *ERC1155MintableCallerSession) IsNonFungible(id *big.Int) (bool, error) {
	return _ERC1155Mintable.Contract.IsNonFungible(&_ERC1155Mintable.CallOpts, id)
}

// IsNonFungibleBaseType is a free data retrieval call binding the contract method 0x7269a327.
//
// Solidity: function isNonFungibleBaseType(uint256 id) pure returns(bool)
func (_ERC1155Mintable *ERC1155MintableCaller) IsNonFungibleBaseType(opts *bind.CallOpts, id *big.Int) (bool, error) {
	var out []interface{}
	err := _ERC1155Mintable.contract.Call(opts, &out, "isNonFungibleBaseType", id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsNonFungibleBaseType is a free data retrieval call binding the contract method 0x7269a327.
//
// Solidity: function isNonFungibleBaseType(uint256 id) pure returns(bool)
func (_ERC1155Mintable *ERC1155MintableSession) IsNonFungibleBaseType(id *big.Int) (bool, error) {
	return _ERC1155Mintable.Contract.IsNonFungibleBaseType(&_ERC1155Mintable.CallOpts, id)
}

// IsNonFungibleBaseType is a free data retrieval call binding the contract method 0x7269a327.
//
// Solidity: function isNonFungibleBaseType(uint256 id) pure returns(bool)
func (_ERC1155Mintable *ERC1155MintableCallerSession) IsNonFungibleBaseType(id *big.Int) (bool, error) {
	return _ERC1155Mintable.Contract.IsNonFungibleBaseType(&_ERC1155Mintable.CallOpts, id)
}

// IsNonFungibleItem is a free data retrieval call binding the contract method 0x5e81b958.
//
// Solidity: function isNonFungibleItem(uint256 id) pure returns(bool)
func (_ERC1155Mintable *ERC1155MintableCaller) IsNonFungibleItem(opts *bind.CallOpts, id *big.Int) (bool, error) {
	var out []interface{}
	err := _ERC1155Mintable.contract.Call(opts, &out, "isNonFungibleItem", id)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsNonFungibleItem is a free data retrieval call binding the contract method 0x5e81b958.
//
// Solidity: function isNonFungibleItem(uint256 id) pure returns(bool)
func (_ERC1155Mintable *ERC1155MintableSession) IsNonFungibleItem(id *big.Int) (bool, error) {
	return _ERC1155Mintable.Contract.IsNonFungibleItem(&_ERC1155Mintable.CallOpts, id)
}

// IsNonFungibleItem is a free data retrieval call binding the contract method 0x5e81b958.
//
// Solidity: function isNonFungibleItem(uint256 id) pure returns(bool)
func (_ERC1155Mintable *ERC1155MintableCallerSession) IsNonFungibleItem(id *big.Int) (bool, error) {
	return _ERC1155Mintable.Contract.IsNonFungibleItem(&_ERC1155Mintable.CallOpts, id)
}

// MaxIndex is a free data retrieval call binding the contract method 0x08d7d469.
//
// Solidity: function maxIndex(uint256 ) view returns(uint256)
func (_ERC1155Mintable *ERC1155MintableCaller) MaxIndex(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _ERC1155Mintable.contract.Call(opts, &out, "maxIndex", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxIndex is a free data retrieval call binding the contract method 0x08d7d469.
//
// Solidity: function maxIndex(uint256 ) view returns(uint256)
func (_ERC1155Mintable *ERC1155MintableSession) MaxIndex(arg0 *big.Int) (*big.Int, error) {
	return _ERC1155Mintable.Contract.MaxIndex(&_ERC1155Mintable.CallOpts, arg0)
}

// MaxIndex is a free data retrieval call binding the contract method 0x08d7d469.
//
// Solidity: function maxIndex(uint256 ) view returns(uint256)
func (_ERC1155Mintable *ERC1155MintableCallerSession) MaxIndex(arg0 *big.Int) (*big.Int, error) {
	return _ERC1155Mintable.Contract.MaxIndex(&_ERC1155Mintable.CallOpts, arg0)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 id) view returns(address)
func (_ERC1155Mintable *ERC1155MintableCaller) OwnerOf(opts *bind.CallOpts, id *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ERC1155Mintable.contract.Call(opts, &out, "ownerOf", id)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 id) view returns(address)
func (_ERC1155Mintable *ERC1155MintableSession) OwnerOf(id *big.Int) (common.Address, error) {
	return _ERC1155Mintable.Contract.OwnerOf(&_ERC1155Mintable.CallOpts, id)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 id) view returns(address)
func (_ERC1155Mintable *ERC1155MintableCallerSession) OwnerOf(id *big.Int) (common.Address, error) {
	return _ERC1155Mintable.Contract.OwnerOf(&_ERC1155Mintable.CallOpts, id)
}

// Create is a paid mutator transaction binding the contract method 0xcc10e401.
//
// Solidity: function create(string uri, bool isNF) returns(uint256 type_)
func (_ERC1155Mintable *ERC1155MintableTransactor) Create(opts *bind.TransactOpts, uri string, isNF bool) (*types.Transaction, error) {
	return _ERC1155Mintable.contract.Transact(opts, "create", uri, isNF)
}

// Create is a paid mutator transaction binding the contract method 0xcc10e401.
//
// Solidity: function create(string uri, bool isNF) returns(uint256 type_)
func (_ERC1155Mintable *ERC1155MintableSession) Create(uri string, isNF bool) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.Create(&_ERC1155Mintable.TransactOpts, uri, isNF)
}

// Create is a paid mutator transaction binding the contract method 0xcc10e401.
//
// Solidity: function create(string uri, bool isNF) returns(uint256 type_)
func (_ERC1155Mintable *ERC1155MintableTransactorSession) Create(uri string, isNF bool) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.Create(&_ERC1155Mintable.TransactOpts, uri, isNF)
}

// CreateWithType is a paid mutator transaction binding the contract method 0x9f4b286a.
//
// Solidity: function createWithType(uint256 type_, string uri) returns()
func (_ERC1155Mintable *ERC1155MintableTransactor) CreateWithType(opts *bind.TransactOpts, type_ *big.Int, uri string) (*types.Transaction, error) {
	return _ERC1155Mintable.contract.Transact(opts, "createWithType", type_, uri)
}

// CreateWithType is a paid mutator transaction binding the contract method 0x9f4b286a.
//
// Solidity: function createWithType(uint256 type_, string uri) returns()
func (_ERC1155Mintable *ERC1155MintableSession) CreateWithType(type_ *big.Int, uri string) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.CreateWithType(&_ERC1155Mintable.TransactOpts, type_, uri)
}

// CreateWithType is a paid mutator transaction binding the contract method 0x9f4b286a.
//
// Solidity: function createWithType(uint256 type_, string uri) returns()
func (_ERC1155Mintable *ERC1155MintableTransactorSession) CreateWithType(type_ *big.Int, uri string) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.CreateWithType(&_ERC1155Mintable.TransactOpts, type_, uri)
}

// MintFungible is a paid mutator transaction binding the contract method 0x78b27221.
//
// Solidity: function mintFungible(uint256 id, address[] to, uint256[] quantities) returns()
func (_ERC1155Mintable *ERC1155MintableTransactor) MintFungible(opts *bind.TransactOpts, id *big.Int, to []common.Address, quantities []*big.Int) (*types.Transaction, error) {
	return _ERC1155Mintable.contract.Transact(opts, "mintFungible", id, to, quantities)
}

// MintFungible is a paid mutator transaction binding the contract method 0x78b27221.
//
// Solidity: function mintFungible(uint256 id, address[] to, uint256[] quantities) returns()
func (_ERC1155Mintable *ERC1155MintableSession) MintFungible(id *big.Int, to []common.Address, quantities []*big.Int) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.MintFungible(&_ERC1155Mintable.TransactOpts, id, to, quantities)
}

// MintFungible is a paid mutator transaction binding the contract method 0x78b27221.
//
// Solidity: function mintFungible(uint256 id, address[] to, uint256[] quantities) returns()
func (_ERC1155Mintable *ERC1155MintableTransactorSession) MintFungible(id *big.Int, to []common.Address, quantities []*big.Int) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.MintFungible(&_ERC1155Mintable.TransactOpts, id, to, quantities)
}

// MintNonFungible is a paid mutator transaction binding the contract method 0xf9419088.
//
// Solidity: function mintNonFungible(uint256 type_, address[] to) returns()
func (_ERC1155Mintable *ERC1155MintableTransactor) MintNonFungible(opts *bind.TransactOpts, type_ *big.Int, to []common.Address) (*types.Transaction, error) {
	return _ERC1155Mintable.contract.Transact(opts, "mintNonFungible", type_, to)
}

// MintNonFungible is a paid mutator transaction binding the contract method 0xf9419088.
//
// Solidity: function mintNonFungible(uint256 type_, address[] to) returns()
func (_ERC1155Mintable *ERC1155MintableSession) MintNonFungible(type_ *big.Int, to []common.Address) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.MintNonFungible(&_ERC1155Mintable.TransactOpts, type_, to)
}

// MintNonFungible is a paid mutator transaction binding the contract method 0xf9419088.
//
// Solidity: function mintNonFungible(uint256 type_, address[] to) returns()
func (_ERC1155Mintable *ERC1155MintableTransactorSession) MintNonFungible(type_ *big.Int, to []common.Address) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.MintNonFungible(&_ERC1155Mintable.TransactOpts, type_, to)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] values, bytes data) returns()
func (_ERC1155Mintable *ERC1155MintableTransactor) SafeBatchTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, ids []*big.Int, values []*big.Int, data []byte) (*types.Transaction, error) {
	return _ERC1155Mintable.contract.Transact(opts, "safeBatchTransferFrom", from, to, ids, values, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] values, bytes data) returns()
func (_ERC1155Mintable *ERC1155MintableSession) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, values []*big.Int, data []byte) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.SafeBatchTransferFrom(&_ERC1155Mintable.TransactOpts, from, to, ids, values, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] values, bytes data) returns()
func (_ERC1155Mintable *ERC1155MintableTransactorSession) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, values []*big.Int, data []byte) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.SafeBatchTransferFrom(&_ERC1155Mintable.TransactOpts, from, to, ids, values, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes data) returns()
func (_ERC1155Mintable *ERC1155MintableTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, id *big.Int, value *big.Int, data []byte) (*types.Transaction, error) {
	return _ERC1155Mintable.contract.Transact(opts, "safeTransferFrom", from, to, id, value, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes data) returns()
func (_ERC1155Mintable *ERC1155MintableSession) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, value *big.Int, data []byte) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.SafeTransferFrom(&_ERC1155Mintable.TransactOpts, from, to, id, value, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 value, bytes data) returns()
func (_ERC1155Mintable *ERC1155MintableTransactorSession) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, value *big.Int, data []byte) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.SafeTransferFrom(&_ERC1155Mintable.TransactOpts, from, to, id, value, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ERC1155Mintable *ERC1155MintableTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _ERC1155Mintable.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ERC1155Mintable *ERC1155MintableSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.SetApprovalForAll(&_ERC1155Mintable.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ERC1155Mintable *ERC1155MintableTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _ERC1155Mintable.Contract.SetApprovalForAll(&_ERC1155Mintable.TransactOpts, operator, approved)
}

// ERC1155MintableApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the ERC1155Mintable contract.
type ERC1155MintableApprovalForAllIterator struct {
	Event *ERC1155MintableApprovalForAll // Event containing the contract specifics and raw log

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
func (it *ERC1155MintableApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1155MintableApprovalForAll)
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
		it.Event = new(ERC1155MintableApprovalForAll)
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
func (it *ERC1155MintableApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1155MintableApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1155MintableApprovalForAll represents a ApprovalForAll event raised by the ERC1155Mintable contract.
type ERC1155MintableApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ERC1155Mintable *ERC1155MintableFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*ERC1155MintableApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ERC1155Mintable.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &ERC1155MintableApprovalForAllIterator{contract: _ERC1155Mintable.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ERC1155Mintable *ERC1155MintableFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *ERC1155MintableApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ERC1155Mintable.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1155MintableApprovalForAll)
				if err := _ERC1155Mintable.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ERC1155Mintable *ERC1155MintableFilterer) ParseApprovalForAll(log types.Log) (*ERC1155MintableApprovalForAll, error) {
	event := new(ERC1155MintableApprovalForAll)
	if err := _ERC1155Mintable.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC1155MintableTransferBatchIterator is returned from FilterTransferBatch and is used to iterate over the raw logs and unpacked data for TransferBatch events raised by the ERC1155Mintable contract.
type ERC1155MintableTransferBatchIterator struct {
	Event *ERC1155MintableTransferBatch // Event containing the contract specifics and raw log

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
func (it *ERC1155MintableTransferBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1155MintableTransferBatch)
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
		it.Event = new(ERC1155MintableTransferBatch)
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
func (it *ERC1155MintableTransferBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1155MintableTransferBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1155MintableTransferBatch represents a TransferBatch event raised by the ERC1155Mintable contract.
type ERC1155MintableTransferBatch struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Ids      []*big.Int
	Values   []*big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferBatch is a free log retrieval operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_ERC1155Mintable *ERC1155MintableFilterer) FilterTransferBatch(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*ERC1155MintableTransferBatchIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC1155Mintable.contract.FilterLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC1155MintableTransferBatchIterator{contract: _ERC1155Mintable.contract, event: "TransferBatch", logs: logs, sub: sub}, nil
}

// WatchTransferBatch is a free log subscription operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_ERC1155Mintable *ERC1155MintableFilterer) WatchTransferBatch(opts *bind.WatchOpts, sink chan<- *ERC1155MintableTransferBatch, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC1155Mintable.contract.WatchLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1155MintableTransferBatch)
				if err := _ERC1155Mintable.contract.UnpackLog(event, "TransferBatch", log); err != nil {
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

// ParseTransferBatch is a log parse operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_ERC1155Mintable *ERC1155MintableFilterer) ParseTransferBatch(log types.Log) (*ERC1155MintableTransferBatch, error) {
	event := new(ERC1155MintableTransferBatch)
	if err := _ERC1155Mintable.contract.UnpackLog(event, "TransferBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC1155MintableTransferSingleIterator is returned from FilterTransferSingle and is used to iterate over the raw logs and unpacked data for TransferSingle events raised by the ERC1155Mintable contract.
type ERC1155MintableTransferSingleIterator struct {
	Event *ERC1155MintableTransferSingle // Event containing the contract specifics and raw log

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
func (it *ERC1155MintableTransferSingleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1155MintableTransferSingle)
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
		it.Event = new(ERC1155MintableTransferSingle)
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
func (it *ERC1155MintableTransferSingleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1155MintableTransferSingleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1155MintableTransferSingle represents a TransferSingle event raised by the ERC1155Mintable contract.
type ERC1155MintableTransferSingle struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Id       *big.Int
	Value    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferSingle is a free log retrieval operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_ERC1155Mintable *ERC1155MintableFilterer) FilterTransferSingle(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*ERC1155MintableTransferSingleIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC1155Mintable.contract.FilterLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC1155MintableTransferSingleIterator{contract: _ERC1155Mintable.contract, event: "TransferSingle", logs: logs, sub: sub}, nil
}

// WatchTransferSingle is a free log subscription operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_ERC1155Mintable *ERC1155MintableFilterer) WatchTransferSingle(opts *bind.WatchOpts, sink chan<- *ERC1155MintableTransferSingle, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC1155Mintable.contract.WatchLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1155MintableTransferSingle)
				if err := _ERC1155Mintable.contract.UnpackLog(event, "TransferSingle", log); err != nil {
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

// ParseTransferSingle is a log parse operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_ERC1155Mintable *ERC1155MintableFilterer) ParseTransferSingle(log types.Log) (*ERC1155MintableTransferSingle, error) {
	event := new(ERC1155MintableTransferSingle)
	if err := _ERC1155Mintable.contract.UnpackLog(event, "TransferSingle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC1155MintableURIIterator is returned from FilterURI and is used to iterate over the raw logs and unpacked data for URI events raised by the ERC1155Mintable contract.
type ERC1155MintableURIIterator struct {
	Event *ERC1155MintableURI // Event containing the contract specifics and raw log

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
func (it *ERC1155MintableURIIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1155MintableURI)
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
		it.Event = new(ERC1155MintableURI)
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
func (it *ERC1155MintableURIIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1155MintableURIIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1155MintableURI represents a URI event raised by the ERC1155Mintable contract.
type ERC1155MintableURI struct {
	Value string
	Id    *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterURI is a free log retrieval operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_ERC1155Mintable *ERC1155MintableFilterer) FilterURI(opts *bind.FilterOpts, id []*big.Int) (*ERC1155MintableURIIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _ERC1155Mintable.contract.FilterLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return &ERC1155MintableURIIterator{contract: _ERC1155Mintable.contract, event: "URI", logs: logs, sub: sub}, nil
}

// WatchURI is a free log subscription operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_ERC1155Mintable *ERC1155MintableFilterer) WatchURI(opts *bind.WatchOpts, sink chan<- *ERC1155MintableURI, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _ERC1155Mintable.contract.WatchLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1155MintableURI)
				if err := _ERC1155Mintable.contract.UnpackLog(event, "URI", log); err != nil {
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

// ParseURI is a log parse operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_ERC1155Mintable *ERC1155MintableFilterer) ParseURI(log types.Log) (*ERC1155MintableURI, error) {
	event := new(ERC1155MintableURI)
	if err := _ERC1155Mintable.contract.UnpackLog(event, "URI", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
