package server

import (
	"errors"
	"io"

	"github.com/boomfunc/base/server/application"
	"github.com/boomfunc/base/server/request"
	"github.com/boomfunc/base/server/transport"
)

type Server struct {
	transport transport.Interface
	app       application.Interface

	inputCh  chan io.ReadWriteCloser
	errCh    chan error
	outputCh chan request.Stat
}

// this function will be passed to dispatcher system
// and will be run at parallel
// TODO move to interface as a link to the dispatcher system
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
	// TODO https://rcrowley.org/articles/golang-graceful-stop.html

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
				task := func() { srv.handle(input) }
				TaskChannel <- task

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

	PerfomanceLog(numWorkers)

	// GOROUTINE 1 (main)
	// This is thread blocking procedure - infinity loop
	srv.transport.Serve()
}
