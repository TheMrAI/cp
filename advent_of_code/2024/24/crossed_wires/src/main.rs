use std::collections::{HashMap, HashSet};
use std::io;

use regex::Regex;

#[derive(Debug, Copy, Clone)]
enum GateTy {
    AND,
    OR,
    XOR,
}

#[derive(Clone, Debug)]
struct Gate {
    pub lhs: String,
    pub ty: GateTy,
    pub rhs: String,
    pub result: String,
}

fn parse_input() -> (Vec<(String, bool)>, Vec<Gate>) {
    let mut input = io::stdin().lines();

    let init_re = Regex::new(r"(\w\d{2}): (\d)").unwrap();
    let init_values = input
        .by_ref()
        .take_while(|line| !line.as_ref().unwrap().is_empty())
        .map(|line| {
            let line = line.unwrap();
            let captures = init_re.captures(&line).unwrap();
            let set = if captures.get(2).unwrap().as_str() == "0" {
                false
            } else {
                true
            };
            (captures.get(1).unwrap().as_str().to_owned(), set)
        })
        .collect();

    let gates_re = Regex::new(r"(.+) (AND|XOR|OR) (.+) -> (.+)").unwrap();
    let gates = input
        .map(|line| {
            let line = line.unwrap();
            let captures = gates_re.captures(&line).unwrap();
            let gate_str = captures.get(2).unwrap().as_str();
            let gate = match gate_str {
                "AND" => GateTy::AND,
                "OR" => GateTy::OR,
                "XOR" => GateTy::XOR,
                _ => panic!("Unknown gate: {}", gate_str),
            };

            Gate {
                lhs: captures.get(1).unwrap().as_str().to_owned(),
                ty: gate,
                rhs: captures.get(3).unwrap().as_str().to_owned(),
                result: captures.get(4).unwrap().as_str().to_owned(),
            }
        })
        .collect();

    (init_values, gates)
}

fn part_one(init_values: &[(String, bool)], ordered_gates: &[Gate]) -> usize {
    let mut value_of = HashMap::new();
    for init_value in init_values {
        value_of.insert(init_value.0.clone(), init_value.1);
    }

    for gate in ordered_gates {
        let lhs = value_of.get(&gate.lhs).unwrap();
        let rhs = value_of.get(&gate.rhs).unwrap();
        let result = match gate.ty {
            GateTy::AND => lhs & rhs,
            GateTy::OR => lhs | rhs,
            GateTy::XOR => lhs ^ rhs,
        };
        value_of.insert(gate.result.clone(), result);
    }

    let mut decimal_number = 0;
    for i in (0..ordered_gates.len()).rev() {
        if !ordered_gates[i].result.starts_with('z') {
            break;
        }
        decimal_number <<= 1;
        decimal_number |= if *value_of.get(&ordered_gates[i].result).unwrap() {
            1
        } else {
            0
        };
    }
    decimal_number
}

fn reorder_gates(init_values: &Vec<(String, bool)>, gates: &Vec<Gate>) -> Vec<Gate> {
    let mut resolved = init_values
        .iter()
        .map(|init| init.0.clone())
        .collect::<HashSet<String>>();

    let mut outputs = Vec::new();
    let mut calculate_gates = Vec::new();
    for gate in gates {
        if gate.result.starts_with('z') {
            outputs.push(gate.clone());
        } else {
            calculate_gates.push(gate.clone());
        }
    }
    // simply sort them so we don't have to worry about bit order no more
    outputs.sort_by(|lhs, rhs| lhs.result.cmp(&rhs.result));

    let mut reordered_gates = Vec::new();
    while !calculate_gates.is_empty() {
        let mut next_calculates_gates = Vec::new();
        for gate in &calculate_gates {
            if resolved.contains(&gate.lhs) && resolved.contains(&gate.rhs) {
                resolved.insert(gate.result.clone());
                reordered_gates.push(gate.clone());
                continue;
            }
            next_calculates_gates.push(gate.clone());
        }

        std::mem::swap(&mut calculate_gates, &mut next_calculates_gates);
    }

    for output in outputs {
        reordered_gates.push(output);
    }

    reordered_gates
}

fn main() {
    let (init_values, gates) = parse_input();

    let gates = reorder_gates(&init_values, &gates);

    println!(
        "The decimal number output on the wires is: {}",
        part_one(&init_values, &gates)
    );
    // println!("Similarity score: {}", part_two(&lhs, &rhs));
}
