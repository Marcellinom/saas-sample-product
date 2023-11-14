package web

import (
	"context"
	"errors"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"its.ac.id/base-go/bootstrap/config"
	"its.ac.id/base-go/pkg/session"
	"its.ac.id/base-go/pkg/session/adapters"
)

var ErrProjectIDNotConfigured = errors.New("firestore project ID not configured. please set SESSION_FIRESTORE_PROJECT_ID in .env file")
var ErrInvalidSqlServerConfig = errors.New("invalid SQL Server configuration")
var ErrInvalidPostgreSqlConfig = errors.New("invalid PostgreSQL configuration")
var ErrUnknownSessionDriver = errors.New("unknown session driver")

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
	log.Printf("Configured session storage driver: %s\n", cfg.Driver)

	switch cfg.Driver {
	case "firestore":
		return setupFirestoreSessionAdapter(cfg)
	case "sqlite":
		// Contoh penggunaan adapter GORM dengan SQLite
		path := cfg.SQLiteDB
		log.Println("Connecting to SQLite database...")
		db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("SQLite connection error: %w", err)
		}
		log.Println("Successfully connected to SQLite database!")
		return adapters.NewGorm(db), nil
	case "sqlserver":
		// Contoh penggunaan adapter GORM dengan SQL Server
		username := cfg.SQLServerUsername
		password := cfg.SQLServerPassword
		host := cfg.SQLServerHost
		port := cfg.SQLServerPort
		database := cfg.SQLServerDatabase

		if username == "" || password == "" || host == "" || port == "" || database == "" {
			return nil, ErrInvalidSqlServerConfig
		}

		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", username, password, host, port, database)
		log.Println("Connecting to SQL Server database...")
		db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("SQL Server connection error: %w", err)
		}
		log.Println("Successfully connected to SQL Server database!")
		return adapters.NewGorm(db), nil
	case "postgres":
		// Contoh penggunaan adapter GORM dengan PostgreSQL
		username := cfg.PostgreSQLUsername
		password := cfg.PostgreSQLPassword
		host := cfg.PostgreSQLHost
		port := cfg.PostgreSQLPort
		database := cfg.PostgreSQLDatabase

		if username == "" || password == "" || host == "" || port == "" || database == "" {
			return nil, ErrInvalidPostgreSqlConfig
		}

		dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=verify-full", username, password, host, port, database)
		log.Println("Connecting to PostgreSQL database...")
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("PostgreSQL connection error: %w", err)
		}
		log.Println("Successfully connected to PostgreSQL database!")
		return adapters.NewGorm(db), nil
	}

	return nil, ErrUnknownSessionDriver
}
