package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/pkg/app/common/errors"
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

		msg := fmt.Sprintf("current user active role (%s) doesn't have permission to access this resource", u.ActiveRoleName())
		// details := fmt.Sprintf("allowed role to access this resource are: %s", strings.Join(roles, ", "))
		details := ""
		ctx.Error(errors.NewForbiddenError(msg, details))
		ctx.Abort()
	}
}

func ActiveRoleHasPermission(neededPermission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		u := services.User(ctx)
		if u.HasPermission(neededPermission) {
			ctx.Next()
			return
		}

		msg := fmt.Sprintf("current user active role (%s) doesn't have permission to access this resource", u.ActiveRoleName())
		// details := fmt.Sprintf("permission to access this resource is: %s", neededPermission)
		details := ""
		ctx.Error(errors.NewForbiddenError(msg, details))
		ctx.Abort()
	}
}
