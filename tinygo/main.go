package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/concrete/api"
	"github.com/ethereum/go-ethereum/concrete/lib"
	"github.com/ethereum/go-ethereum/tinygo"
	"github.com/therealbytes/concrete-sort/quicksort"
)

type quicksortPrecompile struct {
	lib.BlankPrecompile
}

func (t *quicksortPrecompile) Run(env api.Environment, input []byte) ([]byte, error) {
	b := quicksort.NewQuicksortBenchmark()
	checksum := b.Benchmark()
	checksumBN := big.NewInt(int64(checksum))
	return common.BigToHash(checksumBN).Bytes(), nil
}

func init() {
	tinygo.WasmWrap(&quicksortPrecompile{})
}

// main is REQUIRED for TinyGo to compile to WASM
func main() {}
