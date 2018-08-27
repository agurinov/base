package poller

import (
	"testing"
)

func TestNewPollerHeap(t *testing.T) {
	heap, err := Heap()
	if err != nil {
		t.Fatal(err)
	}
	if heap == nil {
		t.Fatal("Unexpected heap (expected *pollerHeap)")
	}
}

func pollerHeapLen(t *testing.T, heap *pollerHeap, pending, ready int) {
	t.Run("ready/Len", func(t *testing.T) {
		// ready via .Len()
		if n := heap.Len(); n != ready {
			t.Fatalf("Expected \"%d\", got \"%d\"", ready, n)
		}
	})

	t.Run("ready/direct", func(t *testing.T) {
		// ready direct
		if n := len(heap.ready); n != ready {
			t.Fatalf("Expected \"%d\", got \"%d\"", ready, n)
		}
	})

	t.Run("pending/direct", func(t *testing.T) {
		// pending direct
		if n := len(heap.pending); n != pending {
			t.Fatalf("Expected \"%d\", got \"%d\"", pending, n)
		}
	})
}

func TestPollerHeapLen(t *testing.T) {
	h, _ := Heap()
	heap := h.(*pollerHeap)
	pollerHeapLen(t, heap, 0, 0)

	t.Run("Push", func(t *testing.T) {
		heap.Push(&HeapItem{Fd: uintptr(1), Value: "foobar"})
		pollerHeapLen(t, heap, 1, 0)
	})

	t.Run("pushReady", func(t *testing.T) {
		heap.pushReady(uintptr(1))
		pollerHeapLen(t, heap, 1, 1)
	})
}

func TestPollerHeapPush(t *testing.T) {
	h, _ := Heap()
	heap := h.(*pollerHeap)

	t.Run("item", func(t *testing.T) {
		heap.Push(&HeapItem{Fd: uintptr(1), Value: "foobar"})
		pollerHeapLen(t, heap, 1, 0)
	})

	t.Run("fake", func(t *testing.T) {
		heap.Push("foobar")
		pollerHeapLen(t, heap, 1, 0)
	})
}

func TestPollerHeapPopReady(t *testing.T) {
	h, _ := Heap()
	heap := h.(*pollerHeap)

	t.Run("empty", func(t *testing.T) {
		x := heap.popReady()
		if x != nil {
			t.Fatalf("Expected \"%v\", got %q", nil, x)
		}
	})

	t.Run("exists", func(t *testing.T) {
		heap.Push(&HeapItem{Fd: uintptr(1), Value: "foobar"})
		heap.pushReady(uintptr(1))

		x := heap.popReady()
		if x != "foobar" {
			t.Fatalf("Expected %q, got %q", "foobar", x)
		}
	})
}
