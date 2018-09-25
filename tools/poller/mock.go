package poller

import (
	"time"
)

type mockEvent struct {
	se uintptr
}

// Fd implements base Event interface
func (ev mockEvent) Fd() uintptr {
	return ev.se
}

type mock struct {
	delay time.Duration
}

func (p *mock) Add(fd uintptr) error {
	return nil
}

func (p *mock) Del(fd uintptr) error {
	return nil
}

func (p *mock) Events() ([]Event, []Event, []Event, error) {
	// imitation polling
	time.Sleep(time.Second * 3)
	return []Event{mockEvent{1}, mockEvent{2}}, []Event{}, []Event{}, nil
}
