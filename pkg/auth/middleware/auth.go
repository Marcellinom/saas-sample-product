package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/pkg/app/common"
	"its.ac.id/base-go/pkg/auth/contracts"
	"its.ac.id/base-go/pkg/auth/internal/utils"
	"its.ac.id/base-go/pkg/session"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sess := session.Default(ctx)
		idIf, ok := sess.Get("user.id")
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
			return
		}
		// TODO: Unserialize roles
		// activeRoleIf, ok := sess.Get("user.active_role")
		// if !ok {
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		// 	return
		// }
		// rolesJsonIf, ok := sess.Get("user.roles")
		// if !ok {
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		// 	return
		// }
		id, ok := idIf.(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
			return
		}
		// activeRole, ok := activeRoleIf.(string)
		// if !ok {
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, common.UnauthorizedResponse)
		// 	return
		// }

		u := contracts.NewUser(id)

		ctx.Set(utils.UserKey, u)
		ctx.Next()
	}
}
