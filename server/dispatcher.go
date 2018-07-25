package server

// TODO look carefully
type Task func()

var TaskChannel = make(chan Task)

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Task
	maxWorkers int
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	d := &Dispatcher{
		WorkerPool: make(chan chan Task, maxWorkers),
		maxWorkers: maxWorkers,
	}
	d.Prepare()

	return d
}

func (d *Dispatcher) Prepare() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		NewWorker(d.WorkerPool).Start()
	}
}

func (d *Dispatcher) Dispatch() {
	for {
		select {
		case task := <-TaskChannel:
			// a Task request has been received
			go func(task Task) {
				// try to obtain a worker Task channel that is available.
				// this will block until a worker is idle
				workerTaskChannel := <-d.WorkerPool
				// dispatch the Task to the worker Task channel
				workerTaskChannel <- task
			}(task)
		}
	}
}

// Worker represents the worker that executes the Task
type Worker struct {
	WorkerPool  chan chan Task
	TaskChannel chan Task
	quit        chan bool
}

func NewWorker(workerPool chan chan Task) *Worker {
	return &Worker{
		WorkerPool:  workerPool,
		TaskChannel: make(chan Task),
		quit:        make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.TaskChannel

			select {
			case task := <-w.TaskChannel:
				task()

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
