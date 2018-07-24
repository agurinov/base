package application

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	// "strings"
	"net/http"
	"time"

	"github.com/boomfunc/base/conf"
	"github.com/boomfunc/base/server/request"
	"github.com/google/uuid"
)

type Interface interface {
	// Parse(io.ReadWriter) (request.Interface, error)
	Handle(io.ReadWriter) request.Stat
}

// type Extension interface {
// 	FromReader(io.Reader) io.Reader
// }

type Extension func(io.Reader) (io.WriterTo, error)

type Application struct {
	router *conf.Router
	// extensions []Extension
}

func (app *Application) parse(rw io.ReadWriter) (request.Interface, error) {
	var r JSONRequest
	r.uuid = uuid.New()

	intermediate := struct {
		Url   string
		Input string
	}{}

	decoder := json.NewDecoder(rw)
	if err := decoder.Decode(&intermediate); err != nil {
		return nil, err
	}

	r.url = intermediate.Url
	r.input = intermediate.Input

	return &r, nil
}

func (app *Application) Handle(rw io.ReadWriter) (stat request.Stat) {
	var req request.Interface
	var begin time.Time
	var err error
	var written int64

	defer func() {
		// end measuring and collect data
		stat.Duration = time.Since(begin)
		stat.Request = req
		stat.Error = err
		stat.Len = written
	}()

	// Start measuring
	begin = time.Now()

	// Parse request
	req, err = app.parse(rw)
	if err != nil {
		return
	}

	// Resolve view
	route, err := app.router.Match(req.Url())
	if err != nil {
		return
	}

	// get base (initial) writerTo interface
	// Pipeline output
	var wt io.WriterTo
	var buf bytes.Buffer

	err = route.Run(req.Input(), &buf)
	if err != nil {
		return
	}

	// Apply extensions
	wt, err = HTTPExtension(&buf)
	if err != nil {
		return
	}
	// for i, ext := range app.extensions {
	// for i, ext := range []Extension{JSONExtension} {
	// 	// wt, err = ext.FromReader(r)
	// 	wt, err = ext(wt)
	// 	if err != nil {
	// 		return stat
	// 	}
	// }

	// write data to rwc only if all success
	written, err = wt.WriteTo(rw)

	return
}

func HTTPExtension(reader io.Reader) (io.WriterTo, error) {
	var b bytes.Buffer

	response := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		// Proto:         "HTTP/1.1",
		Body: ioutil.NopCloser(reader),
		// ContentLength: int64(len(body)),
		// Request:       req,
	}

	response.Write(&b)

	return &b, nil
}

// var r io.Reader = req.Input()
//
// for i, ext := range app.extensions {
// 	r = ext.FromReader(r)
// }
//
// return r
