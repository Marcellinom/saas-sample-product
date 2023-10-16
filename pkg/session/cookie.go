package session

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/bootstrap/config"
)

func AddCookieToResponse(cfg config.SessionConfig, ctx *gin.Context, sessionId string) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	// Set session cookie
	ctx.SetCookie(cfg.CookieName, sessionId, cfg.Lifetime, cfg.CookiePath, cfg.Domain, cfg.Secure, true)
	sess := Default(ctx)
	ctx.SetCookie("CSRF-TOKEN", sess.csrfToken, cfg.Lifetime, cfg.CookiePath, cfg.Domain, cfg.Secure, false)
}
