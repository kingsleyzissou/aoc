use pcre2::bytes::{Captures, Regex};
use std::fs::File;
use std::io::{BufRead, BufReader};

pub fn load(input: String) -> Vec<String> {
    let file = File::open(input).unwrap();
    return BufReader::new(file).lines().map(|l| l.unwrap()).collect();
}

pub fn star1(input: Vec<String>) -> u32 {
    return input
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

fn parse_numbers(capture: Captures<'_>) -> u32 {
    let s = std::str::from_utf8(capture.get(1).unwrap().as_bytes()).unwrap();
    return match s {
        "one" => 1,
        "two" => 2,
        "three" => 3,
        "four" => 4,
        "five" => 5,
        "six" => 6,
        "seven" => 7,
        "eight" => 8,
        "nine" => 9,
        _ => s.parse::<u32>().unwrap(),
    };
}

pub fn star2(input: Vec<String>) -> u32 {
    let pattern = r"(?=(\d|one|two|three|four|five|six|seven|eight|nine))";
    let regex = Regex::new(pattern).unwrap();

    let sum = input
        .iter()
        .map(|line| {
            let digits = regex
                .captures_iter(&line.as_bytes())
                .map(|capture| return parse_numbers(capture.unwrap()))
                .collect::<Vec<u32>>();
            return digits.first().unwrap() * 10 + digits.last().unwrap();
        })
        .sum::<u32>();
    return sum;
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_star1() {
        let input = load("input/test_d1s1.txt".to_string());
        assert_eq!(star1(input), 142);
    }

    #[test]
    fn test_star2() {
        let input = load("input/test_d1s2.txt".to_string());
        assert_eq!(star2(input), 281);
    }
}
