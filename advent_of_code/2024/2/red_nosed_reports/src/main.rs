use std::collections::HashMap;
use std::io;

fn parse_input() -> Vec<Vec<usize>> {
    let mut data = Vec::new();

    for line in io::stdin().lines() {
        let levels = line
            .unwrap()
            .split_ascii_whitespace()
            .map(|value| value.parse::<usize>().unwrap())
            .collect();
        data.push(levels);
    }
    return data;
}

fn safe(levels: &[usize]) -> bool {
    if levels.len() < 2 {
        return false;
    }
    let increase = if levels[0] < levels[1] {
        true
    } else if levels[0] > levels[1] {
        false
    } else {
        return false;
    };
    levels
        .windows(2)
        .map(|check| {
            let delta = if increase {
                check[1].checked_sub(check[0])
            } else {
                check[0].checked_sub(check[1])
            };
            match delta {
                Some(magnitude) => {
                    if 1 <= magnitude && magnitude <= 3 {
                        true
                    } else {
                        false
                    }
                }
                None => false,
            }
        })
        .all(|v| v)
}

fn part_one(data: &[Vec<usize>]) -> usize {
    data.iter()
        .map(|levels| if safe(levels) { 1 } else { 0 })
        .sum()
}

// This is the type of problem where it is entirely okay to implement a "bad"
// algorithm. The running time of this is O(n^2), but we don't really care as
// for my inputs no vector is longer then 10 elements. We could likely get an
// O(n) solution but it would get much more complicated. First we would have to
// somehow vote if we are looking at a decreasing or increasing array, because
// this could not be determined from the first two elements (what if the first or second
// one should be removed). Then we would have to go over the array again, find the first
// invalid entry, remove it, then run through the array again to check it.
// Counting the first initial check assuming there is nothing to remove, we would have to
// go through the array 4 times.
// Much simpler to type down a rock simple for loop that removes elements and check if it
// is valid.
fn part_two(data: &[Vec<usize>]) -> usize {
    data.iter()
        .map(|levels| {
            if safe(levels) {
                return 1;
            }
            for i in 0..levels.len() {
                let filtered = levels
                    .iter()
                    .enumerate()
                    .filter_map(|(index, val)| {
                        if index != i {
                            return Some(*val);
                        }
                        None
                    })
                    .collect::<Vec<usize>>();
                if safe(&filtered) {
                    return 1;
                }
            }
            0
        })
        .sum()
}

fn main() {
    let data = parse_input();

    println!("The number of safe reports is: {}", part_one(&data));
    println!(
        "The number of safe reports with a problem dampener is: {}",
        part_two(&data)
    );
}
