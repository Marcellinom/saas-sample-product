package providers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/bootstrap/event"
	moduleConfig "its.ac.id/base-go/modules/auth/internal/app/config"
	"its.ac.id/base-go/modules/auth/internal/presentation/controllers"
	"its.ac.id/base-go/pkg/oidc"
)

func RegisterDependencies(i *do.Injector, cfg config.Config, moduleCfg moduleConfig.AuthConfig, eventHook *event.EventHook, g *gin.Engine) {
	ctx := context.Background()
	// Libraries
	oidcCfg := moduleCfg.Oidc()
	do.Provide[*oidc.Client](i, func(i *do.Injector) (*oidc.Client, error) {
		return oidc.NewClient(
			ctx,
			oidcCfg.Provider,
			oidcCfg.ClientID,
			oidcCfg.ClientSecret,
			oidcCfg.RedirectURL,
			oidcCfg.Scopes,
		)
	})

	// Queries

	// Repositories

	// Controllers
	do.Provide[*controllers.AuthController](i, func(i *do.Injector) (*controllers.AuthController, error) {
		return controllers.NewAuthController(cfg, moduleCfg, do.MustInvoke[*oidc.Client](i)), nil
	})
}
