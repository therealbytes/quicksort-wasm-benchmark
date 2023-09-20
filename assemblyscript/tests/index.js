import assert from "assert";
import { run } from "../build/debug.js";
assert.strictEqual(run(7n, 1000n, 100n), 49760n);
console.log("ok");
