package transport

import (
	"bytes"
	"time"

	"github.com/boomfunc/base/server/flow"
	"github.com/boomfunc/base/tools"
)

type tick struct {
	ticker *time.Ticker

	// server integration
	inputCh chan *flow.Data
	errCh   chan error
}

func (tr *tick) Connect(inputCh chan *flow.Data, errCh chan error) {
	tr.inputCh = inputCh
	tr.errCh = errCh
}

func (tr *tick) Serve() {
	for {
		<-tr.ticker.C

		flow := flow.New(
			tools.ReadWriteCloser(bytes.NewBufferString("GET /ping HTTP/1.0\r\n\r\n")),
		)
		tr.inputCh <- flow
	}
}
