package eventbus

import (
	"github.com/ottemo/commerce/env"
)

// RegisterEvent informs the system on available event
func (it *DefaultEventBus) RegisterEvent(event string, description string) error {

	if description, present := it.events[event]; present {
		return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "93384d10-c5a9-43f4-a96a-c5dab83a3869", "'" + event + "' event was already registered: " + description )
	} else {
		it.events[event] = description

		lastChar := len(event) - 1
		for charIdx, char := range event {
			if charIdx == 0 || charIdx == lastChar || char == '.' {
				levelEvent := event
				if charIdx != lastChar {
					levelEvent = event[0:charIdx]
				}

				if value, present := it.events[levelEvent]; present {
					if value[0:1] == "+" {
						value += ", " + event
					}
				} else {
					it.events[levelEvent] = "+ " + event
				}
			}
		}
	}
	return nil
}

// RegisterListener adds listener to event handling stack
//   - event listening is patch based, "" - global listener on any event, "api.product" - will listen for app events starts with api.product.[...])
func (it *DefaultEventBus) RegisterListener(event string, id string, handler env.FuncEventHandler) error {

	if it.initialized {
		if _, present := it.events[event]; !present {
			env.ErrorNew(ConstErrorModule, ConstErrorLevel, "fd5e78ed-ff4e-42cf-93d8-a4124f86c514", "there is no event " + event)
		}
	}

	listener := new(DefaultEventListener)
	listener.id = id
	listener.handler = handler
	listener.priority = 1

	if value, present := it.listeners[event]; present {
		it.listeners[event] = append(value, listener)
	} else {
		it.listeners[event] = []env.InterfaceEventListener{listener}
	}

	return nil
}

// Handle performs event handling cross registered listeners
func (it *DefaultEventBus) Handle(event string, data map[string]interface{}) error {

	if data == nil {
		data = make(map[string]interface{})
	}

	eventObject := new(DefaultEvent)
	eventObject.id = event
	eventObject.data = data
	eventObject.propagate = false

	// loop over top level events
	// (i.e. "api.checkout.success" event will notify following listeners: "", "api", "api.checkout", "api.checkout.success")
	lastChar := len(event) - 1
	for charIdx, char := range event {
		if charIdx == 0 || charIdx == lastChar || char == '.' {
			levelEvent := event
			if charIdx != lastChar {
				levelEvent = event[0:charIdx]
			}

			// processing listeners withing level if present
			if listeners, present := it.listeners[levelEvent]; present {
				for _, listener := range listeners {

					// processing listener, if it wants to stop handling - doing this
					if err := listener.Handle(eventObject); err != nil {
						env.ErrorDispatch(err)
					}

				}
			}
		}
	}
	return nil
}

// GetListeners returns a list of registered listeners
func (it *DefaultEventBus) GetListeners(event string) []env.InterfaceEventListener {
	if value, present := it.listeners[event]; present {
		return value
	}
	return []env.InterfaceEventListener {}
}
