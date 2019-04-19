package kubeconf

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"mobingi/ocean/pkg/constants"

	"github.com/pkg/errors"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	certutil "k8s.io/client-go/util/cert"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"

	"mobingi/ocean/pkg/config"
	kubeconfigutil "mobingi/ocean/pkg/util/kubeconfig"
	pkiutil "mobingi/ocean/pkg/util/pki"
)

type kubeconfigSpec struct {
	CACert         *x509.Certificate
	APIServer      string
	ClientName     string
	ClientCertAuth *clientCertAuth
}

type clientCertAuth struct {
	CAKey         *rsa.PrivateKey
	Organizations []string
}

func CreateKubeconf(cfg *config.Config, caCert *x509.Certificate, caKey *rsa.PrivateKey) (map[string][]byte, error) {
	specs, err := getKubeconfigSpecs(cfg, caCert, caKey)
	if err != nil {
		return nil, err
	}

	kubeconfigFileNames := []string{
		kubeadmconstants.AdminKubeConfigFileName,
		kubeadmconstants.ControllerManagerKubeConfigFileName,
		kubeadmconstants.SchedulerKubeConfigFileName,
	}

	kubeconfs := make(map[string][]byte)

	for _, kubeconfigFileName := range kubeconfigFileNames {
		spec, exists := specs[kubeconfigFileName]
		if !exists {
			return nil, errors.Errorf("could't retrive kubeconfigSpec for %s", kubeconfigFileName)
		}

		config, err := buildKubeconfigFromSpec(spec, cfg.ClusterName)
		if err != nil {
			return nil, err
		}
		content, err := clientcmd.Write(*config)
		kubeconfs[kubeconfigFileName] = content
	}

	return kubeconfs, nil
}

func getKubeconfigSpecs(cfg *config.Config, caCert *x509.Certificate, caKey *rsa.PrivateKey) (map[string]*kubeconfigSpec, error) {
	// TODO get port from config
	publicEndpoint := fmt.Sprintf("https://%s:6443", cfg.PublicIP)
	privateEndpoint := fmt.Sprintf("https://%s:6443", cfg.AdvertiseAddress)
	var kubeconfigSepcs = map[string]*kubeconfigSpec{
		kubeadmconstants.AdminKubeConfigFileName: {
			CACert:     caCert,
			APIServer:  publicEndpoint,
			ClientName: "kubernetes-admin",
			ClientCertAuth: &clientCertAuth{
				CAKey:         caKey,
				Organizations: []string{constants.MastersGroup},
			},
		},
		kubeadmconstants.ControllerManagerKubeConfigFileName: {
			CACert:     caCert,
			APIServer:  privateEndpoint,
			ClientName: kubeadmconstants.ControllerManagerUser,
			ClientCertAuth: &clientCertAuth{
				CAKey: caKey,
			},
		},
		kubeadmconstants.SchedulerKubeConfigFileName: {
			CACert:     caCert,
			APIServer:  privateEndpoint,
			ClientName: kubeadmconstants.SchedulerUser,
			ClientCertAuth: &clientCertAuth{
				CAKey: caKey,
			},
		},
	}

	return kubeconfigSepcs, nil
}

func buildKubeconfigFromSpec(spec *kubeconfigSpec, clusterName string) (*clientcmdapi.Config, error) {
	clientCertConfig := certutil.Config{
		CommonName:   spec.ClientName,
		Organization: spec.ClientCertAuth.Organizations,
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	clientCert, clientKey, err := pkiutil.NewCertAndKeyFromCA(spec.CACert, spec.ClientCertAuth.CAKey, &clientCertConfig)
	if err != nil {
		return nil, err
	}

	return kubeconfigutil.CreateWithCerts(
		spec.APIServer,
		clusterName,
		spec.ClientName,
		pkiutil.EncodeCertPEM(spec.CACert),
		pkiutil.EncodePrivateKeyPEM(clientKey),
		pkiutil.EncodeCertPEM(clientCert),
	), nil

	return nil, nil
}
