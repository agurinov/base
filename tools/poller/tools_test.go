package poller

import (
	"fmt"
	"testing"
)

type fakeEvent struct {
	fd uintptr
}

func (ev fakeEvent) Fd() uintptr {
	return ev.fd
}

func TestSliceEqual(t *testing.T) {
	tableTests := []struct {
		a  []uintptr
		b  []uintptr
		eq bool
	}{
		{[]uintptr{}, []uintptr{}, true},
		{[]uintptr{1}, []uintptr{2}, false},
		{[]uintptr{1}, []uintptr{1}, true},
		{[]uintptr{1, 2}, []uintptr{2, 1}, false},
		{[]uintptr{1, 2}, []uintptr{1, 2}, true},
	}

	for i, tt := range tableTests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if eq := sliceEqual(tt.a, tt.b); eq != tt.eq {
				t.Fatalf(".sliceEqual() Expected %t, got %t", tt.eq, eq)
			}
		})
	}
}

func TestEventsToFds(t *testing.T) {
	tableTests := []struct {
		events []Event
		fds    []uintptr
	}{
		{[]Event{}, []uintptr{}},
		{[]Event{fakeEvent{1}, fakeEvent{2}, fakeEvent{3}}, []uintptr{1, 2, 3}},
	}

	for i, tt := range tableTests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if fds := EventsToFds(tt.events...); !sliceEqual(fds, tt.fds) {
				t.Fatalf("Expected %v, got %v", tt.fds, fds)
			}
		})
	}
}
