package providers

import "bitbucket.org/dptsi/its-go/contracts"

func LoadAppProviders(application contracts.Application) {
	services := application.Services()

	extendAuth(application)
	registerEvents(application)
	registerMiddlewares(application)
	registerModules(services.Module)
}
