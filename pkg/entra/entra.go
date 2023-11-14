package entra

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/pkg/auth/contracts"
	"its.ac.id/base-go/pkg/oidc"
	"its.ac.id/base-go/pkg/session"
)

type entraIDClaim struct {
	ObjectId          string   `json:"oid"`
	Name              string   `json:"name"`
	Email             string   `json:"email"`
	PreferredUsername string   `json:"preferred_username"`
	Roles             []string `json:"roles"`
}

func GetUserFromAuthorizationCode(ctx *gin.Context, oidcClient *oidc.Client, sess *session.Data, code string, state string) (*contracts.User, error) {
	_, IDToken, err := oidcClient.ExchangeCodeForToken(ctx, sess, code, state)
	if err != nil {
		return nil, fmt.Errorf("get user from entra id failed: %w", err)
	}

	var claims entraIDClaim
	if err := IDToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("get user from entra id failed: %w", err)
	}

	user := contracts.NewUser(claims.ObjectId)
	user.SetName(claims.Name)
	user.SetPreferredUsername(claims.PreferredUsername)
	user.SetEmail(claims.Email)
	for i, r := range claims.Roles {
		user.AddRole(r, make([]string, 0), i == 0)
	}

	return user, nil
}
