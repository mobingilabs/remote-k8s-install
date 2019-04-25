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
	CreateCerts() error
	AllCerts() (map[string][]byte, error)
	GetCert(name string) ([]byte, error)
}

type Kubeconf interface {
	CreateKubeconfs(caCert []byte, caKey []byte) error
	AllKubeconfs() (map[string][]byte, error)
	GetKubeconf(name string) ([]byte, error)
}

func NewStorage(driver Cluster, cfg *config.Config) (*Cluster, error) {
	driver.Init(cfg)
	return &driver, nil
}
