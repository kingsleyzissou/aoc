package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type hasLength interface {
	len() int
}

type calories struct {
	calories []int
	hasLength
}

func (c calories) len() int {
	length := 0
	for _, v := range c.calories {
		length = length + v
	}
	return length
}

func stringsToCalories(s []string) calories {
	var c calories
	for _, v := range s {
		i, _ := strconv.Atoi(v)
		c.calories = append(c.calories, i)
	}
	return c
}

func getCalories(original []string) []calories {
	prevIndex := 0
	var result []calories
	for i, v := range original {
		if v == "" {
			ints := stringsToCalories(original[prevIndex:i])
			result = append(result, ints)
			prevIndex = i + 1
		}
	}
	return result
}

func readFile() []string {
	file, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(fmt.Sprintf("Error reading file: %v", err))
	}
	return strings.Split(string(file), "\n")
}

func main() {
	contents := readFile()

	subSlice := getCalories(contents)
	sort.Slice(subSlice, func(i, j int) bool {
		return subSlice[i].len() > subSlice[j].len()
	})

	first := subSlice[0].len()
	second := subSlice[1].len()
	third := subSlice[2].len()
	fmt.Println("First: ", first)
	fmt.Println("Second: ", second)
	fmt.Println("Third: ", third)
	fmt.Println("Top 3: ", first+second+third)
}
