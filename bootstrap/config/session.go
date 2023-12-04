package config

import (
	"bitbucket.org/dptsi/go-framework/contracts"
	"bitbucket.org/dptsi/go-framework/database"
	"bitbucket.org/dptsi/go-framework/sessions/adapters"
	"github.com/samber/do"
)

func SetupSession(i *do.Injector) {
	dbMgr := do.MustInvoke[*database.Manager](i)
	do.Provide[contracts.SessionStorage](i, func(i *do.Injector) (contracts.SessionStorage, error) {
		db := dbMgr.GetDefault()
		return adapters.NewGorm(db.DB), nil
	})
}
