package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/pkg/auth/services"
)

func ActiveRoleIn(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		u := services.User(ctx)
		for _, role := range roles {
			if role == u.ActiveRole() {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "forbidden",
			"data":    nil,
		})
	}
}
