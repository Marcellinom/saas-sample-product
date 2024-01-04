package providers

import (
	"os"
	"strings"

	"bitbucket.org/dptsi/its-go/app"
	"bitbucket.org/dptsi/its-go/contracts"
	"bitbucket.org/dptsi/its-go/oidc"
	"its.ac.id/base-go/modules/auth/internal/presentation/controllers"
)

func RegisterDependencies(mod contracts.Module) {
	// Libraries
	app.Bind[*oidc.Client](mod.App(), "modules.auth.oidc_client", func(application contracts.Application) (*oidc.Client, error) {
		sessionsService := application.Services().Session
		return oidc.NewClient(
			application.Context(),
			sessionsService,
			os.Getenv("OIDC_PROVIDER"),
			os.Getenv("OIDC_CLIENT_ID"),
			os.Getenv("OIDC_CLIENT_SECRET"),
			os.Getenv("OIDC_REDIRECT_URL"),
			strings.Split(os.Getenv("OIDC_SCOPES"), ","),
		)
	})

	// Queries

	// Repositories

	// Controllers
	app.Bind[*controllers.AuthController](mod.App(), "modules.auth.controllers.auth", func(application contracts.Application) (*controllers.AuthController, error) {
		services := mod.App().Services()
		oidcClient := app.MustMake[*oidc.Client](application, "modules.auth.oidc_client")

		return controllers.NewAuthController(
			services.Session,
			services.Auth,
			oidcClient,
		), nil
	})
}
