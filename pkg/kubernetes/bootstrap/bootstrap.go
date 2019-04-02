package bootstrap

import (
	"fmt"
	"mobingi/ocean/pkg/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"

	"k8s.io/client-go/tools/clientcmd"
)

// Bootstrap new token and post secret to apiserver
// create bootstrap-kubelet.conf to cache
// rolebindings for this token
func Bootstrap(adminconf []byte, cfg *config.Config, caCrt []byte) ([]byte, error) {
	client, err := newK8sClientFromConf(adminconf)
	if err != nil {
		return nil, err
	}

	bt, err := NewBootstrapToken()
	if err != nil {
		return nil, err
	}

	secret := bt.ToSecret()
	if _, err := client.CoreV1().Secrets(metav1.NamespaceSystem).Create(secret); err != nil {
		return nil, fmt.Errorf("can not create secret:%s", err)
	}

	bootstrapConf, err := BuildBootstrapKubeletConf(cfg, bt.Token.String(), caCrt)
	if err != nil {
		return nil, err
	}

	if err := AllowBootstrapTokensToPostCSRs(client); err != nil {
		return nil, err
	}

	if err := AutoApproveNodeBootstrapTokens(client); err != nil {
		return nil, err
	}

	return bootstrapConf, nil
}

func newK8sClientFromConf(conf []byte) (clientset.Interface, error) {
	config, err := clientcmd.Load(conf)
	if err != nil {
		return nil, err
	}

	clientConfig, err := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, err
	}

	client, err := clientset.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}
