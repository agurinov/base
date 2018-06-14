package conf

import (
	"errors"
	"regexp"

	"app/pattql"
	"app/pipeline"
)

// TODO look at Pipeline.UnmarshalYAML and remake this to type []Route
type Router struct {
	Collection []Route
}

func (rc *Router) Match(uri string) (*Route, error) {
	for _, route := range rc.Collection {
		if route.Match(uri) {
			return &route, nil
		}
	}

	return nil, errors.New("conf: Route not found")
}

type Route struct {
	regexp   *regexp.Regexp
	pipeline *pipeline.Pipeline
}

func (r *Route) Match(uri string) bool {
	return r.regexp.MatchString(uri)
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
