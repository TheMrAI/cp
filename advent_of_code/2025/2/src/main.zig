const std = @import("std");

const Range = struct { from: u64, to: u64 };

pub fn main() !void {
    var threaded: std.Io.Threaded = .init_single_threaded;
    const io = threaded.io();

    var stdin_buffer: [1024]u8 = undefined;
    var stdin: std.fs.File = std.fs.File.stdin();
    var stdin_reader: std.fs.File.Reader = stdin.reader(io, &stdin_buffer);

    const allocator = std.heap.page_allocator;
    var ranges = std.ArrayList(Range).empty;
    defer ranges.deinit(allocator);

    while (stdin_reader.interface.takeDelimiterInclusive('\n')) |line| {
        const trimmed = std.mem.trimRight(u8, line, "\r\n");

        var entries = std.mem.splitScalar(u8, trimmed, ',');
        while (entries.next()) |entry| {
            const hyphen = std.mem.findScalarPos(u8, entry, 0, '-').?;

            const from = try std.fmt.parseInt(u64, entry[0..hyphen], 10);
            const to = try std.fmt.parseInt(u64, entry[hyphen + 1 ..], 10);
            try ranges.append(allocator, Range{ .from = from, .to = to });
        }
    } else |err| {
        if (err != error.EndOfStream) {
            std.debug.print("an error occurred: {any}\n", .{err});
        }
    }

    std.debug.print("Adding up the invalid IDs brute gives sum: {d}\n", .{try part_one_brute_force(ranges.items)});
    std.debug.print("Adding up the invalid IDs brute for any repeating block gives sum: {d}\n", .{try part_two_brute_force(ranges.items)});
}

fn digit_count(value: u64) u8 {
    var count: u8 = 1;
    var v = value;
    while (v != 0) {
        v /= 10;
        if (v != 0) {
            count += 1;
        }
    }

    return count;
}

fn part_one_brute_force(ranges: []Range) !u64 {
    var invalid_id_sum: u64 = 0;

    for (ranges) |range| {
        var checking = range.from;
        while (checking <= range.to) {
            var buf: [20]u8 = undefined;
            const as_string = try std.fmt.bufPrint(&buf, "{d}", .{checking});

            if (repeating_blocks(as_string, 2)) {
                invalid_id_sum += checking;
            }

            checking += 1;
        }
    }

    return invalid_id_sum;
}

fn part_two_brute_force(ranges: []Range) !u64 {
    var invalid_id_sum: u64 = 0;

    for (ranges) |range| {
        var checking = range.from;
        while (checking <= range.to) {
            var buf: [20]u8 = undefined;
            const as_string = try std.fmt.bufPrint(&buf, "{d}", .{checking});

            if (repeating_block_range_found(as_string)) {
                invalid_id_sum += checking;
            }

            checking += 1;
        }
    }

    return invalid_id_sum;
}

fn repeating_block_range_found(buf: []u8) bool {
    for (2..buf.len + 1) |block_count| {
        if (repeating_blocks(buf, block_count)) {
            return true;
        }
    }
    return false;
}

fn repeating_blocks(buf: []u8, block_count: usize) bool {
    if (buf.len % block_count != 0) {
        return false;
    }

    const block_size = buf.len / block_count;
    var front: usize = 0;
    const next: usize = front + block_size;

    while (front < next) {
        // check all block starts for the character match
        for (1..block_count) |block| {
            if (buf[front] != buf[front + (block * block_size)]) {
                return false;
            }
        }
        front += 1;
    }
    return true;
}
