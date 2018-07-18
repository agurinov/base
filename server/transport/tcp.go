package transport

import (
	"net"
)

type tcp struct {
	listener *net.TCPListener
	connCh   chan net.Conn
	errCh    chan error
}

func (tcp *tcp) Serve() {
	for {
		conn, err := tcp.listener.AcceptTCP()
		if err != nil {
			// handle error
			tcp.Error(err)
			continue
		}

		// handle successful connection
		tcp.Conn(conn)
	}
}

func (tcp *tcp) Conn(conn net.Conn) {
	tcp.connCh <- conn
}

func (tcp *tcp) Error(err error) {
	tcp.errCh <- err
}
