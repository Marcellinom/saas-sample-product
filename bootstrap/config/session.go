package config

import (
	"bitbucket.org/dptsi/base-go-libraries/contracts"
	"bitbucket.org/dptsi/base-go-libraries/database"
	"bitbucket.org/dptsi/base-go-libraries/sessions/adapters"
	"github.com/samber/do"
)

func SetupSession(i *do.Injector) {
	dbMgr := do.MustInvoke[*database.Manager](i)
	do.Provide[contracts.SessionStorage](i, func(i *do.Injector) (contracts.SessionStorage, error) {
		db := dbMgr.GetDefault()
		return adapters.NewGorm(db.DB), nil
	})
}
