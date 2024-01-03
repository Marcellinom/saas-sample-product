package config

func Config() map[string]interface{} {
	return map[string]interface{}{
		"cors":       corsConfig(),
		"crypt":      cryptConfig(),
		"csrf":       csrfConfig(),
		"database":   databaseConfig(),
		"middleware": middlewareConfig(),
		"sessions":   sessionsConfig(),
		"web":        webConfig(),
	}
}
