package transport

import (
	"fmt"
	"net"
)

func TCP(ip net.IP, port int, connCh chan net.Conn, errCh chan error) (Interface, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, err
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}

	tcp := &tcp{
		listener: tcpListener,
		errCh:    errCh,
		connCh:   connCh,
	}
	return tcp, nil
}
