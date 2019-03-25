package bootstrap

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	kubeconfigutil "mobingi/ocean/pkg/util/kubeconfig"
)

// BuildBootstrapKubeletConf push the bootstrap-kubelet.conf to cache
func BuildBootstrapKubeletConf(cfg *config.Config, token string, caCrt []byte) ([]byte, error) {
	masterEndpoint := fmt.Sprintf("https://%s:6443", cfg.Masters[0].PrivateIP)

	clientConfig := kubeconfigutil.CreateWithToken(masterEndpoint, "kubernetes", constants.NodeBootstrapTokenAuthGroup, caCrt, token)
	content, err := clientcmd.Write(*clientConfig)
	if err != nil {
		return nil, err
	}

	return content, nil
}
