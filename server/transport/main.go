package transport

import (
	"container/heap"
)

// Interface describes main idea of transport layer.
// Receiving data from outside (internet, cron, etc) and pass to some heap (created by server)
// All Kind or errors to channel (also created by server)
type Interface interface {
	Connect(heap.Interface, chan error)
	Serve()
}
