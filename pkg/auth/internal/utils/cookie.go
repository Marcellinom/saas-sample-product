package utils

import (
	"strings"

	"github.com/gosimple/slug"
	"github.com/samber/do"
	"its.ac.id/base-go/bootstrap/config"
)

func GetCookieName() string {
	cfg := do.MustInvoke[config.Config](do.DefaultInjector)
	appCfg := cfg.App()
	name := appCfg.Name
	name = strings.ReplaceAll(slug.Make(name), "-", "_")
	name += "_token"

	return name
}
