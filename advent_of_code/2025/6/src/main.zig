const std = @import("std");

const Op = enum { plus, mul };

pub fn main() !void {
    var threaded: std.Io.Threaded = .init_single_threaded;
    const io = threaded.io();

    var stdin_buffer: [10240]u8 = undefined;
    var stdin: std.fs.File = std.fs.File.stdin();
    var stdin_reader: std.fs.File.Reader = stdin.reader(io, &stdin_buffer);

    const allocator = std.heap.page_allocator;
    var numbers = std.ArrayList(std.ArrayList(u64)).empty;
    defer numbers.deinit(allocator);

    var ops = std.ArrayList(Op).empty;
    defer ops.deinit(allocator);

    var work_sheet = std.ArrayList(std.ArrayList(u8)).empty;
    defer work_sheet.deinit(allocator);

    while (stdin_reader.interface.takeDelimiterInclusive('\n')) |line| {
        const trimmed = std.mem.trim(u8, line, " \r\n");

        if (trimmed[0] == '+' or trimmed[0] == '-' or trimmed[0] == '*' or trimmed[0] == '/') {
            for (trimmed) |op| {
                if (op == ' ') {
                    continue;
                }
                const operation = switch (op) {
                    '+' => Op.plus,
                    '*' => Op.mul,
                    else => unreachable,
                };
                try ops.append(allocator, operation);
            }
        } else {
            // part one
            var it = std.mem.splitScalar(u8, trimmed, ' ');
            var number_line = std.ArrayList(u64).empty;
            while (it.next()) |value| {
                if (value.len == 0) {
                    continue;
                }
                try number_line.append(allocator, try std.fmt.parseInt(u64, value, 10));
            }
            try numbers.append(allocator, number_line);
        }
        // part two
        const sheet_line = std.mem.trimRight(u8, line, "\r\n");
        var entry = std.ArrayList(u8).empty;
        try entry.appendSlice(allocator, sheet_line);
        try work_sheet.append(allocator, entry);
    } else |err| {
        if (err != error.EndOfStream) {
            std.debug.print("an error occurred: {any}\n", .{err});
        }
    }

    const number_columns = try transpose(allocator, numbers);

    std.debug.print("The grand total is: {d}\n", .{part_one(number_columns, ops.items)});
    std.debug.print("The grand total with cephalopod math is: {d}\n", .{try part_two(work_sheet)});
}

fn transpose(allocator: std.mem.Allocator, numbers: std.ArrayList(std.ArrayList(u64))) error{OutOfMemory}!std.ArrayList(std.ArrayList(u64)) {
    var number_columns = std.ArrayList(std.ArrayList(u64)).empty;
    for (0..numbers.items[0].items.len) |j| {
        try number_columns.append(allocator, std.ArrayList(u64).empty);
        for (0..numbers.items.len) |i| {
            try number_columns.items[j].append(allocator, numbers.items[i].items[j]);
        }
    }
    return number_columns;
}

fn part_one(numbers: std.ArrayList(std.ArrayList(u64)), ops: []const Op) u64 {
    var grand_total: u64 = 0;

    for (0..numbers.items.len) |i| {
        const op = ops[i];
        var sub_total: u64 = numbers.items[i].items[0];
        for (1..numbers.items[i].items.len) |j| {
            switch (op) {
                Op.plus => sub_total += numbers.items[i].items[j],
                Op.mul => sub_total *= numbers.items[i].items[j],
            }
        }
        grand_total += sub_total;
    }

    return grand_total;
}

fn part_two(work_sheet: std.ArrayList(std.ArrayList(u8))) error{OutOfMemory}!u64 {
    var grand_total: u64 = 0;

    // At every index where an operator is present, the corresponding
    // values start in the same column and they either go until the end
    // of the line or 1 index before the next operator.
    // So it is best to get the operator positions first.
    const allocator = std.heap.page_allocator;
    var operator_positions = std.ArrayList(usize).empty;
    defer operator_positions.deinit(allocator);

    const last_row = work_sheet.items[work_sheet.items.len - 1];
    for (0..last_row.items.len) |i| {
        if (last_row.items[i] != ' ') {
            try operator_positions.append(allocator, i);
        }
    }

    for (0..operator_positions.items.len) |op_i| {
        const values_from = operator_positions.items[op_i];
        const values_until = if (op_i < operator_positions.items.len - 1) operator_positions.items[op_i + 1] else last_row.items.len + 1;

        var sub_total = number_from_column(work_sheet.items[0 .. work_sheet.items.len - 1], values_from);

        for (values_from + 1..values_until - 1) |column| {
            const value = number_from_column(work_sheet.items[0 .. work_sheet.items.len - 1], column);

            if (last_row.items[values_from] == '*') {
                sub_total *= value;
            } else {
                sub_total += value;
            }
        }

        grand_total += sub_total;
    }

    return grand_total;
}

fn number_from_column(work_sheet: []const std.ArrayList(u8), column: usize) u64 {
    var number: u64 = 0;

    for (0..work_sheet.len) |i| {
        const c = work_sheet[i].items[column];
        const digit = if (c == ' ') continue else c - '0';
        number *= 10;
        number += digit;
    }

    return number;
}
