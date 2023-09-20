package main

import (
	_ "embed"
	"encoding/binary"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
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

var (
	seed             uint = 7
	arrLen                = 1000
	iter                  = 100
	expectedChecksum uint = 0
)

func init() {
	benchmark := quicksort.NewQuicksortBenchmark(seed)
	expectedChecksum = benchmark.Run(arrLen, iter)
}

func validResult(checksum uint) bool {
	return checksum == expectedChecksum
}

func reportCodeMetadata(b *testing.B, code []byte) {
	b.ReportMetric(float64(len(code)), "bytes")
}

func BenchmarkGo(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmark := quicksort.NewQuicksortBenchmark(seed)
		checksum := benchmark.Run(arrLen, iter)
		if !validResult(checksum) {
			b.Fatal("invalid checksum:", checksum)
		}
		reportCodeMetadata(b, []byte{})
	}
}

//go:embed testdata/solidity.evm
var evmBytecodeHex []byte

func BenchmarkEVM(b *testing.B) {
	var (
		address   = common.HexToAddress("0xc0ffee")
		origin    = common.HexToAddress("0xc0ffee0001")
		bytecode  = common.Hex2Bytes(string(evmBytecodeHex)[2:])
		gasLimit  = uint64(1e9)
		txContext = vm.TxContext{
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

	input := common.Hex2Bytes("24b912e5")
	input = append(input, math.U256Bytes(big.NewInt(int64(seed)))...)
	input = append(input, math.U256Bytes(big.NewInt(int64(arrLen)))...)
	input = append(input, math.U256Bytes(big.NewInt(int64(iter)))...)

	var ret []byte
	var gasLeft uint64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ret, gasLeft, err = evm.Call(vm.AccountRef(origin), address, input, gasLimit, common.Big0)
		if err != nil {
			b.Fatal(err)
		}
		checksum := uint(new(big.Int).SetBytes(ret).Int64())
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

func BenchmarkWasmTinygo(b *testing.B) {
	runtimes := []struct {
		name string
		code []byte
		pc   concrete.Precompile
	}{
		{"wazero/interpreted/o2", tinygoWasmBytecode_o2, wasm.NewWazeroPrecompileWithConfig(tinygoWasmBytecode_o2, wazero.NewRuntimeConfigInterpreter())},
		{"wazero/interpreted/oz", tinygoWasmBytecode_oz, wasm.NewWazeroPrecompileWithConfig(tinygoWasmBytecode_oz, wazero.NewRuntimeConfigInterpreter())},
		{"wazero/compiled/o2", tinygoWasmBytecode_o2, wasm.NewWazeroPrecompileWithConfig(tinygoWasmBytecode_o2, wazero.NewRuntimeConfigCompiler())},
		{"wazero/compiled/oz", tinygoWasmBytecode_oz, wasm.NewWazeroPrecompileWithConfig(tinygoWasmBytecode_oz, wazero.NewRuntimeConfigCompiler())},
		{"wasmer/singlepass/o2", tinygoWasmBytecode_o2, wasm.NewWasmerPrecompileWithConfig(tinygoWasmBytecode_o2, wasmer.NewConfig().UseSinglepassCompiler())},
		{"wasmer/singlepass/oz", tinygoWasmBytecode_oz, wasm.NewWasmerPrecompileWithConfig(tinygoWasmBytecode_oz, wasmer.NewConfig().UseSinglepassCompiler())},
		{"wasmer/cranelift/o2", tinygoWasmBytecode_o2, wasm.NewWasmerPrecompileWithConfig(tinygoWasmBytecode_o2, wasmer.NewConfig().UseCraneliftCompiler())},
		{"wasmer/cranelift/oz", tinygoWasmBytecode_oz, wasm.NewWasmerPrecompileWithConfig(tinygoWasmBytecode_oz, wasmer.NewConfig().UseCraneliftCompiler())},
	}

	input := make([]byte, 9)
	input[0] = byte(seed)
	binary.BigEndian.PutUint32(input[1:5], uint32(arrLen))
	binary.BigEndian.PutUint32(input[5:9], uint32(iter))

	for _, runtime := range runtimes {
		b.Run(runtime.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				ret, err := runtime.pc.Run(nil, input)
				if err != nil {
					b.Fatal(err)
				}
				checksum := uint(binary.BigEndian.Uint32(ret))
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
		ret, err := run(int64(seed), int64(arrLen), int64(iter))
		if err != nil {
			b.Fatal(err)
		}
		_checksum, ok := ret.(int64)
		if !ok {
			b.Fatal("can not convert return value to int64")
		}
		checksum := uint(_checksum)
		if !validResult(checksum) {
			b.Fatal("invalid checksum:", checksum)
		}
		reportCodeMetadata(b, code)
	}
}

//go:embed testdata/rust.wasm
var rustWasmBytecode []byte

func BenchmarkWasmRust(b *testing.B) {
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

//go:embed testdata/assemblyscript.wasm
var assemblyScriptBytecode []byte

func BenchmarkWasmAssemblyScript(b *testing.B) {
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

//go:embed testdata/zig_fast.wasm
var zigBytecode_fast []byte

//go:embed testdata/zig_small.wasm
var zigBytecode_small []byte

func BenchmarkWasmZig(b *testing.B) {
	benchCases := []struct {
		name     string
		code     []byte
		instance *wasmer.Instance
	}{
		{"wasmer/singlepass/fast", zigBytecode_fast, newBenchWasmerInstance(b, zigBytecode_fast, wasmer.NewConfig().UseSinglepassCompiler())},
		{"wasmer/singlepass/small", zigBytecode_small, newBenchWasmerInstance(b, zigBytecode_small, wasmer.NewConfig().UseSinglepassCompiler())},
		{"wasmer/cranelift/fast", zigBytecode_fast, newBenchWasmerInstance(b, zigBytecode_fast, wasmer.NewConfig().UseCraneliftCompiler())},
		{"wasmer/cranelift/small", zigBytecode_small, newBenchWasmerInstance(b, zigBytecode_small, wasmer.NewConfig().UseCraneliftCompiler())},
	}
	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			b.ResetTimer()
			benchWasmerInstance(b, bc.instance, bc.code)
		})
	}
}
