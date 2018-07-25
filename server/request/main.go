package request

import (
	"io"
	"time"

	"github.com/google/uuid"
)

type Request struct {
	UUID  uuid.UUID
	Url   string
	Input io.Reader
}

func New(url string, input io.Reader) *Request {
	return &Request{
		UUID:  uuid.New(),
		Url:   url,
		Input: input,
	}
}

type Stat struct {
	Duration time.Duration
	Request  *Request
	Error    error
	Len      int64
}

func (stat Stat) Successful() bool {
	return stat.Error == nil
}
