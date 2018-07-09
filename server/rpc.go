package server

import (
	"fmt"
	"net"
	"io"
	"net/rpc"
	"net/rpc/jsonrpc"
	"net/http"
	"errors"

	"github.com/boomfunc/base/conf"
	"github.com/boomfunc/log"
)




type Args struct {
    A, B int
}

type Quotient struct {
    Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
    *reply = args.A * args.B
    return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
    if args.B == 0 {
        return errors.New("divide by zero")
    }
    quo.Quo = args.A / args.B
    quo.Rem = args.A % args.B
    return nil
}





type RPCServerWrapper struct {
	listener *net.TCPListener
	router   *conf.Router
	rpc   *rpc.Server
}

func New(ip net.IP, port int, filename string) (*RPCServerWrapper, error) {
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
	arith := new(Arith)
	err = server.Register(arith)
	if err != nil {
		// cannot register methods to rpc server
		return nil, err
	}

	// All wright, some basic logging
	log.Infof("TCP server up and running on %s", log.Wrap(fmt.Sprintf("%s", tcpAddr), log.Bold, log.Blink))
	log.Infof("Spawned config file: %s", log.Wrap(filename, log.Bold))
	log.Debugf("Enabled %s mode", log.Wrap("DEBUG", log.Bold, log.Blink))

	wrapper := &RPCServerWrapper{
		listener: tcpListener,
		router: router,
		rpc: server,
	}
	return wrapper, nil
}

func (wrp *RPCServerWrapper) ServeTCP() {
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
	// TODO defer ch.Close()
	// TODO defer s.conn.Close()
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
	// https://rcrowley.org/articles/golang-graceful-stop.html
	for {
		// Listen for an incoming connection.
		conn, err := wrp.listener.Accept()
		if err != nil {
			panic(err)
		}

		go wrp.rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

func (wrp *RPCServerWrapper) ServeHTTP() {

	http.HandleFunc("/rpc", func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		res := NewRPCRequest(req.Body).Call()
		io.Copy(w, res)
	})

	http.Serve(wrp.listener, nil)
}
