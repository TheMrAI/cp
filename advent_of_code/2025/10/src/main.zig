const std = @import("std");
const builtin = @import("builtin");

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
    std.debug.print("Tthe fewest button presses required to correctly configure the joltage level counters on all of the machines: {d}\n", .{try part_two(machines)});
}

fn part_one(machines: std.ArrayList(Machine)) error{OutOfMemory}!u32 {
    var configuration_steps: u32 = 0;

    for (machines.items) |machine| {
        configuration_steps += try fewest_steps_to_match(machine.start_state, machine.buttons.items);
    }

    return configuration_steps;
}

// This is a BFS algorithm, because it doesn't need recursion and the possible search space is
// small. For each machine the search space will be the 2^(# indicator lights). For my input
// this seems to be max 2^10 == 1024, which is rather small.
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

fn part_two(machines: std.ArrayList(Machine)) error{OutOfMemory}!u32 {
    var button_presses: u32 = 0;

    const allocator = std.heap.page_allocator;
    for (machines.items) |machine| {
        var seen_with_steps = std.HashMap([]u32, u32, SliceHashContext, std.hash_map.default_max_load_percentage).init(allocator);
        defer seen_with_steps.deinit();

        var current_joltages = std.ArrayList(u32).empty;
        defer current_joltages.deinit(allocator);
        try current_joltages.appendNTimes(allocator, 0, machine.jolts.items.len);

        try fewest_button_presses_for_joltage_dfs(allocator, &seen_with_steps, machine.jolts.items, machine.buttons.items, current_joltages, 0);

        button_presses += seen_with_steps.get(machine.jolts.items).?;
    }

    return button_presses;
}

// The maximum button press needed for joltage targets is a few hundred (~300) and there are a about maximum 10
// numbers for a given machine, which means that in the worst case we have 300^10 possible combinations
// if we just wanted to brute force this. The number is btw: 5,904,900,000,000,000,000,000,000 (can't even pronounce it
// so brute won't work)
// Then we can look it perhaps differently at the problem.
// There are these buttons witch specify how the N counters are increased, and there are M buttons.
// If we look at the problem like this, then for example:
// [.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7} means that we have
// (3) => {0, 0, 0, 1}*A1
// (1,3) => {0, 1, 0, 1}*A2
// (2) => {0, 0, 1, 0}*A3
// (2,3) => {0, 0, 1, 1}*A4
// (0,2) => {1, 0, 1, 0}*A5
// (0,1) => {1, 1, 0, 0}*A6
// And then we are looking for sum(Ax) such that it is minimal and sum(Ax) = {3,5,4,7}.
// Notice we aren't interested in any specific A value, only that for a given sum A, can we get the target.
// Which makes me wonder if we can state the problem as
// can_you_solve_it_with_exactly_this_sum_for_target
// and then binary search the space. With lower bound being 1 and upper bound being the sum of the jolts.
// This could maybe be good as log(SumA) for almost any number will be trivially small, but not sure how we
// can estimate solving the problem for a given number.
// Well, no matter how much I think it seems that solving for one given number too is extremely complicated,
// as you have to basically check all multiples for each button within between 0 and the smallest value it
// has to provide in the final joltage.
fn fewest_button_presses_for_joltage_dfs(allocator: std.mem.Allocator, seen_with_steps: *std.HashMap([]u32, u32, SliceHashContext, std.hash_map.default_max_load_percentage), joltages: []u32, operations: []u32, current_joltages: std.ArrayList(u32), current_step: u32) error{OutOfMemory}!void {
    for (0..joltages.len) |i| {
        if (current_joltages.items[i] > joltages[i]) {
            return;
        }
    }
    if (seen_with_steps.contains(current_joltages.items)) {
        const seen_at = seen_with_steps.get(current_joltages.items).?;
        if (seen_at <= current_step) {
            return;
        }
    }
    try seen_with_steps.put(current_joltages.items, current_step);

    for (operations) |operation| {
        // std.debug.print("ITER\n", .{});
        var next_joltages = try current_joltages.clone(allocator);

        var increment_indicator = operation;
        for (0..joltages.len) |i| {
            if (increment_indicator & 1 == 1) {
                next_joltages.items[joltages.len - 1 - i] += 1;
            }
            increment_indicator >>= 1;
        }

        // std.debug.print("Dumping keys:\n", .{});
        // var iter = seen.keyIterator();
        // while (iter.next()) |key| {
        //     std.debug.print("{any}\n", .{key.*});
        // }
        // std.debug.print("\nCandidate: {any}\n", .{next_press_counts.items});

        try fewest_button_presses_for_joltage_dfs(allocator, seen_with_steps, joltages, operations, next_joltages, current_step + 1);
        // std.debug.print("put\n", .{});
    }
}

const SliceHashContext = struct {
    pub fn hash(_: SliceHashContext, counters: []u32) u64 {
        var h = std.hash.Wyhash.init(0);
        for (counters) |counter| {
            var as_bytes: [4]u8 = undefined;
            std.mem.writeInt(u32, &as_bytes, counter, builtin.cpu.arch.endian());
            h.update(&as_bytes);
        }
        const hashs = h.final();
        // std.debug.print("{any} {}\n", .{ counters, hashs });

        return hashs;
    }

    pub fn eql(_: SliceHashContext, lhs: []u32, rhs: []u32) bool {
        return std.mem.eql(u32, lhs, rhs);
    }
};
