package config

func Config() map[string]interface{} {
	return map[string]interface{}{
		"cors":       corsConfig(),
		"database":   databaseConfig(),
		"middleware": middlewareConfig(),
		"sessions":   sessionsConfig(),
		"web":        webConfig(),
	}
}
