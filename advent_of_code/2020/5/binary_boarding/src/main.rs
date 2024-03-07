use std::io;

#[derive(Debug)]
pub struct Seat {
    pub row: u8,
    pub column: u8,
}

impl Seat {
    pub fn id(&self) -> u16 {
        8 * self.row as u16 + self.column as u16
    }
}

impl TryFrom<&str> for Seat {
    type Error = &'static str;

    fn try_from(seat_code: &str) -> Result<Self, Self::Error> {
        if seat_code.len() != 10 {
            return Err("Get fucked");
        }
        let mut row_low = 0;
        let mut row_high = 128;
        let mut column_low = 0;
        let mut column_high = 8;
        for (i, c) in seat_code.chars().enumerate() {
            if i < 7 {
                let mid = row_low + ((row_high - row_low) / 2);
                match c {
                    'F' => row_high = mid,
                    'B' => row_low = mid,
                    _ => panic!("Unexpected character received in seat code: '{}'", c),
                }
                continue;
            }
            let mid = column_low + ((column_high - column_low) / 2);
            match c {
                'L' => column_high = mid,
                'R' => column_low = mid,
                _ => panic!("Unexpected character received in seat code: '{}'", c),
            }
        }

        return Ok(Seat {
            row: row_low,
            column: column_low,
        });
    }
}

fn parse_input() -> Vec<Seat> {
    let mut seats = Vec::new();

    for line in io::stdin().lines() {
        let seat = match Seat::try_from(&line.unwrap()[..]) {
            Ok(seat) => seat,
            Err(e) => panic!("{}", e),
        };
        seats.push(seat);
    }

    seats
}

fn part_one(seats: &Vec<Seat>) -> u16 {
    seats
        .iter()
        .map(|seat| seat.id())
        .max()
        .expect("No seats where given.")
}

fn part_two(seats: &Vec<Seat>) -> Option<u16> {
    let mut seat_ids: Vec<u16> = seats.iter().map(|seat| seat.id()).collect();
    seat_ids.sort_unstable();
    let mut prev = 0;
    for (i, &v) in seat_ids.iter().enumerate() {
        if i == 0 {
            prev = v;
            continue;
        }
        if v - 1 != prev {
            return Some(v - 1);
        }
        prev = v;
    }
    return None;
}

fn main() {
    let seats = parse_input();
    println!(
        "The highest seat ID on a boarding pass: {}",
        part_one(&seats)
    );
    println!("Your seat ID is: {}", part_two(&seats).unwrap());
}
