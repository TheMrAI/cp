use std::collections::HashMap;
use std::io;

fn parse_input() -> Vec<Vec<char>> {
    let mut data = Vec::new();

    for line in io::stdin().lines() {
        let line = line.unwrap();
        data.push(line.chars().collect());
    }
    data
}

fn count_match(
    data: &Vec<Vec<char>>,
    i: usize,
    j: usize,
    target: &Vec<char>,
    offset_i: isize,
    offset_j: isize,
) -> usize {
    let mut forward_match = true;
    let mut backward_match = true;

    let mut check_i = i as isize;
    let mut check_j = j as isize;
    for index in 0..target.len() {
        if check_i < 0
            || check_i >= data.len() as isize
            || check_j < 0
            || check_j >= data[0].len() as isize
        {
            return 0;
        }

        if data[check_i as usize][check_j as usize] != target[index] {
            forward_match = false;
        }
        if data[check_i as usize][check_j as usize] != target[target.len() - 1 - index] {
            backward_match = false;
        }

        if !forward_match && !backward_match {
            return 0;
        }
        check_i += offset_i;
        check_j += offset_j;
    }

    if forward_match && backward_match {
        2
    } else if forward_match || backward_match {
        1
    } else {
        0
    }
}

fn part_one(data: &Vec<Vec<char>>) -> usize {
    let mut total = 0usize;
    let target = "XMAS".chars().collect();
    for i in 0..data.len() {
        for j in 0..data[0].len() {
            total += count_match(data, i, j, &target, 0, 1)
                + count_match(data, i, j, &target, 1, 0)
                + count_match(data, i, j, &target, 1, 1)
                + count_match(data, i, j, &target, -1, 1);
        }
    }
    total
}

fn part_two(data: &Vec<Vec<char>>) -> usize {
    let mut total = 0usize;
    let target = "MAS".chars().collect();
    for i in 0..data.len() {
        for j in 0..data[0].len()-2 {
            if count_match(data, i, j, &target, 1, 1) > 0
                && count_match(data, i, j+2, &target, 1, -1) > 0
            {
                total += 1;
            }
        }
    }
    total
}

fn main() {
    let data = parse_input();

    println!("Total number of XMASes: {}", part_one(&data));
    println!("Total number of X-MASes: {}", part_two(&data));
}
