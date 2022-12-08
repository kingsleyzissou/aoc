package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type forest struct {
	trees [][]tree
}

type neighbours map[string]int

func newNeighbours() neighbours {
	return neighbours{
		"up":    1,
		"down":  1,
		"left":  1,
		"right": 1,
	}
}

type tree struct {
	row        int
	col        int
	visible    bool
	size       int
	neighbours neighbours
}

func (t tree) viewingScore() int {
	prod := 1
	for _, v := range t.neighbours {
		prod *= v
	}
	return prod
}

func (t *tree) setViewingDistance(key string, value int) {
	if value <= 1 {
		return
	}
	t.neighbours[key] = value
}

func stringToInt(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}

func isEdge(curr, min, max int) bool {
	return (curr == min || curr == max)
}

func buildTrees(edge bool, index int, row []string) []tree {
	var trees []tree
	for i, item := range row {
		visible := isEdge(i, 0, len(row)-1) || edge
		trees = append(trees, tree{
			col:        i,
			row:        index,
			visible:    visible,
			size:       stringToInt(item),
			neighbours: newNeighbours(),
		})
	}
	return trees
}

func buildForest(input []string, forest *forest) {
	// do -1 to get rid of empty line break at end
	input = input[:len(input)-1]
	for index, line := range input {
		// if line == "" {
		// break
		// }
		row := strings.Split(line, "")
		edge := isEdge(index, 0, len(input)-1)
		trees := buildTrees(edge, index, row)
		forest.trees = append(forest.trees, trees)
	}
}

// create an iterator function so that since all the
// loop conditions are slightly different. This gives us
// a reusable function that we can then run the same conditions
// inside of (since the conditions we are looking for don't change too much)
func (f *forest) iterator(curr *tree, direction string) func() chan tree {
	if direction == "up" {
		return func() chan tree {
			ch := make(chan tree)
			go func() {
				for i := curr.row - 1; i >= 0; i-- {
					ch <- f.trees[i][curr.col]
				}
				close(ch)
			}()
			return ch
		}
	}
	if direction == "down" {
		return func() chan tree {
			ch := make(chan tree)
			go func() {
				for i := curr.row + 1; i < len(f.trees); i++ {
					ch <- f.trees[i][curr.col]
				}
				close(ch)
			}()
			return ch
		}
	}
	if direction == "left" {
		return func() chan tree {
			ch := make(chan tree)
			go func() {
				for j := curr.col - 1; j >= 0; j-- {
					ch <- f.trees[curr.row][j]
				}
				close(ch)
			}()
			return ch
		}
	}
	return func() chan tree {
		ch := make(chan tree)
		go func() {
			for j := curr.col + 1; j < len(f.trees[0]); j++ {
				ch <- f.trees[curr.row][j]
			}
			close(ch)
		}()
		return ch
	}
}

func (f *forest) travel(curr *tree, direction string) {
	distance := 0
	visible := true
	it := f.iterator(curr, direction)
	for t := range it() {
		distance++
		if t.size >= curr.size {
			visible = false
			break
		}
	}
	curr.visible = curr.visible || visible
	curr.setViewingDistance(direction, distance)
}

func (f *forest) search(callback func(i, j int)) {
	for i := 0; i < len(f.trees); i++ {
		for j := 0; j < len(f.trees[0]); j++ {
			f.travel(&f.trees[i][j], "up")
			f.travel(&f.trees[i][j], "down")
			f.travel(&f.trees[i][j], "left")
			f.travel(&f.trees[i][j], "right")
			callback(i, j)
		}
	}
}

func readFile() []string {
	file, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	return strings.Split(string(file), "\n")
}

func main() {
	forest := &forest{}
	contents := readFile()
	buildForest(contents, forest)

	sum, score := 0, 0
	forest.search(func(i, j int) {
		t := forest.trees[i][j]
		// part one
		if t.visible {
			sum++
		}
		// part two
		vs := t.viewingScore()
		if vs > score {
			score = vs
		}
	})
	fmt.Println("Visible trees: ", sum)
	fmt.Println("Max viewing score: ", score)
}
