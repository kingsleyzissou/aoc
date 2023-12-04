use pcre2::bytes::{Captures, Regex};
use std::fs::File;
use std::io::{BufRead, BufReader};

fn load(input: String) -> Vec<String> {
    let file = File::open(input).unwrap();
    return BufReader::new(file).lines().map(|l| l.unwrap()).collect();
}

fn main() {
    let input = load("input/d1.txt".to_string());

    let star1 = star1(input.clone());
    println!("Star 1: {}", star1);

    let star2 = star2(input.clone());
    println!("Star 2: {}", star2);
}

fn star1(input: Vec<String>) -> u32 {
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

fn star2(input: Vec<String>) -> u32 {
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

    fn to_string_vec(input: Vec<&str>) -> Vec<String> {
        return input.iter().map(|s| s.to_string()).collect();
    }

    #[test]
    fn test_star1() {
        let input = to_string_vec(vec!["1abc2", "pqr3stu8vwx", "a1b2c3d4e5f", "treb7uchet"]);
        assert_eq!(star1(input), 142);
    }

    #[test]
    fn test_star2() {
        let input = to_string_vec(vec![
            "two1nine",
            "eightwothree",
            "abcone2threexyz",
            "xtwone3four",
            "4nineeightseven2",
            "zoneight234",
            "7pqrstsixteen",
        ]);
        assert_eq!(star2(input), 281);
    }
}
