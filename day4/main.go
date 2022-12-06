package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type contains interface {
	Contains(o section) bool
}

type overlaps interface {
	Overlaps(o section) bool
}

type section struct {
	upper int
	lower int
	contains
	overlaps
}

func (s section) Contains(o section) bool {
	return o.upper <= s.upper && o.lower >= s.lower
}

func (s section) Overlaps(o section) bool {
	if o.lower <= s.upper && o.lower >= s.lower {
		return true
	}
	if o.upper >= s.lower && o.upper <= s.upper {
		return true
	}
	return false
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func toSection(s string) section {
	split := strings.Split(s, "-")
	lower := parseInt(split[0])
	upper := parseInt(split[1])
	return section{upper: upper, lower: lower}
}

type pair struct {
	a section
	b section
}

func toPair(s string) pair {
	split := strings.Split(s, ",")
	a := toSection(split[0])
	b := toSection(split[1])
	return pair{a: a, b: b}
}

func toPairs(s []string) []pair {
	var pairs []pair
	for _, p := range s {
		if p == "" {
			break
		}
		pairs = append(pairs, toPair(p))
	}
	return pairs
}

func readFile() []string {
	file, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	return strings.Split(string(file), "\n")
}

func main() {
	contents := readFile()
	pairs := toPairs(contents)
	sum := 0
	for _, p := range pairs {
		if p.a.Overlaps(p.b) || p.b.Overlaps(p.a) {
			sum += 1
		}
	}
	fmt.Println("Sum: ", sum, " Length: ", len(pairs))
}
