package server

import (
	"errors"

	"github.com/boomfunc/base/server/context"
	"github.com/boomfunc/base/server/flow"
	"github.com/boomfunc/base/tools"
)

type Task struct {
	flow *flow.Data
}

// Solve implements dispatcher.Task interface
// this function will be passed to dispatcher system
// and will be run at parallel
func (task Task) Solve() {
	srvInterface, err := context.GetMeta(task.flow.Ctx, "srv")
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

	defer task.flow.RWC.Close()

	task.flow.Chronometer.Enter("app")
	srv.app.Handle(task.flow)
	task.flow.Chronometer.Exit("app")

	srv.outputCh <- task.flow
}
