use std::collections::HashMap;
use std::io;

fn parse_input() -> (Vec<[usize; 5]>, Vec<[usize; 5]>) {
    let mut locks = Vec::new();
    let mut keys = Vec::new();

    let lines = io::stdin()
        .lines()
        .filter_map(|line| {
            let line = line.unwrap();
            if line.is_empty() {
                None
            } else {
                Some(line.chars().collect())
            }
        })
        .collect::<Vec<Vec<char>>>();
    for block in lines.chunks(7) {
        let mut counter = [0; 5];
        for j in 0..5 {
            for i in 1..6 {
                if block[i][j] == '#' {
                    counter[j] += 1;
                }
            }
        }

        if block[0].iter().all(|c| *c == '#') {
            locks.push(counter);
        } else {
            keys.push(counter);
        }
    }

    (locks, keys)
}

fn part_one(locks: &[[usize; 5]], keys: &[[usize; 5]]) -> usize {
    let mut counter = 0;
    for lock in locks {
        for key in keys {
            if lock.iter().zip(key.iter()).all(|(l, k)| l + k <= 5) {
                counter += 1;
            }
        }
    }
    counter
}

fn part_two(lhs: &[usize], rhs: &[usize]) -> usize {
    let mut counter = HashMap::new();
    for value in rhs {
        counter
            .entry(value)
            .and_modify(|count| *count += 1)
            .or_insert(1);
    }
    lhs.iter()
        .map(|number| number * counter.get(number).or(Some(&0)).unwrap())
        .sum()
}

fn main() {
    let (locks, keys) = parse_input();

    println!("{:?}\n{:?}", locks, keys);

    println!(
        "The number of unique lock-key pairs that fit together: {}",
        part_one(&locks, &keys)
    );
    // println!("Similarity score: {}", part_two(&lhs, &rhs));
}
