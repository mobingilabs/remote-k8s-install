package main

import (
	"mobingi/ocean/app"
	"mobingi/ocean/pkg/kubernetes/client"
	"mobingi/ocean/pkg/log"
	"mobingi/ocean/pkg/services/tencent"
	"mobingi/ocean/pkg/storage"
)

func main() {
	// TODO Move to main init func
	storage.NewMongoClient()
	tencent.Init()

	err := client.InitClustersAndNodes()
	if err != nil {
		log.Error(err)
		return
	}
	err = client.InitClustersMonitor()
	if err != nil {
		log.Error(err)
		return
	}

	if err := app.Start(); err != nil {
		panic(err)
	}
}
