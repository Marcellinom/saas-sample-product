package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"

	"github.com/dptsi/its-go/app"
	"github.com/dptsi/its-go/database"
	"github.com/dptsi/its-go/providers"
	"github.com/dptsi/its-go/web"
	"github.com/gin-gonic/gin"
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

// @securityDefinitions.apikey	CSRF Token
// @in							header
// @name						x-csrf-token
// @description 				CSRF token yang didapatkan dari browser -> inspect element -> application -> storage -> cookies -> CSRF-TOKEN (Untuk firefox, storage berada pada tab tersendiri)

// @externalDocs.description  Dokumentasi Base Project
// @externalDocs.url          http://localhost:8080/doc/project
func main() {
	//  log.Println("Loading environment variables from .env")
	if err := godotenv.Load(); err != nil {
		log.Panic("Error loading .env file")
	}
	//  log.Println("Environment variables successfully loaded!")

	//  log.Println("Creating application instance...")
	ctx := context.Background()
	application := app.NewApplication(ctx, do.DefaultInjector, config.Config())
	//  log.Println("Application instance successfully created!")

	//  log.Println("Loading framework providers...")
	if err := providers.LoadProviders(application); err != nil {
		panic(err)
	}
	//  log.Println("Framework providers loaded!")

	services := application.Services()

	//  log.Println("Loading application providers...")
	appProviders.LoadAppProviders(application)
	//  log.Println("Application providers loaded!")

	engine := services.WebEngine
	engine.GET("/csrf-cookie", CSRFCookieRoute)

	// programmatically set swagger info
	if os.Getenv("APP_ENV") == "local" {
		//  log.Println("Local environment detected, setting up swagger...")
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
		//  log.Println("Swagger successfully set up!")
	}

	if os.Getenv("APP_DEBUG") == "true" {
		serviceList := application.ListProvidedServices()
		sort.Strings(serviceList)
	
	}
	db := services.Database.GetDefault()
	controller := Controller{db: db}
	engine.GET("/", controller.Test)
	engine.GET("/status", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	webConfig := config.Config()["web"].(web.Config)
	engine.Run(fmt.Sprintf(":%s", webConfig.Port))
}

type Controller struct {
	db *database.Database
}

func (c *Controller) Test(ctx *gin.Context) {
	var users []struct {
		Id int64 `json:"id"`
		Name string `json:"name"`
	} 
	err := c.db.Table("users").Find(&users).Error
	if err != nil {
		ctx.Error(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "success",
		"data":    users,
	})
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
