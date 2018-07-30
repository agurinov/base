package transport

import (
	"fmt"
	"net"
	// "syscall"

	"github.com/boomfunc/log"
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

	// epfd, e := syscall.EpollCreate1(0)
	// if e != nil {
	// 	fmt.Println("epoll_create1: ", e)
	// 	os.Exit(1)
	// }
	// defer syscall.Close(epfd)

	// syscall.EPOLLIN

	// log.Debug("CUSTOM FUNC", fd&syscall.EPOLLIN)
	log.Debug("CUSTOM FUNC", uint16(fd))
	log.Debug("CUSTOM FUNC", fd&0x1)
	log.Debugf("CUSTOM FUNC: %+v", fd)
	return false
}
