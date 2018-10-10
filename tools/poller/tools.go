package poller

func pendingFilterClosed(pending []*HeapItem, close []uintptr) []*HeapItem {
	new := make([]*HeapItem, 0)

OUTER:
	for _, item := range pending {
		for _, cfd := range close {
			if item.Fd == cfd {
				// fd from heap is closed
				// not relevant, delete it from future .Pop()
				// delete = skip from appending
				continue OUTER
			}
		}
		// not closed - use it
		new = append(new, item)
	}

	return new
}

func pendingMapReady(pending []*HeapItem, ready []uintptr) []*HeapItem {
OUTER:
	for _, item := range pending {
		for _, rfd := range ready {
			if item.Fd == rfd {
				item.ready = true
				continue OUTER
			}
		}
	}

	return pending
}

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
