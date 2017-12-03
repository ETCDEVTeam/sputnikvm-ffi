package sputnikvm

// #include "../../c/sputnikvm.h"
// #include <stdlib.h>
//
// sputnikvm_address sputnikvm_require_value_read_account(sputnikvm_require_value v) {
//   return v.account;
// }
//
// sputnikvm_require_value_account_storage sputnikvm_require_value_read_account_storage(sputnikvm_require_value v) {
//   return v.account_storage;
// }
//
// sputnikvm_u256 sputnikvm_require_value_read_blockhash(sputnikvm_require_value v) {
//   return v.blockhash;
// }
import "C"

import (
	"unsafe"
	"math/big"
	"github.com/ethereumproject/go-ethereum/common"
)

type RequireType int
const (
	RequireNone = iota
	RequireAccount
	RequireAccountCode
	RequireAccountStorage
	RequireBlockhash
)

type Require struct {
	c C.sputnikvm_require
}

func (require *Require) Typ() RequireType {
	switch require.c.typ {
	case C.require_none:
		return RequireNone
	case C.require_account:
		return RequireAccount
	case C.require_account_code:
		return RequireAccountCode
	case C.require_account_storage:
		return RequireAccountStorage
	case C.require_blockhash:
		return RequireBlockhash
	default:
		panic("unreachable")
	}
}

func (require *Require) Address() common.Address {
	switch require.Typ() {
	case RequireAccount, RequireAccountCode:
		return FromCAddress(C.sputnikvm_require_value_read_account(require.c.value))
	case RequireAccountStorage:
		return FromCAddress(C.sputnikvm_require_value_read_account_storage(require.c.value).address)
	default:
		panic("incorrect usage")
	}
}

func (require *Require) StorageKey() *big.Int {
	switch require.Typ() {
	case RequireAccountStorage:
		storage := C.sputnikvm_require_value_read_account_storage(require.c.value)
		return FromCU256(storage.key)
	default:
		panic("incorrect usage")
	}
}

func (require *Require) BlockNumber() *big.Int {
	switch require.Typ() {
	case RequireBlockhash:
		number := C.sputnikvm_require_value_read_blockhash(require.c.value)
		return FromCU256(number)
	default:
		panic("incorrect usage")
	}
}

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

