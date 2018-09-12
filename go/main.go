package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	//"github.com/ETCDEVTeam/sputnikvm-ffi/go/sputnikvm"
	"github.com/GheisMohammadi/sputnikvm-ffi/go/sputnikvm"
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

func main() {

	deployCode, _ := hex.DecodeString("60806040526212d6876001556216e36060005560e9806100206000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146058575b600080fd5b6056600435607c565b005b348015606357600080fd5b50606a60b7565b60408051918252519081900360200190f35b60018190556040805182815290517f23f9887eb044d32dba99d7b0b753c61c3c3b72d70ff0addb9a843542fd7642129181900360200190a150565b600154905600a165627a7a7230582013452558cc58a514b8056c0b45a3f1ab8c5f736b2e087c65e615650b562415ff0029")

	//Setup Some Test Accounts
	fmt.Println("\n=================\nSetup Accounts...\n=================")
	setupAccounts()

	//deploy
	fmt.Println("\n=================\nDeploying...\n=================")
	setAccountCode(addrCallee, deployCode)
	executeVM(addrcaller, addrCallee, addrbenef, 0, 1000000, 0, 0, 0, 0, 0, 1000000000, deployCode, true)

	//Call Set Method
	fmt.Println("\n=================\nCall Set Method by 1234567...\n=================")
	setMethod, _ := hex.DecodeString("60fe47b100000000000000000000000000000000000000000000000000000000000001c8")
	executeVM(addrcaller, addrCallee, addrbenef, 0, 1000000, 10, 1, 1, 0, 0, 1000000, setMethod, false)

	//Call Get Method
	fmt.Println("\n=================\nCall Get Method...\n=================")
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
	fmt.Println("Get Account by Address : ", hex.EncodeToString(addr.Bytes()))

	for n := 0; n < 3; n++ {
		if isEqualAddress(accounts[n].address, addr) {
			fmt.Println("Account is Found!")
			return &accounts[n]
		}
	}
	return ret
}

func isExist(addr common.Address) bool {

	fmt.Println("Check Account By Address : ", hex.EncodeToString(addr.Bytes()))

	for n := 0; n < 3; n++ {
		if isEqualAddress(accounts[n].address, addr) {
			fmt.Println("account is found!")
			return true
		}
	}

	return false
}

func setAccountCode(addr common.Address, code []byte) {

	fmt.Println("Set Account Code : ", hex.EncodeToString(addr.Bytes()))

	for n := 0; n < 3; n++ {
		if isEqualAddress(accounts[n].address, addr) {
			kod := hex.EncodeToString(code)
			accounts[n].code, _ = hex.DecodeString(kod)
			fmt.Println("Code is set successfully!")
			break
		}
	}
}

func getAccountStorageValue(addr common.Address, key big.Int) *big.Int {

	ret := new(big.Int).SetInt64(0)
	fmt.Println("Get Account Storage : ", hex.EncodeToString(addr.Bytes()))

	for n := 0; n < 3; n++ {
		if isEqualAddress(accounts[n].address, addr) {
			//fmt.Println("Account is found for get storage!")
			nstorage := len(accounts[n].storage)
			for ns := 0; ns < nstorage; ns++ {
				if hex.EncodeToString(accounts[n].storage[ns].key.Bytes()) == hex.EncodeToString(key.Bytes()) {
					fmt.Println("key is found successfuly!")
					return &accounts[n].storage[ns].value
				}
			}
		}
	}
	return ret
}

func setAccountStorageValue(addr common.Address, key big.Int, value big.Int) {
	fmt.Println("Set Account Storage Func: ", hex.EncodeToString(addr.Bytes()))

	for n := 0; n < 3; n++ {
		if isEqualAddress(accounts[n].address, addr) {
			//fmt.Println("Account is found for set storage!")
			nstorage := len(accounts[n].storage)
			for ns := 0; ns < nstorage; ns++ {
				if hex.EncodeToString(accounts[n].storage[ns].key.Bytes()) == hex.EncodeToString(key.Bytes()) {
					accounts[n].storage[ns].value = value
					fmt.Println("key is set successfuly!")
					return
				}
			}
			fmt.Println("key is creating...")
			addToAccountStorage(addr, key, value)
		}
	}
}

func addToAccountStorage(addr common.Address, key big.Int, value big.Int) {
	fmt.Println("Add Account Storage : ", hex.EncodeToString(addr.Bytes()))

	for n := 0; n < 3; n++ {
		if isEqualAddress(accounts[n].address, addr) {
			accounts[n].storage = append(accounts[n].storage, storageItem{key, value})
			fmt.Println("key is added to storage!")
			break
		}
	}
}

func getBlockHash(blocknumber *big.Int) common.Hash {
	fmt.Println("Get Block Hash : ", blocknumber.Uint64())
	return *new(common.Hash)
}

func addBalance(addr common.Address, amount *big.Int) {
	for n := 0; n < 3; n++ {
		if isEqualAddress(accounts[n].address, addr) {
			accounts[n].balance.Add(&accounts[n].balance, amount)
			break
		}
	}
	fmt.Println("Add Balance : ", hex.EncodeToString(addr.Bytes()))
}

func subBalance(addr common.Address, amount *big.Int) {
	for n := 0; n < 3; n++ {
		if isEqualAddress(accounts[n].address, addr) {
			accounts[n].balance.Sub(&accounts[n].balance, amount)
			break
		}
	}
	fmt.Println("Sub Balance : ", hex.EncodeToString(addr.Bytes()))
}

