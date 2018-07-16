package server

import (
	// "bytes"
	"io"
	// "net/http"
	"strings"
	"net"

	"github.com/google/uuid"
	"github.com/boomfunc/base/conf"
)

type Request struct {
	conn net.Conn  // The operation to perform.
	router *conf.Router
	c  chan io.Reader    // The channel to return the result.
	uuid uuid.UUID
}

func (req *Request) UUID() uuid.UUID {
	return req.uuid
}

func (req *Request) Url() string {
	return "ping"
}

func (req *Request) Body() io.Reader {
	return strings.NewReader("")
}

func (req *Request) Handle() (err error) {
	// logging and error handling block
	// this defer must be invoked last (first in) for recovering all available panics and errors
	defer func() {
		var status = "SUCCESS"

		if err != nil {
			ErrorLog(err)
			status = "ERROR"
		}
		// log ANY kind result
		AccessLog(req, status)
	}()

	// Phase 1. Resolve view
	route, err := req.router.Match(req.Url())
	if err != nil {
		return err
	}

	// Phase 2. Write answer to output
	return route.Run(req.Body(), req.conn)
}


// type Request interface {
// 	UUID() uuid.UUID
// 	Url() string
// 	Body() io.Reader
//
// 	// Serve()
// }
//
func NewRequest(conn net.Conn, router *conf.Router) Request {
	return Request{
		conn: conn,
		router: router,
		c: make(chan io.Reader),
		uuid: uuid.New(),
	}
}

// func NewRequest2(req interface{}) (*RPCRequest, error) {
// 	switch typed := req.(type) {
// 	case *RPCArgs:
// 		return &RPCRequest{uuid.New(), typed}, nil
// 	// case *http.Request:
// 		// return &HTTPRequest{uuid.New(), typed}, nil
// 	default:
// 		return nil, errors.New("server: Unknown underlying request type")
// 	}
// }

// RPCRequest is wrapper type for rpc args to be Request interface
// type RPCRequest struct {
// 	uuid uuid.UUID
// 	*RPCArgs
// }
//
// func (r *RPCRequest) Url() string {
// 	return r.RPCArgs.Url
// }
//
// func (r *RPCRequest) Body() io.Reader {
// 	return bytes.NewReader(r.RPCArgs.Body)
// }
//
// func (r *RPCRequest) UUID() uuid.UUID {
// 	return r.uuid
// }
//
// // HTTPRequest is wrapper type for http to be Request interface
// type HTTPRequest struct {
// 	uuid uuid.UUID
// 	*http.Request
// }
//
// func (r *HTTPRequest) Url() string {
// 	return r.Request.URL.RequestURI()
// }
//
// func (r *HTTPRequest) Body() io.Reader {
// 	// TODO TMP HARDCODE FOR geoservice
// 	if r.Request.URL.Path == "/geo" || r.Request.URL.Path == "/geo/" {
// 		return strings.NewReader(r.Request.URL.Query().Get("ip"))
// 	}
//
// 	return r.Request.Body
// }
//
// func (r *HTTPRequest) UUID() uuid.UUID {
// 	return r.uuid
// }
