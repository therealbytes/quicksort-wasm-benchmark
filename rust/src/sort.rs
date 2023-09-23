pub struct QuicksortBenchmark {
    seed: usize,
}

impl QuicksortBenchmark {
    pub fn new(seed: usize) -> QuicksortBenchmark {
        QuicksortBenchmark { seed: seed }
    }

    fn random(&mut self) -> u32 {
        self.seed = (1103515245 * self.seed + 12345) % (1 << 31);
        self.seed as u32
    }

    fn randomize_array(&mut self, arr: &mut Vec<u32>) {
        for x in arr.iter_mut() {
            *x = self.random() % 1000;
        }
    }

    fn quick_sort(&mut self, arr: &mut Vec<u32>, left: usize, right: usize) {
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

    pub fn run(&mut self, arr_len: usize, iter: usize) -> usize {
        let mut checksum: u32 = 0;
        let mut arr = vec![0; arr_len];
        for _ in 0..iter {
            self.randomize_array(&mut arr);
            self.quick_sort(&mut arr, 0, arr_len - 1);
            checksum += arr[arr_len/2];
        }
        checksum as usize
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
            let mut qs = QuicksortBenchmark::new(7);
            let checksum = qs.run(1000, 100);
            assert_eq!(checksum, 49760);
        });
    }
}
