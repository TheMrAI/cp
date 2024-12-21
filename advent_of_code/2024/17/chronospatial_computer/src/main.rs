#[derive(Debug, Clone)]
struct Machine {
    pub reg_a: u64,
    pub reg_b: u64,
    pub reg_c: u64,
    pub ip: usize,
    pub instructions: Vec<u64>,
}

impl Machine {
    pub fn combo(&self, operand: u64) -> u64 {
        match operand {
            0 | 1 | 2 | 3 => operand,
            4 => self.reg_a,
            5 => self.reg_b,
            6 => self.reg_c,
            _ => panic!("Unexpected combo operand: {}", operand),
        }
    }

    pub fn simulate(&mut self) -> Option<u64> {
        while self.ip < self.instructions.len() {
            let instruction = self.instructions[self.ip];
            let operand = self.instructions[self.ip + 1];
            self.ip += 2;

            match instruction {
                0 => self.reg_a = self.reg_a / 2u64.pow(self.combo(operand) as u32),
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
                    return Some(self.combo(operand) as u64 % 8);
                }
                6 => self.reg_b = self.reg_a / 2u64.pow(self.combo(operand) as u32),
                7 => self.reg_c = self.reg_a / 2u64.pow(self.combo(operand) as u32),
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

fn part_two(machine: Machine) -> u64 {
    // we immediately jump to the first number that will have the required digits
    let mut candidate_a = 8u64.pow(machine.instructions.len() as u32 - 1);

    let mut cracking_digit = machine.instructions.len() as u64 - 1;
    loop {
        loop {
            let mut checking_machine = machine.clone();
            checking_machine.reg_a = candidate_a;

            // simulate
            let mut output_values = Vec::new();

            while let Some(value) = checking_machine.simulate() {
                output_values.push(value);
            }

            if output_values[cracking_digit as usize]
                == machine.instructions[cracking_digit as usize]
            {
                break;
            } else {
                candidate_a += 8u64.pow(cracking_digit as u32);
            }
        }

        if let Some(next_digit) = cracking_digit.checked_sub(1) {
            cracking_digit = next_digit;
        } else {
            break;
        }
    }

    // we are close, increase by one until we hit the target
    loop {
        let mut checking_machine = machine.clone();
        checking_machine.reg_a = candidate_a;

        // simulate
        let mut output_values = Vec::new();

        while let Some(value) = checking_machine.simulate() {
            output_values.push(value);
        }

        if output_values == machine.instructions {
            break;
        } else {
            candidate_a += 1;
        }
    }

    candidate_a
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
