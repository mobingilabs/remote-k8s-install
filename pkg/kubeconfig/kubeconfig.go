package kubeconfig

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"mobingi/ocean/pkg/constants"
	"path/filepath"

	"github.com/pkg/errors"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	certutil "k8s.io/client-go/util/cert"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/ssh"
	"mobingi/ocean/pkg/tools/cache"
	cmdutil "mobingi/ocean/pkg/util/cmd"
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

func CreateKubeconfigFiles(c ssh.Client, cfg *config.Config) error {
	specs, err := getKubeconfigSpecs(c, cfg)
	if err != nil {
		return err
	}

	//TODO other location just mkdir all once
	cmd := cmdutil.NewMkdirAllCmd(constants.WorkDir)
	_, err = c.Do(cmd)
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

		if err := writeKubeconfigFile(c, constants.WorkDir, kubeconfigFileName, config); err != nil {
			return errors.Wrap(err, "kubeconfig")
		}
	}

	return nil
}

func getKubeconfigSpecs(c ssh.Client, cfg *config.Config) (map[string]*kubeconfigSpec, error) {
	// TODO load from cache
	caCert, caKey, err := pkiutil.TryLoadCertAndKeyFromDisk(c, constants.PKIDir, kubeadmconstants.CACertAndKeyBaseName)
	if err != nil {
		return nil, err
	}

	// TODO port from config
	publicEndpoint := fmt.Sprintf("https://%s:6443", cfg.AdvertiseAddress)
	privateEndpoint := fmt.Sprintf("https://%s:6443", cfg.Masters[0].PrivateIP)

	var kubeconfigSepcs = map[string]*kubeconfigSpec{
		kubeadmconstants.AdminKubeConfigFileName: {
			CACert:     caCert,
			APIServer:  publicEndpoint,
			ClientName: "kubernetes-admin",
			ClientCertAuth: &clientCertAuth{
				CAKey:         caKey,
				Organizations: []string{kubeadmconstants.MastersGroup},
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

func writeKubeconfigFile(c ssh.Client, dir, name string, cfg *clientcmdapi.Config) error {
	content, err := clientcmd.Write(*cfg)
	if err != nil {
		return err
	}

	filename := filepath.Join(dir, name)
	cache.Put(name, content)

	cmd := cmdutil.NewWriteCmd(filename, string(content))
	_, err = c.Do(cmd)
	if err != nil {
		return err
	}

	return nil
}
