package poller

import (
	"container/heap"
)

type HeapItem struct {
	Fd    uintptr
	index int
	Value interface{}
}

type pollerHeap struct {
	items  []*HeapItem
	poller Interface
}

func Heap() (*pollerHeap, error) {
	poller, err := New()
	if err != nil {
		return nil, err
	}

	h := &pollerHeap{
		items:  make([]*HeapItem, 0),
		poller: poller,
	}

	heap.Init(h)
	return h, nil
}

func (h pollerHeap) Len() int {
	return len(h.items)
}

func (h pollerHeap) Less(i, j int) bool {
	// Less reports whether the element with
	// index i should sort before the element with index j.
	return h.items[i].Fd < h.items[j].Fd
}

func (h pollerHeap) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
	h.items[i].index = i
	h.items[j].index = j
}

func (h *pollerHeap) Push(x interface{}) {
	// h.poller.Add(fd)

	// push to heap (success case)
	n := len(h.items)
	item := x.(*HeapItem)
	item.index = n
	h.items = append(h.items, item)
	// re-establish
	heap.Fix(h, n)
}

func (h *pollerHeap) Pop() interface{} {
	// for {
	// 	// blocking mode
	// 	re, _, err := h.poller.Events()
	// 	if err != nil {
	// 		continue
	// 	}
	//
	// 	// iterate over read ready events
	// 	// for _, event := range re {
	// 	// 	// key := event.Fd()
	// 	//
	// 	//
	// 	//
	// 	// 	// if conn, ok := h.items[key]; ok {
	// 	// 	// 	// rm from epoll
	// 	// 	// 	h.poller.Del(key)
	// 	// 	// 	// rm from heap and return
	// 	// 	// 	delete(ph.items, key)
	// 	// 	// 	return conn
	// 	// 	// }
	// 	// }
	// 	// break infinity loop, ready events are in heap
	// 	break
	// }

	// pop first priority and ready fd (minimal fd)
	old := *h
	n := len(old.items)
	item := old.items[n-1]
	item.index = -1 // for safety
	h.items = old.items[0 : n-1]
	return item.Value
}
