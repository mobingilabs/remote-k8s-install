package pki

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"path/filepath"
	"time"

	certutil "k8s.io/client-go/util/cert"
)

const (
	PrivateKeyBlockType = "PRIVATE KEY"
	PublicKeyBlockType  = "PUBLIC KEY"

	CertificateBlockType   = "CERTIFICATE"
	RSAPrivateKeyBlockType = "RSA PRIVATE KEY"

	rsaKeySize   = 2048
	duration365d = time.Hour * 24 * 365
)

func TryLoadCertAndKeyFromDisk(pkiPath, name string) (*x509.Certificate, *rsa.PrivateKey, error) {
	cert, err := tryLoadCertFromDisk(pkiPath, name)
	if err != nil {
		return nil, nil, err
	}
	key, err := tryLoadKeyFromDisk(pkiPath, name)
	if err != nil {
		return nil, nil, err
	}
	return cert, key, nil
}

func tryLoadCertFromDisk(pkiPath, name string) (*x509.Certificate, error) {
	certPath := pathForCert(pkiPath, name)
	certs, err := certutil.CertsFromFile(certPath)
	if err != nil {
		return nil, err
	}

	return certs[0], nil
}

func tryLoadKeyFromDisk(pkiPath, name string) (*rsa.PrivateKey, error) {
	keyPath := pathForKey(pkiPath, name)
	data, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	key, err := parsePrivateKeyPEM(data)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// TODO duplicate with certs/util.go
func pathForCert(pkiPath, name string) string {
	return filepath.Join(pkiPath, fmt.Sprintf("%s.crt", name))
}

func pathForKey(pkiPath, name string) string {
	return filepath.Join(pkiPath, fmt.Sprintf("%s.key", name))
}

func parsePrivateKeyPEM(data []byte) (*rsa.PrivateKey, error) {
	privateKeyPemBlock, _ := pem.Decode(data)
	if parsePrivateKeyPEM == nil {
		return nil, errors.New("can not parse key")
	}
	if privateKeyPemBlock.Type != RSAPrivateKeyBlockType {
		return nil, errors.New("not rsa type")
	}
	key, err := x509.ParsePKCS1PrivateKey(privateKeyPemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func NewCertAndKeyFromCA(caCert *x509.Certificate, caKey *rsa.PrivateKey, certSpec *certutil.Config) (*x509.Certificate, *rsa.PrivateKey, error) {
	key, err := NewPrivateKey()
	if err != nil {
		return nil, nil,err
	}

	cert, err := NewSignedCert(certSpec, key, caCert, caKey)
	if err != nil {
		return nil, nil,err
	}

	return cert, key, nil
}

func NewPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, rsaKeySize)
}

func NewSignedCert(certSpec *certutil.Config, key crypto.Signer, caCert *x509.Certificate, caKey crypto.Signer) (*x509.Certificate, error) {
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

func EncodeCertPEM(cert *x509.Certificate) []byte {
	block := pem.Block{
		Type:  CertificateBlockType,
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(&block)
}

func EncodePrivateKeyPEM(key *rsa.PrivateKey) []byte {
	block := pem.Block{
		Type:  RSAPrivateKeyBlockType,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return pem.EncodeToMemory(&block)
}