package main

// #include "../c/sputnikvm.h"
import "C"

import (
	"fmt"
	"math/big"
	"github.com/ethereumproject/sputnikvm-ffi/go/sputnikvm"
)

func main() {
	goint := big.NewInt(100000000000000)
	cint := sputnikvm.ToCU256(goint)
	fmt.Printf("%v", cint)
	sputnikvm.PrintCU256(cint)

	transaction := C.sputnikvm_default_transaction()
	header := C.sputnikvm_default_header_params()
	vm := C.sputnikvm_new_frontier(transaction, header)
	ret := C.sputnikvm_fire(vm)
	fmt.Printf("%v", ret)
	C.sputnikvm_free(vm)
}
