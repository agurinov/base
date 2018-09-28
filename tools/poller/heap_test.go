package poller

import (
	"container/heap"
	"fmt"
	"sync"
	"testing"
)

func pollerHeapLen(t *testing.T, hp *pollerHeap, ready, pending int) {
	t.Run("ready", func(t *testing.T) {
		// ready via .Len()
		if n := hp.Len(); n != ready {
			t.Fatalf("hp.Len() -> Expected \"%d\", got \"%d\"", ready, n)
		}
	})

	t.Run("pending", func(t *testing.T) {
		// pending direct
		if n := len(hp.pending); n != pending {
			t.Fatalf("len(hp.pending) -> Expected \"%d\", got \"%d\"", pending, n)
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

	t.Run("Len", func(t *testing.T) {
		hp.ready = []uintptr{1, 2, 3}

		for i := 0; i < 5; i++ {
			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				t.Parallel()
				if l := hp.Len(); l != 3 {
					t.Fatalf("Len(). Expected %d, got %d", 3, l)
				}
			})
		}
	})

	t.Run("Pop", func(t *testing.T) {
		poller := MockPoller(
			[]uintptr{1, 2, 3, 4, 5, 6},
			nil,
			false,
		).(*mock)

		hp := &pollerHeap{
			pl:      sync.NewCond(&sync.Mutex{}),
			pending: make(map[uintptr]interface{}, 0),
			ready:   make([]uintptr, 0),
			poller:  poller, // 1 second wait,
		}
		heap.Init(hp)

		// fill
		hp.pending[1] = "foobar1"
		hp.pending[2] = "foobar2"
		hp.pending[3] = "foobar3"
		hp.pending[4] = "foobar4"
		hp.pending[5] = "foobar5"
		hp.pending[6] = "foobar6"

		// concurrent .Pop()
		t.Run("Parallel", func(t *testing.T) {
			for i := 0; i < 6; i++ {
				j := i
				t.Run(fmt.Sprintf("%d", j), func(t *testing.T) {
					t.Parallel()

					if v := heap.Pop(hp); v == nil {
						t.Fatalf("Pop(). Unexpected <nil>")
					}
				})
			}
		})

		// check heap state
		pollerHeapLen(t, hp, 0, 0)
		// if poller.invokes != 1 {
		// 	t.Fatalf("poller invokes. Expected %d, got %d", 1, poller.invokes)
		// }
	})

	t.Run("Push", func(t *testing.T) {
		// clear
		hp.ready = []uintptr{}
		hp.pending = map[uintptr]interface{}{}

		t.Run("real", func(t *testing.T) {
			heap.Push(hp, &HeapItem{Fd: uintptr(1), Value: "foobar"})
			pollerHeapLen(t, hp, 0, 1)
		})

		// clear
		hp.ready = []uintptr{}
		hp.pending = map[uintptr]interface{}{}

		t.Run("fake", func(t *testing.T) {
			heap.Push(hp, "foobar")
			pollerHeapLen(t, hp, 0, 0)
		})

		t.Run("poller/error", func(t *testing.T) {
			hp := &pollerHeap{
				pl:      sync.NewCond(&sync.Mutex{}),
				pending: make(map[uintptr]interface{}, 0),
				ready:   make([]uintptr, 0),
				poller:  MockPoller(nil, nil, true).(*mock),
			}
			heap.Init(hp)

			// real value, but error from poller -> 0
			heap.Push(hp, &HeapItem{Fd: uintptr(1), Value: "foobar"})
			pollerHeapLen(t, hp, 0, 0)
		})
	})
}

func TestHeapPrivate(t *testing.T) {
	heapInterface, _ := Heap()
	hp, _ := heapInterface.(*pollerHeap)

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
				hp.ready = tt.ready

				if l := hp.len(); l != tt.len {
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
				hp.ready = tt.ready

				if x := hp.less(tt.i, tt.j); x != tt.cond {
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
				hp.ready = tt.ready

				hp.swap(tt.i, tt.j)
				if !sliceEqual(hp.ready, tt.swapped) {
					t.Fatalf("swap(). Expected %v, got %v", tt.swapped, hp.ready)
				}
			})
		}
	})

	t.Run("poll", func(t *testing.T) {
		tableTests := []struct {
			poller  Interface
			invokes int
			re      []uintptr
			ce      []uintptr
		}{
			// in case error or e,pty we return hardcoded values (100, 500), (200,300) otherwise operation will block
			{MockPoller([]uintptr{1, 2}, []uintptr{3, 4}, true), 2, []uintptr{100, 500}, []uintptr{200, 300}}, // mock return empty because error
			{MockPoller([]uintptr{}, []uintptr{}, false), 2, []uintptr{100, 500}, []uintptr{200, 300}},        // mock return empty because empty (second continue)
			{MockPoller([]uintptr{1, 2}, []uintptr{3, 4}, false), 1, []uintptr{1, 2}, []uintptr{3, 4}},        // normal return in one polling
		}

		for i, tt := range tableTests {
			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				hp := &pollerHeap{
					pl:      sync.NewCond(&sync.Mutex{}),
					pending: make(map[uintptr]interface{}, 0),
					ready:   make([]uintptr, 0),
					poller:  tt.poller,
				}
				heap.Init(hp)

				re, ce := hp.poll()
				if invokes := tt.poller.(*mock).invokes; invokes != tt.invokes {
					t.Fatalf("invokes. Expected %v, got %v", tt.invokes, invokes)
				}
				if !sliceEqual(re, tt.re) {
					t.Fatalf("re. Expected %v, got %v", tt.re, re)
				}
				if !sliceEqual(ce, tt.ce) {
					t.Fatalf("ce. Expected %v, got %v", tt.ce, ce)
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
				hp.ready = tt.ready
				hp.pending = tt.pending

				// pop and check returned value and counters
				if x := hp.pop(); x != tt.x {
					t.Fatalf("pop(). Expected %q, got %q", tt.x, x)
				}
				pollerHeapLen(t, hp, tt.countReady, tt.countPending)
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
				hp.ready = tt.ready
				hp.pending = tt.pending

				// actualize and check returned value and counters
				hp.actualize(tt.pollerReady, tt.pollerClose)
				pollerHeapLen(t, hp, tt.countReady, tt.countPending)
			})
		}
	})
}
