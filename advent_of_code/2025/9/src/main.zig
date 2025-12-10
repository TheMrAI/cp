const std = @import("std");

const Pos = struct {
    i: i64,
    j: i64,
};

pub fn main() !void {
    var threaded: std.Io.Threaded = .init_single_threaded;
    const io = threaded.io();

    var stdin_buffer: [1024]u8 = undefined;
    var stdin: std.fs.File = std.fs.File.stdin();
    var stdin_reader: std.fs.File.Reader = stdin.reader(io, &stdin_buffer);

    const allocator = std.heap.page_allocator;
    var positions = std.ArrayList(Pos).empty;
    defer positions.deinit(allocator);

    while (stdin_reader.interface.takeDelimiterInclusive('\n')) |line| {
        const trimmed = std.mem.trimRight(u8, line, "\r\n");

        var coordinates = std.mem.splitScalar(u8, trimmed, ',');
        const i = try std.fmt.parseInt(i64, coordinates.next().?, 10);
        const j = try std.fmt.parseInt(i64, coordinates.next().?, 10);
        std.debug.assert(coordinates.next() == null);

        try positions.append(allocator, Pos{ .i = i, .j = j });
    } else |err| {
        if (err != error.EndOfStream) {
            std.debug.print("an error occurred: {any}\n", .{err});
        }
    }

    std.debug.print("The largest rectangle possible: {d}\n", .{part_one(positions)});
    // std.debug.print("The multiple of the X coordinate of the last two junction boxes: {d}\n", .{try part_two(allocator, junction_boxes, distances)});
}

// We could do some scan line based shenanigans, but that is a little more code and brute force
// is good enough for this.
fn part_one(positions: std.ArrayList(Pos)) u64 {
    var biggest_rectangle: u64 = 0;

    for (0..positions.items.len - 1) |i| {
        const corner_A = positions.items[i];
        for (i + 1..positions.items.len) |j| {
            const corner_B = positions.items[j];

            const delta_i = @abs(corner_A.i - corner_B.i) + 1;
            const delta_j = @abs(corner_A.j - corner_B.j) + 1;

            const area = delta_i * delta_j;
            if (biggest_rectangle < area) {
                biggest_rectangle = area;
            }
        }
    }

    return biggest_rectangle;
}

// fn part_two(allocator: std.mem.Allocator, junction_boxes: std.ArrayList(Coordinate), distances: std.ArrayList(Distance)) error{OutOfMemory}!u64 {
//     var dsu_parents = std.ArrayList(usize).empty;
//     defer dsu_parents.deinit(allocator);
//     try dsu_parents.appendNTimes(allocator, 0, junction_boxes.items.len);
//     for (0..junction_boxes.items.len) |i| {
//         make_set(dsu_parents, i);
//     }

//     var dsu_ranks = std.ArrayList(usize).empty;
//     defer dsu_ranks.deinit(allocator);
//     try dsu_ranks.appendNTimes(allocator, 0, junction_boxes.items.len);

//     var fully_connected = false;
//     var distance_idx: usize = 0;
//     while (!fully_connected) {
//         const distance = distances.items[distance_idx];
//         distance_idx += 1;

//         union_sets(dsu_parents, dsu_ranks, distance.from, distance.to);
//         // make sure to compress the tree
//         for (0..dsu_parents.items.len) |i| {
//             const p = find_set(dsu_parents, i);
//             _ = p;
//         }
//         // inefficient check but hey, good enough, not in the mood to make it faster
//         fully_connected = true;
//         const value = dsu_parents.items[0];
//         for (1..dsu_parents.items.len) |i| {
//             if (dsu_parents.items[i] != value) {
//                 fully_connected = false;
//                 break;
//             }
//         }
//     }

//     // distance_idx-1 because the loop iterates it before ending
//     const last_connection = distances.items[distance_idx - 1];
//     return @as(u64, @intFromFloat(junction_boxes.items[last_connection.from].x)) * @as(u64, @intFromFloat(junction_boxes.items[last_connection.to].x));
// }
