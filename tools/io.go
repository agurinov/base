package tools

import (
	"io"
)

type readCloser struct {
	io.Reader
}

func (readCloser) Close() error { return nil }

type writeCloser struct {
	io.Writer
}

func (writeCloser) Close() error { return nil }

type readWriteCloser struct {
	io.ReadWriter
}

func (readWriteCloser) Close() error { return nil }

func ReadCloser(r io.Reader) io.ReadCloser {
	if rc, ok := r.(io.ReadCloser); ok {
		return rc
	}
	return &readCloser{r}
}

func WriteCloser(w io.Writer) io.WriteCloser {
	if wc, ok := w.(io.WriteCloser); ok {
		return wc
	}
	return &writeCloser{w}
}

func ReadWriteCloser(rw io.ReadWriter) io.ReadWriteCloser {
	if rwc, ok := rw.(io.ReadWriteCloser); ok {
		return rwc
	}
	return &readWriteCloser{rw}
}
