.PHONY: all evm solidity wasm tinygo rust benchmark

all: evm wasm benchmark

evm: solidity

solidity:
	forge build --optimizer-runs 1000
	jq -r '.deployedBytecode.object' out/Quicksort.sol/Quicksort.json > testdata/quicksort.evm

wasm: tinygo rust

tinygo:
	tinygo build -opt=2 -o testdata/tinygo.wasm -target=wasi ./tinygo/main.go

rust:
	rustc -O -o testdata/rust-simple.wasm --target wasm32-unknown-unknown --crate-type cdylib rust/src/main.rs

benchmark:
	go test -bench . > benchmark.txt
	echo "\nNative Rust" >> benchmark.txt
	cd rust && cargo +nightly bench >> ../benchmark.txt
