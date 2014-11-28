package fsmedia

import (
	"github.com/ottemo/foundation/env"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew("can't obtain config")
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMediaImageSizes,
		Value:       "small: 75x75, thumb: 260x300, big: 560x650",
		Type:        "string",
		Editor:      "line_text",
		Options:     "",
		Label:       "Image sizes",
		Description: "Predefined image sizes in format ([sizeName]: [maxWidth:maxHeight], ...)",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
