package event

import (
	"log"

	"bitbucket.org/dptsi/base-go-libraries/app"
	"github.com/mikestefanello/hooks"
)

var hookEvent = hooks.NewHook[app.Event]("event")

type EventHook struct {
}

func (e *EventHook) Dispatch(ev app.Event) {
	hookEvent.Dispatch(ev)
}

func (e *EventHook) Listen(fn func(ev app.Event)) {
	hookEvent.Listen(func(event hooks.Event[app.Event]) {
		fn(event.Msg)
	})
}

func SetupEventHook() *EventHook {
	hooks.SetLogger(func(format string, args ...any) {
		log.Printf(format+"\n", args...)
	})

	return &EventHook{}
}
