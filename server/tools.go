package server

import (
	"errors"
	"io"
	"net"

	"github.com/boomfunc/base/conf"
	"github.com/boomfunc/base/server/application"
	"github.com/boomfunc/base/server/request"
	"github.com/boomfunc/base/server/transport"
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
		return nil, errors.New("server: Unknown server application")
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
		return nil, errors.New("server: Unknown server transport")
	}

	srv := &Server{
		transport: tr,
		app:       app,
		inputCh:   inputCh,
		errCh:     errCh,
		outputCh:  outputCh,
	}
	return srv, nil
}
