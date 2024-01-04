package providers

import (
	"bitbucket.org/dptsi/its-go/contracts"
	"its.ac.id/base-go/modules/auth"
)

func registerModules(service contracts.ModuleService) {
	service.Register("auth", auth.SetupModule)
}
