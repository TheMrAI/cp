use std::io;

use regex::Regex;

#[derive(Copy, Clone, Debug)]
struct V2 {
    pub x: i64,
    pub y: i64,
}

fn parse_input() -> Vec<[V2; 2]> {
    let re_numbers = Regex::new(r"p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)").unwrap();

    io::stdin()
        .lines()
        .map(|line| {
            let line = line.unwrap();
            re_numbers
                .captures(&line)
                .map(|caps| {
                    let (_, [p_x, p_y, v_x, v_y]) = caps.extract();
                    [
                        V2 {
                            x: p_x.parse().unwrap(),
                            y: p_y.parse().unwrap(),
                        },
                        V2 {
                            x: v_x.parse().unwrap(),
                            y: v_y.parse().unwrap(),
                        },
                    ]
                })
                .unwrap()
        })
        .collect()
}

fn part_one(mut robots: Vec<[V2; 2]>, width: u32, height: u32) -> usize {
    for robot in &mut robots {
        robot[0].x = (robot[0].x + (robot[1].x * 100)) % width as i64;
        // just turning all the potential negative remainders into the positive range
        if robot[0].x < 0 {
            robot[0].x += width as i64;
        }
        robot[0].y = (robot[0].y + (robot[1].y * 100)) % height as i64;
        if robot[0].y < 0 {
            robot[0].y += height as i64;
        }
    }

    let vertical = width / 2;
    let horizontal = height / 2;
    let mut quadrant_count = [0usize; 4];
    for robot in &robots {
        if robot[0].x < vertical as i64 {
            if robot[0].y < horizontal as i64 {
                quadrant_count[0] += 1;
            } else if robot[0].y > horizontal as i64 {
                quadrant_count[1] += 1;
            }
        } else if robot[0].x > vertical as i64 {
            if robot[0].y < horizontal as i64 {
                quadrant_count[2] += 1;
            } else if robot[0].y > horizontal as i64 {
                quadrant_count[3] += 1;
            }
        }
    }

    quadrant_count.iter().fold(1, |lhs, rhs| lhs * rhs)
}

// This was a little tricky as it was difficult to have a good idea for filtering
// states. Tried looking for symmetry, then if a certain band has more blocks etc.
// Had no clue what type of tree I should even be looking for.
// Well, next time remember that the first problems usually hints at the solution
// for the second.
// So here I needed to just use the quadrants again and stop if one of them had more than
// 70% of the robots.
fn part_two(mut robots: Vec<[V2; 2]>, width: u32, height: u32) {
    let mut i = 1;
    loop {
        for robot in &mut robots {
            robot[0].x = (robot[0].x + (robot[1].x)) % width as i64;
            // just turning all the potential negative remainders into the positive range
            if robot[0].x < 0 {
                robot[0].x += width as i64;
            }
            robot[0].y = (robot[0].y + (robot[1].y)) % height as i64;
            if robot[0].y < 0 {
                robot[0].y += height as i64;
            }
        }

        let vertical = width / 2;
        let horizontal = height / 2;
        let mut quadrant_count = [0usize; 4];
        for robot in &robots {
            if robot[0].x < vertical as i64 {
                if robot[0].y < horizontal as i64 {
                    quadrant_count[0] += 1;
                } else if robot[0].y > horizontal as i64 {
                    quadrant_count[1] += 1;
                }
            } else if robot[0].x > vertical as i64 {
                if robot[0].y < horizontal as i64 {
                    quadrant_count[2] += 1;
                } else if robot[0].y > horizontal as i64 {
                    quadrant_count[3] += 1;
                }
            }
        }
        // if any quadrant contains more than 80% of robots let's call it
        if quadrant_count
            .iter()
            .any(|count| *count > (robots.len() * 70) / 100)
        {
            break;
        }
        i += 1;
    }

    println!("{i} seconds\n");
    let mut area = vec![vec![0usize; width as usize]; height as usize];
    for robot in &robots {
        area[robot[0].y as usize][robot[0].x as usize] += 1;
    }
    for i in 0..area.len() {
        for j in 0..area[0].len() {
            if area[i][j] == 0 {
                print!(" ");
            } else {
                print!("{}", area[i][j]);
            }
        }
        println!();
    }
}

fn main() {
    let robots = parse_input();

    // println!(
    //     "Fewest tokens to win all possible prizes: {}",
    //     part_one(robots, 11, 7)
    // );
    println!(
        "Safety factor after a 100 elapsed seconds: {}",
        part_one(robots.clone(), 101, 103)
    );
    println!("The robots arrange themselves into a Christmas tree after:");
    part_two(robots, 101, 103);
}
