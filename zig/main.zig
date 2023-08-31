const std = @import("std");
const bindings = @import("bindings.zig");
const memory = @import("memory.zig");
const mem_pointer = memory.mem_pointer;
const sort = @import("sort.zig");
const expect = std.testing.expect;

extern "env" fn concrete_Environment(pointer: u64) u64;

const wrapped_sort_precompile = bindings.precompile_wasm_wrapper.init(&sort_precompile);

const sort_precompile = bindings.precompile{
    .is_static = bindings.blank_is_static,
    .finalise = bindings.blank_finalise,
    .commit = bindings.blank_commit,
    .run = sort_run,
};

fn sort_run(_: *const bindings.environment, _: []const u8) anyerror![]const u8 {
    const checksum: u64 = 42;
    // var b = sort.quick_sort_benchmark.init(7);
    // const checksum: u64 = b.benchmark();
    var buf: [8]u8 = undefined;
    std.mem.writeIntBig(u64, &buf, checksum);
    return buf[0..];
}

export fn concrete_Malloc(length: u64) u64 {
    return memory.concrete_Malloc(length);
}

export fn concrete_Free(_pointer: u64) void {
    memory.concrete_Free(_pointer);
}

export fn concrete_Prune() void {
    memory.concrete_Prune();
}

export fn concrete_IsStatic(pointer: u64) u64 {
    return wrapped_sort_precompile.is_static(pointer);
}

export fn concrete_Finalise() u64 {
    return wrapped_sort_precompile.finalise();
}

export fn concrete_Commit() u64 {
    return wrapped_sort_precompile.commit();
}

export fn concrete_Run(pointer: u64) u64 {
    return wrapped_sort_precompile.run(pointer);
}