func FromCU256(v C.sputnikvm_u256) *big.Int {
	bytes := new([32]byte)
	for i := 0; i < 32; i++ {
		bytes[i] = byte(v.data[i])
	}
	i := new(big.Int)
	i.SetBytes(bytes[0:32])
	return i
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

func FromCGas(v C.sputnikvm_gas) *big.Int {
	bytes := new([32]byte)
	for i := 0; i < 32; i++ {
		bytes[i] = byte(v.data[i])
	}
	i := new(big.Int)
	i.SetBytes(bytes[0:32])
	return i
}

func ToCAddress(v common.Address) C.sputnikvm_address {
	caddress := new(C.sputnikvm_address)
	for i := 0; i < 20; i++ {
		caddress.data[i] = C.uchar(v[i])
	}
	return *caddress
}

func FromCAddress(v C.sputnikvm_address) common.Address {
	address := new(common.Address)
	for i := 0; i < 20; i++ {
		address[i] = byte(v.data[i])
	}
	return *address
}

func ToCH256(v common.Hash) C.sputnikvm_h256 {
	chash := new(C.sputnikvm_h256)
	for i := 0; i < 32; i++ {
		chash.data[i] = C.uchar(v[i])
	}
	return *chash
}

func FromCH256(v C.sputnikvm_h256) common.Hash {
	hash := new(common.Hash)
	for i := 0; i < 32; i++ {
		hash[i] = byte(v.data[i])
	}
	return *hash
}

func toCTransaction(transaction *Transaction) (*C.sputnikvm_transaction, unsafe.Pointer) {
	// Malloc input length memory and must be freed manually.

	ctransaction := new(C.sputnikvm_transaction)
	cinput := C.malloc(C.size_t(len(transaction.Input)))
	for i := 0; i < len(transaction.Input); i++ {
		i_cinput := unsafe.Pointer(uintptr(cinput) + uintptr(i))
		*(*C.uchar)(i_cinput) = C.uchar(transaction.Input[i])
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

	return ctransaction, cinput
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

func NewFrontier(transaction *Transaction, header *HeaderParams) *VM {
	ctransaction, cinput := toCTransaction(transaction)
	cheader := ToCHeaderParams(header)

	cvm := C.sputnikvm_new_frontier(*ctransaction, *cheader)
	C.free(cinput)

	vm := new(VM)
	vm.c = cvm

	return vm
}

func NewHomestead(transaction *Transaction, header *HeaderParams) *VM {
	ctransaction, cinput := toCTransaction(transaction)
	cheader := ToCHeaderParams(header)

	cvm := C.sputnikvm_new_homestead(*ctransaction, *cheader)
	C.free(cinput)

	vm := new(VM)
	vm.c = cvm

	return vm
}

func NewEIP150(transaction *Transaction, header *HeaderParams) *VM {
	ctransaction, cinput := toCTransaction(transaction)
	cheader := ToCHeaderParams(header)

	cvm := C.sputnikvm_new_eip150(*ctransaction, *cheader)
	C.free(cinput)

	vm := new(VM)
	vm.c = cvm

	return vm
}

func NewEIP160(transaction *Transaction, header *HeaderParams) *VM {
	ctransaction, cinput := toCTransaction(transaction)
	cheader := ToCHeaderParams(header)

	cvm := C.sputnikvm_new_eip160(*ctransaction, *cheader)
	C.free(cinput)

	vm := new(VM)
	vm.c = cvm

	return vm
}

func (vm *VM) Fire() Require {
	return Require {
		c: C.sputnikvm_fire(vm.c),
	}
}

func (vm *VM) Free() {
	C.sputnikvm_free(vm.c)
}

func (vm *VM) CommitAccount(address common.Address, nonce *big.Int, balance *big.Int, code []byte) {
	caddress := ToCAddress(address)
	cnonce := ToCU256(nonce)
	cbalance := ToCU256(balance)
	ccode := C.malloc(C.size_t(len(code)))
	for i := 0; i < len(code); i++ {
		i_ccode := unsafe.Pointer(uintptr(ccode) + uintptr(i))
		*(*C.uchar)(i_ccode) = C.uchar(code[i])
	}

	C.sputnikvm_commit_account(vm.c, caddress, cnonce, cbalance, (*C.uchar)(ccode), C.uint(len(code)))
	C.free(ccode)
}

func (vm *VM) CommitAccountCode(address common.Address, code []byte) {
	caddress := ToCAddress(address)
	ccode := C.malloc(C.size_t(len(code)))
	for i := 0; i < len(code); i++ {
		i_ccode := unsafe.Pointer(uintptr(ccode) + uintptr(i))
		*(*C.uchar)(i_ccode) = C.uchar(code[i])
	}

	C.sputnikvm_commit_account_code(vm.c, caddress, (*C.uchar)(ccode), C.uint(len(code)))
	C.free(ccode)
}

func (vm *VM) CommitAccountStorage(address common.Address, key *big.Int, value *big.Int) {
	caddress := ToCAddress(address)
	ckey := ToCU256(key)
	cvalue := ToCU256(value)

	C.sputnikvm_commit_account_storage(vm.c, caddress, ckey, cvalue)
}

func (vm *VM) CommitNonexist(address common.Address) {
	caddress := ToCAddress(address)
	C.sputnikvm_commit_nonexist(vm.c, caddress)
}

func (vm *VM) CommitBlockhash(number *big.Int, hash common.Hash) {
	cnumber := ToCU256(number)
	chash := ToCH256(hash)
	C.sputnikvm_commit_blockhash(vm.c, cnumber, chash)
}

func (vm *VM) UsedGas() *big.Int {
	cgas := C.sputnikvm_used_gas(vm.c)
	return FromCGas(cgas)
}
