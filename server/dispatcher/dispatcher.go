package dispatcher

// Task is abstract job
// something that the worker can do
type Task interface {
	Solve()
}

type Dispatcher struct {
	// A pool of workers that are registered with the dispatcher
	WorkerPool chan *Worker
	MaxWorkers int
}

// New returns new Dispatcher instance with all channels linked
func New(MaxWorkers int) *Dispatcher {
	return &Dispatcher{
		WorkerPool: make(chan *Worker, MaxWorkers),
		MaxWorkers: MaxWorkers,
	}
}

func (d *Dispatcher) Prepare() {
	// starting n number of workers
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker()
		d.AttachWorker(worker)
		go worker.Start()
	}

	StartupLog(d.MaxWorkers)
}

// OccupyWorker returns free worker from dispatcher system
// this will block until all worker is idle
func (d *Dispatcher) OccupyWorker() *Worker {
	return <-d.WorkerPool
}

// ReleaseWorker returns busy worker to dispatcher system
// and releases its resources
// may be used for initial state or for `release` worker
func (d *Dispatcher) AttachWorker(w *Worker) {
	w.attach(d)
}
