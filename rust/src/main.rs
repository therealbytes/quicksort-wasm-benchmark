#![cfg_attr(not(target_arch = "wasm32"), feature(test))]

#[cfg(not(target_arch = "wasm32"))]
extern crate test;

mod sort;

#[cfg_attr(all(target_arch = "wasm32"), export_name = "run")]
#[no_mangle]
pub extern "C" fn _run(seed: i64, arr_len: i64, iter: i64) -> i64 {
    let mut qs = sort::QuicksortBenchmark::new(seed as usize);
    let checksum = qs.run(arr_len as usize, iter as usize);
    checksum as i64
}

pub fn main() {}
