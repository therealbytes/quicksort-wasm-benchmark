#![cfg_attr(not(target_arch = "wasm32"), feature(test))]

#[cfg(not(target_arch = "wasm32"))]
extern crate test;

use crate::sort::SEED;

mod sort;

#[cfg_attr(all(target_arch = "wasm32"), export_name = "run")]
#[no_mangle]
pub extern "C" fn _run() -> u64 {
    let mut qs = sort::Quicksort::new(SEED);
    let checksum = qs.benchmark();
    checksum
}

pub fn main() {}
