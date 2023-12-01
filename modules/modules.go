package modules

import (
	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/bootstrap/event"
)

func RegisterModules(g *gin.Engine, eventHook *event.EventHook) {
	// register modules here
	// e.g.:
	// auth.SetupModule(cfg, g, eventHook)
}
