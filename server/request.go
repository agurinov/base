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

	// url
	url, _, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}

	req := request{
		uuid:   uuid.New(),
		url:    string(url),
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
