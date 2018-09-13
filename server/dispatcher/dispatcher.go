package dispatcher

type Task interface {
	Solve()
}

var TaskChannel = make(chan Task)

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Task
	MaxWorkers int
}

func New(MaxWorkers int) *Dispatcher {
	return &Dispatcher{
		WorkerPool: make(chan chan Task, MaxWorkers),
		MaxWorkers: MaxWorkers,
	}
}

func (d *Dispatcher) Prepare() {
	// TODO clear current stack of TaskChannel
	// starting n number of workers
	for i := 0; i < d.MaxWorkers; i++ {
		NewWorker(d.WorkerPool).Start()
	}

	StartupLog(d.MaxWorkers)
}

// FreeTaskChannel returns channel of free worker
// this will block until a worker is idle
func (d *Dispatcher) FreeTaskChannel() chan Task {
	return <-d.WorkerPool
}

// Dispatch is infinity loop for listening global TaskChannel
// NOTE: possible leaks when allocating some data before sending task to channel
// because worker may hang
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
