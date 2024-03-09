use core::panic;
use regex::Regex;
use std::collections::{HashMap, HashSet};

fn parse_input(source: std::io::Stdin) -> HashMap<String, HashMap<String, u32>> {
    let mut bag_tree = HashMap::<String, HashMap<String, u32>>::new();

    let bag_name_regexp = Regex::new(r"(\D+) bags").unwrap();
    let bag_counts = Regex::new(r"(\d+) (\D+) bags?[,.]{1}").unwrap();
    for line in source.lines() {
        let line = line.unwrap();
        let parts: Vec<&str> = line.split(" contain ").collect();
        if parts.len() != 2 {
            panic!("Unexpected input received. Got: '{}'", line);
        }
        let Some((_, [bag_name])) = bag_name_regexp
            .captures(parts[0])
            .map(|caps| caps.extract())
        else {
            panic!("Failed to parse bag name from input: '{}'", parts[0]);
        };
        if parts[1] == "no other bags." {
            bag_tree.insert(bag_name.into(), HashMap::<String, u32>::new());
            continue;
        }

        let mut bags = HashMap::<String, u32>::new();
        for (_, [count_str, bag_type]) in bag_counts
            .captures_iter(parts[1])
            .map(|caps| caps.extract())
        {
            bags.insert(bag_type.into(), count_str.parse().unwrap());
        }
        bag_tree.insert(bag_name.into(), bags);
    }

    return bag_tree;
}

fn part_one(bag_tree: &HashMap<String, HashMap<String, u32>>, target_bag_name: &str) -> usize {
    let mut to_check = bag_tree
        .iter()
        .filter(|entry| entry.1.contains_key(target_bag_name))
        .map(|entry| entry.0.clone())
        .collect::<HashSet<String>>();

    let mut bags_containing_target = HashSet::<String>::new();
    while to_check.len() != 0 {
        let mut check_next = HashSet::<String>::new();
        for checking in to_check.into_iter() {
            if !bags_containing_target.contains(&checking) {
                check_next = bag_tree
                    .iter()
                    .filter(|entry| entry.1.contains_key(&checking))
                    .fold(check_next, |mut acc, entry| {
                        acc.insert(entry.0.clone());
                        return acc;
                    });
            }
            bags_containing_target.insert(checking);
        }
        to_check = check_next;
    }

    bags_containing_target.len()
}

fn part_two(bag_tree: &HashMap<String, HashMap<String, u32>>, target_bag_name: &str) -> usize {
    number_of_bags_with(bag_tree, target_bag_name) - 1
}

fn number_of_bags_with(
    bag_tree: &HashMap<String, HashMap<String, u32>>,
    target_bag_name: &str,
) -> usize {
    bag_tree
        .get(target_bag_name)
        .unwrap()
        .iter()
        .map(|(contained_bag, times)| {
            (*times as usize) * number_of_bags_with(bag_tree, contained_bag)
        })
        .sum::<usize>()
        + 1
}

fn main() {
    let bag_tree = parse_input(std::io::stdin());
    println!(
        "The number of distinct bags that eventually contain a shiny gold bag: {}",
        part_one(&bag_tree, "shiny gold")
    );
    println!(
        "The number of individual bags inside a shiny gold bag is: {}",
        part_two(&bag_tree, "shiny gold")
    );
}
