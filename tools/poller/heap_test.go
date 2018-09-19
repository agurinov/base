package poller

import (
	"fmt"
	"testing"
)

func pollerHeapLen(t *testing.T, heap *pollerHeap, ready, pending int) {
	t.Run("ready/Len", func(t *testing.T) {
		// ready via .Len()
		if n := heap.Len(); n != ready {
			t.Fatalf("heap.Len() -> Expected \"%d\", got \"%d\"", ready, n)
		}
	})

	t.Run("ready/direct", func(t *testing.T) {
		// ready direct
		if n := len(heap.ready); n != ready {
			t.Fatalf("len(heap.ready) -> Expected \"%d\", got \"%d\"", ready, n)
		}
	})

	t.Run("pending/direct", func(t *testing.T) {
		// pending direct
		if n := len(heap.pending); n != pending {
			t.Fatalf("len(heap.pending) -> Expected \"%d\", got \"%d\"", pending, n)
		}
	})
}

func TestHeap(t *testing.T) {
	heapInterface, err := Heap()
	heap, ok := heapInterface.(*pollerHeap)

	t.Run("New", func(t *testing.T) {
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("Unexpected heap type (expected *pollerHeap)")
		}
		pollerHeapLen(t, heap, 0, 0)
	})

	t.Run("pop", func(t *testing.T) {
		tableTests := []struct {
			ready        []uintptr               // ready slice for heap (initial)
			pending      map[uintptr]interface{} // pending for heap (initial)
			x            interface{}             // value returned by pop()
			countReady   int                     // ready slice for heap (count)
			countPending int                     // pending for heap (count)

		}{
			{[]uintptr{}, map[uintptr]interface{}{}, nil, 0, 0},
			{[]uintptr{}, map[uintptr]interface{}{1: "some"}, nil, 0, 1},
			{[]uintptr{1, 2, 3}, map[uintptr]interface{}{4: "some"}, nil, 0, 1},
			{[]uintptr{1, 2, 3}, map[uintptr]interface{}{}, nil, 0, 0},
			{[]uintptr{1, 2, 3}, map[uintptr]interface{}{2: "some"}, "some", 1, 0},
			{[]uintptr{1, 2, 3}, map[uintptr]interface{}{4: "some", 3: "foobar"}, "foobar", 0, 1},
		}

		for i, tt := range tableTests {
			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				// fill heap
				heap.ready = tt.ready
				heap.pending = tt.pending

				// pop and check returned value and counters
				if x := heap.pop(); x != tt.x {
					t.Fatalf("pop(). Expected %q, got %q", tt.x, x)
				}
				pollerHeapLen(t, heap, tt.countReady, tt.countPending)
			})
		}
	})

	t.Run("actualize", func(t *testing.T) {
		tableTests := []struct {
			ready        []uintptr               // ready slice for heap (initial)
			pending      map[uintptr]interface{} // pending for heap (initial)
			pollerReady  []uintptr               // ready events from poller
			pollerClose  []uintptr               // close events from poller
			countReady   int                     // ready slice for heap (count)
			countPending int                     // pending for heap (count)

		}{
			{[]uintptr{}, map[uintptr]interface{}{}, []uintptr{}, []uintptr{}, 0, 0},
			{[]uintptr{}, map[uintptr]interface{}{}, []uintptr{1, 2, 3}, []uintptr{2, 3, 4, 5}, 0, 0},
			{[]uintptr{1, 2}, map[uintptr]interface{}{4: "some"}, []uintptr{3}, []uintptr{4}, 2, 0},
			{[]uintptr{1, 2}, map[uintptr]interface{}{4: "some", 6: "some6"}, []uintptr{4, 6, 7}, []uintptr{4, 2, 5}, 2, 1},
			{[]uintptr{1, 2, 3, 4}, map[uintptr]interface{}{1: "some", 2: "some", 3: "some", 4: "some"}, []uintptr{}, []uintptr{4, 2, 1}, 1, 1},
		}

		for i, tt := range tableTests {
			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				// fill heap
				heap.ready = tt.ready
				heap.pending = tt.pending

				// actualize and check returned value and counters
				heap.actualize(tt.pollerReady, tt.pollerClose)
				pollerHeapLen(t, heap, tt.countReady, tt.countPending)
			})
		}
	})

	t.Run("Pop", func(t *testing.T) {
		// TODO mocking poller
	})

	t.Run("Push", func(t *testing.T) {
		// clear
		heap.ready = []uintptr{}
		heap.pending = map[uintptr]interface{}{}

		t.Run("real", func(t *testing.T) {
			heap.Push(&HeapItem{Fd: uintptr(1), Value: "foobar"})
			pollerHeapLen(t, heap, 0, 1)
		})

		t.Run("fake", func(t *testing.T) {
			heap.Push("foobar")
			pollerHeapLen(t, heap, 0, 1)
		})

		t.Run("poller/error", func(t *testing.T) {
			// TODO poller error
		})
	})

}
