package bootstrap

import (
	"k8s.io/client-go/tools/clientcmd"

	"mobingi/ocean/pkg/constants"
	kubeconfigutil "mobingi/ocean/pkg/util/kubeconfig"
)

// BuildBootstrapKubeletConf push the bootstrap-kubelet.conf to cache
func BuildBootstrapKubeletConf(apiserverURL string, token string, caCrt []byte) ([]byte, error) {
	clientConfig := kubeconfigutil.CreateWithToken(apiserverURL, "kubernetes", constants.NodeBootstrapTokenAuthGroup, caCrt, token)
	content, err := clientcmd.Write(*clientConfig)
	if err != nil {
		return nil, err
	}

	return content, nil
}
