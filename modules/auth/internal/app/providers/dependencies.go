package providers

import (
	"context"
	"os"
	"strings"

	"bitbucket.org/dptsi/base-go-libraries/auth"
	"bitbucket.org/dptsi/base-go-libraries/contracts"
	"bitbucket.org/dptsi/base-go-libraries/oidc"
	"bitbucket.org/dptsi/base-go-libraries/sessions"
	"github.com/samber/do"
	"its.ac.id/base-go/modules/auth/internal/presentation/controllers"
)

func RegisterDependencies(ctx context.Context, i *do.Injector) {
	// Libraries
	do.Provide[*oidc.Client](i, func(i *do.Injector) (*oidc.Client, error) {
		return oidc.NewClient(
			ctx,
			do.MustInvoke[contracts.SessionStorage](i),
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
	do.Provide[*controllers.AuthController](i, func(i *do.Injector) (*controllers.AuthController, error) {
		return controllers.NewAuthController(
			do.MustInvoke[*oidc.Client](i),
			do.MustInvoke[contracts.SessionStorage](i),
			do.MustInvoke[*auth.Service](i),
			do.MustInvoke[*sessions.CookieUtil](i),
		), nil
	})
}
