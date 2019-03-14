package bootstrap

import (
	"fmt"

	"github.com/pkg/errors"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	bootstrapapi "k8s.io/cluster-bootstrap/token/api"

	"mobingi/ocean/pkg/tools/cache"
)

// CreateBootstrapConfigMapIfNotExists creates the kube-public ConfigMap if it doesn't exist already
func CreateBootstrapConfigMapIfNotExists(client clientset.Interface, cacheKey string) error {
	adminConfigBytes, exists := cache.Get(cacheKey)
	if !exists {
		return fmt.Errorf("can not get admin config from cache,key:%s", cacheKey)
	}

	adminConfig, err := clientcmd.Load(adminConfigBytes.([]byte))
	if err != nil {
		return errors.Wrap(err, "failed to load admin kubeconfig")
	}

	adminCluster := adminConfig.Contexts[adminConfig.CurrentContext].Cluster
	// Copy the cluster from admin.conf to the bootstrap kubeconfig, contains the CA cert and the server URL
	bootstrapConfig := &clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			"": adminConfig.Clusters[adminCluster],
		},
	}
	bootstrapBytes, err := clientcmd.Write(*bootstrapConfig)
	if err != nil {
		return err
	}

	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      bootstrapapi.ConfigMapClusterInfo,
			Namespace: metav1.NamespacePublic,
		},
		Data: map[string]string{
			bootstrapapi.KubeConfigKey: string(bootstrapBytes),
		},
	}
	if _, err := client.CoreV1().ConfigMaps(configMap.ObjectMeta.Namespace).Create(configMap); err != nil {
		return err
	}

	return nil
}
