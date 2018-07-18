package application

import (
	"io"
	"time"

	"github.com/boomfunc/base/conf"
	"github.com/boomfunc/base/server/request"
)

// TODO (written int64, err error) at return
func handleRequest(req request.Request, router *conf.Router, output io.Writer) request.Response {
	var response request.Response
	var begin time.Time
	var err error

	defer func() {
		response.Duration = time.Since(begin)
		response.Request = req
		response.Status = "SUCCESS"

		if err != nil {
			// ch <- err
			response.Status = "ERROR"
		}
	}()

	// Start measuring
	begin = time.Now()

	// Phase 1. Resolve view
	route, err := router.Match(req.Url())
	if err != nil {
		return response
	}

	// Phase 2. Write answer to output
	err = route.Run(req.Body(), output)

	return response
}
