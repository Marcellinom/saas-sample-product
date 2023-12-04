package event

import (
	"log"

	"bitbucket.org/dptsi/go-framework/contracts"
	"github.com/mikestefanello/hooks"
)

var hookEvent = hooks.NewHook[contracts.Event]("event")

type EventHook struct {
}

func (e *EventHook) Dispatch(ev contracts.Event) {
	hookEvent.Dispatch(ev)
}

func (e *EventHook) Listen(fn func(ev contracts.Event)) {
	hookEvent.Listen(func(event hooks.Event[contracts.Event]) {
		fn(event.Msg)
	})
}

func SetupEventHook() *EventHook {
	hooks.SetLogger(func(format string, args ...any) {
		log.Printf(format+"\n", args...)
	})

	return &EventHook{}
}
