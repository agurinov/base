package server

import (
	"bytes"

	"github.com/boomfunc/base/conf"
	"github.com/google/uuid"
)

type Args struct {
	Url  string
	Body []byte
}

type pipelineRPC struct {
	router *conf.Router
}

func (rpc *pipelineRPC) Run(args *Args, reply *[]byte) error {
	var output bytes.Buffer
	request := &ArgsRequest{
		uuid.New(),
		args,
	}

	if err := handleRequest(request, rpc.router, &output); err != nil {
		return err
	}

	// write answer back
	*reply = output.Bytes()

	return nil
}
