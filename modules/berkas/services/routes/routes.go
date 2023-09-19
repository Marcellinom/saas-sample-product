package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mikestefanello/hooks"
	"its.ac.id/base-go/services/routes"
)

func init() {
	routes.HookBuildRouter.Listen(func(event hooks.Event[*gin.Engine]) {
		r := event.Msg
		g := r.Group("/berkas")

		g.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "berkas",
			})
		})
	})
}
