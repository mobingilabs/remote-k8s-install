package kubeconfig

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"k8s.io/client-go/util/keyutil"

	certutil "k8s.io/client-go/util/cert"

	"k8s.io/client-go/tools/clientcmd"

	"mobingi/ocean/pkg/constants"
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
	// now is the machine's hostname
	Nodes []string

	CaCert []byte
	CaKey  []byte

	// https://45.0.20.1:6443
	ExternalEndpoint string
	InternalEndpoint string

	ClusterName string
}

// NewKubeconfigs return confs
func NewKubeconfigs(o Options) (map[string][]byte, error) {
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
	certs, _ := certutil.ParseCertsPEM(o.CaCert)
	caCert := certs[0]
	key, _ := keyutil.ParsePrivateKeyPEM(o.CaKey)
	caKey := key.(*rsa.PrivateKey)
	specs := map[string]*kubeconfigSpec{
		constants.AdminConf: {
			CACert:     caCert,
			APIServer:  o.ExternalEndpoint,
			ClientName: constants.AdminUser,
			ClientCertAuth: &clientCertAuth{
				CAKey:         caKey,
				Organizations: []string{constants.MastersGroup},
			},
		},
		constants.ControllerManagerConf: {
			CACert:     caCert,
			APIServer:  o.InternalEndpoint,
			ClientName: constants.ControllerManagerUser,
			ClientCertAuth: &clientCertAuth{
				CAKey: caKey,
			},
		},
		constants.SchedulerConf: {
			CACert:     caCert,
			APIServer:  o.InternalEndpoint,
			ClientName: constants.SchedulerUser,
			ClientCertAuth: &clientCertAuth{
				CAKey: caKey,
			},
		},
	}

	for _, v := range o.Nodes {
		specs[v] = &kubeconfigSpec{
			CACert:     caCert,
			APIServer:  o.InternalEndpoint,
			ClientName: fmt.Sprintf("system:node:%s", v),
			ClientCertAuth: &clientCertAuth{
				CAKey:         caKey,
				Organizations: []string{constants.NodesGroup},
			},
		}
	}

	return specs
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