func removeAccount(addr common.Address) {
	fmt.Println("Removed : ", hex.EncodeToString(addr.Bytes()))
}
func createAccount(addr common.Address) {
	fmt.Println("Created : ", hex.EncodeToString(addr.Bytes()))
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

	fmt.Println("\nsputnikvm VM starting...")

Loop:
	for {
		require := vm.Fire()

		switch require.Typ() {

		case sputnikvm.RequireNone:
			fmt.Println("<- sputnikvm.RequireNone")
			fmt.Println("-> exit ok from loop!")
			break Loop

		case sputnikvm.RequireAccount:

			fmt.Printf("<- sputnikvm.RequireAccount  addr:%s \n", hex.EncodeToString(require.Address().Bytes())) // require.Address().Str())

			if require.Address().IsEmpty() {
				vm.CommitNonexist(require.Address())
				fmt.Println("empty account -> commited non exist!")
			} else {
				var acc *account
				acc = getAccountByAddress(require.Address())
				vm.CommitAccount(require.Address(), new(big.Int).SetUint64(0), &acc.balance, acc.code)
				fmt.Println("-> commited account data!")
			}

		case sputnikvm.RequireAccountCode:
			fmt.Printf("<- sputnikvm.RequireAccountCode  addr:%s\n", hex.EncodeToString(require.Address().Bytes()))
			var acc *account
			acc = getAccountByAddress(require.Address())
			vm.CommitAccountCode(require.Address(), acc.code)

		case sputnikvm.RequireAccountStorage:
			fmt.Printf("<- sputnikvm.RequireAccountStorage  addr:%s\n", hex.EncodeToString(require.Address().Bytes()))
			v := getAccountStorageValue(require.Address(), *require.StorageKey())
			vm.CommitAccountStorage(require.Address(), require.StorageKey(), v)

		case sputnikvm.RequireBlockhash:
			fmt.Printf("<- sputnikvm.RequireBlockhash  addr:%s\n", hex.EncodeToString(require.Address().Bytes()))
			vm.CommitBlockhash(require.BlockNumber(), getBlockHash(require.BlockNumber()))

		default:
			fmt.Printf("<- Require Default  addr:%s\n", hex.EncodeToString(require.Address().Bytes()))
			panic("Panic : unreachable!")
		}
	}

	fmt.Printf("\nUsed Gas:%v\n", vm.UsedGas())
	fmt.Printf("\nLogs: %v\n", vm.Logs())
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
		}
	}

	changedaccs := vm.AccountChanges()
	lacc := len(changedaccs)
	fmt.Println("\nlen accounts is", lacc)

	for i := 0; i < lacc; i++ {
		acc1 := changedaccs[i]
		fmt.Println("\naddress is ", hex.EncodeToString(acc1.Address().Bytes()))

		switch acc1.Typ() {

		case sputnikvm.AccountChangeIncreaseBalance:
			fmt.Println("AccountChangeIncreaseBalance")
			amount := acc1.ChangedAmount()
			addBalance(acc1.Address(), amount)
			fmt.Println("IncreaseBalance: ", amount)

		case sputnikvm.AccountChangeDecreaseBalance:
			fmt.Println("AccountChangeDecreaseBalance")
			amount := acc1.ChangedAmount()
			subBalance(acc1.Address(), amount)
			fmt.Println("DecreaseBalance: ", amount)

		case sputnikvm.AccountChangeRemoved:
			removeAccount(acc1.Address())
			fmt.Println("AccountChangeRemoved")

		case sputnikvm.AccountChangeFull, sputnikvm.AccountChangeCreate:
			cod := hex.EncodeToString(acc1.Code())
			setAccountCode(acc1.Address(), acc1.Code())
			fmt.Println("code is ", cod)

			if acc1.Typ() == sputnikvm.AccountChangeFull {
				fmt.Println("AccountChangeFull")
				changeStorage := acc1.ChangedStorage()
				fmt.Println("Size of changed storage: ", len(changeStorage))
				for i := 0; i < len(changeStorage); i++ {
					fmt.Println("Key ", common.BigToHash(changeStorage[i].Key).Hex(), "=", common.BigToHash(changeStorage[i].Value).Hex())
					setAccountStorageValue(acc1.Address(), *changeStorage[i].Key, *changeStorage[i].Value)
				}
			} else {
				fmt.Println("AccountChangeCreate")
				changeStorage := acc1.Storage()
				fmt.Println("Size of changed storage: ", len(changeStorage))
				for i := 0; i < len(changeStorage); i++ {
					fmt.Println("Key ", common.BigToHash(changeStorage[i].Key).Hex(), "=", common.BigToHash(changeStorage[i].Value).Hex())
					addToAccountStorage(acc1.Address(), *changeStorage[i].Key, *changeStorage[i].Value)
				}
			}

		default:
			fmt.Println("default")
			panic("Panic :unreachable!")
		}

	}

	if !vm.Failed() && isDeploying {
		setAccountCode(calleeaddr, vm.Output())
	}

	println("\nVM Output Len: ", vm.OutLen())
	fmt.Printf("VM Output Array: %s\n", hex.EncodeToString(vm.Output()))
	println("VM Successfuly : ", !vm.Failed())

	vm.Free()

}
