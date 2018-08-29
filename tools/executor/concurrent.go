package executor

import (
	"context"
	"sync"
)

func concurrent(ctx context.Context, fns ...func(context.Context) error) error {
	errCh := make(chan error, 1)
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(ctx)

	// Make sure it's called to release resources even if no errors
	defer cancel()

	for _, fn := range fns {
		wg.Add(1)

		go func(fn func(context.Context) error) {
			defer wg.Done()

			// TODO defer panic?

			// check for another func in another goroutine failed
			// and cancelled all waterfall
			select {
			case <-ctx.Done():
				// context cancelled by another function
				// no need for starting execution of this fn
				return
			default:
				// Default is must to avoid blocking
				// we can start atomic function execution
				if err := fn(ctx); err != nil {
					errCh <- err
					cancel()
					return
				}
			}
		}(fn)
	}

	// wait until all completed
	wg.Wait()

	// returning errors (from channel if exists or nil)
	select {
	case err := <-errCh:
		// context cancelled by another function
		// no need for starting execution of this fn
		return err
	default:
		// Default is must to avoid blocking
		return nil
	}
}
