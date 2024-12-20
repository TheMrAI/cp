use std::{collections::VecDeque, io, usize, vec};

static DIRS: [(i32, i32); 4] = [(0, 1), (1, 0), (0, -1), (-1, 0)];

fn parse_input() -> Vec<Vec<char>> {
    io::stdin()
        .lines()
        .map(|line| {
            let line = line.unwrap();
            line.chars().collect()
        })
        .collect()
}

struct State {
    pub pos: (i32, i32),
    pub cost: usize,
}

// From the problem description we know that there is only one route from
// start to finish, and from looking at the example it seems that all open
// space is used up to get from Start to End, so we don't even need to know
// where the path goes, only mark out how many steps it takes to get to each
// field.
fn part_one(racetrack: &Vec<Vec<char>>) -> usize {
    let mut start_pos = (0, 0);
    let mut target_pos = (0, 0);
    for i in 0..racetrack.len() {
        for j in 0..racetrack[0].len() {
            if racetrack[i][j] == 'S' {
                start_pos = (i as i32, j as i32);
            } else if racetrack[i][j] == 'E' {
                target_pos = (i as i32, j as i32);
            }
        }
    }

    let cheapest_route_costs = find_cheapest_route_costs(&racetrack, start_pos, target_pos);

    let mut saving_at_least_100_pico_count = 0;

    // now we just need trace the route and attempt a cheat in the four cardinal directions
    let mut current_pos = start_pos;
    while current_pos != target_pos {
        let current_cost = cheapest_route_costs[current_pos.0 as usize][current_pos.1 as usize];
        let mut next_pos = current_pos;

        for dir in DIRS {
            let next_pos_candidate = (current_pos.0 + dir.0, current_pos.1 + dir.1);
            if racetrack[next_pos_candidate.0 as usize][next_pos_candidate.1 as usize] == '#' {
                // check for glitch
                let glitch_candidate = (next_pos_candidate.0 + dir.0, next_pos_candidate.1 + dir.1);
                if glitch_candidate.0 < 0
                    || glitch_candidate.0 >= racetrack.len() as i32
                    || glitch_candidate.1 < 0
                    || glitch_candidate.1 >= racetrack[0].len() as i32
                    || racetrack[glitch_candidate.0 as usize][glitch_candidate.1 as usize] == '#'
                {
                    continue;
                }
                if let Some(cost_saving) = cheapest_route_costs[glitch_candidate.0 as usize]
                    [glitch_candidate.1 as usize]
                    .checked_sub(current_cost)
                {
                    if cost_saving <= 2 {
                        continue;
                    }
                    // account for the cheat cost of 2 picoseconds
                    if cost_saving - 2 >= 100 {
                        saving_at_least_100_pico_count += 1;
                    }
                }
            } else {
                // check for next position
                if let Some(delta) = cheapest_route_costs[next_pos_candidate.0 as usize]
                    [next_pos_candidate.1 as usize]
                    .checked_sub(current_cost)
                {
                    if delta == 1 {
                        next_pos = next_pos_candidate;
                    }
                }
            }
        }

        current_pos = next_pos;
    }

    saving_at_least_100_pico_count
}

fn find_cheapest_route_costs(
    racetrack: &Vec<Vec<char>>,
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
    let mut best_costs = vec![vec![usize::MAX; racetrack[0].len()]; racetrack.len()];
    while !search_states.is_empty() {
        let state = search_states.pop_front().unwrap();

        let field = racetrack[state.pos.0 as usize][state.pos.1 as usize];
        // we don't set min for walls, so we don't have to worry about them later
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
                || next_pos.0 >= racetrack.len() as i32
                || next_pos.1 < 0
                || next_pos.1 >= racetrack.len() as i32
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

// Essentially the same algorithm as for part one, with the only
// difference that now we will check all points within a 20 hamming
// distance of the current position for a possible cheat.
// We assume that while we are warping, we don't actually have to move
// through wall tiles.
fn part_two(racetrack: &Vec<Vec<char>>) -> usize {
    let mut start_pos = (0, 0);
    let mut target_pos = (0, 0);
    for i in 0..racetrack.len() {
        for j in 0..racetrack[0].len() {
            if racetrack[i][j] == 'S' {
                start_pos = (i as i32, j as i32);
            } else if racetrack[i][j] == 'E' {
                target_pos = (i as i32, j as i32);
            }
        }
    }

    let cheapest_route_costs = find_cheapest_route_costs(&racetrack, start_pos, target_pos);

    let mut saving_at_least_100_pico_count = 0;

    // now we just need trace the route and attempt a cheat in the four cardinal directions
    let mut current_pos = start_pos;
    while current_pos != target_pos {
        let current_cost = cheapest_route_costs[current_pos.0 as usize][current_pos.1 as usize];
        let mut next_pos = current_pos;

        for dir in DIRS {
            let next_pos_candidate = (current_pos.0 + dir.0, current_pos.1 + dir.1);
            if racetrack[next_pos_candidate.0 as usize][next_pos_candidate.1 as usize] != '#' {
                // check for next position
                if let Some(delta) = cheapest_route_costs[next_pos_candidate.0 as usize]
                    [next_pos_candidate.1 as usize]
                    .checked_sub(current_cost)
                {
                    if delta == 1 {
                        next_pos = next_pos_candidate;
                    }
                }
            }
        }

        // we don't care about walls or nothing in this case
        // the only thing that matters is, if we warp in a 20 hamming distance radius
        // do we land on a non-wall and what is our saving
        // a bit crude, but it ads only a constant cost to the algorithm so it doesn't matter
        for delta_i in -20..21 {
            for delta_j in -20..21 {
                let warp_candidate = (current_pos.0 + delta_i, current_pos.1 + delta_j);
                let hamming_distance = (current_pos.0 - warp_candidate.0).abs() as usize
                    + (current_pos.1 - warp_candidate.1).abs() as usize;
                if hamming_distance > 20
                    || warp_candidate.0 < 0
                    || warp_candidate.0 >= racetrack.len() as i32
                    || warp_candidate.1 < 0
                    || warp_candidate.1 >= racetrack.len() as i32
                    || racetrack[warp_candidate.0 as usize][warp_candidate.1 as usize] == '#'
                {
                    continue;
                }

                if let Some(cost_saving) = cheapest_route_costs[warp_candidate.0 as usize]
                    [warp_candidate.1 as usize]
                    .checked_sub(current_cost)
                {
                    // account for the cheat cost of 2 picoseconds
                    if cost_saving > hamming_distance && cost_saving - hamming_distance >= 100 {
                        saving_at_least_100_pico_count += 1;
                    }
                }
            }
        }

        current_pos = next_pos;
    }

    saving_at_least_100_pico_count
}

fn main() {
    let racetrack = parse_input();

    println!(
        "The number of cheats that would save at least a 100 picoseconds: {}",
        part_one(&racetrack)
    );
    println!(
        "The number of cheats that would save at least a 100 picoseconds, with 20 picoseconds of warp time: {}",
        part_two(&racetrack)
    );
}
