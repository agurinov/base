package executor

import (
	"context"
)

// https://github.com/rafaeldias/async/blob/master/async.go

type executor struct {
	fns []func(context.Context) error
	ctx context.Context
}

// func New(fns ...func(context.Context) error) *executor {
// 	return &executor{
// 		fns: fns,
// 	}
// }

func NewWithContext(ctx context.Context, fns ...func(context.Context) error) *executor {
	return &executor{
		fns: fns,
		ctx: ctx,
	}
}

func (exctr *executor) Add(fn func(context.Context) error) {
	exctr.fns = append(exctr.fns, fn)
}

func (exctr *executor) Concurrent() error {
	return concurrent(exctr.ctx, exctr.fns...)
}
