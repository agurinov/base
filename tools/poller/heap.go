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
	// poller integration
	poller       Interface
	pollerLocked bool
	pl           *sync.Cond // poller locking

	mux sync.RWMutex // this mutex guards variables below

	pending map[uintptr]interface{}
	ready   []uintptr
}

func Heap() (heap.Interface, error) {
	poller, err := New()
	if err != nil {
		return nil, err
	}

	h := &pollerHeap{
		pl:      sync.NewCond(&sync.Mutex{}),
		pending: make(map[uintptr]interface{}, 0),
		ready:   make([]uintptr, 0),
		poller:  poller,
	}

	heap.Init(h)
	return h, nil
}

func (h *pollerHeap) len() int {
	return len(h.ready)
}

func (h *pollerHeap) Len() int {
	h.mux.RLock()
	n := h.len()
	h.mux.RUnlock()

	return n
}

func (h *pollerHeap) less(i, j int) bool {
	// Less reports whether the element with
	// index i should sort before the element with index j.
	return h.ready[i] < h.ready[j]
}

func (h *pollerHeap) Less(i, j int) bool {
	h.mux.RLock()
	b := h.less(i, j)
	h.mux.RUnlock()

	return b
}

func (h *pollerHeap) swap(i, j int) {
	if h.len() >= 2 {
		// there is something to swap
		h.ready[i], h.ready[j] = h.ready[j], h.ready[i]
	}
}

func (h *pollerHeap) Swap(i, j int) {
	h.mux.Lock()
	h.swap(i, j)
	h.mux.Unlock()
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
	for {
		// is polling process running now?
		// there is only one place per unit of time
		h.mux.RLock()
		locked := h.pollerLocked
		h.mux.RUnlock()

		if !locked {
			// lock this IF statement for another goroutines
			h.mux.Lock()
			h.pollerLocked = true
			h.mux.Unlock()

			// run poll with actualizing behind the scenes
			go func() {
				h.poll()

				// unlock parent IF statement for another goroutines
				h.mux.Lock()
				h.pollerLocked = false
				h.mux.Unlock()

				// all events polled, data refreshed
				// unlock all waiting goroutines (server with .Pop() method invoked)
				h.pl.Broadcast()
			}()
		}

		// wait if necessary
		h.pl.L.Lock()
		if h.Len() == 0 {
			// case when nothing to fetch -> ask poller and wait as long as necessary
			// fill ready slice from poller signal
			h.pl.Wait()
		}
		h.pl.L.Unlock()

		// pop ready and return
		h.mux.Lock()
		value := h.pop()
		h.mux.Unlock()

		if value != nil {
			return value
		} else {
			// repeat this
			continue
		}
	}
}

func (h *pollerHeap) poll() {
	var re, ce []Event

	for {
		var err error
		// fetching events from poller
		// blocking mode !!!
		re, _, ce, err = h.poller.Events()
		if err != nil {
			// some error from poller -> poll again
			continue
		}

		if len(re)+len(ce) == 0 {
			// not required events came -> poll again
			continue
		}

		// all fetched without errors
		break
	}

	// events are received (and they are!)
	// push ready, excluding closed
	h.mux.Lock()
	h.actualize(
		EventsToFds(re...),
		EventsToFds(ce...),
	)
	h.mux.Unlock()
}

func (h *pollerHeap) actualize(ready []uintptr, close []uintptr) {
	// Phase 1. delete from heap if fd closed
	var nready []uintptr
OUTER1:
	for _, fd := range h.ready {
		for _, cfd := range close {
			if fd == cfd {
				// fd from heap is closed
				// not relevant, delete it from future .Pop()
				continue OUTER1
			}
		}
		// if this code reachable -> no continue in close loop
		// no occurences in close - save
		nready = append(nready, fd)
	}
	h.ready = nready

	// Phase 2. add to heap if
	// fd has some data (in ready) and not closed (not in close)
	// and fd has associated in pending
OUTER2:
	for _, rfd := range ready {
		// check ready event has associated in pending
		// if not - not relevant (we have not got any returnable value)
		if _, ok := h.pending[rfd]; !ok {
			continue OUTER2
		}

		// ready event can be returned, check next for closing at same time
		for _, cfd := range close {
			if rfd == cfd {
				// ready event is closed -> not relevant, skip
				continue OUTER2
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
	if h.len() == 0 {
		return nil
	}

	// there is something to pop (at first sight)
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
