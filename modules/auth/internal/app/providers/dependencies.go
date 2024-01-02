package providers

import (
	"os"
	"strings"

	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/module"
	"bitbucket.org/dptsi/go-framework/oidc"
	"its.ac.id/base-go/modules/auth/internal/presentation/controllers"
)

func RegisterDependencies(mod contracts.Module) {
	// Libraries
	module.Bind[*oidc.Client](mod, "oidc_client", func(mod contracts.Module) (*oidc.Client, error) {
		sessionsService := mod.App().Services().Session
		return oidc.NewClient(
			mod.App().Context(),
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
	module.Bind[*controllers.AuthController](mod, "controllers.auth", func(mod contracts.Module) (*controllers.AuthController, error) {
		services := mod.App().Services()
		oidcClient := module.MustMake[*oidc.Client](mod, "oidc_client", module.DependencyScopeModule)

		return controllers.NewAuthController(
			services.Session,
			services.Auth,
			oidcClient,
		), nil
	})
}
