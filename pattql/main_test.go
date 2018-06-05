package pattql

import (
	"fmt"
	"testing"
)

func TestMatch(t *testing.T) {
	pattern := "{data|static}/{*}.{jpg|png|gif|jpeg}"

	tableTests := []struct {
		uri string // uri for matching
		out bool   // does match
	}{
		{"", false},
		{"/foo/data/foobar.jpeg", false},
		{"data/barbaz.gif", true},
		{"/data/foobar.png", true},
		{"nedata/skr.jpg", false},
	}

	for i, tt := range tableTests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if matched := Match(pattern, tt.uri); matched != tt.out {
				t.Errorf("Expected \"%t\", got \"%t\"", tt.out, matched)
			}
		})
	}
}
