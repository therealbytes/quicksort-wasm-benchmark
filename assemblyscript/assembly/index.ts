import { Quicksort, SEED } from "./sort";

export function run(): u64 {
  let qs = new Quicksort(SEED);
  const checksum = qs.benchmark();
  return checksum;
}
