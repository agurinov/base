package conf

import (
	"bytes"
	"io"
	"testing"

	"gopkg.in/yaml.v2"
)

type readCloser struct {
	io.Reader
}

func (readCloser) Close() error { return nil }

type writeCloser struct {
	io.Writer
}

func (writeCloser) Close() error { return nil }

const YAML = `collection:

# Route
- pattern: "/foo/bar"
  pipeline:

    - type: tcp
      address: tcp://geoiphost/foo/bar

# Route
- pattern: "/data/{*}.jpg"
  pipeline:

    - type: process
      cmd: "echo 'HEAD / HTTP/1.1\r\n\r\n'"

    - type: tcp
      address: golang.org:80

    - type: process
      cmd: "cat /dev/stdin"`

func TestRouteUnmarshalYAML(t *testing.T) {
	var rc RouteCollection

	if err := yaml.Unmarshal([]byte(YAML), &rc); err != nil {
		t.Fatal(err)
	}

	route, err := rc.Match("data/foobar.jpg")
	if err != nil {
		t.Fatal(err)
	}

	input := readCloser{bytes.NewBuffer([]byte("foobar"))}
	output := writeCloser{bytes.NewBuffer([]byte{})}

	err = route.pipeline.Run(input, output)
	if err != nil {
		t.Fatal(err)
	}

	// t.Log(output)

}
