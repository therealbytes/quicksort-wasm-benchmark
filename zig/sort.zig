const std = @import("std");
const allocator = std.heap.page_allocator;
const expect = std.testing.expect;

pub const quick_sort_benchmark = struct {
    seed: u64,

    pub fn random(self: *quick_sort_benchmark) usize {
        self.seed = (1103515245 * self.seed + 12345) % (1 << 31);
        return @as(usize, @intCast(self.seed));
    }

    pub fn randomize_array(self: *quick_sort_benchmark, arr: []usize) void {
        for (arr) |*value| {
            value.* = self.random();
        }
    }

    pub fn quick_sort(self: *quick_sort_benchmark, arr: []usize, left: usize, right: usize) void {
        if (left >= right) return;
        var i: usize = left;
        var j: usize = right;
        const pivot = arr[left + (right - left) / 2];
        while (i <= j) {
            while (arr[i] < pivot) : (i += 1) {}
            while (arr[j] > pivot) : (j -= 1) {}
            if (i <= j) {
                std.mem.swap(usize, &arr[i], &arr[j]);
                i += 1;
                if (j > 0) j -= 1;
            }
        }
        if (left < j) self.quick_sort(arr, left, j);
        if (i < right) self.quick_sort(arr, i, right);
    }

    pub fn benchmark(self: *quick_sort_benchmark) u64 {
        var checksum: usize = 0;
        var arr: [1000]usize = undefined;
        var i: usize = 0;
        while (i < 1) : (i += 1) {
            self.randomize_array(&arr);
            self.quick_sort(&arr, 0, arr.len - 1);
            checksum += arr[100];
        }
        return checksum;
    }

    pub fn init(seed: u64) quick_sort_benchmark {
        return quick_sort_benchmark{ .seed = seed };
    }
};

test "check checksum" {
    var qs = quick_sort_benchmark.init(7);
    const result = qs.benchmark();
    try expect(result == 21880255009);
}
