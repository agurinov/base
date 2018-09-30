package server

import (
	"container/heap"

	"github.com/boomfunc/base/server/application"
	"github.com/boomfunc/base/server/context"
	"github.com/boomfunc/base/server/dispatcher"
	"github.com/boomfunc/base/server/flow"
	"github.com/boomfunc/base/server/transport"
	"github.com/boomfunc/base/tools"
	"github.com/boomfunc/base/tools/chronometer"
)

type Server struct {
	transport  transport.Interface // TODO https://github.com/boomfunc/base/issues/20
	app        application.Interface
	dispatcher *dispatcher.Dispatcher

	heap     heap.Interface
	errCh    chan error
	outputCh chan *flow.Data
}

func (srv *Server) engine() {
	for {
		// Phase 1. get worker
		// try to fetch empty worker (to be precise, his channel)
		// blocking mode!
		node := chronometer.NewNode()
		taskChannel := srv.dispatcher.FreeWorkerTaskChannel()
		node.Exit()

		// Phase 2. Obtain socket with data from heap/poller
		// blocking mode!
		if flow, ok := heap.Pop(srv.heap).(*flow.Data); ok {
			flow.Chronometer.Exit("transport")
			flow.Chronometer.AddNode("dispatcher", node)
			context.SetMeta(flow.Ctx, "srv", srv)
			// send to worker's channel
			taskChannel <- Task{flow}
		} else {
			// something wrong received
			srv.errCh <- ErrWrongFlow
			// release worker
			// TODO https://github.com/boomfunc/base/issues/19
		}
	}
}

// this function listen all server channels
// TODO listen OS signals and gracefully close server
// check for errors additionally -> errors to log
// response.Stat to log
func (srv *Server) listen() {
	for {
		select {
		case err := <-srv.errCh:
			if err != nil {
				tools.ErrorLog(err)
			}

		case flow := <-srv.outputCh:
			// ready response from dispatcher system
			// log ANY kind of result
			AccessLog(flow)
			// and errors
			if err := flow.Stat.Error; err != nil {
				// TODO think about it
				// TODO https://play.golang.org/p/dNV2qI90EKQ
				go func() {
					srv.errCh <- flow.Stat.Error
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
	go srv.listen()

	// GOROUTINE 3 (main engine)
	// bridge between worker(from dispatcher) and task(from heap)
	go srv.engine()

	// Here we can test some of our system requirements and performance recommendations
	PerformanceLog(srv.dispatcher.MaxWorkers)

	// GOROUTINE 1 (main) - this goroutine
	// This is thread blocking procedure - infinity loop
	srv.transport.Serve()
}
