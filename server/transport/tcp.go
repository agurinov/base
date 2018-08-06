package transport

import (
	"io"
	"net"
	"time"
)

var (
	// TODO parametrize
	readTimeout  = time.Second * 2
	writeTimeout = time.Second * 5
)

type tcp struct {
	listener *net.TCPListener
	inputCh  chan io.ReadWriteCloser
	errCh    chan error
}

func (tr *tcp) Connect(inputCh chan io.ReadWriteCloser, errCh chan error) {
	tr.inputCh = inputCh
	tr.errCh = errCh
}

// https://habr.com/company/mailru/blog/331784/
// before 3.3.1
func (tr *tcp) Serve() {
	for {
		conn, err := tr.listener.AcceptTCP()
		if err != nil {
			// handle error
			tr.errCh <- err
			continue
		}

		// netpoll block (TODO)
		raw, err := conn.SyscallConn()
		if err != nil {
			// handle error
			tr.errCh <- err
			continue
		}

		// wait for data in socket
		// TODO THIS BLOCK MUST BE REFACTORED
		go func(conn *net.TCPConn) {
			// get signal about incoming data (blocking mode)
			if err := raw.Read(tcpDetectRead); err != nil {
				// handle error
				tr.errCh <- err
			} else {
				// handle successful and ready connection
				// conn.SetReadDeadline(time.Now().Add(readTimeout))
				// conn.SetWriteDeadline(time.Now().Add(writeTimeout))
				tr.inputCh <- conn
			}
		}(conn)
	}
}
