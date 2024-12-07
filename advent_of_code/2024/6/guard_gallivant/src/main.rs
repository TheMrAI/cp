use std::{collections::HashSet, io};

fn parse_input() -> Vec<Vec<char>> {
    let mut area = Vec::<Vec<char>>::new();

    for line in io::stdin().lines() {
        let line = line.unwrap();
        area.push(line.chars().collect());
    }
    return area;
}

fn part_one(mut area: Vec<Vec<char>>) -> usize {
    let (mut p_i, mut p_j, mut delta_i) = find_starting_pos_orientation(&area);
    let delta = [(-1, 0), (0, 1), (1, 0), (0, -1)];
    let mut distinct_areas_visited = 1;
    area[p_i as usize][p_j as usize] = 'X';

    loop {
        let next_i = p_i + delta[delta_i].0;
        let next_j = p_j + delta[delta_i].1;

        if next_i < 0 || next_i >= area.len() as i32 || next_j < 0 || next_j >= area[0].len() as i32
        {
            break;
        }
        if area[next_i as usize][next_j as usize] == '#' {
            delta_i = (delta_i + 1) % 4;
        } else {
            p_i = next_i;
            p_j = next_j;
            if area[next_i as usize][next_j as usize] == '.' {
                distinct_areas_visited += 1;
                area[next_i as usize][next_j as usize] = 'X';
            }
        }
    }

    distinct_areas_visited
}

fn find_starting_pos_orientation(area: &Vec<Vec<char>>) -> (i32, i32, usize) {
    let mut p_i = 0i32;
    let mut p_j = 0i32;
    let mut delta_i = 0;

    for i in 0..area.len() {
        for j in 0..area[0].len() {
            if area[i][j] == '^' || area[i][j] == '>' || area[i][j] == 'v' || area[i][j] == '<' {
                p_i = i as i32;
                p_j = j as i32;
                if area[i][j] == '^' {
                    delta_i = 0;
                } else if area[i][j] == '>' {
                    delta_i = 1;
                } else if area[i][j] == 'v' {
                    delta_i = 2;
                } else {
                    delta_i = 3;
                }
            }
        }
    }
    (p_i, p_j, delta_i)
}

fn part_two(mut area: Vec<Vec<char>>) -> usize {
    let (mut p_i, mut p_j, mut delta_i) = find_starting_pos_orientation(&area);
    area[p_i as usize][p_j as usize] = 'X';
    let delta = [(-1, 0), (0, 1), (1, 0), (0, -1)];

    let mut obstacle_count = 0;

    loop {
        let next_i = p_i + delta[delta_i].0;
        let next_j = p_j + delta[delta_i].1;

        if next_i < 0 || next_i >= area.len() as i32 || next_j < 0 || next_j >= area[0].len() as i32
        {
            break;
        }
        if area[next_i as usize][next_j as usize] == '#' {
            delta_i = (delta_i + 1) % 4;
        } else {
            if area[next_i as usize][next_j as usize] != 'X' {
                // put an obstacle there and check
                let mut simulated_area = area.clone();
                simulated_area[next_i as usize][next_j as usize] = '#';
                if ends_in_loop(simulated_area, p_i, p_j, delta_i) {
                    obstacle_count += 1;
                }
            }
            p_i = next_i;
            p_j = next_j;
            area[next_i as usize][next_j as usize] = 'X';
        }
    }

    obstacle_count
}

// end is either detected loop or we fell off the grid
// loop will return true
fn ends_in_loop(area: Vec<Vec<char>>, mut p_i: i32, mut p_j: i32, mut delta_i: usize) -> bool {
    let delta = [(-1, 0), (0, 1), (1, 0), (0, -1)];

    let mut obstacle_collision = HashSet::<(i32, i32, usize)>::new();

    loop {
        let next_i = p_i + delta[delta_i].0;
        let next_j = p_j + delta[delta_i].1;

        // fell off the map
        if next_i < 0 || next_i >= area.len() as i32 || next_j < 0 || next_j >= area[0].len() as i32
        {
            return false;
        }
        if area[next_i as usize][next_j as usize] == '#' {
            if !obstacle_collision.insert((next_i, next_j, delta_i)) {
                return true;
            }
            delta_i = (delta_i + 1) % 4;
        } else {
            p_i = next_i;
            p_j = next_j;
        }
    }
}

fn main() {
    let area = parse_input();

    println!(
        "The number of distinct positions visited: {}",
        part_one(area.clone())
    );
    println!(
        "The number of distinct positions for obstructions: {}",
        part_two(area.clone())
    );
}

// 1618 not right

// ....#.....
// .........#
// ..........
// ..#.......
// .......#..
// ..........
// .#..^.....
// ........#.
// #.........
// ......#...

// ....#.....
// ....+XXX+#
// ....X...X.
// ..#.X...X.
// ....X..#X.
// ....X...X.
// .#.O^XXX+.
// ........#.
// #.........
// ......#...

// ....#.....
// ....+XXX+#
// ....X...X.
// ..#.X...X.
// ..+XXX+#X.
// ..X.X.X.X.
// .#+XXXXX+.
// ......O.#.
// #.........
// ......#...

// ....#.....
// ....+XXX+#
// ....X...X.
// ..#.X...X.
// ..+XXX+#X.
// ..X.X.X.X.
// .#+XXXXX+.
// ......X.#.
// #O+XXX+...
// ......#...

// ....#.....
// ....+XXX+#
// ....X...X.
// ..#.X...X.
// ..+XXX+#X.
// ..X.X.X.X.
// .#+XXXXX+.
// ......X.#.
// #..O+X+...
// ......#...

// ....#.....
// ....+XXX+#
// ....X...X.
// ..#.X...X.
// ..+XXX+#X.
// ..X.X.X.X.
// .#+XXXXX+.
// .+XXXXX+#.
// #+XXXX++..
// ......#O..
