package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type fsNode struct {
	name     string
	kind     string
	size     *int
	parent   *fsNode
	children *[]fsNode
}

func (fs *fsNode) findChild(name string) *fsNode {
	for _, child := range *fs.children {
		if child.name == name {
			return &child
		}
	}
	return nil
}

func (fs fsNode) getRoot() fsNode {
	if fs.name == "/" {
		return fs
	}
	return fs.parent.getRoot()
}

func (fs *fsNode) addNode(node fsNode) {
	if fs.children == nil {
		fs.children = new([]fsNode)
	}
	*fs.children = append(*fs.children, node)
}

func (fs *fsNode) traverse(callback func(node *fsNode)) {
	callback(fs)
	for _, c := range *fs.children {
		if c.children != nil {
			c.traverse(callback)
		}
	}
}

func newRoot() *fsNode {
	return &fsNode{
		name:     "/",
		parent:   nil,
		size:     new(int),
		children: new([]fsNode),
	}

}

func updateSize(size int, parent *fsNode) {
	if parent == nil {
		return
	}
	*parent.size += size
	updateSize(size, parent.parent)
}

func nextNode(dir string, cwd *fsNode) *fsNode {
	if dir == "/" {
		return newRoot()
	}

	if dir == ".." {
		return cwd.parent
	}

	child := cwd.findChild(dir)
	if child == nil {
		panic(fmt.Sprintf("Some massive error happened here: %s (%s)", dir, cwd.name))
	}

	return child
}

func fsEntryToNode(s string, cwd *fsNode) fsNode {
	item := strings.Split(s, " ")
	if strings.HasPrefix(item[0], "dir") {
		return fsNode{
			name:     item[1],
			kind:     "d",
			parent:   cwd,
			size:     new(int),
			children: new([]fsNode),
		}
	}

	size, err := strconv.Atoi(item[0])
	if err != nil {
		panic(err)
	}

	entry := fsNode{
		name:     item[1],
		kind:     "f",
		parent:   cwd,
		size:     &size,
		children: nil,
	}

	// udpate the size through the
	// tree recursively. This way we
	// don't have to traverse the tree
	// again at the end to calculate the
	// size of each dir
	updateSize(size, cwd)

	return entry
}

func readFile() []string {
	file, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	return strings.Split(string(file), "\n")
}

func getDirsUnderValue(root fsNode, limit int) int {
	sum := 0
	root.traverse(func(node *fsNode) {
		if node.kind == "d" && *node.size <= limit {
			sum += *node.size
		}
	})
	return sum
}

func getSmallestEntryOverValue(root fsNode, threshold int) fsNode {
	var entries []fsNode
	root.traverse(func(node *fsNode) {
		if node.kind == "d" && *node.size >= threshold {
			entries = append(entries, *node)
		}
	})

	// sort the slice to get the smallest element
	// that satisfies the criteria
	sort.Slice(entries, func(i, j int) bool {
		return *entries[i].size < *entries[j].size
	})

	return entries[0]
}

func parseInput(input []string) *fsNode {
	cwd := &fsNode{}
	for _, line := range input {
		if line == "" {
			break
		}
		args := strings.Split(line, " ")
		if args[1] == "cd" {
			cwd = nextNode(args[2], cwd)
		} else if args[1] == "ls" {
			continue
		} else {
			entry := fsEntryToNode(line, cwd)
			cwd.addNode(entry)
		}
	}
	return cwd
}

func main() {
	// parse the input and build the
	// fs tree
	tree := parseInput(readFile())

	// go back to the root
	// so we can traverse the tree
	root := tree.getRoot()

	// part one
	sum := getDirsUnderValue(root, 100000)
	fmt.Println("Sum: ", sum)

	//part two
	unused := 70000000 - *root.size
	needed := 30000000 - unused
	entry := getSmallestEntryOverValue(root, needed)
	fmt.Println("Entry: ", *entry.size)
}
