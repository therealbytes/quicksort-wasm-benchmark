import assert from "assert";
import { run } from "../build/debug.js";
assert.strictEqual(run(), 107829970005n);
console.log("ok");
