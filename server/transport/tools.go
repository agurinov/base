package transport

import (
	"fmt"
	"net"
)

func TCP(ip net.IP, port int) (Interface, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, err
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}

	poller, err := NewPoller()
	if err != nil {
		return nil, err
	}
	heap := &ConnHeap{
		items:  make(map[uintptr]*net.TCPConn),
		poller: poller,
	}

	tcp := &tcp{
		listener: tcpListener,
		heap:     heap,
	}
	return tcp, nil
}
