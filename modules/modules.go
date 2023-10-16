package modules

import (
	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/bootstrap/event"
	"its.ac.id/base-go/modules/auth"
)

func RegisterModules(cfg config.Config, g *gin.Engine, eventHook *event.EventHook) {
	// register modules here
	// e.g.:
	auth.SetupModule(cfg, g, eventHook)
}
