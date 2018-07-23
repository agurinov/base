package request

import (
	"io"
	"time"

	"github.com/google/uuid"
)

type Interface interface {
	UUID() uuid.UUID
	Url() string
	Input() io.Reader
	Output() io.Writer
}

type Response struct {
	Duration time.Duration
	Request  Interface
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
