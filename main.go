package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"bitbucket.org/dptsi/base-go-libraries/auth"
	"bitbucket.org/dptsi/base-go-libraries/contracts"
	"bitbucket.org/dptsi/base-go-libraries/database"
	"bitbucket.org/dptsi/base-go-libraries/sessions"
	sessionsMiddleware "bitbucket.org/dptsi/base-go-libraries/sessions/middleware"
	"bitbucket.org/dptsi/base-go-libraries/web"
	webMiddleware "bitbucket.org/dptsi/base-go-libraries/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/samber/do"
	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
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
	createObjects(i)
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
	server.GET("/csrf-cookie", handleCSRFCookie)

	log.Println("Setting up database...")
	config.SetupDatabase(i)
	log.Println("Database successfully set up!")

	log.Println("Setting up session...")
	config.SetupSession(i)
	log.Println("Session successfully set up!")

	log.Println("Setting up event hook...")
	eventHook := event.SetupEventHook()
	log.Println("Event hook successfully set up!")
	do.Provide[*event.EventHook](i, func(i *do.Injector) (*event.EventHook, error) {
		return eventHook, nil
	})

	log.Println("Registering modules...")
	modules.RegisterModules(i, server, eventHook)
	log.Println("All modules successfully registered!")

	providedServices := i.ListProvidedServices()
	log.Printf("registered %d dependencies: %v", len(providedServices), providedServices)

	server.Run()
}

func createObjects(i *do.Injector) {
	do.Provide[*webMiddleware.HandleCors](i, func(i *do.Injector) (*webMiddleware.HandleCors, error) {
		return &webMiddleware.HandleCors{
			AllowedOrigins:   []string{"http://localhost:3000"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"content-type", "x-csrf-token"},
			ExposedHeaders:   []string{},
			AllowCredentials: true,
			MaxAge:           0,
		}, nil
	})
	do.Provide[*middleware.MiddlewareGroup](i, func(i *do.Injector) (*middleware.MiddlewareGroup, error) {
		return middleware.NewMiddlewareGroup(i), nil
	})
	do.Provide[*database.Manager](i, func(i *do.Injector) (*database.Manager, error) {
		return database.NewManager(), nil
	})
	sessionMaxAge, err := strconv.Atoi(os.Getenv("SESSION_MAX_AGE"))
	if err != nil {
		sessionMaxAge = 86400
	}
	sessionConfig := sessions.SessionsConfig{
		Name:           os.Getenv("SESSION_NAME"),
		CsrfCookieName: os.Getenv("SESSION_CSRF_COOKIE_NAME"),
		MaxAge:         sessionMaxAge,
		Path:           os.Getenv("SESSION_PATH"),
		Domain:         os.Getenv("SESSION_DOMAIN"),
		Secure:         os.Getenv("SESSION_SECURE") == "true",
	}
	do.Provide[*sessions.CookieUtil](i, func(i *do.Injector) (*sessions.CookieUtil, error) {
		return sessions.NewCookieUtil(sessionConfig), nil
	})

	do.Provide[*sessionsMiddleware.StartSession](i, func(i *do.Injector) (*sessionsMiddleware.StartSession, error) {
		return sessionsMiddleware.NewStartSession(
			sessionConfig,
			do.MustInvoke[contracts.SessionStorage](i),
			*(do.MustInvoke[*sessions.CookieUtil](i)),
		), nil
	})
	do.Provide[*sessionsMiddleware.VerifyCSRFToken](i, func(i *do.Injector) (*sessionsMiddleware.VerifyCSRFToken, error) {
		return sessionsMiddleware.NewVerifyCSRFToken(), nil
	})
	do.Provide[*auth.Service](i, func(i *do.Injector) (*auth.Service, error) {
		return auth.NewService(
			do.MustInvoke[contracts.SessionStorage](i),
		), nil
	})
}

// CSRF cookie godoc
// @Summary		Rute dummy untuk set CSRF-TOKEN cookie
// @Router		/csrf-cookie [get]
// @Tags		CSRF Protection
// @Produce		json
// @Success		200 {object} responses.GeneralResponse{code=int,message=string} "Cookie berhasil diset"
// @Header      default {string} Set-Cookie "CSRF-TOKEN=00000000-0000-0000-0000-000000000000; Path=/"
func handleCSRFCookie(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    nil,
	})
}
