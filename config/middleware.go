package config

import "bitbucket.org/dptsi/its-go/http/middleware"

func middlewareConfig() middleware.Config {
	return middleware.Config{
		Groups: map[string][]string{
			"global": {"cors", "start_session", "verify_csrf_token"},
		},
	}
}
