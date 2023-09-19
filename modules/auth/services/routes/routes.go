package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mikestefanello/hooks"
	"its.ac.id/base-go/pkg/auth/contracts"
	"its.ac.id/base-go/pkg/auth/services"
	"its.ac.id/base-go/services/web"
)

func init() {
	web.HookBuildRouter.Listen(func(event hooks.Event[*gin.Engine]) {
		r := event.Msg
		g := r.Group("/auth")

		g.GET("/login", func(c *gin.Context) {
			s := services.NewLoginService(c)
			u := contracts.NewUser("123")
			u.AddRole("admin", []string{"admin"}, true)

			s.Login(u)
			c.JSON(200, gin.H{
				"message": "logged in",
			})
		})
	})
}
