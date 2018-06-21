package server

import (
	"bufio"
	"io"

	"github.com/boomfunc/log"
)

type request struct {
	url *string

	reader io.Reader
}

func NewRequest(reader io.Reader) *request {
	return &request{reader: reader}
}

func (r *request) Url() (string, error) {
	if r.url == nil {
		// not parsed -> get first line
		scanner := bufio.NewScanner(r.reader)
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			return "", err
		}

		url := scanner.Text()
		log.Debug(url)
		r.url = &url
	}

	return *r.url, nil
}

func (r *request) Reader() io.Reader {
	return r.reader
}
