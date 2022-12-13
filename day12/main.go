package main

import (
	h "aoc/helpers"
	"fmt"
	"sort"
	"strings"
)

var letters = map[string]int{
	"a": 0, "b": 1, "c": 2, "d": 3, "e": 4, "f": 5, "g": 6, "h": 7, "i": 8, "j": 9, "k": 10, "l": 11, "m": 12, "n": 13, "o": 14, "p": 15,
	"q": 16, "r": 17, "s": 18, "t": 19, "u": 20, "v": 21, "w": 22, "x": 23, "y": 24, "z": 25,
	"S": 0, "E": 25,
}

type position struct {
	x int
	y int
}

func (p *position) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
}

type Edge struct {
	height   int
	possible bool
}

type Node struct {
	value      int
	char       string
	position   *position
	neighbours map[*Node]*Edge
}

func (n1 *Node) addNeighbour(n2 *Node, edge *Edge) {
	n1.neighbours[n2] = edge
}

func (n *Node) equals(o *Node) bool {
	return n.position.x == o.position.x &&
		n.position.y == o.position.y
}

func newNode(x, y int, char string) *Node {
	v, ok := letters[char]
	if !ok {
		panic("not okay")
	}
	return &Node{
		char:       char,
		value:      v,
		position:   &position{x: x, y: y},
		neighbours: make(map[*Node]*Edge),
	}
}

type Graph struct {
	start     *position
	end       *position
	nodes     map[*position]*Node
	edges     map[*Node][]*Edge
	positions map[*position]*position
}

func (g *Graph) addNode(node *Node) {
	if node.char == "S" {
		g.start = node.position
	}
	if node.char == "E" {
		g.end = node.position
	}
	g.nodes[node.position] = node
}

func possible(cost int) bool {
	return cost <= 1
}

func newEdge(source, destination *Node) *Edge {
	if source == destination {
		panic("source can't be the same as destination")
	}
	height := destination.value - source.value
	edge := &Edge{
		height:   height,
		possible: possible(height),
	}
	source.addNeighbour(destination, edge)
	return edge
}

func (g *Graph) addEdge(n1, n2 *Node) {
	g.edges[n1] = append(g.edges[n1], newEdge(n1, n2))
	g.edges[n2] = append(g.edges[n2], newEdge(n1, n2))
}

func (g *Graph) addPosition(p *position) {
	g.positions[p] = p
}

func (g *Graph) getPosition(x, y int) *position {
	for k := range g.positions {
		if k.x == x && k.y == y {
			return k
		}
	}
	return &position{x: x, y: y}
}

func (g *Graph) getNode(x, y int) *Node {
	p := g.getPosition(x, y)
	return g.nodes[p]
}

func newGraph() *Graph {
	return &Graph{
		edges:     make(map[*Node][]*Edge),
		nodes:     make(map[*position]*Node),
		positions: make(map[*position]*position),
	}
}

type Vertex struct {
	node     *Node
	previous *Vertex
	distance int
}

func (vertex *Vertex) queueNeighbours(queue *Queue, graph *Graph, visited map[*position]bool) {
	distance := vertex.distance
	for k, v := range vertex.node.neighbours {
		if visited[k.position] || !v.possible {
			continue
		}
		next := newVertex(k, vertex, distance+1)
		queue.Enqueue(next)
	}
}

func newVertex(node *Node, previous *Vertex, distance int) *Vertex {
	return &Vertex{
		node:     node,
		previous: previous,
		distance: distance,
	}
}

type Queue struct {
	items []*Vertex
}

func (q *Queue) Contains(item *Vertex) bool {
	for _, i := range q.items {
		if i.node.equals(item.node) {
			return true
		}
	}
	return false
}

func (q *Queue) Enqueue(item *Vertex) {
	if q.Contains(item) {
		return
	}
	q.items = append(q.items, item)
	sort.Slice(q.items, func(i, j int) bool {
		return q.items[i].distance < q.items[j].distance
	})
}

func (q *Queue) Dequeue() *Vertex {
	item := q.items[0]
	q.items = q.items[1:len(q.items)]
	return item
}

func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

func (queue *Queue) process(graph *Graph, start, end *Node) (*Vertex, error) {
	visited := make(map[*position]bool)
	for !queue.IsEmpty() {
		vertex := queue.Dequeue()
		if vertex.node.equals(end) {
			return vertex, nil
		}
		visited[vertex.node.position] = true
		vertex.queueNeighbours(queue, graph, visited)
	}
	return nil, fmt.Errorf("Unable to find path")
}

func NewQueue(item *Node) *Queue {
	vertex := newVertex(item, nil, 0)
	return &Queue{items: []*Vertex{vertex}}
}

func shortestPath(graph *Graph) (*Vertex, error) {
	end := graph.getNode(graph.end.x, graph.end.y)
	start := graph.getNode(graph.start.x, graph.start.y)
	return NewQueue(start).process(graph, start, end)
}

func kShortestPaths(graph *Graph) []*Vertex {
	var kPaths []*Vertex
	for _, node := range graph.nodes {
		if node.char == "a" {
			graph.start = node.position
			k, err := shortestPath(graph)
			if err != nil {
				continue
			}
			kPaths = append(kPaths, k)
		}
	}
	sort.Slice(kPaths, func(i, j int) bool {
		return kPaths[i].distance < kPaths[j].distance
	})
	return kPaths
}

func createGraph(input []string, graph *Graph) {
	input = input[:len(input)-1]
	for y, line := range input {
		str := strings.Split(line, "")
		for x, char := range str {
			node := newNode(x, y, char)
			graph.addNode(node)
			graph.addPosition(node.position)
			if x != 0 {
				left := graph.getNode(x-1, y)
				graph.addEdge(left, node)
				graph.addEdge(node, left) // add reciprocal edge
			}
			if y != 0 {
				up := graph.getNode(x, y-1)
				graph.addEdge(up, node)
				graph.addEdge(node, up) // add reciprocal edge
			}
		}
	}
}

func main() {
	var graph = newGraph()
	contents := h.ReadFile()
	createGraph(contents, graph)
	shortest, _ := shortestPath(graph)
	fmt.Println("Shortest path: ", shortest.distance)
	kShortest := kShortestPaths(graph)
	fmt.Println("Shortest path: ", kShortest[0].distance)
}
