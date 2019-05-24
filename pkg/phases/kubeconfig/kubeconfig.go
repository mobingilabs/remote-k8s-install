package kubeconfig

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"mobingi/ocean/pkg/constants"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	certutil "k8s.io/client-go/util/cert"

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

	ClusterName string
}

func CreateKubeconf(o Options) (map[string][]byte, error) {
	specs := getSpecs(o)

	kubeconfigs := make(map[string][]byte, len(specs))

	for k, v := range specs {
		kubeconfig, err := buildKubeconfigFromSpec(v, o.ClusterName)
		if err != nil {
			return nil, fmt.Errorf("create kubeconfigs:%s err:%v", k, v)
		}
		kubeconfigs[k] = kubeconfig
	}

	return kubeconfigs, nil
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

func buildKubeconfigFromSpec(spec *kubeconfigSpec, clusterName string) ([]byte, error) {
	clientCertConfig := certutil.Config{
		CommonName:   spec.ClientName,
		Organization: spec.ClientCertAuth.Organizations,
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	clientCert, clientKey, err := pkiutil.NewCertAndKeyFromCA(spec.CACert, spec.ClientCertAuth.CAKey, &clientCertConfig)
	if err != nil {
		return nil, err
	}

	config := createWithCerts(
		spec.APIServer,
		clusterName,
		spec.ClientName,
		pkiutil.EncodeCertPEM(spec.CACert),
		pkiutil.EncodePrivateKeyPEM(clientKey),
		pkiutil.EncodeCertPEM(clientCert),
	)

	return clientcmd.Write(*config)
}

func createWithCerts(serverURL, clusterName, userName string, caCert []byte, clientKey []byte, clientCert []byte) *clientcmdapi.Config {
	config := CreateBasic(serverURL, clusterName, userName, caCert)
	config.AuthInfos[userName] = &clientcmdapi.AuthInfo{
		ClientKeyData:         clientKey,
		ClientCertificateData: clientCert,
	}
	return config
}

func CreateBasic(serverURL, clusterName, userName string, caCert []byte) *clientcmdapi.Config {
	// Use the cluster and the username as the context name
	contextName := fmt.Sprintf("%s@%s", userName, clusterName)

	return &clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			clusterName: {
				Server:                   serverURL,
				CertificateAuthorityData: caCert,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			contextName: {
				Cluster:  clusterName,
				AuthInfo: userName,
			},
		},
		AuthInfos:      map[string]*clientcmdapi.AuthInfo{},
		CurrentContext: contextName,
	}
}
