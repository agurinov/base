package tools

import (
	"sync"
	"sync/atomic"
)

// Once is an object that will perform exactly one action
// until Reset is called.
// See http://golang.org/pkg/sync/#Once
type Once struct {
	m    sync.Mutex
	lock uint32
}

// Do simulates sync.Once.Do by executing the specified function
// only once, until Reset is triggered after sucessfull f return.
// See http://golang.org/pkg/sync/#Once
func (o *Once) Do(f func(), reset bool) {
	if atomic.LoadUint32(&o.lock) == 1 {
		return
	}

	o.m.Lock()

	if o.lock == 0 {
		atomic.StoreUint32(&o.lock, 1)
		f()
	}

	if reset {
		atomic.StoreUint32(&o.lock, 0)
	}
	o.m.Unlock()
}

// Reset indicates that the next call to Do should actually be called
// once again.
func (o *Once) Reset() {
	o.m.Lock()
	atomic.StoreUint32(&o.lock, 0)
	o.m.Unlock()
}
