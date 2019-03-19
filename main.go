package main

import (
	"log"

	"mobingi/ocean/app"
)

func main() {
	if err := app.Start(); err != nil {
		panic(err)
	}

	log.Println("sucess")
}
