package app

import (
	"mobingi/ocean/app/phases/master"
	"mobingi/ocean/app/phases/node"
	"mobingi/ocean/pkg/config"
)

func Start() error {
	cfg, err := config.LoadConfigFromFile("config.yml")
	if err != nil {
		return err
	}
	if err := master.Start(cfg); err != nil {
		return err
	}

	if err := node.Start(cfg); err != nil {
		return err
	}

	return nil
}
