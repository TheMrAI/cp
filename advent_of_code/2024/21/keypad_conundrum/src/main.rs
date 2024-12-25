use std::collections::HashMap;
use std::{io, usize};

fn char_to_keypad(c: char) -> (i32, i32) {
    match c {
        '7' => (0, 0),
        '8' => (0, 1),
        '9' => (0, 2),
        '4' => (1, 0),
        '5' => (1, 1),
        '6' => (1, 2),
        '1' => (2, 0),
        '2' => (2, 1),
        '3' => (2, 2),
        ' ' => (3, 0),
        '0' => (3, 1),
        'A' => (3, 2),
        _ => panic!("No such key: {}", c),
    }
}

fn char_to_dir_pad(c: char) -> (i32, i32) {
    match c {
        ' ' => (0, 0),
        '^' => (0, 1),
        'A' => (0, 2),
        '<' => (1, 0),
        'v' => (1, 1),
        '>' => (1, 2),
        _ => panic!("No such key: {}", c),
    }
}

fn num_to_left_right(num: i32) -> char {
    if num < 0 {
        '<'
    } else {
        '>'
    }
}

fn num_to_up_down(num: i32) -> char {
    if num < 0 {
        '^'
    } else {
        'v'
    }
}

fn parse_input() -> Vec<String> {
    io::stdin().lines().map(|line| line.unwrap()).collect()
}

static BANNED_KEYPAD_POS: (i32, i32) = (3, 0);

fn generate_options(
    mut at: (i32, i32),
    target: (i32, i32),
    banned_position: (i32, i32),
) -> (Vec<String>, (i32, i32)) {
    let mut valid_options = Vec::<String>::new();

    let hamming_distance = (target.0 - at.0, target.1 - at.1);
    let horizontal_steps = {
        let char = num_to_left_right(hamming_distance.1);
        std::iter::repeat(char)
            .take(hamming_distance.1.abs() as usize)
            .collect::<String>()
    };
    let vertical_steps = {
        let char = num_to_up_down(hamming_distance.0);
        std::iter::repeat(char)
            .take(hamming_distance.0.abs() as usize)
            .collect::<String>()
    };

    let horizontal_first_midpoint = (at.0, at.1 + hamming_distance.1);
    // the dimension wasn't zero and we don't step on invalid field
    if horizontal_first_midpoint != at && horizontal_first_midpoint != banned_position {
        valid_options.push(horizontal_steps.clone() + &vertical_steps + "A");
    }
    let vertical_first_midpoint = (at.0 + hamming_distance.0, at.1);
    if vertical_first_midpoint != at && vertical_first_midpoint != banned_position {
        valid_options.push(vertical_steps + &horizontal_steps + "A");
    }
    if hamming_distance.0 == 0 && hamming_distance.1 == 0 {
        valid_options.push(String::from("A"));
    }

    at = target;

    (valid_options, at)
}

fn key_code_to_arrow_code(key_code: char, at: (i32, i32)) -> (Vec<String>, (i32, i32)) {
    let target = char_to_keypad(key_code);
    generate_options(at, target, BANNED_KEYPAD_POS)
}

static BANNED_ARROW_POS: (i32, i32) = (0, 0);

fn arrow_code_to_arrow_code(arrow_code: char, at: (i32, i32)) -> (Vec<String>, (i32, i32)) {
    let target = char_to_dir_pad(arrow_code);
    generate_options(at, target, BANNED_ARROW_POS)
}

fn part_one(codes: &Vec<String>) -> usize {
    codes
        .iter()
        .map(|code| {
            let code_val = code[0..code.len() - 1].parse::<usize>().unwrap();

            let encode_1 = encode(code, key_code_to_arrow_code, (3, 2));
            let encode_2 = encode_1
                .into_iter()
                .map(|code| encode(&code, arrow_code_to_arrow_code, (0, 2)))
                .flatten()
                .collect::<Vec<String>>();
            let encode_3 = encode_2
                .into_iter()
                .map(|code| encode(&code, arrow_code_to_arrow_code, (0, 2)))
                .flatten()
                .collect::<Vec<String>>();

            code_val * encode_3.into_iter().map(|entry| entry.len()).min().unwrap()
        })
        .sum()
}

