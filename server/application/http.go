package application

import (
	"bufio"
	"context"
	"io"
	"net/http"

	srvctx "github.com/boomfunc/base/server/context"
	"github.com/boomfunc/base/server/request"
	"github.com/boomfunc/base/tools"
)

// Load test
// seq 1000 | xargs -n 1 -P 250 sh -c "curl -i http://playground.lo:8080/geo?ip=185.86.151.11"
type httpPacker struct {
	request *http.Request
}

func (packer *httpPacker) Unpack(ctx context.Context, r io.Reader) (*request.Request, error) {
	br := bufio.NewReader(r)
	httpRequest, err := http.ReadRequest(br)
	if err != nil {
		return nil, err
	}

	// extend ctx
	// get remote ip and save to context
	srvctx.SetMeta(
		ctx, "ip",
		tools.GetRemoteIP(httpRequest.Header.Get("X-Forwarded-For"), nil),
	)
	// Get http query and save to context
	values, err := srvctx.Values(ctx)
	if err != nil {
		return nil, err
	}
	values.Q = httpRequest.URL.Query()

	packer.request = httpRequest

	return request.New(
		httpRequest.URL.RequestURI(),
		httpRequest.Body,
	)
}

func (packer *httpPacker) Pack(r io.Reader, w io.Writer) (int64, error) {
	response := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Body:       tools.ReadCloser(r),
		Request:    packer.request,
	}

	return 0, response.Write(w)
}
