package request

import (
	"io"
	"net/url"

	"github.com/google/uuid"
)

type Request struct {
	UUID  uuid.UUID
	Url   *url.URL
	Input io.Reader
}

func New(raw string, input io.Reader) (*Request, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	request := &Request{
		UUID:  uuid.New(),
		Url:   u,
		Input: input,
	}

	return request, nil
}
