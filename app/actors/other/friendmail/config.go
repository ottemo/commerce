package friendmail

import (
	"github.com/ottemo/foundation/env"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	if config := env.GetConfig(); config != nil {
		err := config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathGroup,
			Value:       nil,
			Type:        env.ConstConfigTypeGroup,
			Editor:      "",
			Options:     nil,
			Label:       "Stock",
			Description: "Stock management system",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathEmailSubject,
			Value: 		 "Email friend",
			Type:        env.ConstConfigTypeVarchar,
			Editor:      "text",
			Options:     nil,
			Label:       "Email subject",
			Description: "Email subject for the friend form",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}

		config.RegisterItem(env.StructConfigItem{
			Path:        ConstConfigPathEmailTemplate,
			Value: `Dear {{.name}}
<br />
<br />
Your friend sent you an email:
{{.content}}`,
			Type:        env.ConstConfigTypeText,
			Editor:      "multiline_text",
			Options:     nil,
			Label:       "Email Body",
			Description: "Email body template for the friend form",
			Image:       "",
		}, nil)

		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}
