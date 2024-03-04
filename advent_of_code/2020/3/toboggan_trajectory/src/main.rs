use std::io;

fn parse_input() -> Vec<String> {
    io::stdin().lines().map(|l| l.unwrap()).collect()
}

fn part_one(trees: &Vec<String>, slope_x: usize, slope_y: usize) -> usize {
    let mut pos_x = 0;
    let mut pos_y = 0;
    let mut trees_hit = 0;

    while pos_y < trees.len() {
        if trees[pos_y].chars().nth(pos_x).unwrap() == '#' {
            trees_hit += 1;
        }
        pos_x = (pos_x + slope_x) % trees[0].len();
        pos_y += slope_y;
    }

    trees_hit
}

fn main() {
    let trees = parse_input();
    let trees_on_3_1 =  part_one(&trees, 3, 1);
    println!(
        "The number of trees you would hit on a 3 right, 1 down slope is: {}",
        trees_on_3_1
    );
    println!("The number of trees you would hit, multiplied together, if you went down the random slopes: {}", trees_on_3_1 * part_one(&trees, 1, 1) * part_one(&trees, 5, 1) * part_one(&trees, 7, 1) * part_one(&trees, 1, 2));
}
