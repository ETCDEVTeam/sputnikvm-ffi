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

func ToCU256(v *big.Int) *C.sputnikvm_u256 {
	bytes := v.Bytes()
	cu256 := new(C.sputnikvm_u256)
	for i, b := range bytes {
		cu256.data[i] = C.uchar(b)
	}
	return cu256
}

func ToCGas(v *big.Int) *C.sputnikvm_gas {
	panic("not implemented")
}

func ToCAddress(v common.Address) *C.sputnikvm_address {
	panic("not implemented")
}

func ToCTransaction(transaction *Transaction) *C.sputnikvm_transaction {
	panic("not implemented")
}
