use std::collections::{HashMap, HashSet};
use std::io;

fn parse_input() -> Vec<Vec<String>> {
    let mut answers = Vec::<Vec<String>>::new();

    let mut group = Vec::<String>::new();
    for line in io::stdin().lines() {
        let line = line.unwrap();
        if line.len() == 0 {
            answers.push(group);
            group = Vec::<String>::new();
            continue;
        }
        group.push(line);
    }
    if group.len() != 0 {
        answers.push(group);
    }

    return answers;
}

fn part_one(answers: &Vec<Vec<String>>) -> usize {
    answers
        .iter()
        .map(|group| {
            group
                .iter()
                .flat_map(|yeses| yeses.chars())
                .collect::<HashSet<char>>()
                .len()
        })
        .sum()
}

fn part_two(answers: &Vec<Vec<String>>) -> usize {
    answers
        .iter()
        .map(|group| {
            let mut counter = HashMap::<char, usize>::new();
            counter = group.iter().flat_map(|yeses| yeses.chars()).fold(
                counter,
                |mut counter, character| {
                    if let Some(entry) = counter.get_mut(&character) {
                        *entry += 1;
                    } else {
                        counter.insert(character, 1);
                    }
                    counter
                },
            );
            counter
                .into_values()
                .filter(|count| *count == group.len())
                .count()
        })
        .sum()
}

fn main() {
    let answers = parse_input();
    println!(
        "The sum of distinct yes counts per group is: {}",
        part_one(&answers)
    );
    println!(
        "The sum of yes counts per group where everyone said yes: {}",
        part_two(&answers)
    );
}
