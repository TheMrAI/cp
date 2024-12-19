use std::io;

fn parse_input() -> (Vec<String>, Vec<String>) {
    let mut lines = io::stdin().lines();

    let towels = lines
        .next()
        .unwrap()
        .unwrap()
        .split(',')
        .map(|towel| towel.trim().to_string())
        .collect();
    let patterns = lines
        .filter_map(|line| {
            let line = line.unwrap();
            if line.is_empty() {
                None
            } else {
                Some(line)
            }
        })
        .collect();

    return (towels, patterns);
}

fn part_one(towels: &Vec<Vec<Vec<char>>>, patterns: &Vec<Vec<char>>) -> usize {
    patterns
        .iter()
        .filter(|pattern| {
            let mut possibility_cache = vec![true; pattern.len()];
            is_possible(pattern, towels, 0, &mut possibility_cache)
        })
        .count()
}

fn is_possible(
    pattern: &Vec<char>,
    towels: &Vec<Vec<Vec<char>>>,
    position: usize,
    possibility_cache: &mut Vec<bool>,
) -> bool {
    if position == pattern.len() {
        return true;
    }
    let cache_key = encode(pattern[position]);
    if !possibility_cache[position] {
        return false;
    }

    for towel in &towels[cache_key] {
        let next_pos = position + towel.len();
        if next_pos > pattern.len() {
            // don't even bother with the check if the patterns are longer then remaining space
            break;
        }
        if *towel == pattern[position..next_pos] {
            let ok = is_possible(pattern, towels, next_pos, possibility_cache);
            if ok {
                return true;
            }
        }
    }
    possibility_cache[position] = false;

    false
}

fn part_two(towels: &Vec<Vec<Vec<char>>>, patterns: &Vec<Vec<char>>) -> usize {
    patterns
        .iter()
        .map(|pattern| {
            // plus one for our convenience
            // for the possibility_cache it has no real purpose we just don't over index
            // for the counter it helps conceptually to say that if we have reached pass the end of the
            // final character position that we can match nothing one way, then the DP resolves itself
            let mut possibility_cache = vec![true; pattern.len() + 1];
            let mut counter = vec![0; pattern.len() + 1];
            count_possibilities(pattern, towels, 0, &mut possibility_cache, &mut counter);
            counter[0]
        })
        .sum()
}

fn count_possibilities(
    pattern: &Vec<char>,
    towels: &Vec<Vec<Vec<char>>>,
    position: usize,
    possibility_cache: &mut Vec<bool>,
    counter: &mut Vec<usize>,
) -> bool {
    if position == pattern.len() {
        counter[position] = 1;
        return true;
    }
    let cache_key = encode(pattern[position]);
    if !possibility_cache[position] {
        return false;
    }

    let mut ok = false;
    for towel in &towels[cache_key] {
        let next_pos = position + towel.len();
        if next_pos > pattern.len() {
            // don't even bother with the check if the patterns are longer then remaining space
            break;
        }
        if *towel == pattern[position..next_pos] {
            if !(possibility_cache[next_pos] && counter[next_pos] != 0) {
                ok |= count_possibilities(pattern, towels, next_pos, possibility_cache, counter);
            }
            counter[position] += counter[next_pos];
        }
    }
    possibility_cache[position] = ok;

    ok
}

// white (w), blue (u), black (b), red (r), or green (g)
// we are going to encode the char symbols for easier handling
fn encode(c: char) -> usize {
    match c {
        'w' => 0,
        'u' => 1,
        'b' => 2,
        'r' => 3,
        'g' => 4,
        _ => panic!("Unexpected character: {}", c), // shouldn't happen
    }
}

fn main() {
    let (mut towels, patterns) = parse_input();

    // sort them by length, so we can exit early later
    towels.sort_unstable_by(|lhs, rhs| lhs.len().cmp(&rhs.len()));

    let mut char_towels = vec![vec![]; 5];
    for towel in towels {
        let towel_chars = towel.chars().collect::<Vec<char>>();
        let code = encode(towel_chars[0]);
        char_towels[code].push(towel_chars);
    }

    let char_patterns = patterns
        .iter()
        .map(|pattern| pattern.chars().collect())
        .collect::<Vec<Vec<char>>>();

    println!(
        "Number of possible designs: {}",
        part_one(&char_towels, &char_patterns)
    );
    println!(
        "Sum of all possible design patterns: {}",
        part_two(&char_towels, &char_patterns)
    );
}
