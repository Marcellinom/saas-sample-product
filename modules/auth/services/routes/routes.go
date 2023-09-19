package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mikestefanello/hooks"
	"github.com/samber/do"
	"its.ac.id/base-go/modules/auth/internal/controllers"
	"its.ac.id/base-go/pkg/auth/middleware"
	"its.ac.id/base-go/services/web"
)

func init() {
	web.HookBuildRouter.Listen(func(event hooks.Event[*gin.Engine]) {
		r := event.Msg
		g := r.Group("/auth")
		i := do.DefaultInjector
		authController := controllers.NewAuthController(i)

		g.GET("/login", authController.Login)
		g.GET("/user", middleware.Auth(), authController.User)
		g.GET("/logout", middleware.Auth(), authController.Logout)
	})
}
