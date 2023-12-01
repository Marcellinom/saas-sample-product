package routes

import (
	"bitbucket.org/dptsi/base-go-libraries/auth/middleware"
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"its.ac.id/base-go/modules/auth/internal/presentation/controllers"
)

func RegisterRoutes(i *do.Injector, r *gin.Engine) {
	g := r.Group("/auth")

	// Controllers
	authController := do.MustInvoke[*controllers.AuthController](i)

	// Routes
	g.POST("/login", authController.Login)
	g.GET("/callback", authController.Callback)
	g.GET("/user", middleware.Auth(), authController.User)
	g.DELETE("/logout", middleware.Auth(), authController.Logout)
	g.POST("/user/switch-active-role", middleware.Auth(), authController.SwitchActiveRole)
}
