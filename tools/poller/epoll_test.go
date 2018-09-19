// +build linux

package poller

import (
	"fmt"
	"os"
	"testing"

	"golang.org/x/sys/unix"
)

func TestPollerPrivate(t *testing.T) {
	poller, err := New()
	epoll, ok := poller.(*epoll)

	t.Run("New", func(t *testing.T) {
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatal("Unexpected poller type (expected *epoll)")
		}
		if epoll.fd <= 0 {
			t.Fatal("Invalid poller fd")
		}
	})

	t.Run("wait", func(t *testing.T) {
		r, w, _ := os.Pipe()
		epoll.Add(r.Fd())
		epoll.Add(w.Fd())

		events, err := epoll.wait()
		if err != nil {
			t.Fatal(err)
		}
		// now only w part is ready for writing -> len == 1
		if len(events) != 1 {
			t.Fatal("Unexpected number of events, expected 1")
		}
		if events[0].Fd != int32(w.Fd()) {
			t.Fatal("Unexpected event fd")
		}

		fmt.Fprint(w, "some playload")

		events, err = epoll.wait()
		if err != nil {
			t.Fatal(err)
		}
		// now w and r parts is ready -> len == 2
		if len(events) != 2 {
			t.Fatal("Unexpected number of events, expected 2")
		}
		if events[0].Fd != int32(w.Fd()) {
			t.Fatal("Unexpected event fd")
		}
		if events[1].Fd != int32(r.Fd()) {
			t.Fatal("Unexpected event fd")
		}
	})
}

func TestToEvent(t *testing.T) {
	se := unix.EpollEvent{Fd: 6}

	event := toEvent(se)

	if event.Fd() != 6 {
		t.Fatal("Unexpected event fd")
	}
}
