package poller

import (
	"errors"
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
	invokes int
	err     bool
	re      []uintptr
	ce      []uintptr
}

func MockPoller(re, ce []uintptr, err bool) Interface {
	return &mock{
		err: err,
		re:  re,
		ce:  ce,
	}
}

func (p *mock) Add(fd uintptr) error {
	if p.err {
		// case when error from poller
		return errors.New("Error from poller")
	}

	return nil
}

func (p *mock) Del(fd uintptr) error {
	if p.err {
		// case when error from poller
		return errors.New("Error from poller")
	}

	return nil
}

func (p *mock) events() ([]Event, []Event, []Event) {
	var re, ce []Event

	for _, fd := range p.re {
		re = append(re, mockEvent{fd})
	}
	for _, fd := range p.ce {
		ce = append(ce, mockEvent{fd})
	}

	return re, nil, ce
}

func (p *mock) Events() ([]Event, []Event, []Event, error) {
	p.invokes++

	if p.invokes == 1 {
		// FIRST INVOKE - BY RULE
		if p.err {
			// case when error from poller
			return nil, nil, nil, errors.New("Error from poller")
		}
		// imitation polling
		time.Sleep(time.Second * 1)

		re, _, ce := p.events()
		return re, nil, ce, nil
	} else {
		// SECOND - RETURN hardcoded events (otherwise block)
		return []Event{mockEvent{100}, mockEvent{500}}, nil, []Event{mockEvent{200}, mockEvent{300}}, nil
	}
}
