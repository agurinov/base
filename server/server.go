package server

import (
	"errors"
	"io"
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

	inputCh chan io.ReadWriteCloser
	errCh   chan error
	// TODO in next iteration it will be perfect Type)
	outputCh chan request.Stat
}

// this function will be passed to dispatcher system
// and will be run at parallel
// TODO move to interface as a link to the dispatcher ыныеуь
func (srv *Server) handle(input io.ReadWriteCloser) {
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

	defer input.Close()

	srv.outputCh <- srv.app.Handle(input)
}

func (srv *Server) Serve(numWorkers int) {
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
	// TODO defer ch.Close()
	// TODO defer s.conn.Close()
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
	// https://rcrowley.org/articles/golang-graceful-stop.html

	// GOROUTINE 2 (dispatcher - listen TaskChannel)
	NewDispatcher(numWorkers).Run()

	// GOROUTINE 3 (listen server channels)
	go func() {
		for {
			select {
			case err := <-srv.errCh:
				if err != nil {
					ErrorLog(err)
				}

			case input := <-srv.inputCh:
				// input from transport layer (conn, file socket, or something else)
				// transform to ServerRequest
				// send to dispatcher's queue
				TaskChannel <- func() {
					srv.handle(input)
				}

			case stat := <-srv.outputCh:
				// ready response from dispatcher system
				// log ANY kind of result
				AccessLog(stat)
				// and errors
				if err := stat.Error; err != nil {
					// TODO not good, find better solution
					// TODO repeat line56
					ErrorLog(err)
				}
			}
		}
	}()

	// TODO https://insights.sei.cmu.edu/sei_blog/2017/08/multicore-and-virtualization-an-introduction.html
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
