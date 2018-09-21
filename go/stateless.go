package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ETCDEVTeam/sputnikvm-ffi/go/sputnikvm"
	"github.com/ethereumproject/go-ethereum/common"
)

type storageItem struct {
	key   big.Int
	value big.Int
}

type account struct {
	address common.Address
	balance big.Int
	code    []byte
	storage []storageItem
}

var accounts []account

var addrcaller common.Address
var callerStr string

var addrCallee common.Address
var calleeStr string

var addrbenef common.Address
var benefStr string

var caleeiscreated bool

func runStatelessSample() {

	//Test Smart Contract
	/*
		pragma solidity ^0.4.24;

		contract SimpleStorage {

			uint private _balance;
			uint private _storedData;
			event notifyStorage(uint x);

			constructor() public payable {
				_storedData = 0x12d687;
				_balance = 1500000;
			}

			function set(uint x) public payable {
				_storedData = x;
				emit notifyStorage(x);
			}

			function get() public view returns (uint) {
				return _storedData;
			}
		}
	*/

	deployCode, _ := hex.DecodeString("60806040526212d6876001556216e36060005560e9806100206000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146058575b600080fd5b6056600435607c565b005b348015606357600080fd5b50606a60b7565b60408051918252519081900360200190f35b60018190556040805182815290517f23f9887eb044d32dba99d7b0b753c61c3c3b72d70ff0addb9a843542fd7642129181900360200190a150565b600154905600a165627a7a7230582013452558cc58a514b8056c0b45a3f1ab8c5f736b2e087c65e615650b562415ff0029")

	//Setup Some Test Accounts
	fmt.Println("\n========================================\nSetup Accounts...\n========================================")
	setupAccounts()

	//deploy
	fmt.Println("\n========================================\nDeploying...\n========================================")
	setAccountCode(addrCallee, deployCode)
	executeVM(addrcaller, addrCallee, addrbenef, 0, 1000000, 0, 0, 0, 0, 0, 1000000, deployCode, true)

	//Call Set Method
	fmt.Println("\n========================================\nCall Set Method by 1234567...\n========================================")
	setMethod, _ := hex.DecodeString("60fe47b100000000000000000000000000000000000000000000000000000000000001c8")
	executeVM(addrcaller, addrCallee, addrbenef, 0, 1000000, 10, 1, 1, 0, 0, 1000000, setMethod, false)

	//Call Get Method
	fmt.Println("\n========================================\nCall Get Method...\n========================================")
	getMethod, _ := hex.DecodeString("6d4ce63c")
	executeVM(addrcaller, addrCallee, addrbenef, 0, 1000000, 0, 1, 1, 0, 0, 1000000, getMethod, false)
}

func setupAccounts() {

	callerbytes, _ := hex.DecodeString("d282d4146c77e6a586c005d7bedf0d137a8b8eb0")
	addrcaller.SetBytes(callerbytes)
	callerStr = hex.EncodeToString(callerbytes)

	calleebytes, _ := hex.DecodeString("6439623964626135363661663438626138643438")
	addrCallee.SetBytes(calleebytes)
	calleeStr = hex.EncodeToString(calleebytes)

	benefbytes, _ := hex.DecodeString("4a9274e0d453b4bbb6694065e561d546da48bbbd")
	addrbenef.SetBytes(benefbytes)
	benefStr = hex.EncodeToString(benefbytes)

	fmt.Println("caller: ", hex.EncodeToString(addrcaller.Bytes()))
	fmt.Println("callee: ", hex.EncodeToString(addrCallee.Bytes()))
	fmt.Println("benef: ", hex.EncodeToString(addrbenef.Bytes()))

	accounts = make([]account, 3)

	accounts[0].address = addrcaller
	accounts[0].balance = *new(big.Int).SetUint64(250000)
	accounts[0].code = []byte{}

	accounts[1].address = addrCallee
	accounts[1].balance = *new(big.Int).SetUint64(150000)
	accounts[1].code = []byte{}

	accounts[2].address = addrbenef
	accounts[2].balance = *new(big.Int).SetUint64(320000)
	accounts[2].code = []byte{}
	fmt.Println("3 Accounts are ready now!")
}

func isEqualAddress(addr1 common.Address, addr2 common.Address) bool {
	a1 := hex.EncodeToString(addr1.Bytes())
	a2 := hex.EncodeToString(addr2.Bytes())
	return a1 == a2
}

