package server

import (
	"fmt"
	"net"
	"runtime/debug"

	"github.com/boomfunc/log"

	"app/conf"
)

type TCPServer struct {
	listener *net.TCPListener
	router   *conf.Router
}

func NewTCP(ip net.IP, port int, filename string) (*TCPServer, error) {
	// Phase 1. get config for server routing
	router, err := conf.LoadFile(filename)
	// cannot load server config
	if err != nil {
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

	log.Infof("TCP server up and running on %s", log.Wrap(fmt.Sprintf("%s", tcpAddr), log.Bold, log.Blink))
	log.Infof("Spawned config file: %s", log.Wrap(filename, log.Bold))
	log.Debugf("Enabled %s mode", log.Wrap("DEBUG", log.Bold, log.Blink))

	server := &TCPServer{listener: tcpListener, router: router}
	return server, nil
}

func (s *TCPServer) handle(conn net.Conn) {
	var url string
	var status = "SUCCESS"

	// logging and error handling block
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("%s\n%s", err, debug.Stack())
			status = "ERROR"
			conn.Close()
		}
		// log ANY kind result
		log.Infof("%s\t-\t%s", url, status)
	}()

	// TODO some layer -> separate uri from request body
	// TODO need some internal style of requests
	// TODO separate uri and request body
	// TODO timeoutDuration := 5 * time.Second
	request, err := NewRequest(conn)
	if err != nil {
		panic(err)
	}

	url = request.Url()

	route, err := s.router.Match(url)
	if err != nil {
		panic(err)
	}

	// important!
	// input and output is io.ReadCloser and io.WriteCloser
	// after route.Run completion they will be closed
	if err := route.Run(request.Reader(), conn); err != nil {
		panic(err)
	}
}

func (s *TCPServer) Serve() {
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe
	// TODO defer ch.Close()
	// TODO defer s.conn.Close()
	// TODO unreachable https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe

	// Phase 1. Listen infinitely TCP connection for incoming requests
	for {
		// Listen for an incoming connection.
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}
		// Handle connections in a new goroutine.
		// Phase 3. Message received, resolve this shit concurrently!
		go s.handle(conn)
	}
}
