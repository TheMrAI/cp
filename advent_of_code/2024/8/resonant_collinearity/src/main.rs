use std::collections::HashMap;
use std::io;

fn parse_input() -> Vec<Vec<char>> {
    let mut area = Vec::new();

    for line in io::stdin().lines() {
        let line = line.unwrap();
        area.push(line.chars().collect());
    }
    return area;
}

fn part_one(area: &Vec<Vec<char>>) -> usize {
    let mut antinode_counter = Vec::new();
    let mut row = Vec::new();
    row.resize(area[0].len(), '.');
    antinode_counter.resize(area.len(), row);

    let mut signal_map = HashMap::<char, Vec<(i32, i32)>>::new();
    for (i, row) in area.iter().enumerate() {
        for (j, c) in row.iter().enumerate() {
            if *c == '.' {
                continue;
            }
            // mark where an antinode cannot be placed
            signal_map
                .entry(*c)
                .and_modify(|locations| locations.push((i as i32, j as i32)))
                .or_insert(vec![(i as i32, j as i32)]);
        }
    }

    // generate antinodes
    for towers in signal_map.values() {
        for i in 0..towers.len() - 1 {
            for j in i + 1..towers.len() {
                let from = towers[i];
                let to = towers[j];
                let delta = (to.0 - from.0, to.1 - from.1);
                let antinodes = vec![
                    (from.0 - delta.0, from.1 - delta.1),
                    (to.0 + delta.0, to.1 + delta.1),
                ];
                for antinode in antinodes {
                    if antinode.0 < 0
                        || antinode.0 >= area.len() as i32
                        || antinode.1 < 0
                        || antinode.1 >= area[0].len() as i32
                    {
                        continue;
                    }

                    antinode_counter[antinode.0 as usize][antinode.1 as usize] = '#';
                }
            }
        }
    }

    antinode_counter
        .iter()
        .map(|row| row.iter().filter(|marker| **marker == '#').count())
        .sum()
}

fn part_two(area: &Vec<Vec<char>>) -> usize {
    let mut antinode_counter = Vec::new();
    let mut row = Vec::new();
    row.resize(area[0].len(), '.');
    antinode_counter.resize(area.len(), row);

    let mut signal_map = HashMap::<char, Vec<(i32, i32)>>::new();
    for (i, row) in area.iter().enumerate() {
        for (j, c) in row.iter().enumerate() {
            if *c == '.' {
                continue;
            }
            // mark where an antinode cannot be placed
            signal_map
                .entry(*c)
                .and_modify(|locations| locations.push((i as i32, j as i32)))
                .or_insert(vec![(i as i32, j as i32)]);
        }
    }

    // generate antinodes
    for towers in signal_map.values() {
        for i in 0..towers.len() - 1 {
            for j in i + 1..towers.len() {
                let from = towers[i];
                let to = towers[j];
                let delta = (to.0 - from.0, to.1 - from.1);

                let mut curr = from;
                while curr.0 >= 0
                    && curr.0 < area.len() as i32
                    && curr.1 >= 0
                    && curr.1 < area[0].len() as i32
                {
                    antinode_counter[curr.0 as usize][curr.1 as usize] = '#';
                    curr = (curr.0 - delta.0, curr.1 - delta.1);
                }
                curr = to;
                while curr.0 >= 0
                    && curr.0 < area.len() as i32
                    && curr.1 >= 0
                    && curr.1 < area[0].len() as i32
                {
                    antinode_counter[curr.0 as usize][curr.1 as usize] = '#';
                    curr = (curr.0 + delta.0, curr.1 + delta.1);
                }
            }
        }
    }

    antinode_counter
        .iter()
        .map(|row| row.iter().filter(|marker| **marker == '#').count())
        .sum()
}

fn main() {
    let area = parse_input();

    println!(
        "The number of unique antinode locations: {}",
        part_one(&area)
    );
    println!(
        "The number of unique antinode locations, with harmonics: {}",
        part_two(&area)
    );
}
