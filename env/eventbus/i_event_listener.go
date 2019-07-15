package eventbus

import "github.com/ottemo/commerce/env"

// GetID returns identifier of event listener
func (it *DefaultEventListener) GetID() string {
	return it.id
}

// Handle performs given event handle
func (it *DefaultEventListener) Handle(event env.InterfaceEvent) error {
	return it.handler(event)
}

// GetPriority returns priority of event listener
func (it *DefaultEventListener) GetPriority() float64 {
	return it.priority
}

// SetPriority changes priority of event listener
func (it *DefaultEventListener) SetPriority(value float64) {
	it.priority = value
}
