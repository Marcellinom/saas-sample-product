package main

import (
	"log"

	"bitbucket.org/dptsi/base-go-libraries/web"
	"bitbucket.org/dptsi/base-go-libraries/web/middleware"
	"github.com/joho/godotenv"
	"github.com/samber/do"
	// Services
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

	i := do.DefaultInjector
	createObjects(i)
	log.Println("Setting up web server...")
	server, err := web.SetupServer()
	if err != nil {
		log.Panic(err)
	}
	log.Println("Web server successfully set up!")

	// log.Println("Setting up event hook...")
	// eventHook := event.SetupEventHook()
	// log.Println("Event hook successfully set up!")

	// log.Println("Registering modules...")
	// modules.RegisterModules(server.Engine(), eventHook)
	// log.Println("All modules successfully registered!")

	providedServices := i.ListProvidedServices()
	log.Printf("registered %d dependencies: %v", len(providedServices), providedServices)

	server.Run()
}

func createObjects(i *do.Injector) {
	do.ProvideNamed[*middleware.HandleCors](i, "HandleCorsMiddleware", func(i *do.Injector) (*middleware.HandleCors, error) {
		return &middleware.HandleCors{
			AllowedOrigins:   []string{"http://localhost:3000"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"content-type", "x-csrf-token"},
			ExposedHeaders:   []string{},
			AllowCredentials: true,
			MaxAge:           0,
		}, nil
	})

}
