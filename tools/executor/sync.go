package executor

import (
	"context"
)

func Sync(fns ...func() error) (int, error) {
	for i, fn := range fns {
		if err := fn(); err != nil {
			// some function return error
			return i, err
		}
	}
	// no error
	return len(fns), nil
}

func SyncWithContext(ctx context.Context, fns ...func(context.Context) error) (int, error) {
	for i, fn := range fns {
		if err := fn(ctx); err != nil {
			// some function return error
			return i, err
		}
	}
	// no error
	return len(fns), nil
}
