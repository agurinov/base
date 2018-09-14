package tools

import (
	"fmt"
	"strings"
	"time"
)

type Timing struct {
	enter map[string]time.Time
	exit  map[string]time.Time
}

func NewTiming() *Timing {
	return &Timing{
		enter: make(map[string]time.Time),
		exit:  make(map[string]time.Time),
	}
}

// Enter starts measuring node
func (t *Timing) Enter(node string) {
	t.enter[node] = time.Now()
}

// Exit stops measuring node
func (t *Timing) Exit(node string) {
	t.exit[node] = time.Now()
}

// Duration returns measured node life time
// returns time.Duration
func (t *Timing) Duration(node string) time.Duration {
	entered, ok := t.enter[node]
	if !ok {
		return 0
	}

	exited, ok := t.exit[node]
	if !ok {
		return 0
	}

	return exited.Sub(entered)
}

// Total is time.Duration between first Enter and last Exit
// time beetween extremums
func (t *Timing) Total(node string) time.Duration {
	return 0
}

func (t *Timing) String() string {
	nodes := []string{}

	for node, _ := range t.enter {
		nodes = append(nodes, fmt.Sprintf("%s: %s", node, t.Duration(node)))
	}

	return strings.Join(nodes, ", ")
}
