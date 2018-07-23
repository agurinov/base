package application

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/boomfunc/base/conf"
	"github.com/boomfunc/base/server/request"
	"github.com/google/uuid"
)

// Load test
// JS='{"url":"geo","input":"185.86.151.11"}'
// seq 1000 | xargs -n 1 -P 250 sh -c "echo '$JS' | nc playground.lo 8080"

// simple JSON {"url":"", "input":"..."} incoming
type JsonApplicationLayer struct {
	router *conf.Router
}

func (app *JsonApplicationLayer) Parse(request io.ReadWriter) (request.Interface, error) {
	var r JSONRequest
	r.uuid = uuid.New()
	r.rw = request

	intermediate := struct {
		Url   string
		Input string
	}{}

	decoder := json.NewDecoder(request)
	if err := decoder.Decode(&intermediate); err != nil {
		return nil, err
	}

	r.url = intermediate.Url
	r.input = intermediate.Input

	return &r, nil
}

func (app *JsonApplicationLayer) Handle(request request.Interface) request.Stat {
	return handle(request, app.router)
}

type JSONRequest struct {
	uuid  uuid.UUID
	url   string
	input string
	rw    io.ReadWriter
}

func (req *JSONRequest) UUID() uuid.UUID {
	return req.uuid
}

func (req *JSONRequest) Url() string {
	return req.url
}

func (req *JSONRequest) Input() io.Reader {
	return strings.NewReader(req.input)
}

func (req *JSONRequest) Output() io.Writer {
	return req.rw
}
