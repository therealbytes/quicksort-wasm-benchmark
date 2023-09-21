import { QuicksortBenchmark } from "./sort";

export function run(seed: i32, arrLen: i32, iter: i32): i32 {
  let qs = new QuicksortBenchmark(<u32>seed);
  const checksum = qs.run(<u32>arrLen, <u32>iter);
  return <i32>checksum;
}
