package server

import (
	"io"

	"github.com/boomfunc/base/conf"
)

// TODO (written int64, err error) at return
func handleRequest(req Request, router *conf.Router, output io.Writer) (err error) {
	// logging and error handling block
	// this defer must be invoked last (first in) for recovering all available panics and errors
	defer func() {
		var status = "SUCCESS"

		if err != nil {
			errorLog(err)
			status = "ERROR"
		}
		// log ANY kind result
		accessLog(req, status)
	}()

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
