package chronometer

import (
	"fmt"
	"sort"
	"strings"
)

type Chronometer map[string]*Node

func New() *Chronometer {
	m := make(map[string]*Node, 0)
	ch := Chronometer(m)
	return &ch
}

// Enter starts measuring node
func (ch *Chronometer) Enter(name string) {
	(*ch)[name] = NewNode()
}

// Exit stops measuring node
func (ch *Chronometer) Exit(name string) {
	if node, ok := (*ch)[name]; ok {
		node.Exit()
	}
}

// AddNode adds node created outside .Enter() to chronometer collection
func (ch *Chronometer) AddNode(name string, node *Node) {
	(*ch)[name] = node
}

func (ch *Chronometer) String() string {
	nodes := []string{}

	for name, node := range *ch {
		d := node.Duration()
		if d != 0 {
			nodes = append(nodes, fmt.Sprintf("%s: %s", name, d))
		}
	}

	sort.Strings(nodes)
	return strings.Join(nodes, ", ")
}
