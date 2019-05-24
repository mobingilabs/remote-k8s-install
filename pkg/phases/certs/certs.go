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

	"mobingi/ocean/pkg/constants"
)

const (
	PrivateKeyBlockType = "PRIVATE KEY"
	PublicKeyBlockType  = "PUBLIC KEY"

	CertificateBlockType   = "CERTIFICATE"
	RSAPrivateKeyBlockType = "RSA PRIVATE KEY"

	rsaKeySize     = 2048
	duration10year = 10 * time.Hour * 24 * 365
)

type config struct {
	AdvertiseAddress string
	PublicIP         string

	SANs []string //now is master machines private ip
}

type Options struct {
	InternalEndpoint string
	ExternalEndpoint string
	SANs             []string
}

func NewRootCACert() ([]byte, []byte, error) {
	key, err := newPrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create private key:%s", err.Error())
	}

	cfg := certutil.Config{
		CommonName: "kubernetes",
		Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}
	cert, err := certutil.NewSelfSignedCACert(cfg, key)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create private key:%s", err.Error())
	}

	return certToByte(cert), keyToByte(key), nil
}

func NewPKIAssets(o Options, caCert []byte, caKey []byte) (map[string][]byte, error) {
	return nil, nil
}

// CreatePKIAssets will create all pki file(includ etcd)
func CreatePKIAssets(o Options) (map[string][]byte, error) {
	certTree, err := getDefaultCerts().asMap().certTree()
	if err != nil {
		return nil, err
	}
	cfg := &config{
		AdvertiseAddress: o.InternalEndpoint,
		SANs:             o.SANs,
		PublicIP:         o.ExternalEndpoint,
	}
	certs, err := certTree.createTree(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "error creating PKI assets")
	}

	privateKey, err := newPrivateKey()
	if err != nil {
		return nil, err
	}
	certs[pathForKey(constants.ServiceAccountKeyBaseName)] = keyToByte(privateKey)
	pubKeyByte, err := pubKeyToByte(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	certs[pathForPub(constants.ServiceAccountKeyBaseName)] = pubKeyByte

	return certs, nil
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
		return nil, nil, fmt.Errorf("unable to create cert:%s", certSpec.CommonName)
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
		NotAfter:     time.Now().Add(duration10year).UTC(),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  certSpec.Usages,
	}, nil
}
