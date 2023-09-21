package main

import (
	"github.com/therealbytes/concrete-sort/quicksort"
)

//export run
func run(seed int32, arrLen int32, iter int32) int32 {
	benchmark := quicksort.NewQuicksortBenchmark(uint(seed))
	checksum := benchmark.Run(int(arrLen), int(iter))
	return int32(checksum)
}

// main is REQUIRED for TinyGo to compile to WASM
func main() {}
