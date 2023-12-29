package config

import (
	"os"

	"bitbucket.org/dptsi/go-framework/web"
)

func webConfig() web.Config {
	return web.Config{
		IsDebugMode: os.Getenv("APP_DEBUG") == "true",
		Environment: os.Getenv("APP_ENV"),
	}
}
