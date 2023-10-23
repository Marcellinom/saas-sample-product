package myitssso

import (
	"github.com/gin-gonic/gin"
	"its.ac.id/base-go/pkg/auth/contracts"
	"its.ac.id/base-go/pkg/oidc"
	"its.ac.id/base-go/pkg/session"
)

type myITSSSOClaim struct {
	Sub               string   `json:"sub"`
	Name              string   `json:"name"`
	Email             string   `json:"email"`
	PreferredUsername string   `json:"preferred_username"`
	Roles             []string `json:"roles"`
}

func GetUserFromAuthorizationCode(ctx *gin.Context, oidcClient *oidc.Client, sess *session.Data, code string, state string) (*contracts.User, error) {
	token, _, err := oidcClient.ExchangeCodeForToken(ctx, sess, code, state)
	if err != nil {
		return nil, err
	}
	// fmt.Println("token", token.AccessToken)
	var claims myITSSSOClaim
	userInfo, err := oidcClient.UserInfo(ctx, token)
	if err != nil {
		return nil, err
	}
	if err := userInfo.Claims(&claims); err != nil {
		return nil, err
	}

	user := contracts.NewUser(claims.Sub)
	user.SetName(claims.Name)
	user.SetPreferredUsername(claims.PreferredUsername)
	user.SetEmail(claims.Email)
	for i, r := range claims.Roles {
		user.AddRole(r, make([]string, 0), i == 0)
	}

	return user, nil
}
