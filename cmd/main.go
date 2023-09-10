package main

import (
	"log"
	"orderservice/internal/app"
	"orderservice/internal/config"
)

// @title OrderService
// @version 1.0
// @description API Server for OrderService Application

// @host localhost:8080
// @BasePath /
func main() {

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	app.Run(cfg)

}
