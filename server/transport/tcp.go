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
	heap  heap.Interface // also connect with poller
	errCh chan error
}

func (tr *tcp) Connect(heap heap.Interface, errCh chan error) {
	tr.heap = heap
	tr.errCh = errCh
}

func (tr *tcp) Serve() {
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
		flow.Chronometer.Enter("transport")
		heap.Push(tr.heap, item)
	}
}
