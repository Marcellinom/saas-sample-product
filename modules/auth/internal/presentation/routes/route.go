package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"its.ac.id/base-go/modules/auth/internal/presentation/controllers"
	"its.ac.id/base-go/pkg/auth/middleware"
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
	g.POST("/user/switch-active-role", middleware.Auth(), middleware.ActiveRoleIn("mahasiswa", "administrator"), authController.SwitchActiveRole)
}
