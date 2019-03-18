package bootstrap

import (
	"fmt"
	"mobingi/ocean/pkg/config"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

// Bootstrap new token and post secret to apiserver
// create bootstrap-kubelet.conf to cache
// rolebindings for this token
func Bootstrap(client clientset.Interface, cfg *config.Config) error {
	bt, err := NewBootstrapToken()
	if err != nil {
		return err
	}

	secret := bt.ToSecret()
	if _, err := client.CoreV1().Secrets(metav1.NamespaceSystem).Create(secret); err != nil {
		return fmt.Errorf("can not create secret:%s", err)
	}

	if err := BuildBootstrapKubeletConf(cfg, bt.Token.String()); err != nil {
		return err
	}

	if err := AllowBootstrapTokensToPostCSRs(client); err != nil {
		return err
	}

	if err := AutoApproveNodeBootstrapTokens(client); err != nil {
		return err
	}

	return nil
}
