extern crate libc;
extern crate bigint;
extern crate sputnikvm;

mod common;

pub use common::{c_address, c_gas, c_u256};

use std::slice;
use std::rc::Rc;
use libc::{c_uchar, c_uint, c_longlong};
use sputnikvm::{TransactionAction, ValidTransaction, HeaderParams, SeqTransactionVM, Patch};

type c_action = c_uchar;
pub const CALL_ACTION: c_action = 0;
pub const CREATE_ACTION: c_action = 1;

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

#[no_mangle]
fn sputnikvm_new<P: Patch>(
    transaction: c_transaction, header: c_header_params
) -> *mut SeqTransactionVM<P> {
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
            let s = unsafe {
                slice::from_raw_parts(transaction.input, transaction.input_len as usize)
            };
            let mut r = Vec::new();
            for v in s {
                r.push(*v);
            }
            Rc::new(r)
        },
        nonce: transaction.nonce.into(),
    };

    unimplemented!()
}
