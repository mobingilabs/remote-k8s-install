package main

import (
	"mobingi/ocean/app"
	"mobingi/ocean/pkg/storage"
)

func main() {
	// TODO Move to main init func
	storage.NewMongoClient()

	if err := app.Start(); err != nil {
		panic(err)
	}
}
