package providers

import "github.com/dptsi/its-go/contracts"

type Middleware struct {
	name        string
	constructor contracts.MiddlewareConstructor
}

var middlewares []Middleware = []Middleware{}

func registerMiddlewares(application contracts.Application) {
	service := application.Services().Middleware
	for _, m := range middlewares {
		service.Register(m.name, m.constructor)
	}
}
