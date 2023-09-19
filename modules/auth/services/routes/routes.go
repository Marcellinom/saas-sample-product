package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikestefanello/hooks"
	"its.ac.id/base-go/pkg/auth/contracts"
	"its.ac.id/base-go/pkg/auth/middleware"
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
			c.JSON(http.StatusOK, gin.H{
				"message": "logged in",
			})
		})
		g.GET("/user", middleware.Auth(), func(ctx *gin.Context) {
			u := services.User(ctx)
			var roles []gin.H
			for _, r := range u.Roles() {
				roles = append(roles, gin.H{
					"name":        r.Name,
					"permissions": r.Permissions,
					"is_default":  r.IsDefault,
				})
			}

			ctx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": "user",
				"data": gin.H{
					"id":          u.Id(),
					"active_role": u.ActiveRole(),
					"roles":       roles,
				},
			})
		})
	})
}
