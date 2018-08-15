package transport

import (
	"net"

	"golang.org/x/sys/unix"
)

type ConnHeap struct {
	items  map[uintptr]*net.TCPConn
	poller *epoll
}

func (h *ConnHeap) Push(conn *net.TCPConn) error {
	raw, err := conn.SyscallConn()
	if err != nil {
		return err
	}

	f := func(fd uintptr) bool {
		// push to epoll
		if err := h.poller.Add(int32(fd)); err == nil {
			// push to map
			h.items[fd] = conn
		}
		return true
	}

	if err := raw.Read(f); err != nil {
		return err
	}

	return nil
}

func (h *ConnHeap) Pop() *net.TCPConn {
	for {
		// blocking mode
		events, err := h.poller.Wait()
		if err != nil {
			continue
		}
		for _, event := range events {
			// check for 'event read ready'
			if event.Events&(unix.EPOLLRDHUP|unix.EPOLLERR) != 0 {
				// closed, error, etc -> no return, just pop from epoll
			} else if event.Events&(unix.EPOLLIN) != 0 {
				// no error, data is coming
				key := uintptr(event.Fd)
				if conn, ok := h.items[key]; ok {
					// rm from epoll
					h.poller.Del(event.Fd)
					// rm from heap and return
					delete(h.items, key)
					return conn
				}
			}
		}
	}
}
