package config

import (
	"fmt"
	"os"

	"github.com/dptsi/its-go/sessions"
	"github.com/stoewer/go-strcase"
)

func sessionsConfig() sessions.Config {
	return sessions.Config{
		Storage:    "database",
		Connection: "default",
		Table:      "sessions",
		Cookie: sessions.CookieConfig{
			Name:           fmt.Sprintf("%s_session", strcase.SnakeCase(os.Getenv("APP_NAME"))),
			CsrfCookieName: "CSRF-TOKEN",
			Path:           "/",
			Domain:         "",
			Secure:         false,
			Lifetime:       60,
		},
	}
}
