package server

import (
	"bufio"
	"io"

	"github.com/google/uuid"
)

type request struct {
	uuid   uuid.UUID
	url    string
	reader io.Reader
}

func NewRequest(r io.Reader) (*request, error) {
	reader := bufio.NewReader(r)

	var url string

	// url
	buf, _, err := reader.ReadLine()
	if err != nil {
		if err != io.EOF {
			return nil, err
		} else {
			url = "ping"
		}
	} else {
		url = string(buf)
	}

	req := request{
		uuid:   uuid.New(),
		url:    url,
		reader: reader,
	}

	return &req, nil
}

func (r *request) Url() string {
	return r.url
}

func (r *request) Body() io.Reader {
	return r.reader
}

func (r *request) Id() uuid.UUID {
	return r.uuid
}
