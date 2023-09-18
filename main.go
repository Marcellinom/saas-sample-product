package main

import (
	"github.com/samber/do"
	"its.ac.id/base-go/pkg/app"

	// Services
	routes "its.ac.id/base-go/services/routes"
)

func main() {
	i := app.Boot()

	server := do.MustInvoke[routes.Server](i)
	server.Start()
}
