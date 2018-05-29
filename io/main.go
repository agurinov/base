package io

type Pipeliner interface {
	setStdin(reader io.ReadCloser)
	setStdout(writer io.WriteCloser)
}

type RunCloser interface {
	run() (err error)
	close() (err error)
}

type PipeExec interface {
	Pipeliner
	RunCloser
}