func getAccountByAddress(addr common.Address) *account {
	ret := new(account)

	for n := 0; n < 3; n++ {
		if isEqualAddress(accounts[n].address, addr) {
			return &accounts[n]
		}
	}
	return ret
}

func isExist(addr common.Address) bool {
	for n := 0; n < 3; n++ {
		if isEqualAddress(accounts[n].address, addr) {
			return true
		}
	}
	return false
}

func setAccountCode(addr common.Address, code []byte) {
	for n := 0; n < 3; n++ {
		if isEqualAddress(accounts[n].address, addr) {
			kod := hex.EncodeToString(code)
			accounts[n].code, _ = hex.DecodeString(kod)
			break
		}
	}
}

func getAccountStorageValue(addr common.Address, key big.Int) *big.Int {
	ret := new(big.Int).SetInt64(0)

	for n := 0; n < 3; n++ {
		if isEqualAddress(accounts[n].address, addr) {
			nstorage := len(accounts[n].storage)
			for ns := 0; ns < nstorage; ns++ {
				if hex.EncodeToString(accounts[n].storage[ns].key.Bytes()) == hex.EncodeToString(key.Bytes()) {
					return &accounts[n].storage[ns].value
				}
			}
		}
	}
	return ret
}

func setAccountStorageValue(addr common.Address, key big.Int, value big.Int) {
	for n := 0; n < 3; n++ {
		if isEqualAddress(accounts[n].address, addr) {
			nstorage := len(accounts[n].storage)
			for ns := 0; ns < nstorage; ns++ {
				if hex.EncodeToString(accounts[n].storage[ns].key.Bytes()) == hex.EncodeToString(key.Bytes()) {
					accounts[n].storage[ns].value = value
					return
				}
			}
			addToAccountStorage(addr, key, value)
		}
	}
}

func addToAccountStorage(addr common.Address, key big.Int, value big.Int) {
	for n := 0; n < 3; n++ {
		if isEqualAddress(accounts[n].address, addr) {
			accounts[n].storage = append(accounts[n].storage, storageItem{key, value})
			break
		}
	}
}

func getBlockHash(blocknumber *big.Int) common.Hash {
	return *new(common.Hash)
}

func addBalance(addr common.Address, amount *big.Int) {
	for n := 0; n < 3; n++ {
		if isEqualAddress(accounts[n].address, addr) {
			accounts[n].balance.Add(&accounts[n].balance, amount)
			break
		}
	}
}

func subBalance(addr common.Address, amount *big.Int) {
	for n := 0; n < 3; n++ {
		if isEqualAddress(accounts[n].address, addr) {
			accounts[n].balance.Sub(&accounts[n].balance, amount)
			break
		}
	}
	fmt.Println("Sub Balance from : ", hex.EncodeToString(addr.Bytes()))
	fmt.Println()
}

func removeAccount(addr common.Address) {
	fmt.Println("Removed : ", hex.EncodeToString(addr.Bytes()))
	//Put your remove account code here!
	fmt.Println()
}

func createAccount(addr common.Address) {
	fmt.Println("Created : ", hex.EncodeToString(addr.Bytes()))
	//Put your create account code here!
	fmt.Println()
}

