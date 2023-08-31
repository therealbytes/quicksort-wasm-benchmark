const sort = @import("sort.zig");

export fn run(seed: u64) u64 {
    var b = sort.quick_sort_benchmark.init(seed);
    const checksum: u64 = b.benchmark();
    return checksum;
}
