package conf

import (
	"context"
	"errors"
	"io"
	"net/url"
	"regexp"

	"github.com/boomfunc/base/pattql"
	"github.com/boomfunc/base/pipeline"
)

var (
	ErrRouteNotFound = errors.New("conf: Route not found")
)

// TODO look at Pipeline.UnmarshalYAML and remake this to type []Route
type Router struct {
	Collection []Route
}

func (rc *Router) Match(url *url.URL) (*Route, error) {
	for _, route := range rc.Collection {
		if route.match(url.RequestURI()) {
			return &route, nil
		}
	}

	return nil, ErrRouteNotFound
}

type Route struct {
	regexp   *regexp.Regexp
	pipeline *pipeline.Pipeline
}

func (r *Route) match(uri string) bool {
	// TODO here extend context for URl values
	return r.regexp.MatchString(uri)
}

func (r *Route) Run(ctx context.Context, input io.Reader, output io.Writer) error {
	return r.pipeline.Run(ctx, input, output)
}

func (r *Route) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// inner struct for accepting strings
	var route struct {
		Pattern  string
		Pipeline *pipeline.Pipeline
	}

	if err := unmarshal(&route); err != nil {
		return err
	}

	// yaml valid, transform it
	r.regexp = pattql.Regexp(route.Pattern)
	r.pipeline = route.Pipeline

	return nil
}
