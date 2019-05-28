package manifests

import (
	"io/ioutil"
	"testing"
)

const ip = "192.168.50.173"

func TestGetStaticPodMainfests(t *testing.T) {
	o := Options{
		IP:        ip,
		NodeName:  "node0",
		EtcdToken: "token",
		Companions: map[string]string{
			"node1": "192.168.1.2",
			"node2": "192.168.1.3",
		},
		EtcdImage:              "k8s.gcr.io/etcd:3.3.10",
		APIServerImage:         "k8s.gcr.io/kube-apiserver:v1.14.2",
		ControllerManagerImage: "k8s.gcr.io/kube-controller-manager:v1.14.2",
		SchedulerImage:         "k8s.gcr.io/kube-scheduler:v1.14.2",
		ServiceIPRange:         "10.96.0.0/12",
	}
	mainfests := GetStaticPodMainfests(o)
	for k, v := range mainfests {
		ioutil.WriteFile(k, v, 0400)
	}
}
