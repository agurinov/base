package io

type Pipeliner interface {
	pipe(reader io.ReadCloser, writer io.WriteCloser)
}

type RunCloser interface {
	run() (err error)
	close() (err error)
}

type PipeExec interface {
	Pipeliner
	RunCloser
}
