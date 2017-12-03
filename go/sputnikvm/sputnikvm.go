package sputnikvm

// #include "../../c/sputnikvm.h"
// #include <stdlib.h>
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

type HeaderParams struct {
	Beneficiary common.Address
	Timestamp uint64
	Number *big.Int
	Difficulty *big.Int
	GasLimit *big.Int
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

func ToCAddress(v common.Address) C.sputnikvm_address {
	caddress := new(C.sputnikvm_address)
	for i := 0; i < 20; i++ {
		caddress.data[i] = C.uchar(v[i])
	}
	return *caddress
}

func toCTransaction(transaction *Transaction) *C.sputnikvm_transaction {
	// Malloc input length memory and must be freed manually.

	ctransaction := new(C.sputnikvm_transaction)
	cinput := C.malloc(C.size_t(len(transaction.Input)))
	for i := 0; i < len(transaction.Input); i++ {
		(*(*[]C.uchar)(cinput))[i] = C.uchar(transaction.Input[i])
	}
	ctransaction.caller = ToCAddress(transaction.Caller)
	ctransaction.gas_price = ToCGas(transaction.GasPrice)
	ctransaction.gas_limit = ToCGas(transaction.GasLimit)
	if transaction.Address == nil {
		ctransaction.action = C.sputnikvm_action(C.CREATE_ACTION)
	} else {
		ctransaction.action = C.sputnikvm_action(C.CALL_ACTION)
		ctransaction.action_address = ToCAddress(*transaction.Address)
	}
	ctransaction.value = ToCU256(transaction.Value)
	ctransaction.input = (*C.uchar)(cinput)
	ctransaction.input_len = C.uint(len(transaction.Input))
	ctransaction.nonce = ToCU256(transaction.Nonce)

	return ctransaction
}

func ToCHeaderParams(header *HeaderParams) *C.sputnikvm_header_params {
	cheader := new(C.sputnikvm_header_params)
	cheader.beneficiary = ToCAddress(header.Beneficiary)
	cheader.timestamp = C.ulonglong(header.Timestamp)
	cheader.number = ToCU256(header.Number)
	cheader.difficulty = ToCU256(header.Difficulty)
	cheader.gas_limit = ToCGas(header.GasLimit)

	return cheader
}
