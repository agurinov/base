package server

import (
	"errors"
	"net"

	"github.com/boomfunc/base/conf"
	"github.com/boomfunc/base/server/application"
	"github.com/boomfunc/base/server/request"
	"github.com/boomfunc/base/server/transport"
	"github.com/boomfunc/log"
)

type Server struct {
	transport transport.Interface
	app       application.Interface

	connCh chan net.Conn
	errCh  chan error

	// move to application layer
	router *conf.Router
}

// // logging and error handling block
// // this defer must be invoked last (first in) for recovering all available panics and errors
// defer func() {
// 	// var status = "SUCCESS"
//
// 	if err != nil {
// 		ErrorLog(err)
// 		// status = "ERROR"
// 	}
// 	// log ANY kind result
// 	// AccessLog(req, status)
// }()

func (srv *Server) Serve() {
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
	// TODO defer ch.Close()
	// TODO defer s.conn.Close()
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
	// https://rcrowley.org/articles/golang-graceful-stop.html

	// First goroutine - listen RequestChannel
	NewDispatcher(4).Run()

	// second goroutine - listen transport channels
	go func() {
		for {
			select {
			case err := <-srv.errCh:
				// error from transport
				log.Error(err)
			case conn := <-srv.connCh:
				// connection from transport
				r := Request{
					under:  request.New(conn),
					server: srv,
				}
				RequestChannel <- r
			}
		}
	}()

	// This is thread blocking procedure - infinity loop
	srv.transport.Serve()
}

func New(transportName string, ip net.IP, port int, filename string) (*Server, error) {
	// Phase 1. Prepare light application layer things
	// router
	router, err := conf.LoadFile(filename)
	if err != nil {
		// cannot load server config
		return nil, err
	}

	// Phase 2. Prepare main application layer
	connCh := make(chan net.Conn)
	errCh := make(chan error)

	// Phase 3. Prepare transport layer
	var tr transport.Interface

	switch transportName {
	case "tcp":
		tr, err = transport.TCP(ip, port, connCh, errCh)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("server: Unknown server transport")
	}

	srv := &Server{
		transport: tr,
		connCh:    connCh,
		errCh:     errCh,
		router:    router,
	}
	return srv, nil
}
