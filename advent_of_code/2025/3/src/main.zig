const std = @import("std");

pub fn main() !void {
    var threaded: std.Io.Threaded = .init_single_threaded;
    const io = threaded.io();

    var stdin_buffer: [1024]u8 = undefined;
    var stdin: std.fs.File = std.fs.File.stdin();
    var stdin_reader: std.fs.File.Reader = stdin.reader(io, &stdin_buffer);

    const allocator = std.heap.page_allocator;
    var battery_banks = std.ArrayList(std.ArrayList(u8)).empty;
    defer battery_banks.deinit(allocator);

    while (stdin_reader.interface.takeDelimiterInclusive('\n')) |line| {
        const trimmed = std.mem.trimRight(u8, line, "\r\n");
        const numbers = try to_number_array(allocator, trimmed);
        try battery_banks.append(allocator, numbers);
    } else |err| {
        if (err != error.EndOfStream) {
            std.debug.print("an error occurred: {any}\n", .{err});
        }
    }

    std.debug.print("The total output joltage is: {d}\n", .{part_one(battery_banks)});
    std.debug.print("The total output joltage while sending 12 batteries is: {d}\n", .{part_two(battery_banks)});
}

fn to_number_array(allocator: std.mem.Allocator, ascii_chars: []const u8) error{OutOfMemory}!std.ArrayList(u8) {
    var numbers = std.ArrayList(u8).empty;

    for (ascii_chars) |ascii| {
        try numbers.append(allocator, ascii - '0');
    }

    return numbers;
}

fn part_one(battery_banks: std.ArrayList(std.ArrayList(u8))) u64 {
    var total_output_joltage: u64 = 0;

    for (battery_banks.items) |bank| {
        total_output_joltage += maximal_jolt_in_bank_for_selection_of(bank.items, 0, 2);
    }

    return total_output_joltage;
}

fn part_two(battery_banks: std.ArrayList(std.ArrayList(u8))) u64 {
    var total_output_joltage: u64 = 0;

    for (battery_banks.items) |bank| {
        total_output_joltage += maximal_jolt_in_bank_for_selection_of(bank.items, 0, 12);
    }

    return total_output_joltage;
}

fn maximal_jolt_in_bank_for_selection_of(bank: []const u8, from: usize, selection_size: u8) u64 {
    var max_num_idx: usize = from;
    const remaining = selection_size - 1;

    for (from + 1..bank.len - remaining) |battery_idx| {
        if (bank[battery_idx] > bank[max_num_idx]) {
            max_num_idx = battery_idx;
        }
    }
    if (remaining > 0) {
        return bank[max_num_idx] * std.math.pow(u64, 10, remaining) + maximal_jolt_in_bank_for_selection_of(bank, max_num_idx + 1, remaining);
    }
    return bank[max_num_idx];
}
