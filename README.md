**Quicksort L=1000 N=10**

| name                    | runs | ns/op     | /native-go |
|-------------------------|------|-----------|------------|
| Native Go               | 2391 | 496255    | 1.00       |
| Native Rust             | 7    | 502429    | 1.01       |
| Wasmer Tinygo Cranelift | 2047 | 604185    | 1.22       |
| Wasmer Rust Cranelift   | 2082 | 569174    | 1.15       |
| Wazero Tinygo           | 1038 | 1148140   | 2.31       |
| Wasmer Tinygo Singlepass| 1056 | 1145900   | 2.31       |
| Wasmer Rust Singlepass  | 1117 | 1081074   | 2.18       |
| EVM                     | 7    | 167516639 | 337.56     |

**Quicksort L=1000 N=100**

| name                    | runs | ns/op     | /native-go |
|-------------------------|------|-----------|------------|
| Native Go               | 237  | 5012137   | 1.00       |
| Native Rust             | 1    | 5235500   | 1.04       |
| Wasmer Tinygo Cranelift | 210  | 5706034   | 1.14       |
| Wasmer Rust Cranelift   | 208  | 5711718   | 1.14       |
| Wazero Tinygo           | 99   | 12183878  | 2.43       |
| Wasmer Tinygo Singlepass| 81   | 12725835  | 2.54       |
| Wasmer Rust Singlepass  | 100  | 11274474  | 2.25       |
| EVM                     | 1    | 1761466683| 351.44     |

**Quicksort L=100 N=100**

| name                    | runs | ns/op     | /native-go |
|-------------------------|------|-----------|------------|
| Native Go               | 2391 | 496255    | 1.00       |
| Native Rust             | 7    | 502429    | 1.01       |
| Wasmer Tinygo Cranelift | 2047 | 604185    | 1.22       |
| Wasmer Rust Cranelift   | 2082 | 569174    | 1.15       |
| Wazero Tinygo           | 1038 | 1148140   | 2.31       |
| Wasmer Tinygo Singlepass| 1056 | 1145900   | 2.31       |
| Wasmer Rust Singlepass  | 1117 | 1081074   | 2.18       |
| EVM                     | 7    | 167516639 | 337.56     |
