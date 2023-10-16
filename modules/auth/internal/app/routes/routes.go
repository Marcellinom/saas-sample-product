package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"its.ac.id/base-go/modules/auth/internal/presentation/controllers"
	"its.ac.id/base-go/pkg/auth/middleware"
)

func RegisterRoutes(i *do.Injector, r *gin.Engine) {
	g := r.Group("/auth")
	authController := do.MustInvoke[*controllers.AuthController](i)

	g.POST("/login", authController.Login)
	g.GET("/callback", authController.Callback)
	g.GET("/user", middleware.Auth(), authController.User)
	g.DELETE("/logout", middleware.Auth(), authController.Logout)
}
