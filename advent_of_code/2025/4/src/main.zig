const std = @import("std");

const Field = enum {
    roll,
    free,
};

pub fn main() !void {
    var threaded: std.Io.Threaded = .init_single_threaded;
    const io = threaded.io();

    var stdin_buffer: [1024]u8 = undefined;
    var stdin: std.fs.File = std.fs.File.stdin();
    var stdin_reader: std.fs.File.Reader = stdin.reader(io, &stdin_buffer);

    const allocator = std.heap.page_allocator;
    var warehouse = std.ArrayList(std.ArrayList(Field)).empty;
    defer warehouse.deinit(allocator);

    while (stdin_reader.interface.takeDelimiterInclusive('\n')) |line| {
        const trimmed = std.mem.trimRight(u8, line, "\r\n");
        const fields = try to_field_array(allocator, trimmed);

        try warehouse.append(allocator, fields);
    } else |err| {
        if (err != error.EndOfStream) {
            std.debug.print("an error occurred: {any}\n", .{err});
        }
    }

    std.debug.print("Accessible roll count: {d}\n", .{part_one(warehouse)});
    std.debug.print("The number of removable paper roll count: {d}\n", .{part_two(warehouse)});
}

fn to_field_array(allocator: std.mem.Allocator, ascii_chars: []const u8) error{OutOfMemory}!std.ArrayList(Field) {
    var fields = std.ArrayList(Field).empty;

    for (ascii_chars) |ascii| {
        const field = switch (ascii) {
            '.' => Field.free,
            '@' => Field.roll,
            else => unreachable,
        };
        try fields.append(allocator, field);
    }

    return fields;
}

fn part_one(warehouse: std.ArrayList(std.ArrayList(Field))) u64 {
    var accessible_count: u64 = 0;

    for (0..warehouse.items.len) |i| {
        for (0..warehouse.items[0].items.len) |j| {
            if (warehouse.items[i].items[j] == Field.free) {
                continue;
            }

            if (accessible(warehouse, i, j)) {
                accessible_count += 1;
            }
        }
    }

    return accessible_count;
}

fn part_two(warehouse: std.ArrayList(std.ArrayList(Field))) u64 {
    var removed_paper_count: u64 = 0;

    var checking = warehouse;
    var after_remove = warehouse;

    while (true) {
        var removed = false;
        for (0..checking.items.len) |i| {
            for (0..checking.items[0].items.len) |j| {
                if (checking.items[i].items[j] == Field.free) {
                    continue;
                }

                if (accessible(checking, i, j)) {
                    after_remove.items[i].items[j] = Field.free;
                    removed_paper_count += 1;
                    removed = true;
                }
            }
        }
        checking = after_remove;
        if (!removed) {
            break;
        }
    }

    return removed_paper_count;
}

fn accessible(warehouse: std.ArrayList(std.ArrayList(Field)), i: usize, j: usize) bool {
    const di = [8]i8{ 1, 0, -1, -1, -1, 0, 1, 1 };
    const dj = [8]i8{ -1, -1, -1, 0, 1, 1, 1, 0 };

    var free_count: u8 = 0;
    for (0..8) |delta_idx| {
        const ti = @as(i64, @intCast(i)) + di[delta_idx];
        const tj = @as(i64, @intCast(j)) + dj[delta_idx];

        if ((ti < 0 or ti >= warehouse.items.len) or (tj < 0 or tj >= warehouse.items[0].items.len)) {
            free_count += 1;
        } else {
            if (warehouse.items[@as(usize, @intCast(ti))].items[@as(usize, @intCast(tj))] == Field.free) {
                free_count += 1;
            }
        }
    }
    if (free_count >= 5) {
        return true;
    }
    return false;
}
