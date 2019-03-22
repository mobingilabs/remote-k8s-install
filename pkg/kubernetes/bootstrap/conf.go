package bootstrap

import (
	"errors"
	"fmt"
	pkiutil "mobingi/ocean/pkg/util/pki"

	"k8s.io/client-go/tools/clientcmd"

	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/constants"
	"mobingi/ocean/pkg/tools/cache"
	kubeconfigutil "mobingi/ocean/pkg/util/kubeconfig"
)

// BuildBootstrapKubeletConf push the bootstrap-kubelet.conf to cache
func BuildBootstrapKubeletConf(cfg *config.Config, token string) error {
	masterEndpoint := fmt.Sprintf("https://%s:6443", cfg.Masters[0].PrivateIP)
	caCert, exists := cache.GetOne(constants.CertPrefix, pkiutil.NameForCert(constants.CACertAndKeyBaseName))
	if !exists {
		return errors.New("ca.crt not exist in cache")
	}

	clientConfig := kubeconfigutil.CreateWithToken(masterEndpoint, "kubernetes", constants.NodeBootstrapTokenAuthGroup, caCert.([]byte), token)
	content, err := clientcmd.Write(*clientConfig)
	if err != nil {
		return err
	}

	cache.Put(constants.KubeconfPrefix, constants.BootstrapKubeletConfName, content)

	return nil
}
