// +build !linux,!darwin,!dragonfly,!freebsd,!netbsd,!openbsd

package poller

import (
	"errors"
)

var (
	ErrNotImplemented = errors.New("Poller not implemented for this OS")
)

func New() (Interface, error) {
	return nil, ErrNotImplemented
}
