package main

import (
	_ "embed"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/concrete/precompiles"
	"github.com/ethereum/go-ethereum/concrete/wasm"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/therealbytes/concrete-sort/quicksort"
	"github.com/wasmerio/wasmer-go/wasmer"
)

func validResult(checksum int64) bool {
	return checksum == quicksort.CHECKSUM
}

func BenchmarkGo(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qs := quicksort.NewQuicksortBenchmark()
		checksum := int64(qs.Benchmark())
		if !validResult(checksum) {
			b.Fatal("invalid checksum:", checksum)
		}
	}
}

//go:embed testdata/quicksort.evm
var evmBytecodeHex []byte

func BenchmarkEVM(b *testing.B) {
	var (
		address        = common.HexToAddress("0xc0ffee")
		origin         = common.HexToAddress("0xc0ffee0001")
		bytecode       = common.Hex2Bytes(string(evmBytecodeHex)[2:])
		benchmarkInput = common.Hex2Bytes("8903c5a2")
		gasLimit       = uint64(1e9)
		txContext      = vm.TxContext{
			Origin:   origin,
			GasPrice: common.Big1,
		}
		context = vm.BlockContext{
			CanTransfer: core.CanTransfer,
			Transfer:    core.Transfer,
			Coinbase:    common.Address{},
			BlockNumber: common.Big1,
			Time:        1,
			Difficulty:  common.Big1,
			GasLimit:    uint64(1e8),
		}
	)

	statedb, err := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	if err != nil {
		b.Fatal(err)
	}

	statedb.CreateAccount(address)
	statedb.SetCode(address, bytecode)
	statedb.AddAddressToAccessList(address)
	statedb.CreateAccount(origin)
	statedb.SetBalance(origin, big.NewInt(1e18))

	evm := vm.NewEVM(context, txContext, statedb, params.TestChainConfig, vm.Config{})

	var ret []byte
	var gasLeft uint64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ret, gasLeft, err = evm.Call(vm.AccountRef(origin), address, benchmarkInput, gasLimit, common.Big0)
		if err != nil {
			b.Fatal(err)
		}
		b.ReportMetric(float64(gasLimit-gasLeft), "gas")
		checksum := new(big.Int).SetBytes(ret).Int64()
		if !validResult(checksum) {
			b.Fatal("invalid checksum:", checksum)
		}
	}
}

//go:embed testdata/tinygo.wasm
var tinygoWasmBytecode []byte

func BenchmarkTinygoQuicksort(b *testing.B) {
	runtimes := []struct {
		name string
		pc   precompiles.Precompile
	}{
		{"wazero", wasm.NewWazeroPrecompile(tinygoWasmBytecode)},
		{"wasmer/singlepass", wasm.NewWasmerPrecompileWithConfig(tinygoWasmBytecode, wasmer.NewConfig().UseSinglepassCompiler())},
		{"wasmer/cranelift", wasm.NewWasmerPrecompileWithConfig(tinygoWasmBytecode, wasmer.NewConfig().UseCraneliftCompiler())},
	}
	for _, runtime := range runtimes {
		b.Run(runtime.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ret, err := runtime.pc.Run(nil, nil)
				if err != nil {
					b.Fatal(err)
				}
				checksum := new(big.Int).SetBytes(ret).Int64()
				if !validResult(checksum) {
					b.Fatal("invalid checksum:", checksum)
				}
			}
		})
	}
}

//go:embed testdata/rust-simple.wasm
var rustWasmBytecode []byte

func BenchmarkRustQuicksort(b *testing.B) {
	compiler := []struct {
		name string
	}{
		{"cranelift"},
		{"singlepass"},
	}
	for _, compiler := range compiler {
		b.Run(fmt.Sprintf("wasmer/%s", compiler.name), func(b *testing.B) {
			config := wasmer.NewConfig()
			if compiler.name == "singlepass" {
				config.UseSinglepassCompiler()
			} else if compiler.name == "cranelift" {
				config.UseCraneliftCompiler()
			} else {
				b.Fatal("invalid compiler:", compiler.name)
			}

			engine := wasmer.NewEngineWithConfig(config)
			store := wasmer.NewStore(engine)
			module, err := wasmer.NewModule(store, rustWasmBytecode)
			if err != nil {
				b.Fatal(err)
			}
			importObject := wasmer.NewImportObject()
			instance, err := wasmer.NewInstance(module, importObject)
			if err != nil {
				b.Fatal(err)
			}
			run, err := instance.Exports.GetFunction("run")
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ret, err := run()
				if err != nil {
					b.Fatal(err)
				}
				checksum, ok := ret.(int64)
				if !ok {
					b.Fatal("can not convert return value to int64:", ret)
				}
				if !validResult(checksum) {
					b.Fatal("invalid checksum:", checksum)
				}
			}
		})
	}
}
