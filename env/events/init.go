package events

import (
	"github.com/ottemo/foundation/env"
)

func init() {
	instance := new(DefaultEventBus)
	var _ env.I_EventBus = instance

	env.RegisterEventBus(instance)
}
