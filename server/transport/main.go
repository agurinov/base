package transport

import (
	"io"

	"golang.org/x/sys/unix"
)

type Interface interface {
	Connect(chan io.ReadWriteCloser, chan error)
	Serve()
}

type Poller interface {
	Add(fd int32) error
	Del(fd int32) error
	Wait() ([]unix.EpollEvent, error)
}
