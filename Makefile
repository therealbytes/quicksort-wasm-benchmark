.PHONY: evm solidity wasm tinygo rust

all: evm wasm

evm: solidity

solidity:
	forge build --optimizer-runs 1000
	jq -r '.deployedBytecode.object' out/QuickSort.sol/QuickSort.json > testdata/quicksort.evm

wasm: tinygo rust

tinygo:
	tinygo build -opt=2 -o testdata/tinygo.wasm -target wasi ./tinygo/main.go

rust:
	rustc -O -o testdata/rust-simple.wasm --target wasm32-unknown-unknown --crate-type cdylib rust/src/main.rs

benchmark:
	go test -bench . > benchmark.txt
