export const SEED: u64 = 7;
export const L: i32 = 1000;
export const N: i32 = 100;
export const CHECKSUM: u64 = 107829970005;

export class Quicksort {
  seed: u64;

  constructor(seed: u64) {
    this.seed = seed;
  }

  random(): u32 {
    this.seed = (1103515245 * this.seed + 12345) % (1 << 31);
    return <u32>(this.seed % (<u64>u32.MAX_VALUE + 1));
  }

  randomizeArray(arr: u32[]): void {
    for (let i = 0; i < arr.length; i++) {
      arr[i] = this.random();
    }
  }

  quickSort(arr: u32[], left: i32, right: i32): void {
    let i: i32 = left;
    let j: i32 = right;

    if (i == j) {
      return;
    }

    let pivot: u32 = arr[left + (right - left) / 2];

    while (i <= j) {
      while (arr[i] < pivot) {
        i++;
      }
      while (pivot < arr[j]) {
        j--;
      }
      if (i <= j) {
        let temp: u32 = arr[i];
        arr[i] = arr[j];
        arr[j] = temp;
        i++;
        j--;
      }
    }

    if (left < j) {
      this.quickSort(arr, left, j);
    }
    if (i < right) {
      this.quickSort(arr, i, right);
    }
  }

  benchmark(): u64 {
    let checksum: u64 = 0;
    let arr = new Array<u32>(L).fill(0);
    for (let _ = 0; _ < N; _++) {
      this.randomizeArray(arr);
      this.quickSort(arr, 0, L - 1);
      checksum += <u64>arr[L / 2];
    }
    return checksum;
  }
}
