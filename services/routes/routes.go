package routes

import (
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
}

func NewGinServer(i *do.Injector) (Server, error) {
	cfg := do.MustInvoke[config.Config](i).App()
	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	return &GinServer{r}, nil
}

func (g *GinServer) Start() {
	g.engine.Run()
}
