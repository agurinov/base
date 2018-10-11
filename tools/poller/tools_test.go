package poller

import (
	"fmt"
	"reflect"
	"testing"
)

type fakeEvent struct {
	fd uintptr
}

func (ev fakeEvent) Fd() uintptr {
	return ev.fd
}

func TestEventsToFds(t *testing.T) {
	tableTests := []struct {
		events []Event
		out    []uintptr
	}{
		{[]Event{}, []uintptr{}},
		{[]Event{fakeEvent{1}, fakeEvent{2}, fakeEvent{3}}, []uintptr{1, 2, 3}},
	}

	for i, tt := range tableTests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if out := EventsToFds(tt.events...); !reflect.DeepEqual(out, tt.out) {
				t.Fatalf("Expected %v, got %v", tt.out, out)
			}
		})
	}
}

func TestPendingFilterClosed(t *testing.T) {
	tableTests := []struct {
		pending []*HeapItem
		close   []uintptr
		out     []*HeapItem
	}{
		// nil values
		{nil, []uintptr{1, 2, 3}, nil},
		{[]*HeapItem{}, []uintptr{1, 2, 3}, nil},
		// real values
		{[]*HeapItem{&HeapItem{Fd: 1, ready: true}, &HeapItem{Fd: 4}}, []uintptr{1, 2, 3}, []*HeapItem{&HeapItem{Fd: 4}}},
		{[]*HeapItem{&HeapItem{Fd: 1, ready: true}, &HeapItem{Fd: 4}}, []uintptr{5, 6, 7}, []*HeapItem{&HeapItem{Fd: 1, ready: true}, &HeapItem{Fd: 4}}},
	}

	for i, tt := range tableTests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if out := pendingFilterClosed(tt.pending, tt.close); !reflect.DeepEqual(out, tt.out) {
				t.Fatalf("Expected %v, got %v", tt.out, out)
			}
		})
	}
}

func TestPendingMapReady(t *testing.T) {
	tableTests := []struct {
		pending []*HeapItem
		ready   []uintptr
		out     []*HeapItem
	}{
		// nil values
		{nil, []uintptr{1, 2, 3}, nil},
		{[]*HeapItem{}, []uintptr{1, 2, 3}, nil},
		// real values
		{[]*HeapItem{&HeapItem{Fd: 1, ready: true}, &HeapItem{Fd: 4}}, []uintptr{1, 2, 3, 4}, []*HeapItem{&HeapItem{Fd: 1, ready: true}, &HeapItem{Fd: 4, ready: true}}},
	}

	for i, tt := range tableTests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if out := pendingMapReady(tt.pending, tt.ready); !reflect.DeepEqual(out, tt.out) {
				t.Fatalf("Expected %v, got %v", tt.out, out)
			}
		})
	}
}
