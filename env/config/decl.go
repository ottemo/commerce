// Package "errors" is a default implementation for "I_Config" interface.
package config

import (
	"github.com/ottemo/foundation/env"
)

// I_Config implementer class
type DefaultConfig struct {
	configValues     map[string]interface{}
	configTypes      map[string]string
	configValidators map[string]env.F_ConfigValueValidator
}
