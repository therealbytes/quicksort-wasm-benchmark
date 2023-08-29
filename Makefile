.PHONY: evm solidity wasm tinygo

evm: solidity

solidity:
	forge build --optimizer-runs 1000
	jq -r '.deployedBytecode.object' out/QuickSort.sol/QuickSort.json > testdata/quicksort.evm

wasm: tinygo

tinygo:
	tinygo build -opt=2 -o testdata/tinygo.wasm -target wasi ./tinygo/main.go
