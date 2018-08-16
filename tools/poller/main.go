package poller

// https://people.freebsd.org/~jlemon/papers/kqueue.pdf

type Event interface {
	Fd() uintptr
}

type Interface interface {
	Add(fd uintptr) error
	Del(fd uintptr) error
	Events() (re []Event, we []Event, err error)
}
