// Package pattql is some layer between regex and our simplified url and path format
package pattql

func Match(pattern, uri string) bool {
	return getRegexp(pattern).MatchString(uri)
}
