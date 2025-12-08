const std = @import("std");

const Distance = struct {
    from: usize,
    to: usize,
    distance: f64,
};

const Coordinate = struct {
    x: f64,
    y: f64,
    z: f64,
};

pub fn main() !void {
    var threaded: std.Io.Threaded = .init_single_threaded;
    const io = threaded.io();

    var stdin_buffer: [1024]u8 = undefined;
    var stdin: std.fs.File = std.fs.File.stdin();
    var stdin_reader: std.fs.File.Reader = stdin.reader(io, &stdin_buffer);

    const allocator = std.heap.page_allocator;
    var junction_boxes = std.ArrayList(Coordinate).empty;
    defer junction_boxes.deinit(allocator);

    while (stdin_reader.interface.takeDelimiterInclusive('\n')) |line| {
        const trimmed = std.mem.trimRight(u8, line, "\r\n");

        var coordinates = std.mem.splitScalar(u8, trimmed, ',');
        const x = try std.fmt.parseFloat(f64, coordinates.next().?);
        const y = try std.fmt.parseFloat(f64, coordinates.next().?);
        const z = try std.fmt.parseFloat(f64, coordinates.next().?);
        std.debug.assert(coordinates.next() == null);

        try junction_boxes.append(allocator, Coordinate{ .x = x, .y = y, .z = z });
    } else |err| {
        if (err != error.EndOfStream) {
            std.debug.print("an error occurred: {any}\n", .{err});
        }
    }

    var distances = try calculate_distances(allocator, junction_boxes);
    const S = struct {
        fn sort_by_distance_ascending(context: void, lhs: Distance, rhs: Distance) bool {
            _ = context;
            return lhs.distance < rhs.distance;
        }
    };
    std.mem.sort(Distance, distances.items, {}, S.sort_by_distance_ascending);

    std.debug.print("The 3 largest circuit size multiple is: {d}\n", .{try part_one(allocator, junction_boxes, distances, 1000)});
    std.debug.print("The multiple of the X coordinate of the last two junction boxes: {d}\n", .{try part_two(allocator, junction_boxes, distances)});
}

fn calculate_distances(allocator: std.mem.Allocator, coordinates: std.ArrayList(Coordinate)) error{ OutOfMemory, Overflow, Underflow }!std.ArrayList(Distance) {
    var distances = std.ArrayList(Distance).empty;

    for (0..coordinates.items.len - 1) |from| {
        for (from + 1..coordinates.items.len) |to| {
            const from_coord = coordinates.items[from];
            const to_coord = coordinates.items[to];

            const distance = std.math.sqrt(std.math.pow(f64, to_coord.x - from_coord.x, 2) + std.math.pow(f64, to_coord.y - from_coord.y, 2) + std.math.pow(f64, to_coord.z - from_coord.z, 2));
            try distances.append(allocator, Distance{ .from = from, .to = to, .distance = distance });
        }
    }
    return distances;
}

