// +build linux

package transport

import (
	"golang.org/x/sys/unix"
)

func NewPoller() (*epoll, error) {
	// create epoll
	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	ep := &epoll{
		epfd: epfd,
		ch:   make(chan []unix.EpollEvent),
	}
	return ep, nil
}

type epoll struct {
	epfd int
	ch   chan []unix.EpollEvent
}

func (p *epoll) Add(fd int32) error {
	event := &unix.EpollEvent{
		Events: unix.EPOLLIN | unix.EPOLLRDHUP | unix.EPOLLET,
		Fd:     fd,
	}

	return unix.EpollCtl(p.epfd, unix.EPOLL_CTL_ADD, int(fd), event)
}

func (p *epoll) Del(fd int32) error {
	return unix.EpollCtl(p.epfd, unix.EPOLL_CTL_DEL, int(fd), nil)
}

func (p *epoll) Wait() ([]unix.EpollEvent, error) {
	events := make([]unix.EpollEvent, 32)

	_, err := unix.EpollWait(p.epfd, events, -1)
	if err != nil {
		return nil, err
	}

	return events, nil

	// for {
	//
	//
	// 	// for i := 0; i < n; i++ {
	// 	// 	ev := events[i]
	// 	// 	// the fd is ready to be read from
	// 	// 	// log.Debugf("EPOLL WAIT. OUR: %d, ENVENTFD: %d", int32(fd), ev.Fd)
	// 	// 	if ev.Fd == int32(fd) {
	// 	// 		// TODO filter
	// 	// 		if ev.Events&(unix.EPOLLIN|unix.EPOLLHUP|unix.EPOLLERR) != 0 {
	// 	// 			return nil
	// 	// 		}
	// 	// 	}
	// 	//
	// 	// 	// the console is ready to be written to
	// 	// 	// if ev.Events&(unix.EPOLLOUT|unix.EPOLLHUP|unix.EPOLLERR) != 0 {
	// 	// 	// 	if epfile := e.getConsole(int(ev.Fd)); epfile != nil {
	// 	// 	// 		epfile.signalWrite()
	// 	// 	// 	}
	// 	// 	// }
	// 	// }
	//
	// }
}
