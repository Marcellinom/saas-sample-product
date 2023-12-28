package config

var Config = map[string]interface{}{
	"database": databasesConfig,
	"sessions": sessionsConfig,
	"web":      webConfig,
}
