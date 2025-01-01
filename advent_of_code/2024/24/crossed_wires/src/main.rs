use std::collections::{HashMap, HashSet};
use std::io;

use regex::Regex;

#[derive(Debug, Copy, Clone, Eq, PartialEq, Hash)]
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

// In hindsight I see that this wasn't entirely necessary.
// We could have done it with a simple DFS starting from the individual z bits.
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
    for gate in &gates {
        println!("{:?}", gate);
    }

    println!(
        "The decimal number output on the wires is: {}",
        part_one(&init_values, &gates)
    );

    let mut counter = HashMap::new();
    for gate in &gates {
        counter
            .entry(gate.ty)
            .and_modify(|count| *count += 1)
            .or_insert(1);
    }
    for (k, v) in counter {
        println!("{:?}: {}", k, v);
    }

    let mut tracer = HashMap::new();
    for gate in &gates {
        tracer.insert(
            gate.result.clone(),
            (gate.ty, [gate.lhs.clone(), gate.rhs.clone()]),
        );
    }
    for i in (0..gates.len()).rev() {
        if !gates[i].result.starts_with('z') {
            break;
        }
        println!();
        dfs(&tracer, &gates[i].result, 0);
    }
    // println!("Similarity score: {}", part_two(&lhs, &rhs));
}

fn dfs(tracer: &HashMap<String, (GateTy, [String; 2])>, gate: &String, depth: usize) {
    if depth > 3 {
        return;
    }
    let mut tab = String::new();
    for _ in 0..depth {
        tab += " ";
    }
    if let Some(d) = tracer.get(gate) {
        println!("{}{} {:?}", tab, gate, d.0);
        dfs(tracer, &d.1[0], depth + 1);
        dfs(tracer, &d.1[1], depth + 1);
    } else {
        println!("{}{}", tab, gate);
    }
}

// A general one bit adder uses gates:
// 2 XOR
// 2 AND
// 1 OR
// for every single bit.
// We have 45 bits 45*5 = 225
// but because the first bit doesn't need the incoming carry bit logic
// it can do be done with only 2 gates (XOR for the input and OR for the carry)
// 225 - 3 = 222
// And we have 222 gates in the input.
// This means that we can be pretty sure that the input describes just an
// optimal n bit adder (as we started with the minimal possible gate count necessary, there is no room
// for shenanigans as the exercise says this adder should work on all inputs).
// This means that we can pretty quickly figure out which gates are definitely incorrect
// by going back the chain and looking for the expected pattern.

// 0-th bit needs:
// XOR
//   x0
//   y0

// n-th bit needs:
// XOR
//  XOR
//    x[n]
//    y[n]
//  OR
//    AND
//      (the same two legs as the top XOR from previous step, the carry in and xor of the inputs)
//    AND
//      x[n-1]
//      y[n-1]

// last bit needs (the carry out complex of the "previous" bit):
//  OR
//    AND
//      (the same two legs as the top XOR from previous step)
//    AND
//      x[n-1]
//      y[n-1]

// AND: 89
// XOR: 89
// OR: 44

// Gate { lhs: "y30", ty: AND, rhs: "x30", result: "rtq" }
// Gate { lhs: "rmb", ty: AND, rhs: "rpt", result: "kcn" }
// z06 - hwk
// hpc - z31
// cgr - z37
// tnt - qmd

// cgr,hpc,hwk,qmd,tnt,z06,z31,z37