package server

import (
	"net"
)

type Request struct {
	conn   net.Conn
	server *Server
}
