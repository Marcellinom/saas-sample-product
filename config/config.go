package config

func Config() map[string]interface{} {
	return map[string]interface{}{
		"cors":       corsConfig(),
		"crypt":      cryptConfig(),
		"csrf":       csrfConfig(),
		"database":   databaseConfig(),
		"logging":    loggingConfig(),
		"middleware": middlewareConfig(),
		"sessions":   sessionsConfig(),
		"sso":        ssoConfig(),
		"web":        webConfig(),
	}
}
