package request

import (
	"io"
	"net"
	"strings"

	"github.com/google/uuid"
)

type r struct {
	conn net.Conn
	uuid uuid.UUID
}

func (req *r) UUID() uuid.UUID {
	return req.uuid
}

func (req *r) Url() string {
	return "geo"
}

func (req *r) Body() io.Reader {
	return strings.NewReader("185.86.151.11")
}

func (req *r) Conn() net.Conn {
	return req.conn
}

func New(conn net.Conn) Request {
	return &r{
		conn: conn,
		uuid: uuid.New(),
	}
}

// func (req *Request) Handle() (err error) {
// 	// logging and error handling block
// 	// this defer must be invoked last (first in) for recovering all available panics and errors
// 	defer func() {
// 		var status = "SUCCESS"
//
// 		if err != nil {
// 			ErrorLog(err)
// 			status = "ERROR"
// 		}
// 		// log ANY kind result
// 		AccessLog(req, status)
// 	}()
//
// 	// Firstly - close connection
// 	defer func() {
// 		err = req.conn.Close()
// 	}()
//
// 	// Phase 1. Resolve view
// 	route, err := req.router.Match(req.Url())
// 	if err != nil {
// 		return err
// 	}
//
// 	var output bytes.Buffer
//
// 	// Phase 2. Write answer to output
// 	if err = route.Run(req.Body(), &output); err != nil {
// 		return err
// 	}
//
// 	io.Copy(req.conn, &output)
//
// 	return nil
// }

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
