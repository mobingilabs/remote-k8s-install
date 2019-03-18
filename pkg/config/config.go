package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	workDir = "/etc/kubernetes"
	pkiDir  = "/etc/kubernetes/pki"
)

type MachineRole int

type Config struct {
	ClusterName      string
	AdvertiseAddress string

	PKIDir  string
	WorkDir string

	Masters []Machine
	Nodes   []Machine
}

// TODO more ssh auth method support
type Machine struct {
	Addr     string
	User     string
	Password string
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

	// ? config it or const set
	conf.PKIDir = pkiDir
	conf.WorkDir = workDir

	return conf, nil
}
