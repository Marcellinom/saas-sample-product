package config

import (
	"os"

	"github.com/dptsi/its-go/database"
)

func databaseConfig() database.Config {
	return database.Config{
		Connections: map[string]database.ConnectionConfig{
			"default": {
				Driver:   os.Getenv("DB_DRIVER"),
				User:     os.Getenv("DB_USER"),
				Password: os.Getenv("DB_PASSWORD"),
				Host:     os.Getenv("DB_HOST"),
				Port:     os.Getenv("DB_PORT"),
				Database: os.Getenv("DB_DATABASE"),
			},
		},
	}
}
