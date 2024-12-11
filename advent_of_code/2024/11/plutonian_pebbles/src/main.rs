use std::{collections::HashMap, io};

fn parse_input() -> Vec<usize> {
    let mut line = String::new();
    io::stdin().read_line(&mut line).unwrap();
    line.split_ascii_whitespace()
        .map(|number| number.parse::<usize>().unwrap())
        .collect()
}

fn part_one(mut numbers: Vec<usize>) -> usize {
    for _ in 0..25 {
        let mut new_numbers = Vec::new();
        for num in &numbers {
            new_numbers.extend(blink(*num));
        }
        std::mem::swap(&mut new_numbers, &mut numbers);
    }
    numbers.len()
}

fn blink(num: usize) -> Vec<usize> {
    let mut new_numbers = Vec::new();
    let number_of_digits = number_of_digits(num);
    if num == 0 {
        new_numbers.push(1);
    } else if number_of_digits % 2 == 0 {
        let power_of_ten = 10usize.pow((number_of_digits / 2) as u32);
        let upper = num / power_of_ten;
        new_numbers.push(upper);
        new_numbers.push(num - (upper * power_of_ten));
    } else {
        new_numbers.push(num * 2024);
    }
    new_numbers
}

fn number_of_digits(num: usize) -> usize {
    let mut digit_count = 1;
    while num >= 10usize.pow(digit_count) {
        digit_count += 1;
    }
    digit_count as usize
}

fn part_two(numbers: &Vec<usize>) -> usize {
    let mut five_blink_cache = HashMap::<usize, Vec<usize>>::new();

    let mut working_set = HashMap::<usize, usize>::new();
    for num in numbers {
        working_set.insert(*num, 1);
    }
    let mut blink_count = 0;
    while blink_count != 75 {
        let mut next_working_set = HashMap::<usize, usize>::new();
        for (num, count) in &working_set {
            let in_five = five_blink_cache.entry(*num).or_insert_with(|| {
                let mut numbers = Vec::new();
                numbers.push(*num);
                for _ in 0..5 {
                    let mut new_numbers = Vec::new();
                    for num in &numbers {
                        new_numbers.extend(blink(*num));
                    }
                    std::mem::swap(&mut new_numbers, &mut numbers);
                }
                numbers
            });
            for value in in_five {
                next_working_set
                    .entry(*value)
                    .and_modify(|counter| *counter += count)
                    .or_insert(*count);
            }
        }
        std::mem::swap(&mut working_set, &mut next_working_set);
        blink_count += 5;
    }

    working_set.values().sum()
}

fn main() {
    let numbers = parse_input();

    println!(
        "The number of stones after blinking 25 times: {}",
        part_one(numbers.clone())
    );
    println!(
        "The number of stones after blinking 75 times: {}",
        part_two(&numbers)
    );
}
