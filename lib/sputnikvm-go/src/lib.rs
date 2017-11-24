extern crate libc;
extern crate bigint;
extern crate sputnikvm;

mod common;

pub use common::{c_address, c_gas, c_u256};

use std::slice;
use std::ptr;
use std::rc::Rc;
use std::ops::DerefMut;
use libc::{c_uchar, c_uint, c_longlong};
use sputnikvm::{TransactionAction, ValidTransaction, HeaderParams, SeqTransactionVM, Patch,
                MainnetFrontierPatch, MainnetHomesteadPatch, MainnetEIP150Patch, MainnetEIP160Patch,
                VM, RequireError};

type c_action = c_uchar;
#[no_mangle]
pub static CALL_ACTION: c_action = 0;
#[no_mangle]
pub static CREATE_ACTION: c_action = 1;

#[repr(C)]
pub struct c_transaction {
    pub caller: c_address,
    pub gas_price: c_gas,
    pub gas_limit: c_gas,
    pub action: c_action,
    pub action_address: c_address,
    pub value: c_u256,
    pub input: *const c_uchar,
    pub input_len: c_uint,
    pub nonce: c_u256,
}

#[repr(C)]
pub struct c_header_params {
    pub beneficiary: c_address,
    pub timestamp: c_longlong,
    pub number: c_u256,
    pub difficulty: c_u256,
    pub gas_limit: c_gas,
}

#[repr(C)]
pub struct c_require {
    pub typ: c_require_type,
    pub value: c_require_value,
}

#[repr(C)]
pub enum c_require_type {
    none,
    account,
    account_code,
    account_storage,
    blockhash
}

#[repr(C)]
pub union c_require_value {
    pub account: c_address,
    pub account_storage: c_require_value_account_storage,
    pub blockhash: c_u256,
}

#[repr(C)]
#[derive(Clone, Copy)]
pub struct c_require_value_account_storage {
    pub address: c_address,
    pub key: c_u256,
}

fn sputnikvm_new<P: Patch + 'static>(
    transaction: c_transaction, header: c_header_params
) -> *mut Box<VM> {
    let transaction = ValidTransaction {
        caller: Some(transaction.caller.into()),
        gas_price: transaction.gas_price.into(),
        gas_limit: transaction.gas_limit.into(),
        action: if transaction.action == CALL_ACTION {
            TransactionAction::Call(transaction.action_address.into())
        } else if transaction.action == CREATE_ACTION {
            TransactionAction::Create
        } else {
            panic!()
        },
        value: transaction.value.into(),
        input: {
            if transaction.input.is_null() {
                Rc::new(Vec::new())
            } else {
                let s = unsafe {
                    slice::from_raw_parts(transaction.input, transaction.input_len as usize)
                };
                let mut r = Vec::new();
                for v in s {
                    r.push(*v);
                }
                Rc::new(r)
            }
        },
        nonce: transaction.nonce.into(),
    };

    let header = HeaderParams {
        beneficiary: header.beneficiary.into(),
        timestamp: header.timestamp as u64,
        number: header.number.into(),
        difficulty: header.difficulty.into(),
        gas_limit: header.gas_limit.into(),
    };

    let vm = SeqTransactionVM::<P>::new(transaction, header);
    Box::into_raw(Box::new(Box::new(vm)))
}

#[no_mangle]
pub extern "C" fn sputnikvm_new_frontier(
    transaction: c_transaction, header: c_header_params
) -> *mut Box<VM> {
    sputnikvm_new::<MainnetFrontierPatch>(transaction, header)
}

#[no_mangle]
pub extern "C" fn sputnikvm_new_homestead(
    transaction: c_transaction, header: c_header_params
) -> *mut Box<VM> {
    sputnikvm_new::<MainnetHomesteadPatch>(transaction, header)
}

#[no_mangle]
pub extern "C" fn sputnikvm_new_eip150(
    transaction: c_transaction, header: c_header_params
) -> *mut Box<VM> {
    sputnikvm_new::<MainnetEIP150Patch>(transaction, header)
}

#[no_mangle]
pub extern "C" fn sputnikvm_new_eip160(
    transaction: c_transaction, header: c_header_params
) -> *mut Box<VM> {
    sputnikvm_new::<MainnetEIP160Patch>(transaction, header)
}

#[no_mangle]
pub extern "C" fn sputnikvm_free(
    vm: *mut Box<VM>
) {
    if vm.is_null() { return; }
    unsafe { Box::from_raw(vm); }
}

#[no_mangle]
pub extern "C" fn sputnikvm_fire(
    vm: *mut Box<VM>
) -> c_require {
    let mut vm_box = unsafe { Box::from_raw(vm) };
    let ret;
    {
        let vm: &mut VM = vm_box.deref_mut().deref_mut();
        match vm.fire() {
            Ok(()) => {
                ret = c_require {
                    typ: c_require_type::none,
                    value: c_require_value {
                        account: c_address::default(),
                    }
                };
            },
            Err(RequireError::Account(address)) => {
                ret = c_require {
                    typ: c_require_type::account,
                    value: c_require_value {
                        account: address.into(),
                    }
                };
            },
            Err(RequireError::AccountCode(address)) => {
                ret = c_require {
                    typ: c_require_type::account_code,
                    value: c_require_value {
                        account: address.into(),
                    }
                };
            },
            Err(RequireError::AccountStorage(address, key)) => {
                ret = c_require {
                    typ: c_require_type::account_storage,
                    value: c_require_value {
                        account_storage: c_require_value_account_storage {
                            address: address.into(),
                            key: key.into(),
                        },
                    }
                };
            },
            Err(RequireError::Blockhash(number)) => {
                ret = c_require {
                    typ: c_require_type::blockhash,
                    value: c_require_value {
                        blockhash: number.into(),
                    },
                };
            },
        }
    }
    Box::into_raw(vm_box);
    ret
}

#[no_mangle]
pub extern "C" fn sputnikvm_default_transaction() -> c_transaction {
    c_transaction {
        caller: c_address::default(),
        gas_price: c_gas::default(),
        gas_limit: c_gas::default(),
        action: CALL_ACTION,
        action_address: c_address::default(),
        value: c_u256::default(),
        input: ptr::null(),
        input_len: 0,
        nonce: c_u256::default(),
    }
}

#[no_mangle]
pub extern "C" fn sputnikvm_default_header_params() -> c_header_params {
    c_header_params {
        beneficiary: c_address::default(),
        timestamp: 0,
        number: c_u256::default(),
        difficulty: c_u256::default(),
        gas_limit: c_gas::default(),
    }
}
