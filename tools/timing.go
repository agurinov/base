package tools

import (
	"time"
)

type timing struct {
	enter map[string]time.Time
	exit  map[string]time.Time
}

func NewTiming() *timing {
	return &timing{
		enter: make(map[string]time.Time),
		exit:  make(map[string]time.Time),
	}
}

// Enter starts measuring node
func (t *timing) Enter(node string) {
	t.enter[node] = time.Now()
}

// Exit stops measuring node
func (t *timing) Exit(node string) {
	t.exit[node] = time.Now()
}

// Duration returns measured node life time
// returns time.Duration
func (t *timing) Duration(node string) time.Duration {
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
func (t *timing) Total(node string) time.Duration {
	return 0
}
