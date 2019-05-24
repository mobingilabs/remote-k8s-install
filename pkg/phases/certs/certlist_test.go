/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package certs

import (
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/assert"
	certutil "k8s.io/client-go/util/cert"
)

func TestGetDefaultCerts(t *testing.T) {
	needCerts := map[string]struct{}{
		"ca":                       struct{}{},
		"apiserver":                struct{}{},
		"apiserver-kubelet-client": struct{}{},

		"front-proxy-ca":     struct{}{},
		"front-proxy-client": struct{}{},

		"etcd-ca":                 struct{}{},
		"etcd-server":             struct{}{},
		"etcd-peer":               struct{}{},
		"etcd-healthcheck-client": struct{}{},
		"apiserver-etcd-client":   struct{}{},
	}
	certList := getDefaultCerts()
	for _, v := range certList {
		delete(needCerts, v.Name)
	}
	assert.EqualValuesf(t, 0, len(needCerts), "cert:%v lost", needCerts)
}

func TestNewCertAndKeyFromCA(t *testing.T) {
	certspec := certutil.Config{}
	rootCa, rootKey, err := newCACertAndKey(&certspec)
	assert.NoErrorf(t, err, "new ca cert and key err:%v", err)

	// TODO missig priKey test
	cert, _, err := certAPIServer.newCertAndKeyFromCA(newTestConfig(), rootCa, rootKey)
	assert.NoErrorf(t, err, "new ca cert and key from ca err:%v", err)

	certPool := x509.NewCertPool()
	certPool.AddCert(rootCa)
	if _, err := cert.Verify(x509.VerifyOptions{
		Roots: certPool,
	}); err != nil {
		t.Logf("cert verify err:%v", err)
	}
}

func newTestConfig() *config {
	return &config{
		AdvertiseAddress: "192.168.1.1",
		PublicIP:         "47.102.23.2",

		SANs: []string{},
	}
}
