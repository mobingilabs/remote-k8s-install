package main

import (
	"mobingi/ocean/app"
	_ "net/http/pprof"
)

func main() {
	if err := app.Start(); err != nil {
		panic(err)
	}
}
