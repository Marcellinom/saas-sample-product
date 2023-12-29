package config

import "bitbucket.org/dptsi/go-framework/sessions"

func sessionsConfig() sessions.Config {
	return sessions.Config{
		Storage:    "database",
		Connection: "default",
		Table:      "sessions",
		Cookie: sessions.CookieConfig{
			Name:           "myits_academics_session",
			CsrfCookieName: "CSRF-TOKEN",
			Path:           "/",
			Domain:         "",
			Secure:         false,
			Lifetime:       60,
		},
	}
}
