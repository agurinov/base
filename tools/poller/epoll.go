// +build linux

package poller

import (
	"golang.org/x/sys/unix"
)

type epoll struct {
	epfd int
}

func New() (Interface, error) {
	// create epoll
	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	ep := &epoll{
		epfd: epfd,
	}
	return ep, nil
}

func (p *epoll) Add(fd int32) error {
	event := &unix.EpollEvent{
		Events: unix.EPOLLIN | unix.EPOLLOUT | unix.EPOLLRDHUP | unix.EPOLLET,
		Fd:     fd,
	}

	return unix.EpollCtl(p.epfd, unix.EPOLL_CTL_ADD, int(fd), event)
}

func (p *epoll) Del(fd int32) error {
	return unix.EpollCtl(p.epfd, unix.EPOLL_CTL_DEL, int(fd), nil)
}

func (p *epoll) Events() ([]unix.EpollEvent, []unix.EpollEvent, error) {
	events, err := p.wait()
	if err != nil {
		return nil, nil, err
	}

	// something received, try it
	var re, we []unix.EpollEvent
	for _, event := range events {
		if event.Events&(unix.EPOLLRDHUP) != 0 {
			// closed by peer
			// nothing to do:
			// http://man7.org/linux/man-pages/man7/epoll.7.html
		} else if event.Events&(unix.EPOLLIN) != 0 {
			// ready to read
			re = append(re, event)
		} else if event.Events&(unix.EPOLLOUT) != 0 {
			// ready to write
			we = append(we, event)
		}
	}

	return re, we, nil

}

func (p *epoll) wait() ([]unix.EpollEvent, error) {
	// TODO max events???
	events := make([]unix.EpollEvent, 32)

	_, err := unix.EpollWait(p.epfd, events, -1)
	if err != nil {
		return nil, err
	}

	return events, nil
}
