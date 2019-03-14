package bootstrap

import (
	"errors"
	"fmt"

	"k8s.io/client-go/tools/clientcmd"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/cache"
	kubeconfigutil "mobingi/ocean/pkg/util/kubeconfig"
)

// BuildBootstrapKubeletConf push the bootstrap-kubelet.conf to cache
func BuildBootstrapKubeletConf(cfg *config.Config, token string) error {
	masterEndpoint := fmt.Sprintf("https://%s:6443", cfg.AdvertiseAddress)
	caCert, exists := cache.Get("ca.crt")
	if !exists {
		return errors.New("ca.crt not exist in cache")
	}

	clientConfig := kubeconfigutil.CreateWithToken(masterEndpoint, "kubernetes", constants.NodeBootstrapTokenAuthGroup, caCert.([]byte), token)

	content, err := clientcmd.Write(*clientConfig)
	if err != nil {
		return err
	}

	cache.Put("bootstrap-kubelet.conf", content)
	fmt.Println(string(content))

	return nil
}
