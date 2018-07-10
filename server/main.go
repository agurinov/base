package server

import (
	"net"
)

type Wrapper interface {
	Serve()
	Addr() net.Addr

	// log
	startupLog()
	accessLog()
	errorLog()
}
