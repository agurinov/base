package server

import (
	"io"
)

type Request struct {
	server *Server
	rw     io.ReadWriter
}

func NewRequest(server *Server, rw io.ReadWriter) Request {
	return Request{server, rw}
}
