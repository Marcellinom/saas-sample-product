package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikestefanello/hooks"
	"github.com/samber/do"
	"its.ac.id/base-go/modules/auth/internal/controllers"
	"its.ac.id/base-go/pkg/auth/contracts"
	"its.ac.id/base-go/pkg/auth/middleware"
	"its.ac.id/base-go/pkg/auth/services"
	"its.ac.id/base-go/services/web"
)

func init() {
	web.HookBuildRouter.Listen(func(event hooks.Event[*gin.Engine]) {
		r := event.Msg
		g := r.Group("/auth")
		i := do.DefaultInjector
		authController := controllers.NewAuthController(i)

		g.GET("/login", func(c *gin.Context) {
			u := contracts.NewUser("123")
			u.AddRole("admin", []string{"admin"}, true)

			services.Login(c, u)
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

		g.GET("/logout", middleware.Auth(), authController.Logout)
	})
}
