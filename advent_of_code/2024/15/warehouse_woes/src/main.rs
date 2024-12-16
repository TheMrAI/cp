use std::{collections::HashSet, io, vec};

fn parse_input() -> (Vec<Vec<char>>, Vec<char>) {
    let mut warehouse = Vec::new();
    let mut movements = Vec::new();

    let mut warehouse_parsed = false;
    for line in io::stdin().lines() {
        let line = line.unwrap();
        if line.is_empty() {
            warehouse_parsed = true;
        }

        if !warehouse_parsed {
            warehouse.push(line.chars().collect());
        } else {
            movements.extend(line.chars());
        }
    }

    (warehouse, movements)
}

fn part_one(mut warehouse: Vec<Vec<char>>, movements: &Vec<char>) -> usize {
    let mut robot_at = (0, 0);
    for i in 0..warehouse.len() {
        for j in 0..warehouse[0].len() {
            if warehouse[i][j] == '@' {
                robot_at.0 = i as i32;
                robot_at.1 = j as i32;
            }
        }
    }

    for step in movements {
        let offset = step_to_offset(*step);
        robot_at = single_line_move(&mut warehouse, offset, robot_at);
    }

    let mut total_gps = 0;
    for i in 0..warehouse.len() {
        for j in 0..warehouse[0].len() {
            if warehouse[i][j] == 'O' {
                total_gps += i * 100 + j;
            }
        }
    }

    total_gps
}

fn single_line_move(
    warehouse: &mut Vec<Vec<char>>,
    mut offset: (i32, i32),
    robot_at: (i32, i32),
) -> (i32, i32) {
    let mut move_chain_last = robot_at;
    while !(warehouse[move_chain_last.0 as usize][move_chain_last.1 as usize] == '.'
        || warehouse[move_chain_last.0 as usize][move_chain_last.1 as usize] == '#')
    {
        move_chain_last = (move_chain_last.0 + offset.0, move_chain_last.1 + offset.1);
    }

    if warehouse[move_chain_last.0 as usize][move_chain_last.1 as usize] == '#' {
        return robot_at;
    }
    let next_robot_pos = (robot_at.0 + offset.0, robot_at.1 + offset.1);

    offset = (offset.0 * -1, offset.1 * -1);
    // move the chain
    while move_chain_last != robot_at {
        let move_chain_next = (move_chain_last.0 + offset.0, move_chain_last.1 + offset.1);
        // a simple swap, but we can't do std::mem::swap because of borrow checker
        let tmp = warehouse[move_chain_last.0 as usize][move_chain_last.1 as usize];
        warehouse[move_chain_last.0 as usize][move_chain_last.1 as usize] =
            warehouse[move_chain_next.0 as usize][move_chain_next.1 as usize];
        warehouse[move_chain_next.0 as usize][move_chain_next.1 as usize] = tmp;
        move_chain_last = move_chain_next;
    }
    next_robot_pos
}

fn step_to_offset(step: char) -> (i32, i32) {
    match step {
        '^' => (-1, 0),
        '>' => (0, 1),
        'v' => (1, 0),
        '<' => (0, -1),
        _ => panic!("not happening"),
    }
}

fn part_two(mut warehouse: Vec<Vec<char>>, movements: &Vec<char>) -> usize {
    let mut robot_at = (0i32, 0i32);
    warehouse = warehouse
        .iter()
        .enumerate()
        .map(|(i, line)| {
            let mut expanded_line = vec!['.'; line.len() * 2];
            for j in 0..line.len() {
                if line[j] == '#' {
                    expanded_line[2 * j] = '#';
                    expanded_line[2 * j + 1] = '#';
                } else if line[j] == '.' {
                    expanded_line[2 * j] = '.';
                    expanded_line[2 * j + 1] = '.';
                } else if line[j] == 'O' {
                    expanded_line[2 * j] = '[';
                    expanded_line[2 * j + 1] = ']';
                } else {
                    expanded_line[2 * j] = '@';
                    expanded_line[2 * j + 1] = '.';
                    robot_at = (i as i32, 2 * j as i32);
                }
            }
            expanded_line
        })
        .collect();

    warehouse = simulate_cascade_movement(warehouse, movements, robot_at);

    let mut total_gps = 0;
    for i in 0..warehouse.len() {
        for j in 0..warehouse[0].len() {
            if warehouse[i][j] == '[' {
                total_gps += i * 100 + j;
            }
        }
    }

    total_gps
}

