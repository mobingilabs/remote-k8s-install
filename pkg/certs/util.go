package certs

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func writeCertAndKey(pkiPath, name string, cert *x509.Certificate, key *rsa.PrivateKey) error {
	keyBlock := pem.Block{
		Type:  RSAPrivateKeyBlockType,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	keyData := pem.EncodeToMemory(&keyBlock)
	keyPath := pathForPrivateKey(pkiPath, name)
	if err := writeKey(keyPath, keyData); err != nil {
		return err
	}

	certBlock := pem.Block{
		Type:  CertificateBlockType,
		Bytes: cert.Raw,
	}
	certData := pem.EncodeToMemory(&certBlock)
	certPath := pathForCert(pkiPath, name)
	if err := writeKey(certPath, certData); err != nil {
		return err
	}

	return nil
}

func writePrivateKey(keyPath string, key *rsa.PrivateKey) error {
	block := pem.Block{
		Type:  RSAPrivateKeyBlockType,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	data := pem.EncodeToMemory(&block)

	return writeKey(keyPath, data)
}

func writePublicKey(keyPath string, key *rsa.PublicKey) error {
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return err
	}

	block := pem.Block{
		Type:  PublicKeyBlockType,
		Bytes: der,
	}

	data := pem.EncodeToMemory(&block)

	return writeKey(keyPath, data)
}

// writeKey writes the data to disk
func writeKey(keyPath string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(keyPath), os.FileMode(0755)); err != nil {
		return err
	}

	return ioutil.WriteFile(keyPath, data, os.FileMode(0600))
}

func pathForCert(pkiPath, name string) string {
	return filepath.Join(pkiPath, fmt.Sprintf("%s.crt", name))
}

func pathForPrivateKey(pkiPath, name string) string {
	return filepath.Join(pkiPath, fmt.Sprintf("%s.key", name))
}

func pathForPublicKey(pkiPath, name string) string {
	return filepath.Join(pkiPath, fmt.Sprintf("%s.pub", name))
}

