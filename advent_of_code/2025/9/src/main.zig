const std = @import("std");

const Pos = struct {
    i: i64,
    j: i64,
};

const Line = struct {
    a: Pos,
    b: Pos,
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
        const j = try std.fmt.parseInt(i64, coordinates.next().?, 10);
        const i = try std.fmt.parseInt(i64, coordinates.next().?, 10);
        std.debug.assert(coordinates.next() == null);

        try positions.append(allocator, Pos{ .i = i, .j = j });
    } else |err| {
        if (err != error.EndOfStream) {
            std.debug.print("an error occurred: {any}\n", .{err});
        }
    }

    std.debug.print("The largest rectangle possible: {d}\n", .{part_one(positions)});

    var horizontal_lines = std.ArrayList(Line).empty;
    defer horizontal_lines.deinit(allocator);

    var vertical_lines = std.ArrayList(Line).empty;
    defer vertical_lines.deinit(allocator);

    try line_segments_from_points(allocator, positions, &horizontal_lines, &vertical_lines);

    // sort horizontal in increasing J order
    std.mem.sort(Line, horizontal_lines.items, {}, struct {
        fn lessThanFn(context: void, lhs: Line, rhs: Line) bool {
            _ = context;
            return lhs.a.j < rhs.a.j;
        }
    }.lessThanFn);

    // sort vertical lines in increasing I order
    std.mem.sort(Line, vertical_lines.items, {}, struct {
        fn lessThanFn(context: void, lhs: Line, rhs: Line) bool {
            _ = context;
            return lhs.a.i < rhs.a.i;
        }
    }.lessThanFn);

    var timer = try std.time.Timer.start();
    std.debug.print("The biggest possible rectangle within the polygon: {d}\n", .{try part_two(allocator, positions, vertical_lines, horizontal_lines)});
    std.debug.print("Part 2 execution time: {} ns\n", .{timer.lap()});
    // for me: Part 2 execution time: 5241200 ns => 0.0053487 s, quite good for a general solution
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

fn line_segments_from_points(allocator: std.mem.Allocator, positions: std.ArrayList(Pos), horizontal_lines: *std.ArrayList(Line), vertical_lines: *std.ArrayList(Line)) error{OutOfMemory}!void {
    for (0..positions.items.len) |i| {
        var a = positions.items[i];

        var j: usize = 0;
        if (i < positions.items.len - 1) {
            j = i + 1;
        }
        var b = positions.items[j];

        if (a.i == b.i) {
            // horizontal
            if (a.j > b.j) {
                std.mem.swap(Pos, &a, &b);
            }
            try horizontal_lines.append(allocator, Line{ .a = a, .b = b });
        } else if (a.j == b.j) {
            // vertical
            if (a.i > b.i) {
                std.mem.swap(Pos, &a, &b);
            }
            try vertical_lines.append(allocator, Line{ .a = a, .b = b });
        } else {
            unreachable;
        }
    }
}

fn intervals_overlap(lhs_start: i64, lhs_end: i64, rhs_start: i64, rhs_end: i64) bool {
    // Make sure the values are in ascending order.
    var l_s = lhs_start;
    var l_e = lhs_end;
    if (l_e < l_s) {
        std.mem.swap(i64, &l_s, &l_e);
    }

    var r_s = rhs_start;
    var r_e = rhs_end;
    if (r_e < r_s) {
        std.mem.swap(i64, &r_s, &r_e);
    }

    if (l_e <= r_s or l_s >= r_e) {
        return false;
    }

    return true;
}

const Span = struct {
    from: i64,
    until: i64,
};

const Rectangle = struct {
    h: Span,
    v: Span,
};

fn part_two(allocator: std.mem.Allocator, positions: std.ArrayList(Pos), vertical_lines: std.ArrayList(Line), horizontal_lines: std.ArrayList(Line)) error{OutOfMemory}!u64 {
    const horizontal_internal_spans = try generate_horizontal_internal_spans(allocator, vertical_lines);

    const vertical_internal_spans = try generate_vertical_internal_spans(allocator, horizontal_lines);

    var biggest_internal_rectangle_area: u64 = 0;

    for (0..positions.items.len - 1) |i| {
        const corner_A = positions.items[i];
        for (i + 1..positions.items.len) |j| {
            const corner_B = positions.items[j];

            const vertical_span = to_span(corner_A.i, corner_B.i);
            const horizontal_span = to_span(corner_A.j, corner_B.j);

            const valid_horizontally = check_horizontal_bounds(horizontal_span, vertical_span, horizontal_internal_spans);

            const valid_vertically = check_vertical_bounds(horizontal_span, vertical_span, vertical_internal_spans);

            const delta_i = @as(u64, @intCast(vertical_span.until - vertical_span.from)) + 1;
            const delta_j = @as(u64, @intCast(horizontal_span.until - horizontal_span.from)) + 1;
            const area = delta_i * delta_j;

            if (!(valid_horizontally and valid_vertically)) {
                continue;
            }

            if (biggest_internal_rectangle_area < area) {
                biggest_internal_rectangle_area = area;
            }
        }
    }

    return biggest_internal_rectangle_area;
}

