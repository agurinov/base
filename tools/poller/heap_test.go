package poller

import (
	"container/heap"
	"fmt"
	"reflect"
	"testing"
)

func pollerHeapState(t *testing.T, hp *pollerHeap, pending int) {
	t.Run("pending", func(t *testing.T) {
		// pending private
		if n := hp.len(); n != pending {
			t.Fatalf("hp.len() -> Expected \"%d\", got \"%d\"", pending, n)
		}
		// pending direct
		if n := len(hp.pending); n != pending {
			t.Fatalf("len(hp.pending) -> Expected \"%d\", got \"%d\"", pending, n)
		}
	})
}

func TestHeap(t *testing.T) {
	t.Run("Heap", func(t *testing.T) {
		heapInterface, err := Heap()
		hp, ok := heapInterface.(*pollerHeap)

		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("hp.type -> Expected *pollerHeap")
		}
		pollerHeapState(t, hp, 0)
	})

	t.Run("HeapWithPoller", func(t *testing.T) {
		heapInterface, err := HeapWithPoller(MockPoller(nil, nil, false))
		hp, ok := heapInterface.(*pollerHeap)

		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("hp.type -> Expected *mock")
		}
		pollerHeapState(t, hp, 0)
	})
}

func TestHeapPublic(t *testing.T) {
	heapInterface, _ := Heap()
	hp, _ := heapInterface.(*pollerHeap)

	t.Run("Len", func(t *testing.T) {
		hp.pending = []*HeapItem{
			&HeapItem{Fd: 1},
			&HeapItem{Fd: 2},
			&HeapItem{Fd: 3},
		}

		for i := 0; i < 5; i++ {
			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				t.Parallel()
				if l := hp.Len(); l != 3 {
					t.Fatalf("hp.Len() -> Expected %d, got %d", 3, l)
				}
			})
		}
	})

	// t.Run("Less", func(t *testing.T) {})

	// t.Run("Swap", func(t *testing.T) {})

	t.Run("Pop", func(t *testing.T) {
		infinityPoller := MockPoller(nil, nil, false).(*mock)
		infinityPoller.inifinity = true

		tableTests := []struct {
			poller  Interface
			pending []*HeapItem
		}{
			{
				MockPoller([]uintptr{1, 2, 3, 4, 5, 6}, nil, false),
				[]*HeapItem{
					&HeapItem{Fd: 1, Value: "1"},
					&HeapItem{Fd: 2, Value: "2"},
					&HeapItem{Fd: 3, Value: "3"},
					&HeapItem{Fd: 4, Value: "4"},
					&HeapItem{Fd: 5, Value: "5"},
					&HeapItem{Fd: 6, Value: "6"},
				},
			},
			{
				infinityPoller,
				[]*HeapItem{
					&HeapItem{Fd: 1, Value: "1", ready: true},
					&HeapItem{Fd: 2, Value: "2", ready: true},
					&HeapItem{Fd: 3, Value: "3", ready: true},
					&HeapItem{Fd: 4, Value: "4", ready: true},
					&HeapItem{Fd: 5, Value: "5", ready: true},
					&HeapItem{Fd: 6, Value: "6", ready: true},
				},
			},
		}

		for i, tt := range tableTests {
			// wg := new(sync.WaitGroup)

			t.Run(fmt.Sprintf("Parallel/%d", i), func(t *testing.T) {
				heapInterface, _ := HeapWithPoller(tt.poller)
				hp, _ := heapInterface.(*pollerHeap)
				hp.pending = tt.pending

				for i := 0; i < len(tt.pending); i++ {
					j := i
					t.Run(fmt.Sprintf("%d", j), func(t *testing.T) {
						t.Parallel()

						// wg.Add(1)
						// defer wg.Done()

						if v := heap.Pop(hp); v == nil {
							t.Fatalf("hp.Pop() -> Unexpected <nil>")
						}
					})
				}
			})

			// t.Run(fmt.Sprintf("PollerInvokes/%d", i), func(t *testing.T) {
			// 	wg.Wait()
			// 	t.Fatal("DMKLDNDJNDK", tt.poller.(*mock).invokes)
			// })
		}
	})

	t.Run("Push", func(t *testing.T) {
		// clear
		hp.pending = []*HeapItem{}

		t.Run("real", func(t *testing.T) {
			heap.Push(hp, &HeapItem{Fd: uintptr(1), Value: "1"})
			pollerHeapState(t, hp, 1)
		})

		// clear
		hp.pending = []*HeapItem{}

		t.Run("fake", func(t *testing.T) {
			heap.Push(hp, "foobar")
			pollerHeapState(t, hp, 0)
		})

		t.Run("poller/error", func(t *testing.T) {
			heapInterface, _ := HeapWithPoller(MockPoller(nil, nil, true))
			hp, _ := heapInterface.(*pollerHeap)
			// real value, but error from poller -> 0
			heap.Push(hp, &HeapItem{Fd: uintptr(1), Value: "foobar"})
			pollerHeapState(t, hp, 0)
		})
	})
}

