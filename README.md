# Quicksort WASM benchmark

A benchmark for WebAssembly implementations of the Quicksort algorithm compared to Native Go and the EVM (Solidity).

**Languages:**

- TinyGo (opt=2, opt=s)
- Rust (opt-level=2, opt-level=s)
- AssemblyScript (optimizeLevel=3 shrinkLevel=1)
- Zig (ReleaseSmall, ReleaseFast)

**WASM Runtimes:**

- Wasmer (singlepass, cranelift)
- Wazero (interpreter, compiler)
- Wasm3
