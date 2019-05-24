package certs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePKIAssets(t *testing.T) {
	o := Options{
		InternalEndpoint: "192.168.1.1",
		ExternalEndpoint: "47.40.40.1",
		SANs:             []string{},
	}
	_, err := CreatePKIAssets(o)
	assert.NoError(t, err, "create pki assests err:%v", err)

}