fn check_horizontal_bounds(horizontal: Span, vertical: Span, internal_horizontals: std.ArrayList(OrderedSpan)) bool {
    var idx = std.sort.lowerBound(OrderedSpan, internal_horizontals.items, vertical.from, struct {
        fn compareOrder(context: i64, item: OrderedSpan) std.math.Order {
            return std.math.order(context, item.order);
        }
    }.compareOrder);

    while (idx < internal_horizontals.items.len) {
        const ordered_span = internal_horizontals.items[idx];

        if (!span_contained_by(ordered_span.span, horizontal)) {
            return false;
        }
        if (ordered_span.order == vertical.until) {
            break;
        }

        idx += 1;
    }

    return true;
}

fn check_vertical_bounds(horizontal: Span, vertical: Span, internal_verticals: std.ArrayList(OrderedSpan)) bool {
    var idx = std.sort.lowerBound(OrderedSpan, internal_verticals.items, horizontal.from, struct {
        fn compareOrder(context: i64, item: OrderedSpan) std.math.Order {
            return std.math.order(context, item.order);
        }
    }.compareOrder);

    while (idx < internal_verticals.items.len) {
        const ordered_span = internal_verticals.items[idx];

        if (!span_contained_by(ordered_span.span, vertical)) {
            return false;
        }
        if (ordered_span.order == horizontal.until) {
            break;
        }

        idx += 1;
    }

    return true;
}

fn span_contained_by(by: Span, candidate: Span) bool {
    return by.from <= candidate.from and candidate.until <= by.until;
}

fn to_span(a: i64, b: i64) Span {
    if (b < a) {
        return Span{ .from = b, .until = a };
    }
    return Span{ .from = a, .until = b };
}

const OrderedSpan = struct {
    order: i64,
    span: Span,
};

fn generate_horizontal_internal_spans(allocator: std.mem.Allocator, vertical_lines: std.ArrayList(Line)) error{OutOfMemory}!std.ArrayList(OrderedSpan) {
    var ordered_internal_spans = std.ArrayList(OrderedSpan).empty;

    var active_lines = std.ArrayList(Line).empty;
    defer active_lines.deinit(allocator);

    var next_active_lines = std.ArrayList(Line).empty;
    defer next_active_lines.deinit(allocator);

    var internal_spans = std.ArrayList(Span).empty;
    defer internal_spans.deinit(allocator);

    var i: usize = 0;
    var at_order = vertical_lines.items[0].a.i;
    while (true) {
        next_active_lines.clearRetainingCapacity();
        // evict all active line segments that are no longer in range
        for (active_lines.items) |active| {
            if (active.b.i < at_order) {
                continue;
            }
            try next_active_lines.append(allocator, active);
        }
        // add unused line segments that start at "at_order"
        while (i < vertical_lines.items.len) {
            const candidate = vertical_lines.items[i];
            if (candidate.a.i != at_order) {
                break;
            }
            const S = struct {
                fn horizontalPosition(context: Line, item: Line) std.math.Order {
                    if (context.a.j < item.a.j) {
                        return std.math.Order.lt;
                    } else if (context.a.j > item.a.j) {
                        return std.math.Order.gt;
                    } else {
                        return std.math.Order.eq;
                    }
                }
            };

            const insert_position = std.sort.upperBound(Line, next_active_lines.items, candidate, S.horizontalPosition);
            try next_active_lines.insert(allocator, insert_position, candidate);

            i += 1;
        }

        // If at this point we have no next_active lines, that means
        // there are no more internal sections to find.
        // This is because we have only one polygon, which
        // consists only of horizontal or vertical lines.
        // (Won't work if there are multiple disjoint polygons.)
        if (next_active_lines.items.len == 0) {
            break;
        }

        // go through the active line segments by pair and log down the internal spans
        internal_spans.clearRetainingCapacity();
        var n_i: usize = 0;
        while (n_i < next_active_lines.items.len - 1) {
            const from = next_active_lines.items[n_i];
            const until = next_active_lines.items[n_i + 1];

            try internal_spans.append(allocator, Span{ .from = from.a.j, .until = until.a.j });

            n_i += 1;
        }
        // now merge all internal spans that overlap and add them to the ordered_internal_spans
        n_i = 0;
        while (n_i < internal_spans.items.len) {
            const from = internal_spans.items[n_i].from;
            var until = internal_spans.items[n_i].until;

            n_i += 1;
            while (n_i < internal_spans.items.len) {
                if (until != internal_spans.items[n_i].from) {
                    break;
                }
                until = internal_spans.items[n_i].until;
                n_i += 1;
            }
            try ordered_internal_spans.append(allocator, OrderedSpan{ .order = at_order, .span = Span{ .from = from, .until = until } });
        }
        // find the next at_order level, which is the nearest closest from the current at_order
        var next_at_order: i64 = std.math.maxInt(i64);
        if (i < vertical_lines.items.len) {
            if (vertical_lines.items[i].a.i < next_at_order) {
                next_at_order = vertical_lines.items[i].a.i;
            }
        }
        for (next_active_lines.items) |active| {
            if (active.b.i != at_order and active.b.i < next_at_order) {
                next_at_order = active.b.i;
            }
        }
        at_order = next_at_order;

        std.mem.swap(std.ArrayList(Line), &active_lines, &next_active_lines);
    }

    return ordered_internal_spans;
}

