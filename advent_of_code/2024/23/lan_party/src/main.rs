use std::collections::{HashMap, HashSet};
use std::io;

fn parse_input() -> Vec<(String, String)> {
    io::stdin()
        .lines()
        .map(|line| {
            let line = line.unwrap();
            let (lhs, rhs) = line.split_once('-').unwrap();
            (lhs.to_string(), rhs.to_string())
        })
        .collect()
}

fn to_adjacency_list(connections: &Vec<(String, String)>) -> HashMap<String, HashSet<String>> {
    let mut adjacency_list = HashMap::<String, HashSet<String>>::new();

    for connection in connections {
        adjacency_list
            .entry(connection.0.clone())
            .and_modify(|neighbors| {
                neighbors.insert(connection.1.clone());
            })
            .or_insert_with(|| {
                let mut set = HashSet::new();
                set.insert(connection.1.clone());
                set
            });
        adjacency_list
            .entry(connection.1.clone())
            .and_modify(|neighbors| {
                neighbors.insert(connection.0.clone());
            })
            .or_insert_with(|| {
                let mut set = HashSet::new();
                set.insert(connection.0.clone());
                set
            });
    }
    adjacency_list
}

fn part_one(adjacency_list: &HashMap<String, HashSet<String>>) -> usize {
    let mut triplets = HashSet::<[String; 3]>::new();

    for (machine1, m1_connections) in adjacency_list {
        for machine2 in m1_connections {
            for machine3 in adjacency_list.get(machine2).unwrap() {
                if !adjacency_list.get(machine3).unwrap().contains(machine1) {
                    continue;
                }
                let mut triplet = [machine1.clone(), machine2.clone(), machine3.clone()];
                triplet.sort();
                triplets.insert(triplet);
            }
        }
    }

    triplets
        .iter()
        .filter(|triplet| triplet.iter().any(|machine| machine.starts_with('t')))
        .count()
}

// A nice little NP hard problem for the end.
// Well, the Bron-Kerbosch is simple to implement and does the job.
fn part_two(adjacency_list: &HashMap<String, HashSet<String>>) -> String {
    let candidate_set = adjacency_list
        .keys()
        .map(|key| key.clone())
        .collect::<HashSet<String>>();

    let mut maximal_cliques = Vec::new();
    bron_kerbosch(
        HashSet::new(),
        candidate_set,
        HashSet::new(),
        &mut maximal_cliques,
        adjacency_list,
    );

    let mut max_i = 0;
    for i in 1..maximal_cliques.len() {
        if maximal_cliques[i].len() > maximal_cliques[max_i].len() {
            max_i = i;
        }
    }

    let mut password_parts = maximal_cliques[max_i]
        .iter()
        .map(|k| k.clone())
        .collect::<Vec<String>>();
    password_parts.sort();
    password_parts
        .iter()
        .fold(String::new(), |mut password, part| {
            if password.is_empty() {
                password += part;
            } else {
                password += &format!(",{part}");
            }
            password
        })
}

fn bron_kerbosch(
    clique: HashSet<String>,
    mut candidates: HashSet<String>,
    mut excluded: HashSet<String>,
    maximal_cliques: &mut Vec<HashSet<String>>,
    adjacency_list: &HashMap<String, HashSet<String>>,
) {
    if candidates.is_empty() && excluded.is_empty() {
        maximal_cliques.push(clique);
        return;
    }
    for v in candidates
        .iter()
        .map(|k| k.clone())
        .collect::<Vec<String>>()
    {
        let mut expanded_clique = clique.clone();
        expanded_clique.insert(v.clone());

        let neighbors_of_v = adjacency_list.get(&v).unwrap();
        let next_candidates = candidates
            .intersection(neighbors_of_v)
            .map(|k| k.clone())
            .collect();
        let next_excluded = excluded
            .intersection(neighbors_of_v)
            .map(|k| k.clone())
            .collect();

        bron_kerbosch(
            expanded_clique,
            next_candidates,
            next_excluded,
            maximal_cliques,
            adjacency_list,
        );
        candidates.remove(&v);
        excluded.insert(v);
    }
}

fn main() {
    let connections = parse_input();

    let adjacency_list = to_adjacency_list(&connections);
    println!("The number of connected computer triplets that contain at least one with a name that starts with t: {}", part_one(&adjacency_list));

    println!(
        "The password to the LAN party is: {}",
        part_two(&adjacency_list)
    );
}
