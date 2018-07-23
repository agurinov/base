package application

import (
	"io"

	"github.com/boomfunc/base/server/request"
)

type Interface interface {
	Parse(request io.ReadWriter) (request.Interface, error)
	Handle(request request.Interface) request.Stat
}
