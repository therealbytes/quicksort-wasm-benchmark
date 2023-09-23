package main

import (
	"context"
	_ "embed"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/matiasinsaurralde/go-wasm3"
	"github.com/tetratelabs/wazero"
	wz_api "github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"github.com/therealbytes/quicksort-wasm-benchmark/quicksort"
	"github.com/wasmerio/wasmer-go/wasmer"
)

type LanguageName string

const (
	Go             LanguageName = "Go"
	Sol            LanguageName = "Solidity"
	TinyGo         LanguageName = "TinyGo"
	Rust           LanguageName = "Rust"
	AssemblyScript LanguageName = "AssemblyScript"
	Zig            LanguageName = "Zig"
)

type WasmRuntimeName string

const (
	Wazero WasmRuntimeName = "Wazero"
	Wasmer WasmRuntimeName = "Wasmer"
	Wasm3  WasmRuntimeName = "Wasm3"
)

var (
	Seed   = getEnvVarInt("SEED", 0)
	ArrLen = getEnvVarInt("ARR_LEN", 1000)
	Iter   = getEnvVarInt("ITER", 100)
)

var (
	// Native
	BenchNative = getEnvVarBool("NATIVE", true)
	BenchGo     = getEnvVarBool("GO", BenchNative)
	// EVM
	BenchEVM = getEnvVarBool("EVM", true)
	BenchSol = getEnvVarBool("SOLIDITY", BenchEVM)
	// Wasm
	BenchAllLangs = getEnvVarBool("ALL_LANGS", true)
	BenchTinyGo   = getEnvVarBool("TINYGO", BenchAllLangs)
	BenchRust     = getEnvVarBool("RUST", BenchAllLangs)
	BenchAs       = getEnvVarBool("ASSEMBLYSCRIPT", BenchAllLangs)
	BenchZig      = getEnvVarBool("ZIG", BenchAllLangs)
	// Runtimes
	BenchAllRuntimes = getEnvVarBool("ALL_RUNTIMES", true)
	BenchWazero      = getEnvVarBool("WAZERO", BenchAllRuntimes)
	BenchWasmer      = getEnvVarBool("WASMER", BenchAllRuntimes)
	BenchWasm3       = getEnvVarBool("WASM3", BenchAllRuntimes)
)

func getEnvVarInt(name string, defaultValue int) int {
	strVal := os.Getenv(name)
	if strVal == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		panic(err)
	}
	return intVal
}

func getEnvVarBool(name string, defaultValue bool) bool {
	strVal := os.Getenv(name)
	if strVal == "" {
		return defaultValue
	}
	boolVal, err := strconv.ParseBool(strVal)
	if err != nil {
		panic(err)
	}
	return boolVal
}

//go:embed testdata/solidity.evm
var evmBytecodeHex []byte

//go:embed testdata/tinygo_o2.wasm
var tinygoWasmBytecode_o2 []byte

//go:embed testdata/tinygo_oz.wasm
var tinygoWasmBytecode_oz []byte

//go:embed testdata/rust_o2.wasm
var rustWasmBytecode_o2 []byte

//go:embed testdata/rust_os.wasm
var rustWasmBytecode_os []byte

//go:embed testdata/assemblyscript.wasm
var assemblyScriptBytecode []byte

//go:embed testdata/zig_fast.wasm
var zigBytecode_fast []byte

//go:embed testdata/zig_small.wasm
var zigBytecode_small []byte

type Runtime struct {
	Name      WasmRuntimeName
	ConfigStr string
}

type Binary struct {
	Code               []byte
	CompilerOptionsStr string
}

type Benchmark struct {
	Language LanguageName
	Runtimes []Runtime
	Binaries []Binary
}

type BenchmarkRunner interface {
	Run(seed int, arrLen int, iter int) (uint, error)
}

var (
	WazeroCompiled    = Runtime{Name: Wazero, ConfigStr: "Compiled"}
	WazeroInterpreted = Runtime{Name: Wazero, ConfigStr: "Interpreted"}
	WasmerCranelift   = Runtime{Name: Wasmer, ConfigStr: "Cranelift"}
	WasmerSinglepass  = Runtime{Name: Wasmer, ConfigStr: "Singlepass"}
	Wasm3Interpreted  = Runtime{Name: Wasm3, ConfigStr: "Interpreted"}
)

var (
	TinyGoO2 = Binary{Code: tinygoWasmBytecode_o2, CompilerOptionsStr: "o2"}
	TinyGoOz = Binary{Code: tinygoWasmBytecode_oz, CompilerOptionsStr: "oz"}
	RustO2   = Binary{Code: rustWasmBytecode_o2, CompilerOptionsStr: "o2"}
	RustOs   = Binary{Code: rustWasmBytecode_os, CompilerOptionsStr: "os"}
	As       = Binary{Code: assemblyScriptBytecode, CompilerOptionsStr: "optimize=3shrink=1"}
	ZigFast  = Binary{Code: zigBytecode_fast, CompilerOptionsStr: "ReleaseFast"}
	ZigSmall = Binary{Code: zigBytecode_small, CompilerOptionsStr: "ReleaseSmall"}
)

