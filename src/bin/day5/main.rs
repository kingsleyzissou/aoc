use itertools::Itertools;
use std::cmp::Ord;
use std::fs::File;
use std::io::{self, Read};

#[derive(Clone, Copy, Ord, PartialOrd, Eq, PartialEq, Debug)]
struct Range {
    start: u64,
    end: u64,
}

impl Range {
    fn new(start: u64, end: u64) -> Range {
        return Range { start, end };
    }

    fn intersection(&self, other: Range) -> Vec<Range> {
        let start = std::cmp::max(self.start, other.start);
        let end = std::cmp::min(self.end, other.end);
        let mut tmp: Vec<Range> = Vec::new();

        if start < end {
            let intersection = Range::new(start, end);
            tmp.push(intersection);

            // perfect intersection
            if self.start == start && self.end == end {
                return tmp;
            }

            // left overlap
            if self.start < start {
                tmp.push(Range::new(self.start, start - 1));
            }

            // right overlap
            if end < self.end {
                tmp.push(Range::new(end + 1, self.end))
            }

            return tmp;
        }

        // return an empty vector
        return tmp;
    }
}

#[derive(Debug, Clone)]
struct Mapping {
    src: Range,
    dest: Range,
}

impl Mapping {
    fn new(src: Range, dest: Range) -> Mapping {
        return Mapping { src, dest };
    }

    fn transform(&self, range: Range) -> Range {
        // we're dealing with unsigned ints!
        if self.src.start >= self.dest.start {
            let delta = self.src.start - self.dest.start;
            return Range::new(range.start - delta, range.end - delta);
        }
        let delta = self.dest.start - self.src.start;
        return Range::new(range.start + delta, range.end + delta);
    }
}

fn contains_or(maps: &Vec<Mapping>, needle: Range) -> Range {
    for map in maps.clone() {
        if map.src.start < needle.start && needle.start < map.src.end {
            let diff = map.dest.start + needle.start - map.src.start;
            return Range::new(diff, diff);
        }
    }
    return needle;
}

fn intersects_or(mappings: &Vec<Mapping>, needle: Range) -> Vec<Range> {
    let intersections: Vec<Range> = mappings
        .iter()
        .map(|mapping| needle.intersection(mapping.src))
        .filter(|mapping| mapping.len() > 0)
        .map(|ranges| {
            let mut transformations = Vec::new();
            for r in ranges {
                for m in mappings {
                    // off by one error that I can't figure out.
                    // so I'll just include the end \o/
                    if r.start >= m.src.start && r.end <= m.src.end {
                        transformations.push(m.transform(r));
                    }
                }
            }
            return transformations;
        })
        .flat_map(|range| range.into_iter())
        .collect();

    if !intersections.is_empty() {
        return intersections;
    }

    return vec![needle];
}

fn parse_mapping(line: &str) -> Mapping {
    let n = line.split(" ");
    let dest = n.clone().nth(0).unwrap().parse::<u64>().unwrap();
    let src = n.clone().nth(1).unwrap().parse::<u64>().unwrap();
    let size = n.clone().nth(2).unwrap().parse::<u64>().unwrap();
    return Mapping::new(Range::new(src, src + size), Range::new(dest, dest + size));
}

fn parse_mappings(input: Vec<String>) -> Vec<Vec<Mapping>> {
    return input
        .iter()
        .map(|line| {
            return line
                .split("\n")
                .skip(1) // ignore the title
                .map(parse_mapping)
                .sorted_by(|a, b| Ord::cmp(&a.src, &b.src))
                .collect::<Vec<Mapping>>();
        })
        .filter(|m| m.len() > 0)
        .collect();
}

fn parse_seeds(input: &String) -> Vec<u64> {
    return input
        .split(": ")
        .nth(1)
        .expect("string split should have two items")
        .split(" ")
        .map(|s| s.parse::<u64>().expect("string integer"))
        .collect();
}

fn star1(input: Vec<String>) -> u64 {
    let mappings: Vec<Vec<Mapping>> = parse_mappings(input.clone());

    return parse_seeds(input.iter().nth(0).expect("input to have multiple lines"))
        .chunks_exact(1)
        .map(|chunk| Range::new(chunk[0], chunk[0] + 1))
        .fold(u64::MAX, |acc, seed| {
            let mut result = seed;
            for mapping in mappings.iter() {
                result = contains_or(mapping, result);
            }
            return acc.min(result.start);
        });
}

fn star2(input: Vec<String>) -> u64 {
    let seeds: Vec<Range> = parse_seeds(input.iter().nth(0).expect("input to have multiple lines"))
        .chunks_exact(2)
        // this is partially responsible for the off by 1 error
        .map(|chunk| Range::new(chunk[0], chunk[0] + chunk[1] - 1))
        .collect();

    return parse_mappings(input.clone())
        .iter()
        .fold(seeds.clone(), |acc, mapping| {
            let mut intersections = Vec::new();
            for seed in acc.iter() {
                let i = intersects_or(mapping, *seed);
                intersections.extend(i);
            }
            return intersections;
        })
        .into_iter()
        .sorted_by(|a, b| Ord::cmp(&a.start, &b.start))
        .collect::<Vec<Range>>()
        .first()
        .unwrap()
        .start;
}

fn load_input() -> io::Result<Vec<String>> {
    let file = File::open("input/d5.txt")?;
    let mut input = String::new();

    let mut reader = io::BufReader::new(file);
    reader.read_to_string(&mut input)?;

    return Ok(input
        .split("\n\n")
        .map(|s| String::from(s.trim()))
        .collect());
}

fn main() {
    let input = load_input().expect("input should load fine");
    println!("star 1: {}", star1(input.clone()));
    println!("star 2: {}", star2(input.clone()));
}

#[cfg(test)]
mod tests {
    use super::*;

    const TEST_INPUT: &str = "
seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4
";

    fn to_string_vec(input: Vec<&str>) -> Vec<String> {
        return input.iter().map(|s| s.to_string()).collect();
    }

    #[test]
    fn test_star1() {
        let input = to_string_vec(
            TEST_INPUT
                .split("\n\n")
                .map(|s| s.trim())
                .collect::<Vec<&str>>(),
        );
        assert_eq!(star1(input), 35);
    }

    #[test]
    fn test_star2() {
        let input = to_string_vec(
            TEST_INPUT
                .split("\n\n")
                .map(|s| s.trim())
                .collect::<Vec<&str>>(),
        );
        assert_eq!(star2(input), 46);
    }
}
