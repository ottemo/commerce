package config

import (
	"github.com/ottemo/foundation/env"
)

type DefaultConfig struct {
	configValues     map[string]interface{}
	configValidators map[string]env.F_ConfigValueValidator
}
