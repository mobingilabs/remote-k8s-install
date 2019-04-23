package main

import "mobingi/ocean/app"

func main() {
	if err := app.Start(); err != nil {
		panic(err)
	}
}
