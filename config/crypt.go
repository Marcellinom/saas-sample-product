package config

import (
	"os"

	"github.com/dptsi/its-go/crypt"
)

func cryptConfig() crypt.Config {
	return crypt.Config{
		Key: os.Getenv("APP_KEY"),
	}
}
