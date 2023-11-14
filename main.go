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

// @contact.name   Direktorat Pengembangan Teknologi dan Sistem Informasi (DPTSI) - ITS
// @contact.url    http://its.ac.id/dptsi
// @contact.email  dptsi@its.ac.id

// @securityDefinitions.apikey	Session
// @in							cookie
// @name						akademik_its_ac_id_session
// @securityDefinitions.apikey	CSRF Token
// @in							header
// @name						x-csrf-token

// @externalDocs.description  Dokumentasi Base Project
// @externalDocs.url          http://localhost:8080/doc/project
func main() {
	log.Println("Loading environment variables from .env")
	if err := godotenv.Load(); err != nil {
		log.Panic("Error loading .env file")
	}
	log.Println("Environment variables successfully loaded!")

	log.Println("Loading application configurations...")
	cfg, err := config.SetupAppConfig()
	if err != nil {
		log.Panic(err)
	}
	log.Println("Application configurations successfully loaded!")

	log.Println("Setting up web server...")
	server, err := web.SetupServer(cfg)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Web server successfully set up!")

	log.Println("Setting up event hook...")
	eventHook := event.SetupEventHook()
	log.Println("Event hook successfully set up!")

	log.Println("Registering modules...")
	modules.RegisterModules(cfg, server.Engine(), eventHook)
	log.Println("All modules successfully registered!")

	i := do.DefaultInjector
	providedServices := i.ListProvidedServices()
	log.Printf("registered %d dependencies: %v", len(providedServices), providedServices)

	server.Start()
}
