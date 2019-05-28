package service

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func WaitingForApiserverStart(conf []byte) error {
	config, err := clientcmd.RESTConfigFromKubeConfig(conf)
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
retry:
	_, err = clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		time.Sleep(1 * time.Second)
		goto retry
	}
	return nil
}
