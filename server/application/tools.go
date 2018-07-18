package application

import (
	"io"
	"net"
	"time"

	"github.com/boomfunc/base/conf"
	"github.com/boomfunc/base/server/request"
)

type BaseApplicationLayer struct {
	router *conf.Router
}

func (app *BaseApplicationLayer) HandleRequest(request request.Request, conn net.Conn) request.Response {
	return handleRequest(request, app.router, conn)
}

func New(router *conf.Router) Interface {
	return &BaseApplicationLayer{router}
}

// TODO (written int64, err error) at return
func handleRequest(req request.Request, router *conf.Router, output io.Writer) (response request.Response) {
	var begin time.Time
	var err error

	defer func() {
		// end measuring
		response.Duration = time.Since(begin)
		response.Request = req
		response.Error = err
	}()

	// Start measuring
	begin = time.Now()

	// Phase 1. Resolve view
	route, err := router.Match(req.Url())
	if err != nil {
		return
	}

	// Phase 2. Write answer to output
	err = route.Run(req.Body(), output)

	return
}
