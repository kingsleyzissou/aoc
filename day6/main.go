package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func containsMultiple(needle string, haystack []string) bool {
	count := 0
	for _, line := range haystack {
		if line == needle {
			count += 1
		}
	}
	return count >= 2
}

func checkSequence(sequence []string) bool {
	for _, c := range sequence {
		if containsMultiple(c, sequence) {
			return false
		}
	}
	return true
}

func subroutine(datastream []string, length int) int {
	marker := 0
	for i := length - 1; i < len(datastream); i++ {
		res := checkSequence(datastream[i : i+length])
		if res == true {
			marker = i + length
			break
		}
	}
	return marker
}

func readFile() []string {
	file, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	return strings.Split(string(file), "")
}

func main() {
	datastream := readFile()

	marker := subroutine(datastream, 4)
	fmt.Println("Packet marker: ", marker)
	fmt.Println(datastream[marker-4 : marker])

	marker = subroutine(datastream, 14)
	fmt.Println("Message marker: ", marker)
	fmt.Println(datastream[marker-14 : marker])
}
