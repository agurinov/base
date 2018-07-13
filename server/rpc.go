package server

import (
	"bytes"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/boomfunc/base/conf"
)

type RPCWrapper struct {
	listener net.Listener
	router   *conf.Router
	server   *rpc.Server
}

func newRPCWrapper(listener net.Listener, router *conf.Router) (*RPCWrapper, error) {
	server := rpc.NewServer()
	pipeline := &pipelineRPC{router}
	if err := server.RegisterName("Pipeline", pipeline); err != nil {
		// cannot register methods to rpc server
		return nil, err
	}

	wrapper := &RPCWrapper{
		listener: listener,
		router:   router,
		server:   server,
	}
	return wrapper, nil
}

// echo '{"method":"Pipeline.Run","params":[{"Url":"ping","Body":"pong"}],"id":123}' | nc playground.lo 8080
func (wrp *RPCWrapper) Serve() {
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

type RPCArgs struct {
	Url  string
	Body []byte
}

type pipelineRPC struct {
	router *conf.Router
}

func (rpc *pipelineRPC) Run(args *RPCArgs, reply *[]byte) error {
	var output bytes.Buffer

	request, err := NewRequest(args)

	if err != nil {
		return err
	}

	if err := handleRequest(request, rpc.router, &output); err != nil {
		return err
	}

	// write answer back
	*reply = output.Bytes()

	return nil
}
