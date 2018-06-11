package pattql

import (
	"fmt"
	"testing"
)

func BenchmarkBracketString(b *testing.B) {
	tableTests := []struct {
		source string // source
		out    string // expected value of String() method
	}{
		{"", "(?:)"},
		{"foo|bar|baz", "(?:foo|bar|baz)"},
		{"foo|bar|*", "(?:foo|bar|.*)"},
	}

	for i, tt := range tableTests {
		bracket := bracketFromSource(tt.source)

		b.Run(fmt.Sprintf("%d.String()", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bracket.String()
			}
		})
	}
}

func BenchmarkBracketFromSource(b *testing.B) {
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
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bracketFromSource(tt.source)
			}
		})
	}
}

func BenchmarkRegexp(b *testing.B) {
	tableTests := []struct {
		pattern string // source
		out     string // expected value of String() method
	}{
		{"", "^/"},
		{"^/{data|static}/{*}.{jpg|png|jpeg}", "^/(?:data|static)/(?:.*).(?:jpg|png|jpeg)"},
		{"/{data|static}/{*}.{jpg|png|jpeg}", "^/(?:data|static)/(?:.*).(?:jpg|png|jpeg)"},
		{"^{data|static}/{*}.{jpg|png|jpeg}", "^/(?:data|static)/(?:.*).(?:jpg|png|jpeg)"},
		{"{data|static}/{*}.{jpg|png|jpeg}", "^/(?:data|static)/(?:.*).(?:jpg|png|jpeg)"},
	}

	for i, tt := range tableTests {
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Regexp(tt.pattern)
			}
		})
	}
}
