package web

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mikestefanello/hooks"
	"github.com/samber/do"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/pkg/app"
	"its.ac.id/base-go/pkg/app/common"
	"its.ac.id/base-go/pkg/session"
	"its.ac.id/base-go/pkg/session/adapters"
	"its.ac.id/base-go/pkg/session/middleware"
)

type Server interface {
	Start()
}

func init() {
	app.HookBoot.Listen(func(e hooks.Event[*do.Injector]) {
		do.Provide[session.Storage](e.Msg, func(i *do.Injector) (session.Storage, error) {
			cfg := do.MustInvoke[config.Config](i).Session()

			switch cfg.Driver {
			case "firestore":
				return setupFirestoreSessionAdapter(i)
			case "sqlite":
				// Contoh penggunaan adapter GORM dengan SQLite
				path := cfg.SQLiteDB
				if path == "" {
					panic("invalid SQLite database path for session")
				}

				db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
				if err != nil {
					panic("failed to connect to SQLite database for session")
				}
				return adapters.NewGorm(db), nil
			case "sqlserver":
				// Contoh penggunaan adapter GORM dengan SQL Server
				username := cfg.SQLServerUsername
				password := cfg.SQLServerPassword
				host := cfg.SQLServerHost
				port := cfg.SQLServerPort
				database := cfg.SQLServerDatabase

				if username == "" || password == "" || host == "" || port == "" || database == "" {
					panic("invalid SQL Server configuration for session")
				}

				dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", username, password, host, port, database)
				db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
				if err != nil {
					panic("failed to connect SQL Server database for session")
				}
				return adapters.NewGorm(db), nil
			}

			panic("unknown session driver")
		})
		do.Provide[Server](e.Msg, NewGinServer)
	})
}

type GinServer struct {
	engine *gin.Engine
	cfg    config.Config
}

func NewGinServer(i *do.Injector) (Server, error) {
	cfg := do.MustInvoke[config.Config](i)
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

	s := &GinServer{r, cfg}
	s.buildRouter()

	return s, nil
}

func (g *GinServer) Start() {
	g.engine.Run(":" + strconv.Itoa(g.cfg.HTTP().Port))
}

// HookBuildRouter allows modules the ability to build on the web router
var HookBuildRouter = hooks.NewHook[*gin.Engine]("router.build")

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
	g.engine.Use(middleware.StartSession())
	g.engine.Use(middleware.VerifyCSRFToken())
	g.engine.Use(g.initiateCorsMiddleware())
	g.engine.GET("/csrf-cookie", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"code":    200,
			"message": "success",
			"data":    nil,
		})
	})

	HookBuildRouter.Dispatch(g.engine)
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
