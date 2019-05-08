package main

import (
	"mobingi/ocean/app"
	"mobingi/ocean/pkg/storage"
)

func main() {
	// TODO Move to main init func
	storage.NewMongoClient()

	// storage := storage.NewStorage()
	// kubeconfig, err := storage.GetKubeconf("kubernetes", "admin.conf")
	// if err != nil {
	// 	log.Error(err)
	// }
	// err = client.Init(kubeconfig)
	// if err != nil {
	// 	log.Error(err)
	// }

	// nodes, err := client.GetNode()
	// fmt.Println(nodes.Items)
	if err := app.Start(); err != nil {
		panic(err)
	}
}
