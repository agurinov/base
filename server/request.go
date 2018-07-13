package server

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Request interface {
	UUID() uuid.UUID
	Url() string
	Body() io.Reader
}

func NewRequest(req interface{}) (Request, error) {
	switch typed := req.(type) {
	case *Args:
		return &ArgsRequest{uuid.New(), typed}, nil
	case *http.Request:
		return &HTTPRequest{uuid.New(), typed}, nil
	default:
		return nil, errors.New("server: Unknown underlying request type")
	}
}

// ArgsRequest is wrapper type for rpc args to be Request interface
type ArgsRequest struct {
	uuid uuid.UUID
	*Args
}

func (r *ArgsRequest) Url() string {
	return r.Args.Url
}

func (r *ArgsRequest) Body() io.Reader {
	return bytes.NewReader(r.Args.Body)
}

func (r *ArgsRequest) UUID() uuid.UUID {
	return r.uuid
}

// HTTPRequest is wrapper type for http to be Request interface
type HTTPRequest struct {
	uuid uuid.UUID
	*http.Request
}

func (r *HTTPRequest) Url() string {
	return r.Request.URL.RequestURI()
}

func (r *HTTPRequest) Body() io.Reader {
	// TODO TMP HARDCODE FOR geoservice
	if r.Request.URL.Path == "/geo" || r.Request.URL.Path == "/geo/" {
		return strings.NewReader(r.Request.URL.Query().Get("ip"))
	}

	return r.Request.Body
}

func (r *HTTPRequest) UUID() uuid.UUID {
	return r.uuid
}
