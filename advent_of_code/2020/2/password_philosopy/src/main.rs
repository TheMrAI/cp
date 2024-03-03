use regex::Regex;
use std::io;

#[derive(Debug)]
struct Password {
    count_from: usize,
    count_until: usize,
    character: char,
    sequence: String,
}

fn parse_input() -> Vec<Password> {
    let mut data = Vec::new();

    let data_regex = Regex::new(r"(\d+)-(\d+) (\w): (\w+)").unwrap();
    for line in io::stdin().lines() {
        let Some((_, [from, until, char_key, pass])) = data_regex
            .captures(line.as_ref().unwrap())
            .map(|caps| caps.extract())
        else {
            panic!("Failed to parse password from input: '{}'", line.unwrap());
        };
        data.push(Password {
            count_from: from
                .parse()
                .expect(&format!("u32 conversion failed on input: '{}'", from)),
            count_until: until
                .parse()
                .expect(&format!("u32 conversion failed on input: '{}'", until)),
            character: char_key.chars().nth(0).unwrap(),
            sequence: pass.to_string(),
        });
    }
    data
}

fn part_one(passwords: &Vec<Password>) -> usize {
    passwords
        .iter()
        .filter(|password| -> bool {
            let match_count = password
                .sequence
                .chars()
                .filter(|c| *c == password.character)
                .count();
            password.count_from <= match_count && match_count <= password.count_until
        })
        .count()
}

fn part_two(passwords: &Vec<Password>) -> usize {
    passwords
        .iter()
        .filter(|password| -> bool {
            let match_one = password
                .sequence
                .chars()
                .nth(password.count_from - 1)
                .unwrap()
                == password.character;
            let match_two = password
                .sequence
                .chars()
                .nth(password.count_until - 1)
                .unwrap()
                == password.character;
            match_one ^ match_two
        })
        .count()
}

fn main() {
    let passwords = parse_input();

    println!("The number of valid passwords is: {}", part_one(&passwords));
    println!(
        "The number of valid passwords with Toboggan rules: {}",
        part_two(&passwords)
    );
}
