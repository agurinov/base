package server

import (
	"errors"
	"net"
	"runtime"

	"github.com/boomfunc/base/conf"
	"github.com/boomfunc/base/server/application"
	"github.com/boomfunc/base/server/request"
	"github.com/boomfunc/base/server/transport"
	"github.com/boomfunc/log"
)

type Server struct {
	transport transport.Interface
	app       application.Interface

	connCh     chan net.Conn
	errCh      chan error
	responseCh chan request.Response

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

func (srv *Server) Serve(numWorkers int) {
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
	// TODO defer ch.Close()
	// TODO defer s.conn.Close()
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
	// https://rcrowley.org/articles/golang-graceful-stop.html

	// GOROUTINE 2 (dispatcher - listen RequestChannel)
	NewDispatcher(numWorkers).Run()

	// GOROUTINE 3 (listen server channels)
	go func() {
		for {
			select {
			case err := <-srv.errCh:
				if err != nil {
					ErrorLog(err)
				}

			case conn := <-srv.connCh:
				// connection from transport
				// send to dispatcher's queue
				RequestChannel <- Request{conn, srv}

			case response := <-srv.responseCh:
				// ready response from worker
				// log ANY kind of result
				AccessLog(response)
				// and errors
				if err := response.Error; err != nil {
					// TODO not good, find better solution
					// TODO repeat line56
					ErrorLog(err)
				}
			}
		}
	}()

	log.Debugf("Spawned %d goroutines", runtime.NumGoroutine())
	if runtime.NumGoroutine() != numWorkers+3 {
		log.Warnf(
			"Unexpected number of goroutines, possibly an issue. Expected: %d, Got: %d",
			numWorkers+3,
			runtime.NumGoroutine(),
		)
	}
	log.Debugf("Detected %d CPU cores", runtime.NumCPU())
	if runtime.NumCPU() < numWorkers {
		log.Warnf(
			"Possible overloading of CPU cores. Detected: %[1]d CPU. Recommended worker number: %[1]d (Current: %[2]d)",
			runtime.NumCPU(), numWorkers,
		)
	} else if runtime.NumCPU() > numWorkers {
		log.Warnf(
			"Possible performance improvements. Increase worker number. Detected: %[1]d CPU. Recommended worker number: %[1]d (Current: %[2]d)",
			runtime.NumCPU(), numWorkers,
		)
	}

	// GOROUTINE 1 (main)
	// This is thread blocking procedure - infinity loop
	srv.transport.Serve()
}

func New(transportName string, applicationName string, ip net.IP, port int, filename string) (*Server, error) {
	// Phase 1. Prepare light application layer things
	// router
	router, err := conf.LoadFile(filename)
	if err != nil {
		// cannot load server config
		return nil, err
	}

	// Phase 2. Prepare main application layer
	// TODO
	app := application.New(router)
	connCh := make(chan net.Conn)
	errCh := make(chan error)
	responseCh := make(chan request.Response)

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
		transport:  tr,
		app:        app,
		connCh:     connCh,
		errCh:      errCh,
		responseCh: responseCh,
		router:     router,
	}
	return srv, nil
}
