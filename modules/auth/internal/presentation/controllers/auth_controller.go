package controllers

import (
	"its.ac.id/base-go/modules/auth/internal/presentation/responses"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"its.ac.id/base-go/modules/auth/internal/app/config"
	"its.ac.id/base-go/pkg/auth/contracts"
	"its.ac.id/base-go/pkg/auth/services"
	"its.ac.id/base-go/pkg/oidc"
	"its.ac.id/base-go/pkg/session"
)

type AuthController struct {
	i   *do.Injector
	cfg config.AuthConfig
}

func NewAuthController() *AuthController {
	i := do.DefaultInjector
	cfg := do.MustInvoke[config.AuthConfig](i)

	return &AuthController{i, cfg}
}

func (c *AuthController) Logout(ctx *gin.Context) {
	op, err := c.getOidcClient(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "unable_to_get_oidc_client",
			"data":    nil,
		})
		return
	}
	cfg := c.cfg.Oidc()
	endSessionEndpoint, err := op.RPInitiatedLogout(cfg.EndSessionEndpoint, cfg.PostLogoutRedirectURI)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "unable_to_get_end_session_endpoint",
			"data":    nil,
		})
		return
	}

	err = services.Logout(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "logout_failed",
			"data":    nil,
		})
		return
	}
	sess := session.Default(ctx)
	sess.Invalidate()
	sess.RegenerateCSRFToken()
	if err := sess.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "unable_to_save_session",
			"data":    nil,
		})
		return
	}

	session.AddCookieToResponse(ctx, sess.Id())

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "logout_success",
		"data":    endSessionEndpoint,
	})
}

// User godoc
// @Summary		get user info
// @Description	get user information
// @Router		/user [get]
// @Tags		auth
// @Param		username query string true "username"
// @Param		password query string true "password"
// @Produce		json
// @Success		200 {object} responses.GeneralResponse{code=int,message=string,data=responses.User} "Success"
// @Failure		500 {object} responses.GeneralResponse{code=int,message=string} "Internal Server Error"
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

	ctx.JSON(http.StatusOK, &responses.GeneralResponse{
		Code:    http.StatusOK,
		Message: "user",
		Data: gin.H{
			"id":          u.Id(),
			"active_role": activeRole,
			"roles":       roles,
		},
	})
}

func (c *AuthController) getOidcClient(ctx *gin.Context) (*oidc.Client, error) {
	cfg := c.cfg.Oidc()
	sess := session.Default(ctx)
	op, err := oidc.NewClient(ctx, cfg.Provider, sess, ctx)
	if err != nil {
		return nil, err
	}

	return op, nil
}

// Login godoc
// @Summary		login user
// @Description	call oidc login function
// @Router		/login [post]
// @Tags		auth
// @Accept		json
// @Produce		json
// @Success		200 {object} responses.GeneralResponse{code=int,message=string,data=string} "Success"
// @Failure		500 {object} responses.GeneralResponse{code=int,message=string} "Internal Server Error"
func (c *AuthController) Login(ctx *gin.Context) {
	op, err := c.getOidcClient(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &responses.GeneralResponse{
			Code:    http.StatusInternalServerError,
			Message: "login_failed",
		})
	}
	cfg := c.cfg.Oidc()
	url := op.RedirectURL(cfg.ClientID, cfg.ClientSecret, cfg.RedirectURL, cfg.Scopes)
	ctx.JSON(http.StatusOK, &responses.GeneralResponse{
		Code:    http.StatusOK,
		Message: "login_url",
		Data:    url,
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
	sess := session.Default(ctx)
	sess.Regenerate()
	if err := sess.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "unable_to_save_session",
			"data":    nil,
		})
		return
	}

	session.AddCookieToResponse(ctx, sess.Id())
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "login_success",
		"data":    nil,
	})
}
