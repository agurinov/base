package server

import (
	"container/heap"

	"github.com/boomfunc/log"
)

type Balancer struct {
	pool Pool
	done chan *Worker
}

func NewBalancer() *Balancer {
	nWorker := 2

	done := make(chan *Worker, nWorker)
	b := &Balancer{
		pool: make(Pool, 0, nWorker),
		done: done
	}

	for i := 0; i < nWorker; i++ {
		w := &Worker{requests: make(chan Request)}
		heap.Push(&b.pool, w)
		go w.work(b.done)
	}

	return b
}

func (b *Balancer) dispatch(req Request) {
	log.Debug("b.dispatch()")
	// Grab the least loaded worker...
	w := heap.Pop(&b.pool).(*Worker)
	// ...send it the task.
	w.requests <- req
	// One more in its work queue.
	w.pending++
	// Put it into its place on the heap.
	heap.Push(&b.pool, w)
}

// Job is complete; update heap
func (b *Balancer) completed(w *Worker) {
	// One fewer in the queue.
	w.pending--
	// Remove it from heap.
	heap.Remove(&b.pool, w.index)
	// Put it into its place on the heap.
	heap.Push(&b.pool, w)
}

func (b *Balancer) balance(work chan Request) {
	for {
		select {
		case req := <-work: // received a Request...
			b.dispatch(req) // ...so send it to a Worker
		case w := <-b.done: // a worker has finished ...
			b.completed(w)  // ...so update its info
		}
	}
}
