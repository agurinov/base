package server

import (
	"github.com/boomfunc/base/server/request"
)

type Request struct {
	under  request.Request
	server *Server
}
