package conf

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"
)

const YAML = `collection:

# Route
- pattern: "/foo/bar"
  pipeline:

    - type: socket
      address: tcp://geoiphost/foo/bar

# Route
- pattern: "/data/{*}.jpg"
  pipeline:

    - type: socket
      address: tcp://geoiphost

    - type: process
      name: "echo -n 'No answer'"
      foo: bar`

func TestRouteUnmarshalYAML(t *testing.T) {
	var rc RouteCollection

	if err := yaml.Unmarshal([]byte(YAML), &rc); err != nil {
		t.Error(err)
	}

	fmt.Printf("GHJKHGHJK   %+v\n", rc.Collection[0])

}
