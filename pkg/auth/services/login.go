package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/samber/do"
	"its.ac.id/base-go/pkg/auth/contracts"
	"its.ac.id/base-go/pkg/auth/internal/utils"
	"its.ac.id/base-go/services/config"
)

const (
	ErrBuildToken = "error_build_token"
)

func Login(ctx *gin.Context, u *contracts.User) error {
	cfg := do.MustInvoke[config.Config](do.DefaultInjector)
	appCfg := cfg.App()
	authCfg := cfg.Auth()
	maxAge := authCfg.MaxAge

	token, err := jwt.NewBuilder().
		Issuer(appCfg.URL).
		Expiration(time.Now().Add(time.Duration(maxAge)*time.Second)).
		Subject(u.Id()).
		IssuedAt(time.Now()).
		Claim("active_role", u.ActiveRole()).
		Claim("roles", u.Roles()).
		Build()

	if err != nil {
		return errors.New(ErrBuildToken)
	}

	signed, err := jwt.Sign(token, jwa.HS256, []byte(appCfg.Key))
	if err != nil {
		return err
	}

	httpCfg := cfg.HTTP()

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		utils.GetCookieName(),
		string(signed),
		maxAge,
		authCfg.CookiePath,
		authCfg.CookieDomain,
		httpCfg.Secure,
		true,
	)

	return nil
}
