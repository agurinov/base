package transport

import (
	"container/heap"
	"io"
	"net"
	"time"

	"github.com/boomfunc/base/tools/poller"
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
	heap heap.Interface
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
			if rwc, ok := heap.Pop(tr.heap).(io.ReadWriteCloser); ok {
				tr.inputCh <- rwc
			}
		}
	}()

	for {
		conn, err := tr.listener.AcceptTCP()
		if err != nil {
			tr.errCh <- err
			continue
		}

		fd, err := tcpFD(conn)
		if err != nil {
			tr.errCh <- err
			continue
		}

		// push incoming connection to heap
		item := &poller.HeapItem{Fd: fd, Value: conn}
		heap.Push(tr.heap, item)
	}
}
