package main

/*
#cgo LDFLAGS: ../c/libsputnikvm.a -ldl
#include "../c/sputnikvm.h"
*/
import "C"

import (
	"fmt"
	"math/big"
	"github.com/ethereumproject/sputnikvm-ffi/go/sputnikvm"
)

func main() {
	goint := big.NewInt(5)
	cint := sputnikvm.ToCU256(goint)
	fmt.Printf("%v", cint)
	transaction := C.sputnikvm_default_transaction()
	header := C.sputnikvm_default_header_params()
	vm := C.sputnikvm_new_frontier(transaction, header)
	ret := C.sputnikvm_fire(vm)
	fmt.Printf("%v", ret)
	C.sputnikvm_free(vm)
}
