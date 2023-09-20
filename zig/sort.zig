const std = @import("std");
const allocator = std.heap.page_allocator;
const expect = std.testing.expect;

pub const quick_sort_benchmark = struct {
    seed: u64,

    pub fn init(seed: usize) quick_sort_benchmark {
        return quick_sort_benchmark{ .seed = @as(u64, seed) };
    }

    pub fn random(self: *quick_sort_benchmark) usize {
        self.seed = (1103515245 * self.seed + 12345) % (1 << 31);
        return @as(usize, @truncate(self.seed));
    }

    pub fn randomize_array(self: *quick_sort_benchmark, arr: []usize) void {
        for (arr) |*value| {
            value.* = self.random() % 1000;
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

    pub fn run(self: *quick_sort_benchmark, arr_len: usize, iter: usize) !usize {
        var checksum: usize = 0;
        const arr = try allocator.alloc(usize, arr_len);
        defer allocator.free(arr);

        var i: usize = 0;
        while (i < iter) : (i += 1) {
            self.randomize_array(arr);
            self.quick_sort(arr, 0, arr.len - 1);
            checksum += arr[arr_len / 2];
        }
        return checksum;
    }
};

test "check checksum" {
    var qs = quick_sort_benchmark.init(7);
    const result = try qs.run(1000, 100);
    try expect(result == 49760); // You may need to adjust this expectation
}
