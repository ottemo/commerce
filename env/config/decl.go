// Package config is a default implementation of I_Config declared in
// "github.com/ottemo/foundation/env" package
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
