package auth

import (
	"context"

	"github.com/samber/do"
	"its.ac.id/base-go/modules/auth/internal/app/providers"
	"its.ac.id/base-go/modules/auth/internal/presentation/routes"
)

func SetupModule(ctx context.Context, i *do.Injector) {
	i = i.Clone()

	providers.RegisterDependencies(ctx, i)

	routes.RegisterRoutes(i)
}
