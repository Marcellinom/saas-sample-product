package config

import (
	"os"

	"bitbucket.org/dptsi/go-framework/database"
)

var DatabasesConfig map[string]database.Config = map[string]database.Config{
	"default": {
		Driver:   os.Getenv("DB_DRIVER"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
	},
}
