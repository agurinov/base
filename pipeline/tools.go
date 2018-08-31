package pipeline

import (
	"context"
	"io"

	"github.com/boomfunc/base/tools/executor"
)

// piping establishes pipe connections between IO processes (Able)
// the first obj accepts as stdin the input buffer
// the last obj puts into output buffer his stdout
func piping(input io.ReadCloser, output io.WriteCloser, objs ...Able) error {
	// main logic that create pairs of (io.ReadCloser, io.WriteCloser)
	// but with offset to another obj
	// for example
	// obj 1: (input, io.WriteCloser 1)
	// obj 2: (io.ReadCloser 1, io.WriteCloser 2)
	// obj 3: (io.ReadCloser 2, io.WriteCloser 3)
	// obj 4: (io.ReadCloser 3, output)
	for i := 0; i < len(objs); i++ {
		if i == 0 {
			// case this obj first
			objs[i].setStdin(input)
		}
		if i == len(objs)-1 {
			// case this obj last
			objs[i].setStdout(output)
		} else {
			// this is intermediate obj, need piping
			r, w := io.Pipe()
			objs[i].setStdout(w)
			objs[i+1].setStdin(r)
		}
	}

	return nil
}

func executeFunc(obj Exec) executor.OperationFunc {
	return func(ctx context.Context) error {
		if err := obj.run(ctx); err != nil {
			obj.close(ctx)
			return err
		}
		return obj.close(ctx)
	}
}

func prepareFunc(obj Exec) executor.OperationFunc {
	return func(ctx context.Context) error {
		if err := obj.prepare(ctx); err != nil {
			return err
		}
		return obj.check(ctx)
	}
}

func run(ctx context.Context, objs ...Exec) error {
	// Phase 1. PREPARE AND CHECK
	// in case of error it will be rolled back to initial incoming state
	prepareUp := make([]executor.OperationFunc, len(objs))
	prepareDown := make([]executor.OperationFunc, len(objs))
	executeUp := make([]executor.OperationFunc, len(objs))

	for i, obj := range objs {
		prepareUp[i] = prepareFunc(obj)
		executeUp[i] = executeFunc(obj)
		prepareDown[i] = obj.close
	}

	return executor.New(
		ctx,
		executor.Operation(prepareUp, prepareDown, false),
		executor.Operation(executeUp, nil, true),
	).Run()
}
