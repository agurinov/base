package server

import (
	"context"
	"io"

	"github.com/boomfunc/base/server/application"
	"github.com/boomfunc/base/server/request"
	"github.com/boomfunc/base/server/transport"
	"github.com/boomfunc/base/tools"
)

type Server struct {
	transport transport.Interface
	app       application.Interface
	ctx       context.Context

	inputCh  chan io.ReadWriteCloser
	errCh    chan error
	outputCh chan request.Stat
}

func (srv *Server) listenOS() {

}

// this function listen all server channels and proxying
// errors to log
// io.RWC to dispatcher system
// response.Stat to log and check for errors additionally
func (srv *Server) listenCh() {
	for {
		select {
		case err := <-srv.errCh:
			if err != nil {
				tools.ErrorLog(err)
			}

		case input := <-srv.inputCh:
			// input from transport layer (conn, file socket, or something else)
			// create flow context, fill from server info
			// send to dispatcher's queue
			TaskChannel <- Task{srv.ctx, input}

		case stat := <-srv.outputCh:
			// ready response from dispatcher system
			// log ANY kind of result
			AccessLog(stat)
			// and errors
			if err := stat.Error; err != nil {
				// TODO not good, find better solution
				// TODO repeat line56
				tools.ErrorLog(err)
			}
		}
	}
}

func (srv *Server) Serve(numWorkers int) {
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
	// TODO defer ch.Close()
	// TODO defer s.conn.Close()
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
	// TODO https://rcrowley.org/articles/golang-graceful-stop.html

	// GOROUTINE 2 (dispatcher - listen TaskChannel)
	dispatcher := NewDispatcher(numWorkers)
	go dispatcher.Dispatch()

	// GOROUTINE 3 (listen server channels)
	go srv.listenCh()

	// TODO GOROUTINE 4 (listen for os signals and gracefully close server)
	// go srv.listenOS()

	// Here we can test some of our system requirements and performance recommendations
	PerfomanceLog(numWorkers)

	// GOROUTINE 1 (main) - this goroutine
	// This is thread blocking procedure - infinity loop
	srv.transport.Serve()
}
