package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mikestefanello/hooks"
	"github.com/samber/do"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/pkg/app"
)

type Server interface {
	Start()
}

func init() {
	app.HookBoot.Listen(func(e hooks.Event[*do.Injector]) {
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
	g.engine.Use(gin.Recovery())
	g.engine.Use(g.initiateCorsMiddleware())

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
