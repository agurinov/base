package transport

import (
	"io"
	"net"
	"time"

	"github.com/boomfunc/base/tools/poller/heap"
)

var (
	// TODO parametrize
	readTimeout  = time.Second * 2
	writeTimeout = time.Second * 5
)

type tcp struct {
	listener *net.TCPListener

	// server integration
	inputCh chan io.ReadWriteCloser
	errCh   chan error

	// poller integration
	// TODO implement golang heap
	heap *heap.PollerHeap
}

func (tr *tcp) Connect(inputCh chan io.ReadWriteCloser, errCh chan error) {
	tr.inputCh = inputCh
	tr.errCh = errCh
}

// https://habr.com/company/mailru/blog/331784/
// before 3.3.1
func (tr *tcp) Serve() {
	go func() {
		for {
			// Obtain socket with data from heap/poller
			// (blocking mode)
			tr.inputCh <- tr.heap.Pop()
		}
	}()

	for {
		conn, err := tr.listener.AcceptTCP()
		if err != nil {
			tr.errCh <- err
			continue
		}

		if err := tr.heap.Push(conn); err != nil {
			tr.errCh <- err
		}
	}
}
