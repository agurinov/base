package server

import (
	"context"
	"errors"
	"io"

	srvctx "github.com/boomfunc/base/server/context"
	"github.com/boomfunc/base/tools"
)

type Task struct {
	ctx   context.Context
	input io.ReadWriteCloser
}

// Solve implements dispatcher.Task interface
// this function will be passed to dispatcher system
// and will be run at parallel
func (task Task) Solve() {
	srvInterface, err := srvctx.GetMeta(task.ctx, "srv")
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

	defer task.input.Close()

	srv.outputCh <- srv.app.Handle(task.ctx, task.input)
}
