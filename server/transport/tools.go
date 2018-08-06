package transport

import (
	"fmt"
	"net"

	"github.com/boomfunc/base/tools"
	"golang.org/x/sys/unix"
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

	tcp := &tcp{listener: tcpListener}
	return tcp, nil
}

func tcpDetectRead(fd uintptr) (done bool) {
	var event unix.EpollEvent
	var events [32]unix.EpollEvent

	// create epoll
	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		tools.FatalLog(err)
	}

	// add conn to epoll
	// for listen incoming data
	event.Events = unix.EPOLLIN | unix.EPOLLET
	event.Fd = int32(fd)
	if err := unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, int(fd), &event); err != nil {
		tools.FatalLog(err)
	}

	_, err = unix.EpollWait(epfd, events[:], -1)
	return err == nil

	// wait epoll
	// for {
	// 	_, err := unix.EpollWait(epfd, events[:], -1)
	// 	if err != nil {
	// 		tools.FatalLog(err)
	// 		break
	// 	}
	//
	// 	return true
	// }

	// return false
}
