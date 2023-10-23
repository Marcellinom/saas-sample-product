package services

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/pkg/auth/contracts"
	internalContract "its.ac.id/base-go/pkg/auth/internal/contracts"
	"its.ac.id/base-go/pkg/session"
)

func Login(ctx *gin.Context, u *contracts.User) error {
	sess := session.Default(ctx)
	userData := internalContract.UserSessionData{
		Id:                strings.ToLower(u.Id()),
		ActiveRole:        u.ActiveRole(),
		Name:              u.Name(),
		PreferredUsername: u.PreferredUsername(),
		Email:             u.Email(),
		Roles:             u.Roles(),
	}
	userJson, err := json.Marshal(userData)
	if err != nil {
		return err
	}
	sess.Set("user", string(userJson))
	sess.Save()

	return nil
}
