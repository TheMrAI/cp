use std::collections::HashMap;
use std::io;

fn parse_input() -> (Vec<usize>, Vec<usize>) {
    let mut lhs = Vec::new();
    let mut rhs = Vec::new();

    for line in io::stdin().lines() {
        let line = line.unwrap();
        let parts = line.split_ascii_whitespace().collect::<Vec<&str>>();

        lhs.push(parts[0].parse::<usize>().unwrap());
        rhs.push(parts[1].parse::<usize>().unwrap());
    }
    return (lhs, rhs);
}

fn part_one(lhs: &[usize], rhs: &[usize]) -> usize {
    lhs.iter()
        .zip(rhs.iter())
        .map(|(lhs, rhs)| if lhs > rhs { lhs - rhs } else { rhs - lhs })
        .sum()
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
    let (mut lhs, mut rhs) = parse_input();
    lhs.sort_unstable();
    rhs.sort_unstable();

    println!("Total distance between two lists: {}", part_one(&lhs, &rhs));
    println!("Similarity score: {}", part_two(&lhs, &rhs));
}
