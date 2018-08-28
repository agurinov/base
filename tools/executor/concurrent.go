package executor

import (
	"context"
	"sync"
)

// https://stackoverflow.com/questions/45500836/close-multiple-goroutine-if-an-error-occurs-in-one-in-go

func concurrent(ctx context.Context, fns ...func(context.Context) error) error {
	errCh := make(chan error, 1)
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // Make sure it's called to release resources even if no errors

	wg.Add(len(fns))

	for _, fn := range fns {
		go func(fn func(context.Context) error) {
			defer wg.Done()

			// TODO defer panic?

			select {
			case <-ctx.Done():
				// context cancelled by another function
				// no need for starting execution of this fn
				return
			default: // Default is must to avoid blocking
			}

			// we can start atomic function execution
			if err := fn(ctx); err != nil {
				errCh <- err
				cancel()
				return
			}
		}(fn)
	}

	// wait until all completed
	wg.Wait()

	if ctx.Err() != nil {
		return <-errCh
	}

	return nil
}
