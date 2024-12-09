use std::collections::HashMap;
use std::io;

fn parse_input() -> Vec<u8> {
    let mut line = String::new();
    io::stdin().read_line(&mut line).unwrap();
    line.bytes().map(|b| b - '0' as u8).collect()
}

fn part_one(mut disk_map: Vec<u8>) -> usize {
    let mut front = 0;
    let mut front_id = 0;
    let (mut back, mut back_id) = if disk_map.len() % 2 != 0 {
        (disk_map.len() - 1, disk_map.len() / 2)
    } else {
        (disk_map.len() - 2, disk_map.len() - 1 / 2)
    };

    let mut compressed = vec![];
    while front <= back {
        if front % 2 == 0 {
            while disk_map[front] > 0 {
                compressed.push(front_id);
                disk_map[front] -= 1;
            }
            front += 1;
            // front_id += 1;
        } else {
            if front == back {
                // final case where front has reached the back and
                // we must carry the remaining elements
                while disk_map[back] > 0 {
                    compressed.push(back_id);
                    disk_map[back] -= 1;
                }
                break;
            }
            while disk_map[front] > 0 && disk_map[back] > 0 {
                compressed.push(back_id);
                disk_map[front] -= 1;
                disk_map[back] -= 1;
            }
            if disk_map[front] == 0 {
                front += 1;
                front_id += 1;
            }
            if disk_map[back] == 0 {
                back -= 2;
                back_id -= 1;
            }
        }
    }
    compressed.iter().enumerate().map(|(i, v)| i * v).sum()
}

fn part_two(disk_map: Vec<u8>) -> usize {
    let mut block_compressed_disk_map = disk_map
        .iter()
        .enumerate()
        .map(|(i, v)| {
            if i % 2 == 0 {
                ((i + 1) / 2, *v, true)
            } else {
                (0, *v, false)
            }
        })
        .collect::<Vec<(usize, u8, bool)>>();

    let mut block_position = block_compressed_disk_map
        .iter()
        .rposition(|elem| elem.2)
        .unwrap();
    let mut next_block_id = block_compressed_disk_map[block_position].0;
    // zeroth block never has to move
    while next_block_id > 0 {
        if block_compressed_disk_map[block_position].0 == next_block_id {
            let candidate_space = block_compressed_disk_map.iter().enumerate().find_map(
                |(i, (_, length, is_block))| {
                    if !is_block && i < block_position && *length >= block_compressed_disk_map[block_position].1 {
                        Some(i)
                    } else {
                        None
                    }
                },
            );
            if let Some(i) = candidate_space {
                block_compressed_disk_map[i].1 -= block_compressed_disk_map[block_position].1;
                block_compressed_disk_map.insert(i, block_compressed_disk_map[block_position]);
                // then turn the moved block into zero size
                // note the +1 offsetting the index
                // block_compressed_disk_map[block_position + 1].1 = 0;
                block_compressed_disk_map[block_position + 1].2 = false;
            } else {
                block_position -= 1;
            }
            next_block_id -= 1;
        } else {
            // remember all spaces have id of zero, but we also never disturb the zeroth
            // block so it is enough to check the block id here and we can be assured the next
            // potentially movable space is found
            while block_compressed_disk_map[block_position].0 != next_block_id {
                block_position -= 1;
            }
        }
    }

    let mut counter = 0;
    block_compressed_disk_map
        .iter()
        .map(|(id, count, is_block)| {
            let mut total = 0;
            for _ in 0..*count {
                if *is_block {
                    total += id * counter;
                }
                counter += 1;
            }
            total
        })
        .sum()
}

fn main() {
    let disk_map = parse_input();

    println!(
        "Checksum of compressed filesystem: {}",
        part_one(disk_map.clone())
    );
    println!(
        "Checksum of compressed filesystem using whole blocks: {}",
        part_two(disk_map)
    );
}
