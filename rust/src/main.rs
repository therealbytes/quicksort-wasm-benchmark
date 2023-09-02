#![cfg_attr(not(target_arch = "wasm32"), feature(test))]

#[cfg(not(target_arch = "wasm32"))]
extern crate test;

mod sort;

#[cfg_attr(all(target_arch = "wasm32"), export_name = "run")]
#[no_mangle]
pub extern "C" fn _run(seed: u64) -> u64 {
    let mut qs = sort::Quicksort::new(seed);
    let checksum = qs.benchmark();
    checksum
}

pub fn main() {}
