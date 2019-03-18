package main

import (
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	client, err := NewK8SClientFromConf()
	if err != nil {
		panic(err)
	}

	resp, err := client.CoreV1().Secrets("kube-system").List(v1.ListOptions{})
	fmt.Println(resp.Items, err)
}

func NewK8SClientFromConf() (clientset.Interface, error) {
	config, err := clientcmd.LoadFromFile("admin.conf")
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
