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
