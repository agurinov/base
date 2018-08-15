package transport

import (
	"fmt"
	"net"

	"github.com/boomfunc/base/tools/poller/heap"
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

	heap, err := heap.NewPollerHeap()
	if err != nil {
		return nil, err
	}

	tcp := &tcp{
		listener: tcpListener,
		heap:     heap,
	}
	return tcp, nil
}
