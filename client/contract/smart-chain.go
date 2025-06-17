// Code generated - DO NOT EDIT.

package contract

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

var SmartMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var SmartABI = SmartMetaData.ABI

type Smart struct {
	SmartCaller     
	SmartTransactor 
	SmartFilterer   
}

type SmartCaller struct {
	contract *bind.BoundContract 
}

type SmartTransactor struct {
	contract *bind.BoundContract 
}

type SmartFilterer struct {
	contract *bind.BoundContract 
}

type SmartSession struct {
	Contract     *Smart            
	CallOpts     bind.CallOpts     
	TransactOpts bind.TransactOpts 
}

type SmartCallerSession struct {
	Contract *SmartCaller  
	CallOpts bind.CallOpts 
}

type SmartTransactorSession struct {
	Contract     *SmartTransactor  
	TransactOpts bind.TransactOpts 
}

type SmartRaw struct {
	Contract *Smart 
}

type SmartCallerRaw struct {
	Contract *SmartCaller 
}

type SmartTransactorRaw struct {
	Contract *SmartTransactor 
}

func NewSmart(address common.Address, backend bind.ContractBackend) (*Smart, error) {
	contract, err := bindSmart(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Smart{SmartCaller: SmartCaller{contract: contract}, SmartTransactor: SmartTransactor{contract: contract}, SmartFilterer: SmartFilterer{contract: contract}}, nil
}

func NewSmartCaller(address common.Address, caller bind.ContractCaller) (*SmartCaller, error) {
	contract, err := bindSmart(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SmartCaller{contract: contract}, nil
}

func NewSmartTransactor(address common.Address, transactor bind.ContractTransactor) (*SmartTransactor, error) {
	contract, err := bindSmart(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SmartTransactor{contract: contract}, nil
}

func NewSmartFilterer(address common.Address, filterer bind.ContractFilterer) (*SmartFilterer, error) {
	contract, err := bindSmart(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SmartFilterer{contract: contract}, nil
}

func bindSmart(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SmartABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

func (_Smart *SmartRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Smart.Contract.SmartCaller.contract.Call(opts, result, method, params...)
}

func (_Smart *SmartRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Smart.Contract.SmartTransactor.contract.Transfer(opts)
}

func (_Smart *SmartRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Smart.Contract.SmartTransactor.contract.Transact(opts, method, params...)
}

func (_Smart *SmartCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Smart.Contract.contract.Call(opts, result, method, params...)
}

func (_Smart *SmartTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Smart.Contract.contract.Transfer(opts)
}

func (_Smart *SmartTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Smart.Contract.contract.Transact(opts, method, params...)
}


func (_Smart *SmartCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Smart.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}


func (_Smart *SmartSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Smart.Contract.Allowance(&_Smart.CallOpts, owner, spender)
}


func (_Smart *SmartCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Smart.Contract.Allowance(&_Smart.CallOpts, owner, spender)
}


func (_Smart *SmartCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Smart.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}


func (_Smart *SmartSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Smart.Contract.BalanceOf(&_Smart.CallOpts, account)
}


func (_Smart *SmartCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Smart.Contract.BalanceOf(&_Smart.CallOpts, account)
}


func (_Smart *SmartCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Smart.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}


func (_Smart *SmartSession) Decimals() (uint8, error) {
	return _Smart.Contract.Decimals(&_Smart.CallOpts)
}


func (_Smart *SmartCallerSession) Decimals() (uint8, error) {
	return _Smart.Contract.Decimals(&_Smart.CallOpts)
}


func (_Smart *SmartCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Smart.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}


func (_Smart *SmartSession) Name() (string, error) {
	return _Smart.Contract.Name(&_Smart.CallOpts)
}


func (_Smart *SmartCallerSession) Name() (string, error) {
	return _Smart.Contract.Name(&_Smart.CallOpts)
}


func (_Smart *SmartCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Smart.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}


func (_Smart *SmartSession) Symbol() (string, error) {
	return _Smart.Contract.Symbol(&_Smart.CallOpts)
}


func (_Smart *SmartCallerSession) Symbol() (string, error) {
	return _Smart.Contract.Symbol(&_Smart.CallOpts)
}


func (_Smart *SmartCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Smart.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}


func (_Smart *SmartSession) TotalSupply() (*big.Int, error) {
	return _Smart.Contract.TotalSupply(&_Smart.CallOpts)
}


func (_Smart *SmartCallerSession) TotalSupply() (*big.Int, error) {
	return _Smart.Contract.TotalSupply(&_Smart.CallOpts)
}


func (_Smart *SmartTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Smart.contract.Transact(opts, "approve", spender, amount)
}


func (_Smart *SmartSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Smart.Contract.Approve(&_Smart.TransactOpts, spender, amount)
}


func (_Smart *SmartTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Smart.Contract.Approve(&_Smart.TransactOpts, spender, amount)
}


func (_Smart *SmartTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Smart.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}


func (_Smart *SmartSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Smart.Contract.DecreaseAllowance(&_Smart.TransactOpts, spender, subtractedValue)
}


func (_Smart *SmartTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Smart.Contract.DecreaseAllowance(&_Smart.TransactOpts, spender, subtractedValue)
}


func (_Smart *SmartTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Smart.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}


func (_Smart *SmartSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Smart.Contract.IncreaseAllowance(&_Smart.TransactOpts, spender, addedValue)
}


func (_Smart *SmartTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Smart.Contract.IncreaseAllowance(&_Smart.TransactOpts, spender, addedValue)
}


func (_Smart *SmartTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Smart.contract.Transact(opts, "transfer", recipient, amount)
}


func (_Smart *SmartSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Smart.Contract.Transfer(&_Smart.TransactOpts, recipient, amount)
}


func (_Smart *SmartTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Smart.Contract.Transfer(&_Smart.TransactOpts, recipient, amount)
}


func (_Smart *SmartTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Smart.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}


func (_Smart *SmartSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Smart.Contract.TransferFrom(&_Smart.TransactOpts, sender, recipient, amount)
}


func (_Smart *SmartTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Smart.Contract.TransferFrom(&_Smart.TransactOpts, sender, recipient, amount)
}

type SmartApprovalIterator struct {
	Event *SmartApproval 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}

func (it *SmartApprovalIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SmartApproval)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(SmartApproval)
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

func (it *SmartApprovalIterator) Error() error {
	return it.fail
}

func (it *SmartApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SmartApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log 
}


func (_Smart *SmartFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*SmartApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Smart.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &SmartApprovalIterator{contract: _Smart.contract, event: "Approval", logs: logs, sub: sub}, nil
}


func (_Smart *SmartFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *SmartApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Smart.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(SmartApproval)
				if err := _Smart.contract.UnpackLog(event, "Approval", log); err != nil {
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


func (_Smart *SmartFilterer) ParseApproval(log types.Log) (*SmartApproval, error) {
	event := new(SmartApproval)
	if err := _Smart.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

type SmartTransferIterator struct {
	Event *SmartTransfer 

	contract *bind.BoundContract 
	event    string              

	logs chan types.Log        
	sub  ethereum.Subscription 
	done bool                  
	fail error                 
}

func (it *SmartTransferIterator) Next() bool {
	
	if it.fail != nil {
		return false
	}
	
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SmartTransfer)
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
	
	select {
	case log := <-it.logs:
		it.Event = new(SmartTransfer)
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

func (it *SmartTransferIterator) Error() error {
	return it.fail
}

func (it *SmartTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

type SmartTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log 
}


func (_Smart *SmartFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SmartTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Smart.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SmartTransferIterator{contract: _Smart.contract, event: "Transfer", logs: logs, sub: sub}, nil
}


func (_Smart *SmartFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *SmartTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Smart.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				
				event := new(SmartTransfer)
				if err := _Smart.contract.UnpackLog(event, "Transfer", log); err != nil {
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


func (_Smart *SmartFilterer) ParseTransfer(log types.Log) (*SmartTransfer, error) {
	event := new(SmartTransfer)
	if err := _Smart.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
