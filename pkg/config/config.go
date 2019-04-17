package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ClusterName      string `yaml:"clusterName"`
	AdvertiseAddress string `yaml:"advertiseAddress"`
	DownloadBinSite  string `yaml:"downloadBinSite"`

	Masters []Machine
	Nodes   []Machine
}


// TODO more ssh auth method support
type Machine struct {
	PublicIP  string `yaml:"publicIP"`
	PrivateIP string `yaml:"privateIP"`
	User      string
	Password  string
}

func LoadConfigFromFile(name string) (*Config, error) {
	confByte, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	conf := &Config{}

	if err := yaml.Unmarshal(confByte, conf); err != nil {
		return nil, err
	}

	return conf, nil
}
