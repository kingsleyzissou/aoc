package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Score interface {
	score() int
}

type shape interface {
	int() int
}

type opponent struct {
	value string
	shape
}

func (o opponent) int() int {
	if o.value == "A" {
		return 1
	}
	if o.value == "B" {
		return 2
	}
	return 3
}

type yours struct {
	value string
	shape
}

func (y yours) int() int {
	if y.value == "X" {
		return 1
	}
	if y.value == "Y" {
		return 2
	}
	return 3

}

type pair struct {
	opponent opponent
	yours    yours
	Score
}

//lint:ignore deadcode
func calculateFromScore(x, y string) int {
	if x == "A" { // they play rock
		if y == "X" {
			return 3 + 1 // I play rock (draw + 1 for rock)
		}
		if y == "Y" {
			return 6 + 2 // I play paper (win + 2 for paper)
		}
		return 3 // I play scissors (lose + 3 for scissors)
	}

	if x == "B" {
		if y == "X" {
			return 1
		}
		if y == "Y" {
			return 3 + 2
		}
		return 6 + 3
	}

	if y == "X" {
		return 6 + 1
	}

	if y == "Y" {
		return 2
	}

	return 3 + 3
}

func calculateFromResult(x, y string) int {
	if x == "A" { // they play rock
		if y == "X" {
			return 3 // to lose I play scissors
		}
		if y == "Y" {
			return 1 + 3 // to draw I play rock
		}
		return 2 + 6 // to win I play paper
	}

	if x == "B" {
		if y == "X" {
			return 1
		}
		if y == "Y" {
			return 2 + 3
		}
		return 3 + 6
	}

	if y == "X" {
		return 2
	}
	if y == "Y" {
		return 3 + 3
	}

	return 1 + 6
}

func (p pair) score() int {
	x, y := p.opponent.value, p.yours.value
	score := calculateFromResult(x, y)
	return score
}

func getPair(s string) pair {
	var p = pair{}
	r := strings.Split(s, " ")
	p.opponent.value = r[0]
	p.yours.value = r[1]
	return p
}

type summable interface {
	sum() int
}

type Pairs struct {
	pairs []pair
	summable
}

func (p Pairs) sum() int {
	sum := 0
	for _, v := range p.pairs {
		sum = sum + v.score()
	}
	return sum
}

func readFile() []string {
	file, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	return strings.Split(string(file), "\n")
}

func getPairs(s []string) Pairs {
	var p = Pairs{}
	for _, v := range s {
		if v != "" {
			p.pairs = append(p.pairs, getPair(v))
		}
	}
	return p
}

func main() {
	contents := readFile()
	pairs := getPairs(contents)
	fmt.Println("Score: ", pairs.sum())
}
