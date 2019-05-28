package certs

import (
	"crypto/x509"
	"fmt"
	"net"

	"mobingi/ocean/pkg/constants"

	certutil "k8s.io/client-go/util/cert"
	"k8s.io/kubernetes/pkg/registry/core/service/ipallocator"
)

type cert struct {
	name   string
	config certutil.Config
}

func newAPIServerCert(o Options) (*cert, error) {
	internalAddress := net.ParseIP(o.InternalEndpoint)
	if internalAddress == nil {
		return nil, fmt.Errorf("unable to parse internal endpoint:%q", o.InternalEndpoint)
	}
	externalAddress := net.ParseIP(o.ExternalEndpoint)
	if externalAddress == nil {
		return nil, fmt.Errorf("unable to parse external endpoint:%q", o.ExternalEndpoint)
	}

	_, svcSubnet, err := net.ParseCIDR(o.ServiceSubnet)
	if err != nil {
		return nil, fmt.Errorf("parse service subnet %q error:%v", svcSubnet, err)
	}

	apiserverVirutalIP, err := ipallocator.GetIndexedIP(svcSubnet, 1)
	if err != nil {
		return nil, fmt.Errorf("unable to get first ip address")
	}

	return &cert{
		name: "apiserver",
		config: certutil.Config{
			CommonName: constants.APIServerCertCommonName,
			Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			AltNames: certutil.AltNames{
				DNSNames: []string{
					"kubernetes",
					"kubernetes.default",
				},
				IPs: []net.IP{
					apiserverVirutalIP,
					internalAddress,
					externalAddress,
				},
			},
		},
	}, nil
}

func newAPIServerKubeletClientCert() *cert {
	return &cert{
		name: "apiserver-kubelet-client",
		config: certutil.Config{
			CommonName:   constants.APIServerKubeletClientCertCommonName,
			Organization: []string{constants.MastersGroup},
			Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
	}
}

func newAPIServerEtcdClientCert() *cert {
	return &cert{
		name: "apiserver-etcd-client",
		config: certutil.Config{
			CommonName:   constants.APIServerEtcdClientCertCommonName,
			Organization: []string{constants.MastersGroup},
			Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
	}
}

func newEtcdServerCert(o Options) *cert {
	ips := make([]net.IP, 0, len(o.SANs))
	for _, v := range o.SANs {
		ips = append(ips, net.ParseIP(v))
	}
	return &cert{
		name: "etcd-server",
		config: certutil.Config{
			Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
			AltNames: certutil.AltNames{
				IPs: ips,
			},
		},
	}
}

func newEtcdPeerCert(o Options) *cert {
	ips := make([]net.IP, 0, len(o.SANs))
	for _, v := range o.SANs {
		ips = append(ips, net.ParseIP(v))
	}

	return &cert{
		name: "etcd-peer",
		config: certutil.Config{
			Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
			AltNames: certutil.AltNames{
				IPs: ips,
			},
		},
	}
}

func getCertList(o Options) ([]*cert, error) {
	apiserverCert, err := newAPIServerCert(o)
	if err != nil {
		return nil, err
	}

	return []*cert{
		apiserverCert,
		newAPIServerKubeletClientCert(),
		newAPIServerEtcdClientCert(),
		newEtcdServerCert(o),
		newEtcdPeerCert(o),
	}, nil
}
