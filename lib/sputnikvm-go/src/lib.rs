extern crate libc;

use libc::{c_uchar, c_uint, c_longlong};

pub struct c_address {
    pub data: [c_uchar; 20],
}
pub struct c_gas {
    pub data: [c_uchar; 32],
}
pub struct c_u256 {
    pub data: [c_uchar; 32],
}

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
