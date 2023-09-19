package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"its.ac.id/base-go/pkg/auth/contracts"
	"its.ac.id/base-go/pkg/auth/services"
)

type AuthController struct {
	i *do.Injector
}

func NewAuthController(i *do.Injector) *AuthController {
	return &AuthController{i: i}
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
	var roles []gin.H
	for _, r := range u.Roles() {
		roles = append(roles, gin.H{
			"name":        r.Name,
			"permissions": r.Permissions,
			"is_default":  r.IsDefault,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "user",
		"data": gin.H{
			"id":          u.Id(),
			"active_role": u.ActiveRole(),
			"roles":       roles,
		},
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	u := contracts.NewUser("123")
	u.AddRole("admin", []string{"admin"}, true)

	err := services.Login(ctx, u)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "login_failed",
			"data":    nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "logged in",
		"data":    nil,
	})
}
