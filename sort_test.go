package main

import (
	_ "embed"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/concrete"
	"github.com/ethereum/go-ethereum/concrete/wasm"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/tetratelabs/wazero"
	"github.com/therealbytes/concrete-sort/quicksort"
	"github.com/wasmerio/wasmer-go/wasmer"
)

func validResult(checksum int64) bool {
	return checksum == quicksort.CHECKSUM
}

func reportCodeMetadata(b *testing.B, code []byte) {
	b.ReportMetric(float64(len(code)), "bytes")
}

func BenchmarkGo(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qs := quicksort.NewQuicksortBenchmark()
		checksum := int64(qs.Benchmark())
		if !validResult(checksum) {
			b.Fatal("invalid checksum:", checksum)
		}
		reportCodeMetadata(b, []byte{})
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
		checksum := new(big.Int).SetBytes(ret).Int64()
		if !validResult(checksum) {
			b.Fatal("invalid checksum:", checksum)
		}
		// b.ReportMetric(float64(gasLimit-gasLeft), "gas")
		_ = gasLeft
		reportCodeMetadata(b, evmBytecodeHex)
	}
}

//go:embed testdata/tinygo_o2.wasm
var tinygoWasmBytecode_o2 []byte

//go:embed testdata/tinygo_oz.wasm
var tinygoWasmBytecode_oz []byte

func BenchmarkTinygoQuicksort(b *testing.B) {
	runtimes := []struct {
		name string
		pc   concrete.Precompile
		code []byte
	}{
		{"wazero/interpreted/o2", wasm.NewWazeroPrecompileWithConfig(tinygoWasmBytecode_o2, wazero.NewRuntimeConfigInterpreter()), tinygoWasmBytecode_o2},
		{"wazero/interpreted/oz", wasm.NewWazeroPrecompileWithConfig(tinygoWasmBytecode_oz, wazero.NewRuntimeConfigInterpreter()), tinygoWasmBytecode_oz},
		{"wazero/compiled/o2", wasm.NewWazeroPrecompileWithConfig(tinygoWasmBytecode_o2, wazero.NewRuntimeConfigCompiler()), tinygoWasmBytecode_o2},
		{"wazero/compiled/oz", wasm.NewWazeroPrecompileWithConfig(tinygoWasmBytecode_oz, wazero.NewRuntimeConfigCompiler()), tinygoWasmBytecode_oz},
		{"wasmer/singlepass/o2", wasm.NewWasmerPrecompileWithConfig(tinygoWasmBytecode_o2, wasmer.NewConfig().UseSinglepassCompiler()), tinygoWasmBytecode_o2},
		{"wasmer/singlepass/oz", wasm.NewWasmerPrecompileWithConfig(tinygoWasmBytecode_oz, wasmer.NewConfig().UseSinglepassCompiler()), tinygoWasmBytecode_oz},
		{"wasmer/cranelift/o2", wasm.NewWasmerPrecompileWithConfig(tinygoWasmBytecode_o2, wasmer.NewConfig().UseCraneliftCompiler()), tinygoWasmBytecode_o2},
		{"wasmer/cranelift/oz", wasm.NewWasmerPrecompileWithConfig(tinygoWasmBytecode_oz, wasmer.NewConfig().UseCraneliftCompiler()), tinygoWasmBytecode_oz},
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
				reportCodeMetadata(b, runtime.code)
			}
		})
	}
}

func newWasmerInstance(code []byte, config *wasmer.Config) (*wasmer.Instance, error) {
	engine := wasmer.NewEngineWithConfig(config)
	store := wasmer.NewStore(engine)
	module, err := wasmer.NewModule(store, code)
	if err != nil {
		return nil, err
	}

	importObject := wasmer.NewImportObject()
	importObject.Register(
		"env", map[string]wasmer.IntoExtern{
			"abort": wasmer.NewFunction(
				store,
				wasmer.NewFunctionType(
					wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32),
					wasmer.NewValueTypes(),
				),
				func([]wasmer.Value) ([]wasmer.Value, error) {
					return nil, fmt.Errorf("abort")
				},
			),
		},
	)

	instance, err := wasmer.NewInstance(module, importObject)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func newBenchWasmerInstance(b *testing.B, code []byte, config *wasmer.Config) *wasmer.Instance {
	instance, err := newWasmerInstance(code, config)
	if err != nil {
		b.Fatal(err)
	}
	return instance
}

func benchWasmerInstance(b *testing.B, instance *wasmer.Instance, code []byte) {
	for i := 0; i < b.N; i++ {
		run, err := instance.Exports.GetFunction("run")
		if err != nil {
			b.Fatal(err)
		}
		ret, err := run()
		if err != nil {
			b.Fatal(err)
		}
		checksum, ok := ret.(int64)
		if !ok {
			b.Fatal("can not convert return value to int64")
		}
		if !validResult(checksum) {
			b.Fatal("invalid checksum:", checksum)
		}
		reportCodeMetadata(b, code)
	}
}

//go:embed testdata/rust-simple.wasm
var rustWasmBytecode []byte

//go:embed testdata/assemblyscript.wasm
var assemblyScriptBytecode []byte

func BenchmarkWasmRustQuicksort(b *testing.B) {
	benchCases := []struct {
		name     string
		instance *wasmer.Instance
	}{
		{"wasmer/singlepass", newBenchWasmerInstance(b, rustWasmBytecode, wasmer.NewConfig().UseSinglepassCompiler())},
		{"wasmer/cranelift", newBenchWasmerInstance(b, rustWasmBytecode, wasmer.NewConfig().UseCraneliftCompiler())},
	}
	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			b.ResetTimer()
			benchWasmerInstance(b, bc.instance, rustWasmBytecode)
		})
	}
}

func BenchmarkWasmAssemblyScriptQuicksort(b *testing.B) {
	benchCases := []struct {
		name     string
		instance *wasmer.Instance
	}{
		{"wasmer/singlepass", newBenchWasmerInstance(b, assemblyScriptBytecode, wasmer.NewConfig().UseSinglepassCompiler())},
		{"wasmer/cranelift", newBenchWasmerInstance(b, assemblyScriptBytecode, wasmer.NewConfig().UseCraneliftCompiler())},
	}
	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			b.ResetTimer()
			benchWasmerInstance(b, bc.instance, assemblyScriptBytecode)
		})
	}
}
