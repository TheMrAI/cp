use std::{
    cmp::Ordering,
    collections::{BinaryHeap, HashSet, VecDeque},
    io,
    time::Instant,
    usize, vec,
};

static DIRS: [(i32, i32); 4] = [(0, 1), (1, 0), (0, -1), (-1, 0)];

fn parse_input() -> Vec<Vec<char>> {
    let mut maze = Vec::new();

    for line in io::stdin().lines() {
        let line = line.unwrap();
        maze.push(line.chars().collect());
    }

    maze
}

// state of a search route
#[derive(Eq, PartialEq)]
struct State {
    pub pos: (i32, i32),
    pub dir_i: usize,
    pub cost: usize,
    pub hamming_distance: usize,
}

impl Ord for State {
    fn cmp(&self, other: &Self) -> std::cmp::Ordering {
        if self.hamming_distance == other.hamming_distance {
            // well I know, but I don't want to wrap every element in Reverse
            // to have a min heap
            self.cost.cmp(&other.cost).reverse()
        } else {
            self.hamming_distance.cmp(&other.hamming_distance).reverse()
        }
    }
}

impl PartialOrd for State {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

fn find_cheapest_route_cost(
    maze: &Vec<Vec<char>>,
    start_pos: (i32, i32),
    start_dir_i: usize,
    target_pos: (i32, i32),
) -> Vec<Vec<[usize; 4]>> {
    let mut search_states = BinaryHeap::<State>::new();

    search_states.push(State {
        pos: start_pos,
        dir_i: start_dir_i,
        cost: 0,
        hamming_distance: hamming_distance(start_pos, target_pos),
    });

    // let mut best_cost = usize::MAX;
    // another layer of optimization, so that we can stop if we ever observe a
    // state that we could have already reached for cheaper
    let mut best_costs = vec![vec![[usize::MAX; 4]; maze[0].len()]; maze.len()];

    while !search_states.is_empty() {
        let state = search_states.pop().unwrap();
        // println!("{} {}", state.pos.0, state.pos.1);

        let field = maze[state.pos.0 as usize][state.pos.1 as usize];
        // we don't set min for walls, so we don't have to worry about later
        if field == '#' {
            continue;
        }

        let current_best = best_costs[state.pos.0 as usize][state.pos.1 as usize][state.dir_i];
        if current_best < state.cost {
            continue;
        } else {
            best_costs[state.pos.0 as usize][state.pos.1 as usize][state.dir_i] = state.cost;
        }
        if field == 'E' {
            continue;
        }

        let current_dir = DIRS[state.dir_i];
        for d_i in 0..4 {
            let next_dir = DIRS[d_i];
            let cost = if next_dir == current_dir {
                1
            } else if next_dir == (current_dir.0 * -1, current_dir.1 * -1) {
                2 * 1000 + 1
            } else {
                1000 + 1
            };
            let next_pos = (state.pos.0 + next_dir.0, state.pos.1 + next_dir.1);
            let start_state = State {
                pos: next_pos,
                dir_i: d_i,
                cost: state.cost + cost,
                hamming_distance: hamming_distance(next_pos, target_pos),
            };
            search_states.push(start_state);
        }
    }

    best_costs
}

fn hamming_distance(from: (i32, i32), to: (i32, i32)) -> usize {
    let delta_i = (to.0 - from.0).abs();
    let delta_j = (to.1 - from.1).abs();
    (delta_i + delta_j) as usize
}

fn part_two(
    mut best_costs: Vec<Vec<[usize; 4]>>,
    start_pos: (i32, i32),
    end_pos: (i32, i32),
) -> usize {
    let mut min_path_points = HashSet::<(i32, i32)>::new();

    // rotate all best costs by 2 so all direction information is reversed
    for costs in &mut best_costs {
        for dir_costs in costs {
            // left or right rotate it doesn't matter
            dir_costs.rotate_right(2);
        }
    }
    let mut target_dir = 0;
    for i in 1..4 {
        if best_costs[end_pos.0 as usize][end_pos.1 as usize][i]
            < best_costs[end_pos.0 as usize][end_pos.1 as usize][target_dir]
        {
            target_dir = i;
        }
    }

    // follow all minimal cost paths back to the start
    let mut tracing = VecDeque::<(i32, i32, usize)>::new();
    tracing.push_back((end_pos.0, end_pos.1, target_dir));

    while !tracing.is_empty() {
        let state = tracing.pop_front().unwrap();
        min_path_points.insert((state.0, state.1));

        // original starting position
        if state == (start_pos.0, start_pos.1, 0) {
            continue;
        }

        // while tracing back we always move first, then make a rotation
        // going forward we rotated and stepped, not it is in reverse
        let state_dir = DIRS[state.2];
        let next_pos = (state.0 + state_dir.0, state.1 + state_dir.1);
        for dir_i in 0..4 {
            let dir = DIRS[dir_i];
            // never go backwards
            if dir == (state_dir.0 * -1, state_dir.1 * -1) {
                continue;
            }
            let move_cost = if dir == state_dir { 1 } else { 1000 + 1 };
            if let Some(cost) =
                best_costs[state.0 as usize][state.1 as usize][state.2].checked_sub(move_cost)
            {
                if cost == best_costs[next_pos.0 as usize][next_pos.1 as usize][dir_i] {
                    tracing.push_back((next_pos.0, next_pos.1, dir_i));
                }
            }
        }
    }

    min_path_points.len()
}

fn main() {
    let now = Instant::now();
    let maze = parse_input();

    let mut start_pos = (0i32, 0i32);
    let mut target_pos = (0i32, 0i32);
    for i in 0..maze.len() {
        for j in 0..maze[0].len() {
            if maze[i][j] == 'S' {
                start_pos = (i as i32, j as i32);
            }
            if maze[i][j] == 'E' {
                target_pos = (i as i32, j as i32);
            }
        }
    }

    let best_costs = find_cheapest_route_cost(&maze, start_pos, 0, target_pos);

    println!(
        "The lowest score a Reindeer could possibly get: {}",
        best_costs[target_pos.0 as usize][target_pos.1 as usize]
            .iter()
            .min()
            .unwrap()
    );
    println!(
        "The sum of all box GPS coordinates, in the bigger warehouse: {}",
        part_two(best_costs, start_pos, target_pos)
    );
    println!("Elapsed: {} ms", now.elapsed().as_millis());
}
