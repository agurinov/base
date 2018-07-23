package transport

import (
	"io"
)

type Interface interface {
	Connect(inputCh chan io.ReadWriteCloser, errCh chan error)
	Serve()
}
