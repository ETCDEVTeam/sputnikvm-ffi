package sputnikvm

// #include "../../c/sputnikvm.h"
import "C"

import (
	"math/big"
	"github.com/ethereumproject/go-ethereum/common"
)

type VM struct {
	c *C.sputnikvm_vm_t
}

type Transaction struct {
	Caller common.Address
	GasPrice *big.Int
	GasLimit *big.Int
	Address *common.Address // If it is nil, then we take it as a Create transaction.
	Value *big.Int
	Input []byte
	Nonce *big.Int
}

func PrintCU256(v C.sputnikvm_u256) {
	C.print_u256(v)
}

func ToCU256(v *big.Int) C.sputnikvm_u256 {
	bytes := v.Bytes()
	cu256 := new(C.sputnikvm_u256)
	for i := 0; i < 32; i++ {
		if i < (32 - len(bytes)) {
			continue
		}
		cu256.data[i] = C.uchar(bytes[i - (32 - len(bytes))])
	}
	return *cu256
}

func ToCGas(v *big.Int) C.sputnikvm_gas {
	bytes := v.Bytes()
	cgas := new(C.sputnikvm_gas)
	for i := 0; i < 32; i++ {
		if i < (32 - len(bytes)) {
			continue
		}
		cgas.data[i] = C.uchar(bytes[i - (32 - len(bytes))])
	}
	return *cgas
}

func ToCAddress(v common.Address) *C.sputnikvm_address {
	panic("not implemented")
}

func ToCTransaction(transaction *Transaction) *C.sputnikvm_transaction {
	panic("not implemented")
}
