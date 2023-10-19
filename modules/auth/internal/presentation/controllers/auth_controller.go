package controllers

import (
	"errors"
	"net/http"
	"strings"

	"its.ac.id/base-go/modules/auth/internal/presentation/responses"

	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/bootstrap/config"
	moduleConfig "its.ac.id/base-go/modules/auth/internal/app/config"
	"its.ac.id/base-go/pkg/auth/contracts"
	"its.ac.id/base-go/pkg/auth/services"
	"its.ac.id/base-go/pkg/oidc"
	"its.ac.id/base-go/pkg/session"
)

const entraIDPrefix = "https://login.microsoftonline.com"

type AuthController struct {
	cfg       config.Config
	moduleCfg moduleConfig.AuthConfig

	oidcClient *oidc.Client
}

func NewAuthController(appCfg config.Config, cfg moduleConfig.AuthConfig, oidcClient *oidc.Client) *AuthController {
	return &AuthController{appCfg, cfg, oidcClient}
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
	url, err := c.oidcClient.RedirectURL(session.Default(ctx))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, &responses.GeneralResponse{
			Code:    http.StatusInternalServerError,
			Message: "unable_to_get_login_url",
		})
	}
	ctx.JSON(http.StatusOK, &responses.GeneralResponse{
		Code:    http.StatusOK,
		Message: "login_url",
		Data:    url,
	})
}

func (c *AuthController) Callback(ctx *gin.Context) {
	var queryParams struct {
		Code  string `form:"code" binding:"required"`
		State string `form:"state" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "missing_code_or_state",
			"data":    nil,
		})
		return
	}

	sess := session.Default(ctx)
	_, IDToken, err := c.oidcClient.ExchangeCodeForToken(ctx, sess, queryParams.Code, queryParams.State)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, oidc.ErrInvalidState) || errors.Is(err, oidc.ErrInvalidNonce) || errors.Is(err, oidc.ErrInvalidIdToken) {
			status = http.StatusBadRequest
		}
		ctx.JSON(status, gin.H{
			"code":    status,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	userID := IDToken.Subject
	if c.isEntraID() {
		type EntraIDClaim struct {
			ObjectId string `json:"oid"`
		}
		var claims EntraIDClaim
		if err := IDToken.Claims(&claims); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "invalid_id_token",
				"data":    nil,
			})
			return
		}

		userID = claims.ObjectId
	}

	user := contracts.NewUser(userID)
	err = services.Login(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "login_failed",
			"data":    nil,
		})
		return
	}
	sess.Regenerate()
	if err := sess.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "unable_to_save_session",
			"data":    nil,
		})
		return
	}

	session.AddCookieToResponse(c.cfg.Session(), ctx, sess.Id())

	frontendUrl := c.cfg.App().FrontendURL
	if frontendUrl != "" {
		ctx.Redirect(http.StatusFound, frontendUrl)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "login_success",
		"data":    nil,
	})
}

// User godoc
// @Summary		get user info
// @Description	get user information
// @Router		/user [get]
// @Tags		auth
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

func (c *AuthController) Logout(ctx *gin.Context) {
	cfg := c.moduleCfg.Oidc()
	endSessionEndpoint, err := c.oidcClient.RPInitiatedLogout(session.Default(ctx), cfg.PostLogoutRedirectURI)
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

	session.AddCookieToResponse(c.cfg.Session(), ctx, sess.Id())

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "logout_success",
		"data":    endSessionEndpoint,
	})
}

func (c *AuthController) isEntraID() bool {
	return strings.HasPrefix(c.moduleCfg.Oidc().Provider, entraIDPrefix)
}
