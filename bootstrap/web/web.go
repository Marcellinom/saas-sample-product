package web

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

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
	storage, err := setupSessionStorage(cfg.Session())
	if err != nil {
		return nil, err
	}
	return newGinServer(cfg, storage)
}

type GinServer struct {
	engine         *gin.Engine
	cfg            config.Config
	sessionStorage session.Storage
}

func newGinServer(cfg config.Config, sessionStorage session.Storage) (Server, error) {
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
	g.engine.StaticFile("/oas3.yml", "./oas3.yml")
	g.engine.Static("/doc/api", "./static/swagger-ui")
	g.engine.Static("/doc/project", "./static/mkdocs")
	g.engine.Use(middleware.StartSession(g.cfg.Session(), g.sessionStorage))
	g.engine.Use(middleware.VerifyCSRFToken())
	g.engine.Use(g.initiateCorsMiddleware())
	g.engine.GET("/csrf-cookie", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success",
			"data":    nil,
		})
	})

	return g.engine
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
