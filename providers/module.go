package providers

import (
	"bitbucket.org/dptsi/go-framework/contracts"
	"its.ac.id/base-go/modules/auth"
)

func RegisterModules(service contracts.ModuleService) {
	service.Register("auth", auth.SetupModule)
}
