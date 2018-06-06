package pattql

import (
	"fmt"
	"regexp"
	"strings"
)

var bracketsRegexp = regexp.MustCompile("{.*?}") // lazy!!

// Bracket is some dynamic object in our pattern
// that must be converted to valid regex expression
type Bracket struct {
	source string
}

func bracketFromSource(source string) *Bracket {
	return &Bracket{
		// match will include '{}'
		source: strings.Trim(source, "{}"),
	}
}

// String returns 'regex ready' source representation
func (b Bracket) String() string {
	return fmt.Sprintf("(?:%s)", strings.NewReplacer(
		"*", ".*", // {*} --> {.*}
	).Replace(b.source))
}

// getRegexp returns complex regex for incoming pattern
// some several bases
// 1. expression will always starts with ^/
// 2. bracket source will be replaced by some condition (look at Bracket.String())
func getRegexp(pattern string) *regexp.Regexp {
	expr := bracketsRegexp.ReplaceAllStringFunc(pattern, func(source string) string {
		return bracketFromSource(source).String()
	})

	expr = strings.TrimLeft(expr, "/^") // trim from left all special regex chars
	expr = fmt.Sprintf("^/?%s", expr)   // string begin with non required leading /

	return regexp.MustCompile(expr)
}
