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

	tcp := &tcp{
		listener: tcpListener,
	}
	return tcp, nil
}

func tcpFD(conn *net.TCPConn) (uintptr, error) {
	raw, err := conn.SyscallConn()
	if err != nil {
		return 0, err
	}

	var fd uintptr

	// TODO for now it is some kind of workaround
	// TODO inner error not visible!
	f := func(innerFd uintptr) bool {
		fd = innerFd
		return true
	}

	if err := raw.Read(f); err != nil {
		return 0, err
	}

	return fd, nil
}
