package main

import (
	"fmt"
	"math/big"
	"github.com/ethereumproject/go-ethereum/common"
	"github.com/ethereumproject/sputnikvm-ffi/go/sputnikvm"
)

func main() {
	goint := big.NewInt(100000000000000)
	cint := sputnikvm.ToCU256(goint)
	fmt.Printf("%v", cint)
	sputnikvm.PrintCU256(cint)

	transaction := sputnikvm.Transaction {
		Caller: *new(common.Address),
		GasPrice: new(big.Int),
		GasLimit: new(big.Int),
		Address: new(common.Address),
		Value: new(big.Int),
		Input: []byte{1, 2, 3, 4, 5},
		Nonce: new(big.Int),
	}

	header := sputnikvm.HeaderParams {
		Beneficiary: *new(common.Address),
		Timestamp: 0,
		Number: new(big.Int),
		Difficulty: new(big.Int),
		GasLimit: new(big.Int),
	}

	vm := sputnikvm.NewFrontier(&transaction, &header)
	ret := vm.Fire()
	fmt.Printf("%v", ret)
	vm.Free()
}
