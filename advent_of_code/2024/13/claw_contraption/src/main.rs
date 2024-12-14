use std::io;

use regex::Regex;

#[derive(Copy, Clone, Debug)]
struct V2 {
    pub x: i64,
    pub y: i64,
}

fn parse_input() -> Vec<[V2; 3]> {
    let move_re = Regex::new(r"X(\+\d+), Y(\+\d+)").unwrap();
    let prize_re = Regex::new(r"X=(\d+), Y=(\d+)").unwrap();

    io::stdin()
        .lines()
        .filter_map(|line| match line {
            Ok(data) => {
                if !data.is_empty() {
                    Some(data)
                } else {
                    None
                }
            }
            Err(e) => panic!("{}", e),
        })
        .collect::<Vec<String>>()
        .chunks(3)
        .map(|data| {
            let a = move_re
                .captures(&data[0])
                .map(|caps| {
                    let (_, [x_delta, y_delta]) = caps.extract();
                    V2 {
                        x: x_delta.parse().unwrap(),
                        y: y_delta.parse().unwrap(),
                    }
                })
                .unwrap();
            let b = move_re
                .captures(&data[1])
                .map(|caps| {
                    let (_, [x_delta, y_delta]) = caps.extract();
                    V2 {
                        x: x_delta.parse().unwrap(),
                        y: y_delta.parse().unwrap(),
                    }
                })
                .unwrap();
            let prize = prize_re
                .captures(&data[2])
                .map(|caps| {
                    let (_, [x, y]) = caps.extract();
                    V2 {
                        x: x.parse().unwrap(),
                        y: y.parse().unwrap(),
                    }
                })
                .unwrap();
            [a, b, prize]
        })
        .collect::<Vec<[V2; 3]>>()
}

fn part_one(machines: &Vec<[V2; 3]>) -> usize {
    machines
        .iter()
        .filter_map(|machine| {
            if let Some(solution) = solve(machine) {
                Some(solution)
            } else {
                None
            }
        })
        .sum()
}

// Linear algebra: Ax = v
// x = A^-1 * v
fn solve(machine: &[V2; 3]) -> Option<usize> {
    // A matrix is: machine[0].x, machine[1].x
    //              machine[0].y, machine[1].y
    // For a 2x2 the inverse has a simple formula.
    // From: https://www.mathcentre.ac.uk/resources/uploaded/sigma-matrices7-2009-1.pdf
    let a = machine[0].x;
    let c = machine[0].y;
    let b = machine[1].x;
    let d = machine[1].y;
    let determinant = (a * d) - (b * c);
    if determinant == 0 {
        return None;
    }

    let multiple_x = ((d * machine[2].x) + (-b * machine[2].y)) / determinant;
    let multiple_y = ((-c * machine[2].x) + (a * machine[2].y)) / determinant;

    let check_x = (machine[0].x * multiple_x) + (machine[1].x * multiple_y);
    let check_y = (machine[0].y * multiple_x) + (machine[1].y * multiple_y);
    if check_x == machine[2].x && check_y == machine[2].y {
        return Some(((multiple_x * 3) + multiple_y) as usize);
    }
    None
}

fn part_two(mut machines: Vec<[V2; 3]>) -> usize {
    for machine in &mut machines {
        machine[2].x += 10000000000000;
        machine[2].y += 10000000000000;
    }
    part_one(&machines)
}

fn main() {
    let machines = parse_input();

    println!(
        "Fewest tokens to win all possible prizes: {}",
        part_one(&machines)
    );
    println!(
        "Fewest tokens to win all possible prizes with shifted targets: {}",
        part_two(machines)
    );
}
