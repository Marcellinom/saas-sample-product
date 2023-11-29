package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"its.ac.id/base-go/modules/auth/internal/presentation/responses"

	commonErrors "bitbucket.org/dptsi/base-go-libraries/app/errors"
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
		ctx.Error(fmt.Errorf("unable to get login url: %w", err))
		return
	}
	ctx.JSON(http.StatusOK, &responses.GeneralResponse{
		Code:    statusCode[successMessage],
		Message: successMessage,
		Data:    url,
	})
}

func (c *AuthController) Callback(ctx *gin.Context) {
	var queryParams struct {
		Code  string `form:"code" binding:"required"`
		State string `form:"state" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.Error(err)
		return
	}

	sess := session.Default(ctx)
	var user *contracts.User
	var err error
	if c.isEntraID() {
		user, err = entra.GetUserFromAuthorizationCode(ctx, c.oidcClient, sess, queryParams.Code, queryParams.State)
	} else {
		user, err = myitssso.GetUserFromAuthorizationCode(ctx, c.oidcClient, sess, queryParams.Code, queryParams.State)
	}

	isBadRequest := errors.Is(err, oidc.ErrInvalidState) || errors.Is(err, oidc.ErrInvalidNonce) || errors.Is(err, oidc.ErrInvalidCodeChallenge)
	var message string
	if errors.Is(err, oidc.ErrInvalidState) {
		message = invalidState
	} else if errors.Is(err, oidc.ErrInvalidNonce) {
		message = invalidNonce
	} else if errors.Is(err, oidc.ErrInvalidCodeChallenge) {
		message = invalidCodeChallenge
	}
	if isBadRequest {
		data := map[string]interface{}{}
		if c.cfg.App().Env != "production" {
			data["hint"] = "Jika anda menggunakan Postman saat memanggil endpoint /auth/login, maka copy URL dari halaman ini dan buat request ke URL ini melalui Postman. Jika masih gagal, ulangi sekali lagi."
		}
		ctx.Error(commonErrors.NewBadRequest(commonErrors.BadRequestParam{
			Code:    statusCode[message],
			Message: message,
			Data:    data,
		}))
		return
	}

	if err := services.Login(ctx, user); err != nil {
		ctx.Error(err)
		return
	}
	sess.Regenerate()
	if err := sess.Save(); err != nil {
		ctx.Error(err)
		return
	}

	session.AddCookieToResponse(c.cfg.Session(), ctx, sess.Id())

	frontendUrl := c.cfg.App().FrontendURL
	if frontendUrl != "" {
		ctx.Redirect(http.StatusFound, frontendUrl)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    statusCode[successMessage],
		"message": successMessage,
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
			"id":          r.Id,
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
	if u.ActiveRoleName() != "" {
		data["active_role_name"] = u.ActiveRoleName()
	}
	data["roles"] = roles

	ctx.JSON(http.StatusOK, &responses.GeneralResponse{
		Code:    statusCode[successMessage],
		Message: successMessage,
		Data:    data,
	})
}

// @Summary		Rute untuk logout
// @Router		/auth/logout [delete]
// @Tags		Authentication & Authorization
// @Security	Session
// @Security	CSRF Token
// @Produce		json
// @Success		200 {object} responses.GeneralResponse{code=int,message=string,data=string} "Logout berhasil"
func (c *AuthController) Logout(ctx *gin.Context) {
	cfg := c.moduleCfg.Oidc()
	endSessionEndpoint, err := c.oidcClient.RPInitiatedLogout(session.Default(ctx), cfg.PostLogoutRedirectURI)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = services.Logout(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	sess := session.Default(ctx)
	sess.Invalidate()
	sess.RegenerateCSRFToken()
	if err := sess.Save(); err != nil {
		ctx.Error(err)
		return
	}
	session.AddCookieToResponse(c.cfg.Session(), ctx, sess.Id())

	ctx.JSON(http.StatusOK, gin.H{
		"code":    statusCode[successMessage],
		"message": successMessage,
		"data":    endSessionEndpoint,
	})
}

func (c *AuthController) isEntraID() bool {
	return strings.HasPrefix(c.moduleCfg.Oidc().Provider, entraIDPrefix)
}

type switchActiveRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// @Summary		Rute untuk mengubah active role user
// @Router		/auth/user/switch-active-role [post]
// @Tags		Authentication & Authorization
// @Security	Session
// @Security	CSRF Token
// @Accept		json
// @Produce		json
// @Param		body body switchActiveRoleRequest	true	"ID role yang akan dijadikan active role"
// @Success		200 {object} responses.GeneralResponse{code=int,message=string,data=string} "Active role berhasil diubah"
// @Failure		400 {object} responses.GeneralResponse{code=int,message=string,data=string} "Missing role"
// @Failure		400 {object} responses.GeneralResponse{code=int,message=string,data=string} "User tidak memiliki role tersebut"
func (c *AuthController) SwitchActiveRole(ctx *gin.Context) {
	type request switchActiveRoleRequest
	user := services.User(ctx)
	var req request
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	if err := user.SetActiveRole(req.Role); err != nil {
		ctx.Error(commonErrors.NewBadRequest(commonErrors.BadRequestParam{
			Code:    statusCode[userDoesNotHaveThisRole],
			Message: userDoesNotHaveThisRole,
		}))
		return
	}
	if err := services.Login(ctx, user); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, &responses.GeneralResponse{
		Code:    statusCode[successMessage],
		Message: successMessage,
	})
}
