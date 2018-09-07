package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	//"github.com/ETCDEVTeam/sputnikvm-ffi/go/sputnikvm"
	"github.com/GheisMohammadi/sputnikvm-ffi/go/sputnikvm"
	"github.com/ethereumproject/go-ethereum/common"
)

func main() {

	var addrcaller common.Address
	addrcaller.SetBytes([]byte("0x50ad1920e22425b10f32d9b9dba566af48ba8d48"))

	var addrCallee common.Address
	//addrCallee.SetBytes([]byte("0x6439623964626135363661663438626138643438"))

	var addrbenef common.Address
	addrbenef.SetBytes([]byte("0x4a9274e0d453b4bbb6694065e561d546da48bbbd"))

	fmt.Println("caller: ", addrcaller.Str())
	fmt.Println("callee ", addrCallee.Str())
	fmt.Println("benef ", addrbenef.Str())

	testCode, _ := hex.DecodeString("60806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146064575b600080fd5b348015605857600080fd5b5060626004356088565b005b348015606f57600080fd5b506076608d565b60408051918252519081900360200190f35b600055565b6212d687905600a165627a7a723058209baea01664908f42f35914eeeee9092206e36832612009367f5e9c730337a1990029")

	//setMethod, _ := hex.DecodeString("60fe47b100000000000000000000000000000000000000000000000000000000000001c8")
	getMethod, _ := hex.DecodeString("6d4ce63c")
	println("\nMethod: ", getMethod)

	transaction := sputnikvm.Transaction{
		Caller:   addrcaller,
		GasPrice: new(big.Int).SetUint64(0),
		GasLimit: new(big.Int).SetUint64(10000000),
		Address:  &addrCallee,               //new(common.Address),
		Value:    new(big.Int).SetUint64(0), // The Amount that should transfer to the smart contract
		// For this reason the smart contract must have atleast one payable method!
		Input: getMethod, //[]byte{},
		Nonce: new(big.Int).SetUint64(0),
	}

	header := sputnikvm.HeaderParams{
		Beneficiary: addrbenef, //*new(common.Address),
		Timestamp:   0,
		Number:      new(big.Int).SetUint64(1),
		Difficulty:  new(big.Int).SetUint64(0),
		GasLimit:    new(big.Int).SetUint64(1000000000),
	}

	vm := sputnikvm.NewMordenEIP160(&transaction, &header)
	storageValue := new(big.Int).SetUint64(0)
Loop:
	for {
		require := vm.Fire()

		switch require.Typ() {

		case sputnikvm.RequireNone:
			fmt.Println("sputnikvm.RequireNone")
			fmt.Println("-> exit ok from loop!")
			break Loop

		case sputnikvm.RequireAccount:
			fmt.Printf("sputnikvm.RequireAccount  addr:%s \n", require.Address().Str())
			if require.Address().Str() == addrcaller.Str() {
				fmt.Println("-> commited some data!")
				vm.CommitAccount(require.Address(), new(big.Int).SetUint64(2), new(big.Int).SetUint64(1500000), []byte(""))
			} else if require.Address().Str() == addrCallee.Str() { //.IsEmpty() {
				fmt.Println("-> commited code!")
				vm.CommitAccount(require.Address(), new(big.Int).SetUint64(10), new(big.Int).SetUint64(100), testCode)
			} else {
				fmt.Println("-> commited non exist!")
				vm.CommitNonexist(require.Address())
			}

		case sputnikvm.RequireAccountCode:
			fmt.Printf("sputnikvm.RequireAccountCode  addr:%s\n", require.Address().Str())
			vm.CommitAccountCode(require.Address(), testCode)

		case sputnikvm.RequireAccountStorage:
			fmt.Printf("sputnikvm.RequireAccountStorage  addr:%s\n", require.Address().Str())
			vm.CommitAccountStorage(require.Address(), require.StorageKey(), storageValue)

		case sputnikvm.RequireBlockhash:
			fmt.Printf("sputnikvm.RequireBlockhash  addr:%s\n", require.Address().Str())
			vm.CommitBlockhash(require.BlockNumber(), *new(common.Hash))

		default:
			fmt.Printf("Require Default  addr:%s\n", require.Address().Str())
			panic("Panic : unreachable!")
		}
	}

	//must return 12D687

	fmt.Printf("\nUsed Gas:%v\n", vm.UsedGas())
	fmt.Printf("\nLogs: %v\n", vm.Logs())
	fmt.Printf("\nAccount Changes: %v\n", vm.AccountChanges())
	fmt.Printf("\nlen log is %d\n", len(vm.Logs()))

	for _, log := range vm.Logs() {
		for i := 0; i < len(log.Topics); i++ {
			println("LOG address: ", log.Address.Str())
			println("LOG data: ", log.Data)
		}
	}

	accounts := vm.AccountChanges()
	lacc := len(accounts)
	fmt.Println("\nlen accounts is", lacc)

	for i := 0; i < lacc; i++ {
		account := accounts[i]
		fmt.Println("\naddress is ", account.Address().Str())

		switch account.Typ() {

		case sputnikvm.AccountChangeIncreaseBalance:
			fmt.Println("AccountChangeIncreaseBalance")
			amount := account.ChangedAmount()
			fmt.Println("IncreaseBalance: ", amount)

		case sputnikvm.AccountChangeDecreaseBalance:
			fmt.Println("AccountChangeDecreaseBalance")
			amount := account.ChangedAmount()
			fmt.Println("DecreaseBalance: ", amount)

		case sputnikvm.AccountChangeRemoved:
			fmt.Println("AccountChangeRemoved")

		case sputnikvm.AccountChangeFull, sputnikvm.AccountChangeCreate:
			cod := hex.EncodeToString(account.Code())
			fmt.Println("code is ", cod)

			if account.Typ() == sputnikvm.AccountChangeFull {
				fmt.Println("AccountChangeFull")
				changeStorage := account.ChangedStorage()
				fmt.Println("Size of changed storage: ", len(changeStorage))
				for i := 0; i < len(changeStorage); i++ {
					fmt.Println("Key ", common.BigToHash(changeStorage[i].Key).Hex(), "=", common.BigToHash(changeStorage[i].Value).Hex())
				}
			} else {
				fmt.Println("AccountChangeCreate")
				changeStorage := account.Storage()
				fmt.Println("Size of changed storage: ", len(changeStorage))
				for i := 0; i < len(changeStorage); i++ {
					fmt.Println("Key ", common.BigToHash(changeStorage[i].Key).Hex(), "=", common.BigToHash(changeStorage[i].Value).Hex())
				}
			}

		default:
			fmt.Println("default")
			panic("Panic :unreachable!")
		}

	}
	println("VM Output: ", vm.OutLen())
	var i uint
	for i = 0; i < 32; i++ {
		println("VM Output: ", vm.Out(i))
	}

	fmt.Printf("\nVM Output Array: %v\n", vm.Output())

	println("VM Successfuly : ", !vm.Failed())
	vm.Free()

}
