package main

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/concrete"
	"github.com/ethereum/go-ethereum/concrete/lib"
	"github.com/ethereum/go-ethereum/tinygo"
	"github.com/therealbytes/concrete-sort/quicksort"
)

type benchmarkPrecompile struct {
	lib.BlankPrecompile
}

func (t *benchmarkPrecompile) Run(env concrete.Environment, input []byte) ([]byte, error) {
	var (
		seed   = uint(input[0])
		arrLen = int(binary.BigEndian.Uint32(input[1:5]))
		iter   = int(binary.BigEndian.Uint32(input[5:9]))
	)
	benchmark := quicksort.NewQuicksortBenchmark(seed)
	checksum := benchmark.Run(arrLen, iter)
	output := make([]byte, 4)
	binary.BigEndian.PutUint32(output, uint32(checksum))
	return output, nil
}

func init() {
	tinygo.WasmWrap(&benchmarkPrecompile{})
}

// main is REQUIRED for TinyGo to compile to WASM
func main() {}
