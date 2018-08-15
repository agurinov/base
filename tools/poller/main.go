package poller

import (
	"golang.org/x/sys/unix"
)

type Interface interface {
	Add(fd int32) error
	Del(fd int32) error
	Events() (re []unix.EpollEvent, we []unix.EpollEvent, err error)
}
