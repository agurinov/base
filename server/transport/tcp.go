package transport

import (
	"io"
	"net"
)

type tcp struct {
	listener *net.TCPListener
	inputCh  chan io.ReadWriteCloser
	errCh    chan error
}

func (tcp *tcp) Connect(inputCh chan io.ReadWriteCloser, errCh chan error) {
	tcp.inputCh = inputCh
	tcp.errCh = errCh
}

func (tcp *tcp) Serve() {
	for {
		conn, err := tcp.listener.AcceptTCP()
		if err != nil {
			// handle error
			tcp.errCh <- err
			continue
		}

		// handle successful connection
		tcp.inputCh <- conn
	}
}
