// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package goBind

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
	_ = abi.ConvertType
)

// FogmetaRebuilderMetaData contains all meta data concerning the FogmetaRebuilder contract.
var FogmetaRebuilderMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"requested\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"available\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AddressBalance\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transferToken\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// FogmetaRebuilderABI is the input ABI used to generate the binding from.
// Deprecated: Use FogmetaRebuilderMetaData.ABI instead.
var FogmetaRebuilderABI = FogmetaRebuilderMetaData.ABI

// FogmetaRebuilder is an auto generated Go binding around an Ethereum contract.
type FogmetaRebuilder struct {
	FogmetaRebuilderCaller     // Read-only binding to the contract
	FogmetaRebuilderTransactor // Write-only binding to the contract
	FogmetaRebuilderFilterer   // Log filterer for contract events
}

// FogmetaRebuilderCaller is an auto generated read-only Go binding around an Ethereum contract.
type FogmetaRebuilderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FogmetaRebuilderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FogmetaRebuilderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FogmetaRebuilderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FogmetaRebuilderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FogmetaRebuilderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FogmetaRebuilderSession struct {
	Contract     *FogmetaRebuilder // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FogmetaRebuilderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FogmetaRebuilderCallerSession struct {
	Contract *FogmetaRebuilderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// FogmetaRebuilderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FogmetaRebuilderTransactorSession struct {
	Contract     *FogmetaRebuilderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// FogmetaRebuilderRaw is an auto generated low-level Go binding around an Ethereum contract.
type FogmetaRebuilderRaw struct {
	Contract *FogmetaRebuilder // Generic contract binding to access the raw methods on
}

// FogmetaRebuilderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FogmetaRebuilderCallerRaw struct {
	Contract *FogmetaRebuilderCaller // Generic read-only contract binding to access the raw methods on
}

// FogmetaRebuilderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FogmetaRebuilderTransactorRaw struct {
	Contract *FogmetaRebuilderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFogmetaRebuilder creates a new instance of FogmetaRebuilder, bound to a specific deployed contract.
func NewFogmetaRebuilder(address common.Address, backend bind.ContractBackend) (*FogmetaRebuilder, error) {
	contract, err := bindFogmetaRebuilder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FogmetaRebuilder{FogmetaRebuilderCaller: FogmetaRebuilderCaller{contract: contract}, FogmetaRebuilderTransactor: FogmetaRebuilderTransactor{contract: contract}, FogmetaRebuilderFilterer: FogmetaRebuilderFilterer{contract: contract}}, nil
}

// NewFogmetaRebuilderCaller creates a new read-only instance of FogmetaRebuilder, bound to a specific deployed contract.
func NewFogmetaRebuilderCaller(address common.Address, caller bind.ContractCaller) (*FogmetaRebuilderCaller, error) {
	contract, err := bindFogmetaRebuilder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FogmetaRebuilderCaller{contract: contract}, nil
}

// NewFogmetaRebuilderTransactor creates a new write-only instance of FogmetaRebuilder, bound to a specific deployed contract.
func NewFogmetaRebuilderTransactor(address common.Address, transactor bind.ContractTransactor) (*FogmetaRebuilderTransactor, error) {
	contract, err := bindFogmetaRebuilder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FogmetaRebuilderTransactor{contract: contract}, nil
}

// NewFogmetaRebuilderFilterer creates a new log filterer instance of FogmetaRebuilder, bound to a specific deployed contract.
func NewFogmetaRebuilderFilterer(address common.Address, filterer bind.ContractFilterer) (*FogmetaRebuilderFilterer, error) {
	contract, err := bindFogmetaRebuilder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FogmetaRebuilderFilterer{contract: contract}, nil
}

// bindFogmetaRebuilder binds a generic wrapper to an already deployed contract.
func bindFogmetaRebuilder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FogmetaRebuilderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FogmetaRebuilder *FogmetaRebuilderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FogmetaRebuilder.Contract.FogmetaRebuilderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FogmetaRebuilder *FogmetaRebuilderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FogmetaRebuilder.Contract.FogmetaRebuilderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FogmetaRebuilder *FogmetaRebuilderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FogmetaRebuilder.Contract.FogmetaRebuilderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FogmetaRebuilder *FogmetaRebuilderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FogmetaRebuilder.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FogmetaRebuilder *FogmetaRebuilderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FogmetaRebuilder.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FogmetaRebuilder *FogmetaRebuilderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FogmetaRebuilder.Contract.contract.Transact(opts, method, params...)
}

