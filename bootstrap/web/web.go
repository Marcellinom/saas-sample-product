package web

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"its.ac.id/base-go/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/pkg/app/common"
	"its.ac.id/base-go/pkg/session"
	"its.ac.id/base-go/pkg/session/middleware"
)

type Server interface {
	Start()
	Engine() *gin.Engine
}

func SetupServer(cfg config.Config) (Server, error) {
	log.Println("Setting up session storage...")
	storage, err := setupSessionStorage(cfg.Session())
	if err != nil {
		return nil, fmt.Errorf("setup session storage: %w", err)
	}
	log.Println("Session storage successfully set up!")
	return newGinServer(cfg, storage)
}

type GinServer struct {
	engine         *gin.Engine
	cfg            config.Config
	sessionStorage session.Storage
}

func newGinServer(cfg config.Config, sessionStorage session.Storage) (Server, error) {
	log.Println("Setting up Gin server...")
	appCfg := cfg.App()
	if appCfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	s := &GinServer{r, cfg, sessionStorage}
	s.buildRouter()

	log.Println("Gin server successfully set up!")
	return s, nil
}

func (g *GinServer) Start() {
	g.engine.Run(":" + strconv.Itoa(g.cfg.HTTP().Port))
}

func (g *GinServer) Engine() *gin.Engine {
	return g.engine
}

func (g *GinServer) buildRouter() *gin.Engine {
	// Custom Handlers
	g.engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "not_found",
			"data":    nil,
		})
	})
	g.engine.HandleMethodNotAllowed = true
	g.engine.NoMethod(func(ctx *gin.Context) {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{
			"code":    http.StatusMethodNotAllowed,
			"message": "method_not_allowed",
			"data":    nil,
		})
	})

	// Global middleware
	g.engine.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		c.JSON(http.StatusInternalServerError, common.InternalServerErrorResponse)
	}))
	isLocal := g.cfg.App().Env == "local"
	isStaging := g.cfg.App().Env == "staging"
	if isLocal || isStaging {
		g.engine.Static("/doc/project", "./static/mkdocs")
	}
	g.engine.Use(middleware.StartSession(g.cfg.Session(), g.sessionStorage))
	g.engine.Use(middleware.VerifyCSRFToken())
	g.engine.Use(g.initiateCorsMiddleware())
	g.engine.GET("/csrf-cookie", g.handleCSRFCookie)

	appURL, err := url.Parse(g.cfg.App().URL)
	if err != nil {
		appURL, _ = url.Parse("http://localhost:8080")
	}

	// programmatically set swagger info
	if isLocal || isStaging {
		docs.SwaggerInfo.Title = g.cfg.App().Name
		docs.SwaggerInfo.Description = g.cfg.App().Description
		docs.SwaggerInfo.Version = g.cfg.App().Version
		docs.SwaggerInfo.Host = appURL.Host
		docs.SwaggerInfo.BasePath = ""
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
		g.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return g.engine
}

// CSRF cookie godoc
// @Summary		Rute dummy untuk set CSRF-TOKEN cookie
// @Router		/csrf-cookie [get]
// @Tags		CSRF Protection
// @Produce		json
// @Success		200 {object} responses.GeneralResponse{code=int,message=string} "Cookie berhasil diset"
// @Header      default {string} Set-Cookie "CSRF-TOKEN=00000000-0000-0000-0000-000000000000; Path=/"
func (g *GinServer) handleCSRFCookie(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    nil,
	})
}

func (g *GinServer) initiateCorsMiddleware() gin.HandlerFunc {
	cfg := g.cfg.Cors()
	cors := cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     cfg.AllowedMethods,
		AllowHeaders:     cfg.AllowedHeaders,
		ExposeHeaders:    cfg.ExposedHeaders,
		AllowCredentials: cfg.SupportCred,
		MaxAge:           time.Duration(cfg.MaxAge) * time.Second,
	})

	return cors
}
