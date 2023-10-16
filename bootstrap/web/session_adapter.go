package web

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/pkg/session"
	"its.ac.id/base-go/pkg/session/adapters"
)

var ErrProjectIDNotConfigured = errors.New("firestore project ID not configured. please set SESSION_FIRESTORE_PROJECT_ID in .env file")

func setupFirestoreSessionAdapter(cfg config.SessionConfig) (*adapters.Firestore, error) {
	ctx := context.Background()
	if cfg.FirestoreProjectID == "" {
		return nil, ErrProjectIDNotConfigured
	}

	client, err := firestore.NewClient(ctx, cfg.FirestoreProjectID)
	if err != nil {
		return nil, err
	}
	return adapters.NewFirestore(client, cfg.FirestoreCollection), nil
}

func setupSessionStorage(cfg config.SessionConfig) (session.Storage, error) {
	switch cfg.Driver {
	case "firestore":
		return setupFirestoreSessionAdapter(cfg)
	case "sqlite":
		// Contoh penggunaan adapter GORM dengan SQLite
		path := cfg.SQLiteDB
		if path == "" {
			panic("invalid SQLite database path for session")
		}

		db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			panic("failed to connect to SQLite database for session")
		}
		return adapters.NewGorm(db), nil
	case "sqlserver":
		// Contoh penggunaan adapter GORM dengan SQL Server
		username := cfg.SQLServerUsername
		password := cfg.SQLServerPassword
		host := cfg.SQLServerHost
		port := cfg.SQLServerPort
		database := cfg.SQLServerDatabase

		if username == "" || password == "" || host == "" || port == "" || database == "" {
			panic("invalid SQL Server configuration for session")
		}

		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", username, password, host, port, database)
		db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect SQL Server database for session")
		}
		return adapters.NewGorm(db), nil
	}

	panic("unknown session driver")
}
