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

	certutil "k8s.io/client-go/util/cert"

	"mobingi/ocean/pkg/constants"
)

type cert struct {
	name   string
	config certutil.Config
}

func newAPIServerCert(o Options) *cert {
	return &cert{
		name: "apiserver",
		config: certutil.Config{
			CommonName: constants.APIServerCertCommonName,
			Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		// TODO makeAltNamesMutator(getAPIServerAltNames),
	}
}

func newAPIServerKubeletClientCert(o Options) *cert {
	return &cert{
		name: "apiserver-kubelet-client",
		config: certutil.Config{
			CommonName:   constants.APIServerKubeletClientCertCommonName,
			Organization: []string{constants.MastersGroup},
			Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
	}
}

func newAPIServerEtcdClientCert(o Options) *cert {
	return &cert{
		name: "apiserver-etcd-client",
		config: certutil.Config{
			CommonName:   constants.APIServerEtcdClientCertCommonName,
			Organization: []string{constants.MastersGroup},
			Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		/*
			configMutators: []configMutatorsFunc{
				makeAltNamesMutator(getEtcdAltNames),
			},*/
	}
}

func newEtcdServerCert(o Options) *cert {
	return &cert{
		name: "etcd-server",
		config: certutil.Config{
			//TODO CommonName: "xxx",
			Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		},
		// TODO makeAltNamesMutator(getEtcdAltNames),
	}
}

func newEtcdPeerCert(o Options) *cert {
	return &cert{
		name: "etcd-peer",
		config: certutil.Config{
			Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		},
		// TODO		makeAltNamesMutator(getEtcdAltNames),setCommonNameToNodeName(),
	}
}

func getCertList(o Options) []*cert {
	return []*cert{
		newAPIServerCert(o),
		newAPIServerKubeletClientCert(o),
		newAPIServerEtcdClientCert(o),
		newEtcdServerCert(o),
		newEtcdPeerCert(o),
	}
}

/*
func setCommonNameToNodeName() configMutatorsFunc {
	return func(cc *certutil.Config, cfg *config) error {
		//TODO	cc.CommonName = cfg.NodeRegistration.Name
		cc.CommonName = "etcd"
		return nil
	}
}
*/
