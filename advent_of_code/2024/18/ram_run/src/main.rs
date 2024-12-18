use std::{collections::VecDeque, io, usize, vec};

static DIRS: [(i32, i32); 4] = [(0, 1), (1, 0), (0, -1), (-1, 0)];

fn parse_input() -> Vec<(i32, i32)> {
    io::stdin()
        .lines()
        .map(|line| {
            let line = line.unwrap();
            let parts = line.split_once(',').unwrap();
            (parts.0.parse().unwrap(), parts.1.parse().unwrap())
        })
        .collect()
}

struct State {
    pub pos: (i32, i32),
    pub cost: usize,
}

fn part_one(byte_positions: &Vec<(i32, i32)>, side_length: usize, fall_count: usize) -> usize {
    let mut memory_area = vec![vec!['.'; side_length]; side_length];

    for i in 0..fall_count {
        let pos = byte_positions[i];
        memory_area[pos.0 as usize][pos.1 as usize] = '#';
    }

    let cheapest_route_costs = find_cheapest_route_costs(
        &memory_area,
        (0, 0),
        (side_length as i32 - 1, side_length as i32 - 1),
    );

    cheapest_route_costs[side_length - 1][side_length - 1]
}

fn find_cheapest_route_costs(
    memory_area: &Vec<Vec<char>>,
    start_pos: (i32, i32),
    target_pos: (i32, i32),
) -> Vec<Vec<usize>> {
    let mut search_states = VecDeque::<State>::new();

    search_states.push_back(State {
        pos: start_pos,
        cost: 0,
    });

    // In this case we just simply do a BFS, and just step. Since the cost of a move
    // is only one and thus the cheapest path between two points is simply the shortest,
    // this will very quickly map out the whole search space.
    let mut best_costs = vec![vec![usize::MAX; memory_area[0].len()]; memory_area.len()];
    while !search_states.is_empty() {
        let state = search_states.pop_front().unwrap();

        let field = memory_area[state.pos.0 as usize][state.pos.1 as usize];
        // we don't set min for walls, so we don't have to worry about later
        if field == '#' {
            continue;
        }

        let current_best = best_costs[state.pos.0 as usize][state.pos.1 as usize];
        if current_best <= state.cost {
            continue;
        } else {
            best_costs[state.pos.0 as usize][state.pos.1 as usize] = state.cost;
        }
        if state.pos == target_pos {
            continue;
        }

        for dir in DIRS {
            let next_pos = (state.pos.0 + dir.0, state.pos.1 + dir.1);
            if next_pos.0 < 0
                || next_pos.0 >= memory_area.len() as i32
                || next_pos.1 < 0
                || next_pos.1 >= memory_area.len() as i32
            {
                continue;
            }
            let start_state = State {
                pos: next_pos,
                cost: state.cost + 1,
            };
            search_states.push_back(start_state);
        }
    }

    best_costs
}

// Here we are just run an upper bound binary search on the fallen block count.
fn part_two(byte_positions: &Vec<(i32, i32)>, side_length: usize) -> (i32, i32) {
    let mut low = 0;
    let mut high = byte_positions.len();

    while high - low > 0 {
        let mid = ((high - low) / 2) + low;

        let is_reacheable = if part_one(byte_positions, side_length, mid) == usize::MAX {
            false
        } else {
            true
        };
        if is_reacheable {
            low = mid + 1;
        } else {
            high = mid;
        }
    }

    byte_positions[low - 1]
}

fn main() {
    let byte_positions = parse_input();

    println!(
        "The minimum number of steps needed to reach the exit: {}",
        part_one(&byte_positions, 71, 1024)
    );

    println!(
        "The coordinates of the first byte that will prevent the exit from being reachable: {:?}",
        part_two(&byte_positions, 71)
    );
}
