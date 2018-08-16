package heap

import (
	"net"

	"github.com/boomfunc/base/tools/poller"
)

type PollerHeap struct {
	items  map[uintptr]*net.TCPConn
	poller poller.Interface
}

func NewPollerHeap() (*PollerHeap, error) {
	poller, err := poller.New()
	if err != nil {
		return nil, err
	}

	heap := &PollerHeap{
		items:  make(map[uintptr]*net.TCPConn),
		poller: poller,
	}

	return heap, nil
}

func (ph *PollerHeap) Push(conn *net.TCPConn) error {
	raw, err := conn.SyscallConn()
	if err != nil {
		return err
	}

	// TODO for now it is some kind of workaround
	f := func(fd uintptr) bool {
		// push to epoll
		if err := ph.poller.Add(fd); err == nil {
			// push to map
			ph.items[fd] = conn
		}
		return true
	}

	if err := raw.Read(f); err != nil {
		return err
	}

	return nil
}

func (ph *PollerHeap) Pop() *net.TCPConn {
	for {
		// blocking mode
		re, _, err := ph.poller.Events()
		if err != nil {
			continue
		}

		// iterate over read ready events
		for _, event := range re {
			key := event.Fd()

			if conn, ok := ph.items[key]; ok {
				// rm from epoll
				ph.poller.Del(key)
				// rm from heap and return
				delete(ph.items, key)
				return conn
			}
		}
	}
}
