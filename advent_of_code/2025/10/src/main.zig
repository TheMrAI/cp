const std = @import("std");

const Machine = struct {
    start_state: u32,
    buttons: std.ArrayList(u32),
    jolts: std.ArrayList(u32),
};

pub fn main() !void {
    var threaded: std.Io.Threaded = .init_single_threaded;
    const io = threaded.io();

    var stdin_buffer: [1024]u8 = undefined;
    var stdin: std.fs.File = std.fs.File.stdin();
    var stdin_reader: std.fs.File.Reader = stdin.reader(io, &stdin_buffer);

    const allocator = std.heap.page_allocator;
    var machines = std.ArrayList(Machine).empty;
    defer machines.deinit(allocator);

    while (stdin_reader.interface.takeDelimiterInclusive('\n')) |line| {
        const trimmed = std.mem.trimRight(u8, line, "\r\n");

        var sections = std.mem.splitScalar(u8, trimmed, ' ');

        // no input verification, it is expected to be correct
        var start_state: u32 = 0;
        var buttons = std.ArrayList(u32).empty;
        var jolts = std.ArrayList(u32).empty;

        const start_part = sections.next().?;
        for (1..start_part.len - 1) |i| {
            start_state <<= 1;
            if (start_part[i] == '#') {
                start_state |= 1;
            }
        }
        const mask_width_minus_one = start_part.len - 3;
        std.debug.print("{}\n", .{mask_width_minus_one});

        while (sections.next()) |section| {
            var parts = std.mem.splitScalar(u8, section[1 .. section.len - 1], ',');

            if (section[0] == '(') {
                var button: u32 = 0;
                while (parts.next()) |button_text| {
                    const value = try std.fmt.parseInt(u32, button_text, 10);
                    const offset: u5 = @intCast(mask_width_minus_one - value);
                    const one: u32 = 1;
                    button |= one << offset;
                }
                try buttons.append(allocator, button);
            } else if (section[0] == '{') {
                while (parts.next()) |jolt_text| {
                    const jolt = try std.fmt.parseInt(u32, jolt_text, 10);
                    try jolts.append(allocator, jolt);
                }
            } else {
                unreachable;
            }
        }
        std.debug.assert(sections.next() == null);

        try machines.append(allocator, Machine{ .start_state = start_state, .buttons = buttons, .jolts = jolts });
    } else |err| {
        if (err != error.EndOfStream) {
            std.debug.print("an error occurred: {any}\n", .{err});
        }
    }

    for (machines.items) |machine| {
        std.debug.print("\nMachine: {b}\n", .{machine.start_state});
        std.debug.print("Buttons:", .{});
        for (machine.buttons.items) |button| {
            std.debug.print("{b},", .{button});
        }
        std.debug.print("\nJolts: {any}\n", .{machine.jolts});
    }

    std.debug.print("The fewest button presses required to correctly configure the indicator lights on all of the machines: {d}\n", .{try part_one(machines)});
    // std.debug.print("The multiple of the X coordinate of the last two junction boxes: {d}\n", .{try part_two(allocator, junction_boxes, distances)});
}

fn part_one(machines: std.ArrayList(Machine)) error{OutOfMemory}!u32 {
    var configuration_steps: u32 = 0;

    for (machines.items) |machine| {
        configuration_steps += try fewest_steps_to_match(machine.start_state, machine.buttons.items);
    }

    return configuration_steps;
}

fn fewest_steps_to_match(desired_state: u32, operations: []u32) error{OutOfMemory}!u32 {
    var at_step: u32 = 1;

    const allocator = std.heap.page_allocator;

    var seen = std.AutoHashMap(u32, u1).init(allocator);
    defer seen.deinit();

    var next_states = std.ArrayList(u32).empty;
    defer next_states.deinit(allocator);

    var current_states = std.ArrayList(u32).empty;
    defer current_states.deinit(allocator);
    try current_states.append(allocator, 0);

    while (current_states.items.len != 0) {
        next_states.clearRetainingCapacity();

        for (current_states.items) |state| {
            for (operations) |operation| {
                const next_state = state ^ operation;
                if (next_state == desired_state) {
                    return at_step;
                }
                if (seen.contains(next_state)) {
                    continue;
                }
                try seen.put(next_state, 0);
                try next_states.append(allocator, next_state);
            }
        }
        at_step += 1;

        std.mem.swap(std.ArrayList(u32), &current_states, &next_states);
    }

    return std.math.maxInt(u32);
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
