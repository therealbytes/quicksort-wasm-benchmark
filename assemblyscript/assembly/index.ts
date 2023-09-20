import { QuicksortBenchmark } from "./sort";

export function run(seed: i64, arrLen: i64, iter: i64): i64 {
  let qs = new QuicksortBenchmark(<u32>seed);
  const checksum = qs.run(<u32>arrLen, <u32>iter);
  return <i64>checksum;
}
