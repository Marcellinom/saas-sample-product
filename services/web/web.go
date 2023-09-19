package web

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mikestefanello/hooks"
	"github.com/samber/do"
	"its.ac.id/base-go/pkg/app"
	"its.ac.id/base-go/services/config"
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
	g.engine.Run()
}

// HookBuildRouter allows modules the ability to build on the web router
var HookBuildRouter = hooks.NewHook[*gin.Engine]("router.build")

func (g *GinServer) buildRouter() *gin.Engine {
	// Global middleware
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
