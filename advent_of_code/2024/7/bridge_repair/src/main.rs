use std::collections::HashMap;
use std::io;

fn parse_input() -> Vec<(i64, Vec<i64>)> {
    let mut data = Vec::<(i64, Vec<i64>)>::new();
    for line in io::stdin().lines() {
        let line = line.unwrap();
        let colon_at = line.find(':').unwrap();
        let total = line[0..colon_at].parse::<i64>().unwrap();
        let values = line[colon_at + 1..]
            .split_ascii_whitespace()
            .map(|value| value.parse::<i64>().unwrap())
            .collect::<Vec<i64>>();
        data.push((total, values));
    }
    data
}

fn part_one(equations: &Vec<(i64, Vec<i64>)>) -> i64 {
    equations
        .iter()
        .filter_map(|equation| {
            if valid_using_add_mul(&equation.0, 0, &equation.1) {
                Some(equation.0)
            } else {
                None
            }
        })
        .sum()
}

fn valid_using_add_mul(expected: &i64, have: i64, values: &[i64]) -> bool {
    if values.is_empty() {
        return *expected == have;
    }

    if valid_using_add_mul(expected, have + values[0], &values[1..])
        || valid_using_add_mul(expected, have * values[0], &values[1..])
    {
        return true;
    }
    false
}

fn part_two(equations: &Vec<(i64, Vec<i64>)>) -> i64 {
    equations
        .iter()
        .filter_map(|equation| {
            if valid_using_add_mul_concat(&equation.0, 0, &equation.1) {
                Some(equation.0)
            } else {
                None
            }
        })
        .sum()
}

fn valid_using_add_mul_concat(expected: &i64, have: i64, values: &[i64]) -> bool {
    if values.is_empty() {
        return *expected == have;
    }

    if valid_using_add_mul_concat(expected, have + values[0], &values[1..])
        || valid_using_add_mul_concat(expected, have * values[0], &values[1..])
        || valid_using_add_mul_concat(expected, concat_numbers(have, values[0]), &values[1..])
    {
        return true;
    }
    false
}

fn concat_numbers(lhs: i64, rhs: i64) -> i64 {
    let mut tennex = 10;
    while rhs / tennex != 0 {
        tennex *= 10;
    }
    lhs * tennex + rhs
}

fn main() {
    let equations = parse_input();
    println!("{:?}", equations);

    println!(
        "Total calibration result using add, mul: {}",
        part_one(&equations)
    );
    println!(
        "Total calibration result using add, mul, concat: {}",
        part_two(&equations)
    );
}
