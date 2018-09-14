package dispatcher

// Task is abstract job
// something that the worker can do
type Task interface {
	Solve()
}

type Dispatcher struct {
	// TaskChannel is global channel for all incoming tasks
	// used in Dispatch() mode (forever loop)
	TaskChannel chan Task
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Task
	MaxWorkers int
}

// New returns new Dispatcher instance with all channels linked
func New(MaxWorkers int) *Dispatcher {
	return &Dispatcher{
		TaskChannel: make(chan Task),
		WorkerPool:  make(chan chan Task, MaxWorkers),
		MaxWorkers:  MaxWorkers,
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

// FreeWorkerTaskChannel returns channel of free worker
// this will block until a worker is idle
func (d *Dispatcher) FreeWorkerTaskChannel() chan Task {
	return <-d.WorkerPool
}

// Dispatch is infinity loop for listening global TaskChannel
// NOTE: possible leaks when allocating some data before sending task to channel
// because worker may hang
func (d *Dispatcher) Dispatch() {
	for {
		select {
		case task := <-d.TaskChannel:
			// a Task request has been received
			go func(task Task) {
				// try to obtain a worker Task channel that is available.
				// and then send the Task to the worker Task channel
				d.FreeWorkerTaskChannel() <- task
			}(task)
		}
	}
}
