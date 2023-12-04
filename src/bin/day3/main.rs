use regex::Regex;
use std::collections::HashMap;
use std::fs::File;
use std::io::{BufRead, BufReader};

const DIRECTIONS: [Point; 9] = [
    Point { x: 0, y: 0 },   // current
    Point { x: 0, y: -1 },  // up
    Point { x: 1, y: 0 },   // right
    Point { x: 0, y: 1 },   // down
    Point { x: -1, y: 0 },  // left
    Point { x: -1, y: -1 }, // up left
    Point { x: 1, y: -1 },  // up right
    Point { x: -1, y: 1 },  // down left
    Point { x: 1, y: 1 },   // down right
];

#[derive(Eq, Debug, PartialEq, Hash, Clone, Copy)]
struct Point {
    x: i32,
    y: i32,
}

fn insert_point(points: &mut HashMap<Point, char>, point: Point, value: char) {
    if value != '.' && !value.is_digit(10) {
        points.insert(point, value);
    }
}

fn get_points(input: Vec<String>) -> HashMap<Point, char> {
    let mut points = HashMap::new();
    for (y, line) in input.iter().enumerate() {
        for (x, c) in line.chars().enumerate() {
            let point = Point {
                x: x as i32,
                y: y as i32,
            };
            insert_point(&mut points, point, c);
        }
    }
    return points;
}

fn insert_component(
    components: &mut HashMap<Point, Vec<i32>>,
    points: &HashMap<Point, char>,
    start: Point,
    end: Point,
    value: i32,
) {
    for x in start.x..end.x {
        for direction in DIRECTIONS.iter() {
            let point = Point {
                x: x as i32 + direction.x,
                y: start.y as i32 + direction.y,
            };
            let entry = components.entry(point).or_insert(Vec::new());
            // check for duplicate values
            if !entry.contains(&value) && points.contains_key(&point) {
                entry.push(value);
            }
        }
    }
}

fn get_components(input: Vec<String>, points: HashMap<Point, char>) -> HashMap<Point, Vec<i32>> {
    let pattern = r"(\d+)";
    let regex = Regex::new(pattern).unwrap();
    let mut components = HashMap::new();
    for (y, line) in input.iter().enumerate() {
        regex.find_iter(line).for_each(|m| {
            let value = m.as_str().parse::<i32>().unwrap();
            let start = Point {
                x: m.start() as i32,
                y: y as i32,
            };
            let end = Point {
                x: m.end() as i32,
                y: y as i32,
            };
            insert_component(&mut components, &points, start, end, value);
        });
    }
    return components;
}

fn star1(input: Vec<String>) -> i32 {
    let points = get_points(input.clone());
    return get_components(input.clone(), points.clone())
        .iter()
        .map(|(_, v)| v.iter().sum::<i32>())
        .sum::<i32>();
}

fn main() {
    let input = File::open("input/d3.txt").unwrap();
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
            "467..114..",
            "...*......",
            "..35..633.",
            "......#...",
            "617*......",
            ".....+.58.",
            "..592.....",
            "......755.",
            "...$.*....",
            ".664.598..",
        ]);
        assert_eq!(star1(input), 4361);
    }
}
