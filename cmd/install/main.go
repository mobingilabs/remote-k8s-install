package main

import (
	"io/ioutil"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/tools/certs"
	"mobingi/ocean/pkg/tools/kubeconf"
	"mobingi/ocean/pkg/util/pki"
)

const ip = "192.168.50.173"

func main() {
	Init()
}

func Init() {
	certs, err := certs.CreatePKIAssets(ip, ip, nil)
	if err != nil {
		panic(err)
	}
	for k, v := range certs {
		ioutil.WriteFile(k, v, 0444)
	}

	cfg := config.Config{
		AdvertiseAddress: ip,
		PublicIP:         ip,
	}

	cert, _ := pki.ParseCertPEM(certs["ca.crt"])
	key, _ := pki.ParsePrivateKeyPEM(certs["ca.key"])

	confs, err := kubeconf.CreateKubeconf(&cfg, cert, key)
	if err != nil {
		panic(err)
	}
	for k, v := range confs {
		ioutil.WriteFile(k, v, 0444)
	}
}
