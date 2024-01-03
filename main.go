package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"bitbucket.org/dptsi/its-go/app"
	"bitbucket.org/dptsi/its-go/providers"
	"bitbucket.org/dptsi/its-go/web"
	"github.com/joho/godotenv"
	"github.com/samber/do"
	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
	"its.ac.id/base-go/config"
	"its.ac.id/base-go/docs"
	appProviders "its.ac.id/base-go/providers"
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

	log.Println("Creating application instance...")
	ctx := context.Background()
	application := app.NewApplication(ctx, do.DefaultInjector, config.Config())
	log.Println("Application instance successfully created!")

	log.Println("Loading framework providers...")
	if err := providers.LoadProviders(application); err != nil {
		panic(err)
	}
	log.Println("Framework providers loaded!")

	services := application.Services()

	log.Println("Loading modules...")
	appProviders.RegisterModules(services.Module)
	log.Println("Modules loaded!")

	engine := services.WebEngine

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
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		log.Println("Swagger successfully set up!")
	}

	// log.Println("Registering modules...")
	// ctx := context.Background()
	// modules.RegisterModules(ctx, i)
	// log.Println("All modules successfully registered!")

	serviceList := application.ListProvidedServices()
	sort.Strings(serviceList)
	log.Printf(
		"registered %d dependencies: \n%s",
		len(serviceList),
		strings.Join((func() []string {
			arr := make([]string, len(serviceList))
			for i, s := range serviceList {
				arr[i] = fmt.Sprintf("- %s", s)
			}

			return arr
		})(), "\n"),
	)

	engine.GET("/csrf-cookie", CSRFCookieRoute)
	engine.Run()
}

// CSRF cookie godoc
// @Summary		Rute dummy untuk set CSRF-TOKEN cookie
// @Router		/csrf-cookie [get]
// @Tags		CSRF Protection
// @Produce		json
// @Success		200 {object} responses.GeneralResponse{code=int,message=string} "Cookie berhasil diset"
// @Header      default {string} Set-Cookie "CSRF-TOKEN=00000000-0000-0000-0000-000000000000; Path=/"
func CSRFCookieRoute(ctx *web.Context) {
	ctx.JSON(http.StatusOK, web.H{
		"code":    0,
		"message": "success",
		"data":    nil,
	})
}
