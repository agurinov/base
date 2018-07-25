package application

import (
	"bufio"
	"io"
	"net/http"

	"github.com/boomfunc/base/server/request"
)

// Load test
// seq 1000 | xargs -n 1 -P 250 sh -c "curl -i http://playground.lo:8080/geo?ip=185.86.151.11"
type httpPacker struct {
	request *http.Request
}

func (packer *httpPacker) Unpack(r io.Reader) (*request.Request, error) {
	br := bufio.NewReader(r)
	httpRequest, err := http.ReadRequest(br)

	if err != nil {
		return nil, err
	}

	packer.request = httpRequest

	req := request.New(
		httpRequest.URL.RequestURI(),
		httpRequest.Body,
	)

	return req, nil
}

func (packer *httpPacker) Pack(rc io.ReadCloser, w io.Writer) (int64, error) {
	response := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Body:       rc,
		Request:    packer.request,
	}

	return 0, response.Write(w)
}