func TestHeapPrivate(t *testing.T) {
	heapInterface, _ := Heap()
	hp, _ := heapInterface.(*pollerHeap)

	t.Run("len", func(t *testing.T) {
		tableTests := []struct {
			pending []*HeapItem // pending slice for heap (initial)
			len     int
		}{
			{[]*HeapItem{}, 0},
			{[]*HeapItem{&HeapItem{Fd: 1}, &HeapItem{Fd: 2}}, 2},
			{[]*HeapItem{&HeapItem{Fd: 1}, &HeapItem{Fd: 2}, &HeapItem{Fd: 3}}, 3},
		}

		for i, tt := range tableTests {
			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				hp.pending = tt.pending
				pollerHeapState(t, hp, tt.len)
			})
		}
	})

	// t.Run("less", func(t *testing.T) {})

	// t.Run("swap", func(t *testing.T) {})

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
				heapInterface, _ := HeapWithPoller(tt.poller)
				hp, _ := heapInterface.(*pollerHeap)

				re, ce := hp.poll()
				if invokes := tt.poller.(*mock).invokes; invokes != tt.invokes {
					t.Fatalf("poller.invokes -> Expected %v, got %v", tt.invokes, invokes)
				}
				if !reflect.DeepEqual(re, tt.re) {
					t.Fatalf("poller.re -> Expected %v, got %v", tt.re, re)
				}
				if !reflect.DeepEqual(ce, tt.ce) {
					t.Fatalf("poller.ce -> Expected %v, got %v", tt.ce, ce)
				}
			})
		}
	})

	t.Run("pop", func(t *testing.T) {
		tableTests := []struct {
			before []*HeapItem // pending before .pop()
			x      interface{} // value returned by pop()
			after  []*HeapItem // pending after .pop()
		}{
			{[]*HeapItem{}, nil, []*HeapItem{}},
			{[]*HeapItem{&HeapItem{Fd: 1, Value: "1"}}, nil, []*HeapItem{&HeapItem{Fd: 1, Value: "1"}}},
			{[]*HeapItem{&HeapItem{Fd: 1, Value: "1"}, &HeapItem{Fd: 2, Value: "2", ready: true}}, "2", []*HeapItem{&HeapItem{Fd: 1, Value: "1"}}},
			{[]*HeapItem{&HeapItem{Fd: 1, Value: "1", ready: true}, &HeapItem{Fd: 2, Value: "2", ready: true}}, "1", []*HeapItem{&HeapItem{Fd: 2, Value: "2", ready: true}}},
		}

		for i, tt := range tableTests {
			t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
				// fill heap
				hp.pending = tt.before

				// pop and check returned value and counters
				if x := hp.pop(); x != tt.x {
					t.Fatalf("hp.pop() -> Expected %q, got %q", tt.x, x)
				}

				if !reflect.DeepEqual(hp.pending, tt.after) {
					t.Fatalf("hp.pending -> Expected %v, got %v", tt.after, hp.pending)
				}
			})
		}
	})
}
