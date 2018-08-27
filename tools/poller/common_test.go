package poller

import (
	"os"
	"testing"
)

func TestPollerAddDelete(t *testing.T) {
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
