package config

import (
	"github.com/joeshaw/envdecode"
)

type OidcConfig struct {
	Provider              string   `env:"OIDC_PROVIDER,required"`
	ClientID              string   `env:"OIDC_CLIENT_ID,required"`
	ClientSecret          string   `env:"OIDC_CLIENT_SECRET,required"`
	RedirectURL           string   `env:"OIDC_REDIRECT_URL,required"`
	Scopes                []string `env:"OIDC_SCOPES,default=openid,email,profile,groups"`
	PostLogoutRedirectURI string   `env:"OIDC_POST_LOGOUT_REDIRECT_URI"`
}

type AuthConfig interface {
	Oidc() OidcConfig
}

type AuthConfigImpl struct {
	oidc OidcConfig
}

func (c AuthConfigImpl) Oidc() OidcConfig {
	return c.oidc
}

func SetupConfig() (AuthConfig, error) {
	var oidc OidcConfig
	err := envdecode.StrictDecode(&oidc)
	if err != nil {
		return nil, err
	}

	return AuthConfigImpl{oidc}, nil
}
