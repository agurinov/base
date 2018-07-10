package server

import (
	"github.com/google/uuid"
)

type Request interface {
	UUID() uuid.UUID
	Url() string
	Body() []byte
}

// ArgsRequest is wrapper type for rpc args to be Request interface
type ArgsRequest struct {
	uuid uuid.UUID
	*Args
}

func (a *ArgsRequest) Url() string {
	return a.Args.Url
}

func (a *ArgsRequest) Body() []byte {
	return a.Args.Body
}

func (a *ArgsRequest) UUID() uuid.UUID {
	return a.uuid
}
