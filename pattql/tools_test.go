package pattql

import (
	"fmt"
	"testing"
)

func TestBracketString(t *testing.T) {
	tableTests := []struct {
		source string // source
		out    string // expected value of String() method
	}{
		{"", "(?:)"},
		{"foo|bar|baz", "(?:foo|bar|baz)"},
		{"foo|bar|*", "(?:foo|bar|.*)"},
	}

	for i, tt := range tableTests {
		bracket := Bracket{tt.source}

		t.Run(fmt.Sprintf("%d.String()", i), func(t *testing.T) {
			if bracket.String() != tt.out {
				t.Errorf("Expected %q, got %q", tt.out, bracket.String())
			}
		})

		t.Run(fmt.Sprintf("fmt.Sprint(%d)", i), func(t *testing.T) {
			if fmt.Sprint(bracket) != tt.out {
				t.Errorf("Expected %q, got %q", tt.out, fmt.Sprint(bracket))
			}
		})
	}
}

func TestBracketFromSource(t *testing.T) {
	tableTests := []struct {
		source string // source
		out    string // expected value of String() method
	}{
		{"", ""},
		{"{}", ""},
		{"{foo|bar|baz}", "foo|bar|baz"},
		{"{foo|bar|baz", "foo|bar|baz"},
	}

	for i, tt := range tableTests {
		bracket := bracketFromSource(tt.source)

		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if bracket.source != tt.out {
				t.Errorf("Expected %q, got %q", tt.out, bracket.source)
			}
		})
	}
}

func TestRegexp(t *testing.T) {
	tableTests := []struct {
		pattern string // source
		out     string // expected value of String() method
	}{
		{"", "^/?"},
		{"^/{data|static}/{*}.{jpg|png|jpeg}", "^/?(?:data|static)/(?:.*).(?:jpg|png|jpeg)"},
		{"/{data|static}/{*}.{jpg|png|jpeg}", "^/?(?:data|static)/(?:.*).(?:jpg|png|jpeg)"},
		{"^{data|static}/{*}.{jpg|png|jpeg}", "^/?(?:data|static)/(?:.*).(?:jpg|png|jpeg)"},
		{"{data|static}/{*}.{jpg|png|jpeg}", "^/?(?:data|static)/(?:.*).(?:jpg|png|jpeg)"},
	}

	for i, tt := range tableTests {
		finalExpr := Regexp(tt.pattern).String()

		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if finalExpr != tt.out {
				t.Errorf("Expected %q, got %q", tt.out, finalExpr)
			}
		})
	}
}
