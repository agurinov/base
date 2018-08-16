package poller

type Event interface {
	Fd() uintptr
}

type Interface interface {
	Add(fd uintptr) error
	Del(fd uintptr) error
	Events() (re []Event, we []Event, err error)
}
