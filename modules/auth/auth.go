package auth

import (
	"github.com/samber/do"
	"its.ac.id/base-go/modules/auth/internal/app/providers"
	"its.ac.id/base-go/modules/auth/internal/presentation/routes"
)

func SetupModule(i *do.Injector) {
	i = i.Clone()

	providers.RegisterDependencies(i)

	routes.RegisterRoutes(i)
}
