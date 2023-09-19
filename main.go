package main

import (
	"github.com/samber/do"
	"its.ac.id/base-go/pkg/app"

	// Services
	_ "its.ac.id/base-go/services/config"
	routes "its.ac.id/base-go/services/web"

	// Modules
	_ "its.ac.id/base-go/modules/auth"
	_ "its.ac.id/base-go/modules/berkas"
)

func main() {
	i := app.Boot()

	server := do.MustInvoke[routes.Server](i)
	server.Start()
}
