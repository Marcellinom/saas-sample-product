package event

import (
	"log"

	"github.com/mikestefanello/hooks"
	"its.ac.id/base-go/pkg/app/common"
)

var hookEvent = hooks.NewHook[common.Event]("event")

type EventHook struct {
}

func (e *EventHook) Dispatch(ev common.Event) {
	hookEvent.Dispatch(ev)
}

func (e *EventHook) Listen(fn func(ev common.Event)) {
	hookEvent.Listen(func(event hooks.Event[common.Event]) {
		fn(event.Msg)
	})
}

func SetupEventHook() *EventHook {
	hooks.SetLogger(func(format string, args ...any) {
		log.Printf(format+"\n", args...)
	})

	return &EventHook{}
}
