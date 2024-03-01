use std::collections::{HashMap, HashSet};
use std::io;

fn parse_input() -> Vec<u32> {
    let mut data = Vec::new();

    for line in io::stdin().lines() {
        let value: u32 = match line.as_ref().unwrap().parse() {
            Ok(v) => v,
            Err(_) => panic!("u32 parsing failed on input: '{}'", line.unwrap()),
        };
        data.push(value);
    }
    return data;
}

fn part_one(numbers: &Vec<u32>, sum_to: u32) -> u32 {
    let mut seeking_remainder = HashSet::new();
    for number in numbers.iter() {
        let seeking = sum_to - number;
        if seeking_remainder.get(&seeking).is_some() {
            return number * seeking;
        } else {
            seeking_remainder.insert(number);
        }
    }
    return 0;
}

fn part_two(numbers: &Vec<u32>, sum_to: u32) -> u32 {
    let mut seeking_remainder = HashMap::new();
    for (i, number) in numbers.iter().enumerate() {
        if i == 0 {
            continue;
        }
        let seeking = sum_to - number;
        if let Some(stored_multiple) = seeking_remainder.get(&seeking) {
            return number * stored_multiple;
        } else {
            // insert all possible pair sums using current number
            for j in 0..i {
                seeking_remainder.insert(number + numbers[j], number * numbers[j]);
            }
        }
    }
    return 0;
}

fn main() {
    let data = parse_input();

    println!(
        "The two entries that sum up to 2020 multiplied together give: {}",
        part_one(&data, 2020)
    );
    println!(
        "The three entries that sum up to 2020 multiplied together give: {}",
        part_two(&data, 2020)
    );
}
