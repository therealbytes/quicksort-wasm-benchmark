const std = @import("std");
const allocator = std.heap.page_allocator;
const memory = @import("memory.zig");

pub const environment = struct {};

pub const precompile = struct {
    is_static: fn (input: []const u8) bool,
    finalise: fn (env: *const environment) anyerror!void,
    commit: fn (env: *const environment) anyerror!void,
    run: fn (env: *const environment, input: []const u8) anyerror![]const u8,
};

pub const blank_precompile = precompile{
    .is_static = blank_is_static,
    .finalise = blank_finalise,
    .commit = blank_commit,
    .run = blank_run,
};

pub fn blank_is_static(input: []const u8) bool {
    _ = input;
    return true;
}

pub fn blank_finalise(env: *const environment) anyerror!void {
    _ = env;
}

pub fn blank_commit(env: *const environment) anyerror!void {
    _ = env;
}

pub fn blank_run(env: *const environment, input: []const u8) anyerror![]const u8 {
    _ = input;
    _ = env;
    return &[_]u8{};
}

fn encode_pointer(pointer: u64) []const u8 {
    var buffer: [8]u8 = undefined;
    std.mem.writeIntBig(u64, &buffer, pointer);
    return buffer[0..];
}

const error_nil = &[_]u8{0x00};
const error_empty = &[_]u8{0x01};

pub const precompile_wasm_wrapper = struct {
    pc: *const precompile,

    pub fn is_static(comptime self: *const precompile_wasm_wrapper, pointer: u64) u64 {
        const input = memory.read(pointer);
        if (self.pc.is_static(input)) {
            return 1;
        } else {
            return 0;
        }
    }

    pub fn finalise(comptime self: *const precompile_wasm_wrapper) u64 {
        const env = &environment{};
        self.pc.finalise(env) catch {
            return memory.write(error_empty);
        };
        return memory.write(error_nil);
    }

    pub fn commit(comptime self: *const precompile_wasm_wrapper) u64 {
        const env = &environment{};
        self.pc.commit(env) catch {
            return memory.write(error_empty);
        };
        return memory.write(error_nil);
    }

    pub fn run(comptime self: *const precompile_wasm_wrapper, pointer: u64) u64 {
        const input = memory.read(pointer);
        const env = &environment{};
        const output = self.pc.run(env, input) catch {
            const err_ptr = memory.write(error_empty);
            return memory.write(encode_pointer(err_ptr));
        };
        const err_ptr = memory.write(error_empty);
        const err_ptr_enc = encode_pointer(err_ptr);
        const out_ptr = memory.write(output);
        const out_ptr_enc = encode_pointer(out_ptr);
        var pack: [16]u8 = undefined;
        pack[0..8].* = out_ptr_enc[0..8].*;
        pack[8..16].* = err_ptr_enc[0..8].*;
        return memory.write(pack[0..]);
    }

    pub fn init(comptime pc: *const precompile) precompile_wasm_wrapper {
        return precompile_wasm_wrapper{ .pc = pc };
    }
};
