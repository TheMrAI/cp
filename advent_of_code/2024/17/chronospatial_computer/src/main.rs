#[derive(Debug, Clone)]
struct Machine {
    pub reg_a: i32,
    pub reg_b: i32,
    pub reg_c: i32,
    pub ip: usize,
    pub instructions: Vec<i32>,
}

impl Machine {
    pub fn combo(&self, operand: i32) -> i32 {
        match operand {
            0 | 1 | 2 | 3 => operand,
            4 => self.reg_a,
            5 => self.reg_b,
            6 => self.reg_c,
            _ => panic!("Unexpected combo operand: {}", operand),
        }
    }

    pub fn simulate(&mut self) -> Option<i32> {
        while self.ip < self.instructions.len() {
            let instruction = self.instructions[self.ip];
            let operand = self.instructions[self.ip + 1];
            self.ip += 2;

            match instruction {
                0 => self.reg_a = self.reg_a / 2i32.pow(self.combo(operand) as u32),
                1 => self.reg_b = self.reg_b ^ operand,
                2 => self.reg_b = self.combo(operand) % 8,
                3 => {
                    if self.reg_a != 0 {
                        self.ip = operand as usize;
                        continue;
                    }
                }
                4 => self.reg_b = self.reg_b ^ self.reg_c,
                5 => {
                    return Some(self.combo(operand) % 8);
                }
                6 => self.reg_b = self.reg_a / 2i32.pow(self.combo(operand) as u32),
                7 => self.reg_c = self.reg_a / 2i32.pow(self.combo(operand) as u32),
                _ => panic!("Unexpected instruction: {}", instruction),
            }
        }
        None
    }
}

// reading the input is by far the worst part of AOC
fn parse_input() -> Machine {
    let processed = std::io::stdin()
        .lines()
        .map(|line| line.unwrap())
        .filter_map(|line| {
            if line.is_empty() {
                None
            } else {
                Some(line.split_once(": ").unwrap().1.to_string())
            }
        })
        .collect::<Vec<String>>();

    Machine {
        reg_a: processed[0].parse().unwrap(),
        reg_b: processed[1].parse().unwrap(),
        reg_c: processed[2].parse().unwrap(),
        ip: 0,
        instructions: processed[3]
            .split(',')
            .into_iter()
            .map(|v| v.parse().unwrap())
            .collect(),
    }
}

fn part_one(mut machine: Machine) -> String {
    let mut output_values = Vec::new();

    while let Some(value) = machine.simulate() {
        output_values.push(value);
    }

    output_values.into_iter().fold(String::new(), |acc, val| {
        if acc.is_empty() {
            acc + &val.to_string()
        } else {
            acc + "," + &val.to_string()
        }
    })
}

fn part_two(machine: Machine) -> i32 {
    let mut working_machine = machine.clone();

    let mut i = 0;
    loop {
        working_machine = machine.clone();
        working_machine.reg_a = i;

        let mut index = 0;
        while let Some(value) = working_machine.simulate() {
            if machine.instructions[index] != value {
                break;
            }
            index += 1;
        }
        if index == machine.instructions.len() {
            break;
        }
        i += 1;
    }

    i
}

fn main() {
    let machine = parse_input();

    println!(
        "The output values produced by the program: {}",
        part_one(machine.clone())
    );
    println!(
        "The lowest positive initial value for register A that causes the program to output a copy of itself: {}",
        part_two(machine)
    );
}
