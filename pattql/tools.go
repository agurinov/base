package pattql

import (
	"fmt"
	"regexp"
	"strings"
)

var bracketsRegexp = regexp.MustCompile("{.*?}") // lazy!!

type Bracket struct {
	source string
}

func BracketFromSource(source string) *Bracket {
	return &Bracket{
		// match will include '{}'
		source: strings.Trim(source, "{}"),
	}
}

func (b Bracket) String() string {
	return fmt.Sprintf("(?:%s)", strings.NewReplacer(
		"*", ".*", // {*} --> {.*}
	).Replace(b.source))
}

func getRegexp(pattern string) *regexp.Regexp {
	expr := bracketsRegexp.ReplaceAllStringFunc(pattern, func(source string) string {
		return BracketFromSource(source).String()
	})

	expr = strings.TrimLeft(expr, "/^") // trim from left all special regex chars
	expr = fmt.Sprintf("^/%s", expr)

	return regexp.MustCompile(expr)
}