fn encode(
    code: &String,
    encode_with: impl Fn(char, (i32, i32)) -> (Vec<String>, (i32, i32)),
    at: (i32, i32),
) -> Vec<String> {
    let mut encodings = Vec::new();
    encodings.push((String::new(), at));
    for c in code.chars() {
        let mut next_encodings = Vec::new();
        for encoding in &encodings {
            let (possibilities, at) = encode_with(c, encoding.1);
            for possibility in possibilities {
                next_encodings.push((encoding.0.clone() + &possibility, at));
            }
        }
        std::mem::swap(&mut encodings, &mut next_encodings);
    }
    encodings.into_iter().map(|entry| entry.0).collect()
}

fn part_two(codes: &Vec<String>) -> usize {
    codes
        .iter()
        .map(|code| {
            let code_val = code[0..code.len() - 1].parse::<usize>().unwrap();

            let mut minimal_cost = usize::MAX;
            let mut cache = HashMap::new();
            let arrows = encode(code, key_code_to_arrow_code, (3, 2));
            for candidate in arrows {
                let mut cost = 0;
                // Adding a starting "A" so, the starting position is set appropriately
                let candidate = String::from("A") + &candidate;
                for window in candidate.chars().collect::<Vec<char>>().windows(2) {
                    cost += minimal_move_cost_after(24, window[0], window[1], &mut cache);
                }
                minimal_cost = std::cmp::min(minimal_cost, cost);
            }

            code_val * minimal_cost
        })
        .sum()
}

// Depth of zero still means 1 conversion.
// So you need to call it with target depth-1.
fn minimal_move_cost_after(
    depth: usize,
    from: char,
    to: char,
    cache: &mut HashMap<(usize, char, char), usize>,
) -> usize {
    let cache_key = (depth, from, to);
    if let Some(minimal_cost) = cache.get(&cache_key) {
        return *minimal_cost;
    }

    let at = char_to_dir_pad(from);

    let arrows = arrow_code_to_arrow_code(to, at).0;
    let minimal_cost = if depth == 0 {
        arrows.iter().map(|encoding| encoding.len()).min().unwrap()
    } else {
        let mut minimal_cost = usize::MAX;
        for candidate in arrows {
            let mut cost = 0;
            // Adding a starting "A" so, the starting position is set appropriately
            let candidate = String::from("A") + &candidate;
            for window in candidate.chars().collect::<Vec<char>>().windows(2) {
                cost += minimal_move_cost_after(depth - 1, window[0], window[1], cache);
            }
            minimal_cost = std::cmp::min(minimal_cost, cost);
        }
        minimal_cost
    };
    cache.insert(cache_key, minimal_cost);

    minimal_cost
}

fn main() {
    let codes = parse_input();

    println!("The sum of the complexities: {}", part_one(&codes));
    println!(
        "The sum of the complexities, after 25 levels of indirection: {}",
        part_two(&codes)
    );
}

// 3
// ^A
// <A>A
// v<<A>>^AvA^A

// <v<A>>^AvA^A
// <v<A>>^AvA^A

// 37
// ^A ^^<<A
// <A>A <AAv<AA>>^A
// <v<A>>^AvA^A v<<A>>^AAv<A<A>>^AAvAA<^A>A

// <v<A>>^AvA^A <vA<AA>>^AAvA<^A>AAvA^A<vA>^AA<

// ^^<<A => <AAv<AA>>^A
// <<^^A => <vA<AA>>^AA

// v<<A>>^AAv<A<A>>^AAvAA<^A>A
// v<A<AA>>^AAvA<^A>AAvA^A