fn simulate_cascade_movement(
    mut warehouse: Vec<Vec<char>>,
    movements: &Vec<char>,
    mut robot_at: (i32, i32),
) -> Vec<Vec<char>> {
    for step in movements {
        let offset = step_to_offset(*step);

        // left - right
        robot_at = if offset == (0, -1) || offset == (0, 1) {
            single_line_move(&mut warehouse, offset, robot_at)
        } else {
            // up - down
            cascade_move(&mut warehouse, offset, robot_at)
        };
    }

    warehouse
}

// only handles up-down movement
fn cascade_move(
    warehouse: &mut Vec<Vec<char>>,
    offset: (i32, i32),
    robot_at: (i32, i32),
) -> (i32, i32) {
    debug_assert!(offset == (-1, 0) || offset == (1, 0));

    let mut move_stack = Vec::new();
    let mut moves = HashSet::new();
    moves.insert(robot_at);
    move_stack.push(moves);

    while !(move_stack.last().unwrap().iter().any(|(i, j)| {
        let entry = warehouse[*i as usize][*j as usize];
        entry == '#'
    }) || (move_stack.last().unwrap().iter().all(|(i, j)| {
        let entry = warehouse[*i as usize][*j as usize];
        entry == '.'
    }))) {
        let mut next_moves = HashSet::new();
        move_stack
            .last()
            .unwrap()
            .iter()
            .filter(|(i, j)| warehouse[*i as usize][*j as usize] != '.')
            .for_each(|(i, j)| {
                let i = i + offset.0;
                let j = j + offset.1;

                // Would have preferred to avoid using hashset here, but then there would have been the issue
                // of tracking for exactly which piece have we prepared a move, which can get complicated.
                if warehouse[i as usize][j as usize] == '[' {
                    next_moves.insert((i, j));
                    next_moves.insert((i, j + 1));
                } else if warehouse[i as usize][j as usize] == ']' {
                    next_moves.insert((i, j));
                    next_moves.insert((i, j - 1));
                } else {
                    next_moves.insert((i, j));
                }
            });
        move_stack.push(next_moves);
    }

    if move_stack.last().unwrap().iter().any(|(i, j)| {
        let entry = warehouse[*i as usize][*j as usize];
        entry == '#'
    }) {
        return robot_at;
    }
    let next_robot_pos = (robot_at.0 + offset.0, robot_at.1 + offset.1);

    // move the chain
    while !move_stack.is_empty() {
        let moves = move_stack
            .pop()
            .unwrap()
            .into_iter()
            .filter(|(i, j)| warehouse[*i as usize][*j as usize] != '.')
            .collect::<Vec<(i32, i32)>>();
        // Contrary to part one here we push the change instead of pulling it, otherwise we
        // could potentially pull apart other blocks.
        for movement in moves.iter() {
            // a simple swap, but we can't do std::mem::swap because of borrow checker
            let i = movement.0 as usize;
            let j = movement.1 as usize;
            let n_i = (movement.0 + offset.0) as usize;
            let n_j = (movement.1 + offset.1) as usize;
            let tmp = warehouse[i][j];
            warehouse[i][j] = warehouse[n_i][n_j];
            warehouse[n_i][n_j] = tmp;
        }
    }
    next_robot_pos
}

fn main() {
    let (warehouse, movements) = parse_input();

    println!(
        "The sum of all box GPS coordinates: {}",
        part_one(warehouse.clone(), &movements)
    );
    println!(
        "The sum of all box GPS coordinates, in the bigger warehouse: {}",
        part_two(warehouse, &movements)
    );
}
