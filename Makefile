.PHONY: all evm solidity wasm tinygo rust assemblyscript benchmark

all: evm wasm benchmark

evm: solidity

solidity:
	forge build --optimizer-runs 1000 --sizes
	jq -r '.deployedBytecode.object' out/Quicksort.sol/Quicksort.json > testdata/quicksort.evm

wasm: tinygo rust

tinygo:
	tinygo build -opt=2 -no-debug -o testdata/tinygo_o2.wasm -target=wasi ./tinygo/main.go
	tinygo build -opt=z -no-debug -o testdata/tinygo_oz.wasm -target=wasi ./tinygo/main.go

rust:
	rustc -O -o testdata/rust-simple.wasm --target wasm32-unknown-unknown --crate-type cdylib rust/src/main.rs

assemblyscript:
	cd assemblyscript && npm run asbuild
	cp assemblyscript/build/release.wasm testdata/assemblyscript.wasm

benchmark:
	go test -bench . > benchmark.txt
	echo "\nNative Rust" >> benchmark.txt
	cd rust && cargo +nightly bench >> ../benchmark.txt
