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
		u.SetName(userData.Name)
		u.SetNickname(userData.Nickname)
		u.SetEmail(userData.Email)
		u.SetEmailVerified(bool(userData.EmailVerified))
		u.SetPicture(userData.Picture)
		u.SetGender(userData.Gender)
		u.SetBirthdate(userData.Birthdate)
		u.SetZoneinfo(userData.Zoneinfo)
		u.SetLocale(userData.Locale)
		u.SetPhoneNumber(userData.PhoneNumber)
		u.SetPhoneNumberVerified(bool(userData.PhoneNumberVerified))
		u.SetPreferredUsername(userData.PreferredUsername)
		for _, role := range userData.Roles {
			u.AddRole(role.Id, role.Name, role.Permissions, role.IsDefault)
		}
		u.SetActiveRole(userData.ActiveRole)

		ctx.Set(utils.UserKey, u)
		ctx.Next()
	}
}
