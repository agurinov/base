package transport

import (
	"container/heap"
	"net"
	"time"

	"github.com/boomfunc/base/server/flow"
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
	inputCh chan *flow.Data
	errCh   chan error

	// poller integration
	heap heap.Interface
}

func (tr *tcp) Connect(inputCh chan *flow.Data, errCh chan error) {
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
			if flow, ok := heap.Pop(tr.heap).(*flow.Data); ok {
				flow.Timing.Exit("poller")
				tr.inputCh <- flow
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
		flow := flow.New(conn)
		item := &poller.HeapItem{Fd: fd, Value: flow}
		flow.Timing.Enter("poller")
		heap.Push(tr.heap, item)
	}
}
