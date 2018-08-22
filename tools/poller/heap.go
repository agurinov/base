package poller

import (
	"container/heap"
)

type HeapItem struct {
	Fd    uintptr
	Value interface{}
}

type pollerHeap struct {
	pending map[uintptr]interface{}
	ready   []uintptr
	poller  Interface
}

func Heap() (*pollerHeap, error) {
	poller, err := New()
	if err != nil {
		return nil, err
	}

	h := &pollerHeap{
		pending: make(map[uintptr]interface{}, 0),
		ready:   make([]uintptr, 0),
		poller:  poller,
	}

	heap.Init(h)
	return h, nil
}

func (h pollerHeap) Len() int {
	return len(h.ready)
}

func (h pollerHeap) Less(i, j int) bool {
	// Less reports whether the element with
	// index i should sort before the element with index j.
	return h.ready[i] < h.ready[j]
}

func (h pollerHeap) Swap(i, j int) {
	if i < 0 || j < 0 {
		return
	}
	h.ready[i], h.ready[j] = h.ready[j], h.ready[i]
}

func (h *pollerHeap) Push(x interface{}) {
	if item, ok := x.(*HeapItem); ok {
		// try to add to poller
		// TODO error not visible! in transport layer
		if err := h.poller.Add(item.Fd); err == nil {
			// fd in poller, store it for .Pop()
			h.pending[item.Fd] = item.Value
		}
	}
}

func (h *pollerHeap) Pop() interface{} {
	// fill ready slice from poller signal
	h.fillFromPoller()
	// pop ready and return
	return h.popReady()
}

func (h *pollerHeap) fillFromPoller() {
	for {
		// blocking mode
		// fetching ready events from poller
		re, _, err := h.poller.Events()
		if err != nil {
			continue
		}

		// iterate over read ready events
		for _, event := range re {
			key := event.Fd()
			// mark as ready to pop
			h.pushReady(key)
			// rm from epoll
			h.poller.Del(key)
		}

		break
	}
}

func (h *pollerHeap) pushReady(x interface{}) {
	h.ready = append(h.ready, x.(uintptr))
}

func (h *pollerHeap) popReady() interface{} {
	n := len(h.ready)

	if n == 0 {
		return nil
	}

	fd := h.ready[n-1]
	h.ready = h.ready[0 : n-1]

	// get Value by fd
	if value, ok := h.pending[fd]; ok {
		delete(h.pending, fd)
		return value
	}

	return nil
}
