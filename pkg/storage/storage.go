package storage

import "mobingi/ocean/pkg/config"

type Database interface {
	NewClient(config interface{}) error
}

type Cluster interface {
	Init(cfg *config.Config) error
	Cert
	Kubeconf
}

type Cert interface {
	CreateCerts(cfg *config.Config) error
	AllCerts(clusterName string) (map[string][]byte, error)
	GetCert(clusterName, name string) ([]byte, error)
}

type Kubeconf interface {
	CreateKubeconfs(cfg *config.Config, caCert []byte, caKey []byte) error
	AllKubeconfs(clusterName string) (map[string][]byte, error)
	GetKubeconf(clusterName, name string) ([]byte, error)
}

func NewStorage(driver Cluster, cfg *config.Config) (Cluster, error) {
	driver.Init(cfg)
	return driver, nil
}
