package main

import (
	"mobingi/ocean/app/phases"
	"mobingi/ocean/pkg/config"
)

func main() {
	cfg := config.NewConfig()
	if err := phases.Init(cfg); err != nil {
		panic(err)
	}
}
