package application

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/boomfunc/base/conf"
	"github.com/boomfunc/base/server/request"
	"github.com/google/uuid"
)

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

func (app *JsonApplicationLayer) Handle(request request.Interface) request.Response {
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
