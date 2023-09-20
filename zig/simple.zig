const sort = @import("sort.zig");

export fn run(seed: i64, arr_len: i64, iter: i64) i64 {
    var qs = sort.quick_sort_benchmark.init(@as(usize, @intCast(seed)));
    const checksum_result = qs.run(@as(usize, @intCast(arr_len)), @as(usize, @intCast(iter))) catch {
        return -1;
    };
    return @as(i64, @intCast(checksum_result));
}
