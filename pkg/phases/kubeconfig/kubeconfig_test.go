package kubeconfig

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKubeconf(t *testing.T) {
	o := newOptions()
	datas, err := NewKubeconf(o)
	assert.NoError(t, err, "new kubeconf err:", err)
	for k, v := range datas {
		ioutil.WriteFile(k, v, 0444)
	}
}

func newOptions() Options {
	ca, _ := ioutil.ReadFile("testdata/ca.crt")
	key, _ := ioutil.ReadFile("testdata/ca.key")
	return Options{
		CaCert:           ca,
		CaKey:            key,
		ExternalEndpoint: "https://47.42.1.2:6443",
		InternalEndpoint: "https://192.168.1.1:6443",
		Nodes:            []string{"node0", "node1", "node2"},
		ClusterName:      "kubernetes",
	}
}
