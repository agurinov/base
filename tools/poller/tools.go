package poller

func EventsToFds(events ...Event) []uintptr {
	fds := make([]uintptr, len(events))

	for i, event := range events {
		fds[i] = event.Fd()
	}

	return fds
}
