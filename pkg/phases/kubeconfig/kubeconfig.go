package kubeconfig

import (
	"crypto/rsa"
	"crypto/x509"
	"mobingi/ocean/pkg/constants"

	"github.com/pkg/errors"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	certutil "k8s.io/client-go/util/cert"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"

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

type Options struct {
	CaCert *x509.Certificate
	CaKey  *rsa.PrivateKey

	// https://45.0.20.1:6443
	ExternalEndpoint string
	InternalEndpoint string
}

func CreateKubeconf(o Options) (map[string][]byte, error) {
	specs, err := getKubeconfigSpecs(cfg, caCert, caKey)
	if err != nil {
		return nil, err
	}

	kubeconfigFileNames := []string{
		kubeadmconstants.AdminKubeConfigFileName,
		kubeadmconstants.ControllerManagerKubeConfigFileName,
		kubeadmconstants.SchedulerKubeConfigFileName,
		"kubelet.conf",
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

func getSpecs(o Options) map[string]*kubeconfigSpec {
	return map[string]*kubeconfigSpec{
		constants.AdminConf: {
			CACert:     o.CaCert,
			APIServer:  o.ExternalEndpoint,
			ClientName: constants.AdminUser,
			ClientCertAuth: &clientCertAuth{
				CAKey:         o.CaKey,
				Organizations: []string{constants.MastersGroup},
			},
		},
		constants.ControllerManagerConf: {
			CACert:     o.CaCert,
			APIServer:  o.InternalEndpoint,
			ClientName: constants.ControllerManagerUser,
			ClientCertAuth: &clientCertAuth{
				CAKey: o.CaKey,
			},
		},
		constants.SchedulerConf: {
			CACert:     o.CaCert,
			APIServer:  o.InternalEndpoint,
			ClientName: constants.SchedulerUser,
			ClientCertAuth: &clientCertAuth{
				CAKey: o.CaKey,
			},
		},
	}
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
