package chronometer

import (
	"testing"
	"time"
)

func TestNodePublic(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		node := NewNode()

		if node.enter.IsZero() {
			t.Fatalf("node.enter. Unexpected zero dt")
		}

		if !node.exit.IsZero() {
			t.Fatalf("node.exit. Unexpected non zero dt")
		}
	})

	t.Run("Duration", func(t *testing.T) {
		t.Run("NotEntered", func(t *testing.T) {
			node := NewNode()
			node.enter = time.Time{}

			if node.Duration() != 0 {
				t.Fatal("Unexpected duration, expected 0")
			}
		})

		t.Run("NotClosed", func(t *testing.T) {
			node := NewNode()

			if node.Duration() != 0 {
				t.Fatal("Unexpected duration, expected 0")
			}
		})

		t.Run("Both", func(t *testing.T) {
			node := Node{}

			if node.Duration() != 0 {
				t.Fatal("Unexpected duration, expected 0")
			}
		})

		t.Run("Concrete", func(t *testing.T) {
			node := Node{
				enter: time.Date(2018, 9, 17, 10, 0, 0, 0, time.UTC),
				exit:  time.Date(2018, 9, 17, 10, 5, 0, 0, time.UTC),
			}

			if d := node.Duration(); d.String() != "5m0s" {
				t.Fatalf("Unexpected duration, expected: %q, got: %q", "5m0s", d.String())
			}
		})
	})

	t.Run("Enter", func(t *testing.T) {
		node := Node{}
		if !node.enter.IsZero() {
			t.Fatalf("node.enter. Unexpected non zero dt")
		}
		node.Enter()
		if node.enter.IsZero() {
			t.Fatalf("node.enter. Unexpected zero dt")
		}
	})

	t.Run("Enter", func(t *testing.T) {
		node := Node{}
		if !node.exit.IsZero() {
			t.Fatalf("node.enter. Unexpected non zero dt")
		}
		node.Exit()
		if node.exit.IsZero() {
			t.Fatalf("node.enter. Unexpected zero dt")
		}
	})
}
