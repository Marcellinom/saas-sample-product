package middleware

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/pkg/auth/contracts"
	internalContract "its.ac.id/base-go/pkg/auth/internal/contracts"
	"its.ac.id/base-go/pkg/auth/internal/utils"
	"its.ac.id/base-go/pkg/session"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sess := session.Default(ctx)
		userIf, ok := sess.Get("user")
		if !ok {
			ctx.Error(unauthorizedError)
			ctx.Abort()
			return
		}
		userJson, ok := userIf.(string)
		if !ok {
			ctx.Error(unauthorizedError)
			ctx.Abort()
			return
		}
		var userData internalContract.UserSessionData
		err := json.Unmarshal([]byte(userJson), &userData)
		if err != nil {
			ctx.Error(unauthorizedError)
			ctx.Abort()
			return
		}

		u := contracts.NewUser(userData.Id)
		u.SetEmail(userData.Email)
		u.SetName(userData.Name)
		u.SetPreferredUsername(userData.PreferredUsername)
		u.SetPicture(userData.Picture)
		for _, role := range userData.Roles {
			u.AddRole(role.Name, role.Permissions, role.IsDefault)
		}
		u.SetActiveRole(userData.ActiveRole)

		ctx.Set(utils.UserKey, u)
		ctx.Next()
	}
}
