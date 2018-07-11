package server

import (
	"net"
)

type Wrapper interface {
	Serve()
	Addr() net.Addr
	ServeConn(conn net.Conn)

	// log
	startupLog()
	accessLog()
	errorLog()
}
