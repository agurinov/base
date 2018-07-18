package application

import (
	"io"

	"github.com/boomfunc/base/conf"
	"github.com/boomfunc/base/server/request"
)

// TODO (written int64, err error) at return
func handleRequest(req request.Request, router *conf.Router, output io.Writer) error {
	// Phase 1. Resolve view
	route, err := router.Match(req.Url())
	if err != nil {
		return err
	}

	// Phase 2. Write answer to output
	if err := route.Run(req.Body(), output); err != nil {
		return err
	}

	return nil
}
