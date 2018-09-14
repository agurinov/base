package server

import (
	"github.com/boomfunc/base/server/application"
	"github.com/boomfunc/base/server/context"
	"github.com/boomfunc/base/server/dispatcher"
	"github.com/boomfunc/base/server/flow"
	"github.com/boomfunc/base/server/request"
	"github.com/boomfunc/base/server/transport"
	"github.com/boomfunc/base/tools"
)

type Server struct {
	transport  transport.Interface
	app        application.Interface
	dispatcher *dispatcher.Dispatcher

	inputCh  chan *flow.Data
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

		case flow := <-srv.inputCh:
			go func() {
				// input from transport layer (conn, file socket, or something else)
				// try to fetch empty worker (to be precise, his channel)
				// blocking mode!
				flow.Timing.Enter("dispatcher")
				taskChannel := srv.dispatcher.FreeWorkerTaskChannel()
				flow.Timing.Exit("dispatcher")
				// create request own flow context, fill server part of data
				context.SetMeta(flow.Ctx, "srv", srv)
				// send to worker's channel
				taskChannel <- Task{flow}
			}()

		case stat := <-srv.outputCh:
			// ready response from dispatcher system
			// log ANY kind of result
			AccessLog(stat)
			// and errors
			if err := stat.Error; err != nil {
				go func() {
					srv.errCh <- stat.Error
				}()
			}
		}
	}
}

func (srv *Server) Serve() {
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
	// TODO defer ch.Close()
	// TODO defer s.conn.Close()
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
	// TODO https://rcrowley.org/articles/golang-graceful-stop.html

	// create real worker instances
	srv.dispatcher.Prepare()

	// GOROUTINE 2 (listen server channels)
	go srv.listenCh()

	// TODO GOROUTINE 3 (listen for os signals and gracefully close server)
	// go srv.listenOS()

	// Here we can test some of our system requirements and performance recommendations
	PerformanceLog(srv.dispatcher.MaxWorkers)

	// GOROUTINE 1 (main) - this goroutine
	// This is thread blocking procedure - infinity loop
	srv.transport.Serve()
}
