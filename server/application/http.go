package application

import (
	"bufio"
	"io"
	"net/http"

	"github.com/boomfunc/base/server/request"
)

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
		Proto:      "HTTP/1.1",
		StatusCode: 200,
		Body:       rc,
		Request:    packer.request,
	}

	return 0, response.Write(w)
}
