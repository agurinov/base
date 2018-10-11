package poller

import (
	"container/heap"
	"sync"
)

type HeapItem struct {
	Fd    uintptr
	Value interface{}
	ready bool
}

type pollerHeap struct {
	// poller integration
	poller Interface
	once   sync.Once // poller locking

	// mutex and state
	mutex   sync.Mutex // this mutex guards variables below
	cond    *sync.Cond
	pending []*HeapItem
}

func Heap() (heap.Interface, error) {
	poller, err := New()
	if err != nil {
		return nil, err
	}

	h := new(pollerHeap)
	h.poller = poller
	h.cond = sync.NewCond(&h.mutex)
	h.pending = make([]*HeapItem, 0)

	heap.Init(h)
	return h, nil
}

func (h *pollerHeap) len() int {
	return len(h.pending)
}

func (h *pollerHeap) Len() int {
	h.mutex.Lock()
	n := h.len()
	h.mutex.Unlock()

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
	// h.mutex.RLock()
	// b := h.less(i, j)
	// h.mutex.RUnlock()
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
	// h.mutex.Lock()
	// h.swap(i, j)
	// h.mutex.Unlock()
}

// Pop implements heap.Interface
// adds flow to poller
func (h *pollerHeap) Push(x interface{}) {
	if item, ok := x.(*HeapItem); ok {
		// try to add to poller
		// TODO error not visible! in transport layer
		if err := h.poller.Add(item.Fd); err == nil {
			// fd in poller, store it for .Pop()
			h.mutex.Lock()
			h.pending = append(h.pending, item)
			h.mutex.Unlock()
		}
	}
}

// Pop implements heap.Interface
// pops first in flow with `ready` status
func (h *pollerHeap) Pop() interface{} {
	for {
		go h.Poll()

		// try to pop ready
		h.mutex.Lock()

		value := h.pop()
		if value == nil {
			h.cond.Wait()
		}

		h.mutex.Unlock()

		if value != nil {
			return value
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

// Poll is thread safety operation for waiting events from poller
// and actualize heap data
func (h *pollerHeap) Poll() {
	// f is poll with actualizing heap data
	// 1 running instance of this function per time across all workers
	f := func() {
		// blocking mode operation !!
		re, ce := h.poll()
		// events are received (and they are!)
		h.actualize(re, ce) // push ready, excluding closed
		h.cond.Broadcast()  // release all .Wait()
	}

	// f invokes with mutex locking on once.Do layer -> thread safety
	h.once.Do(f)

	// invalidate once for future
	h.mutex.Lock()
	h.once = sync.Once{} // reset sync.Once
	h.mutex.Unlock()
}

// actualize called after success polling process finished
// purpose: update state (add new ready, delete closed)
func (h *pollerHeap) actualize(ready []uintptr, close []uintptr) {
	filtered := pendingFilterClosed(h.pending, close)
	mapped := pendingMapReady(filtered, ready)

	h.pending = mapped
}

// pop searches first entry in `pending` slice
// which has `ready` flag == true
func (h *pollerHeap) pop() interface{} {
	if h.len() == 0 {
		return nil
	}

	// there is something to pop (at first sight)
	// get first fd from heap, available in pending
	for i, item := range h.pending {
		if item.ready {
			h.pending = append(h.pending[:i], h.pending[i+1:]...)
			return item.Value
		}
	}

	// nobody ready
	return nil
}
