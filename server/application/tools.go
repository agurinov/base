package application

import (
	"time"

	"github.com/boomfunc/base/conf"
	"github.com/boomfunc/base/server/request"
)

func New(router *conf.Router) Interface {
	return &JsonApplicationLayer{router}
}

func handle(req request.Interface, router *conf.Router) (response request.Response) {
	var begin time.Time
	var err error
	var written int64

	defer func() {
		// end measuring and collect data
		response.Duration = time.Since(begin)
		response.Request = req
		response.Error = err
		response.Len = written
	}()

	// Start measuring
	begin = time.Now()

	url := req.Url()
	input := req.Input()
	output := req.Output()

	// Phase 1. Resolve view
	route, err := router.Match(url)
	if err != nil {
		return
	}

	// Phase 2. Write answer to output
	err = route.Run(input, output)

	return
}
