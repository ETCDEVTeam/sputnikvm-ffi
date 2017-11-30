package main

/*
#cgo LDFLAGS: ../c/libsputnikvm.a -ldl
#include "../c/sputnikvm.h"
*/
import "C"

import (
	"fmt"
	// "github.com/ethereumproject/sputnikvm-ffi/go/sputnikvm"
)

func main() {

	transaction := C.sputnikvm_default_transaction()
	header := C.sputnikvm_default_header_params()
	vm := C.sputnikvm_new_frontier(transaction, header)
	ret := C.sputnikvm_fire(vm)
	fmt.Printf("%v", ret)
	C.sputnikvm_free(vm)
}
