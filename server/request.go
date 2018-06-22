package server

import (
	"bufio"
	"io"
	"strings"
	"io/ioutil"

	"github.com/boomfunc/log"
)

type request struct {
	url    string
	body    string
	// reader io.Reader
}

func NewRequest(r io.Reader) (*request, error) {
	reader := bufio.NewReader(r)

	// url
	url, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(reader)

	log.Debug(string(body))

	req := request{
		url:    string(url),
		body: string(body),
		// reader: reader,
	}

	return &req, nil
}

func (r *request) Url() string {
	return r.url
}

func (r *request) Reader() io.Reader {
	return strings.NewReader(r.body)
	// return r.reader
}
