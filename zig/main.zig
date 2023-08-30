const bindings = @import("bindings.zig");
const memory = @import("memory.zig");
const sort = @import("sort.zig");

extern "env" fn concrete_Environment(pointer: u64) u64;

const wrapped_sort_precompile = bindings.precompile_wasm_wrapper.init(&sort_precompile);

const sort_precompile = bindings.precompile{
    .is_static = bindings.blank_is_static,
    .finalise = bindings.blank_finalise,
    .commit = bindings.blank_commit,
    .run = sort_run,
};

fn sort_run(_: *const bindings.environment, _: []const u8) anyerror![]const u8 {
    var b = sort.quick_sort_benchmark.init(7);
    _ = b.benchmark();
    return &[_]u8{};
}

export fn concrete_Malloc(length: usize) ?[*]u8 {
    return memory.concrete_Malloc(length);
}

export fn concrete_Free(ptr: [*]u8) void {
    memory.concrete_Free(ptr);
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
