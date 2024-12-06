use std::io;

fn parse_input() -> (Vec<(usize, usize)>, Vec<Vec<usize>>) {
    return (parse_orderings(), parse_updates());
}

fn parse_orderings() -> Vec<(usize, usize)> {
    let mut ordering = Vec::new();

    for line in io::stdin().lines() {
        let line = line.unwrap();
        if line.is_empty() {
            break;
        }
        // little hack as all seem to be exactly 2 digit wide
        ordering.push((line[0..2].parse().unwrap(), line[3..].parse().unwrap()));
    }
    ordering
}

fn parse_updates() -> Vec<Vec<usize>> {
    io::stdin()
        .lines()
        .map(|line| {
            let line = line.unwrap();
            line.split(',')
                .map(|number| number.parse().unwrap())
                .collect()
        })
        .collect()
}

// The rules seem to fully specify the order for each pair of elements,
// so we won't bother checking if any of them were broken. Only check if
// we see a number unexpected according to our currently active rules.
fn part_one(ordering_list: &[[bool; 90]; 90], updates: &Vec<Vec<usize>>) -> usize {
    updates
        .iter()
        .filter_map(|update| {
            if valid(ordering_list, update) {
                Some(update[update.len() / 2])
            } else {
                None
            }
        })
        .sum()
}

fn valid(ordering_list: &[[bool; 90]; 90], update: &Vec<usize>) -> bool {
    let mut counter = [0usize; 90];
    for (i, value) in update.iter().enumerate() {
        if counter[value - 10] != i {
            return false;
        }
        for (mark, can_appear_after) in ordering_list[value - 10].iter().enumerate() {
            if *can_appear_after {
                counter[mark] += 1;
            }
        }
    }
    true
}

fn part_two(ordering_list: &[[bool; 90]; 90], updates: &Vec<Vec<usize>>) -> usize {
    updates
        .iter()
        .filter_map(|update| {
            if valid(ordering_list, update) {
                None
            } else {
                Some(middle_in_correct_order(ordering_list, update))
            }
        })
        .sum()
}

// There is no need to create a correct order, just find the one that would be in the middle
// if a correct order was restored.
fn middle_in_correct_order(ordering_list: &[[bool; 90]; 90], update: &Vec<usize>) -> usize {
    let rule_count_for_middle = update.len() / 2;
    for value in update {
        let mut matching_rule_count = 0;
        for rule_for in update {
            if ordering_list[value - 10][rule_for - 10] {
                matching_rule_count += 1;
            }
        }
        if rule_count_for_middle == matching_rule_count {
            return *value;
        }
    }
    0 // this will not happen, as there will be an appropriate value, otherwise we gracefully explode
}

fn main() {
    let (orderings, updates) = parse_input();

    // Didn't want to use maps or sets.
    let mut ordering_list = [[false; 90]; 90];
    for ordering in &orderings {
        ordering_list[ordering.0 - 10][ordering.1 - 10] = true;
    }
    println!(
        "Middle page number sum of correctly ordered updates: {}",
        part_one(&ordering_list, &updates)
    );
    println!(
        "Middle page number sum of incorrectly ordered updates after their correction: {}",
        part_two(&ordering_list, &updates)
    );
}
