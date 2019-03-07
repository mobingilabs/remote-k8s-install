package main

import (
	"mobingi/ocean/pkg/certs"
	"mobingi/ocean/pkg/config"
)

func main() {
	cfg := config.NewConfig()
	err := certs.CreatePKIAssets(cfg)
	if err != nil {
		panic(err)
	}
}
