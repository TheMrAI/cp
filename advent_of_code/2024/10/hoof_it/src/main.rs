use std::collections::HashSet;
use std::io;

fn parse_input() -> Vec<Vec<u32>> {
    io::stdin()
        .lines()
        .map(|line| {
            line.unwrap()
                .chars()
                .map(|c| c.to_digit(10).unwrap())
                .collect()
        })
        .collect()
}

fn part_one(area: &Vec<Vec<u32>>) -> usize {
    let mut total_score = 0;

    for i in 0..area.len() {
        for j in 0..area[0].len() {
            if area[i][j] == 0 {
                total_score += collect_trail_ends(area, i, j)
                    .into_iter()
                    .collect::<HashSet<(i32, i32)>>()
                    .len();
            }
        }
    }

    total_score
}

fn part_two(area: &Vec<Vec<u32>>) -> usize {
    let mut total_score = 0;

    for i in 0..area.len() {
        for j in 0..area[0].len() {
            if area[i][j] == 0 {
                total_score += collect_trail_ends(area, i, j).len();
            }
        }
    }

    total_score
}

fn collect_trail_ends(area: &Vec<Vec<u32>>, i: usize, j: usize) -> Vec<(i32, i32)> {
    let mut trail_ends = Vec::new();

    let dirs = [(0, 1), (1, 0), (0, -1), (-1, 0)];

    let mut trails = Vec::new();
    trails.push((i as i32, j as i32));
    while !trails.is_empty() {
        let mut next_trails = Vec::new();

        for pos in &trails {
            for dir in dirs {
                let n_i = pos.0 + dir.0;
                let n_j = pos.1 + dir.1;

                if n_i < 0 || n_i >= area.len() as i32 || n_j < 0 || n_j >= area[0].len() as i32 {
                    continue;
                }

                if area[pos.0 as usize][pos.1 as usize] + 1 == area[n_i as usize][n_j as usize] {
                    if area[n_i as usize][n_j as usize] == 9 {
                        trail_ends.push((n_i, n_j));
                        continue;
                    }
                    next_trails.push((n_i, n_j));
                }
            }
        }

        std::mem::swap(&mut trails, &mut next_trails);
    }

    trail_ends
}

fn main() {
    let area = parse_input();

    println!("The sum of all trailhead scores: {}", part_one(&area));
    println!("The sum of all trailhead ratings: {}", part_two(&area));
}
