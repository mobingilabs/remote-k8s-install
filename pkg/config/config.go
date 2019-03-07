package config

type Config struct {
	PKIDir                string
	AdvertiseAddress      string
	ServiceClusterIPRange string
}

func NewConfig() *Config {
	return &Config{
		PKIDir:                   "/etc/kubernetes/pki",
		AdvertiseAddress:      "192.168.0.218",
		ServiceClusterIPRange: "10.0.0.1",
	}
}
