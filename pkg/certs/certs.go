/*
Copyright 2016 The Kubernetes Authors.

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
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/pkg/errors"
	certutil "k8s.io/client-go/util/cert"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"

	"mobingi/ocean/pkg/config"
)

const (
	PrivateKeyBlockType = "PRIVATE KEY"
	PublicKeyBlockType  = "PUBLIC KEY"

	CertificateBlockType   = "CERTIFICATE"
	RSAPrivateKeyBlockType = "RSA PRIVATE KEY"

	rsaKeySize   = 2048
	duration365d = time.Hour * 24 * 365
)

// CreatePKIAssets will create and write to disk all PKI assets necessary to establish the control plane.
// If the PKI assets already exists in the target folder, they are used only if evaluated equal; otherwise an error is returned.
func CreatePKIAssets(cfg *config.Config) error {
	certTree, err := getDefaultCertList().asMap().certTree()
	if err != nil {
		return err
	}
	if err := certTree.createTree(cfg); err != nil {
		return errors.Wrap(err, "error creating PKI assets")
	}

	// Service accounts are not x509 certs, so handled separately
	return createServiceAccountKeyAndPublicKeyFiles(cfg.Dir)
}

// CreateServiceAccountKeyAndPublicKeyFiles create a new public/private key files for signing service account users.
// If the sa public/private key files already exists in the target folder, they are used only if evaluated equals; otherwise an error is returned.
func createServiceAccountKeyAndPublicKeyFiles(certsDir string) error {
	privateKey, err := newPrivateKey()
	if err != nil {
		return err
	}

	privateKeyPath := pathForPrivateKey(certsDir, kubeadmconstants.ServiceAccountKeyBaseName)
	if err := writePrivateKey(privateKeyPath, privateKey); err != nil {
		return err
	}

	publicKeyPath := pathForPublicKey(certsDir, kubeadmconstants.ServiceAccountKeyBaseName)
	return writePublicKey(publicKeyPath, &privateKey.PublicKey)
}

func newPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, rsaKeySize)
}

func newCACertAndKey(certSpec *certutil.Config) (*x509.Certificate, *rsa.PrivateKey, error) {
	key, err := newPrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create private key:%s", err.Error())
	}

	cert, err := certutil.NewSelfSignedCACert(*certSpec, key)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create cert:%s", cert)
	}

	return cert, key, nil
}

func newSignedCert(certSpec *certutil.Config, key crypto.Signer, caCert *x509.Certificate, caKey crypto.Signer) (*x509.Certificate, error) {
	if err := validateCertSpec(certSpec); err != nil {
		return nil, err
	}
	certTmpl, err := createCertTmpl(certSpec, caCert.NotBefore)
	if err != nil {
		return nil, err
	}

	certDERBytes, err := x509.CreateCertificate(rand.Reader, certTmpl, caCert, key.Public(), caKey)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(certDERBytes)
}

func validateCertSpec(certSpec *certutil.Config) error {
	if len(certSpec.CommonName) == 0 {
		return errors.New("must specify a CommonName")
	}

	if len(certSpec.Usages) == 0 {
		return errors.New("must specify at least one ExtKeyUsage")
	}

	return nil
}

func createCertTmpl(certSpec *certutil.Config, notBefore time.Time) (*x509.Certificate, error) {
	serial, err := rand.Int(rand.Reader, new(big.Int).SetInt64(math.MaxInt64))
	if err != nil {
		return nil, err
	}
	return &x509.Certificate{
		Subject: pkix.Name{
			CommonName:   certSpec.CommonName,
			Organization: certSpec.Organization,
		},
		DNSNames:     certSpec.AltNames.DNSNames,
		IPAddresses:  certSpec.AltNames.IPs,
		SerialNumber: serial,
		NotBefore:    notBefore,
		NotAfter:     time.Now().Add(duration365d).UTC(),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  certSpec.Usages,
	}, nil
}
