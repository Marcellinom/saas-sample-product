package web

import (
	"github.com/gin-gonic/gin"
)

func buildRouter(r *gin.Engine) {
	// Custom Handlers

	// Global middleware
	// r.Use(middleware.StartSession(
	// 	r.sessionStorage,
	// 	sessions.AddSessionCookieToResponseAttributes{
	// 		Name:           r.cfg.Session().CookieName,
	// 		Path:           r.cfg.Session().CookiePath,
	// 		Domain:         r.cfg.Session().Domain,
	// 		MaxAge:         r.cfg.Session().Lifetime,
	// 		Secure:         r.cfg.Session().Secure,
	// 		CsrfCookieName: "CSRF-TOKEN",
	// 	},
	// ))
	// r.Use(middleware.VerifyCSRFToken())

	// Global routes
	// isLocal := r.cfg.App().Env == "local"
	// isStaging := r.cfg.App().Env == "staging"
	// if isLocal || isStaging {
	// 	r.Static("/doc/project", "./static/mkdocs")
	// }
	// r.GET("/csrf-cookie", r.handleCSRFCookie)

	// Swagger
	// appURL, err := url.Parse(r.cfg.App().URL)
	// if err != nil {
	// 	appURL, _ = url.Parse("http://localhost:8080")
	// }

	// programmatically set swagger info
	// if isLocal || isStaging {
	// 	docs.SwaggerInfo.Title = r.cfg.App().Name
	// 	docs.SwaggerInfo.Description = r.cfg.App().Description
	// 	docs.SwaggerInfo.Version = r.cfg.App().Version
	// 	docs.SwaggerInfo.Host = appURL.Host
	// 	docs.SwaggerInfo.BasePath = ""
	// 	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	// 	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// }

}

// CSRF cookie godoc
// @Summary		Rute dummy untuk set CSRF-TOKEN cookie
// @Router		/csrf-cookie [get]
// @Tags		CSRF Protection
// @Produce		json
// @Success		200 {object} responses.GeneralResponse{code=int,message=string} "Cookie berhasil diset"
// @Header      default {string} Set-Cookie "CSRF-TOKEN=00000000-0000-0000-0000-000000000000; Path=/"
// func (g *GinServer) handleCSRFCookie(ctx *gin.Context) {
// 	ctx.JSON(200, gin.H{
// 		"code":    0,
// 		"message": "success",
// 		"data":    nil,
// 	})
// }
