package request

import (
	"io"
	"net"
	"time"

	"github.com/google/uuid"
)

type Request interface {
	UUID() uuid.UUID
	Url() string
	Body() io.Reader
	Conn() net.Conn
}

type Response struct {
	Duration time.Duration
	Request  Request
	Error    error
	Len      int64
}

func (r Response) Successful() bool {
	return r.Error == nil
}

// type Response interface {
// 	Duration() time.Duration
// 	Request() Request
// 	Status() bool
// 	Len() int64
// }
