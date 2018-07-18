package request

import (
	"io"
	"time"

	"github.com/google/uuid"
)

type Request interface {
	UUID() uuid.UUID
	Url() string
	Body() io.Reader
}

type Response interface {
	Duration() time.Duration
	Request() Request
	Status() bool
	Len() int64
}
