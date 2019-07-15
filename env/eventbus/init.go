package eventbus

import (
	"github.com/ottemo/commerce/api"
	"github.com/ottemo/commerce/app"
	"github.com/ottemo/commerce/env"
)

// init makes package self-initialization
func init() {
	eventBus = new(DefaultEventBus)
	eventBus.events = make(map[string]string)
	eventBus.listeners = make(map[string][]env.InterfaceEventListener)

	app.OnAppInit(func() error {
		for event, listeners := range eventBus.listeners {
			listenersString := ""
			for _, listener := range listeners {
				listenersString += listener.GetID() + "; "
			}

			if _, present := eventBus.events[event]; !present {
				return env.ErrorNew(ConstErrorModule, ConstErrorLevel, "e1d6fb25-e25c-419d-a068-a99e1941605e", "'" + event + "' does not exist, following listeners will not be handled:" + listenersString)
			}
		}
		eventBus.initialized = true
		return nil
	})

	api.RegisterOnRestServiceStart(func() error {
		return setupAPI()
	})

	var _ env.InterfaceEventBus = eventBus

	if err := env.RegisterEventBus(eventBus); err != nil {
		_ = env.ErrorDispatch(err)
	}

}
