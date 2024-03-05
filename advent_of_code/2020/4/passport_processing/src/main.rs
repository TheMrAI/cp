use std::io;

#[derive(Debug)]
struct DataSet {
    birth_year: Option<String>,
    issue_year: Option<String>,
    expiration_year: Option<String>,
    height: Option<String>,
    hair_color: Option<String>,
    eye_color: Option<String>,
    passport_id: Option<String>,
    country_id: Option<String>,
}

impl Default for DataSet {
    fn default() -> Self {
        DataSet {
            birth_year: None,
            issue_year: None,
            expiration_year: None,
            height: None,
            hair_color: None,
            eye_color: None,
            passport_id: None,
            country_id: None,
        }
    }
}

fn parse_input() -> Vec<DataSet> {
    let mut data_set = Vec::new();
    let mut data = DataSet::default();
    for line in io::stdin().lines() {
        if line.as_ref().unwrap().len() == 0 {
            data_set.push(data);
            data = DataSet::default();
            continue;
        }
        for segment in line.unwrap().split_ascii_whitespace() {
            let parts = match segment.split_once(':') {
                Some(parts) => parts,
                None => panic!("Malformed input detected: `{}`", segment),
            };

            match parts.0 {
                "byr" => data.birth_year = Some(parts.1.to_string()),
                "iyr" => data.issue_year = Some(parts.1.to_string()),
                "eyr" => data.expiration_year = Some(parts.1.to_string()),
                "hgt" => data.height = Some(parts.1.to_string()),
                "hcl" => data.hair_color = Some(parts.1.to_string()),
                "ecl" => data.eye_color = Some(parts.1.to_string()),
                "pid" => data.passport_id = Some(parts.1.to_string()),
                "cid" => data.country_id = Some(parts.1.to_string()),
                _ => panic!("Unexpected field name recived: `{}`", parts.0),
            };
        }
    }
    data_set.push(data);
    data_set
}

fn part_one(data_set: &Vec<DataSet>) -> usize {
    data_set
        .iter()
        .filter(|data| {
            data.birth_year.is_some()
                && data.issue_year.is_some()
                && data.expiration_year.is_some()
                && data.height.is_some()
                && data.hair_color.is_some()
                && data.eye_color.is_some()
                && data.passport_id.is_some()
        })
        .count()
}

fn part_two(data_set: &Vec<DataSet>) -> usize {
    data_set
        .iter()
        .filter(|data| {
            data.birth_year.is_some()
                && data.issue_year.is_some()
                && data.expiration_year.is_some()
                && data.height.is_some()
                && data.hair_color.is_some()
                && data.eye_color.is_some()
                && data.passport_id.is_some()
        })
        .filter(|data| {
            let birth_year: i32 = match data.birth_year.as_ref().unwrap().parse() {
                Ok(v) => v,
                Err(_) => return false,
            };
            if birth_year < 1920 || birth_year > 2002 {
                return false;
            }

            let issue_year: i32 = match data.issue_year.as_ref().unwrap().parse() {
                Ok(v) => v,
                Err(_) => return false,
            };
            if issue_year < 2010 || issue_year > 2020 {
                return false;
            }

            let expiration_year: i32 = match data.expiration_year.as_ref().unwrap().parse() {
                Ok(v) => v,
                Err(_) => return false,
            };
            if expiration_year < 2020 || expiration_year > 2030 {
                return false;
            }

            let height_text = data.height.as_ref().unwrap();
            if height_text.ends_with("cm") {
                let in_cm: i32 = match &height_text[0..height_text.len() - 2].parse() {
                    Ok(v) => *v,
                    Err(_) => return false,
                };
                if in_cm < 150 || in_cm > 193 {
                    return false;
                }
            } else if height_text.ends_with("in") {
                let in_in: i32 = match &height_text[0..height_text.len() - 2].parse() {
                    Ok(v) => *v,
                    Err(_) => return false,
                };
                if in_in < 59 || in_in > 76 {
                    return false;
                }
            } else {
                return false;
            };

            let hair_color = data.hair_color.as_ref().unwrap();
            if hair_color.len() != 7
                || hair_color.chars().nth(0).unwrap() != '#'
                || !&hair_color[1..]
                    .chars()
                    .all(|c| c.is_ascii_digit() || (c.is_ascii_alphabetic() && c.is_lowercase()))
            {
                return false;
            }

            let eye_color = data.eye_color.as_ref().unwrap();
            if eye_color != "amb"
                && eye_color != "blu"
                && eye_color != "brn"
                && eye_color != "gry"
                && eye_color != "grn"
                && eye_color != "hzl"
                && eye_color != "oth"
            {
                return false;
            }
            let passport_id = data.passport_id.as_ref().unwrap();
            if passport_id.len() != 9 || !passport_id.chars().all(|c| c.is_ascii_digit()) {
                return false;
            }

            true
        })
        .count()
}

fn main() {
    println!("Hello, world!");
    let data_set = parse_input();
    for entry in data_set.iter() {
        println!("{:?}", entry);
    }
    println!("The number of valid passports: {}", part_one(&data_set));
    println!(
        "The number of valid passports, with field checks: {}",
        part_two(&data_set)
    );
}
