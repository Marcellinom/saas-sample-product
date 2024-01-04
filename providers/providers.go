package providers

import "github.com/dptsi/its-go/contracts"

func LoadAppProviders(application contracts.Application) {
	services := application.Services()

	extendAuth(application)
	registerEvents(application)
	registerMiddlewares(application)
	registerModules(services.Module)
}
