package fsmedia

import (
	"github.com/ottemo/foundation/env"
)

// setupConfig setups package configuration values for a system
func (it *FilesystemMediaStorage) setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew("can't obtain config")
	}

	imageSizesValidator := func(newValue interface{}) (interface{}, error) {
		if newValue, ok := newValue.(string); ok && newValue != "" {
			err := it.UpdateSizeNames(newValue)
			if err != nil {
				return ConstDefaultImageSizes, env.ErrorDispatch(err)
			}
			return newValue, nil
		}
		return ConstDefaultImageSizes, env.ErrorNew("unexpected image sizes value")
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMediaImageSizes,
		Value:       ConstDefaultImageSizes,
		Type:        "string",
		Editor:      "line_text",
		Options:     "",
		Label:       "Image sizes",
		Description: "Predefined image sizes in format ([sizeName]: [maxWidth:maxHeight], ...)",
		Image:       "",
	}, imageSizesValidator)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	imageSizesValidator(env.ConfigGetValue(ConstConfigPathMediaImageSizes))

	return nil
}
