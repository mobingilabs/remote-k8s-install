package config

import (
	"mobingi/ocean/app/proto"
)

type Config struct {
	ClusterName      string `yaml:"clusterName"`
	AdvertiseAddress string `yaml:"advertiseAddress"`
	DownloadBinSite  string `yaml:"downloadBinSite"`
	PublicIP         string `yaml:"publicIP"`

	Masters []Machine
	Nodes   []Machine
}

func (c *Config) GetSANs() []string {
	sans := make([]string, 0, len(c.Masters))
	for _, v := range c.Masters {
		sans = append(sans, v.PrivateIP)
	}

	return sans
}

func (c *Config) GetMasterPrivateIPs() []string {
	privateIPs := make([]string, 0, len(c.Masters))
	for _, v := range c.Masters {
		privateIPs = append(privateIPs, v.PrivateIP)
	}

	return privateIPs
}

// TODO more ssh auth method support
type Machine struct {
	PublicIP  string `yaml:"publicIP"`
	PrivateIP string `yaml:"privateIP"`
	User      string
	Password  string
}

func LoadConfigFromGrpc(clusterCfg *proto.ClusterConfig) (*Config, error) {
	var masters []Machine
	for _, v := range clusterCfg.Masters {
		machine := Machine{
			PublicIP:  v.PublicIP,
			PrivateIP: v.PrivateIP,
			User:      v.User,
			Password:  v.Password,
		}
		masters = append(masters, machine)
	}
	cfg := &Config{
		ClusterName:      clusterCfg.ClusterName,
		AdvertiseAddress: clusterCfg.AdvertiseAddress,
		DownloadBinSite:  clusterCfg.DownloadBinSite,
		PublicIP:         clusterCfg.PublicIP,
		Masters:          []Machine(masters),
	}
	return cfg, nil
}
