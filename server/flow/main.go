package flow

import (
	"context"
	"io"

	srvctx "github.com/boomfunc/base/server/context"
	"github.com/boomfunc/base/tools/chronometer"
	"github.com/google/uuid"
)

type Data struct {
	UUID        uuid.UUID
	RWC         io.ReadWriteCloser
	Ctx         context.Context
	Chronometer *chronometer.Chronometer
	Stat        Stat
}

func New(rwc io.ReadWriteCloser) *Data {
	return &Data{
		UUID:        uuid.New(),
		RWC:         rwc,
		Ctx:         srvctx.New(),
		Chronometer: chronometer.New(),
	}
}
