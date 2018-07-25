package application

import (
	"errors"
	"io"
	"time"

	"github.com/boomfunc/base/conf"
	"github.com/boomfunc/base/server/request"
)

var (
	ErrBadRequest  = errors.New("application: cannot parse request")
	ErrNotFound    = errors.New("application: no route found")
	ErrServerError = errors.New("application: internal server error")
)

type Interface interface {
	// Parse(io.ReadWriter) (request.Interface, error)
	Handle(io.ReadWriter) request.Stat
}

type Packer interface {
	Unpack(io.Reader) (*request.Request, error)
	Pack(io.ReadCloser, io.Writer) (int64, error)
}

type Application struct {
	router *conf.Router
	packer Packer
}

func (app *Application) Handle(rw io.ReadWriter) (stat request.Stat) {
	var req *request.Request
	var begin time.Time
	var err error
	var written int64

	defer func() {
		// end measuring and collect data
		stat.Duration = time.Since(begin)
		stat.Request = req
		stat.Error = err
		stat.Len = written
	}()

	// Start measuring
	begin = time.Now()

	// Parse request
	// TODO ErrBadRequest
	req, err = app.packer.Unpack(rw)
	if err != nil {
		return
	}

	// Resolve view
	// TODO ErrNotFound
	route, err := app.router.Match(req.Url)
	if err != nil {
		return
	}

	// Run pipeline (under app layer)
	pr, pw := io.Pipe()
	go func() {
		// close the writer, so the reader knows there's no more data
		defer pw.Close()

		// BUG: race condition
		// TODO ErrServerError
		err = route.Run(req.Input, pw)
	}()

	// write data to rwc only if all success
	// TODO ErrServerError
	written, err = app.packer.Pack(pr, rw)

	return
}
