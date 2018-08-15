package transport

import (
	"io"
)

type Interface interface {
	Connect(chan io.ReadWriteCloser, chan error)
	Serve()
}
