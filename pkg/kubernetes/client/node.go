package client

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNode() (*v1.NodeList, error) {
	nodes, err := client.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

func DeleteNode(name string) error {
	err := client.CoreV1().Nodes().Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
