package certs

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/cache"
)

func writeCertAndKey(baseName string, cert *x509.Certificate, key *rsa.PrivateKey) {
	keyBlock := pem.Block{
		Type:  RSAPrivateKeyBlockType,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	keyData := pem.EncodeToMemory(&keyBlock)
	writeKey(baseName, keyData)

	certBlock := pem.Block{
		Type:  CertificateBlockType,
		Bytes: cert.Raw,
	}
	certData := pem.EncodeToMemory(&certBlock)
	writeCert(baseName, certData)
}

func writePrivateKey(baseName string, key *rsa.PrivateKey) {
	block := pem.Block{
		Type:  RSAPrivateKeyBlockType,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	data := pem.EncodeToMemory(&block)

	writeKey(baseName, data)
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

	writePub(keyPath, data)

	return nil
}

// write put the pki data to cache
func writeKey(baseName string, data []byte) {
	name := fmt.Sprintf("%s.key", baseName)
	cache.Put(constants.CertPrefix, name, data)
}

func writeCert(baseName string, data []byte) {
	name := fmt.Sprintf("%s.crt", baseName)
	cache.Put(constants.CertPrefix, name, data)
}

func writePub(baseName string, data []byte) {
	name := fmt.Sprintf("%s.pub", baseName)
	cache.Put(constants.CertPrefix, name, data)
}

func pathForCert(baseName string) string {
	return fmt.Sprintf("%s.crt", baseName)
}

func pathForKey(baseName string) string {
	return fmt.Sprintf("%s.key", baseName)
}

func pathForPub(baseName string) string {
	return fmt.Sprintf("%s.pub", baseName)
}

func certToByte(cert *x509.Certificate) []byte {
	certBlock := pem.Block{
		Type:  CertificateBlockType,
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(&certBlock)
}

func keyToByte(key *rsa.PrivateKey) []byte {
	keyBlock := pem.Block{
		Type:  RSAPrivateKeyBlockType,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return pem.EncodeToMemory(&keyBlock)
}

func pubKeyToByte(key *rsa.PublicKey) ([]byte, error) {
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, err
	}

	block := pem.Block{
		Type:  PublicKeyBlockType,
		 Bytes: der,
	}
	data := pem.EncodeToMemory(&block)
	return data, nil
}
