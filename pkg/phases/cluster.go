package phases

import (
	"strings"
	"fmt"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/phases/certs"
	"mobingi/ocean/pkg/phases/manifests"
	configstorage "mobingi/ocean/pkg/storage"
)

func Init(cfg *config.Config) (configstorage.Cluster, error) {
	storage := configstorage.NewStorage()
	err := storage.Init(cfg)
	if err != nil {
		return nil, err
	}
	return storage, nil
}

// PrepareMaster generate all files we need
func PrepareMaster(cfg *config.Config) {
}

func newPKIs(cfg config.Config) (map[string][]byte, error) {
	sans := make([]string, 0, len(cfg.Masters))
	for _, v := range cfg.Masters {
		sans = append(sans, v.PublicIP)
	}
	o := certs.Options{
		InternalEndpoint: cfg.AdvertiseAddress,
		ExternalEndpoint: cfg.PublicIP,
		SANs:             sans,
		ServiceSubnet:    "10.96.0.0/12",
	}
	return certs.NewPKIAssets(o)
}

func newManifests(cfg config.Config) (map[string][string][]byte, error) {
	nodeCnt := len(cfg.Masters)
	data := make(map[string]map[string][]byte, 0, len(cfg.Masters))

	companions := make(map[string]map[string]string, 0, len(cfg.Masters))
	for i:=0; i<len(nodeCnt); i++ {
		companions[getNodeName(i)] = make(map[string]string, 0, nodeCnt-1)
	}
	for i, m := range cfg.Masters {
		for k, v := range companions {
			if strings.HasSuffix(k, m.PrivateIP) {
				break	
			}
		}
	}
	o := manifests.Options{
		IP:                     v.PrivateIP,
		EtcdNodeName:           fmt.Sprintf("node:%d", i),
		EtcdToken:              "token",
		EtcdImage:              "n1ce37/etcd:3.3.10",
		APIServerImage:         "n1ce37/kube-apiserver:v1.14.2",
		ControllerManagerImage: "n1ce37/kube-controller-manager:v1.14.2",
		SchedulerImage:         "n1ce37/kube-scheduler:v1.14.2",
		ServiceIPRange:         "10.96.0.0/12",
	}

	for i, v := range cfg.Masters {

		files := manifests.NewStaticPodManifests(o)
	}
	return nil, nil
}

func RunMaster() {
}

// node name is node:{privateip}
func getNodeName(ip string) string{
	return fmt.Sprintf("node:%s", string)
}
