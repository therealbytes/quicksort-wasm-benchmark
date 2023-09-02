const sort = @import("sort.zig");

export fn run() u64 {
    var b = sort.quick_sort_benchmark.init();
    const checksum: u64 = b.benchmark();
    return checksum;
}
