use regex::Regex;
use std::collections::HashMap;
use std::io;

enum Op {
    Mul(i32, i32),
    Do,
    Dont,
}

// didn't see any negative numbers in the part one input, but that doesn't
// mean there won't be
fn parse_input() -> Vec<Op> {
    let mut commands = Vec::new();

    let ops = Regex::new(r"mul\((-?\d+),(-?\d+)\)|(do\(\))|(don't\(\))").unwrap();
    for line in io::stdin().lines() {
        let line = line.unwrap();
        for captures in ops.captures_iter(&line) {
            if &captures[0] == "do()" {
                commands.push(Op::Do);
            } else if &captures[0] == "don't()" {
                commands.push(Op::Dont);
            } else {
                let (lhs, rhs) = (&captures[1], &captures[2]);
                commands.push(Op::Mul(
                    lhs.parse::<i32>().unwrap(),
                    rhs.parse::<i32>().unwrap(),
                ));
            }
        }
    }
    return commands;
}

fn part_one(ops: &[Op]) -> i32 {
    ops.iter()
        .filter_map(|op| {
            if let Op::Mul(lhs, rhs) = op {
                Some(lhs * rhs)
            } else {
                None
            }
        })
        .sum()
}

fn part_two(ops: &[Op]) -> i32 {
    let mut enabled = true;
    let mut total = 0i32;
    for op in ops.iter() {
        match op {
            Op::Do => enabled = true,
            Op::Dont => enabled = false,
            Op::Mul(lhs, rhs) => {
                if enabled {
                    total += lhs * rhs;
                }
            }
        }
    }
    total
}

fn main() {
    let ops = parse_input();
    
    println!("Sum of uncorrupted mul operations: {}", part_one(&ops));
    println!("Sum of enabled uncorrupted mul operations: {}", part_two(&ops));
}
