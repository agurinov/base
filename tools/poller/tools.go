package poller

func sliceEqual(a, b []uintptr) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func EventsToFds(events ...Event) []uintptr {
	fds := make([]uintptr, len(events))

	for i, event := range events {
		fds[i] = event.Fd()
	}

	return fds
}
