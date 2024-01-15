package providers

import (
	"bitbucket.org/dptsi/go-modules/auth"
	"github.com/dptsi/its-go/contracts"
)

func registerModules(service contracts.ModuleService) {
	service.Register("auth", auth.SetupModule)
}
