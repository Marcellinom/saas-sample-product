package routes

import (
	"bitbucket.org/dptsi/go-framework/auth/middleware"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	mg "its.ac.id/base-go/bootstrap/middleware"
	"its.ac.id/base-go/modules/auth/internal/presentation/controllers"
)

func RegisterRoutes(i *do.Injector) {
	middlewareGroup := do.MustInvoke[*mg.MiddlewareGroup](i)
	r := do.MustInvoke[*gin.Engine](i)
	g := r.Group("/auth")
	g.Use(middlewareGroup.WebMiddleware()...)

	// Controllers
	authController := do.MustInvoke[*controllers.AuthController](i)

	// Routes
	g.POST("/login", authController.Login)
	g.GET("/callback", authController.Callback)
	g.GET("/user", middleware.Auth(), authController.User)
	g.DELETE("/logout", middleware.Auth(), authController.Logout)
	g.POST("/user/switch-active-role", middleware.Auth(), authController.SwitchActiveRole)
}
