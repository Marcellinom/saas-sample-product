package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/bootstrap/event"
	moduleConfig "its.ac.id/base-go/modules/auth/internal/app/config"
	"its.ac.id/base-go/modules/auth/internal/app/providers"
	"its.ac.id/base-go/modules/auth/internal/presentation/routes"
)

func SetupModule(cfg config.Config, g *gin.Engine, eventHook *event.EventHook) {
	i := do.New()

	moduleCfg, err := moduleConfig.SetupConfig()
	if err != nil {
		panic(err)
	}

	providers.RegisterDependencies(i, cfg, moduleCfg, eventHook, g)

	routes.RegisterRoutes(i, g)
}
