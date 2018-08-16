// +build linux

package poller

import (
	"golang.org/x/sys/unix"
)

type epollEvent struct {
	se unix.EpollEvent // source event
}

// Fd implements base Event interface
func (ev epollEvent) Fd() uintptr {
	return uintptr(ev.se.Fd)
}

type epoll struct {
	fd int
}

func New() (Interface, error) {
	// create epoll
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	poller := &epoll{fd}
	return poller, nil
}

func (p *epoll) Add(fd uintptr) error {
	// TODO we tracking creation of fd - no matter for this event
	event := &unix.EpollEvent{
		Events: unix.EPOLLIN | unix.EPOLLOUT | unix.EPOLLRDHUP | unix.EPOLLET,
		Fd:     int32(fd),
	}

	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_ADD, int(fd), event)
}

func (p *epoll) Del(fd uintptr) error {
	return unix.EpollCtl(p.fd, unix.EPOLL_CTL_DEL, int(fd), nil)
}

func (p *epoll) Events() ([]Event, []Event, error) {
	events, err := p.wait()
	if err != nil {
		return nil, nil, err
	}

	// something received, try it
	var re, we []Event
	for _, event := range events {
		if event.Events&(unix.EPOLLRDHUP) != 0 {
			// closed by peer
			// nothing to do:
			// http://man7.org/linux/man-pages/man7/epoll.7.html
		} else if event.Events&(unix.EPOLLIN) != 0 {
			// ready to read
			re = append(re, toEvent(event))
		} else if event.Events&(unix.EPOLLOUT) != 0 {
			// ready to write
			we = append(we, toEvent(event))
		}
	}

	return re, we, nil
}

func (p *epoll) wait() ([]unix.EpollEvent, error) {
	// TODO max events???
	events := make([]unix.EpollEvent, 32)

	// blocking mode
	n, err := unix.EpollWait(p.fd, events, -1)
	if err != nil {
		return nil, err
	}

	return events[:n], nil
}

// special tool for converting os specific event to interface
func toEvent(event unix.EpollEvent) Event {
	return epollEvent{event}
}