// GetBalance is a free data retrieval call binding the contract method 0x12065fe0.
//
// Solidity: function getBalance() view returns(uint256)
func (_FogmetaRebuilder *FogmetaRebuilderCaller) GetBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FogmetaRebuilder.contract.Call(opts, &out, "getBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalance is a free data retrieval call binding the contract method 0x12065fe0.
//
// Solidity: function getBalance() view returns(uint256)
func (_FogmetaRebuilder *FogmetaRebuilderSession) GetBalance() (*big.Int, error) {
	return _FogmetaRebuilder.Contract.GetBalance(&_FogmetaRebuilder.CallOpts)
}

// GetBalance is a free data retrieval call binding the contract method 0x12065fe0.
//
// Solidity: function getBalance() view returns(uint256)
func (_FogmetaRebuilder *FogmetaRebuilderCallerSession) GetBalance() (*big.Int, error) {
	return _FogmetaRebuilder.Contract.GetBalance(&_FogmetaRebuilder.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FogmetaRebuilder *FogmetaRebuilderCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FogmetaRebuilder.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FogmetaRebuilder *FogmetaRebuilderSession) Owner() (common.Address, error) {
	return _FogmetaRebuilder.Contract.Owner(&_FogmetaRebuilder.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FogmetaRebuilder *FogmetaRebuilderCallerSession) Owner() (common.Address, error) {
	return _FogmetaRebuilder.Contract.Owner(&_FogmetaRebuilder.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FogmetaRebuilder *FogmetaRebuilderTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FogmetaRebuilder.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FogmetaRebuilder *FogmetaRebuilderSession) RenounceOwnership() (*types.Transaction, error) {
	return _FogmetaRebuilder.Contract.RenounceOwnership(&_FogmetaRebuilder.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FogmetaRebuilder *FogmetaRebuilderTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _FogmetaRebuilder.Contract.RenounceOwnership(&_FogmetaRebuilder.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0x12514bba.
//
// Solidity: function transfer(uint256 amount) returns()
func (_FogmetaRebuilder *FogmetaRebuilderTransactor) Transfer(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _FogmetaRebuilder.contract.Transact(opts, "transfer", amount)
}

// Transfer is a paid mutator transaction binding the contract method 0x12514bba.
//
// Solidity: function transfer(uint256 amount) returns()
func (_FogmetaRebuilder *FogmetaRebuilderSession) Transfer(amount *big.Int) (*types.Transaction, error) {
	return _FogmetaRebuilder.Contract.Transfer(&_FogmetaRebuilder.TransactOpts, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0x12514bba.
//
// Solidity: function transfer(uint256 amount) returns()
func (_FogmetaRebuilder *FogmetaRebuilderTransactorSession) Transfer(amount *big.Int) (*types.Transaction, error) {
	return _FogmetaRebuilder.Contract.Transfer(&_FogmetaRebuilder.TransactOpts, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FogmetaRebuilder *FogmetaRebuilderTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _FogmetaRebuilder.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FogmetaRebuilder *FogmetaRebuilderSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FogmetaRebuilder.Contract.TransferOwnership(&_FogmetaRebuilder.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FogmetaRebuilder *FogmetaRebuilderTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FogmetaRebuilder.Contract.TransferOwnership(&_FogmetaRebuilder.TransactOpts, newOwner)
}

// TransferToken is a paid mutator transaction binding the contract method 0x799a5359.
//
// Solidity: function transferToken() payable returns()
func (_FogmetaRebuilder *FogmetaRebuilderTransactor) TransferToken(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FogmetaRebuilder.contract.Transact(opts, "transferToken")
}

// TransferToken is a paid mutator transaction binding the contract method 0x799a5359.
//
// Solidity: function transferToken() payable returns()
func (_FogmetaRebuilder *FogmetaRebuilderSession) TransferToken() (*types.Transaction, error) {
	return _FogmetaRebuilder.Contract.TransferToken(&_FogmetaRebuilder.TransactOpts)
}

// TransferToken is a paid mutator transaction binding the contract method 0x799a5359.
//
// Solidity: function transferToken() payable returns()
func (_FogmetaRebuilder *FogmetaRebuilderTransactorSession) TransferToken() (*types.Transaction, error) {
	return _FogmetaRebuilder.Contract.TransferToken(&_FogmetaRebuilder.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns(bool)
func (_FogmetaRebuilder *FogmetaRebuilderTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _FogmetaRebuilder.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns(bool)
func (_FogmetaRebuilder *FogmetaRebuilderSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _FogmetaRebuilder.Contract.Withdraw(&_FogmetaRebuilder.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns(bool)
func (_FogmetaRebuilder *FogmetaRebuilderTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _FogmetaRebuilder.Contract.Withdraw(&_FogmetaRebuilder.TransactOpts, amount)
}

// FogmetaRebuilderAddressBalanceIterator is returned from FilterAddressBalance and is used to iterate over the raw logs and unpacked data for AddressBalance events raised by the FogmetaRebuilder contract.
type FogmetaRebuilderAddressBalanceIterator struct {
	Event *FogmetaRebuilderAddressBalance // Event containing the contract specifics and raw log

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
func (it *FogmetaRebuilderAddressBalanceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FogmetaRebuilderAddressBalance)
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
		it.Event = new(FogmetaRebuilderAddressBalance)
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
func (it *FogmetaRebuilderAddressBalanceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FogmetaRebuilderAddressBalanceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FogmetaRebuilderAddressBalance represents a AddressBalance event raised by the FogmetaRebuilder contract.
type FogmetaRebuilderAddressBalance struct {
	From   common.Address
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterAddressBalance is a free log retrieval operation binding the contract event 0x8df3abdf0e7448cbffada332013b25f9a504477115ce3ab0a3fec5e22f02fba5.
//
// Solidity: event AddressBalance(address from, address to, uint256 amount)
func (_FogmetaRebuilder *FogmetaRebuilderFilterer) FilterAddressBalance(opts *bind.FilterOpts) (*FogmetaRebuilderAddressBalanceIterator, error) {

	logs, sub, err := _FogmetaRebuilder.contract.FilterLogs(opts, "AddressBalance")
	if err != nil {
		return nil, err
	}
	return &FogmetaRebuilderAddressBalanceIterator{contract: _FogmetaRebuilder.contract, event: "AddressBalance", logs: logs, sub: sub}, nil
}

// WatchAddressBalance is a free log subscription operation binding the contract event 0x8df3abdf0e7448cbffada332013b25f9a504477115ce3ab0a3fec5e22f02fba5.
//
// Solidity: event AddressBalance(address from, address to, uint256 amount)
func (_FogmetaRebuilder *FogmetaRebuilderFilterer) WatchAddressBalance(opts *bind.WatchOpts, sink chan<- *FogmetaRebuilderAddressBalance) (event.Subscription, error) {

	logs, sub, err := _FogmetaRebuilder.contract.WatchLogs(opts, "AddressBalance")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FogmetaRebuilderAddressBalance)
				if err := _FogmetaRebuilder.contract.UnpackLog(event, "AddressBalance", log); err != nil {
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

// ParseAddressBalance is a log parse operation binding the contract event 0x8df3abdf0e7448cbffada332013b25f9a504477115ce3ab0a3fec5e22f02fba5.
//
// Solidity: event AddressBalance(address from, address to, uint256 amount)
func (_FogmetaRebuilder *FogmetaRebuilderFilterer) ParseAddressBalance(log types.Log) (*FogmetaRebuilderAddressBalance, error) {
	event := new(FogmetaRebuilderAddressBalance)
	if err := _FogmetaRebuilder.contract.UnpackLog(event, "AddressBalance", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FogmetaRebuilderOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the FogmetaRebuilder contract.
type FogmetaRebuilderOwnershipTransferredIterator struct {
	Event *FogmetaRebuilderOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *FogmetaRebuilderOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FogmetaRebuilderOwnershipTransferred)
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
		it.Event = new(FogmetaRebuilderOwnershipTransferred)
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
func (it *FogmetaRebuilderOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FogmetaRebuilderOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FogmetaRebuilderOwnershipTransferred represents a OwnershipTransferred event raised by the FogmetaRebuilder contract.
type FogmetaRebuilderOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FogmetaRebuilder *FogmetaRebuilderFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FogmetaRebuilderOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FogmetaRebuilder.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FogmetaRebuilderOwnershipTransferredIterator{contract: _FogmetaRebuilder.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FogmetaRebuilder *FogmetaRebuilderFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FogmetaRebuilderOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FogmetaRebuilder.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FogmetaRebuilderOwnershipTransferred)
				if err := _FogmetaRebuilder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_FogmetaRebuilder *FogmetaRebuilderFilterer) ParseOwnershipTransferred(log types.Log) (*FogmetaRebuilderOwnershipTransferred, error) {
	event := new(FogmetaRebuilderOwnershipTransferred)
	if err := _FogmetaRebuilder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
