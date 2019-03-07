package kubeconfig

import (
	"crypto/rsa"
	"crypto/x509"
	"path/filepath"

	"github.com/pkg/errors"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	certutil "k8s.io/client-go/util/cert"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"

	"mobingi/ocean/pkg/config"
	kubeconfigutil "mobingi/ocean/pkg/util/kubeconfig"
	"mobingi/ocean/pkg/util/pki"
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

func CreateKubeconfigFiles(cfg *config.Config) error {
	specs, err := getKubeconfigSpecs(cfg)
	if err != nil {
		return err
	}

	kubeconfigFileNames := []string{
		kubeadmconstants.AdminKubeConfigFileName,
		kubeadmconstants.ControllerManagerKubeConfigFileName,
		kubeadmconstants.SchedulerKubeConfigFileName,
	}

	for _, kubeconfigFileName := range kubeconfigFileNames {
		spec, exists := specs[kubeconfigFileName]
		if !exists {
			return errors.Errorf("could't retrive kubeconfigSpec for %s", kubeconfigFileName)
		}

		config, err := buildKubeconfigFromSpec(spec, cfg.ClusterName)
		if err != nil {
			return err
		}

		if err := writeKubeconfigFile(cfg.PKIDir, kubeconfigFileName, config); err != nil {
			return err
		}
	}

	return nil
}

func getKubeconfigSpecs(cfg *config.Config) (map[string]*kubeconfigSpec, error) {
	caCert, caKey, err := pki.TryLoadCertAndKeyFromDisk(cfg.PKIDir, kubeadmconstants.CACertAndKeyBaseName)
	if err != nil {
		return nil, err
	}

	// TODO read it from config
	masterEndpoint := "https://192.168.0.218:6443"

	var kubeconfigSepcs = map[string]*kubeconfigSpec{
		kubeadmconstants.AdminKubeConfigFileName: {
			CACert:     caCert,
			APIServer:  masterEndpoint,
			ClientName: "kubernetes-admin",
			ClientCertAuth: &clientCertAuth{
				CAKey:         caKey,
				Organizations: []string{kubeadmconstants.MastersGroup},
			},
		},
		kubeadmconstants.ControllerManagerKubeConfigFileName: {
			CACert:     caCert,
			APIServer:  masterEndpoint,
			ClientName: kubeadmconstants.ControllerManagerUser,
			ClientCertAuth: &clientCertAuth{
				CAKey: caKey,
			},
		},
		kubeadmconstants.SchedulerKubeConfigFileName: {
			CACert:     caCert,
			APIServer:  masterEndpoint,
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
	clientCert, clientKey, err := pki.NewCertAndKeyFromCA(spec.CACert, spec.ClientCertAuth.CAKey, &clientCertConfig)
	if err != nil {
		return nil, err
	}

	return kubeconfigutil.CreateWithCerts(
		spec.APIServer,
		clusterName,
		spec.ClientName,
		pki.EncodeCertPEM(spec.CACert),
		pki.EncodePrivateKeyPEM(clientKey),
		pki.EncodeCertPEM(clientCert),
	), nil

	return nil, nil
}

func writeKubeconfigFile(dir, name string, cfg *clientcmdapi.Config) error {
	kubeconfigFilepath := filepath.Join(dir, name)
	return clientcmd.WriteToFile(*cfg, kubeconfigFilepath)
}
