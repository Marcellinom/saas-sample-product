package config

import (
	"os"

	"bitbucket.org/dptsi/base-go-libraries/database"
	"github.com/samber/do"
)

func SetupDatabase(i *do.Injector) {
	dbMgr := do.MustInvoke[*database.Manager](i)
	err := dbMgr.AddDatabase(database.DefaultDatabaseName, database.Config{
		Driver:   os.Getenv("DB_DRIVER"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_DATABASE"),
	})

	if err != nil {
		panic(err)
	}
}
