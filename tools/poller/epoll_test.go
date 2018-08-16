// +build linux

package poller

import (
	"testing"
	"os"
	"fmt"

	"golang.org/x/sys/unix"
)

func TestNew(t *testing.T) {
	poller, err := New()
	if err != nil {
		t.Fatal(err)
	}

	epoll, ok := poller.(*epoll)
	if !ok {
		t.Fatal("Unexpected poller (expected *epoll)")
	}

	if epoll.fd <= 0 {
		t.Fatal("Invalid poller fd")
	}
}

func TestToEvent(t *testing.T) {
	se := unix.EpollEvent{Fd: 6}

	event := toEvent(se)

	if event.Fd() != 6 {
		t.Fatal("Unexpected event fd")
	}
}

func TestEpollAddDelete(t *testing.T) {
	poller, _ := New()
	r, w, _ := os.Pipe()

	// Success adding
	if err := poller.Add(r.Fd()); err != nil {
		t.Fatal(err)
	}
	if err := poller.Add(w.Fd()); err != nil {
		t.Fatal(err)
	}

	// Duplicate adding -> error
	if err := poller.Add(r.Fd()); err == nil {
		t.Fatal("Expecting error during Add(r)")
	}
	if err := poller.Add(w.Fd()); err == nil {
		t.Fatal("Expecting error during Add(w)")
	}

	// Success deleting
	if err := poller.Del(r.Fd()); err != nil {
		t.Fatal(err)
	}
	if err := poller.Del(w.Fd()); err != nil {
		t.Fatal(err)
	}

	// Duplicate adding -> error
	if err := poller.Del(r.Fd()); err == nil {
		t.Fatal("Expecting error during Del(r)")
	}
	if err := poller.Del(w.Fd()); err == nil {
		t.Fatal("Expecting error during Del(w)")
	}
}

func TestEventsWait(t *testing.T) {
	t.Run("wait", func(t *testing.T) {
		poller, _ := New()
		r, w, _ := os.Pipe()
		poller.Add(r.Fd())
		poller.Add(w.Fd())

		events, err := poller.(*epoll).wait()
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

		events, err = poller.(*epoll).wait()
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

	t.Run("Events", func(t *testing.T) {
		poller, _ := New()
		r, w, _ := os.Pipe()
		poller.Add(r.Fd())
		poller.Add(w.Fd())

		re, we, err := poller.Events()
		if err != nil {
			t.Fatal(err)
		}
		// now only w part is ready for writing -> len == 1
		if len(we) != 1 {
			t.Fatal("Unexpected number of write events, expected 1")
		}
		if len(re) != 0 {
			t.Fatal("Unexpected number of read events, expected 0")
		}
		if we[0].Fd() != w.Fd() {
			t.Fatal("Unexpected event fd")
		}

		// write imitation
		fmt.Fprint(w, "some playload")

		re, we, err = poller.Events()
		if err != nil {
			t.Fatal(err)
		}
		// now w and r parts is ready -> len == 2
		if len(we) != 1 {
			t.Fatal("Unexpected number of write events, expected 1")
		}
		if len(re) != 1 {
			t.Fatal("Unexpected number of read events, expected 1")
		}
		if re[0].Fd() != r.Fd() {
			t.Fatal("Unexpected event fd")
		}
		if we[0].Fd() != w.Fd() {
			t.Fatal("Unexpected event fd")
		}
	})
}
