package chronometer

import (
	"testing"
	"time"
)

func checkChronometerLen(t *testing.T, ch *Chronometer, nodes int) {
	if l := len(*ch); l != nodes {
		t.Fatalf("Unexpected len(chronometer): expected %d, got %d", nodes, l)
	}
}

func TestChronometer(t *testing.T) {
	ch := New()

	t.Run("New", func(t *testing.T) {
		checkChronometerLen(t, ch, 0)
	})

	t.Run("Enter", func(t *testing.T) {
		ch.Enter("foo")
		ch.Enter("bar")

		checkChronometerLen(t, ch, 2)
	})

	t.Run("Exit", func(t *testing.T) {
		ch.Exit("foo")
		ch.Exit("baz") // non existence node ignored

		checkChronometerLen(t, ch, 2)
	})

	t.Run("AddNode", func(t *testing.T) {
		ch.AddNode("lolkek", NewNode())

		checkChronometerLen(t, ch, 3)
	})

	t.Run("String", func(t *testing.T) {
		ch := New()

		foo := &Node{enter: time.Date(2018, 9, 17, 10, 0, 0, 0, time.UTC)}
		bar := &Node{enter: time.Date(2018, 9, 17, 10, 0, 0, 0, time.UTC), exit: time.Date(2018, 9, 17, 10, 5, 0, 0, time.UTC)}
		baz := &Node{exit: time.Date(2018, 9, 17, 10, 5, 0, 0, time.UTC)}

		ch.AddNode("foo", foo)
		ch.AddNode("bar", bar)
		ch.AddNode("baz", baz)

		// non entered or not closed nodes ignored
		if log := ch.String(); log != "bar: 5m0s" {
			t.Fatalf("Unexpected log string, expected: %q, got: %q", "bar: 5m0s", log)
		}

		// append missing information
		// close foo
		foo.exit = time.Date(2018, 9, 17, 10, 5, 0, 0, time.UTC)
		// open baz
		baz.enter = time.Date(2018, 9, 17, 10, 0, 0, 0, time.UTC)

		if log := ch.String(); log != "bar: 5m0s, baz: 5m0s, foo: 5m0s" {
			t.Fatalf("Unexpected log string, expected: %q, got: %q", "bar: 5m0s, baz: 5m0s, foo: 5m0s", log)
		}
	})
}
