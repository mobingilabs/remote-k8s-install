package config

import ()

type Config struct {
	PKIDir           string
	AdvertiseAddress string
	ClusterName      string

	Networking Networking
}

func NewConfig() *Config {
	return &Config{
		PKIDir:           "/Users/n1ce/tmp/test/pki",
		AdvertiseAddress: "192.168.0.218",
	}
}

type Networking struct {
	ServiceSubnet string
	PodSubnet     string
	DNSDomain     string
}
