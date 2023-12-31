const std = @import("std");

pub fn build(b: *std.build.Builder) void {
    const target = .{ .cpu_arch = .wasm32, .os_tag = .freestanding };
    const optimize = b.standardOptimizeOption(.{});

    const lib = b.addSharedLibrary(.{
        .name = "zig",
        .root_source_file = .{ .path = "./main.zig" },
        .target = target,
        .optimize = optimize,
    });
    lib.rdynamic = true;
    b.installArtifact(lib);
}
