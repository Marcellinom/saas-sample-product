package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mikestefanello/hooks"
	"github.com/samber/do"
	"its.ac.id/base-go/modules/auth/internal/app/controllers"
	"its.ac.id/base-go/pkg/auth/middleware"
	"its.ac.id/base-go/services/web"
)

func registerRoutes(r *gin.Engine) {
	g := r.Group("/auth")
	i := do.DefaultInjector
	authController := controllers.NewAuthController(i)

	g.POST("/login", authController.Login)
	g.GET("/user", middleware.Auth(), authController.User)
	g.DELETE("/logout", middleware.Auth(), authController.Logout)
}

func init() {
	web.HookBuildRouter.Listen(func(event hooks.Event[*gin.Engine]) {
		registerRoutes(event.Msg)
	})
}
