package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/samber/do"
	"its.ac.id/base-go/pkg/app"

	// Services
	_ "its.ac.id/base-go/services/config"
	routes "its.ac.id/base-go/services/web"

	// Modules
	_ "its.ac.id/base-go/modules/auth"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
	i := app.Boot()

	server := do.MustInvoke[routes.Server](i)
	server.Start()
}
