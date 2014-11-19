// Package config is a default implementation of InterfaceConfig declared in
// "github.com/ottemo/foundation/env" package
package config

import (
	"github.com/ottemo/foundation/env"
)

// InterfaceConfig implementer class
type DefaultConfig struct {
	configValues     map[string]interface{}
	configTypes      map[string]string
	configValidators map[string]env.FuncConfigValueValidator
}
