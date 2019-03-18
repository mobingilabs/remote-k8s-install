package main

import (
	"mobingi/ocean/app/phases"
	"mobingi/ocean/pkg/config"
)

func main() {
	cfg, err := config.LoadConfigFromFile("config.yml")
	if err != nil {
		panic(err)
	}
	if err := phases.Init(cfg); err != nil {
		panic(err)
	}
}
