package server

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/boomfunc/base/conf"
)

type TCPServerWrapper struct {
	listener *net.TCPListener
	router   *conf.Router
	server   *rpc.Server
}

func NewTCP(ip net.IP, port int, filename string) (*TCPServerWrapper, error) {
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

	// Phase 3. Create RPC server
	server := rpc.NewServer()
	pipeline := &pipelineRPC{router}
	if err := server.RegisterName("Pipeline", pipeline); err != nil {
		// cannot register methods to rpc server
		return nil, err
	}

	startupLog("TCP", tcpAddr.String(), filename)

	wrapper := &TCPServerWrapper{
		listener: tcpListener,
		router:   router,
		server:   server,
	}

	return wrapper, nil
}

func (wrp *TCPServerWrapper) Serve() {
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
	// TODO defer ch.Close()
	// TODO defer s.conn.Close()
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
	// https://rcrowley.org/articles/golang-graceful-stop.html

	// Phase 1. Listen infinitely TCP connection for incoming requests
	for {
		// Listen for an incoming connection.
		conn, err := wrp.listener.Accept()
		if err != nil {
			panic(err)
		}
		// Handle connections in a new goroutine.
		// Phase 3. Message received, resolve this shit concurrently!
		// conn will be closed by ServeCodec!
		go wrp.server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
