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
POP:
	if h.Len() == 0 {
		// case when nothing to fetch -> ask poller
		// fill ready slice from poller signal
		h.poll()
	}
	// pop ready and return
	if value := h.pop(); value != nil {
		return value
	} else {
		// repeat this
		goto POP
	}
}

func (h *pollerHeap) poll() {
	for {
		// blocking mode
		// fetching events from poller
		re, _, ce, err := h.poller.Events()
		if err != nil {
			continue
		}

		// push ready, excluding closed
		readyFds := EventsToFds(re...)
		closeFds := EventsToFds(ce...)
		h.actualize(readyFds, closeFds)

		return
	}
}

func (h *pollerHeap) actualize(ready []uintptr, close []uintptr) {
	h.mux.Lock()
	defer h.mux.Unlock()

	// Phase 1. delete from heap if fd closed
	for i, fd := range h.ready {
		for _, cfd := range close {
			if fd == cfd {
				// fd from heap is closed
				// not relevant, delete it from future .Pop()
				h.ready = append(h.ready[0:i], h.ready[i+1:]...)
			}
		}
	}

	// Phase 2. add to heap if
	// fd has some data (in ready) and not closed (not in close)
	// and fd has associated in pending
OUTER:
	for _, rfd := range ready {
		// check ready event has associated in pending
		// if not - not relevant (we have not gor any returnable value)
		if _, ok := h.pending[rfd]; !ok {
			continue OUTER
		}

		// ready event can be returned, check next for closing at same time
		for _, cfd := range close {
			if rfd == cfd {
				// ready event is closed -> not relevant, skip
				continue OUTER
			}
		}
		h.ready = append(h.ready, rfd)
	}

	// Phase 3. actualize heap pending elements
	// closed elements remove from pending (memory release)
	for _, cfd := range close {
		if _, ok := h.pending[cfd]; ok {
			delete(h.pending, cfd)
		}
	}
}

func (h *pollerHeap) pop() interface{} {
	if h.Len() == 0 {
		return nil
	}

	// there is something to pop
	// at first sight
	h.mux.Lock()
	defer h.mux.Unlock()

	var i int
	var value interface{}

	// get first fd from heap, available in pending
	for j, fd := range h.ready {
		// save counter
		i = j
		// get Value by fd from pending
		var ok bool
		value, ok = h.pending[fd]
		if ok {
			// there is associated value in pending
			// delete from pending and return
			delete(h.pending, fd)
			break
		}
	}

	// we got first associated entry
	// all before not relevant, because
	// if associated pending exists - return this (it is this case)
	// otherwise it is orphans -> try again (loop continue)
	h.ready = h.ready[i+1:]

	return value
}
