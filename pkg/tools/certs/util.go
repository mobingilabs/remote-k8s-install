package certs

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"

	"github.com/pkg/errors"

	certutil "k8s.io/client-go/util/cert"
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
	// name := fmt.Sprintf("%s.key", baseName)
	// cache.Put(constants.CertPrefix, name, data)
}

func writeCert(baseName string, data []byte) {
	// name := fmt.Sprintf("%s.crt", baseName)
	// cache.Put(constants.CertPrefix, name, data)
}

func writePub(baseName string, data []byte) {
	// name := fmt.Sprintf("%s.pub", baseName)
	// cache.Put(constants.CertPrefix, name, data)
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

func makeAltNamesMutator(f func(*config) (*certutil.AltNames, error)) configMutatorsFunc {
	return func(cc *certutil.Config, c *config) error {
		altNames, err := f(c)
		if err != nil {
			return err
		}
		cc.AltNames = *altNames
		return nil
	}
}

func getAPIServerAltNames(cfg *config) (*certutil.AltNames, error) {
	advertiseAddress := net.ParseIP(cfg.AdvertiseAddress)
	if advertiseAddress == nil {
		return nil, errors.New("must have advertiseAddress")
	}

	publicIP := net.ParseIP(cfg.PublicIP)
	if publicIP == nil {
		return nil, errors.New("must have public ip")
	}

	altNames := &certutil.AltNames{
		DNSNames: []string{
			"kubernetes",
			"kubernetes.default",
			"kubernetes.default.svc",
			// TODO fix it	fmt.Sprintf("kubernetes.default.svc.%s", cfg.Networking.DNSDomain),
		},
		IPs: []net.IP{
			// TODO fix
			//internalAPIServerVirtualIP,
			advertiseAddress,
			publicIP,
		},
	}

	appendSANsToAltNames(altNames, cfg.SANs)

	//TODO fix up

	return altNames, nil
}

func getEtcdAltNames(cfg *config) (*certutil.AltNames, error) {
	altNames := &certutil.AltNames{}

	appendSANsToAltNames(altNames, cfg.SANs)

	return altNames, nil
}

// appendSANsToAltNames parses SANs from as list of strings and adds them to altNames for use on a specific cert
// altNames is passed in with a pointer, and the struct is modified
// valid IP address strings are parsed and added to altNames.IPs as net.IP's
// RFC-1123 compliant DNS strings are added to altNames.DNSNames as strings
// certNames is used to print user facing warnings and should be the name of the cert the altNames will be used for
func appendSANsToAltNames(altNames *certutil.AltNames, SANs []string) {
	for _, altname := range SANs {
		if ip := net.ParseIP(altname); ip != nil {
			altNames.IPs = append(altNames.IPs, ip)
		} // TODO fix else if len(validation.NameIsDNSSubdomain(altname)) == 0 {
		altNames.DNSNames = append(altNames.DNSNames, altname)
	} /* else {
			fmt.Printf(
				"[certificates] WARNING: '%s' was not added to the SAN, because it is not a valid IP or RFC-1123 compliant DNS entry\n",
				altname,
			)
		}
	}*/
}
