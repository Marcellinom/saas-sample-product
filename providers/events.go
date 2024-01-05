package providers

import (
	"github.com/dptsi/its-go/contracts"
)

type Listener struct {
	eventName            string
	listenersConstructor []func(application contracts.Application) (contracts.EventListener, error)
}

var listen []Listener = []Listener{}

func registerEvents(application contracts.Application) {
	service := application.Services().Event
	for _, l := range listen {
		service.Register(l.eventName, l.listenersConstructor)
	}
}
