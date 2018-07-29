package context

import (
	"context"
	"errors"
	"net/url"
)

var (
	ErrContextBroken = errors.New("server/context: Something went wrong with server context")
)

// meta is soecial map for server and request data
// Q is request query
// Url is object passed from router regex matching
// for example:
// pattern: /foo/{bar}?bar=baz
// url: /foo/lolkek?bar=baz
// q.bar == baz
// url.bar == lolkek
type values struct {
	meta map[string]interface{}
	Q    url.Values
	Url  url.Values
}

func New() context.Context {
	values := new(values)
	values.meta = make(map[string]interface{})
	values.Q = url.Values{}
	values.Url = url.Values{}

	return context.WithValue(context.Background(), "values", values)
}

func Values(ctx context.Context) (*values, error) {
	return fromCtx(ctx)
}

func SetMeta(ctx context.Context, key string, value interface{}) error {
	values, err := fromCtx(ctx)
	if err != nil {
		return err
	}

	values.meta[key] = value
	return nil
}

func GetMeta(ctx context.Context, key string) (interface{}, error) {
	values, err := fromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return values.meta[key], nil
}

func GetQ(ctx context.Context, key string) (string, error) {
	values, err := fromCtx(ctx)
	if err != nil {
		return "", err
	}

	return values.Q.Get(key), nil
}

func GetUrl(ctx context.Context, key string) (string, error) {
	values, err := fromCtx(ctx)
	if err != nil {
		return "", err
	}

	return values.Url.Get(key), nil
}
