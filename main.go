package main

/*
#cgo LDFLAGS: -L./lib -lsputnikvm_go
#include "./lib/sputnikvm-go.h"
*/
import "C"

func main() {
	transaction := C.sputnikvm_default_transaction()
	header := C.sputnikvm_default_header_params()
	vm := C.sputnikvm_new_frontier(transaction, header)
	C.sputnikvm_free(vm)
}
