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
	"its.ac.id/base-go/pkg/entra"
	"its.ac.id/base-go/pkg/myitssso"
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

// @Summary		Rute untuk mendapatkan link login melalui OpenID Connect
// @Router		/auth/login [post]
// @Tags		Authentication & Authorization
// @Produce		json
// @Security 	CSRF Token
// @Success		200 {object} responses.GeneralResponse "Link login berhasil didapatkan"
// @Failure		500 {object} responses.GeneralResponse "Terjadi kesalahan saat menghubungi provider OpenID Connect"
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

// @Summary		Rute untuk logout
// @Router		/auth/logout [delete]
// @Tags		Authentication & Authorization
// @Security	Session
// @Produce		json
// @Success		200 {object} responses.GeneralResponse{code=int,message=string,data=string} "Logout berhasil"
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
	var user *contracts.User
	if c.isEntraID() {
		tmp, err := entra.GetUserFromAuthorizationCode(ctx, c.oidcClient, sess, queryParams.Code, queryParams.State)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "unable_to_get_user",
				"data":    nil,
			})
			return
		}
		user = tmp
	} else {
		tmp, err := myitssso.GetUserFromAuthorizationCode(ctx, c.oidcClient, sess, queryParams.Code, queryParams.State)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "unable_to_get_user",
				"data":    nil,
			})
			return
		}
		user = tmp
	}

	if err := services.Login(ctx, user); err != nil {
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

// @Summary		Rute untuk mendapatkan data user yang sedang login
// @Router		/auth/user [get]
// @Tags		Authentication & Authorization
// @Security	Session
// @Produce		json
// @Success		200 {object} responses.GeneralResponse{code=int,message=string,data=responses.User{roles=[]responses.Role}} "Data user berhasil didapatkan"
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

	data := make(map[string]interface{})
	data["id"] = u.Id()
	data["name"] = nil
	data["email"] = nil
	data["preferred_username"] = nil
	data["picture"] = nil
	data["active_role"] = nil
	if u.Name() != "" {
		data["name"] = u.Name()
	}
	if u.Email() != "" {
		data["email"] = u.Email()
	}
	if u.PreferredUsername() != "" {
		data["preferred_username"] = u.PreferredUsername()
	}
	if u.Picture() != "" {
		data["picture"] = u.Picture()
	}
	if u.ActiveRole() != "" {
		data["active_role"] = u.ActiveRole()
	}
	data["roles"] = roles

	ctx.JSON(http.StatusOK, &responses.GeneralResponse{
		Code:    http.StatusOK,
		Message: "user",
		Data:    data,
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

// @Summary		Rute untuk mengubah active role user
// @Router		/auth/user/switch-active-role [post]
// @Tags		Authentication & Authorization
// @Security	Session
// @Produce		json
// @Param		role	body	string	true	"Nama role yang akan dijadikan active role"
// @Success		200 {object} responses.GeneralResponse{code=int,message=string,data=string} "Active role berhasil diubah"
// @Failure		400 {object} responses.GeneralResponse{code=int,message=string,data=string} "Missing role"
// @Failure		400 {object} responses.GeneralResponse{code=int,message=string,data=string} "User tidak memiliki role tersebut"
func (c *AuthController) SwitchActiveRole(ctx *gin.Context) {
	type request struct {
		Role string `json:"role" binding:"required"`
	}
	user := services.User(ctx)
	var req request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &responses.GeneralResponse{
			Code:    http.StatusBadRequest,
			Message: "missing_role",
		})
		return
	}
	if err := user.SetActiveRole(req.Role); err != nil {
		ctx.JSON(http.StatusBadRequest, &responses.GeneralResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if err := services.Login(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, &responses.GeneralResponse{
			Code:    http.StatusInternalServerError,
			Message: "unable_to_change_active_role",
		})
		return
	}

	ctx.JSON(http.StatusOK, &responses.GeneralResponse{
		Code:    http.StatusOK,
		Message: "active_role_changed",
	})
}
