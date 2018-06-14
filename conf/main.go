package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func LoadFile(filename string) (*Router, error) {
	// Try to read config file
	YAML, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// config file loaded, try to parse yaml into router
	var router Router

	if err := yaml.Unmarshal(YAML, &router); err != nil {
		return nil, err
	}

	return &router, nil
}
