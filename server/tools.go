package server

import (
	"errors"
	"net"

	"github.com/boomfunc/base/conf"
	"github.com/boomfunc/base/server/application"
	"github.com/boomfunc/base/server/dispatcher"
	"github.com/boomfunc/base/server/flow"
	"github.com/boomfunc/base/server/transport"
	"github.com/boomfunc/base/tools/poller"
)

var (
	ErrWrongContext       = errors.New("server: Context without required key")
	ErrUnknownApplication = errors.New("server: Unknown server application")
	ErrUnknownTransport   = errors.New("server: Unknown server transport")
)

func New(transportName string, applicationName string, workers int, ip net.IP, port int, config string) (*Server, error) {
	// Phase 1. Prepare light application layer things
	// router
	router, err := conf.LoadFile(config)
	if err != nil {
		// cannot load server config
		return nil, err
	}

	// Phase 2. Prepare application layer
	var app application.Interface
	switch applicationName {
	case "http":
		app = application.HTTP(router)
	case "json":
		app = application.JSON(router)
	default:
		return nil, ErrUnknownApplication
	}

	// Phase 3. Prepare transport layer
	var tr transport.Interface
	switch transportName {
	case "tcp":
		tr, err = transport.TCP(ip, port)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrUnknownTransport
	}

	// Phase 4. transport layer recognized, we can create support data for connection layers
	heap, err := poller.Heap()
	if err != nil {
		return nil, err
	}
	errCh := make(chan error)
	outputCh := make(chan *flow.Data)
	tr.Connect(heap, errCh)

	// Phase 5. Create server
	srv := new(Server)
	// flow data
	srv.transport = tr
	srv.app = app
	srv.dispatcher = dispatcher.New(workers)
	// channels and data storages
	srv.heap = heap
	srv.errCh = errCh
	srv.outputCh = outputCh

	return srv, nil
}
