package dispatcher

// Worker represents the worker that executes the Task
type Worker struct {
	TaskChannel chan Task
	quit        chan bool
}

func NewWorker() *Worker {
	return &Worker{
		TaskChannel: make(chan Task),
		quit:        make(chan bool),
	}
}

func (w *Worker) attach(d *Dispatcher) {
	// register the current worker into the dispathcer's worker pool
	d.WorkerPool <- w
}

// Start starts worker to wait for incoing task or quit signal
func (w *Worker) Start() {
	for {
		select {
		case task := <-w.TaskChannel:
			// we have received some task to solve!
			task.Solve()

		case <-w.quit:
			// we have received a signal to stop!
			return
		}
	}
}

// Stop signals the worker to stop listening for work requests.
func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
