package providers

import (
	"github.com/samber/do"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/bootstrap/event"
	moduleConfig "its.ac.id/base-go/modules/auth/internal/app/config"
	"its.ac.id/base-go/modules/auth/internal/presentation/controllers"
)

func RegisterDependencies(i *do.Injector, cfg config.Config, moduleCfg moduleConfig.AuthConfig, eventHook *event.EventHook) {
	// Controllers
	do.Provide[*controllers.AuthController](i, func(i *do.Injector) (*controllers.AuthController, error) {
		return controllers.NewAuthController(cfg, moduleCfg), nil
	})

	// Queries

	// Repositories
}
