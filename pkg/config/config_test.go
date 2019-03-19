package config

import (
	"log"
	"testing"
)

func TestLoadFromcfgFromFile(t *testing.T) {
	cfg, err := LoadConfigFromFile("testdata/test_config.yaml")
	if err != nil {
		log.Fatalf("cfg file parse err:%s\n", err.Error())
	}

	if cfg.ClusterName != "kubernetes" {
		log.Fatalf("cluster name parse failed,parsed to:%s\n", cfg.ClusterName)
	}

	if cfg.AdvertiseAddress != "192.168.1.7" {
		log.Fatalf("cluster name parse failed,parsed to:%s\n", cfg.AdvertiseAddress)
	}

	if cfg.Masters[0].PublicIP != "47.52.227.32" || cfg.Masters[0].PrivateIP != "192.168.1.1" ||
		cfg.Masters[0].User != "root" || cfg.Masters[0].Password != "312313" {
		log.Fatalf("cfg masters parsed failed,publicIP:%s,privateIP:%s,user:%s,password:%s",
			cfg.Masters[0].PublicIP, cfg.Masters[0].PrivateIP, cfg.Masters[0].User, cfg.Masters[0].Password)
	}

	if cfg.Nodes[0].PublicIP != "47.52.227.33" || cfg.Nodes[0].PrivateIP != "192.168.1.2" ||
		cfg.Nodes[0].User != "root" || cfg.Nodes[0].Password != "312313" {
		log.Fatalf("cfg masters parsed failed,publicIP:%s,privateIP:%s,user:%s,password:%s",
			cfg.Nodes[0].PublicIP, cfg.Nodes[0].PrivateIP, cfg.Nodes[0].User, cfg.Nodes[0].Password)
	}
}
