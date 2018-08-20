// +build darwin dragonfly freebsd netbsd openbsd

package poller

import (
	"fmt"
	"os"
	"testing"

	"golang.org/x/sys/unix"
)

func TestNew(t *testing.T) {
	poller, err := New()
	if err != nil {
		t.Fatal(err)
	}

	kqueue, ok := poller.(*kqueue)
	if !ok {
		t.Fatal("Unexpected poller (expected *kqueue)")
	}

	if kqueue.fd <= 0 {
		t.Fatal("Invalid poller fd")
	}
}

func TestToEvent(t *testing.T) {
	se := unix.Kevent_t{Ident: 6}

	event := toEvent(se)

	if event.Fd() != 6 {
		t.Fatal("Unexpected event fd")
	}
}

func TestEventsWait(t *testing.T) {
	t.Run("wait", func(t *testing.T) {
		poller, _ := New()
		r, w, _ := os.Pipe()
		poller.Add(r.Fd())
		poller.Add(w.Fd())

		events, err := poller.(*kqueue).wait()
		if err != nil {
			t.Fatal(err)
		}
		// now only w part is ready for writing -> len == 1
		if len(events) != 1 {
			t.Fatal("Unexpected number of events, expected 1")
		}
		if events[0].Ident != uint64(w.Fd()) {
			t.Fatal("Unexpected event fd")
		}

		fmt.Fprint(w, "some playload")

		events, err = poller.(*kqueue).wait()
		if err != nil {
			t.Fatal(err)
		}
		// now w and r parts is ready -> len == 2
		if len(events) != 2 {
			t.Fatal("Unexpected number of events, expected 2")
		}
		if events[0].Ident != uint64(w.Fd()) {
			t.Fatal("Unexpected event fd")
		}
		if events[1].Ident != uint64(r.Fd()) {
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
