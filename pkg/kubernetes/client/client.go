package client

import (
	"mobingi/ocean/pkg/storage"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClient(cluster string) (*kubernetes.Clientset, error) {
	storage := storage.NewStorage()
	kubeconfig, err := storage.GetKubeconf(cluster, "admin.conf")
	if err != nil {
		return nil, err
	}
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
