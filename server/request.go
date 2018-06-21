package server

import (
	"bufio"
	"io"
)

type request struct {
	url    string
	reader io.Reader
}

func NewRequest(r io.Reader) (*request, error) {
	reader := bufio.NewReader(r)

	// url
	url, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}

	req := request{
		url:    string(url),
		reader: reader,
	}

	return &req, nil
}

func (r *request) Url() string {
	return r.url
}

func (r *request) Reader() io.Reader {
	return r.reader
}
