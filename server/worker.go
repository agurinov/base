package server

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
			case request := <-w.RequestChannel:
				request.server.app.HandleRequest(request.under)

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