fn part_one(allocator: std.mem.Allocator, junction_boxes: std.ArrayList(Coordinate), distances: std.ArrayList(Distance), connection_count: usize) error{OutOfMemory}!usize {
    var connections_made: usize = 0;

    var dsu_parents = std.ArrayList(usize).empty;
    defer dsu_parents.deinit(allocator);
    try dsu_parents.appendNTimes(allocator, 0, junction_boxes.items.len);
    for (0..junction_boxes.items.len) |i| {
        make_set(dsu_parents, i);
    }

    var dsu_ranks = std.ArrayList(usize).empty;
    defer dsu_ranks.deinit(allocator);
    try dsu_ranks.appendNTimes(allocator, 0, junction_boxes.items.len);

    var distance_idx: usize = 0;
    while (connections_made < connection_count) {
        const distance = distances.items[distance_idx];
        distance_idx += 1;

        // You would think the problem wanted you to skip shortest connections, which don't modify the set,
        // but it doesn't.
        // Realized this after, seeing that the real input has a 1000 nodes and the task ask you to
        // add a 1000 connections, then multiply the biggest 3 groups. Which can't make sense since,
        // 501 connections would be enough to fully connect everything if these redundant connections
        // are skipped.
        // if (find_set(dsu_parents, distance.from) == find_set(dsu_parents, distance.to)) {
        //     std.debug.print("Already in same set: {} - {}\n", .{ distance.from, distance.to });
        //     continue;
        // }
        union_sets(dsu_parents, dsu_ranks, distance.from, distance.to);
        connections_made += 1;
    }

    // Compress the DSU tree completely, otherwise we can't count properly the sets later.
    for (0..dsu_parents.items.len) |i| {
        const p = find_set(dsu_parents, i);
        _ = p;
    }

    var set_counts = std.AutoHashMap(usize, usize).init(allocator);
    defer set_counts.deinit();

    for (dsu_parents.items) |parent| {
        if (!set_counts.contains(parent)) {
            try set_counts.put(parent, 1);
        } else {
            const entry = set_counts.getPtr(parent).?;
            entry.* += 1;
        }
    }

    var size_ordered = std.ArrayList(usize).empty;
    defer size_ordered.deinit(allocator);

    var iterator = set_counts.iterator();
    while (iterator.next()) |entry| {
        try size_ordered.append(allocator, entry.value_ptr.*);
    }
    const S = struct {
        fn sort_by_size_descending(context: void, lhs: usize, rhs: usize) bool {
            _ = context;
            return lhs > rhs;
        }
    };
    std.mem.sort(usize, size_ordered.items, {}, S.sort_by_size_descending);

    return size_ordered.items[0] * size_ordered.items[1] * size_ordered.items[2];
}

// quick DSU implementation
fn make_set(parent: std.ArrayList(usize), set_id: usize) void {
    parent.items[set_id] = set_id;
}

fn find_set(parent: std.ArrayList(usize), set_id: usize) usize {
    if (parent.items[set_id] == set_id) {
        return set_id;
    }
    const top_parent = find_set(parent, parent.items[set_id]);
    parent.items[set_id] = top_parent;
    return top_parent;
}

fn union_sets(parent: std.ArrayList(usize), rank: std.ArrayList(usize), lhs: usize, rhs: usize) void {
    var parent_lhs = find_set(parent, lhs);
    var parent_rhs = find_set(parent, rhs);
    if (parent_lhs != parent_rhs) {
        if (rank.items[parent_lhs] < rank.items[parent_rhs]) {
            std.mem.swap(usize, &parent_lhs, &parent_rhs);
        }
        parent.items[parent_rhs] = parent_lhs;
        if (rank.items[parent_lhs] == rank.items[parent_rhs]) {
            rank.items[parent_lhs] += 1;
        }
    }
}

fn part_two(allocator: std.mem.Allocator, junction_boxes: std.ArrayList(Coordinate), distances: std.ArrayList(Distance)) error{OutOfMemory}!u64 {
    var dsu_parents = std.ArrayList(usize).empty;
    defer dsu_parents.deinit(allocator);
    try dsu_parents.appendNTimes(allocator, 0, junction_boxes.items.len);
    for (0..junction_boxes.items.len) |i| {
        make_set(dsu_parents, i);
    }

    var dsu_ranks = std.ArrayList(usize).empty;
    defer dsu_ranks.deinit(allocator);
    try dsu_ranks.appendNTimes(allocator, 0, junction_boxes.items.len);

    var fully_connected = false;
    var distance_idx: usize = 0;
    while (!fully_connected) {
        const distance = distances.items[distance_idx];
        distance_idx += 1;

        union_sets(dsu_parents, dsu_ranks, distance.from, distance.to);
        // make sure to compress the tree
        for (0..dsu_parents.items.len) |i| {
            const p = find_set(dsu_parents, i);
            _ = p;
        }
        // inefficient check but hey, good enough, not in the mood to make it faster
        fully_connected = true;
        const value = dsu_parents.items[0];
        for (1..dsu_parents.items.len) |i| {
            if (dsu_parents.items[i] != value) {
                fully_connected = false;
                break;
            }
        }
    }

    // distance_idx-1 because the loop iterates it before ending
    const last_connection = distances.items[distance_idx - 1];
    return @as(u64, @intFromFloat(junction_boxes.items[last_connection.from].x)) * @as(u64, @intFromFloat(junction_boxes.items[last_connection.to].x));
}