var WasmBenchmarks = []Benchmark{
	{
		Language: TinyGo,
		Runtimes: []Runtime{
			WazeroCompiled,
			WazeroInterpreted,
			WasmerCranelift,
			WasmerSinglepass,
		},
		Binaries: []Binary{
			TinyGoO2,
			TinyGoOz,
		},
	},
	{
		Language: Rust,
		Runtimes: []Runtime{
			WazeroCompiled,
			WazeroInterpreted,
			WasmerCranelift,
			WasmerSinglepass,
			Wasm3Interpreted,
		},
		Binaries: []Binary{
			RustO2,
			RustOs,
		},
	},
	{
		Language: AssemblyScript,
		Runtimes: []Runtime{
			WazeroCompiled,
			WazeroInterpreted,
			WasmerCranelift,
			WasmerSinglepass,
		},
		Binaries: []Binary{
			As,
		},
	},
	{
		Language: Zig,
		Runtimes: []Runtime{
			WazeroCompiled,
			WazeroInterpreted,
			WasmerCranelift,
			WasmerSinglepass,
		},
		Binaries: []Binary{
			ZigFast,
			ZigSmall,
		},
	},
}

var (
	ExpectedChecksum = int(quicksort.NewQuicksortBenchmark(uint(Seed)).Run(ArrLen, Iter))
)

func validResult(checksum int) bool {
	return checksum == ExpectedChecksum
}

func reportCodeMetadata(b *testing.B, code []byte) {
	b.ReportMetric(float64(len(code)), "bytes")
}

func BenchmarkGo(b *testing.B) {
	if !BenchGo {
		b.SkipNow()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmark := quicksort.NewQuicksortBenchmark(uint(Seed))
		checksum := benchmark.Run(ArrLen, Iter)
		if !validResult(int(checksum)) {
			b.Fatal("invalid checksum:", checksum)
		}
		reportCodeMetadata(b, []byte{})
	}
}

func BenchmarkEVM(b *testing.B) {
	if !BenchSol {
		b.SkipNow()
	}
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
	input = append(input, math.U256Bytes(big.NewInt(int64(Seed)))...)
	input = append(input, math.U256Bytes(big.NewInt(int64(ArrLen)))...)
	input = append(input, math.U256Bytes(big.NewInt(int64(Iter)))...)

	var ret []byte
	var gasLeft uint64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ret, gasLeft, err = evm.Call(vm.AccountRef(origin), address, input, gasLimit, common.Big0)
		if err != nil {
			b.Fatal(err)
		}
		checksum := int(new(big.Int).SetBytes(ret).Int64())
		if !validResult(checksum) {
			b.Fatal("invalid checksum:", checksum)
		}
		// b.ReportMetric(float64(gasLimit-gasLeft), "gas")
		_ = gasLeft
		reportCodeMetadata(b, evmBytecodeHex)
	}
}

func BenchmarkWasm(b *testing.B) {
	for _, benchmark := range WasmBenchmarks {
		for _, runtime := range benchmark.Runtimes {
			for _, binary := range benchmark.Binaries {
				name := fmt.Sprintf("%s_%s_%s_%s", benchmark.Language, runtime.Name, runtime.ConfigStr, binary.CompilerOptionsStr)
				b.Run(name, func(b *testing.B) {
					switch {
					case benchmark.Language == TinyGo && !BenchTinyGo:
						b.SkipNow()
					case benchmark.Language == Rust && !BenchRust:
						b.SkipNow()
					case benchmark.Language == AssemblyScript && !BenchAs:
						b.SkipNow()
					case benchmark.Language == Zig && !BenchZig:
						b.SkipNow()
					case runtime.Name == Wazero && !BenchWazero:
						b.SkipNow()
					case runtime.Name == Wasmer && !BenchWasmer:
						b.SkipNow()
					case runtime.Name == Wasm3 && !BenchWasm3:
						b.SkipNow()
					}
					var runner BenchmarkRunner
					switch runtime.Name {
					case Wazero:
						var config wazero.RuntimeConfig
						switch runtime.ConfigStr {
						case "Compiled":
							config = wazero.NewRuntimeConfigCompiler()
						case "Interpreted":
							config = wazero.NewRuntimeConfigInterpreter()
						default:
							b.Fatal("unknown runtime config:", runtime.ConfigStr)
						}
						runner = newWazeroRunner(b, binary.Code, config, benchmark.Language)

					case Wasmer:
						var config *wasmer.Config
						switch runtime.ConfigStr {
						case "Cranelift":
							config = wasmer.NewConfig().UseCraneliftCompiler()
						case "Singlepass":
							config = wasmer.NewConfig().UseSinglepassCompiler()
						default:
							b.Fatal("unknown runtime config:", runtime.ConfigStr)
						}
						runner = newWasmerRunner(b, binary.Code, config, benchmark.Language)

					case Wasm3:
						runner = newWasm3Runner(b, binary.Code, benchmark.Language)

					default:
						b.Fatal("unknown runtime:", runtime.Name)
					}

					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						checksum, err := runner.Run(Seed, ArrLen, Iter)
						if err != nil {
							b.Fatal(err)
						}
						if !validResult(int(checksum)) {
							b.Fatal("invalid checksum:", checksum)
						}
						reportCodeMetadata(b, binary.Code)
					}
				})
			}
		}
	}
}

