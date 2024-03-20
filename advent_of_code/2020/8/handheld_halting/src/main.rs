use regex::Regex;
use std::collections::HashSet;

#[derive(Debug, Clone, Copy, Eq, PartialEq)]
enum CommandType {
    Nop,
    Acc,
    Jmp,
}

#[derive(Debug, Clone, Copy)]
struct Instruction {
    command: CommandType,
    value: i64,
}

fn parse_input(source: std::io::Stdin) -> Vec<Instruction> {
    let mut instructions = Vec::new();

    let number_regexp = Regex::new(r"([+-]\d+)").unwrap();
    for line in source.lines() {
        let line = line.unwrap();
        let parts: Vec<&str> = line.split(" ").collect();
        if parts.len() != 2 {
            panic!("Unexpected input received. Got: '{}'", line);
        }
        let Some((_, [value_text])) = number_regexp.captures(parts[1]).map(|caps| caps.extract())
        else {
            panic!("Failed to match value: '{}'", parts[1]);
        };

        let value: i64 = match value_text.parse() {
            Ok(v) => v,
            Err(_) => panic!("Failed to parse integer from input: '{}'", value_text),
        };

        let command = match parts[0] {
            "nop" => CommandType::Nop,
            "acc" => CommandType::Acc,
            "jmp" => CommandType::Jmp,
            _ => panic!("Unexpected command received: '{}'", parts[0]),
        };

        instructions.push(Instruction { command, value });
    }

    instructions
}

fn simulate_program(instructions: &[Instruction]) -> (bool, i64) {
    let mut acc = 0;

    let mut seen = HashSet::new();
    let mut index = 0;
    while index < instructions.len() {
        if seen.contains(&index) {
            return (false, acc);
        }
        seen.insert(index);
        let instruction = &instructions[index];
        match instruction.command {
            CommandType::Nop => index += 1,
            CommandType::Acc => {
                acc += instruction.value;
                index += 1
            }
            CommandType::Jmp => index = (index as i64 + instruction.value) as usize,
        }
    }

    (true, acc)
}

fn part_one(instructions: &[Instruction]) -> i64 {
    simulate_program(instructions).1
}

// Simple stupid brute force. No need to complicate.
fn part_two(instructions: &[Instruction]) -> i64 {
    let mut working_instructions = Vec::from(instructions);

    for index in 0..working_instructions.len() {
        if working_instructions[index].command == CommandType::Acc {
            continue;
        }
        if working_instructions[index].command == CommandType::Nop {
            working_instructions[index].command = CommandType::Jmp;
            let (ok, value) = simulate_program(&working_instructions);
            if ok {
                return value;
            }
            working_instructions[index].command = CommandType::Nop;
        } else if working_instructions[index].command == CommandType::Jmp {
            working_instructions[index].command = CommandType::Nop;
            let (ok, value) = simulate_program(&working_instructions);
            if ok {
                return value;
            }
            working_instructions[index].command = CommandType::Jmp;
        }
    }
    -1
}

fn main() {
    let instructions = parse_input(std::io::stdin());
    println!(
        "The value of the accumulator before the start of the infinite loop: {}",
        part_one(&instructions)
    );
    println!(
        "The value of the accumulator for a terminating program: {}",
        part_two(&instructions)
    );
}
