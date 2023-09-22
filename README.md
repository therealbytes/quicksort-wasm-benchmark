# Quicksort WASM benchmark

A benchmark for WebAssembly runtimes and the EVM based on quicksort.

[Results](./results/benchmark_results_1000.csv) (run on a Intel Core i5 2020 MacBook Pro)

**Languages:**

- Go (Native)
- Solidity (EVM)
- TinyGo (WASM) (opt=2, opt=s)
- Rust (WASM) (opt-level=2, opt-level=s)
- AssemblyScript (WASM) (optimizeLevel=3 shrinkLevel=1)
- Zig (WASM) (ReleaseSmall, ReleaseFast)

**WASM Runtimes:**

- Wasmer (singlepass, cranelift)
- Wazero (interpreter, compiler)
- Wasm3
