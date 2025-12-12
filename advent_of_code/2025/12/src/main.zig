const std = @import("std");

const Shape = [3][3]bool;

const Task = struct {
    width: usize,
    height: usize,
    shape_counts: std.ArrayList(u8),
};

pub fn main() !void {
    var threaded: std.Io.Threaded = .init_single_threaded;
    const io = threaded.io();

    var stdin_buffer: [1024]u8 = undefined;
    var stdin: std.fs.File = std.fs.File.stdin();
    var stdin_reader: std.fs.File.Reader = stdin.reader(io, &stdin_buffer);

    const allocator = std.heap.page_allocator;

    var lines = std.ArrayList(std.ArrayList(u8)).empty;
    defer lines.deinit(allocator);

    while (stdin_reader.interface.takeDelimiterInclusive('\n')) |file_line| {
        const trimmed = std.mem.trimRight(u8, file_line, "\r\n");

        var line = std.ArrayList(u8).empty;
        try line.appendSlice(allocator, trimmed);
        try lines.append(allocator, line);
    } else |err| {
        if (err != error.EndOfStream) {
            std.debug.print("an error occurred: {any}\n", .{err});
        }
    }

    for (lines.items) |line| {
        std.debug.print("{s}\n", .{line.items});
    }

    var shapes = std.ArrayList(Shape).empty;
    defer shapes.deinit(allocator);

    var tasks = std.ArrayList(Task).empty;
    defer tasks.deinit(allocator);

    try parse_input(allocator, lines, &shapes, &tasks);

    var melted_shape_area = std.ArrayList(u8).empty;
    defer melted_shape_area.deinit(allocator);

    for (shapes.items) |shape| {
        var shape_area_count: u8 = 0;
        for (0..3) |i| {
            for (0..3) |j| {
                if (shape[i][j]) {
                    shape_area_count += 1;
                }
            }
        }
        try melted_shape_area.append(allocator, shape_area_count);
    }

    var impossible_count: usize = 0;
    var possible_count: usize = 0;

    for (tasks.items) |task| {
        var melted_area: usize = 0;
        var as_unmergeable_area: usize = 0;
        for (0..task.shape_counts.items.len) |i| {
            const multiple = task.shape_counts.items[i];
            melted_area += @as(usize, multiple) * @as(usize, melted_shape_area.items[i]);

            as_unmergeable_area += @as(usize, multiple) * 9;
        }
        const available_area = task.width * task.height;
        if (melted_area > available_area) {
            impossible_count += 1;
        }
        if (melted_area < available_area) {
            possible_count += 1;
        }
    }

    std.debug.print("The number of present stacks that are\npossible: {}\nimpossible: {}\nundetermined: {}\n", .{ possible_count, impossible_count, tasks.items.len - possible_count - impossible_count });
}

fn parse_input(allocator: std.mem.Allocator, lines: std.ArrayList(std.ArrayList(u8)), shapes: *std.ArrayList(Shape), tasks: *std.ArrayList(Task)) error{ OutOfMemory, Overflow, InvalidCharacter }!void {
    var at_line: usize = 0;

    var parsing_shapes = true;
    while (at_line < lines.items.len) {
        const line = lines.items[at_line];
        if (line.items.len == 0) {
            at_line += 1;
            continue;
        }

        if (line.getLast() == ':') {
            std.debug.assert(parsing_shapes);
            at_line = try parse_shape(allocator, lines, at_line, shapes);
        } else {
            parsing_shapes = false;
            at_line = try parse_task(allocator, lines, at_line, tasks, shapes.items.len);
        }
    }
}

fn parse_shape(allocator: std.mem.Allocator, lines: std.ArrayList(std.ArrayList(u8)), at_line: usize, shapes: *std.ArrayList(Shape)) error{ OutOfMemory, Overflow, InvalidCharacter }!usize {
    const id_line = lines.items[at_line].items;
    const id = try std.fmt.parseUnsigned(usize, id_line[0 .. id_line.len - 1], 10);

    var shape: [3][3]bool = undefined;
    for (0..3) |offset_i| {
        const line = lines.items[at_line + 1 + offset_i];
        std.debug.assert(line.items.len == 3);
        for (0..3) |j| {
            if (line.items[j] == '#') {
                shape[offset_i][j] = true;
            } else if (line.items[j] == '.') {
                shape[offset_i][j] = false;
            } else {
                unreachable;
            }
        }
    }

    std.debug.assert(shapes.items.len == id);
    try shapes.append(allocator, shape);

    return at_line + 4;
}

fn parse_task(allocator: std.mem.Allocator, lines: std.ArrayList(std.ArrayList(u8)), at_line: usize, tasks: *std.ArrayList(Task), expected_count: usize) error{ OutOfMemory, InvalidCharacter, Overflow }!usize {
    const task_line = lines.items[at_line];

    const colon_idx = std.mem.findScalarPos(u8, task_line.items, 0, ':').?;
    const x_idx = std.mem.findScalarPos(u8, task_line.items[0..colon_idx], 0, 'x').?;

    const width = try std.fmt.parseUnsigned(usize, task_line.items[0..x_idx], 10);
    const height = try std.fmt.parseUnsigned(usize, task_line.items[x_idx + 1 .. colon_idx], 10);

    var count_iterator = std.mem.splitScalar(u8, task_line.items[colon_idx + 2 ..], ' ');

    var shape_counts = std.ArrayList(u8).empty;

    while (count_iterator.next()) |count_text| {
        const count = try std.fmt.parseUnsigned(u8, count_text, 10);
        try shape_counts.append(allocator, count);
    }
    std.debug.assert(shape_counts.items.len == expected_count);

    const task = Task{ .width = width, .height = height, .shape_counts = shape_counts };
    try tasks.append(allocator, task);

    return at_line + 1;
}