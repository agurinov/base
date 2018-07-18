package application

import (
	"net"

	"github.com/boomfunc/base/server/request"
)

type Interface interface {
	HandleRequest(request request.Request, conn net.Conn) request.Response
}
