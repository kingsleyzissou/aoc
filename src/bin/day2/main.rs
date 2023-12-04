use std::collections::HashMap;
use std::fs::File;
use std::io::{BufRead, BufReader};

#[derive(Debug)]
struct Game {
    id: i32,
    throws: Vec<Throw>,
}

#[derive(Debug)]
struct Throw {
    red: i32,
    green: i32,
    blue: i32,
}

impl Game {
    fn new(id: i32) -> Game {
        return Game {
            id,
            throws: Vec::new(),
        };
    }
}

fn parse_colour(colour: &str) -> (&str, i32) {
    let colour = colour.split(" ").collect::<Vec<&str>>();
    let value = colour[0].parse::<i32>().unwrap();
    return (colour[1], value);
}

fn parse_throw(throw: &str) -> HashMap<&str, i32> {
    return throw
        .split(", ")
        .map(parse_colour)
        .collect::<HashMap<&str, i32>>();
}

fn parser(input: Vec<String>) -> Vec<Game> {
    return input
        .iter()
        .enumerate()
        .map(|(i, line)| {
            let mut game = Game::new((i + 1) as i32);
            let line = line.split(": ").collect::<Vec<&str>>();
            line[1].split("; ").map(parse_throw).for_each(|throw| {
                let t = Throw {
                    red: *throw.get("red").unwrap_or(&0),
                    green: *throw.get("green").unwrap_or(&0),
                    blue: *throw.get("blue").unwrap_or(&0),
                };
                game.throws.push(t);
            });
            return game;
        })
        .collect();
}

fn filter_throws(game: &Game) -> bool {
    for throw in &game.throws {
        if throw.red > 12 || throw.green > 13 || throw.blue > 14 {
            return false;
        }
    }
    return true;
}

fn star1(input: Vec<String>) -> i32 {
    return parser(input)
        .iter()
        .filter(|game| filter_throws(game))
        .map(|game| game.id)
        .sum::<i32>();
}

fn main() {
    let input = File::open("input/d2.txt").unwrap();
    let input = BufReader::new(input)
        .lines()
        .map(|l| l.unwrap())
        .collect::<Vec<String>>();
    println!("Star 1: {}", star1(input.clone()));
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
            "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green",
            "Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue",
            "Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red",
            "Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red",
            "Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green",
        ]);
        assert_eq!(star1(input), 8);
    }

}
