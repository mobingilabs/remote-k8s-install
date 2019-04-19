package main

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	c, err := clientcmd.LoadFromFile("admin.conf")
	if err != nil {
		panic(err)
	}

	clientConfig, err := clientcmd.NewDefaultClientConfig(*c, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		panic(err)
	}

	client, err := clientset.NewForConfig(clientConfig)
	if err != nil {
		panic(err)
	}

	_, err = client.CoreV1().Pods("default").Get("xx", metav1.GetOptions{})
	if err != nil {
		fmt.Println("client  get pods err:%s", err.Error())
	}
}
