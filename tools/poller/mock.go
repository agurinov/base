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

type mock struct{}

func (p *mock) Add(fd uintptr) error {
	return nil
}

func (p *mock) Del(fd uintptr) error {
	return nil
}

func (p *mock) Events() ([]Event, []Event, []Event, error) {
	// imitation polling
	time.Sleep(time.Second * 3)
	re := []Event{mockEvent{1}, mockEvent{2}, mockEvent{3}, mockEvent{4}, mockEvent{5}, mockEvent{6}}
	we := []Event{}
	ce := []Event{}
	return re, we, ce, nil
}
