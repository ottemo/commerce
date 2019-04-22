package eventbus

import (
	"github.com/ottemo/commerce/env"
)

// DefaultEventBus InterfaceEventBus implementer class
type DefaultEventBus struct {
	listeners map[string][]env.FuncEventListener
}
