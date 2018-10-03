package poller

import (
	"container/heap"
	// "sort"
	"sync"
	// "github.com/boomfunc/log"
)

type HeapItem struct {
	Fd    uintptr
	Value interface{}
	ready bool
}

type pollerHeap struct {
	// poller integration
	poller       Interface
	pollerLocked bool
	pl           *sync.Cond // poller locking

	mux sync.RWMutex // this mutex guards variables below

	pending []*HeapItem
}

func Heap() (heap.Interface, error) {
	poller, err := New()
	if err != nil {
		return nil, err
	}

	h := &pollerHeap{
		poller:  poller,
		pl:      sync.NewCond(&sync.Mutex{}),
		pending: make([]*HeapItem, 0),
	}

	heap.Init(h)
	return h, nil
}

func (h *pollerHeap) len() int {
	return len(h.pending)
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
	if !h.pending[i].ready && h.pending[j].ready {
		return true
	}
	return false
}

func (h *pollerHeap) Less(i, j int) bool {
	return false
	// h.mux.RLock()
	// b := h.less(i, j)
	// h.mux.RUnlock()
	//
	// return b
}

func (h *pollerHeap) swap(i, j int) {
	if h.len() >= 2 {
		// there is something to swap
		h.pending[i], h.pending[j] = h.pending[j], h.pending[i]
	}
}

func (h *pollerHeap) Swap(i, j int) {
	return
	// h.mux.Lock()
	// h.swap(i, j)
	// h.mux.Unlock()
}

func (h *pollerHeap) Push(x interface{}) {
	if item, ok := x.(*HeapItem); ok {
		// try to add to poller
		// TODO error not visible! in transport layer
		if err := h.poller.Add(item.Fd); err == nil {
			// fd in poller, store it for .Pop()
			h.mux.Lock()
			h.pending = append(h.pending, item)
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
				// blocking mode operation !!
				re, ce := h.poll()

				// events are received (and they are!)
				// push ready, excluding closed
				h.mux.Lock()
				h.actualize(re, ce)
				h.mux.Unlock()
				// sort.Sort(h)

				// all events polled, data refreshed
				// unlock all waiting goroutines (server with .Pop() method invoked)
				h.pl.Broadcast()

				// unlock parent IF statement for another goroutines
				h.mux.Lock()
				h.pollerLocked = false
				h.mux.Unlock()
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

// poll is the main link between heap and core poller
// blocking operation until really received some events
// otherwise poll again
func (h *pollerHeap) poll() ([]uintptr, []uintptr) {
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
		return EventsToFds(re...), EventsToFds(ce...)
	}
}

// actualize called after success polling process finished
// purpose: update state (add new ready, delete closed)
func (h *pollerHeap) actualize(ready []uintptr, close []uintptr) {
	var new []*HeapItem

OUTER:
	for _, item := range h.pending {
		// check in closed
		for _, cfd := range close {
			if item.Fd == cfd {
				// fd from heap is closed
				// not relevant, delete it from future .Pop()
				// delete = skip from appending
				continue OUTER
			}
		}

		// this item not closed yet, check for ready state
		for _, rfd := range ready {
			if item.Fd == rfd {
				item.ready = true
			}
		}

		new = append(new, item)
	}

	// update heap items
	h.pending = new
}

// pop searches first entry in `pending` map
// using key from `ready` slice according sort.interface
func (h *pollerHeap) pop() interface{} {
	if h.len() == 0 {
		return nil
	}

	// there is something to pop (at first sight)
	var i int

	defer func() {
		// we return value but we need to cut this value from pending
		// if nothing returns -> i == len(h.pending - 1) and we return same slice
		h.pending = append(h.pending[:i], h.pending[i+1:]...)
	}()

	// get first fd from heap, available in pending
	for j, item := range h.pending {
		if item.ready {
			i = j
			// log.Debugf("RETURNED: i=%d, Fd=%d", i, item.Fd)
			return item.Value
		}
	}

	// nobody ready
	// log.Debug("RETURNED <nil>")
	return nil
}
