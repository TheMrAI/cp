const std = @import("std");

const Rotation = struct {
    dir: u8,
    value: u16,
};

pub fn main() !void {
    var threaded: std.Io.Threaded = .init_single_threaded;
    const io = threaded.io();

    var stdin_buffer: [1024]u8 = undefined;
    var stdin: std.fs.File = std.fs.File.stdin();
    var stdin_reader: std.fs.File.Reader = stdin.reader(io, &stdin_buffer);

    const allocator = std.heap.page_allocator;
    var rotations = std.ArrayList(Rotation).empty;
    defer rotations.deinit(allocator);

    while (stdin_reader.interface.takeDelimiterInclusive('\n')) |line| {
        const trimmed = std.mem.trimRight(u8, line, "\r\n");
        const dir = trimmed[0];
        const value = try std.fmt.parseInt(u16, trimmed[1..], 10);

        try rotations.append(allocator, Rotation{ .dir = dir, .value = value });
    } else |err| {
        if (err != error.EndOfStream) {
            std.debug.print("an error occurred: {any} {any}\n", .{ err, rotations });
        }
    }

    std.debug.print("Door password: {d}\n", .{part_one(50, rotations.items)});
    std.debug.print("Door password with method 0x434C49434B: {d}\n", .{part_two(50, rotations.items)});
}

fn part_one(start_pos: u8, rotations: []Rotation) u16 {
    var zero_count: u16 = 0;

    var pos: i32 = start_pos;
    for (rotations) |rotation| {
        const value = @as(i32, rotation.value);
        if (rotation.dir == 'L') {
            pos -= value;
        } else {
            pos += value;
        }

        pos = @mod(pos, 100);

        if (pos == 0) {
            zero_count += 1;
        }
    }

    return zero_count;
}

fn part_two(start_pos: u8, rotations: []Rotation) u16 {
    var zero_count: u16 = 0;

    var pos: i32 = start_pos;
    for (rotations) |rotation| {
        var value = @as(i32, rotation.value);

        if (rotation.dir == 'L') {
            value = -value;
        }

        zero_count += @truncate(@divTrunc(@abs(value), 100));
        const remainder = @rem(value, 100);
        const stepped = pos + remainder;
        if ((pos != 0 and stepped <= 0) or stepped >= 100) {
            zero_count += 1;
        }

        pos = @mod(stepped, 100);
    }

    return zero_count;
}

// too low: 5611
