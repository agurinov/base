package server

import (
	"github.com/boomfunc/base/server/request"
)

var RequestChannel = make(chan Request)

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Request
	maxWorkers int
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	return &Dispatcher{
		WorkerPool: make(chan chan Request, maxWorkers),
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
		case request := <-RequestChannel:
			// a job request has been received
			go func(request Request) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				workerChannel := <-d.WorkerPool
				// dispatch the job to the worker job channel
				workerChannel <- request
			}(request)
		}
	}
}

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool     chan chan Request
	RequestChannel chan Request
	quit           chan bool
}

func NewWorker(workerPool chan chan Request) *Worker {
	return &Worker{
		WorkerPool:     workerPool,
		RequestChannel: make(chan Request),
		quit:           make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.RequestChannel

			select {
			case r := <-w.RequestChannel:
				server := r.server
				req := request.New(r.conn)
				server.responseCh <- server.app.HandleRequest(req, r.conn)

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