func executeVM(calleraddr common.Address, calleeaddr common.Address, benefaddr common.Address,
	gasprice uint64, transgaslimit uint64, value uint64, nonce uint64,
	blocknumber uint64, timestamp uint64, difficulty uint64, headergaslimit uint64,
	data []byte, isDeploying bool) {

	transaction := sputnikvm.Transaction{
		Caller:   calleraddr,
		GasPrice: new(big.Int).SetUint64(gasprice),
		GasLimit: new(big.Int).SetUint64(transgaslimit),
		Address:  &calleeaddr,
		Value:    new(big.Int).SetUint64(value),
		Input:    data,
		Nonce:    new(big.Int).SetUint64(nonce),
	}

	header := sputnikvm.HeaderParams{
		Beneficiary: benefaddr,
		Timestamp:   timestamp,
		Number:      new(big.Int).SetUint64(blocknumber),
		Difficulty:  new(big.Int).SetUint64(difficulty),
		GasLimit:    new(big.Int).SetUint64(headergaslimit),
	}

	vm := sputnikvm.NewFrontier(&transaction, &header)

	fmt.Println("\nSputnikvm is starting...")

Loop:
	for {
		require := vm.Fire()

		switch require.Typ() {

		case sputnikvm.RequireNone:
			break Loop

		case sputnikvm.RequireAccount:

			if require.Address().IsEmpty() {
				vm.CommitNonexist(require.Address())
			} else {
				var acc *account
				acc = getAccountByAddress(require.Address())
				vm.CommitAccount(require.Address(), new(big.Int).SetUint64(0), &acc.balance, acc.code)
			}

		case sputnikvm.RequireAccountCode:
			var acc *account
			acc = getAccountByAddress(require.Address())
			vm.CommitAccountCode(require.Address(), acc.code)

		case sputnikvm.RequireAccountStorage:
			v := getAccountStorageValue(require.Address(), *require.StorageKey())
			vm.CommitAccountStorage(require.Address(), require.StorageKey(), v)

		case sputnikvm.RequireBlockhash:
			vm.CommitBlockhash(require.BlockNumber(), getBlockHash(require.BlockNumber()))

		default:
			panic("Panic : unreachable!")
		}
	}

	fmt.Printf("\nUsed Gas:%v\n", vm.UsedGas())
	fmt.Printf("\nlen log is %d\n", len(vm.Logs()))

	for _, log := range vm.Logs() {
		for i := 0; i < len(log.Topics); i++ {
			println("LOG address: ", log.Address.Str())
			ltopics := len(log.Topics)
			if ltopics > 0 {
				for nt := 0; nt < ltopics; nt++ {
					topic := log.Topics[nt]
					fmt.Println("LOG Topic ", nt, ": ", topic.Hex())
				}
			}
			println("LOG data: ", hex.EncodeToString(log.Data))
			println()
		}
	}

	changedaccs := vm.AccountChanges()
	lacc := len(changedaccs)

	for i := 0; i < lacc; i++ {
		acc1 := changedaccs[i]

		switch acc1.Typ() {

		case sputnikvm.AccountChangeIncreaseBalance:
			amount := acc1.ChangedAmount()
			addBalance(acc1.Address(), amount)

		case sputnikvm.AccountChangeDecreaseBalance:
			amount := acc1.ChangedAmount()
			subBalance(acc1.Address(), amount)

		case sputnikvm.AccountChangeRemoved:
			removeAccount(acc1.Address())

		case sputnikvm.AccountChangeFull, sputnikvm.AccountChangeCreate:
			cod := hex.EncodeToString(acc1.Code())
			setAccountCode(acc1.Address(), acc1.Code())
			if len(cod) > 0 {
				fmt.Println("\naddress is: ", hex.EncodeToString(acc1.Address().Bytes()))
				fmt.Println("code is: ", cod)
				println()
			}

			if acc1.Typ() == sputnikvm.AccountChangeFull {
				changeStorage := acc1.ChangedStorage()
				if len(changeStorage) > 0 {
					fmt.Println("Size of changed storage: ", len(changeStorage))
					for i := 0; i < len(changeStorage); i++ {
						fmt.Println("Key: ", common.BigToHash(changeStorage[i].Key).Hex(), "=", common.BigToHash(changeStorage[i].Value).Hex())
						setAccountStorageValue(acc1.Address(), *changeStorage[i].Key, *changeStorage[i].Value)
					}
					println()
				}
			} else {
				createAccount(acc1.Address())
				changeStorage := acc1.Storage()
				if len(changeStorage) > 0 {
					fmt.Println("Size of changed storage: ", len(changeStorage))
					for i := 0; i < len(changeStorage); i++ {
						fmt.Println("Key: ", common.BigToHash(changeStorage[i].Key).Hex(), "=", common.BigToHash(changeStorage[i].Value).Hex())
						addToAccountStorage(acc1.Address(), *changeStorage[i].Key, *changeStorage[i].Value)
					}
					println()
				}
			}

		default:
			panic("Panic :unreachable!")
		}

	}

	if !vm.Failed() && isDeploying {
		setAccountCode(calleeaddr, vm.Output())
	}

	fmt.Printf("VM Output : %s\n\n", hex.EncodeToString(vm.Output()))
	println("VM Successfuly : ", !vm.Failed())

	vm.Free()

}
