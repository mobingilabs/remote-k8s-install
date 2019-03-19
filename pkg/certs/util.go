package certs

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"mobingi/ocean/pkg/log"
	"path/filepath"

	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/ssh"
	"mobingi/ocean/pkg/tools/cache"
	cmdutil "mobingi/ocean/pkg/util/cmd"
)

func writeCertAndKey(c ssh.Client, name string, cert *x509.Certificate, key *rsa.PrivateKey) error {
	keyBlock := pem.Block{
		Type:  RSAPrivateKeyBlockType,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	keyData := pem.EncodeToMemory(&keyBlock)
	if err := writeKey(c, name, keyData); err != nil {
		return err
	}

	certBlock := pem.Block{
		Type:  CertificateBlockType,
		Bytes: cert.Raw,
	}
	certData := pem.EncodeToMemory(&certBlock)
	if err := writeCert(c, name, certData); err != nil {
		return err
	}
	log.Infof("write cert and key:%s sucessed", name)

	return nil
}

func writePrivateKey(c ssh.Client, keyPath string, key *rsa.PrivateKey) error {
	block := pem.Block{
		Type:  RSAPrivateKeyBlockType,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	data := pem.EncodeToMemory(&block)

	return writeKey(c, keyPath, data)
}

func writePublicKey(c ssh.Client, keyPath string, key *rsa.PublicKey) error {
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return err
	}

	block := pem.Block{
		Type:  PublicKeyBlockType,
		Bytes: der,
	}

	data := pem.EncodeToMemory(&block)

	return writePub(c, keyPath, data)
}

// writeKey writes the data to disk
func writeKey(c ssh.Client, name string, data []byte) error {
	keyPath := filepath.Join(constants.PKIDir, fmt.Sprintf("%s.key", name))
	cache.Put(name, data)
	cmd := cmdutil.NewWriteCmd(keyPath, string(data))
	// TODO check output exec result, ok or false
	_, err := c.Do(cmd)
	return err
}

func writeCert(c ssh.Client, name string, data []byte) error {
	keyPath := filepath.Join(constants.PKIDir, fmt.Sprintf("%s.crt", name))
	cache.Put(name, data)
	_, err := c.Do(cmdutil.NewWriteCmd(keyPath, string(data)))

	return err
}

func writePub(c ssh.Client, name string, data []byte) error {
	keyPath := filepath.Join(constants.PKIDir, fmt.Sprintf("%s.pub", name))
	cache.Put(name, data)
	_, err := c.Do(cmdutil.NewWriteCmd(keyPath, string(data)))

	return err
}
