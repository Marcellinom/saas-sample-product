package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/samber/do"

	// Services
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/bootstrap/event"
	"its.ac.id/base-go/bootstrap/web"
	"its.ac.id/base-go/modules"
)

// @securityDefinitions.apikey	Session
// @in							cookie
// @name						akademik_its_ac_id_session
// @securityDefinitions.apikey	CSRF Token
// @in							header
// @name						x-csrf-token
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

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

	i := do.DefaultInjector
	providedServices := i.ListProvidedServices()
	log.Printf("registered %d dependencies: %v", len(providedServices), providedServices)

	server.Start()
}
