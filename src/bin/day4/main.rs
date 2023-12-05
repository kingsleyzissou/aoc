use std::collections::HashMap;
use std::fs::File;
use std::io::{BufRead, BufReader};

#[derive(Debug, Clone)]
struct Card {
    index: usize,
    winners: i32,
}

fn parse_numbers(line: &str) -> Vec<i32> {
    let regex = regex::Regex::new(r"(\d+)").unwrap();
    return regex
        .captures_iter(line)
        .map(|c| c.get(1).unwrap().as_str().parse::<i32>().unwrap())
        .collect();
}

fn parse_line(index: usize, line: &String) -> Card {
    let sides = line.split(" | ");

    // winning numbers
    let winning = parse_numbers(sides.clone().nth(0).unwrap());

    // numbers got
    let got = parse_numbers(sides.clone().nth(1).unwrap())
        .iter()
        .filter(|n| winning.contains(n))
        .map(|n| *n)
        .collect::<Vec<i32>>();

    return Card {
        index,
        winners: got.len() as i32,
    };
}
fn star1(input: Vec<String>) -> i32 {
    return input
        .iter()
        .map(|s| s.split(": ").nth(1).unwrap())
        .enumerate()
        .map(|(i, s)| parse_line(i, &String::from(s)))
        .filter(|c| c.winners > 0)
        .map(|c| {
            let base: i32 = 2;
            return base.pow((c.winners - 1) as u32) as i32;
        })
        .sum();
}

fn star2(input: Vec<String>) -> i32 {
    let mut cards = input
        .iter()
        .map(|s| s.split(": ").nth(1).unwrap())
        .enumerate()
        .map(|(i, s)| parse_line(i, &String::from(s)))
        .rev()
        .collect::<Vec<Card>>();

    // this makes things slightly more efficient
    let mut lookup = HashMap::new();
    for card in cards.clone() {
        lookup.insert(card.index, card);
    }

    let mut copies = cards.len() as i32;

    // still not super efficient :[
    while cards.len() > 0 {
        let card = cards.pop().unwrap();
        copies += card.winners;
        for i in card.index + 1..card.index + card.winners as usize + 1 {
            let c = lookup.get(&i).unwrap();
            cards.push(c.clone());
        }
    }

    return copies;
}

fn main() {
    let input = File::open("input/d4.txt").unwrap();
    let input = BufReader::new(input)
        .lines()
        .map(|l| l.unwrap())
        .collect::<Vec<String>>();
    println!("Star 1: {}", star1(input.clone()));
    println!("Star 2: {}", star2(input.clone()));
}

#[cfg(test)]
mod tests {
    use super::*;

    fn to_string_vec(input: Vec<&str>) -> Vec<String> {
        return input.iter().map(|s| s.to_string()).collect();
    }

    #[test]
    fn test_star1() {
        let input = to_string_vec(vec![
            "Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53",
            "Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19",
            "Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1",
            "Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83",
            "Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36",
            "Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11",
        ]);
        assert_eq!(star1(input), 13);
    }

    #[test]
    fn test_star2() {
        let input = to_string_vec(vec![
            "Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53",
            "Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19",
            "Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1",
            "Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83",
            "Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36",
            "Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11",
        ]);
        assert_eq!(star2(input), 30);
    }
}
