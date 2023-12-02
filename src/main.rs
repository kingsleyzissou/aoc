mod day1;

fn main() {
    let input = day1::load("input/d1.txt".to_string());

    let star1 = day1::star1(input.clone());
    println!("Star 1: {}", star1);

    let star2 = day1::star2(input.clone());
    println!("Star 2: {}", star2);
}
