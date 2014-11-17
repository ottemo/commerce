package events

import (
	"github.com/ottemo/foundation/env"
)

// init makes package self-initialization routine
func init() {
	instance := new(DefaultEventBus)
	instance.listeners = make(map[string][]env.F_EventListener)

	var _ env.I_EventBus = instance

	env.RegisterEventBus(instance)
}
