package poller

import (
	"container/heap"
	"sync"
)

type HeapItem struct {
	Fd    uintptr
	Value interface{}
}

type pollerHeap struct {
	pending map[uintptr]interface{}
	ready   []uintptr
	poller  Interface
	mux     *sync.RWMutex
}

func Heap() (heap.Interface, error) {
	poller, err := New()
	if err != nil {
		return nil, err
	}

	h := &pollerHeap{
		pending: make(map[uintptr]interface{}, 0),
		ready:   make([]uintptr, 0),
		poller:  poller,
		mux:     new(sync.RWMutex),
	}

	heap.Init(h)
	return h, nil
}

func (h pollerHeap) Len() int {
	h.mux.RLock()
	defer h.mux.RUnlock()

	return len(h.ready)
}

func (h pollerHeap) Less(i, j int) bool {
	h.mux.RLock()
	defer h.mux.RUnlock()

	// Less reports whether the element with
	// index i should sort before the element with index j.
	return h.ready[i] < h.ready[j]
}

func (h pollerHeap) Swap(i, j int) {
	if h.Len() >= 2 {
		// there is something to swap
		h.mux.Lock()
		h.ready[i], h.ready[j] = h.ready[j], h.ready[i]
		h.mux.Unlock()
	}
}

func (h *pollerHeap) Push(x interface{}) {
	if item, ok := x.(*HeapItem); ok {
		// try to add to poller
		// TODO error not visible! in transport layer
		if err := h.poller.Add(item.Fd); err == nil {
			// fd in poller, store it for .Pop()
			h.mux.Lock()
			h.pending[item.Fd] = item.Value
			h.mux.Unlock()
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

func (h *pollerHeap) pushReady(x uintptr) {
	h.mux.Lock()
	h.ready = append(h.ready, x)
	h.mux.Unlock()
}

func (h *pollerHeap) popReady() interface{} {
	h.mux.RLock()
	n := len(h.ready)
	h.mux.RUnlock()

	if n == 0 {
		return nil
	}

	h.mux.RLock()
	fd := h.ready[n-1]
	h.mux.RUnlock()

	h.mux.Lock()
	h.ready = h.ready[0 : n-1]
	h.mux.Unlock()

	// get Value by fd
	if value, ok := h.pending[fd]; ok {
		h.mux.Lock()
		delete(h.pending, fd)
		h.mux.Unlock()

		return value
	}

	return nil
}
