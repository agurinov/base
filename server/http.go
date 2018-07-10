package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/boomfunc/base/conf"
)

type HTTPServerWrapper struct {
	listener *net.TCPListener
	router   *conf.Router
}

func NewHTTP(ip net.IP, port int, filename string) (*HTTPServerWrapper, error) {
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

	startupLog("HTTP", tcpAddr.String(), filename)

	wrapper := &HTTPServerWrapper{
		listener: tcpListener,
		router:   router,
	}

	return wrapper, nil
}

func (wrp *HTTPServerWrapper) Serve() {

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		request, err := NewRequest(req)

		if err != nil {
			// handle error 500
		}

		handleRequest(request, wrp.router, w)
	})

	http.Serve(wrp.listener, nil)
}
