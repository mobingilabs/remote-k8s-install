package main

import (
	"fmt"
	"io/ioutil"
	"mobingi/ocean/pkg/phases/certs"
	"mobingi/ocean/pkg/phases/kubeconfig"
	"mobingi/ocean/pkg/phases/manifests"
	"path/filepath"
)

const ip = "172.17.81.61"

func main() {
	Init()
}

func bootstra() {
	caCert, _ := ioutil.ReadFile("tmp/pki/ca.crt")
	caKey, _ := ioutil.ReadFile("tmp/pki/ca.key")
	writeConfs(caCert, caKey)
}

func Init() {
	certs := writeCerts()
	writeConfs(certs["ca.crt"], certs["ca.key"])
	writeManifests()
}

func writeCerts() map[string][]byte {
	dir := "pki"
	certOpt := certs.Options{
		InternalEndpoint: ip,
		ExternalEndpoint: ip,
		ServiceSubnet:    "10.96.0.0/12",
		SANs:             []string{ip},
	}
	certs, err := certs.NewPKIAssets(certOpt)
	if err != nil {
		panic(err)
	}
	for k, v := range certs {
		writeFile(k, dir, v)
	}

	return certs
}

func writeConfs(ca []byte, key []byte) {
	dir := "conf"
	confOpt := kubeconfig.Options{
		CaCert:           ca,
		CaKey:            key,
		ExternalEndpoint: fmt.Sprintf("https://%s:6443", ip),
		InternalEndpoint: fmt.Sprintf("https://%s:6443", ip),
		ClusterName:      "kubernetes",
	}
	confs, err := kubeconfig.NewKubeconf(confOpt)
	if err != nil {
		panic(err)
	}

	for k, v := range confs {
		writeFile(k, dir, v)
	}
}

func writeManifests() {
	dir := "manifests"
	manifestsOpt := manifests.Options{
		IP:                     ip,
		NodeName:               "test",
		EtcdToken:              "token",
		EtcdImage:              "n1ce37/etcd:3.3.10",
		APIServerImage:         "n1ce37/kube-apiserver:v1.14.2",
		ControllerManagerImage: "n1ce37/kube-controller-manager:v1.14.2",
		SchedulerImage:         "n1ce37/kube-scheduler:v1.14.2",
		ServiceIPRange:         "10.96.0.0/12",
	}
	files := manifests.NewStaticPodManifests(manifestsOpt)
	for k, v := range files {
		writeFile(k, dir, v)
	}
}

func writeFile(name, dir string, data []byte) {
	ioutil.WriteFile(filepath.Join("tmp", dir, name), data, 0444)
}
