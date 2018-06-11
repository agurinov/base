// Package pattql is some layer between regex and our simplified url and path format
// TODO wants to isolate into a separate module
package pattql

func Match(pattern, uri string) bool {
	return Regexp(pattern).MatchString(uri)
}
