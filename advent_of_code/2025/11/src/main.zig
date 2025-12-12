const std = @import("std");

pub fn main() !void {
    var threaded: std.Io.Threaded = .init_single_threaded;
    const io = threaded.io();

    var stdin_buffer: [1024]u8 = undefined;
    var stdin: std.fs.File = std.fs.File.stdin();
    var stdin_reader: std.fs.File.Reader = stdin.reader(io, &stdin_buffer);

    const allocator = std.heap.page_allocator;

    var lines = std.ArrayList(std.ArrayList(u8)).empty;
    errdefer lines.deinit(allocator);

    while (stdin_reader.interface.takeDelimiterInclusive('\n')) |line| {
        const trimmed = std.mem.trimRight(u8, line, "\r\n");
        var owned_line = std.ArrayList(u8).empty;
        try owned_line.appendSlice(allocator, trimmed);

        try lines.append(allocator, owned_line);
    } else |err| {
        if (err != error.EndOfStream) {
            std.debug.print("an error occurred: {any}\n", .{err});
        }
    }

    var adjacency_list = std.AutoHashMap(usize, std.AutoHashMap(usize, u1)).init(allocator);
    defer adjacency_list.deinit();

    var node_encoder = std.ArrayList(std.ArrayList(u8)).empty;
    defer node_encoder.deinit(allocator);

    try lines_to_encoded_array_list(allocator, lines.items, &node_encoder, &adjacency_list);

    std.debug.print("The many different paths lead from 'you' to 'out': {d}\n", .{try part_one(allocator, adjacency_list, node_encoder)});
    std.debug.print("Paths from 'svr' to 'out' through 'fft' and 'dac': {d}\n", .{try part_two(allocator, node_encoder, adjacency_list)});
}

// Transform the text id of the nodes into proper ids.
// We don't want the hassle of trying to hash strings with zig, nor pay the price for hashing all that nonsense.
fn lines_to_encoded_array_list(allocator: std.mem.Allocator, lines: []std.ArrayList(u8), node_encoder: *std.ArrayList(std.ArrayList(u8)), adjacency_list: *std.AutoHashMap(usize, std.AutoHashMap(usize, u1))) error{OutOfMemory}!void {
    for (lines) |line| {
        var node_name = std.ArrayList(u8).empty;
        try node_name.appendSlice(allocator, line.items[0..3]);
        // we could parse, but it is more code and not necessary here
        try node_encoder.append(allocator, node_name);
    }
    var output_node = std.ArrayList(u8).empty;
    try output_node.appendSlice(allocator, "out");
    try node_encoder.append(allocator, output_node);

    std.mem.sort(std.ArrayList(u8), node_encoder.items, {}, struct {
        fn lessThan(_: void, lhs: std.ArrayList(u8), rhs: std.ArrayList(u8)) bool {
            return std.mem.order(u8, lhs.items, rhs.items) == .lt;
        }
    }.lessThan);

    for (lines) |line| {
        const node_name = line.items[0..3];
        const from_id = find_node_id(node_encoder.*, node_name).?;

        var neighbors = std.AutoHashMap(usize, u1).init(allocator);
        // 5 so we skip the space after :
        var to_nodes = std.mem.splitScalar(u8, line.items[5..], ' ');
        while (to_nodes.next()) |neighbor| {
            const to_id = find_node_id(node_encoder.*, neighbor).?;

            try neighbors.put(to_id, 0);
        }
        try adjacency_list.put(from_id, neighbors);
    }
    try adjacency_list.put(find_node_id(node_encoder.*, "out").?, std.AutoHashMap(usize, u1).init(allocator));
}

// I am quite disappointed Zig doesn't have a general search function either.
fn find_node_id(node_encoder: std.ArrayList(std.ArrayList(u8)), node_name: []const u8) ?usize {
    for (node_encoder.items, 0..node_encoder.items.len) |node, i| {
        if (std.mem.eql(u8, node.items, node_name)) {
            return i;
        }
    }
    return null;
}

// Here we make the reasonable assumption that the input is of the form of a DAG.
// Otherwise there can be no solution anyways.
// Simple BFS tracing all paths at once from 'you' until 'out'.
fn part_one(allocator: std.mem.Allocator, adjacency_list: std.AutoHashMap(usize, std.AutoHashMap(usize, u1)), node_encoder: std.ArrayList(std.ArrayList(u8))) error{OutOfMemory}!u64 {
    const from = find_node_id(node_encoder, "you").?;
    const to = find_node_id(node_encoder, "out").?;

    return find_paths_between(allocator, adjacency_list, from, to);
}

fn find_paths_between(allocator: std.mem.Allocator, adjacency_list: std.AutoHashMap(usize, std.AutoHashMap(usize, u1)), from: usize, to: usize) error{OutOfMemory}!u64 {
    var next_nodes = std.AutoHashMap(usize, u64).init(allocator);
    defer next_nodes.deinit();

    var current_nodes = std.AutoHashMap(usize, u64).init(allocator);
    defer current_nodes.deinit();

    try current_nodes.put(from, 1);

    var path_count_to_out: u64 = 0;

    var seen = std.AutoHashMap(usize, u1).init(allocator);
    defer seen.deinit();

    while (current_nodes.count() != 0) {
        next_nodes.clearRetainingCapacity();

        var current_iterator = current_nodes.iterator();
        while (current_iterator.next()) |entry| {
            const neighbors = adjacency_list.get(entry.key_ptr.*).?;

            var neighbors_id_iterator = neighbors.keyIterator();
            while (neighbors_id_iterator.next()) |neighbor_id| {
                if (seen.contains(neighbor_id.*)) {
                    continue;
                }
                if (!next_nodes.contains(neighbor_id.*)) {
                    try next_nodes.put(neighbor_id.*, 0);
                }

                next_nodes.getPtr(neighbor_id.*).?.* += entry.value_ptr.*;
                if (neighbor_id.* == to) {
                    path_count_to_out += entry.value_ptr.*;
                }
            }
        }

        std.mem.swap(std.AutoHashMap(usize, u64), &current_nodes, &next_nodes);
    }

    return path_count_to_out;
}

fn part_two(allocator: std.mem.Allocator, node_encoder: std.ArrayList(std.ArrayList(u8)), adjacency_list: std.AutoHashMap(usize, std.AutoHashMap(usize, u1))) error{OutOfMemory}!u64 {
    const svr = find_node_id(node_encoder, "svr").?;
    const fft = find_node_id(node_encoder, "fft").?;
    const dac = find_node_id(node_encoder, "dac").?;
    const out = find_node_id(node_encoder, "out").?;

    const from_svr_to_fft = try find_paths_between(allocator, adjacency_list, svr, fft);
    const from_fft_to_dac = try find_paths_between(allocator, adjacency_list, fft, dac);
    const from_dac_to_out = try find_paths_between(allocator, adjacency_list, dac, out);

    // As we have established this cannot happen. Or there would be a loop.
    // const from_svr_to_dac = try find_paths_between(allocator, adjacency_list, svr, dac);
    // const from_dac_to_fft = try find_paths_between(allocator, adjacency_list, dac, fft);
    // const from_fft_to_out = try find_paths_between(allocator, adjacency_list, fft, out);

    const first_fft = from_svr_to_fft * from_fft_to_dac * from_dac_to_out;
    // const first_dac = from_svr_to_dac * from_dac_to_fft * from_fft_to_out;

    return first_fft;
    // return first_fft + first_dac;
}
