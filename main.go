package main

import (
	"log"
	"smply/app"
	"smply/config"
)

func main() {
	config.LoadEnv()

	if err := config.InitDB(); err != nil {
		log.Fatal(err)
	}

	app.Start()
}
