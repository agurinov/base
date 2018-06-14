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
      address: geoiphost:100500

# Route
- pattern: "/data/{*}.jpg"
  pipeline:

    - type: process
      cmd: "echo 'HEAD / HTTP/1.1\r\n\r\n'"

    - type: tcp
      address: golang.org:80

    - type: process
      cmd: "cat /dev/stdin"

# Route
- pattern: "{*}"
  pipeline:

    - type: process
      cmd: "echo 'nobody home...'"`

func TestRouteUnmarshalYAML(t *testing.T) {
	var router Router

	if err := yaml.Unmarshal([]byte(YAML), &router); err != nil {
		t.Fatal(err)
	}

	// exist url
	route, err := router.Match("data/foobar.jpg")
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

	// default route
	route, err = router.Match("/foobar/bar/baz")
	if err != nil {
		t.Fatal(err)
	}

	input = readCloser{bytes.NewBuffer([]byte{})}
	output = writeCloser{bytes.NewBuffer([]byte{})}

	err = route.pipeline.Run(input, output)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output)

}
