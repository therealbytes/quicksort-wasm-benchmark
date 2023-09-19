.PHONY: all evm solidity wasm tinygo rust assemblyscript benchmark benchmark-native-rust

all: evm wasm benchmark

evm: solidity

solidity:
	forge build --optimizer-runs 1000 --sizes
	jq -r '.deployedBytecode.object' out/Quicksort.sol/Quicksort.json > testdata/quicksort.evm

wasm: tinygo rust assemblyscript

tinygo:
	tinygo build -opt=2 -no-debug -o testdata/tinygo_o2.wasm -target=wasi ./tinygo/main.go
	tinygo build -opt=z -no-debug -o testdata/tinygo_oz.wasm -target=wasi ./tinygo/main.go

rust:
	rustc -O -o testdata/rust-simple.wasm --target wasm32-unknown-unknown --crate-type cdylib rust/src/main.rs

assemblyscript:
	cd assemblyscript && npm run asbuild
	cp assemblyscript/build/release.wasm testdata/assemblyscript.wasm

benchmark:
	go test -bench . -benchmem | tee benchmark_output.txt
	echo "Benchmark,Size,Iterations,ns/op,Bytes/op,Allocs/op" > benchmark_results.csv
	awk '/Benchmark/ { print $$1 "," $$5 "," $$2 "," $$3 "," $$7 "," $$9 }' benchmark_output.txt >> benchmark_results.csv
	rm benchmark_output.txt

benchmark-native-rust:
	cd rust && cargo +nightly bench

