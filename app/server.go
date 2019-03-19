package app

import (
	"mobingi/ocean/app/phases/master"
	"mobingi/ocean/app/phases/node"
	"mobingi/ocean/pkg/config"
	"mobingi/ocean/pkg/log"
)

func Start() error {
	cfg, err := config.LoadConfigFromFile("config.yaml")
	if err != nil {
		log.Error(err)
		return err
	}
	if err := master.Start(cfg); err != nil {
		log.Error(err)
		return err
	}

	if err := node.Start(cfg); err != nil {
		return err
	}

	return nil
}
