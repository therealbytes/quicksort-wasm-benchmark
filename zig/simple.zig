const sort = @import("sort.zig");

export fn run(seed: i32, arr_len: i32, iter: i32) i32 {
    var qs = sort.quick_sort_benchmark.init(@as(usize, @intCast(seed)));
    const checksum_result = qs.run(@as(usize, @intCast(arr_len)), @as(usize, @intCast(iter))) catch {
        return -1;
    };
    return @as(i32, @intCast(checksum_result));
}
