package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/pkg/auth/contracts"
	"its.ac.id/base-go/pkg/auth/services"
	"its.ac.id/oidc-client"
)

type AuthController struct {
	i   *do.Injector
	cfg config.Config
}

func NewAuthController() *AuthController {
	i := do.DefaultInjector
	cfg := do.MustInvoke[config.Config](i)

	return &AuthController{i, cfg}
}

func (c *AuthController) Logout(ctx *gin.Context) {
	err := services.Logout(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "logout_failed",
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "logout_success",
		"data":    nil,
	})
}

func (c *AuthController) User(ctx *gin.Context) {
	u := services.User(ctx)
	roles := make([]gin.H, 0)
	for _, r := range u.Roles() {
		roles = append(roles, gin.H{
			"name":        r.Name,
			"permissions": r.Permissions,
			"is_default":  r.IsDefault,
		})
	}
	var activeRole any
	activeRole = nil
	if u.ActiveRole() != "" {
		activeRole = u.ActiveRole()
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "user",
		"data": gin.H{
			"id":          u.Id(),
			"active_role": activeRole,
			"roles":       roles,
		},
	})
}

type OidcCookieProvider struct {
	ctx *gin.Context
	cfg config.Config
}

func (c *OidcCookieProvider) Cookie(name string) (string, error) {
	return c.ctx.Cookie(name)
}

func (c *OidcCookieProvider) SetCookie(name string, value string) {
	cfg := c.cfg.Auth()
	c.ctx.SetCookie(name, value, 0, cfg.CookiePath, cfg.CookieDomain, c.cfg.HTTP().Secure, true)
}

func (c *AuthController) getOidcClient(ctx *gin.Context) (*oidc.Client, error) {
	cp := &OidcCookieProvider{ctx, c.cfg}
	cfg := c.cfg.Oidc()
	op, err := oidc.NewClient(ctx, cfg.Provider, cp, ctx)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func (c *AuthController) Login(ctx *gin.Context) {
	op, err := c.getOidcClient(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "login_failed",
			"data":    nil,
		})
	}
	cfg := c.cfg.Oidc()
	url := op.RedirectURL(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL, cfg.Scopes)
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "login_url",
		"data":    url,
	})
}

func (c *AuthController) Callback(ctx *gin.Context) {
	op, err := c.getOidcClient(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "login_failed",
			"data":    nil,
		})
		return
	}

	userInfo, err := op.UserInfo(c.cfg.Oidc().ClientID, c.cfg.Oidc().ClientSecret, c.cfg.Oidc().RedirectURL, c.cfg.Oidc().Scopes)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == oidc.ErrorRetrieveUserInfo {
			status = http.StatusInternalServerError
		}
		ctx.JSON(status, gin.H{
			"code":    status,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	user := contracts.NewUser(userInfo.Subject)

	err = services.Login(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "login_failed",
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "login_success",
		"data":    nil,
	})
}
