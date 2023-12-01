use std::fs::File;
use std::io::{BufRead, BufReader};

fn load(input: String) -> Vec<String> {
    let file = File::open(input).unwrap();
    return BufReader::new(file).lines().map(|l| l.unwrap()).collect();
}

pub fn star1() -> u32 {
    let lines = load("input/d1.txt".to_string());
    return lines
        .iter()
        .map(|line| {
            let digits = line
                .chars()
                .filter(|c| c.is_digit(10))
                .map(|c| c.to_digit(10).unwrap())
                .collect::<Vec<u32>>();
            return digits.first().unwrap() * 10 + digits.last().unwrap();
        })
        .sum::<u32>();
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_star1() {
        assert_eq!(star1(), 142);
    }
}
