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
	Handle(*flow.Data)
}

type Packer interface {
	Unpack(context.Context, io.Reader) (*request.Request, error)
	Pack(io.Reader, io.Writer) (int64, error)
}

type Application struct {
	router *conf.Router
	packer Packer
}

func (app *Application) Handle(flow *flow.Data) {
	var req *request.Request
	var err error
	var written int64

	defer func() {
		flow.Stat.Request = req
		flow.Stat.Error = err
		flow.Stat.Len = written + 12
	}()

	// Parse request
	// fill context meta part and q part
	// TODO ErrBadRequest
	req, err = app.packer.Unpack(flow.Ctx, flow.Input)
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
		err = route.Run(flow.Ctx, req.Input, pw)
	}()

	// write data to rwc only if all success
	// TODO ErrServerError
	written, err = app.packer.Pack(pr, flow.Input)

	return
}
