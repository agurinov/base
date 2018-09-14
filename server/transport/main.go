package transport

import (
	"github.com/boomfunc/base/server/flow"
)

type Interface interface {
	Connect(chan *flow.Data, chan error)
	Serve()
}
