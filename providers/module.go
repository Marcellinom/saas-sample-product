package providers

import (
	"github.com/dptsi/its-go/contracts"
	"its.ac.id/base-go/modules/auth"
)

func registerModules(service contracts.ModuleService) {
	service.Register("auth", auth.SetupModule)
}