fn generate_vertical_internal_spans(allocator: std.mem.Allocator, horizontal_lines: std.ArrayList(Line)) error{OutOfMemory}!std.ArrayList(OrderedSpan) {
    var ordered_internal_spans = std.ArrayList(OrderedSpan).empty;

    var active_lines = std.ArrayList(Line).empty;
    defer active_lines.deinit(allocator);

    var next_active_lines = std.ArrayList(Line).empty;
    defer next_active_lines.deinit(allocator);

    var internal_spans = std.ArrayList(Span).empty;
    defer internal_spans.deinit(allocator);

    var i: usize = 0;
    var at_order = horizontal_lines.items[0].a.j;
    while (true) {
        next_active_lines.clearRetainingCapacity();
        // evict all active line segments that are no longer in range
        for (active_lines.items) |active| {
            if (active.b.j < at_order) {
                continue;
            }
            try next_active_lines.append(allocator, active);
        }
        // add unused line segments that start at "at_order"
        while (i < horizontal_lines.items.len) {
            const candidate = horizontal_lines.items[i];
            if (candidate.a.j != at_order) {
                break;
            }
            const S = struct {
                fn verticalPosition(context: Line, item: Line) std.math.Order {
                    if (context.a.i < item.a.i) {
                        return std.math.Order.lt;
                    } else if (context.a.i > item.a.i) {
                        return std.math.Order.gt;
                    } else {
                        return std.math.Order.eq;
                    }
                }
            };

            const insert_position = std.sort.upperBound(Line, next_active_lines.items, candidate, S.verticalPosition);
            try next_active_lines.insert(allocator, insert_position, candidate);

            i += 1;
        }

        // If at this point we have no next_active lines, that means
        // there are no more internal sections to find.
        // This is because we have only one polygon, which
        // consists only of horizontal or vertical lines.
        // (Won't work if there are multiple disjoint polygons.)
        if (next_active_lines.items.len == 0) {
            break;
        }

        // go through the active line segments by pair and log down the internal spans
        internal_spans.clearRetainingCapacity();
        var n_i: usize = 0;
        while (n_i < next_active_lines.items.len - 1) {
            const from = next_active_lines.items[n_i];
            const until = next_active_lines.items[n_i + 1];

            try internal_spans.append(allocator, Span{ .from = from.a.i, .until = until.a.i });

            n_i += 1;
        }
        // now merge all internal spans that overlap and add them to the ordered_internal_spans
        n_i = 0;
        while (n_i < internal_spans.items.len) {
            const from = internal_spans.items[n_i].from;
            var until = internal_spans.items[n_i].until;

            n_i += 1;
            while (n_i < internal_spans.items.len) {
                if (until != internal_spans.items[n_i].from) {
                    break;
                }
                until = internal_spans.items[n_i].until;
                n_i += 1;
            }
            try ordered_internal_spans.append(allocator, OrderedSpan{ .order = at_order, .span = Span{ .from = from, .until = until } });
        }
        // find the next at_order level, which is the nearest closest from the current at_order
        var next_at_order: i64 = std.math.maxInt(i64);
        if (i < horizontal_lines.items.len) {
            if (horizontal_lines.items[i].a.j < next_at_order) {
                next_at_order = horizontal_lines.items[i].a.j;
            }
        }
        for (next_active_lines.items) |active| {
            if (active.b.j != at_order and active.b.j < next_at_order) {
                next_at_order = active.b.j;
            }
        }
        at_order = next_at_order;

        std.mem.swap(std.ArrayList(Line), &active_lines, &next_active_lines);
    }

    return ordered_internal_spans;
}

fn abs_diff(a: i64, b: i64) u64 {
    if (a < b) {
        return @as(u64, @intCast(b)) - @as(u64, @intCast(a));
    }
    return @as(u64, @intCast(a)) - @as(u64, @intCast(b));
}

// too low: 223562430
//          1496002500
//          1569262188
// not correct: 2603882976
// too high:4545726068
