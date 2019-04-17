package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRunEtcdJobs(t *testing.T) {
	clusterIPs := []string{"192.168.1.0", "192.168.1.1", "192.168.1.2"}
	jobs, err := NewRunEtcdJobs(clusterIPs)
	assert.Nil(t, err)
	assert.Equal(t, len(clusterIPs), len(jobs))
}

func TestGetEtcdServers(t *testing.T) {
	clusterIPs := []string{"192.168.1.0", "192.168.1.1", "192.168.1.2"}
	etcdServers := GetEtcdServers(clusterIPs)
	expectEtcdServers := fmt.Sprintf("https://192.168.1.0:%d,https://192.168.1.1:%d,https://192.168.1.2:%d", EtcdClientPort, EtcdClientPort, EtcdClientPort)
	assert.Equal(t, expectEtcdServers, etcdServers)
}
