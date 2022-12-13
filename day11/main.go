package main

import (
	h "aoc/helpers"
	"fmt"
	"sort"
	"strings"
)

type actions interface {
	inspect(item uint64) uint64
	bored(item uint64) uint64
	throw(item uint64, other monkey)
	push(item uint64)
}

type monkey struct {
	name        string
	inspections uint64
	divisor     uint64
	operation   func(item uint64) uint64
	test        func(item uint64) bool
	items       []uint64
	candidates  map[bool]*monkey
	actions
}

func equation(a, b, operator string) func(item uint64) uint64 {
	translate := func(item uint64, s string) uint64 {
		if s == "old" {
			return item
		}
		return h.StringToInt64(s)
	}
	switch operator {
	case "+":
		return func(item uint64) uint64 {
			return translate(item, a) + translate(item, b)
		}
	case "*":
		return func(item uint64) uint64 {
			return translate(item, a) * translate(item, b)
		}
	default:
		panic("undefined operation")
	}
}

func test(divisor uint64) func(item uint64) bool {
	return func(item uint64) bool {
		return item%divisor == 0
	}
}

func (m *monkey) push(item uint64) {
	m.items = append(m.items, item)
}

func (m *monkey) inspect(item uint64) uint64 {
	m.inspections++
	return m.operation(item)
}

func (m *monkey) bored(item uint64) uint64 {
	item = uint64(item / 3)
	return item
}

func (m *monkey) decide(b bool) *monkey {
	return m.candidates[b]
}

func (m *monkey) throw(item uint64, other *monkey) {
	other.push(item)
}

func (m *monkey) setItems(items []uint64) {
	m.items = items
}

func getFactorial(monkeys []*monkey) uint64 {
	var factorial uint64 = 1
	for _, m := range monkeys {
		factorial *= m.divisor
	}
	return factorial
}

func (m *monkey) do(decrease bool, factorial uint64) {
	for len(m.items) > 0 {
		item, items := (m.items)[0], (m.items)[1:]
		m.setItems(items)
		item = m.inspect(item)
		if decrease {
			item = m.bored(item)
		}
		// reduce number sizes
		item = item % factorial
		next := m.decide(m.test(item))
		m.throw(item, next)
	}

}

func newMonkey() monkey {
	return monkey{
		inspections: 0,
		candidates:  make(map[bool]*monkey),
	}
}

func newMonkeys() []*monkey {
	var monkeys []*monkey
	for i := 0; i < 8; i++ {
		m := newMonkey()
		m.name = fmt.Sprint(i)
		monkeys = append(monkeys, &m)
	}
	return monkeys
}

func createMonkeys(input []string, monkeys []*monkey) {
	var current *monkey
	for _, line := range input {
		if line == "" {
			continue
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Monkey") {
			s := strings.Split(line, " ")
			index := h.StringToInt(strings.TrimSuffix(s[1], ":"))
			current = monkeys[index]
		}
		if prefix := "Starting items: "; strings.HasPrefix(line, prefix) {
			its := strings.Split(strings.TrimPrefix(line, prefix), ", ")
			var items []uint64
			for _, item := range its {
				i := h.StringToInt64(item)
				items = append(items, i)
			}
			current.items = items
		}
		if prefix := "Operation: new = "; strings.HasPrefix(line, prefix) {
			rhs := strings.Split(strings.TrimPrefix(line, prefix), " ")
			current.operation = equation(rhs[0], rhs[2], rhs[1])
		}
		if prefix := "Test: "; strings.HasPrefix(line, prefix) {
			rhs := strings.Split(strings.TrimPrefix(line, prefix), " ")
			current.divisor = h.StringToInt64(rhs[len(rhs)-1])
			current.test = test(current.divisor)
		}
		if prefix := "If true: "; strings.HasPrefix(line, prefix) {
			rhs := strings.Split(strings.TrimPrefix(line, prefix), " ")
			index := h.StringToInt(rhs[len(rhs)-1])
			candidate := monkeys[index]
			current.candidates[true] = candidate
		}
		if prefix := "If false: "; strings.HasPrefix(line, prefix) {
			rhs := strings.Split(strings.TrimPrefix(line, prefix), " ")
			index := h.StringToInt(rhs[len(rhs)-1])
			candidate := monkeys[index]
			current.candidates[false] = candidate
		}
	}
}

func rounds(monkeys []*monkey, rounds int) {
	decrease := false
	if rounds == 20 {
		decrease = true
	}
	// help reduce the number sizes
	// thanks to:
	// https://nickymeuleman.netlify.app/garden/aoc2022-day11#part-2
	factorial := getFactorial(monkeys)
	for i := 1; i <= rounds; i++ {
		for _, m := range monkeys {
			m.do(decrease, factorial)
		}
	}
}

func main() {
	contents := h.ReadFile()
	monkeys := newMonkeys()
	createMonkeys(contents, monkeys)
	rounds(monkeys, 10000) // edit rounds to 20 for part 1

	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].inspections > monkeys[j].inspections
	})

	fmt.Println("Busiest: ", monkeys[0].inspections)
	fmt.Println("Busiest: ", monkeys[1].inspections)
	fmt.Println(monkeys[0].inspections * monkeys[1].inspections)
}
