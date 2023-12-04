package providers

import (
	"context"

	"bitbucket.org/dptsi/base-go-libraries/auth"
	"bitbucket.org/dptsi/base-go-libraries/contracts"
	"bitbucket.org/dptsi/base-go-libraries/oidc"
	"bitbucket.org/dptsi/base-go-libraries/sessions"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"its.ac.id/base-go/bootstrap/event"
	moduleConfig "its.ac.id/base-go/modules/auth/internal/app/config"
	"its.ac.id/base-go/modules/auth/internal/presentation/controllers"
)

func RegisterDependencies(i *do.Injector, moduleCfg moduleConfig.AuthConfig, eventHook *event.EventHook, g *gin.Engine) {
	ctx := context.Background()
	// Libraries
	oidcCfg := moduleCfg.Oidc()
	do.Provide[*oidc.Client](i, func(i *do.Injector) (*oidc.Client, error) {
		return oidc.NewClient(
			ctx,
			do.MustInvoke[contracts.SessionStorage](i),
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
		return controllers.NewAuthController(
			do.MustInvoke[*oidc.Client](i),
			do.MustInvoke[contracts.SessionStorage](i),
			do.MustInvoke[*auth.Service](i),
			do.MustInvoke[*sessions.CookieUtil](i),
		), nil
	})
}
