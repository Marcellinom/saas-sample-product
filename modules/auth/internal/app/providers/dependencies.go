package providers

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/bootstrap/event"
	moduleConfig "its.ac.id/base-go/modules/auth/internal/app/config"
	"its.ac.id/base-go/modules/auth/internal/presentation/controllers"
	"its.ac.id/base-go/modules/auth/internal/presentation/routes"
)

func RegisterDependencies(i *do.Injector, cfg config.Config, moduleCfg moduleConfig.AuthConfig, eventHook *event.EventHook, g *gin.Engine) {
	// Libraries

	// Queries

	// Repositories

	// Controllers
	authController := controllers.NewAuthController(cfg, moduleCfg)
	r := routes.NewRoutes(g, authController)

	// Route
	do.Provide[*routes.Route](i, func(i *do.Injector) (*routes.Route, error) {
		return r, nil
	})
}
