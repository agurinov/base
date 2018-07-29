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

func (tr *tcp) Connect(inputCh chan io.ReadWriteCloser, errCh chan error) {
	tr.inputCh = inputCh
	tr.errCh = errCh
}

func (tr *tcp) Serve() {
	for {
		conn, err := tr.listener.AcceptTCP()
		if err != nil {
			// handle error
			tr.errCh <- err
			continue
		}

		// handle successful connection
		// TODO maybe send connections obly when the caller starts to write to it?
		tr.inputCh <- conn
	}
}
