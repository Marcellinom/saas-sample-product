package controllers

import (
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
}

func NewAuthController(appCfg config.Config, cfg moduleConfig.AuthConfig) *AuthController {
	return &AuthController{appCfg, cfg}
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
	cfg := c.moduleCfg.Oidc()
	endSessionEndpoint, err := op.RPInitiatedLogout(cfg.PostLogoutRedirectURI)
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

// User godoc
// @Summary		get user info
// @Description	get user information
// @Router		/user/:id [get]
// @Tags		auth
// @Param		id path string false "id"
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
	cfg := c.moduleCfg.Oidc()
	sess := session.Default(ctx)
	op, err := oidc.NewClient(
		ctx,
		cfg.Provider,
		cfg.ClientID,
		cfg.ClientSecret,
		cfg.RedirectURL,
		cfg.Scopes,
		sess,
	)
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
	url := op.RedirectURL()
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

	_, IDToken, err := op.ExchangeCodeForToken(ctx, queryParams.Code, queryParams.State)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == oidc.InvalidState {
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

func (c *AuthController) isEntraID() bool {
	return strings.HasPrefix(c.moduleCfg.Oidc().Provider, entraIDPrefix)
}
