pub const SEED: u64 = 7;
pub const L: usize = 1000;
pub const N: usize = 100;
pub const CHECKSUM: u64 = 107829970005;

pub struct Quicksort {
    seed: u64,
}

impl Quicksort {
    pub fn new(seed: u64) -> Quicksort {
        Quicksort { seed: seed }
    }

    fn random(&mut self) -> usize {
        self.seed = (1103515245 * self.seed + 12345) % (1 << 31);
        (self.seed % (std::u32::MAX as u64 + 1)) as usize
    }

    fn randomize_array(&mut self, arr: &mut Vec<usize>) {
        for x in arr.iter_mut() {
            *x = self.random();
        }
    }

    fn quick_sort(&mut self, arr: &mut Vec<usize>, left: usize, right: usize) {
        if left >= right {
            return;
        }
        let i = left;
        let j = right;
        let pivot = arr[left + (right - left) / 2];

        let mut i = i;
        let mut j = j;
        while i <= j {
            while arr[i] < pivot {
                i += 1;
            }
            while pivot < arr[j] {
                j -= 1;
            }
            if i <= j {
                arr.swap(i, j);
                i += 1;
                j = j.saturating_sub(1);
            }
        }
        if left < j {
            self.quick_sort(arr, left, j);
        }
        if i < right {
            self.quick_sort(arr, i, right);
        }
    }

    pub fn benchmark(&mut self) -> u64 {
        let mut checksum: u64 = 0;
        let mut arr = vec![0; L];
        for _ in 0..N {
            self.randomize_array(&mut arr);
            self.quick_sort(&mut arr, 0, L - 1);
            checksum += arr[L / 2] as u64;
        }
        checksum
    }
}

#[cfg(all(test, not(target_arch = "wasm32")))]
mod tests {
    use super::*;
    extern crate test;
    use test::Bencher;

    #[bench]
    fn benchmark_check_checksum(b: &mut Bencher) {
        b.iter(|| {
            let mut qs = Quicksort::new(SEED);
            let checksum = qs.benchmark();
            assert_eq!(checksum, CHECKSUM);
        });
    }
}
