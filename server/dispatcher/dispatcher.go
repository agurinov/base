package dispatcher

type Task interface {
	Solve()
}

var TaskChannel = make(chan Task)

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Task
	maxWorkers int
}

func New(maxWorkers int) *Dispatcher {
	d := &Dispatcher{
		WorkerPool: make(chan chan Task, maxWorkers),
		maxWorkers: maxWorkers,
	}
	d.Prepare()

	return d
}

func (d *Dispatcher) Prepare() {
	// TODO clear current stack of TaskChannel
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		NewWorker(d.WorkerPool).Start()
	}

	StartupLog(d.maxWorkers)
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
