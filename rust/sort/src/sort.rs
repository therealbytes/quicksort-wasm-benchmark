pub struct QuickSort {
    seed: u64,
}

impl QuickSort {
    pub fn new(seed: u64) -> QuickSort {
        QuickSort { seed: seed }
    }

    fn random(&mut self) -> u64 {
        self.seed = (1103515245 * self.seed + 12345) % (1 << 31);
        self.seed
    }

    fn randomize_array(&mut self, arr: &mut Vec<u64>) {
        for x in arr.iter_mut() {
            *x = self.random();
        }
    }

    fn quick_sort(&mut self, arr: &mut Vec<u64>, left: usize, right: usize) {
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
        let mut arr = vec![0; 1000];
        for _ in 0..100 {
            self.randomize_array(&mut arr);
            self.quick_sort(&mut arr, 0, 999);
            checksum += arr[100];
        }
        checksum
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn check_checksum() {
        let mut qs = QuickSort::new(7);
        let checksum = qs.benchmark();
        assert_eq!(checksum, 21880255009);
    }
}
