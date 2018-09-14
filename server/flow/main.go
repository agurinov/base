package flow

import (
	"context"
	"io"
	// "net/url"

	// "github.com/google/uuid"
	srvctx "github.com/boomfunc/base/server/context"
	"github.com/boomfunc/base/tools"
)

// type Request struct {
// 	UUID  uuid.UUID
// 	Url   *url.URL
// 	Input io.Reader
// }

type Data struct {
	Input  io.ReadWriteCloser
	Ctx    context.Context
	Timing *tools.Timing
}

func New(input io.ReadWriteCloser) *Data {
	return &Data{
		Input:  input,
		Ctx:    srvctx.New(),
		Timing: tools.NewTiming(),
	}
}
