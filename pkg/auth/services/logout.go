package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/pkg/auth/internal/utils"
)

func Logout(ctx *gin.Context) error {
	cfg := do.MustInvoke[config.Config](do.DefaultInjector)
	authCfg := cfg.Auth()
	httpCfg := cfg.HTTP()

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		utils.GetCookieName(),
		"",
		-1,
		authCfg.CookiePath,
		authCfg.CookieDomain,
		httpCfg.Secure,
		true,
	)

	return nil
}
