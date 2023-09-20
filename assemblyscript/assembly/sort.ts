export class QuicksortBenchmark {
  seed: u32;

  constructor(seed: u32) {
    this.seed = seed;
  }

  random(): u32 {
    this.seed = (1103515245 * this.seed + 12345) % (1 << 31);
    return this.seed;
  }

  randomizeArray(arr: u32[]): void {
    for (let i = 0; i < arr.length; i++) {
      arr[i] = this.random() % 1000;
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

  run(arrLen: u32, iter: u32): u32 {
    let checksum: u32 = 0;
    let arr = new Array<u32>(arrLen).fill(0);
    for (let _: u32 = 0; _ < iter; _++) {
      this.randomizeArray(arr);
      this.quickSort(arr, 0, arrLen - 1);
      checksum += arr[arrLen / 2];
    }
    return checksum;
  }
}
