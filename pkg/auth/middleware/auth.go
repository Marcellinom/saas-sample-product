package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/samber/do"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/pkg/app/common"
	"its.ac.id/base-go/pkg/auth/internal/utils"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := utils.GetCookieName()
		cookie, err := ctx.Cookie(name)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, common.UnauthorizedResponse)
			ctx.Abort()
			return
		}
		cfg := do.MustInvoke[config.Config](do.DefaultInjector)
		appCfg := cfg.App()

		token, err := jwt.Parse([]byte(cookie), jwt.WithVerify(jwa.HS256, []byte(appCfg.Key)))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, common.UnauthorizedResponse)
			ctx.Abort()
			return
		}
		u := utils.UserFromToken(token)

		ctx.Set(utils.UserKey, u)
		ctx.Next()
	}
}
