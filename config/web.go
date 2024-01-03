package config

import (
	"os"

	"bitbucket.org/dptsi/its-go/web"
)

func webConfig() web.Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return web.Config{
		IsDebugMode: os.Getenv("APP_DEBUG") == "true",
		Environment: os.Getenv("APP_ENV"),
		Port:        port,
	}
}
