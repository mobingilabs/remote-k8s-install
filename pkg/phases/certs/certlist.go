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
	"crypto/rsa"
	"crypto/x509"

	"github.com/pkg/errors"

	certutil "k8s.io/client-go/util/cert"

	"mobingi/ocean/pkg/constants"
	pkiutil "mobingi/ocean/pkg/util/pki"
)

type configMutatorsFunc func(*certutil.Config, *config) error

// Cert represents a certificate that will create to function properly.
type cert struct {
	Name     string
	BaseName string
	CAName   string
	// Some attributes will depend on the InitConfiguration, only known at runtime.
	// These functions will be run in series, passed both the InitConfiguration and a cert Config.
	configMutators []configMutatorsFunc
	config         certutil.Config
}

// GetConfig returns the definition for the given cert given the provided InitConfiguration
func (c *cert) getConfig(cfg *config) (*certutil.Config, error) {
	for _, f := range c.configMutators {
		if err := f(&c.config, cfg); err != nil {
			return nil, err
		}
	}

	return &c.config, nil
}

func (c *cert) newCertAndKeyFromCA(cfg *config, caCert *x509.Certificate, caKey *rsa.PrivateKey) (*x509.Certificate, *rsa.PrivateKey, error) {
	certSpec, err := c.getConfig(cfg)
	if err != nil {
		return nil, nil, err
	}

	key, err := newPrivateKey()
	if err != nil {
		return nil, nil, err
	}

	cert, err := newSignedCert(certSpec, key, caCert, caKey)
	if err != nil {
		return nil, nil, err
	}

	return cert, key, nil
}

// CertificateTree is represents a one-level-deep tree, mapping a CA to the certs that depend on it.
type ceretificates []*cert
type certificateTree map[*cert]certificates

// CreateTree creates the CAs, certs signed by the CAs, and writes them all to disk.
func (t certificateTree) createTree(cfg *config) (map[string][]byte, error) {
	certs := make(map[string][]byte)

	for ca, leaves := range t {
		certSpec, err := ca.getConfig(cfg)
		if err != nil {
			return nil, err
		}

		caCert, caKey, err := newCACertAndKey(certSpec)
		if err != nil {
			return nil, err
		}
		certs[pathForCert(ca.BaseName)] = pkiutil.EncodeCertPEM(caCert)
		certs[pathForKey(ca.BaseName)] = pkiutil.EncodePrivateKeyPEM(caKey)

		for _, leaf := range leaves {
			cert, key, err := leaf.newCertAndKeyFromCA(cfg, caCert, caKey)
			if err != nil {
				return nil, err
			}
			certs[pathForCert(leaf.BaseName)] = pkiutil.EncodeCertPEM(cert)
			certs[pathForKey(leaf.BaseName)] = pkiutil.EncodePrivateKeyPEM(key)
		}
	}

	return certs, nil
}

// CertificateMap is a flat map of certificates, keyed by Name.
type certificateMap map[string]*cert

// CertTree returns a one-level-deep tree, mapping a CA cert to an array of certificates that should be signed by it.
func (m certificateMap) certTree() (certificateTree, error) {
	caMap := make(certificateTree)

	for _, c := range m {
		if c.CAName == "" {
			if _, ok := caMap[c]; !ok {
				caMap[c] = []*cert{}
			}
		} else {
			ca, ok := m[c.CAName]
			if !ok {
				return nil, errors.Errorf("certificate %q references unknown CA %q", c.Name, c.CAName)
			}
			caMap[ca] = append(caMap[ca], c)
		}
	}

	return caMap, nil
}

// Certificates is a list of Certificates that Kubeadm should create.
type certificates []*cert

// AsMap returns the list of certificates as a map, keyed by name.
func (c certificates) asMap() certificateMap {
	certMap := make(map[string]*cert)
	for _, cert := range c {
		certMap[cert.Name] = cert
	}

	return certMap
}

// GetDefaultCertList returns  all of the certificates kubeadm requires to function.
func getDefaultCertList() certificates {
	return certificates{
		&certRootCA,
		&certAPIServer,
		&certKubeletClient,
		// Front Proxy certs
		&certFrontProxyCA,
		&certFrontProxyClient,
		// etcd certs
		&certEtcdCA,
		&certEtcdServer,
		&certEtcdPeer,
		&certEtcdHealthcheck,
		&certEtcdAPIClient,
	}
}

