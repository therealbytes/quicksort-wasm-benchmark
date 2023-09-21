.PHONY: all prepare evm solidity wasm tinygo rust zig assemblyscript benchmark benchmark-native-rust repro-tinygo-issue

all: prepare evm wasm benchmark-arrlen-many

prepare:
	mkdir -p testdata
	mkdir -p results

evm: solidity

solidity:
	forge build --optimizer-runs 1000 --sizes
	jq -r '.deployedBytecode.object' out/QuicksortBenchmark.sol/QuicksortBenchmark.json > testdata/solidity.evm

wasm: tinygo rust assemblyscript zig

tinygo:
	tinygo build -opt=2 -no-debug -o testdata/tinygo_o2.wasm -target=wasi tinygo/main.go
	tinygo build -opt=s -no-debug -o testdata/tinygo_oz.wasm -target=wasi tinygo/main.go

rust:
	rustc -C opt-level=2 -o testdata/rust_o2.wasm --target wasm32-unknown-unknown --crate-type cdylib rust/src/main.rs
	rustc -C opt-level=s -o testdata/rust_os.wasm --target wasm32-unknown-unknown --crate-type cdylib rust/src/main.rs

zig:
	cd zig && zig build -Doptimize=ReleaseFast
	cp zig/zig-out/lib/zig.wasm testdata/zig_fast.wasm
	cd zig && zig build -Doptimize=ReleaseSmall
	cp zig/zig-out/lib/zig.wasm testdata/zig_small.wasm

assemblyscript:
	cd assemblyscript && npm run asbuild
	cp assemblyscript/build/release.wasm testdata/assemblyscript.wasm

benchmark:
	go test -bench . -benchmem | tee benchmark_output.txt
	echo "Benchmark,Size,Iterations,ns/op,Bytes/op,Allocs/op" > results/benchmark_results.csv
	awk '/Benchmark/ { print $$1 "," $$5 "," $$2 "," $$3 "," $$7 "," $$9 }' benchmark_output.txt >> results/benchmark_results.csv
	rm benchmark_output.txt

benchmark-arrlen: benchmark
	mv results/benchmark_results.csv results/benchmark_results_$(ARR_LEN).csv

benchmark-arrlen-many:
	# ARR_LEN=10 $(MAKE) benchmark-arrlen
	ARR_LEN=100 $(MAKE) benchmark-arrlen
	ARR_LEN=1000 $(MAKE) benchmark-arrlen

benchmark-native-rust:
	cd rust && cargo +nightly bench

repro-tinygo-issue:
	ARR_LEN=10 ITER=1 NATIVE=false EVM=false ALL_LANGS=false ALL_RUNTIMES=false WASMER=true TINYGO=true go test -bench .
