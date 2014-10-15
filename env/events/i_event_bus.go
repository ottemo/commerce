package events

import (
	"github.com/ottemo/foundation/env"
)

// adds listener to event handling stack
func (it *DefaultEventBus) RegisterListener(listener env.F_EventListener) {
	it.listeners = append(it.listeners, listener)
}

// generates new event, with following dispatching
func (it *DefaultEventBus) New(event string, args map[string]interface{}) {
	for _, listener := range it.listeners {
		if listener(event, args) {
			break
		}
	}
}
