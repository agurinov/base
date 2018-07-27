package server

import (
	"errors"
	"io"
	"net"

	"github.com/boomfunc/base/conf"
	"github.com/boomfunc/base/server/application"
	"github.com/boomfunc/base/server/context"
	"github.com/boomfunc/base/server/request"
	"github.com/boomfunc/base/server/transport"
	"github.com/boomfunc/base/tools"
)

var (
	ErrWrongContext       = errors.New("server: Context without required key")
	ErrUnknownApplication = errors.New("server: Unknown server application")
	ErrUnknownTransport   = errors.New("server: Unknown server transport")
)

func New(transportName string, applicationName string, ip net.IP, port int, filename string) (*Server, error) {
	// Phase 1. Prepare light application layer things
	// router
	router, err := conf.LoadFile(filename)
	if err != nil {
		// cannot load server config
		return nil, err
	}

	// Phase 2. Prepare main application layer
	var app application.Interface

	switch applicationName {
	case "http":
		app = application.HTTP(router)
	case "json":
		app = application.JSON(router)
	default:
		return nil, ErrUnknownApplication
	}
	inputCh := make(chan io.ReadWriteCloser)
	errCh := make(chan error)
	outputCh := make(chan request.Stat)

	// Phase 3. Prepare transport layer
	var tr transport.Interface

	switch transportName {
	case "tcp":
		tr, err = transport.TCP(ip, port)
		if err != nil {
			return nil, err
		}
		tr.Connect(inputCh, errCh)
	default:
		return nil, ErrUnknownTransport
	}

	srv := new(Server)
	// flow data
	srv.transport = tr
	srv.app = app
	// channels
	srv.inputCh = inputCh
	srv.errCh = errCh
	srv.outputCh = outputCh

	return srv, nil
}

// this function will be passed to dispatcher system
// and will be run at parallel
func HandleTask(task Task) {
	srvInterface, err := context.GetMeta(task.ctx, "srv")
	if err != nil {
		tools.FatalLog(err)
	}

	srv, ok := srvInterface.(*Server)
	if !ok {
		tools.FatalLog(ErrWrongContext)
	}

	defer func() {
		if r := recover(); r != nil {
			switch typed := r.(type) {
			case error:
				srv.errCh <- typed
			case string:
				srv.errCh <- errors.New(typed)
			}
		}
	}()

	defer task.input.Close()

	srv.outputCh <- srv.app.Handle(task.ctx, task.input)
}
