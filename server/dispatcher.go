package server

import (
	"github.com/boomfunc/base/server/request"
)

var RequestQueue = make(chan request.Request)

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan request.Request
	maxWorkers int
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	return &Dispatcher{
		WorkerPool: make(chan chan request.Request, maxWorkers),
		maxWorkers: maxWorkers,
	}
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		NewWorker(d.WorkerPool).Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case req := <-RequestQueue:
			// a job request has been received
			go func(req request.Request) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				workerChannel := <-d.WorkerPool
				// dispatch the job to the worker job channel
				workerChannel <- req
			}(req)
		}
	}
}
