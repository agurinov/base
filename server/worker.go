package server

type Worker struct {
	requests chan Request // work to do (buffered channel)
	pending  int          // count of pending tasks
	index     int         // index in the heap
}

func NewWorker(requests chan Request) *Worker {
	return &Worker{requests: requests}
}

func (w *Worker) work(done chan *Worker) {
	for {
		req := <-w.requests // get Request from balancer
		req.Handle()   // call fn and send result
		done <- w           // we've finished this request
	}
}

type Pool []*Worker

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p Pool) Len() int {
	return len(p)
}

func (p *Pool) Swap(i, j int) {
	a := *p
	a[i], a[j] = a[j], a[i]
	a[i].index = i
	a[j].index = j
}

func (p *Pool) Push(x interface{}) {
	a := *p
	n := len(a)
	a = a[0 : n+1]
	w := x.(*Worker)
	a[n] = w
	w.index = n
	*p = a
}

func (p *Pool) Pop() interface{} {
	a := *p
	*p = a[0 : len(a)-1]
	w := a[len(a)-1]
	w.index = -1 // for safety
	return w
}
