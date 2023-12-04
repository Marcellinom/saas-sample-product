package main

import (
	"context"
	"log"
	"net/url"
	"os"

	sessionsMiddleware "bitbucket.org/dptsi/go-framework/sessions/middleware"
	"bitbucket.org/dptsi/go-framework/web"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/samber/do"
	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
	"its.ac.id/base-go/bootstrap"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/bootstrap/event"
	"its.ac.id/base-go/bootstrap/middleware"
	"its.ac.id/base-go/docs"
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

	i := do.DefaultInjector
	bootstrap.CreateObjects(i)
	log.Println("Setting up web server...")
	server, err := web.SetupServer(web.Config{
		IsDebugMode: os.Getenv("APP_DEBUG") == "true",
		Environment: os.Getenv("APP_ENV"),
	})
	server.Use(do.MustInvoke[*middleware.MiddlewareGroup](i).GlobalMiddleware()...)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Web server successfully set up!")

	// programmatically set swagger info
	if os.Getenv("APP_ENV") == "local" {
		log.Println("Local environment detected, setting up swagger...")
		appUrlEnv := os.Getenv("APP_URL")
		appURL, err := url.Parse(appUrlEnv)
		if err != nil {
			appURL, _ = url.Parse("http://localhost:8080")
		}
		docs.SwaggerInfo.Title = os.Getenv("APP_NAME")
		// docs.SwaggerInfo.Version = r.cfg.App().Version
		docs.SwaggerInfo.Host = appURL.Host
		docs.SwaggerInfo.BasePath = ""
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
		server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		log.Println("Swagger successfully set up!")
	}
	server.GET("/csrf-cookie", sessionsMiddleware.CSRFCookieRoute)

	log.Println("Setting up database...")
	config.SetupDatabase(i)
	log.Println("Database successfully set up!")

	log.Println("Setting up session...")
	config.SetupSession(i)
	log.Println("Session successfully set up!")

	log.Println("Setting up event hook...")
	eventHook := event.SetupEventHook()
	log.Println("Event hook successfully set up!")

	log.Println("Registering dependencies...")
	do.Provide[*event.EventHook](i, func(i *do.Injector) (*event.EventHook, error) {
		return eventHook, nil
	})
	do.Provide[*gin.Engine](i, func(i *do.Injector) (*gin.Engine, error) {
		return server, nil
	})

	log.Println("Registering modules...")
	ctx := context.Background()
	modules.RegisterModules(ctx, i)
	log.Println("All modules successfully registered!")

	providedServices := i.ListProvidedServices()
	log.Printf("registered %d dependencies: %v", len(providedServices), providedServices)

	server.Run()
}
