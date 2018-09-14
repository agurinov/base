package application

import (
	"context"
	"errors"
	"io"

	"github.com/boomfunc/base/conf"
	// srvctx "github.com/boomfunc/base/server/context"
	"github.com/boomfunc/base/server/flow"
	"github.com/boomfunc/base/server/request"
	// "github.com/boomfunc/log"
)

var (
	ErrBadRequest  = errors.New("server/application: cannot parse request")
	ErrServerError = errors.New("server/application: internal server error")
)

type Interface interface {
	Handle(context.Context, io.ReadWriter) flow.Stat
}

type Packer interface {
	Unpack(context.Context, io.Reader) (*request.Request, error)
	Pack(io.Reader, io.Writer) (int64, error)
}

type Application struct {
	router *conf.Router
	packer Packer
}

func (app *Application) Handle(ctx context.Context, rw io.ReadWriter) (stat flow.Stat) {
	var req *request.Request
	var err error
	var written int64

	defer func() {
		stat.Request = req
		stat.Error = err
		stat.Len = written
	}()

	// Parse request
	// fill context meta part and q part
	// TODO ErrBadRequest
	req, err = app.packer.Unpack(ctx, rw)
	if err != nil {
		return
	}

	// ip, err := srvctx.GetMeta(ctx, "ip")
	// log.Debug("IP:", ip)
	//
	// values, err := srvctx.Values(ctx)
	// log.Debug("Q:", values.Q)

	// Resolve view
	// TODO conf.ErrRouteNotFound
	// fill context url
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
		err = route.Run(ctx, req.Input, pw)
	}()

	// write data to rwc only if all success
	// TODO ErrServerError
	written, err = app.packer.Pack(pr, rw)

	return
}
