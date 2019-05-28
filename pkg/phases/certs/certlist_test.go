package certs

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPKIAssets(t *testing.T) {
	data, err := NewPKIAssets(newOptions())
	assert.NoError(t, err, "new pki assests err")
	os.Mkdir("test", 0777)
	for k, v := range data {
		ioutil.WriteFile(fmt.Sprintf("test/%s", k), v, 0444)
	}
}

func newOptions() Options {
	return Options{
		InternalEndpoint: "192.168.1.5",
		ExternalEndpoint: "47.40.4.4",
		SANs:             []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"},
		ServiceSubnet:    "10.96.0.0/12",
	}
}
