use std::collections::{HashMap, HashSet, VecDeque};
use std::io::stdin;
use std::time::Instant;

fn parse_input() -> Vec<u64> {
    stdin()
        .lines()
        .map(|line| line.unwrap().parse().unwrap())
        .collect()
}

fn part_one(numbers: Vec<u64>) -> u64 {
    numbers
        .iter()
        .map(|number| {
            let mut secret_number = *number;
            for _ in 0..2000 {
                secret_number = next_secret_number(secret_number);
            }
            secret_number
        })
        .sum()
}

fn next_secret_number(mut secret_number: u64) -> u64 {
    secret_number = (secret_number * 64) ^ secret_number;
    secret_number = secret_number % 16_777_216;
    secret_number = (secret_number / 32) ^ secret_number;
    secret_number = secret_number % 16_777_216;
    secret_number = (secret_number * 2048) ^ secret_number;
    secret_number = secret_number % 16_777_216;

    secret_number
}

fn part_two(numbers: Vec<u64>) -> u64 {
    let mut sequence_totals = HashMap::new();

    for number in numbers {
        let mut secret_number = number;
        let mut prev_digit = (secret_number % 10) as i8;
        let mut seen = HashSet::new();
        let mut deltas = VecDeque::new();

        for _ in 0..2000 {
            secret_number = next_secret_number(secret_number);
            let curr_digit = (secret_number % 10) as i8;
            deltas.push_back(curr_digit - prev_digit);

            if deltas.len() > 4 {
                deltas.pop_front();
            }
            if deltas.len() == 4 {
                let sequence_key =
                    std::format!("{},{},{},{}", deltas[0], deltas[1], deltas[2], deltas[3]);
                if seen.insert(sequence_key.clone()) {
                    sequence_totals
                        .entry(sequence_key)
                        .and_modify(|total| *total += curr_digit as u64)
                        .or_insert(curr_digit as u64);
                }
            }

            prev_digit = curr_digit;
        }
    }

    *sequence_totals.values().max().unwrap()
}

fn main() {
    let start = Instant::now();
    let numbers = parse_input();

    println!(
        "The sum of the 2000th secret number for each buyer: {}",
        part_one(numbers.clone())
    );
    println!("The most bananas we can get: {}", part_two(numbers));
    println!(
        "Elapsed time: {} ms",
        Instant::now().duration_since(start).as_millis()
    );
}
