package services

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"its.ac.id/base-go/pkg/auth/internal/utils"
	"its.ac.id/base-go/services/config"
)

func Logout(ctx *gin.Context) error {
	cfg := do.MustInvoke[config.Config](do.DefaultInjector)
	authCfg := cfg.Auth()
	httpCfg := cfg.HTTP()

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
