package routes

import (
	"bitbucket.org/dptsi/its-go/contracts"
	"bitbucket.org/dptsi/its-go/module"
	"its.ac.id/base-go/modules/auth/internal/presentation/controllers"
)

func RegisterRoutes(mod contracts.Module) {
	engine := mod.App().Services().WebEngine
	middlewareService := mod.App().Services().Middleware

	// Routing
	g := engine.Group("/auth")

	// Controllers
	authController := module.MustMake[*controllers.AuthController](mod, "controllers.auth", module.DependencyScopeModule)

	// Routes
	g.POST("/login", authController.Login)
	g.GET("/callback", authController.Callback)
	g.GET("/user", middlewareService.Use("auth", nil), authController.User)
	g.DELETE("/logout", middlewareService.Use("auth", nil), authController.Logout)
	g.POST("/user/switch-active-role", middlewareService.Use("auth", nil), authController.SwitchActiveRole)
}
