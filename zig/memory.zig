const std = @import("std");
const allocator = std.heap.page_allocator;

var allocated_pointers = std.AutoHashMap([*]u8, usize).init(allocator);

pub fn concrete_Malloc(length: u64) ?[*]u8 {
    const length_usize: usize = @as(usize, @intCast(length));
    const buff = allocator.alloc(u8, length_usize) catch unreachable;
    const ptr = buff.ptr;
    allocated_pointers.put(ptr, length_usize) catch unreachable;
    return ptr;
}

pub fn concrete_Free(ptr: [*]u8) void {
    const length = allocated_pointers.get(ptr) orelse return;
    allocator.free(ptr[0..length]);
    _ = allocated_pointers.remove(ptr);
}

pub fn concrete_Prune() void {
    var it = allocated_pointers.iterator();
    while (it.next()) |entry| {
        const ptr: [*]u8 = @as([*]u8, @ptrCast(entry.key_ptr));
        const length: usize = @intFromPtr(entry.value_ptr);
        allocator.free(ptr[0..length]);
    }
    allocated_pointers.deinit();
    allocated_pointers = std.AutoHashMap([*]u8, usize).init(allocator);
}

pub fn write(data: []const u8) u64 {
    if (data.len == 0) {
        return 0;
    }
    const offset: u64 = @intFromPtr(data.ptr);
    const pointer: u64 = offset << 32 | data.len;
    return pointer;
}

pub fn read(pointer: u64) []const u8 {
    if (pointer == 0) {
        return &[_]u8{};
    }
    const offset: u64 = pointer >> 32;
    const offset_usize: usize = @as(usize, @intCast(offset));
    const length: u64 = pointer & 0xffffffff;
    const length_usize: usize = @as(usize, @intCast(length));
    const ptr = @as([*]u8, @ptrFromInt(offset_usize));
    return ptr[0..length_usize];
}

pub fn malloc(size: usize) u64 {
    const offset: u64 = @intFromPtr(concrete_Malloc(size));
    const pointer: u64 = offset << 32 | size;
    return pointer;
}

pub fn free(pointer: u64) void {
    const offset: u64 = pointer >> 32;
    const offset_usize: usize = @as(usize, @intCast(offset));
    const ptr = @as([*]u8, @ptrFromInt(offset_usize));
    concrete_Free(ptr);
}

pub fn prune() void {
    concrete_Prune();
}