type wazeroRunner struct {
	run wz_api.Function
}

func newWazeroRunner(b *testing.B, code []byte, config wazero.RuntimeConfig, lang LanguageName) BenchmarkRunner {
	var err error
	ctx := context.Background()

	r := wazero.NewRuntimeWithConfig(ctx, config)

	if lang == AssemblyScript {
		_, err = r.NewHostModuleBuilder("env").
			NewFunctionBuilder().WithFunc(func(int32, int32, int32, int32) {}).Export("abort").
			Instantiate(ctx)
	} else {
		_, err = r.NewHostModuleBuilder("env").Instantiate(ctx)
	}
	if err != nil {
		b.Fatal(err)
	}
	if lang == TinyGo {
		wasi_snapshot_preview1.MustInstantiate(ctx, r)
	}

	mod, err := r.Instantiate(ctx, code)
	if err != nil {
		b.Fatal(err)
	}
	run := mod.ExportedFunction("run")
	if err != nil {
		b.Fatal(err)
	}

	return &wazeroRunner{
		run: run,
	}
}

func (r *wazeroRunner) Run(seed int, arrLen int, iter int) (uint, error) {
	ctx := context.Background()
	ret, err := r.run.Call(ctx, uint64(seed), uint64(arrLen), uint64(iter))
	if err != nil {
		return 0, err
	}
	checksum := uint(ret[0])
	return checksum, nil
}

type wasmerRunner struct {
	run wasmer.NativeFunction
}

func newWasmerRunner(b *testing.B, code []byte, config *wasmer.Config, lang LanguageName) BenchmarkRunner {
	engine := wasmer.NewEngineWithConfig(config)
	store := wasmer.NewStore(engine)
	module, err := wasmer.NewModule(store, code)
	if err != nil {
		b.Fatal(err)
	}

	var importObject *wasmer.ImportObject
	if lang == TinyGo {
		wasiEnv, err := wasmer.NewWasiStateBuilder("wasi-program").Finalize()
		if err != nil {
			b.Fatal(err)
		}
		importObject, err = wasiEnv.GenerateImportObject(store, module)
		if err != nil {
			b.Fatal(err)
		}
	} else {
		importObject = wasmer.NewImportObject()
	}
	if lang == AssemblyScript {
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
	}

	instance, err := wasmer.NewInstance(module, importObject)
	if err != nil {
		b.Fatal(err)
	}
	run, err := instance.Exports.GetFunction("run")
	if err != nil {
		b.Fatal(err)
	}

	return &wasmerRunner{
		run: run,
	}
}

func (r *wasmerRunner) Run(seed int, arrLen int, iter int) (uint, error) {
	ret, err := r.run(int32(seed), int32(arrLen), int32(iter))
	if err != nil {
		return 0, err
	}
	_checksum, ok := ret.(int32)
	if !ok {
		return 0, fmt.Errorf("can not convert return value to int32")
	}
	checksum := uint(_checksum)
	return checksum, nil
}

type wasm3Runner struct {
	run wasm3.FunctionWrapper
}

func newWasm3Runner(b *testing.B, code []byte, lang LanguageName) BenchmarkRunner {
	runtime := wasm3.NewRuntime(&wasm3.Config{
		Environment: wasm3.NewEnvironment(),
		StackSize:   64 * 1024,
	})
	module, err := runtime.ParseModule(code)
	if err != nil {
		b.Fatal(err)
	}
	_, err = runtime.LoadModule(module)
	if err != nil {
		b.Fatal(err)
	}
	run, err := runtime.FindFunction("run")
	if err != nil {
		b.Fatal(err)
	}
	return &wasm3Runner{
		run: run,
	}
}

func (r *wasm3Runner) Run(seed int, arrLen int, iter int) (uint, error) {
	ret, err := r.run(int(seed), int(arrLen), int(iter))
	if err != nil {
		return 0, err
	}
	checksum := uint(ret)
	return checksum, nil
}
