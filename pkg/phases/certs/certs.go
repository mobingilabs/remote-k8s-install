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

	certutil "k8s.io/client-go/util/cert"
)

const (
	PrivateKeyBlockType = "PRIVATE KEY"
	PublicKeyBlockType  = "PUBLIC KEY"

	CertificateBlockType   = "CERTIFICATE"
	RSAPrivateKeyBlockType = "RSA PRIVATE KEY"

	rsaKeySize     = 2048
	duration10year = 10 * time.Hour * 24 * 365
)

type Options struct {
	// now it is ip
	InternalEndpoint string
	ExternalEndpoint string
	SANs             []string
	ServiceSubnet    string
}

// NewPKIAssets will create all cert and key what we use
func NewPKIAssets(o Options) (map[string][]byte, error) {
	caCert, caKey, err := newRootCA()
	if err != nil {
		return nil, fmt.Errorf("new ca err:%v", err)
	}
	rootCAByte, rootKeyByte := certAndKeyToByte(caCert, caKey)

	certs, err := getCertList(o)
	if err != nil {
		return nil, fmt.Errorf("new pki assests err:%v", err)
	}

	pkiAssets := make(map[string][]byte, len(certs)+1)
	pkiAssets[nameForCert("ca")] = rootCAByte
	pkiAssets[nameForKey("ca")] = rootKeyByte
	for _, v := range certs {
		cert, key, err := newCertAndKeyFromCA(v.config, caCert, caKey)
		if err != nil {
			return nil, fmt.Errorf("create cert and key for %s err:%v", v.name, err)
		}
		caByte, keyByte := certAndKeyToByte(cert, key)

		pkiAssets[nameForCert(v.name)] = caByte
		pkiAssets[nameForKey(v.name)] = keyByte
	}

	priKey, pubKey, err := newServiceAccountKeyPair()
	if err != nil {
		return nil, err
	}
	pkiAssets[nameForKey("sa")] = priKey
	pkiAssets[nameForPub("sa")] = pubKey

	return pkiAssets, nil
}

func newRootCA() (*x509.Certificate, *rsa.PrivateKey, error) {
	key, err := generatePrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create private key:%s", err.Error())
	}

	certCfg := certutil.Config{
		CommonName: "kubernetes",
		Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
	}
	cert, err := certutil.NewSelfSignedCACert(certCfg, key)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create private key:%s", err.Error())
	}

	return cert, key, nil
}

func newServiceAccountKeyPair() ([]byte, []byte, error) {
	privateKey, err := generatePrivateKey()
	if err != nil {
		return nil, nil, err
	}
	priKeyByte := keyToByte(privateKey)
	pubKeyByte, err := pubKeyToByte(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("pub key to byte err:%v", err)
	}

	return priKeyByte, pubKeyByte, nil
}

func newCertAndKeyFromCA(certCfg certutil.Config, caCert *x509.Certificate, caKey *rsa.PrivateKey) (*x509.Certificate, *rsa.PrivateKey, error) {
	key, err := generatePrivateKey()
	if err != nil {
		return nil, nil, fmt.Errorf("new private key err:%v", err)
	}

	cert, err := newSignedCert(certCfg, key, caCert, caKey)
	if err != nil {
		return nil, nil, fmt.Errorf("signed ca err:%v", err)
	}

	return cert, key, nil
}

func generatePrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, rsaKeySize)
}

func newSignedCert(certSpec certutil.Config, key crypto.Signer, caCert *x509.Certificate, caKey crypto.Signer) (*x509.Certificate, error) {
	certTmpl, err := createCertTmpl(&certSpec, caCert.NotBefore)
	if err != nil {
		return nil, err
	}

	certDERBytes, err := x509.CreateCertificate(rand.Reader, certTmpl, caCert, key.Public(), caKey)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(certDERBytes)
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
