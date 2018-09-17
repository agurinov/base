package tools

import (
	"testing"
	"time"
)

func checkTiming(t *testing.T, timing *Timing, enter, exit int) {
	if l := len(timing.enter); l != enter {
		t.Fatalf("Unexpected len(timing.enter): expected %d, got %d", enter, l)
	}
	if l := len(timing.exit); l != exit {
		t.Fatalf("Unexpected len(timing.exit): expected %d, got %d", exit, l)
	}
}

func TestTiming(t *testing.T) {
	timing := NewTiming()

	t.Run("New", func(t *testing.T) {
		checkTiming(t, timing, 0, 0)
	})

	t.Run("Enter", func(t *testing.T) {
		timing.Enter("foo")
		timing.Enter("bar")

		checkTiming(t, timing, 2, 0)
	})

	t.Run("Exit", func(t *testing.T) {
		timing.Exit("foo")

		checkTiming(t, timing, 2, 1)
	})

	t.Run("Duration", func(t *testing.T) {
		t.Run("NotEntered", func(t *testing.T) {
			if d := timing.Duration("baz"); d != 0 {
				t.Fatal("Unexpected duration, expected 0")
			}
		})

		t.Run("NotClosed", func(t *testing.T) {
			if d := timing.Duration("bar"); d != 0 {
				t.Fatal("Unexpected duration, expected 0")
			}
		})

		t.Run("Real", func(t *testing.T) {
			if d := timing.Duration("foo"); d == 0 {
				t.Fatal("Unexpected duration, must be greater than 0")
			}
		})
	})

	t.Run("String", func(t *testing.T) {
		tmng := &Timing{
			enter: make(map[string]time.Time),
			exit:  make(map[string]time.Time),
		}
		// entering
		tmng.enter["foo"] = time.Date(2018, 9, 17, 10, 0, 0, 0, time.UTC)
		tmng.enter["bar"] = time.Date(2018, 9, 17, 10, 0, 0, 0, time.UTC)

		// exiting
		tmng.exit["bar"] = time.Date(2018, 9, 17, 10, 5, 0, 0, time.UTC)
		tmng.exit["baz"] = time.Date(2018, 9, 17, 10, 5, 0, 0, time.UTC)

		// non entered nodes ignored
		if log := tmng.String(); log != "bar: 5m0s, foo: 0s" {
			t.Fatalf("Unexpected log string, expected: %q, got: %q", "bar: 5m0s, foo: 0s", log)
		}

		// append missing information
		// close foo
		tmng.exit["foo"] = time.Date(2018, 9, 17, 10, 5, 0, 0, time.UTC)
		// open baz
		tmng.enter["baz"] = time.Date(2018, 9, 17, 10, 0, 0, 0, time.UTC)

		if log := tmng.String(); log != "bar: 5m0s, baz: 5m0s, foo: 5m0s" {
			t.Fatalf("Unexpected log string, expected: %q, got: %q", "bar: 5m0s, baz: 5m0s, foo: 5m0s", log)
		}
	})
}
