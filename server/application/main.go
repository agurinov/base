package application

import (
	"context"
	"errors"
	"io"

	"github.com/boomfunc/base/conf"
	"github.com/boomfunc/base/server/flow"
)

var (
	ErrBadRequest  = errors.New("server/application: cannot parse request")
	ErrServerError = errors.New("server/application: internal server error")
)

type Interface interface {
	Handle(*flow.Data)
}

type Packer interface {
	Unpack(context.Context, io.Reader) (*flow.Request, error)
	Pack(io.Reader, io.Writer) (int64, error)
}

type Application struct {
	router *conf.Router
	packer Packer
}

func (app *Application) Handle(fl *flow.Data) {
	var req *flow.Request
	var err error
	var written int64

	defer func() {
		fl.Stat.Request = req
		fl.Stat.Error = err
		fl.Stat.Len = written
	}()

	// Parse request
	// fill context meta part and q part
	// TODO ErrBadRequest
	req, err = app.packer.Unpack(fl.Ctx, fl.RWC)
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
		err = route.Run(fl.Ctx, req.Input, pw)
	}()

	// write data to rwc only if all success
	// TODO ErrServerError
	written, err = app.packer.Pack(pr, fl.RWC)

	return
}
