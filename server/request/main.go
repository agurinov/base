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

type Stat struct {
	Duration time.Duration
	Request  Interface
	Error    error
	Len      int64
}

func (stat Stat) Successful() bool {
	return stat.Error == nil
}

// type Response interface {
// 	Duration() time.Duration
// 	Request() Request
// 	Status() bool
// 	Len() int64
// }