var (
	// CertRootCA is the definition of the Kubernetes Root CA for the API Server and kubelet.
	certRootCA = cert{
		Name:     "ca",
		BaseName: constants.CACertAndKeyBaseName,
		config: certutil.Config{
			CommonName: "kubernetes",
		},
	}
	// CertAPIServer is the definition of the cert used to serve the Kubernetes API.
	certAPIServer = cert{
		Name:     "apiserver",
		BaseName: constants.APIServerCertAndKeyBaseName,
		CAName:   "ca",
		config: certutil.Config{
			CommonName: constants.APIServerCertCommonName,
			Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		configMutators: []configMutatorsFunc{
			makeAltNamesMutator(getAPIServerAltNames),
		},
	}

	// certKubeletClient is the definition of the cert used by the API server to access the kubelet.
	certKubeletClient = cert{
		Name:     "apiserver-kubelet-client",
		BaseName: constants.APIServerKubeletClientCertAndKeyBaseName,
		CAName:   "ca",
		config: certutil.Config{
			CommonName:   constants.APIServerKubeletClientCertCommonName,
			Organization: []string{constants.MastersGroup},
			Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
	}

	// certFrontProxyCA is the definition of the CA used for the front end proxy.
	certFrontProxyCA = cert{
		Name:     "front-proxy-ca",
		BaseName: constants.FrontProxyCACertAndKeyBaseName,
		config: certutil.Config{
			CommonName: "front-proxy-ca",
		},
	}

	// certFrontProxyClient is the definition of the cert used by the API server to access the front proxy.
	certFrontProxyClient = cert{
		Name:     "front-proxy-client",
		BaseName: constants.FrontProxyClientCertAndKeyBaseName,
		CAName:   "front-proxy-ca",
		config: certutil.Config{
			CommonName: constants.FrontProxyClientCertCommonName,
			Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
	}

	// KubeadmCertEtcdCA is the definition of the root CA used by the hosted etcd server.
	certEtcdCA = cert{
		Name:     "etcd-ca",
		BaseName: constants.EtcdCACertAndKeyBaseName,
		config: certutil.Config{
			CommonName: "etcd-ca",
		},
	}

	// certEtcdServer is the definition of the cert used to serve etcd to clients.
	certEtcdServer = cert{
		Name:     "etcd-server",
		BaseName: constants.EtcdServerCertAndKeyBaseName,
		CAName:   "etcd-ca",
		config: certutil.Config{
			// TODO: etcd 3.2 introduced an undocumented requirement for ClientAuth usage on the
			// server cert: https://github.com/coreos/etcd/issues/9785#issuecomment-396715692
			// Once the upstream issue is resolved, this should be returned to only allowing
			// ServerAuth usage.
			Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		},
		configMutators: []configMutatorsFunc{
			makeAltNamesMutator(getEtcdAltNames),
			setCommonNameToNodeName(),
		},
	}
	// certEtcdPeer is the definition of the cert used by etcd peers to access each other.
	certEtcdPeer = cert{
		Name:     "etcd-peer",
		BaseName: constants.EtcdPeerCertAndKeyBaseName,
		CAName:   "etcd-ca",
		config: certutil.Config{
			Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		},
		configMutators: []configMutatorsFunc{
			makeAltNamesMutator(getEtcdAltNames),
			setCommonNameToNodeName(),
		},
	}
	// certEtcdHealthcheck is the definition of the cert used by Kubernetes to check the health of the etcd server.
	certEtcdHealthcheck = cert{
		Name:     "etcd-healthcheck-client",
		BaseName: constants.EtcdHealthcheckClientCertAndKeyBaseName,
		CAName:   "etcd-ca",
		config: certutil.Config{
			CommonName:   constants.EtcdHealthcheckClientCertCommonName,
			Organization: []string{constants.MastersGroup},
			Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
	}
	// certEtcdAPIClient is the definition of the cert used by the API server to access etcd.
	certEtcdAPIClient = cert{
		Name:     "apiserver-etcd-client",
		BaseName: constants.APIServerEtcdClientCertAndKeyBaseName,
		CAName:   "etcd-ca",
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
)

func setCommonNameToNodeName() configMutatorsFunc {
	return func(cc *certutil.Config, cfg *config) error {
		//TODO	cc.CommonName = cfg.NodeRegistration.Name
		cc.CommonName = "etcd"
		return nil
	}
}
