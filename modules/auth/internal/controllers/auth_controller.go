package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
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
