package config

import (
	"strings"

	"slices"

	"github.com/gosimple/slug"
	"github.com/joeshaw/envdecode"
)

type SessionConfig struct {
	Lifetime   int    `env:"SESSION_LIFETIME,default=7200"`
	CookieName string `env:"SESSION_NAME,default=base-go"`
	CookiePath string `env:"SESSION_PATH,default=/"`
	Domain     string `env:"SESSION_DOMAIN,default=localhost"`
	Secure     bool   `env:"SESSION_SECURE_COOKIE,default=false"`

	Driver string `env:"SESSION_DRIVER,default="`

	// Firestore session adapter
	FirestoreProjectID  string `env:"SESSION_FIRESTORE_PROJECT_ID"`
	FirestoreCollection string `env:"SESSION_FIRESTORE_COLLECTION,default=sessions"`

	// SQLite session adapter (GORM)
	SQLiteDB string `env:"SESSION_SQLITE_DB"`

	// SQL Server session adapter (GORM)
	SQLServerHost     string `env:"SESSION_SQLSERVER_HOST"`
	SQLServerPort     string `env:"SESSION_SQLSERVER_PORT"`
	SQLServerDatabase string `env:"SESSION_SQLSERVER_DATABASE"`
	SQLServerUsername string `env:"SESSION_SQLSERVER_USERNAME"`
	SQLServerPassword string `env:"SESSION_SQLSERVER_PASSWORD"`

	// PostgreSQL session adapter (GORM)
	PostgreSQLHost     string `env:"SESSION_POSTGRES_HOST"`
	PostgreSQLPort     string `env:"SESSION_POSTGRES_PORT"`
	PostgreSQLDatabase string `env:"SESSION_POSTGRES_DATABASE"`
	PostgreSQLUsername string `env:"SESSION_POSTGRES_USERNAME"`
	PostgreSQLPassword string `env:"SESSION_POSTGRES_PASSWORD"`
}

func (c ConfigImpl) Session() SessionConfig {
	return c.session
}

var availableDrivers = []string{
	"firestore",
	"sqlite",
	"sqlserver",
	"postgres",
}

func setupSessionConfig(appName string) SessionConfig {
	var session SessionConfig
	err := envdecode.StrictDecode(&session)
	if err != nil {
		panic(err)
	}
	if session.CookieName == "" || session.CookieName == "base-go" {
		name := slug.Make(appName)
		name = strings.ReplaceAll(name, "-", "_") + "_session"
		session.CookieName = name
	}
	if session.Driver == "" {
		session.Driver = "sqlite"
	}
	if !slices.Contains(availableDrivers, session.Driver) {
		panic("invalid session driver")
	}

	return session
}
