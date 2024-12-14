use std::collections::HashMap;
use std::io;

// DSU
fn make_set(parents: &mut Vec<usize>, ranks: &mut Vec<usize>, num: usize) {
    parents[num] = num;
    ranks[num] = 0;
}

fn union_sets(parents: &mut Vec<usize>, ranks: &mut Vec<usize>, lhs: usize, rhs: usize) {
    let mut lhs = find_set(parents, lhs);
    let mut rhs = find_set(parents, rhs);

    if lhs != rhs {
        if ranks[lhs] < ranks[rhs] {
            std::mem::swap(&mut lhs, &mut rhs);
        }
        parents[rhs] = lhs;
        if ranks[lhs] == ranks[rhs] {
            ranks[lhs] += 1;
        }
    }
}

fn find_set(parents: &mut Vec<usize>, num: usize) -> usize {
    if num == parents[num] {
        return num;
    }
    parents[num] = find_set(parents, parents[num]);
    parents[num]
}

fn parse_input() -> Vec<Vec<char>> {
    let mut area = Vec::new();

    for line in io::stdin().lines() {
        let line = line.unwrap();
        area.push(line.chars().collect());
    }
    area
}

fn part_one(regions: &Vec<Vec<usize>>) -> usize {
    let dirs = [(0, 1), (1, 0), (0, -1), (-1, 0)];
    let mut fencer = HashMap::<usize, (usize, usize)>::new();
    for i in 0..regions.len() as i32 {
        for j in 0..regions[0].len() as i32 {
            let mut perimeter = 0;
            for dir in &dirs {
                let n_i = i + dir.0;
                let n_j = j + dir.1;
                if n_i < 0
                    || n_i >= regions.len() as i32
                    || n_j < 0
                    || n_j >= regions[0].len() as i32
                {
                    perimeter += 1;
                    continue;
                }
                if regions[i as usize][j as usize] != regions[n_i as usize][n_j as usize] {
                    perimeter += 1;
                }
            }
            fencer
                .entry(regions[i as usize][j as usize])
                .and_modify(|(perim, area)| {
                    *perim += perimeter;
                    *area += 1;
                })
                .or_insert((perimeter, 1));
        }
    }

    fencer
        .values()
        .map(|(perimeter, area)| perimeter * area)
        .sum()
}

// use scan-lines to count the sides for each shape on each axis
fn part_two(regions: &Vec<Vec<usize>>) -> usize {
    let mut fencer = HashMap::<usize, (usize, usize)>::new();

    // fence up-down
    for i in 0..regions.len() {
        let mut fencing_up = false;
        let mut fencing_down = false;
        let mut region = 0; // irrelevant on the first round
        for j in 0..regions[0].len() {
            fencer
                .entry(regions[i][j])
                .and_modify(|(_, area)| *area += 1)
                .or_insert((0, 1));

            if regions[i][j] != region {
                region = regions[i][j];
                // clear the fencing status on region change
                fencing_up = false;
                fencing_down = false;
            }

            if i == 0 || regions[i - 1][j] != region {
                if !fencing_up {
                    fencer.get_mut(&region).unwrap().0 += 1;
                    fencing_up = true;
                }
            } else {
                fencing_up = false;
            }

            if i == regions.len() - 1 || regions[i + 1][j] != region {
                if !fencing_down {
                    fencer.get_mut(&region).unwrap().0 += 1;
                    fencing_down = true;
                }
            } else {
                fencing_down = false;
            }
        }
    }

    // fence left-right
    for j in 0..regions[0].len() {
        let mut fencing_left = false;
        let mut fencing_right = false;
        let mut region = 0; // irrelevant on the first round
        for i in 0..regions.len() {
            if regions[i][j] != region {
                region = regions[i][j];
                // clear the fencing status on region change
                fencing_left = false;
                fencing_right = false;
            }

            if j == 0 || regions[i][j - 1] != region {
                if !fencing_left {
                    fencer.get_mut(&region).unwrap().0 += 1;
                    fencing_left = true;
                }
            } else {
                fencing_left = false;
            }
            if j == regions[0].len() - 1 || regions[i][j + 1] != region {
                if !fencing_right {
                    fencer.get_mut(&region).unwrap().0 += 1;
                    fencing_right = true;
                }
            } else {
                fencing_right = false;
            }
        }
    }

    fencer.values().map(|(side, area)| side * area).sum()
}

fn main() {
    let area = parse_input();

    // identify areas on the map using DSU
    let mut parents = vec![0usize; area.len() * area[0].len()];
    let mut ranks = vec![0usize; area.len() * area[0].len()];
    for i in 0..area.len() {
        for j in 0..area[0].len() {
            let field_id = (i * area[0].len()) + j;
            make_set(&mut parents, &mut ranks, field_id);
        }
    }
    // merge the sets
    let dirs = [(0, 1), (1, 0), (0, -1), (-1, 0)];
    for i in 0..area.len() as i32 {
        for j in 0..area[0].len() as i32 {
            let field_id = (i as usize * area[0].len()) + j as usize;
            for dir in &dirs {
                let n_i = i + dir.0;
                let n_j = j + dir.1;
                if n_i < 0 || n_i >= area.len() as i32 || n_j < 0 || n_j >= area[0].len() as i32 {
                    continue;
                }
                if area[i as usize][j as usize] == area[n_i as usize][n_j as usize] {
                    let neighbor_id = (n_i as usize * area[0].len()) + n_j as usize;
                    union_sets(&mut parents, &mut ranks, field_id, neighbor_id);
                }
            }
        }
    }

    let mut regions = vec![vec![0usize; area[0].len()]; area.len()];
    for i in 0..area.len() {
        for j in 0..area[0].len() {
            let field_id = (i * area[0].len()) + j;
            regions[i][j] = find_set(&mut parents, field_id);
        }
    }

    println!(
        "Total price of fencing for all regions on the map: {}",
        part_one(&regions)
    );
    println!(
        "Total price of fencing for all regions on the map, with bulk pricing: {}",
        part_two(&regions)
    );
}
