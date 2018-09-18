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

func sliceEqual(a, b []uintptr) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
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
