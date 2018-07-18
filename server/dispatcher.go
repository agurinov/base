package server

var RequestQueue = make(chan Request)

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
		case request := <-RequestQueue:
			// a job request has been received
			go func(request Request) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				requestChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				requestChannel <- request
			}(request)
		}
	}
}
