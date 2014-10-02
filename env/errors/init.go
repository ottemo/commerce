package errors

import (
	"github.com/ottemo/foundation/env"
)

func init() {
	instance := &DefaultErrorBus{listeners: make([]env.F_ErrorListener, 0)}
	var _ env.I_ErrorBus = instance

	env.RegisterErrorBus(instance)
}
