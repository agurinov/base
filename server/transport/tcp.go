package transport

import (
	"io"
	"net"
	"time"
	// "golang.org/x/sys/unix"
)

var (
	// TODO parametrize
	readTimeout  = time.Second * 2
	writeTimeout = time.Second * 5
)

type tcp struct {
	// server integration
	listener *net.TCPListener
	inputCh  chan io.ReadWriteCloser
	errCh    chan error
	// poller integration
	heap *ConnHeap
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
			// Obtain socket with data from heap/epoll
			// (blocking mode)
			tr.inputCh <- tr.heap.Pop()
			// Also we want be sure there is worker available for this conn (with data)
			// TODO try to fetch worker
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
