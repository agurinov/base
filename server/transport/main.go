package transport

import (
	"net"
)

type Interface interface {
	Serve()
	Conn(conn net.Conn)
	Error(err error)
}
