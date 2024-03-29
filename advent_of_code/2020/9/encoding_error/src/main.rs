use std::collections::HashSet;

fn parse_input(source: std::io::Stdin) -> Vec<i64> {
    let mut numbers = Vec::new();

    for line in source.lines() {
        let line = line.unwrap();
        let value: i64 = match line.parse() {
            Ok(v) => v,
            Err(_) => panic!("Failed to parse integer from input: '{}'", line),
        };
        numbers.push(value);
    }

    numbers
}

fn part_one(numbers: &Vec<i64>, window: usize) -> Option<i64> {
    let mut distinct_numbers = HashSet::<i64>::new();
    for preamble in numbers.iter().take(window) {
        distinct_numbers.insert(*preamble);
    }

    let non_constructable = numbers.iter().skip(window).enumerate().find(|&entry| {
        let cant_construct = distinct_numbers
            .iter()
            .find(|&value| distinct_numbers.contains(&(entry.1 - *value)))
            .is_none();
        distinct_numbers.remove(&numbers[entry.0]);
        distinct_numbers.insert(*entry.1);
        cant_construct
    });
    match non_constructable {
        Some(data) => Some(*data.1),
        None => None,
    }
}

fn part_two(numbers: &Vec<i64>, invalid_number: i64) -> Option<i64> {
    let mut window_width = 0;
    let mut sum = 0;
    let mut match_sequence = None;
    'labelloopdyloop: for (index, num) in numbers.iter().enumerate() {
        sum += num;
        window_width += 1;
        if sum == invalid_number {
            match_sequence = Some(&numbers[(1 + index - window_width)..index + 1]);
            break 'labelloopdyloop;
        }
        while sum > invalid_number && window_width > 0 {
            sum -= numbers[1 + index - window_width];
            window_width -= 1;
            if sum == invalid_number {
                match_sequence = Some(&numbers[(1 + index - window_width)..index + 1]);
                break 'labelloopdyloop;
            }
        }
    }

    match match_sequence {
        Some(sequence) => {
            let min = sequence.into_iter().min().unwrap();
            let max = sequence.into_iter().max().unwrap();
            Some(*min + *max)
        }
        None => None,
    }
}

fn main() {
    let numbers = parse_input(std::io::stdin());
    let invalid_number = part_one(&numbers, 25).unwrap();
    println!(
        "The first number that can't be built from the prior 25 numbers is: {}",
        invalid_number
    );
    println!(
        "The encryption weakness is: {}",
        part_two(&numbers, invalid_number).unwrap()
    );
}
