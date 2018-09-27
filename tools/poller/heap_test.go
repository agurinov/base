package poller

import (
	"container/heap"
	"fmt"
	"sync"
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

func TestHeapPublic(t *testing.T) {
	heapInterface, err := Heap()
	hp, ok := heapInterface.(*pollerHeap)

	t.Run("New", func(t *testing.T) {
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("Unexpected heap type (expected *pollerHeap)")
		}
		pollerHeapLen(t, hp, 0, 0)
	})

	t.Run("Pop", func(t *testing.T) {
		hp := &pollerHeap{
			pl:      sync.NewCond(&sync.Mutex{}),
			pending: make(map[uintptr]interface{}, 0),
			ready:   make([]uintptr, 0),
			poller:  &mock{}, // 2 seconds wait,
		}
		heap.Init(hp)
		hp.pending[1] = "foobar1"
		hp.pending[2] = "foobar2"
		hp.pending[3] = "foobar3"
		hp.pending[4] = "foobar4"
		hp.pending[5] = "foobar5"
		hp.pending[6] = "foobar6"

		wg := sync.WaitGroup{}
		wg.Add(6)

		for i := 0; i < 6; i++ {
			go func(i int) {
				t.Logf("WINK (%d): %+v", i, heap.Pop(hp))
				wg.Done()
			}(i)
		}

		wg.Wait()
	})

	t.Run("Push", func(t *testing.T) {
		// clear
		hp.ready = []uintptr{}
		hp.pending = map[uintptr]interface{}{}

		t.Run("real", func(t *testing.T) {
			heap.Push(hp, &HeapItem{Fd: uintptr(1), Value: "foobar"})
			pollerHeapLen(t, hp, 0, 1)
		})

		t.Run("fake", func(t *testing.T) {
			heap.Push(hp, "foobar")
			pollerHeapLen(t, hp, 0, 1)
		})

		t.Run("poller/error", func(t *testing.T) {
			// TODO poller error
		})
	})
}

func TestHeapPrivate(t *testing.T) {
	heapInterface, _ := Heap()
	heap, _ := heapInterface.(*pollerHeap)

	t.Run("len", func(t *testing.T) {
		tableTests := []struct {
			ready []uintptr // ready slice for heap (initial)
			len   int
		}{
			{[]uintptr{}, 0},
			{[]uintptr{1, 2}, 2},
			{[]uintptr{1, 2, 3}, 3},
		}

		for i, tt := range tableTests {
			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				// fill heap
				heap.ready = tt.ready

				if l := heap.len(); l != tt.len {
					t.Fatalf("len(). Expected %d, got %d", tt.len, l)
				}
			})
		}
	})

	t.Run("less", func(t *testing.T) {
		tableTests := []struct {
			ready []uintptr // ready slice for heap (initial)
			i, j  int
			cond  bool
		}{
			// {[]uintptr{}, 0, 1, false},
			{[]uintptr{1, 2}, 0, 1, true},
			{[]uintptr{1, 2, 3}, 2, 1, false},
		}

		for i, tt := range tableTests {
			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				// fill heap
				heap.ready = tt.ready

				if x := heap.less(tt.i, tt.j); x != tt.cond {
					t.Fatalf("less(i, j). Expected %t, got %t", tt.cond, x)
				}
			})
		}
	})

	t.Run("swap", func(t *testing.T) {
		tableTests := []struct {
			ready   []uintptr // ready slice for heap (initial)
			i, j    int
			swapped []uintptr
		}{
			{[]uintptr{}, 0, 1, []uintptr{}},
			{[]uintptr{1}, 0, 1, []uintptr{1}},
			{[]uintptr{1, 2}, 0, 1, []uintptr{2, 1}},
			{[]uintptr{1, 2, 3}, 0, 2, []uintptr{3, 2, 1}},
		}

		for i, tt := range tableTests {
			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				// fill heap
				heap.ready = tt.ready

				heap.swap(tt.i, tt.j)
				if !sliceEqual(heap.ready, tt.swapped) {
					t.Fatalf("swap(). Expected %v, got %v", tt.swapped, heap.ready)
				}
			})
		}
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
}
