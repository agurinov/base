package server

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/boomfunc/base/conf"
)

type Wrapper interface {
	Serve()
	// Router()
	// Close()

	// Addr() net.Addr
	// ServeConn(conn net.Conn)

	// log
	// startupLog()
	// accessLog()
	// errorLog()
}

func New(mode string, ip net.IP, port int, filename string) (wrapper Wrapper, err error) {
	// startup logging if no error
	defer func() {
		if err == nil {
			serverStartupLog(strings.ToUpper(mode), fmt.Sprintf("%s:%d", ip, port), filename)
		}
	}()

	// Phase 1. get config for server routing
	router, err := conf.LoadFile(filename)
	if err != nil {
		// cannot load server config
		return nil, err
	}

	// Phase 2. Resolve tcp address and create tcp server listening on provided port
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		// cannot resolve address (invalid options (ip or port))
		return nil, err
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		// cannot establish connection on this addr
		return nil, err
	}

	// Phase 3. Create transport (http or rpc)
	switch mode {
	case "http":
		return newHTTPWrapper(tcpListener, router), nil
	case "rpc":
		return newRPCWrapper(tcpListener, router)
	default:
		return nil, errors.New("server: Unknown server protocol")
	}
}
