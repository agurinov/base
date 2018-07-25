package transport

import (
	"io"
	// "net"
)

type udp struct {
	inputCh chan io.ReadWriteCloser
	errCh   chan error
}

func (tr *udp) Connect(inputCh chan io.ReadWriteCloser, errCh chan error) {
	tr.inputCh = inputCh
	tr.errCh = errCh
}

func (tr *udp) Serve() {
	for {
		// conn, err := net.ListenUDP(network string, laddr *UDPAddr) (*UDPConn, error)
		// if err != nil {
		// 	// handle error
		// 	tr.errCh <- err
		// 	continue
		// }
		//
		// // handle successful connection
		// tr.inputCh <- conn
	}
}
