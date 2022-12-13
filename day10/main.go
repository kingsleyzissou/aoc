package main

import (
	h "aoc/helpers"
	"fmt"
	"strings"
)

type register struct {
	x        int
	cycles   int
	strength map[int]int
}

func (r *register) addX(v int) {
	r.x += v
}

func (r *register) getSignalStrength() int {
	sum := 0
	for k, v := range r.strength {
		sum += k * v
	}
	return sum
}

func in(c int, checkpoints []int) bool {
	for _, v := range checkpoints {
		if v == c {
			return true
		}
	}
	return false
}

func newRegister() register {
	return register{
		x:        1,
		cycles:   0,
		strength: make(map[int]int),
	}
}

func handleCheckpoint(r *register) {
	checkpoints := []int{19, 59, 99, 139, 179, 219}
	if in(r.cycles, checkpoints) {
		r.strength[r.cycles+1] = r.x
	}
	r.cycles++
}

func (r *register) pixelCharacter() string {
	lt := r.cycles%40 <= r.x+1
	gt := r.cycles%40 >= r.x-1
	if gt && lt {
		return "#"
	}
	return "."
}

func (r *register) drawPixel() {
	if r.cycles%40 == 0 {
		fmt.Println()
	}
	fmt.Print(r.pixelCharacter())
	handleCheckpoint(r)
}

func (r *register) getRow() int {
	return int(r.cycles / 40)
}

func cycle(input []string) {
	r := newRegister()
	input = input[:len(input)-1]
	for _, line := range input {
		s := strings.Split(line, " ")
		if s[0] == "noop" {
			r.drawPixel()
			continue
		}
		r.drawPixel()
		r.drawPixel()
		r.addX(h.StringToInt(s[1]))
	}
	fmt.Println()
	fmt.Println("Signal strength: ", r.getSignalStrength())
}

func main() {
	contents := h.ReadFile()
	cycle(contents)
}
