package storage

import "mobingi/ocean/pkg/config"

type Database interface {
	NewClient(config interface{}) error
}

type Cluster interface {
	Init(cfg *config.Config) error
	Exist(name string) (bool, error)
	Drop(cfg *config.Config) error
	All() ([]string, error)
	Cert
	Kubeconf
	EtcdServer
	Bootatrap
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

type EtcdServer interface {
	SetEtcdServers(cfg *config.Config) error
	GetEtcdServers(clusterName string) (string, error)
}

type Bootatrap interface {
	SetBootstrapConf(clusterName string) error
}

func NewStorage() Cluster {
	driver := &ClusterMongo{}
	return driver
}
