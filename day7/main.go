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
	// for parent != nil {
	// parent = parent.parent
	*parent.size += size
	updateSize(size, parent.parent)
	// }
}

func cd(dir string, cwd *fsNode) *fsNode {
	if dir == "/" {
		return newRoot()
	}

	if dir == ".." {
		if cwd.parent != nil {

		}
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

func ls(s []string, cwd *fsNode) *fsNode {
	for _, line := range s {
		if strings.HasPrefix(line, "$") || line == "" {
			return cwd
		}
		entry := fsEntryToNode(line, cwd)
		cwd.addNode(entry)
	}
	return cwd
}

func execute(s string, lines []string, index int, cwd fsNode) *fsNode {
	cmd := strings.Split(s, " ")
	if strings.HasPrefix(cmd[1], "ls") {
		return ls(lines[index:], &cwd)
	}
	return cd(cmd[2], &cwd)
}

func isCommand(s string) bool {
	return strings.HasPrefix(s, "$")
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

func main() {
	cwd := fsNode{}
	lines := readFile()
	for index, line := range lines {
		if isCommand(line) {
			res := execute(line, lines, index+1, cwd)
			cwd = *res
		}
	}

	// go back to the root
	root := cwd.getRoot()

	// part one
	sum := getDirsUnderValue(root, 100000)
	fmt.Println("Sum: ", sum)

	//part two
	unused := 70000000 - *root.size
	needed := 30000000 - unused
	entry := getSmallestEntryOverValue(root, needed)
	fmt.Println("Entry: ", *entry.size)
}
