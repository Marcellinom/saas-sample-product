package routes

import (
	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/modules/auth/internal/presentation/controllers"
	"its.ac.id/base-go/pkg/auth/middleware"
)

type Route struct {
	g              *gin.Engine
	authController *controllers.AuthController
}

func NewRoutes(g *gin.Engine, authController *controllers.AuthController) *Route {
	return &Route{
		g:              g,
		authController: authController,
	}
}

func (r Route) RegisterRoutes() {
	g := r.g.Group("/auth")

	g.POST("/login", r.authController.Login)
	g.GET("/callback", r.authController.Callback)
	g.GET("/user", middleware.Auth(), r.authController.User)
	g.DELETE("/logout", middleware.Auth(), r.authController.Logout)
}
