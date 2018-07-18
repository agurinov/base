package application

import (
	"github.com/boomfunc/base/server/request"
)

type Interface interface {
	HandleRequest(request request.Request) request.Response
}
