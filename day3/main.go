package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

var lookup = map[string]int{
	"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6, "g": 7, "h": 8, "i": 9, "j": 10, "k": 11, "l": 12, "m": 13, "n": 14, "o": 15, "p": 16,
	"q": 17, "r": 18, "s": 19, "t": 20, "u": 21, "v": 22, "w": 23, "x": 24, "y": 25, "z": 26,
	"A": 27, "B": 28, "C": 29, "D": 30, "E": 31, "F": 32, "G": 33, "H": 34, "I": 35, "J": 36, "K": 37, "L": 38, "M": 39, "N": 40, "O": 41, "P": 42,
	"Q": 43, "R": 44, "S": 45, "T": 46, "U": 47, "V": 48, "W": 49, "X": 50, "Y": 51, "Z": 52,
}

type hasValue interface {
	value() int
}

type alphabetic struct {
	key string
	hasValue
}

func (a alphabetic) value() int {
	v, ok := lookup[a.key]
	if !ok {
		panic("No idea what happened")
	}
	return v
}

type compareable interface {
	compare() alphabetic
}

type compartment struct {
	compartment []string
}

func (c compartment) String() string {
	return strings.Join(c.compartment, "")
}

func (c compartment) Len() int {
	return len(c.compartment)
}

type backpack struct {
	a compartment
	b compartment
	compareable
}

func (b backpack) allContents() compartment {
	c := append(b.a.compartment, b.b.compartment...)
	sort.Strings(c)
	return compartment{compartment: c}
}

func (b backpack) compare() alphabetic {
	for _, a := range b.a.compartment {
		for _, b := range b.b.compartment {
			if a == b {
				return alphabetic{key: a}
			}
		}
	}
	return alphabetic{}
}

type group struct {
	one   backpack
	two   backpack
	three backpack
	compareable
}

// this is super ugly, but hey
func (g group) compare() alphabetic {
	for _, a := range g.one.allContents().compartment {
		for _, b := range g.two.allContents().compartment {
			for _, c := range g.three.allContents().compartment {
				if a == b && b == c {
					return alphabetic{key: a}
				}
			}
		}
	}
	return alphabetic{}
}

func readFile() []string {
	file, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	return strings.Split(string(file), "\n")
}

func getCompartment(s string) compartment {
	a := strings.Split(s, "")
	sort.Strings(a)
	return compartment{compartment: a}
}

func getCompartments(s string) backpack {
	var backpack backpack
	half := (len(s) / 2)
	backpack.a = getCompartment(s[:half])
	backpack.b = getCompartment(s[half:])
	return backpack
}

func getBackpacks(s []string) []backpack {
	var backpacks []backpack
	for _, v := range s {
		if v != "" {
			backpacks = append(backpacks, getCompartments(v))
		}
	}
	return backpacks
}

func getGroups(b []backpack) []group {
	var groups []group
	for i := 0; i < len(b); i = i + 3 {
		groups = append(groups, group{
			one:   b[i],
			two:   b[i+1],
			three: b[i+2],
		})
	}
	return groups
}

func main() {
	contents := readFile()
	backpacks := getBackpacks(contents)
	groups := getGroups(backpacks)
	sum := 0

	for _, g := range groups {
		sum = sum + g.compare().value()
	}

	fmt.Println("Sum: ", sum)
}
