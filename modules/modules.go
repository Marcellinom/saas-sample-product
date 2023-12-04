package modules

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"its.ac.id/base-go/bootstrap/event"
	"its.ac.id/base-go/modules/auth"
)

func RegisterModules(i *do.Injector, g *gin.Engine, eventHook *event.EventHook) {
	// register modules here
	// e.g.:
	auth.SetupModule(i, g, eventHook)
}
