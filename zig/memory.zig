const std = @import("std");
const allocator = std.heap.page_allocator;

pub const mem_pointer = struct {
    offset: usize,
    length: usize,

    pub fn is_null(self: mem_pointer) bool {
        return self.offset == 0 and self.length == 0;
    }

    pub fn to_u64(self: mem_pointer) u64 {
        const offset: u64 = @intCast(self.offset);
        const length: u64 = @intCast(self.length);
        return offset << 32 | length;
    }

    pub fn from_u64(_pointer: u64) mem_pointer {
        const offset: usize = @as(usize, @intCast(_pointer >> 32));
        const length: usize = @as(usize, @intCast(_pointer & 0xffffffff));
        return mem_pointer{
            .offset = offset,
            .length = length,
        };
    }

    pub fn encode(self: mem_pointer) []const u8 {
        const _pointer = self.to_u64();
        var buf: [8]u8 = undefined;
        std.mem.writeIntBig(u64, &buf, _pointer);
        return buf[0..];
    }

    pub fn decode(data: []const u8) mem_pointer {
        var buf: [8]u8 = undefined;
        std.mem.copy(u8, &buf, data);
        const _pointer: u64 = std.mem.readIntBig(u64, &buf);
        return mem_pointer.from_u64(_pointer);
    }

    pub fn init(offset: usize, length: usize) mem_pointer {
        return mem_pointer{
            .offset = offset,
            .length = length,
        };
    }
};

fn malloc(length: usize) mem_pointer {
    if (length == 0) {
        return mem_pointer.init(0, 0);
    }
    const buff = allocator.alloc(u8, length) catch unreachable;
    const ptr = buff.ptr;
    const offset: usize = @intFromPtr(ptr);
    const pointer = mem_pointer.init(offset, length);
    return pointer;
}

fn free(pointer: mem_pointer) void {
    if (pointer.is_null()) {
        return;
    }
    const ptr = @as([*]u8, @ptrFromInt(pointer.offset));
    const length = pointer.length;
    allocator.free(ptr[0..length]);
}

fn prune() void {
    unreachable;
}

pub fn concrete_Malloc(length: u64) u64 {
    const length_usize: usize = @as(usize, @intCast(length));
    const pointer = malloc(length_usize);
    return pointer.to_u64();
}

pub fn concrete_Free(_pointer: u64) void {
    const pointer = mem_pointer.from_u64(_pointer);
    free(pointer);
}

pub fn concrete_Prune() void {
    prune();
}

pub fn write(data: []const u8) mem_pointer {
    if (data.len == 0) {
        return mem_pointer.init(0, 0);
    }
    const cp = allocator.alloc(u8, data.len) catch unreachable;
    std.mem.copy(u8, cp, data);
    const offset: usize = @intFromPtr(cp.ptr);
    const pointer = mem_pointer.init(offset, data.len);
    return pointer;
}

pub fn read(pointer: mem_pointer) []const u8 {
    if (pointer.is_null()) {
        return &[_]u8{};
    }
    const offset = pointer.offset;
    const length = pointer.length;
    const ptr = @as([*]u8, @ptrFromInt(offset));
    return ptr[0..length];
}

test "memory" {
    const data = [_]u8{ 'h', 'e', 'l', 'l', 'o' };
    const pointer = write(&data);
    const readData = read(pointer);
    try std.testing.expect(std.mem.eql(u8, &data, readData));
    free(pointer);
}
