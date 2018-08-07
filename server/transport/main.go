package transport

import (
	"io"
)

type Interface interface {
	Connect(chan io.ReadWriteCloser, chan error)
	Serve()
}

type Poller interface {
	Add(fd uintptr)
	Del(fd uintptr)
	Wait(fd uintptr)
}
