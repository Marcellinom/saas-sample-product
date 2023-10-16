package main

import (
	"log"

	"github.com/joho/godotenv"

	// Services
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/bootstrap/event"
	"its.ac.id/base-go/bootstrap/web"
	"its.ac.id/base-go/modules"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	// i := do.DefaultInjector
	// providedServices := i.ListProvidedServices()

	cfg, err := config.SetupAppConfig()
	if err != nil {
		panic(err)
	}

	server, err := web.SetupServer(cfg)
	if err != nil {
		panic(err)
	}

	eventHook := event.SetupEventHook()
	modules.RegisterModules(cfg, server.Engine(), eventHook)

	// log.Printf("registered %d dependencies: %v", len(providedServices), providedServices)

	server.Start()
}
