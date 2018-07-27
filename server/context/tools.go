package context

import (
	"context"
)

func fromCtx(ctx context.Context) (*values, error) {
	// ctx.Value return interface{}
	// non existence key -> nil -> ok will be false
	// wrong type case -> ok will be also false
	values, ok := ctx.Value("values").(*values)
	if !ok {
		return nil, ErrContextBroken
	}

	return values, nil
}
