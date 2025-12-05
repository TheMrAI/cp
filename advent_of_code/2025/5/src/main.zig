const std = @import("std");

const Range = struct {
    from: u64,
    to: u64,
};

pub fn main() !void {
    var threaded: std.Io.Threaded = .init_single_threaded;
    const io = threaded.io();

    var stdin_buffer: [1024]u8 = undefined;
    var stdin: std.fs.File = std.fs.File.stdin();
    var stdin_reader: std.fs.File.Reader = stdin.reader(io, &stdin_buffer);

    const allocator = std.heap.page_allocator;
    var ranges = std.ArrayList(Range).empty;
    defer ranges.deinit(allocator);

    var ids = std.ArrayList(u64).empty;
    defer ids.deinit(allocator);

    var read_ranges = true;
    while (stdin_reader.interface.takeDelimiterInclusive('\n')) |line| {
        const trimmed = std.mem.trimRight(u8, line, "\r\n");
        if (trimmed.len == 0) {
            read_ranges = false;
            continue;
        }

        if (read_ranges) {
            const hyphen = std.mem.findScalarPos(u8, trimmed, 0, '-').?;

            const from = try std.fmt.parseInt(u64, trimmed[0..hyphen], 10);
            const to = try std.fmt.parseInt(u64, trimmed[hyphen + 1 ..], 10);
            const range = Range{ .from = from, .to = to };

            try merge_into_range(allocator, &ranges, range);
        } else {
            const id = try std.fmt.parseInt(u64, trimmed, 10);
            try ids.append(allocator, id);
        }
    } else |err| {
        if (err != error.EndOfStream) {
            std.debug.print("an error occurred: {any}\n", .{err});
        }
    }

    std.debug.print("Total fresh ID count: {d}\n", .{part_one(ranges.items, ids.items)});
    std.debug.print("Total valid ID count: {d}\n", .{part_two(ranges.items)});
}

fn merge_into_range(allocator: std.mem.Allocator, ranges: *std.ArrayList(Range), range: Range) error{OutOfMemory}!void {
    var working_range = range;

    const S = struct {
        fn compareRangeFrom(context: Range, item: Range) std.math.Order {
            return std.math.order(context.from, item.from);
        }

        fn compareRangeToWithFrom(context: Range, item: Range) std.math.Order {
            return std.math.order(context.to, item.from);
        }
    };

    var from_idx = std.sort.lowerBound(Range, ranges.items, working_range, S.compareRangeFrom);
    const to_idx = std.sort.upperBound(Range, ranges.items, working_range, S.compareRangeToWithFrom);

    if (from_idx > 0) {
        if (ranges.items[from_idx - 1].to >= working_range.from) {
            working_range.from = ranges.items[from_idx - 1].from;
            from_idx -= 1;
        }
    }
    if (to_idx > 0) {
        if (ranges.items[to_idx - 1].to >= working_range.to) {
            working_range.to = ranges.items[to_idx - 1].to;
        }
    }
    const replace = [1]Range{working_range};
    try ranges.replaceRange(allocator, from_idx, to_idx - from_idx, &replace);
}

fn part_one_brute(ranges: []const Range, ids: []const u64) u64 {
    var fresh_count: u64 = 0;

    for (ids) |id| {
        for (ranges) |range| {
            if (range.from <= id and id <= range.to) {
                fresh_count += 1;
                break;
            }
        }
    }

    return fresh_count;
}

fn part_one(ranges: []const Range, ids: []const u64) u64 {
    var fresh_count: u64 = 0;

    const S = struct {
        fn compareWithRangeFrom(context: u64, item: Range) std.math.Order {
            if (context < item.from) {
                return std.math.Order.lt;
            } else if (context > item.to) {
                return std.math.Order.gt;
            } else {
                return std.math.Order.eq;
            }
        }
    };

    for (ids) |id| {
        if (std.sort.binarySearch(Range, ranges, id, S.compareWithRangeFrom) != null) {
            fresh_count += 1;
        }
    }

    return fresh_count;
}

fn part_two(ranges: []const Range) u64 {
    var valid_id_count: u64 = 0;
    for (ranges) |range| {
        valid_id_count += 1 + range.to - range.from;
    }
    return valid_id_count;
}
