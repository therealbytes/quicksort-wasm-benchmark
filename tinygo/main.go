package main

import (
	"github.com/therealbytes/concrete-sort/quicksort"
)

//go:wasm-module env
//export concrete_Environment
func run() int64 {
	qs := quicksort.NewQuicksortBenchmark(7)
	result := qs.Benchmark()
	return int64(result)
}

func main() {}
