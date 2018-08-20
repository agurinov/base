package transport

import (
	"container/heap"
	"io"
	"net"
	"time"

	"github.com/boomfunc/base/tools/poller"
	"github.com/boomfunc/log"
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
			// Obtain socket with data from heap/epoll
			// (blocking mode)
			tr.inputCh <- heap.Pop(tr.heap).(io.ReadWriteCloser)
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

		fd, err := tcpFD(conn)
		if err != nil {
			tr.errCh <- err
			continue
		}

		log.Debug("FD: ", fd)

		conn.SetDeadline(time.Now().Add(time.Second * 5))
		heap.Push(tr.heap, &poller.HeapItem{Fd: fd, Value: conn})

	}
}
