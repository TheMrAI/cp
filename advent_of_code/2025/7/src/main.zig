const std = @import("std");

pub fn main() !void {
    var threaded: std.Io.Threaded = .init_single_threaded;
    const io = threaded.io();

    var stdin_buffer: [1024]u8 = undefined;
    var stdin: std.fs.File = std.fs.File.stdin();
    var stdin_reader: std.fs.File.Reader = stdin.reader(io, &stdin_buffer);

    const allocator = std.heap.page_allocator;
    var tachyon_manifold = std.ArrayList(std.ArrayList(u8)).empty;
    defer tachyon_manifold.deinit(allocator);

    while (stdin_reader.interface.takeDelimiterInclusive('\n')) |line| {
        const trimmed = std.mem.trimRight(u8, line, "\r\n");

        var entry = std.ArrayList(u8).empty;
        try entry.appendSlice(allocator, trimmed);
        try tachyon_manifold.append(allocator, entry);
    } else |err| {
        if (err != error.EndOfStream) {
            std.debug.print("an error occurred: {any}\n", .{err});
        }
    }

    var tachyon_manifold_clone = std.ArrayList(std.ArrayList(u8)).empty;
    defer tachyon_manifold_clone.deinit(allocator);
    // deep clone tachyon_manifold, can't believe I have to write this out by hand
    for (0..tachyon_manifold.items.len) |i| {
        const line_clone = try tachyon_manifold.items[i].clone(allocator);
        try tachyon_manifold_clone.append(allocator, line_clone);
    }

    std.debug.print("The number of ray splits is: {d}\n", .{part_one(tachyon_manifold_clone)});
    std.debug.print("The number of different timelines a single tachyon particle would take: {d}\n", .{try part_two(allocator, tachyon_manifold)});
}

fn part_one(tachyon_manifold: std.ArrayList(std.ArrayList(u8))) u64 {
    var split_count: u64 = 0;

    for (1..tachyon_manifold.items.len) |i| {
        for (0..tachyon_manifold.items[i].items.len) |j| {
            const above_state = tachyon_manifold.items[i - 1].items[j];
            if (above_state == 'S' or above_state == '|') {
                switch (tachyon_manifold.items[i].items[j]) {
                    '.' => tachyon_manifold.items[i].items[j] = '|',
                    '^' => {
                        split_count += 1;
                        if (j > 0 and tachyon_manifold.items[i].items[j - 1] != '^') {
                            tachyon_manifold.items[i].items[j - 1] = '|';
                        }
                        if (j < tachyon_manifold.items[i].items.len - 1 and tachyon_manifold.items[i].items[j + 1] != '^') {
                            tachyon_manifold.items[i].items[j + 1] = '|';
                        }
                    },
                    '|' => continue,
                    else => unreachable,
                }
            }
        }
    }

    return split_count;
}

// Part two example resolution:
// .......S.......
// .......1.......
// ......1^1......
// ......1.1......
// .....1^2^1.....
// .....1.2.1.....
// ....1^3^3^1....
// ....1.3.3.1....
// ...1^4^331^1...
// ...1.4.331.1...
// ..1^5^434^2^1..
// ..1.5.434.2.1..
// .1^154^74.21^1.
// .1.154.74.21.1.
// 1^2^a^b^b^211^1
// 1.2.a.b.b.211.1

fn part_two(allocator: std.mem.Allocator, tachyon_manifold: std.ArrayList(std.ArrayList(u8))) error{OutOfMemory}!u64 {
    var prev_timeline_counts = std.ArrayList(u64).empty;
    try prev_timeline_counts.appendNTimes(allocator, 0, tachyon_manifold.items[0].items.len);
    defer prev_timeline_counts.deinit(allocator);

    var timeline_counts = std.ArrayList(u64).empty;
    try timeline_counts.appendNTimes(allocator, 0, tachyon_manifold.items[0].items.len);
    defer timeline_counts.deinit(allocator);

    for (0..tachyon_manifold.items.len) |j| {
        if (tachyon_manifold.items[0].items[j] == 'S') {
            timeline_counts.items[j] = 1;
            break;
        }
    }

    for (1..tachyon_manifold.items.len) |i| {
        std.mem.swap(std.ArrayList(u64), &prev_timeline_counts, &timeline_counts);
        for (0..timeline_counts.items.len) |j| {
            timeline_counts.items[j] = 0;
        }

        for (0..tachyon_manifold.items[i].items.len) |j| {
            const above_state = tachyon_manifold.items[i - 1].items[j];
            if (above_state == 'S' or above_state == '|') {
                switch (tachyon_manifold.items[i].items[j]) {
                    '.' => {
                        tachyon_manifold.items[i].items[j] = '|';
                        timeline_counts.items[j] += prev_timeline_counts.items[j];
                    },
                    '^' => {
                        if (j > 0 and tachyon_manifold.items[i].items[j - 1] != '^') {
                            tachyon_manifold.items[i].items[j - 1] = '|';
                            tachyon_manifold.items[i].items[j] = '|';
                            timeline_counts.items[j - 1] += prev_timeline_counts.items[j];
                        }
                        if (j < tachyon_manifold.items[i].items.len - 1 and tachyon_manifold.items[i].items[j + 1] != '^') {
                            tachyon_manifold.items[i].items[j + 1] = '|';
                            timeline_counts.items[j + 1] += prev_timeline_counts.items[j];
                        }
                    },
                    '|' => timeline_counts.items[j] += prev_timeline_counts.items[j],
                    else => unreachable,
                }
            }
        }
    }

    var total_timelines: u64 = 0;
    for (timeline_counts.items) |count| {
        total_timelines += count;
    }
    return total_timelines;
}
