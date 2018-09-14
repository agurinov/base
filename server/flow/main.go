package flow

import (
	"context"
	"io"
	// "net/url"

	// "github.com/google/uuid"
	srvctx "github.com/boomfunc/base/server/context"
	"github.com/boomfunc/base/server/request"
	"github.com/boomfunc/base/tools"
)

type Stat struct {
	Request *request.Request
	Error   error
	Len     int64
}

func (stat Stat) Successful() bool {
	return stat.Error == nil
}

type Data struct {
	Input  io.ReadWriteCloser
	Ctx    context.Context
	Timing *tools.Timing
	Stat   Stat
}

func New(input io.ReadWriteCloser) *Data {
	return &Data{
		Input:  input,
		Ctx:    srvctx.New(),
		Timing: tools.NewTiming(),
	}
}
