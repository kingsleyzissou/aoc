package main

import (
	"fmt"
	"strings"

	h "aoc/helpers"
)

type position struct {
	x int
	y int
}

func (p *position) horizontal(movement int) {
	p.x += movement
}

func (p *position) vertical(movement int) {
	p.y += movement
}

func (p *position) distance(o *position) position {
	return position{
		x: p.x - o.x,
		y: p.y - o.y,
	}
}

type distance position

func (d *distance) touching() bool {
	return h.Absolute(d.x) < 2 && h.Absolute(d.y) < 2
}

type knot struct {
	position position
	visited  map[position]int
}

func newKnot() knot {
	p := position{
		x: 0,
		y: 0,
	}
	v := make(map[position]int)
	v[p] = 1
	return knot{
		position: p,
		visited:  v,
	}
}

func (k *knot) visit() {
	_, ok := k.visited[k.position]
	if !ok {
		k.visited[k.position] = 1
	}
}

func (tail *knot) move(head *position) {
	diff := distance(head.distance(&tail.position))
	if diff.touching() {
		return
	}
	travelY := h.GetSign(diff.y)
	travelX := h.GetSign(diff.x)
	tail.position.horizontal(travelX)
	tail.position.vertical(travelY)
	tail.visit()
}

func step(direction string, steps int, knots *[]*knot, tail *knot) {
	head := (*knots)[0]
	for count := 0; count < steps; count++ {
		switch direction {
		case "U":
			head.position.vertical(+1)
		case "D":
			head.position.vertical(-1)
		case "L":
			head.position.horizontal(-1)
		case "R":
			head.position.horizontal(+1)
		default:
			panic("Unspecificed direction")
		}
		for i, k := range *knots {
			if i == 0 || i == 9 {
				continue
			}
			k.move(&(*knots)[i-1].position)
		}
		prev := (*knots)[len(*knots)-2]
		tail.move(&prev.position)
	}
}

func initKnots() []*knot {
	var knots []*knot
	for i := 0; i < 10; i++ { // change to i < 2 for part one
		k := newKnot()
		knots = append(knots, &k)
	}
	return knots
}

func (k *knot) count() int {
	return len(k.visited)
}

func start(input []string) {
	knots := initKnots()
	tail := knots[8] // change to 2 for part one
	input = input[:len(input)-1]
	for _, line := range input {
		s := strings.Split(line, " ")
		step(s[0], h.StringToInt(s[1]), &knots, tail)
	}
	count := tail.count()
	fmt.Println("Tail has visited: ", count, " spots")
}

func main() {
	contents := h.ReadFile()
	start(contents)
}
