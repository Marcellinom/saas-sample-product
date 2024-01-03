package auth

import (
	"bitbucket.org/dptsi/its-go/contracts"
	"its.ac.id/base-go/modules/auth/internal/app/providers"
	"its.ac.id/base-go/modules/auth/internal/presentation/routes"
)

func SetupModule(mod contracts.Module) {
	providers.RegisterDependencies(mod)
	routes.RegisterRoutes(mod)
}
