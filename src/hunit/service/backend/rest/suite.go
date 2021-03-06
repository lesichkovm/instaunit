package rest

import (
	"io"
	"io/ioutil"
	"sync"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// A request
type Request struct {
	sync.Mutex
	Methods []string          `yaml:"methods"`
	Path    string            `yaml:"path"`
	Params  map[string]string `yaml:"params"`
	Headers map[string]string `yaml:"headers"`
	Cookies map[string]string `yaml:"cookies"`
	Entity  string            `yaml:"entity"`
}

// A response
type Response struct {
	Status  int               `yaml:"status"`
	Headers map[string]string `yaml:"headers"`
	Cookies map[string]string `yaml:"cookies"`
	Entity  string            `yaml:"entity"`
}

// An endpoint
type Endpoint struct {
	Wait     time.Duration `yaml:"wait"`
	Request  *Request      `yaml:"endpoint"`
	Response *Response     `yaml:"response"`
}

// A test suite
type Suite struct {
	Endpoints []Endpoint `yaml:"service"`
}

// Load a test suite
func LoadSuite(src io.ReadCloser) (*Suite, error) {
	suite := &Suite{}
	var ferr error

	data, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, suite)
	if err != nil {
		ferr = err
	}

	if len(suite.Endpoints) < 1 {
		var endpoints []Endpoint
		err := yaml.Unmarshal(data, &endpoints)
		if err != nil {
			return nil, coalesce(ferr, err)
		} else {
			suite.Endpoints = endpoints
		}
	}

	return suite, nil
}

// Return the first non-nil error or nil if there are none.
func coalesce(err ...error) error {
	for _, e := range err {
		if e != nil {
			return e
		}
	}
	return nil
}
